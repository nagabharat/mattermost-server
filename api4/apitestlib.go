// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/app"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/store"
	"github.com/mattermost/platform/utils"
)

type TestHelper struct {
	Client        *model.Client4
	BasicUser     *model.User
	BasicUser2    *model.User
	TeamAdminUser *model.User
	BasicTeam     *model.Team

	SystemAdminClient *model.Client4
	SystemAdminUser   *model.User
}

func Setup() *TestHelper {
	if app.Srv == nil {
		utils.TranslationsPreInit()
		utils.LoadConfig("config.json")
		utils.InitTranslations(utils.Cfg.LocalizationSettings)
		utils.Cfg.TeamSettings.MaxUsersPerTeam = 50
		*utils.Cfg.RateLimitSettings.Enable = false
		utils.Cfg.EmailSettings.SendEmailNotifications = true
		utils.Cfg.EmailSettings.SMTPServer = "dockerhost"
		utils.Cfg.EmailSettings.SMTPPort = "2500"
		utils.Cfg.EmailSettings.FeedbackEmail = "test@example.com"
		utils.DisableDebugLogForTest()
		app.NewServer()
		app.InitStores()
		InitRouter()
		app.StartServer()
		InitApi(true)
		utils.EnableDebugLogForTest()
		app.Srv.Store.MarkSystemRanUnitTests()

		*utils.Cfg.TeamSettings.EnableOpenServer = true
	}

	th := &TestHelper{}
	th.Client = th.CreateClient()
	th.SystemAdminClient = th.CreateClient()
	return th
}

func (me *TestHelper) InitBasic() *TestHelper {
	me.TeamAdminUser = me.CreateUser()
	me.LoginTeamAdmin()
	me.BasicTeam = me.CreateTeam()
	me.BasicUser = me.CreateUser()
	LinkUserToTeam(me.BasicUser, me.BasicTeam)
	me.BasicUser2 = me.CreateUser()
	LinkUserToTeam(me.BasicUser2, me.BasicTeam)
	app.UpdateUserRoles(me.BasicUser.Id, model.ROLE_SYSTEM_USER.Id)
	me.LoginBasic()

	return me
}

func (me *TestHelper) InitSystemAdmin() *TestHelper {
	me.SystemAdminUser = me.CreateUser()
	app.UpdateUserRoles(me.SystemAdminUser.Id, model.ROLE_SYSTEM_USER.Id+" "+model.ROLE_SYSTEM_ADMIN.Id)
	me.LoginSystemAdmin()

	return me
}

func (me *TestHelper) CreateClient() *model.Client4 {
	return model.NewAPIv4Client("http://localhost" + utils.Cfg.ServiceSettings.ListenAddress)
}

func (me *TestHelper) CreateUser() *model.User {
	return me.CreateUserWithClient(me.Client)
}

func (me *TestHelper) CreateTeam() *model.Team {
	return me.CreateTeamWithClient(me.Client)
}

func (me *TestHelper) CreateTeamWithClient(client *model.Client4) *model.Team {
	id := model.NewId()
	team := &model.Team{
		DisplayName: "dn_" + id,
		Name:        GenerateTestTeamName(),
		Email:       GenerateTestEmail(),
		Type:        model.TEAM_OPEN,
	}

	utils.DisableDebugLogForTest()
	rteam, _ := client.CreateTeam(team)
	utils.EnableDebugLogForTest()
	return rteam
}

func (me *TestHelper) CreateUserWithClient(client *model.Client4) *model.User {
	id := model.NewId()

	user := &model.User{
		Email:     GenerateTestEmail(),
		Username:  GenerateTestUsername(),
		Nickname:  "nn_" + id,
		FirstName: "f_" + id,
		LastName:  "l_" + id,
		Password:  "Password1",
	}

	utils.DisableDebugLogForTest()
	ruser, _ := client.CreateUser(user)
	ruser.Password = "Password1"
	VerifyUserEmail(ruser.Id)
	utils.EnableDebugLogForTest()
	return ruser
}

func (me *TestHelper) LoginBasic() {
	me.LoginBasicWithClient(me.Client)
}

func (me *TestHelper) LoginBasic2() {
	me.LoginBasic2WithClient(me.Client)
}

func (me *TestHelper) LoginTeamAdmin() {
	me.LoginTeamAdminWithClient(me.Client)
}

func (me *TestHelper) LoginSystemAdmin() {
	me.LoginSystemAdminWithClient(me.SystemAdminClient)
}

func (me *TestHelper) LoginBasicWithClient(client *model.Client4) {
	utils.DisableDebugLogForTest()
	client.Login(me.BasicUser.Email, me.BasicUser.Password)
	utils.EnableDebugLogForTest()
}

func (me *TestHelper) LoginBasic2WithClient(client *model.Client4) {
	utils.DisableDebugLogForTest()
	client.Login(me.BasicUser.Email, me.BasicUser.Password)
	utils.EnableDebugLogForTest()
}

func (me *TestHelper) LoginTeamAdminWithClient(client *model.Client4) {
	utils.DisableDebugLogForTest()
	client.Login(me.TeamAdminUser.Email, me.TeamAdminUser.Password)
	utils.EnableDebugLogForTest()
}

func (me *TestHelper) LoginSystemAdminWithClient(client *model.Client4) {
	utils.DisableDebugLogForTest()
	client.Login(me.SystemAdminUser.Email, me.SystemAdminUser.Password)
	utils.EnableDebugLogForTest()
}

func LinkUserToTeam(user *model.User, team *model.Team) {
	utils.DisableDebugLogForTest()

	err := app.JoinUserToTeam(team, user)
	if err != nil {
		l4g.Error(err.Error())
		l4g.Close()
		time.Sleep(time.Second)
		panic(err)
	}

	utils.EnableDebugLogForTest()
}

func GenerateTestEmail() string {
	return strings.ToLower("success+" + model.NewId() + "@simulator.amazonses.com")
}

func GenerateTestUsername() string {
	return "fakeuser" + model.NewId()
}

func GenerateTestTeamName() string {
	return "faketeam" + model.NewId()
}

func VerifyUserEmail(userId string) {
	store.Must(app.Srv.Store.User().VerifyEmail(userId))
}

func CheckUserSanitization(t *testing.T, user *model.User) {
	if user.Password != "" {
		t.Fatal("password wasn't blank")
	}

	if user.AuthData != nil && *user.AuthData != "" {
		t.Fatal("auth data wasn't blank")
	}

	if user.MfaSecret != "" {
		t.Fatal("mfa secret wasn't blank")
	}
}

func CheckEtag(t *testing.T, data interface{}, resp *model.Response) {
	if !reflect.ValueOf(data).IsNil() {
		t.Fatal("etag data was not nil")
	}

	if resp.StatusCode != http.StatusNotModified {
		t.Log("actual: " + strconv.Itoa(resp.StatusCode))
		t.Log("expected: " + strconv.Itoa(http.StatusNotModified))
		t.Fatal("wrong status code for etag")
	}
}

func CheckNoError(t *testing.T, resp *model.Response) {
	if resp.Error != nil {
		t.Fatal(resp.Error)
	}
}

func CheckForbiddenStatus(t *testing.T, resp *model.Response) {
	if resp.Error == nil {
		t.Fatal("should have errored with status:" + strconv.Itoa(http.StatusForbidden))
		return
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Log("actual: " + strconv.Itoa(resp.StatusCode))
		t.Log("expected: " + strconv.Itoa(http.StatusForbidden))
		t.Fatal("wrong status code")
	}
}

func CheckUnauthorizedStatus(t *testing.T, resp *model.Response) {
	if resp.Error == nil {
		t.Fatal("should have errored with status:" + strconv.Itoa(http.StatusUnauthorized))
		return
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Log("actual: " + strconv.Itoa(resp.StatusCode))
		t.Log("expected: " + strconv.Itoa(http.StatusUnauthorized))
		t.Fatal("wrong status code")
	}
}

func CheckNotFoundStatus(t *testing.T, resp *model.Response) {
	if resp.Error == nil {
		t.Fatal("should have errored with status:" + strconv.Itoa(http.StatusNotFound))
		return
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Log("actual: " + strconv.Itoa(resp.StatusCode))
		t.Log("expected: " + strconv.Itoa(http.StatusNotFound))
		t.Fatal("wrong status code")
	}
}

func CheckBadRequestStatus(t *testing.T, resp *model.Response) {
	if resp.Error == nil {
		t.Fatal("should have errored with status:" + strconv.Itoa(http.StatusBadRequest))
		return
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Log("actual: " + strconv.Itoa(resp.StatusCode))
		t.Log("expected: " + strconv.Itoa(http.StatusBadRequest))
		t.Fatal("wrong status code")
	}
}

func CheckErrorMessage(t *testing.T, resp *model.Response, errorId string) {
	if resp.Error == nil {
		t.Fatal("should have errored with message:" + errorId)
		return
	}

	if resp.Error.Id != errorId {
		t.Log("actual: " + resp.Error.Id)
		t.Log("expected: " + errorId)
		t.Fatal("incorrect error message")
	}
}
