// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"github.com/mattermost/platform/model"
	"strings"
	"testing"
)

func TestShortcutsCommand(t *testing.T) {
	th := Setup().InitBasic()
	Client := th.BasicClient
	channel := th.BasicChannel

	rs := Client.Must(Client.Command(channel.Id, "/shortcuts ", false)).Data.(*model.CommandResponse)
	if !strings.Contains(rs.Text, "ALT") {
		t.Fatal("failed to display shortcuts")
	}
}
