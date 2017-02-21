// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/mattermost/platform/app"
	"github.com/mattermost/platform/model"
)

func TestCreatePost(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	post := &model.Post{ChannelId: th.BasicChannel.Id, Message: "#hashtag a" + model.NewId() + "a"}
	rpost, resp := Client.CreatePost(post)
	CheckNoError(t, resp)

	if rpost.Message != post.Message {
		t.Fatal("message didn't match")
	}

	if rpost.Hashtags != "#hashtag" {
		t.Fatal("hashtag didn't match")
	}

	if len(rpost.FileIds) != 0 {
		t.Fatal("shouldn't have files")
	}

	if rpost.EditAt != 0 {
		t.Fatal("newly created post shouldn't have EditAt set")
	}

	post.RootId = rpost.Id
	post.ParentId = rpost.Id
	_, resp = Client.CreatePost(post)
	CheckNoError(t, resp)

	post.RootId = "junk"
	_, resp = Client.CreatePost(post)
	CheckBadRequestStatus(t, resp)

	post.RootId = rpost.Id
	post.ParentId = "junk"
	_, resp = Client.CreatePost(post)
	CheckBadRequestStatus(t, resp)

	post2 := &model.Post{ChannelId: th.BasicChannel2.Id, Message: "a" + model.NewId() + "a", CreateAt: 123}
	rpost2, resp := Client.CreatePost(post2)

	if rpost2.CreateAt == post2.CreateAt {
		t.Fatal("create at should not match")
	}

	post.RootId = rpost2.Id
	post.ParentId = rpost2.Id
	_, resp = Client.CreatePost(post)
	CheckBadRequestStatus(t, resp)

	post.RootId = ""
	post.ParentId = ""
	post.ChannelId = "junk"
	_, resp = Client.CreatePost(post)
	CheckForbiddenStatus(t, resp)

	post.ChannelId = model.NewId()
	_, resp = Client.CreatePost(post)
	CheckForbiddenStatus(t, resp)

	if r, err := Client.DoApiPost("/posts", "garbage"); err == nil {
		t.Fatal("should have errored")
	} else {
		if r.StatusCode != http.StatusBadRequest {
			t.Log("actual: " + strconv.Itoa(r.StatusCode))
			t.Log("expected: " + strconv.Itoa(http.StatusBadRequest))
			t.Fatal("wrong status code")
		}
	}

	Client.Logout()
	_, resp = Client.CreatePost(post)
	CheckUnauthorizedStatus(t, resp)

	post.ChannelId = th.BasicChannel.Id
	post.CreateAt = 123
	rpost, resp = th.SystemAdminClient.CreatePost(post)
	CheckNoError(t, resp)

	if rpost.CreateAt != post.CreateAt {
		t.Fatal("create at should match")
	}
}

func TestGetPostsForChannel(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	post1 := th.CreatePost()
	post2 := th.CreatePost()
	post3 := th.CreatePost()
	post4 := &model.Post{ChannelId: th.BasicChannel.Id, Message: "a" + model.NewId() + "a", RootId: post1.Id}
	post4, _ = Client.CreatePost(post4)

	posts, resp := Client.GetPostsForChannel(th.BasicChannel.Id, 0, 60, "")
	CheckNoError(t, resp)

	if posts.Order[0] != post4.Id {
		t.Fatal("wrong order")
	}

	if posts.Order[1] != post3.Id {
		t.Fatal("wrong order")
	}

	if posts.Order[2] != post2.Id {
		t.Fatal("wrong order")
	}

	if posts.Order[3] != post1.Id {
		t.Fatal("wrong order")
	}

	posts, resp = Client.GetPostsForChannel(th.BasicChannel.Id, 0, 3, resp.Etag)
	CheckEtag(t, posts, resp)

	posts, resp = Client.GetPostsForChannel(th.BasicChannel.Id, 0, 3, "")
	CheckNoError(t, resp)

	if len(posts.Order) != 3 {
		t.Fatal("wrong number returned")
	}

	if _, ok := posts.Posts[post4.Id]; !ok {
		t.Fatal("missing comment")
	}

	if _, ok := posts.Posts[post1.Id]; !ok {
		t.Fatal("missing root post")
	}

	posts, resp = Client.GetPostsForChannel(th.BasicChannel.Id, 1, 1, "")
	CheckNoError(t, resp)

	if posts.Order[0] != post3.Id {
		t.Fatal("wrong order")
	}

	posts, resp = Client.GetPostsForChannel(th.BasicChannel.Id, 10000, 10000, "")
	CheckNoError(t, resp)

	if len(posts.Order) != 0 {
		t.Fatal("should be no posts")
	}

	_, resp = Client.GetPostsForChannel("", 0, 60, "")
	CheckUnauthorizedStatus(t, resp)

	_, resp = Client.GetPostsForChannel("junk", 0, 60, "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetPostsForChannel(model.NewId(), 0, 60, "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetPostsForChannel(model.NewId(), 0, 60, "")
	CheckUnauthorizedStatus(t, resp)

	_, resp = th.SystemAdminClient.GetPostsForChannel(th.BasicChannel.Id, 0, 60, "")
	CheckNoError(t, resp)
}

func TestGetPost(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	post, resp := Client.GetPost(th.BasicPost.Id, "")
	CheckNoError(t, resp)

	if post.Id != th.BasicPost.Id {
		t.Fatal("post ids don't match")
	}

	post, resp = Client.GetPost(th.BasicPost.Id, resp.Etag)
	CheckEtag(t, post, resp)

	_, resp = Client.GetPost("", "")
	CheckNotFoundStatus(t, resp)

	_, resp = Client.GetPost("junk", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetPost(model.NewId(), "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetPost(model.NewId(), "")
	CheckUnauthorizedStatus(t, resp)

	post, resp = th.SystemAdminClient.GetPost(th.BasicPost.Id, "")
	CheckNoError(t, resp)
}

func TestDeletePost(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	_, resp := Client.DeletePost("")
	CheckNotFoundStatus(t, resp)

	_, resp = Client.DeletePost("junk")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.DeletePost(th.BasicPost.Id)
	CheckForbiddenStatus(t, resp)

	Client.Login(th.TeamAdminUser.Email, th.TeamAdminUser.Password)
	_, resp = Client.DeletePost(th.BasicPost.Id)
	CheckNoError(t, resp)

	post := th.CreatePost()
	user := th.CreateUser()

	Client.Logout()
	Client.Login(user.Email, user.Password)

	_, resp = Client.DeletePost(post.Id)
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.DeletePost(model.NewId())
	CheckUnauthorizedStatus(t, resp)

	status, resp := th.SystemAdminClient.DeletePost(post.Id)
	if status == false {
		t.Fatal("post should return status OK")
	}
	CheckNoError(t, resp)
}

func TestGetPostThread(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	post := &model.Post{ChannelId: th.BasicChannel.Id, Message: "a" + model.NewId() + "a", RootId: th.BasicPost.Id}
	post, _ = Client.CreatePost(post)

	list, resp := Client.GetPostThread(th.BasicPost.Id, "")
	CheckNoError(t, resp)

	var list2 *model.PostList
	list2, resp = Client.GetPostThread(th.BasicPost.Id, resp.Etag)
	CheckEtag(t, list2, resp)

	if list.Order[0] != th.BasicPost.Id {
		t.Fatal("wrong order")
	}

	if _, ok := list.Posts[th.BasicPost.Id]; !ok {
		t.Fatal("should have had post")
	}

	if _, ok := list.Posts[post.Id]; !ok {
		t.Fatal("should have had post")
	}

	_, resp = Client.GetPostThread("junk", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetPostThread(model.NewId(), "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetPostThread(model.NewId(), "")
	CheckUnauthorizedStatus(t, resp)

	list, resp = th.SystemAdminClient.GetPostThread(th.BasicPost.Id, "")
	CheckNoError(t, resp)
}

func TestSearchPosts(t *testing.T) {
	th := Setup().InitBasic()
	defer TearDown()
	th.LoginBasic()
	Client := th.Client

	message := "search for post1"
	_ = th.CreateMessagePost(message)

	message = "search for post2"
	post2 := th.CreateMessagePost(message)

	message = "#hashtag search for post3"
	post3 := th.CreateMessagePost(message)

	message = "hashtag for post4"
	_ = th.CreateMessagePost(message)

	posts, resp := Client.SearchPosts(th.BasicTeam.Id, "search", false)
	CheckNoError(t, resp)
	if len(posts.Order) != 3 {
		t.Fatal("wrong search")
	}

	posts, resp = Client.SearchPosts(th.BasicTeam.Id, "post2", false)
	CheckNoError(t, resp)
	if len(posts.Order) != 1 && posts.Order[0] == post2.Id {
		t.Fatal("wrong search")
	}

	posts, resp = Client.SearchPosts(th.BasicTeam.Id, "#hashtag", false)
	CheckNoError(t, resp)
	if len(posts.Order) != 1 && posts.Order[0] == post3.Id {
		t.Fatal("wrong search")
	}

	if posts, resp = Client.SearchPosts(th.BasicTeam.Id, "*", false); len(posts.Order) != 0 {
		t.Fatal("searching for just * shouldn't return any results")
	}

	posts, resp = Client.SearchPosts(th.BasicTeam.Id, "post1 post2", true)
	CheckNoError(t, resp)
	if len(posts.Order) != 2 {
		t.Fatal("wrong search results")
	}

	_, resp = Client.SearchPosts("junk", "#sgtitlereview", false)
	CheckBadRequestStatus(t, resp)

	_, resp = Client.SearchPosts(model.NewId(), "#sgtitlereview", false)
	CheckForbiddenStatus(t, resp)

	_, resp = Client.SearchPosts(th.BasicTeam.Id, "", false)
	CheckBadRequestStatus(t, resp)

	Client.Logout()
	_, resp = Client.SearchPosts(th.BasicTeam.Id, "#sgtitlereview", false)
	CheckUnauthorizedStatus(t, resp)

}

func TestSearchHashtagPosts(t *testing.T) {
	th := Setup().InitBasic()
	defer TearDown()
	th.LoginBasic()
	Client := th.Client

	message := "#sgtitlereview with space"
	_ = th.CreateMessagePost(message)

	message = "#sgtitlereview\n with return"
	_ = th.CreateMessagePost(message)

	message = "no hashtag"
	_ = th.CreateMessagePost(message)

	posts, resp := Client.SearchPosts(th.BasicTeam.Id, "#sgtitlereview", false)
	CheckNoError(t, resp)
	if len(posts.Order) != 2 {
		t.Fatal("wrong search results")
	}

	Client.Logout()
	_, resp = Client.SearchPosts(th.BasicTeam.Id, "#sgtitlereview", false)
	CheckUnauthorizedStatus(t, resp)
}

func TestSearchPostsInChannel(t *testing.T) {
	th := Setup().InitBasic()
	defer TearDown()
	th.LoginBasic()
	Client := th.Client

	channel := th.CreatePublicChannel()

	message := "sgtitlereview with space"
	_ = th.CreateMessagePost(message)

	message = "sgtitlereview\n with return"
	_ = th.CreateMessagePostWithClient(Client, th.BasicChannel2, message)

	message = "other message with no return"
	_ = th.CreateMessagePostWithClient(Client, th.BasicChannel2, message)

	message = "other message with no return"
	_ = th.CreateMessagePostWithClient(Client, channel, message)

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "channel:", false); len(posts.Order) != 0 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "in:", false); len(posts.Order) != 0 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "channel:"+th.BasicChannel.Name, false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "in:"+th.BasicChannel2.Name, false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "channel:"+th.BasicChannel2.Name, false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "ChAnNeL:"+th.BasicChannel2.Name, false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "sgtitlereview", false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "sgtitlereview channel:"+th.BasicChannel.Name, false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "sgtitlereview in: "+th.BasicChannel2.Name, false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "sgtitlereview channel: "+th.BasicChannel2.Name, false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "channel: "+th.BasicChannel2.Name+" channel: "+channel.Name, false); len(posts.Order) != 3 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

}

func TestSearchPostsFromUser(t *testing.T) {
	th := Setup().InitBasic()
	defer TearDown()
	Client := th.Client

	th.LoginTeamAdmin()
	user := th.CreateUser()
	LinkUserToTeam(user, th.BasicTeam)
	app.AddUserToChannel(user, th.BasicChannel)
	app.AddUserToChannel(user, th.BasicChannel2)

	message := "sgtitlereview with space"
	_ = th.CreateMessagePost(message)

	Client.Logout()
	th.LoginBasic2()

	message = "sgtitlereview\n with return"
	_ = th.CreateMessagePostWithClient(Client, th.BasicChannel2, message)

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.TeamAdminUser.Username, false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.BasicUser2.Username, false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.BasicUser2.Username+" sgtitlereview", false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	message = "hullo"
	_ = th.CreateMessagePost(message)

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.BasicUser2.Username+" in:"+th.BasicChannel.Name, false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	Client.Login(user.Email, user.Password)

	// wait for the join/leave messages to be created for user3 since they're done asynchronously
	time.Sleep(100 * time.Millisecond)

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.BasicUser2.Username, false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.BasicUser2.Username+" from: "+user.Username, false); len(posts.Order) != 2 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.BasicUser2.Username+" from: "+user.Username+" in:"+th.BasicChannel2.Name, false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}

	message = "coconut"
	_ = th.CreateMessagePostWithClient(Client, th.BasicChannel2, message)

	if posts, _ := Client.SearchPosts(th.BasicTeam.Id, "from: "+th.BasicUser2.Username+" from: "+user.Username+" in:"+th.BasicChannel2.Name+" coconut", false); len(posts.Order) != 1 {
		t.Fatalf("wrong number of posts returned %v", len(posts.Order))
	}
}
