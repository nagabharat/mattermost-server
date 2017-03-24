// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

func TestCreateTeam(t *testing.T) {
	th := Setup().InitBasic()
	defer TearDown()
	Client := th.Client

	team := &model.Team{Name: GenerateTestUsername(), DisplayName: "Some Team", Type: model.TEAM_OPEN}
	rteam, resp := Client.CreateTeam(team)
	CheckNoError(t, resp)

	if rteam.Name != team.Name {
		t.Fatal("names did not match")
	}

	if rteam.DisplayName != team.DisplayName {
		t.Fatal("display names did not match")
	}

	if rteam.Type != team.Type {
		t.Fatal("types did not match")
	}

	_, resp = Client.CreateTeam(rteam)
	CheckBadRequestStatus(t, resp)

	rteam.Id = ""
	_, resp = Client.CreateTeam(rteam)
	CheckErrorMessage(t, resp, "store.sql_team.save.domain_exists.app_error")
	CheckBadRequestStatus(t, resp)

	rteam.Name = ""
	_, resp = Client.CreateTeam(rteam)
	CheckErrorMessage(t, resp, "model.team.is_valid.characters.app_error")
	CheckBadRequestStatus(t, resp)

	if r, err := Client.DoApiPost("/teams", "garbage"); err == nil {
		t.Fatal("should have errored")
	} else {
		if r.StatusCode != http.StatusBadRequest {
			t.Log("actual: " + strconv.Itoa(r.StatusCode))
			t.Log("expected: " + strconv.Itoa(http.StatusBadRequest))
			t.Fatal("wrong status code")
		}
	}

	Client.Logout()

	_, resp = Client.CreateTeam(rteam)
	CheckUnauthorizedStatus(t, resp)

	// Update permission
	enableTeamCreation := utils.Cfg.TeamSettings.EnableTeamCreation
	defer func() {
		utils.Cfg.TeamSettings.EnableTeamCreation = enableTeamCreation
		utils.SetDefaultRolesBasedOnConfig()
	}()
	utils.Cfg.TeamSettings.EnableTeamCreation = false
	utils.SetDefaultRolesBasedOnConfig()

	th.LoginBasic()
	_, resp = Client.CreateTeam(team)
	CheckForbiddenStatus(t, resp)
}

func TestGetTeam(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	team := th.BasicTeam

	rteam, resp := Client.GetTeam(team.Id, "")
	CheckNoError(t, resp)

	if rteam.Id != team.Id {
		t.Fatal("wrong team")
	}

	_, resp = Client.GetTeam("junk", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeam("", "")
	CheckNotFoundStatus(t, resp)

	_, resp = Client.GetTeam(model.NewId(), "")
	CheckNotFoundStatus(t, resp)

	th.LoginTeamAdmin()

	team2 := &model.Team{DisplayName: "Name", Name: GenerateTestTeamName(), Email: GenerateTestEmail(), Type: model.TEAM_INVITE}
	rteam2, _ := Client.CreateTeam(team2)

	th.LoginBasic()
	_, resp = Client.GetTeam(rteam2.Id, "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetTeam(team.Id, "")
	CheckUnauthorizedStatus(t, resp)

	_, resp = th.SystemAdminClient.GetTeam(rteam2.Id, "")
	CheckNoError(t, resp)
}

func TestGetTeamUnread(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	teamUnread, resp := Client.GetTeamUnread(th.BasicTeam.Id, th.BasicUser.Id)
	CheckNoError(t, resp)
	if teamUnread.TeamId != th.BasicTeam.Id {
		t.Fatal("wrong team id returned for regular user call")
	}

	_, resp = Client.GetTeamUnread("junk", th.BasicUser.Id)
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamUnread(th.BasicTeam.Id, "junk")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamUnread(model.NewId(), th.BasicUser.Id)
	CheckForbiddenStatus(t, resp)

	_, resp = Client.GetTeamUnread(th.BasicTeam.Id, model.NewId())
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetTeamUnread(th.BasicTeam.Id, th.BasicUser.Id)
	CheckUnauthorizedStatus(t, resp)

	teamUnread, resp = th.SystemAdminClient.GetTeamUnread(th.BasicTeam.Id, th.BasicUser.Id)
	CheckNoError(t, resp)
	if teamUnread.TeamId != th.BasicTeam.Id {
		t.Fatal("wrong team id returned")
	}
}

func TestUpdateTeam(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	team := &model.Team{DisplayName: "Name", Description: "Some description", AllowOpenInvite: false, InviteId: "inviteid0", Name: "z-z-" + model.NewId() + "a", Email: "success+" + model.NewId() + "@simulator.amazonses.com", Type: model.TEAM_OPEN}
	team, _ = Client.CreateTeam(team)

	team.Description = "updated description"
	uteam, resp := Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.Description != "updated description" {
		t.Fatal("Update failed")
	}

	team.DisplayName = "Updated Name"
	uteam, resp = Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.DisplayName != "Updated Name" {
		t.Fatal("Update failed")
	}

	team.AllowOpenInvite = true
	uteam, resp = Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.AllowOpenInvite != true {
		t.Fatal("Update failed")
	}

	team.InviteId = "inviteid1"
	uteam, resp = Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.InviteId != "inviteid1" {
		t.Fatal("Update failed")
	}

	team.Name = "Updated name"
	uteam, resp = Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.Name == "Updated name" {
		t.Fatal("Should not update name")
	}

	team.Email = "test@domain.com"
	uteam, resp = Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.Email == "test@domain.com" {
		t.Fatal("Should not update email")
	}

	team.Type = model.TEAM_INVITE
	uteam, resp = Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.Type == model.TEAM_INVITE {
		t.Fatal("Should not update type")
	}

	team.AllowedDomains = "domain"
	uteam, resp = Client.UpdateTeam(team)
	CheckNoError(t, resp)

	if uteam.AllowedDomains == "domain" {
		t.Fatal("Should not update allowed_domains")
	}

	originalTeamId := team.Id
	team.Id = model.NewId()

	if r, err := Client.DoApiPut(Client.GetTeamRoute(originalTeamId), team.ToJson()); err != nil {
		t.Fatal(err)
	} else {
		uteam = model.TeamFromJson(r.Body)
	}

	if uteam.Id != originalTeamId {
		t.Fatal("wrong team id")
	}

	team.Id = "fake"
	_, resp = Client.UpdateTeam(team)
	CheckBadRequestStatus(t, resp)

	Client.Logout()
	_, resp = Client.UpdateTeam(team)
	CheckUnauthorizedStatus(t, resp)

	team.Id = originalTeamId
	_, resp = th.SystemAdminClient.UpdateTeam(team)
	CheckNoError(t, resp)
}

func TestPatchTeam(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	team := &model.Team{DisplayName: "Name", Description: "Some description", CompanyName: "Some company name", AllowOpenInvite: false, InviteId: "inviteid0", Name: "z-z-" + model.NewId() + "a", Email: "success+" + model.NewId() + "@simulator.amazonses.com", Type: model.TEAM_OPEN}
	team, _ = Client.CreateTeam(team)

	patch := &model.TeamPatch{}

	patch.DisplayName = new(string)
	*patch.DisplayName = "Other name"
	patch.Description = new(string)
	*patch.Description = "Other description"
	patch.CompanyName = new(string)
	*patch.CompanyName = "Other company name"
	patch.InviteId = new(string)
	*patch.InviteId = "inviteid1"
	patch.AllowOpenInvite = new(bool)
	*patch.AllowOpenInvite = true

	rteam, resp := Client.PatchTeam(team.Id, patch)
	CheckNoError(t, resp)
	CheckTeamSanitization(t, rteam)

	if rteam.DisplayName != "Other name" {
		t.Fatal("DisplayName did not update properly")
	}
	if rteam.Description != "Other description" {
		t.Fatal("Description did not update properly")
	}
	if rteam.CompanyName != "Other company name" {
		t.Fatal("CompanyName did not update properly")
	}
	if rteam.InviteId != "inviteid1" {
		t.Fatal("InviteId did not update properly")
	}
	if rteam.AllowOpenInvite != true {
		t.Fatal("AllowOpenInvite did not update properly")
	}

	_, resp = Client.PatchTeam("junk", patch)
	CheckBadRequestStatus(t, resp)

	_, resp = Client.PatchTeam(GenerateTestId(), patch)
	CheckForbiddenStatus(t, resp)

	if r, err := Client.DoApiPut("/teams/"+team.Id+"/patch", "garbage"); err == nil {
		t.Fatal("should have errored")
	} else {
		if r.StatusCode != http.StatusBadRequest {
			t.Log("actual: " + strconv.Itoa(r.StatusCode))
			t.Log("expected: " + strconv.Itoa(http.StatusBadRequest))
			t.Fatal("wrong status code")
		}
	}

	Client.Logout()
	_, resp = Client.PatchTeam(team.Id, patch)
	CheckUnauthorizedStatus(t, resp)

	th.LoginBasic2()
	_, resp = Client.PatchTeam(team.Id, patch)
	CheckForbiddenStatus(t, resp)

	_, resp = th.SystemAdminClient.PatchTeam(team.Id, patch)
	CheckNoError(t, resp)
}

func TestGetAllTeams(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	team := &model.Team{DisplayName: "Name", Name: GenerateTestTeamName(), Email: GenerateTestEmail(), Type: model.TEAM_OPEN, AllowOpenInvite: true}
	_, resp := Client.CreateTeam(team)
	CheckNoError(t, resp)

	rrteams, resp := Client.GetAllTeams("", 0, 1)
	CheckNoError(t, resp)

	if len(rrteams) != 1 {
		t.Log(len(rrteams))
		t.Fatal("wrong number of teams - should be 1")
	}

	for _, rt := range rrteams {
		if rt.Type != model.TEAM_OPEN {
			t.Fatal("not all teams are open")
		}
	}

	rrteams1, resp := Client.GetAllTeams("", 1, 0)
	CheckNoError(t, resp)

	if len(rrteams1) != 0 {
		t.Fatal("wrong number of teams - should be 0")
	}

	rrteams2, resp := th.SystemAdminClient.GetAllTeams("", 1, 1)
	CheckNoError(t, resp)

	if len(rrteams2) != 1 {
		t.Fatal("wrong number of teams - should be 1")
	}

	rrteams2, resp = Client.GetAllTeams("", 1, 0)
	CheckNoError(t, resp)

	if len(rrteams2) != 0 {
		t.Fatal("wrong number of teams - should be 0")
	}

	Client.Logout()
	_, resp = Client.GetAllTeams("", 1, 10)
	CheckUnauthorizedStatus(t, resp)
}

func TestGetTeamByName(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	team := th.BasicTeam

	rteam, resp := Client.GetTeamByName(team.Name, "")
	CheckNoError(t, resp)

	if rteam.Name != team.Name {
		t.Fatal("wrong team")
	}

	_, resp = Client.GetTeamByName("junk", "")
	CheckNotFoundStatus(t, resp)

	_, resp = Client.GetTeamByName("", "")
	CheckNotFoundStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetTeamByName(team.Name, "")
	CheckUnauthorizedStatus(t, resp)

	_, resp = th.SystemAdminClient.GetTeamByName(team.Name, "")
	CheckNoError(t, resp)

	th.LoginTeamAdmin()

	team2 := &model.Team{DisplayName: "Name", Name: GenerateTestTeamName(), Email: GenerateTestEmail(), Type: model.TEAM_INVITE}
	rteam2, _ := Client.CreateTeam(team2)

	th.LoginBasic()
	_, resp = Client.GetTeamByName(rteam2.Name, "")
	CheckForbiddenStatus(t, resp)
}

func TestGetTeamsForUser(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	team2 := &model.Team{DisplayName: "Name", Name: GenerateTestTeamName(), Email: GenerateTestEmail(), Type: model.TEAM_INVITE}
	rteam2, _ := Client.CreateTeam(team2)

	teams, resp := Client.GetTeamsForUser(th.BasicUser.Id, "")
	CheckNoError(t, resp)

	if len(teams) != 2 {
		t.Fatal("wrong number of teams")
	}

	found1 := false
	found2 := false
	for _, t := range teams {
		if t.Id == th.BasicTeam.Id {
			found1 = true
		} else if t.Id == rteam2.Id {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Fatal("missing team")
	}

	_, resp = Client.GetTeamsForUser("junk", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamsForUser(model.NewId(), "")
	CheckForbiddenStatus(t, resp)

	_, resp = Client.GetTeamsForUser(th.BasicUser2.Id, "")
	CheckForbiddenStatus(t, resp)

	_, resp = th.SystemAdminClient.GetTeamsForUser(th.BasicUser2.Id, "")
	CheckNoError(t, resp)
}

func TestGetTeamMember(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	team := th.BasicTeam
	user := th.BasicUser

	rmember, resp := Client.GetTeamMember(team.Id, user.Id, "")
	CheckNoError(t, resp)

	if rmember.TeamId != team.Id {
		t.Fatal("wrong team id")
	}

	if rmember.UserId != user.Id {
		t.Fatal("wrong team id")
	}

	_, resp = Client.GetTeamMember("junk", user.Id, "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamMember(team.Id, "junk", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamMember("junk", "junk", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamMember(team.Id, model.NewId(), "")
	CheckNotFoundStatus(t, resp)

	_, resp = Client.GetTeamMember(model.NewId(), user.Id, "")
	CheckForbiddenStatus(t, resp)

	_, resp = th.SystemAdminClient.GetTeamMember(team.Id, user.Id, "")
	CheckNoError(t, resp)
}

func TestGetTeamMembers(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	team := th.BasicTeam
	userNotMember := th.CreateUser()

	rmembers, resp := Client.GetTeamMembers(team.Id, 0, 100, "")
	CheckNoError(t, resp)

	t.Logf("rmembers count %v\n", len(rmembers))

	if len(rmembers) == 0 {
		t.Fatal("should have results")
	}

	for _, rmember := range rmembers {
		if rmember.TeamId != team.Id || rmember.UserId == userNotMember.Id {
			t.Fatal("user should be a member of team")
		}
	}

	rmembers, resp = Client.GetTeamMembers(team.Id, 0, 1, "")
	CheckNoError(t, resp)
	if len(rmembers) != 1 {
		t.Fatal("should be 1 per page")
	}

	rmembers, resp = Client.GetTeamMembers(team.Id, 1, 1, "")
	CheckNoError(t, resp)
	if len(rmembers) != 1 {
		t.Fatal("should be 1 per page")
	}

	rmembers, resp = Client.GetTeamMembers(team.Id, 10000, 100, "")
	CheckNoError(t, resp)
	if len(rmembers) != 0 {
		t.Fatal("should be no member")
	}

	_, resp = Client.GetTeamMembers("junk", 0, 100, "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamMembers(model.NewId(), 0, 100, "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	rmembers, resp = Client.GetTeamMembers(team.Id, 0, 1, "")
	CheckUnauthorizedStatus(t, resp)

	rmembers, resp = th.SystemAdminClient.GetTeamMembers(team.Id, 0, 100, "")
	CheckNoError(t, resp)
}

func TestGetTeamMembersByIds(t *testing.T) {
	th := Setup().InitBasic()
	defer TearDown()
	Client := th.Client

	tm, resp := Client.GetTeamMembersByIds(th.BasicTeam.Id, []string{th.BasicUser.Id})
	CheckNoError(t, resp)

	if tm[0].UserId != th.BasicUser.Id {
		t.Fatal("returned wrong user")
	}

	_, resp = Client.GetTeamMembersByIds(th.BasicTeam.Id, []string{})
	CheckBadRequestStatus(t, resp)

	tm1, resp := Client.GetTeamMembersByIds(th.BasicTeam.Id, []string{"junk"})
	CheckNoError(t, resp)
	if len(tm1) > 0 {
		t.Fatal("no users should be returned")
	}

	tm1, resp = Client.GetTeamMembersByIds(th.BasicTeam.Id, []string{"junk", th.BasicUser.Id})
	CheckNoError(t, resp)
	if len(tm1) != 1 {
		t.Fatal("1 user should be returned")
	}

	tm1, resp = Client.GetTeamMembersByIds("junk", []string{th.BasicUser.Id})
	CheckBadRequestStatus(t, resp)

	tm1, resp = Client.GetTeamMembersByIds(model.NewId(), []string{th.BasicUser.Id})
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetTeamMembersByIds(th.BasicTeam.Id, []string{th.BasicUser.Id})
	CheckUnauthorizedStatus(t, resp)
}

func TestAddTeamMember(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	team := th.BasicTeam
	otherUser := th.CreateUser()

	// by user_id
	th.LoginTeamAdmin()

	tm, resp := Client.AddTeamMember(team.Id, otherUser.Id, "", "", "")
	CheckNoError(t, resp)

	if tm == nil {
		t.Fatal("should have returned team member")
	}

	if tm.UserId != otherUser.Id {
		t.Fatal("user ids should have matched")
	}

	if tm.TeamId != team.Id {
		t.Fatal("team ids should have matched")
	}

	tm, resp = Client.AddTeamMember(team.Id, "junk", "", "", "")
	CheckBadRequestStatus(t, resp)

	if tm != nil {
		t.Fatal("should have not returned team member")
	}

	_, resp = Client.AddTeamMember("junk", otherUser.Id, "", "", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.AddTeamMember(GenerateTestId(), otherUser.Id, "", "", "")
	CheckForbiddenStatus(t, resp)

	_, resp = Client.AddTeamMember(team.Id, GenerateTestId(), "", "", "")
	CheckNotFoundStatus(t, resp)

	th.LoginBasic()

	_, resp = Client.AddTeamMember(team.Id, otherUser.Id, "", "", "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()

	_, resp = Client.AddTeamMember(team.Id, otherUser.Id, "", "", "")
	CheckUnauthorizedStatus(t, resp)

	_, resp = th.SystemAdminClient.AddTeamMember(team.Id, otherUser.Id, "", "", "")
	CheckNoError(t, resp)

	// by hash and data
	Client.Login(otherUser.Email, otherUser.Password)

	dataObject := make(map[string]string)
	dataObject["time"] = fmt.Sprintf("%v", model.GetMillis())
	dataObject["id"] = team.Id

	data := model.MapToJson(dataObject)
	hashed := model.HashPassword(fmt.Sprintf("%v:%v", data, utils.Cfg.EmailSettings.InviteSalt))

	tm, resp = Client.AddTeamMember(team.Id, "", hashed, data, "")
	CheckNoError(t, resp)

	if tm == nil {
		t.Fatal("should have returned team member")
	}

	if tm.UserId != otherUser.Id {
		t.Fatal("user ids should have matched")
	}

	if tm.TeamId != team.Id {
		t.Fatal("team ids should have matched")
	}

	tm, resp = Client.AddTeamMember(team.Id, "", "junk", data, "")
	CheckNotFoundStatus(t, resp)

	if tm != nil {
		t.Fatal("should have not returned team member")
	}

	_, resp = Client.AddTeamMember(team.Id, "", hashed, "junk", "")
	CheckNotFoundStatus(t, resp)

	// expired data of more than 50 hours
	dataObject["time"] = fmt.Sprintf("%v", model.GetMillis()-1000*60*60*50)
	data = model.MapToJson(dataObject)
	hashed = model.HashPassword(fmt.Sprintf("%v:%v", data, utils.Cfg.EmailSettings.InviteSalt))

	tm, resp = Client.AddTeamMember(team.Id, "", hashed, data, "")
	CheckNotFoundStatus(t, resp)

	// invalid team id
	dataObject["id"] = GenerateTestId()
	data = model.MapToJson(dataObject)
	hashed = model.HashPassword(fmt.Sprintf("%v:%v", data, utils.Cfg.EmailSettings.InviteSalt))

	tm, resp = Client.AddTeamMember(team.Id, "", hashed, data, "")
	CheckNotFoundStatus(t, resp)

	// by invite_id
	Client.Login(otherUser.Email, otherUser.Password)

	tm, resp = Client.AddTeamMember(team.Id, "", "", "", team.InviteId)
	CheckNoError(t, resp)

	if tm == nil {
		t.Fatal("should have returned team member")
	}

	if tm.UserId != otherUser.Id {
		t.Fatal("user ids should have matched")
	}

	if tm.TeamId != team.Id {
		t.Fatal("team ids should have matched")
	}

	tm, resp = Client.AddTeamMember(team.Id, "", "", "", "junk")
	CheckNotFoundStatus(t, resp)

	if tm != nil {
		t.Fatal("should have not returned team member")
	}

	_, resp = Client.AddTeamMember(team.Id, "", "", "", "junk")
	CheckNotFoundStatus(t, resp)
}

func TestGetTeamStats(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	team := th.BasicTeam

	rstats, resp := Client.GetTeamStats(team.Id, "")
	CheckNoError(t, resp)

	if rstats.TeamId != team.Id {
		t.Fatal("wrong team id")
	}

	if rstats.TotalMemberCount != 3 {
		t.Fatal("wrong count")
	}

	if rstats.ActiveMemberCount != 3 {
		t.Fatal("wrong count")
	}

	_, resp = Client.GetTeamStats("junk", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamStats(model.NewId(), "")
	CheckForbiddenStatus(t, resp)

	_, resp = th.SystemAdminClient.GetTeamStats(team.Id, "")
	CheckNoError(t, resp)

	// deactivate BasicUser2
	th.UpdateActiveUser(th.BasicUser2, false)

	rstats, resp = th.SystemAdminClient.GetTeamStats(team.Id, "")
	CheckNoError(t, resp)

	if rstats.TotalMemberCount != 3 {
		t.Fatal("wrong count")
	}

	if rstats.ActiveMemberCount != 2 {
		t.Fatal("wrong count")
	}

	// login with different user and test if forbidden
	user := th.CreateUser()
	Client.Login(user.Email, user.Password)
	_, resp = Client.GetTeamStats(th.BasicTeam.Id, "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetTeamStats(th.BasicTeam.Id, "")
	CheckUnauthorizedStatus(t, resp)
}

func TestUpdateTeamMemberRoles(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	SystemAdminClient := th.SystemAdminClient

	const TEAM_MEMBER = "team_user"
	const TEAM_ADMIN = "team_user team_admin"

	// user 1 tries to promote user 2
	ok, resp := Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser2.Id, TEAM_ADMIN)
	CheckForbiddenStatus(t, resp)
	if ok {
		t.Fatal("should have returned false")
	}

	// user 1 tries to promote himself
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser.Id, TEAM_ADMIN)
	CheckForbiddenStatus(t, resp)

	// user 1 tries to demote someone
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.SystemAdminUser.Id, TEAM_MEMBER)
	CheckForbiddenStatus(t, resp)

	// system admin promotes user 1
	ok, resp = SystemAdminClient.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser.Id, TEAM_ADMIN)
	CheckNoError(t, resp)
	if !ok {
		t.Fatal("should have returned true")
	}

	// user 1 (team admin) promotes user 2
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser2.Id, TEAM_ADMIN)
	CheckNoError(t, resp)

	// user 1 (team admin) demotes user 2 (team admin)
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser2.Id, TEAM_MEMBER)
	CheckNoError(t, resp)

	// user 1 (team admin) tries to demote system admin (not member of a team)
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.SystemAdminUser.Id, TEAM_MEMBER)
	CheckBadRequestStatus(t, resp)

	// user 1 (team admin) demotes system admin (member of a team)
	LinkUserToTeam(th.SystemAdminUser, th.BasicTeam)
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.SystemAdminUser.Id, TEAM_MEMBER)
	CheckNoError(t, resp)
	// Note from API v3
	// Note to anyone who thinks this (above) test is wrong:
	// This operation will not affect the system admin's permissions because they have global access to all teams.
	// Their team level permissions are irrelavent. A team admin should be able to manage team level permissions.

	// System admins should be able to manipulate permission no matter what their team level permissions are.
	// system admin promotes user 2
	_, resp = SystemAdminClient.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser2.Id, TEAM_ADMIN)
	CheckNoError(t, resp)

	// system admin demotes user 2 (team admin)
	_, resp = SystemAdminClient.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser2.Id, TEAM_MEMBER)
	CheckNoError(t, resp)

	// user 1 (team admin) tries to promote himself to a random team
	_, resp = Client.UpdateTeamMemberRoles(model.NewId(), th.BasicUser.Id, TEAM_ADMIN)
	CheckForbiddenStatus(t, resp)

	// user 1 (team admin) tries to promote a random user
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, model.NewId(), TEAM_ADMIN)
	CheckBadRequestStatus(t, resp)

	// user 1 (team admin) tries to promote invalid team permission
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser.Id, "junk")
	CheckBadRequestStatus(t, resp)

	// user 1 (team admin) demotes himself
	_, resp = Client.UpdateTeamMemberRoles(th.BasicTeam.Id, th.BasicUser.Id, TEAM_MEMBER)
	CheckNoError(t, resp)
}

func TestGetMyTeamsUnread(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client

	user := th.BasicUser
	Client.Login(user.Email, user.Password)

	teams, resp := Client.GetTeamsUnreadForUser(user.Id, "")
	CheckNoError(t, resp)
	if len(teams) == 0 {
		t.Fatal("should have results")
	}

	teams, resp = Client.GetTeamsUnreadForUser(user.Id, th.BasicTeam.Id)
	CheckNoError(t, resp)
	if len(teams) != 0 {
		t.Fatal("should not have results")
	}

	_, resp = Client.GetTeamsUnreadForUser("fail", "")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetTeamsUnreadForUser(model.NewId(), "")
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetTeamsUnreadForUser(user.Id, "")
	CheckUnauthorizedStatus(t, resp)
}

func TestTeamExists(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()
	defer TearDown()
	Client := th.Client
	team := th.BasicTeam

	th.LoginBasic()

	exists, resp := Client.TeamExists(team.Name, "")
	CheckNoError(t, resp)
	if exists != true {
		t.Fatal("team should exist")
	}

	exists, resp = Client.TeamExists("testingteam", "")
	CheckNoError(t, resp)
	if exists != false {
		t.Fatal("team should not exist")
	}

	Client.Logout()
	_, resp = Client.TeamExists(team.Name, "")
	CheckUnauthorizedStatus(t, resp)
}
