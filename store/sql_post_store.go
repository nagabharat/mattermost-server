// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

type SqlPostStore struct {
	*SqlStore
}

func NewSqlPostStore(sqlStore *SqlStore) PostStore {
	s := &SqlPostStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Post{}, "Posts").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("UserId").SetMaxSize(26)
		table.ColMap("ChannelId").SetMaxSize(26)
		table.ColMap("RootId").SetMaxSize(26)
		table.ColMap("ParentId").SetMaxSize(26)
		table.ColMap("Message").SetMaxSize(4000)
		table.ColMap("Type").SetMaxSize(26)
		table.ColMap("Hashtags").SetMaxSize(1000)
		table.ColMap("Props").SetMaxSize(4000)
		table.ColMap("Filenames").SetMaxSize(4000)
	}

	return s
}

func (s SqlPostStore) UpgradeSchemaIfNeeded() {
}

func (s SqlPostStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_posts_update_at", "Posts", "UpdateAt")
	s.CreateIndexIfNotExists("idx_posts_create_at", "Posts", "CreateAt")
	s.CreateIndexIfNotExists("idx_posts_channel_id", "Posts", "ChannelId")
	s.CreateIndexIfNotExists("idx_posts_root_id", "Posts", "RootId")

	s.CreateFullTextIndexIfNotExists("idx_posts_message_txt", "Posts", "Message")
	s.CreateFullTextIndexIfNotExists("idx_posts_hashtags_txt", "Posts", "Hashtags")
}

func (s SqlPostStore) Save(post *model.Post) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if len(post.Id) > 0 {
			result.Err = model.NewAppError("SqlPostStore.Save",
				"You cannot update an existing Post", "id="+post.Id)
			storeChannel <- result
			close(storeChannel)
			return
		}

		post.PreSave()
		if result.Err = post.IsValid(); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if err := s.GetMaster().Insert(post); err != nil {
			result.Err = model.NewAppError("SqlPostStore.Save", "We couldn't save the Post", "id="+post.Id+", "+err.Error())
		} else {
			time := model.GetMillis()

			if post.Type != model.POST_JOIN_LEAVE {
				s.GetMaster().Exec("UPDATE Channels SET LastPostAt = :LastPostAt, TotalMsgCount = TotalMsgCount + 1 WHERE Id = :ChannelId", map[string]interface{}{"LastPostAt": time, "ChannelId": post.ChannelId})
			} else {
				// don't update TotalMsgCount for unimportant messages so that the channel isn't marked as unread
				s.GetMaster().Exec("UPDATE Channels SET LastPostAt = :LastPostAt WHERE Id = :ChannelId", map[string]interface{}{"LastPostAt": time, "ChannelId": post.ChannelId})
			}

			if len(post.RootId) > 0 {
				s.GetMaster().Exec("UPDATE Posts SET UpdateAt = :UpdateAt WHERE Id = :RootId", map[string]interface{}{"UpdateAt": time, "RootId": post.RootId})
			}

			result.Data = post
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) Update(oldPost *model.Post, newMessage string, newHashtags string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		editPost := *oldPost
		editPost.Message = newMessage
		editPost.UpdateAt = model.GetMillis()
		editPost.Hashtags = newHashtags

		oldPost.DeleteAt = editPost.UpdateAt
		oldPost.UpdateAt = editPost.UpdateAt
		oldPost.OriginalId = oldPost.Id
		oldPost.Id = model.NewId()

		if result.Err = editPost.IsValid(); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if _, err := s.GetMaster().Update(&editPost); err != nil {
			result.Err = model.NewAppError("SqlPostStore.Update", "We couldn't update the Post", "id="+editPost.Id+", "+err.Error())
		} else {
			time := model.GetMillis()
			s.GetMaster().Exec("UPDATE Channels SET LastPostAt = :LastPostAt  WHERE Id = :ChannelId", map[string]interface{}{"LastPostAt": time, "ChannelId": editPost.ChannelId})

			if len(editPost.RootId) > 0 {
				s.GetMaster().Exec("UPDATE Posts SET UpdateAt = :UpdateAt WHERE Id = :RootId", map[string]interface{}{"UpdateAt": time, "RootId": editPost.RootId})
			}

			// mark the old post as deleted
			s.GetMaster().Insert(oldPost)

			result.Data = &editPost
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) Get(id string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}
		pl := &model.PostList{}

		var post model.Post
		err := s.GetReplica().SelectOne(&post, "SELECT * FROM Posts WHERE Id = :Id AND DeleteAt = 0", map[string]interface{}{"Id": id})
		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.GetPost", "We couldn't get the post", "id="+id+err.Error())
		}

		pl.AddPost(&post)
		pl.AddOrder(id)

		rootId := post.RootId

		if rootId == "" {
			rootId = post.Id
		}

		var posts []*model.Post
		_, err = s.GetReplica().Select(&posts, "SELECT * FROM Posts WHERE (Id = :Id OR RootId = :RootId) AND DeleteAt = 0", map[string]interface{}{"Id": rootId, "RootId": rootId})
		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.GetPost", "We couldn't get the post", "root_id="+rootId+err.Error())
		} else {
			for _, p := range posts {
				pl.AddPost(p)
			}
		}

		result.Data = pl

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

type etagPosts struct {
	Id       string
	UpdateAt int64
}

func (s SqlPostStore) GetEtag(channelId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var et etagPosts
		err := s.GetReplica().SelectOne(&et, "SELECT Id, UpdateAt FROM Posts WHERE ChannelId = :ChannelId ORDER BY UpdateAt DESC LIMIT 1", map[string]interface{}{"ChannelId": channelId})
		if err != nil {
			result.Data = fmt.Sprintf("%v.0.%v", model.CurrentVersion, model.GetMillis())
		} else {
			result.Data = fmt.Sprintf("%v.%v.%v", model.CurrentVersion, et.Id, et.UpdateAt)
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) Delete(postId string, time int64) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		_, err := s.GetMaster().Exec("Update Posts SET DeleteAt = :DeleteAt, UpdateAt = :UpdateAt WHERE Id = :Id OR ParentId = :ParentId OR RootId = :RootId", map[string]interface{}{"DeleteAt": time, "UpdateAt": time, "Id": postId, "ParentId": postId, "RootId": postId})
		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.Delete", "We couldn't delete the post", "id="+postId+", err="+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) GetPosts(channelId string, offset int, limit int) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if limit > 1000 {
			result.Err = model.NewAppError("SqlPostStore.GetLinearPosts", "Limit exceeded for paging", "channelId="+channelId)
			storeChannel <- result
			close(storeChannel)
			return
		}

		rpc := s.getRootPosts(channelId, offset, limit)
		cpc := s.getParentsPosts(channelId, offset, limit)

		if rpr := <-rpc; rpr.Err != nil {
			result.Err = rpr.Err
		} else if cpr := <-cpc; cpr.Err != nil {
			result.Err = cpr.Err
		} else {
			posts := rpr.Data.([]*model.Post)
			parents := cpr.Data.([]*model.Post)

			list := &model.PostList{Order: make([]string, 0, len(posts))}

			for _, p := range posts {
				list.AddPost(p)
				list.AddOrder(p.Id)
			}

			for _, p := range parents {
				list.AddPost(p)
			}

			list.MakeNonNil()

			result.Data = list
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) GetPostsSince(channelId string, time int64) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var posts []*model.Post
		_, err := s.GetReplica().Select(&posts,
			`(SELECT
			    *
			FROM
			    Posts
			WHERE
			    (UpdateAt > :Time
			        AND ChannelId = :ChannelId)
			LIMIT 1000)
			UNION
			(SELECT
			    *
			FROM
			    Posts
			WHERE
			    Id
			IN
			    (SELECT * FROM (SELECT
			        RootId
			    FROM
			        Posts
			    WHERE
			        UpdateAt > :Time
			            AND ChannelId = :ChannelId
			    LIMIT 1000) temp_tab))
			ORDER BY CreateAt DESC`,
			map[string]interface{}{"ChannelId": channelId, "Time": time})

		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.GetPostsSince", "We couldn't get the posts for the channel", "channelId="+channelId+err.Error())
		} else {

			list := &model.PostList{Order: make([]string, 0, len(posts))}

			for _, p := range posts {
				list.AddPost(p)
				if p.UpdateAt > time {
					list.AddOrder(p.Id)
				}
			}

			result.Data = list
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) getRootPosts(channelId string, offset int, limit int) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var posts []*model.Post
		_, err := s.GetReplica().Select(&posts, "SELECT * FROM Posts WHERE ChannelId = :ChannelId AND DeleteAt = 0 ORDER BY CreateAt DESC LIMIT :Limit OFFSET :Offset", map[string]interface{}{"ChannelId": channelId, "Offset": offset, "Limit": limit})
		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.GetLinearPosts", "We couldn't get the posts for the channel", "channelId="+channelId+err.Error())
		} else {
			result.Data = posts
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) getParentsPosts(channelId string, offset int, limit int) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var posts []*model.Post
		_, err := s.GetReplica().Select(&posts,
			`SELECT
			    q2.*
			FROM
			    Posts q2
			        INNER JOIN
			    (SELECT DISTINCT
			        q3.RootId
			    FROM
			        (SELECT
			        RootId
			    FROM
			        Posts
			    WHERE
			        ChannelId = :ChannelId1
			            AND DeleteAt = 0
			    ORDER BY CreateAt DESC
			    LIMIT :Limit OFFSET :Offset) q3
			    WHERE q3.RootId != '') q1 ON q1.RootId = q2.Id
			WHERE
			    ChannelId = :ChannelId2
			        AND DeleteAt = 0
			ORDER BY CreateAt`,
			map[string]interface{}{"ChannelId1": channelId, "Offset": offset, "Limit": limit, "ChannelId2": channelId})
		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.GetLinearPosts", "We couldn't get the parent post for the channel", "channelId="+channelId+err.Error())
		} else {
			result.Data = posts
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

var specialSearchChar = []string{
	"<",
	">",
	"+",
	"-",
	"(",
	")",
	"~",
	"@",
}

func (s SqlPostStore) Search(teamId string, userId string, params *model.SearchParams) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		termMap := map[string]bool{}
		terms := params.Terms

		if terms == "" && params.InChannel == "" && params.FromUser == "" {
			result.Data = []*model.Post{}
			storeChannel <- result
			return
		}

		searchType := "Message"
		if params.IsHashtag {
			searchType = "Hashtags"
			for _, term := range strings.Split(terms, " ") {
				termMap[term] = true
			}
		}

		// these chars have speical meaning and can be treated as spaces
		for _, c := range specialSearchChar {
			terms = strings.Replace(terms, c, " ", -1)
		}

		var posts []*model.Post

		if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_POSTGRES {
			// Parse text for wildcards
			if wildcard, err := regexp.Compile("\\*($| )"); err == nil {
				terms = wildcard.ReplaceAllLiteralString(terms, "* ")
			}
		}

		searchQuery := `
			SELECT
				*
			FROM
				Posts
			WHERE
				DeleteAt = 0
				POST_FILTER
				AND ChannelId IN (
					SELECT
						Id
					FROM
						Channels,
						ChannelMembers
					WHERE
						Id = ChannelId
							AND TeamId = :TeamId
							AND UserId = :UserId
							AND DeleteAt = 0
							CHANNEL_FILTER)
				SEARCH_CLAUSE
				ORDER BY CreateAt DESC
			LIMIT 100`

		if params.InChannel != "" {
			searchQuery = strings.Replace(searchQuery, "CHANNEL_FILTER", "AND Name = :InChannel", 1)
		} else {
			searchQuery = strings.Replace(searchQuery, "CHANNEL_FILTER", "", 1)
		}

		if params.FromUser != "" {
			searchQuery = strings.Replace(searchQuery, "POST_FILTER", `
				AND UserId IN (
					SELECT
						Id
					FROM
						Users
					WHERE
						TeamId = :TeamId
						AND Username = :FromUser)`, 1)
		} else {
			searchQuery = strings.Replace(searchQuery, "POST_FILTER", "", 1)
		}

		if terms == "" {
			// we've already confirmed that we have a channel or user to search for
			searchQuery = strings.Replace(searchQuery, "SEARCH_CLAUSE", "", 1)
		} else if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_POSTGRES {
			// Parse text for wildcards
			if wildcard, err := regexp.Compile("\\*($| )"); err == nil {
				terms = wildcard.ReplaceAllLiteralString(terms, ":* ")
			}

			terms = strings.Join(strings.Fields(terms), " | ")

			searchClause := fmt.Sprintf("AND %s @@  to_tsquery(:Terms)", searchType)
			searchQuery = strings.Replace(searchQuery, "SEARCH_CLAUSE", searchClause, 1)
		} else if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
			searchClause := fmt.Sprintf("AND MATCH (%s) AGAINST (:Terms IN BOOLEAN MODE)", searchType)
			searchQuery = strings.Replace(searchQuery, "SEARCH_CLAUSE", searchClause, 1)
		}

		queryParams := map[string]interface{}{
			"TeamId":    teamId,
			"UserId":    userId,
			"Terms":     terms,
			"InChannel": params.InChannel,
			"FromUser":  params.FromUser,
		}

		_, err := s.GetReplica().Select(&posts, searchQuery, queryParams)
		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.Search", "We encounted an error while searching for posts", "teamId="+teamId+", err="+err.Error())
		}

		list := &model.PostList{Order: make([]string, 0, len(posts))}

		for _, p := range posts {
			if searchType == "Hashtags" {
				exactMatch := false
				for _, tag := range strings.Split(p.Hashtags, " ") {
					if termMap[tag] {
						exactMatch = true
					}
				}
				if !exactMatch {
					continue
				}
			}
			list.AddPost(p)
			list.AddOrder(p.Id)
		}

		list.MakeNonNil()

		result.Data = list

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlPostStore) GetForExport(channelId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var posts []*model.Post
		_, err := s.GetReplica().Select(
			&posts,
			"SELECT * FROM Posts WHERE ChannelId = :ChannelId AND DeleteAt = 0",
			map[string]interface{}{"ChannelId": channelId})
		if err != nil {
			result.Err = model.NewAppError("SqlPostStore.GetForExport", "We couldn't get the posts for the channel", "channelId="+channelId+err.Error())
		} else {
			result.Data = posts
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
