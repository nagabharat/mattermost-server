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

func InitSystem() {
	l4g.Debug(utils.T("api.system.init.debug"))

	BaseRoutes.System.Handle("/ping", ApiHandler(getSystemPing)).Methods("GET")
	BaseRoutes.ApiRoot.Handle("/config", ApiSessionRequired(getConfig)).Methods("GET")
	BaseRoutes.ApiRoot.Handle("/email/test", ApiSessionRequired(testEmail)).Methods("POST")
}

func getSystemPing(c *Context, w http.ResponseWriter, r *http.Request) {
	ReturnStatusOK(w)
}

func testEmail(c *Context, w http.ResponseWriter, r *http.Request) {

	if !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	err := app.TestEmail(c.Session.UserId, utils.Cfg)
	if err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}

func getConfig(c *Context, w http.ResponseWriter, r *http.Request) {
	if !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	cfg := app.GetConfig()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Write([]byte(cfg.ToJson()))
}
