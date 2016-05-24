// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"fmt"
	"net/http"

	l4g "github.com/alecthomas/log4go"

	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

func InitGeneral() {
	l4g.Debug(utils.T("api.general.init.debug"))

	BaseRoutes.General.Handle("/client_props", ApiAppHandler(getClientConfig)).Methods("GET")
	BaseRoutes.General.Handle("/log_client", ApiAppHandler(logClient)).Methods("POST")
	BaseRoutes.General.Handle("/ping", ApiAppHandler(ping)).Methods("GET")

}

func getClientConfig(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(model.MapToJson(utils.ClientCfg)))
}

func logClient(c *Context, w http.ResponseWriter, r *http.Request) {
	m := model.MapFromJson(r.Body)

	lvl := m["level"]
	msg := m["message"]

	if len(msg) > 400 {
		msg = msg[0:399]
	}

	if lvl == "ERROR" {
		err := &model.AppError{}
		err.Message = msg
		err.Id = msg
		err.Where = "client"
		c.LogError(err)
	}

	ReturnStatusOK(w)
}

func ping(c *Context, w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["version"] = model.CurrentVersion
	m["server_time"] = fmt.Sprintf("%v", model.GetMillis())
	m["node_id"] = ""
	w.Write([]byte(model.MapToJson(m)))
}
