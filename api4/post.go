// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"

	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/app"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

func InitPost() {
	l4g.Debug(utils.T("api.post.init.debug"))

	BaseRoutes.Posts.Handle("", ApiSessionRequired(createPost)).Methods("POST")
	BaseRoutes.Post.Handle("", ApiSessionRequired(getPost)).Methods("GET")
	BaseRoutes.Post.Handle("/thread", ApiSessionRequired(getPostThread)).Methods("GET")
	BaseRoutes.PostsForChannel.Handle("", ApiSessionRequired(getPostsForChannel)).Methods("GET")
}

func createPost(c *Context, w http.ResponseWriter, r *http.Request) {
	post := model.PostFromJson(r.Body)
	if post == nil {
		c.SetInvalidParam("post")
		return
	}

	post.UserId = c.Session.UserId

	if !app.SessionHasPermissionToChannel(c.Session, post.ChannelId, model.PERMISSION_CREATE_POST) {
		c.SetPermissionError(model.PERMISSION_CREATE_POST)
		return
	}

	if post.CreateAt != 0 && !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		post.CreateAt = 0
	}

	rp, err := app.CreatePostAsUser(post)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(rp.ToJson()))
}

func getPostsForChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	if !app.SessionHasPermissionToChannel(c.Session, c.Params.ChannelId, model.PERMISSION_READ_CHANNEL) {
		c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
		return
	}

	etag := app.GetPostsEtag(c.Params.ChannelId)

	if HandleEtag(etag, "Get Posts", w, r) {
		return
	}

	if list, err := app.GetPostsPage(c.Params.ChannelId, c.Params.Page, c.Params.PerPage); err != nil {
		c.Err = err
		return
	} else {
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
		w.Write([]byte(list.ToJson()))
	}
}

func getPost(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePostId()
	if c.Err != nil {
		return
	}

	if !app.SessionHasPermissionToChannelByPost(c.Session, c.Params.PostId, model.PERMISSION_READ_CHANNEL) {
		c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
		return
	}

	if post, err := app.GetSinglePost(c.Params.PostId); err != nil {
		c.Err = err
		return
	} else if HandleEtag(post.Etag(), "Get Post", w, r) {
		return
	} else {
		w.Header().Set(model.HEADER_ETAG_SERVER, post.Etag())
		w.Write([]byte(post.ToJson()))
	}
}

func getPostThread(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePostId()
	if c.Err != nil {
		return
	}

	if !app.SessionHasPermissionToChannelByPost(c.Session, c.Params.PostId, model.PERMISSION_READ_CHANNEL) {
		c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
		return
	}

	if list, err := app.GetPostThread(c.Params.PostId); err != nil {
		c.Err = err
		return
	} else if HandleEtag(list.Etag(), "Get Post Thread", w, r) {
		return
	} else {
		w.Header().Set(model.HEADER_ETAG_SERVER, list.Etag())
		w.Write([]byte(list.ToJson()))
	}
}
