// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"net/http"
	"strings"

	l4g "code.google.com/p/log4go"
	"github.com/gorilla/mux"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

type CommandProvider interface {
	GetCommand() *model.Command
	DoCommand(c *Context, channelId string, message string) *model.CommandResponse
}

var commandProviders = make(map[string]CommandProvider)

func RegisterCommandProvider(newProvider CommandProvider) {
	commandProviders[newProvider.GetCommand().Trigger] = newProvider
}

func GetCommandProvidersProvider(name string) CommandProvider {
	provider, ok := commandProviders[name]
	if ok {
		return provider
	}

	return nil
}

func InitCommand(r *mux.Router) {
	l4g.Debug("Initializing command api routes")

	sr := r.PathPrefix("/commands").Subrouter()

	sr.Handle("/execute", ApiUserRequired(executeCommand)).Methods("POST")
	sr.Handle("/list", ApiUserRequired(listCommands)).Methods("GET")

	sr.Handle("/create", ApiUserRequired(createCommand)).Methods("POST")
	sr.Handle("/list_team_commands", ApiUserRequired(listTeamCommands)).Methods("GET")
	sr.Handle("/regen_token", ApiUserRequired(regenCommandToken)).Methods("POST")
	sr.Handle("/delete", ApiUserRequired(deleteCommand)).Methods("POST")
}

func listCommands(c *Context, w http.ResponseWriter, r *http.Request) {
	commands := make([]*model.Command, 0, 32)
	seen := make(map[string]bool)
	for _, value := range commandProviders {
		cpy := *value.GetCommand()
		if cpy.AutoComplete && !seen[cpy.Id] {
			cpy.Sanatize()
			seen[cpy.Trigger] = true
			commands = append(commands, &cpy)
		}
	}

	if result := <-Srv.Store.Command().GetByTeam(c.Session.TeamId); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		teamCmds := result.Data.([]*model.Command)
		for _, cmd := range teamCmds {
			if cmd.AutoComplete && !seen[cmd.Id] {
				cmd.Sanatize()
				seen[cmd.Trigger] = true
				commands = append(commands, cmd)
			}
		}
	}

	w.Write([]byte(model.CommandListToJson(commands)))
}

func executeCommand(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)
	command := strings.TrimSpace(props["command"])
	channelId := strings.TrimSpace(props["channelId"])

	if len(command) <= 1 || strings.Index(command, "/") != 0 {
		c.Err = model.NewAppError("command", "Command must start with /", "")
		return
	}

	if len(channelId) > 0 {
		cchan := Srv.Store.Channel().CheckPermissionsTo(c.Session.TeamId, channelId, c.Session.UserId)

		if !c.HasPermissionsToChannel(cchan, "checkCommand") {
			return
		}
	}

	parts := strings.Split(command, " ")
	trigger := parts[0][1:]
	provider := GetCommandProvidersProvider(trigger)

	if provider != nil {
		message := strings.Join(parts[1:], " ")
		response := provider.DoCommand(c, channelId, message)

		if response.ResponseType == model.COMMAND_RESPONSE_TYPE_IN_CHANNEL {
			post := &model.Post{}
			post.ChannelId = channelId
			post.Message = response.Text
			if _, err := CreatePost(c, post, true); err != nil {
				c.Err = model.NewAppError("command", "An error while saving the command response to the channel", "")
			}
		} else if response.ResponseType == model.COMMAND_RESPONSE_TYPE_EPHEMERAL {
			post := &model.Post{}
			post.ChannelId = channelId
			post.Message = "TODO_EPHEMERAL: " + response.Text
			if _, err := CreatePost(c, post, true); err != nil {
				c.Err = model.NewAppError("command", "An error while saving the command response to the channel", "")
			}
		}

		w.Write([]byte(response.ToJson()))
	} else {
		c.Err = model.NewAppError("command", "Command with a trigger of '"+trigger+"' not found", "")
	}
}

func createCommand(c *Context, w http.ResponseWriter, r *http.Request) {
	if !*utils.Cfg.ServiceSettings.EnableCommands {
		c.Err = model.NewAppError("createCommand", "Commands have been disabled by the system admin.", "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	if *utils.Cfg.ServiceSettings.EnableOnlyAdminIntegrations {
		if !(c.IsSystemAdmin() || c.IsTeamAdmin()) {
			c.Err = model.NewAppError("createCommand", "Integrations have been limited to admins only.", "")
			c.Err.StatusCode = http.StatusForbidden
			return
		}
	}

	c.LogAudit("attempt")

	cmd := model.CommandFromJson(r.Body)

	if cmd == nil {
		c.SetInvalidParam("createCommand", "command")
		return
	}

	cmd.CreatorId = c.Session.UserId
	cmd.TeamId = c.Session.TeamId

	if result := <-Srv.Store.Command().Save(cmd); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		c.LogAudit("success")
		rcmd := result.Data.(*model.Command)
		w.Write([]byte(rcmd.ToJson()))
	}
}

func listTeamCommands(c *Context, w http.ResponseWriter, r *http.Request) {
	if !*utils.Cfg.ServiceSettings.EnableCommands {
		c.Err = model.NewAppError("createCommand", "Commands have been disabled by the system admin.", "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	if *utils.Cfg.ServiceSettings.EnableOnlyAdminIntegrations {
		if !(c.IsSystemAdmin() || c.IsTeamAdmin()) {
			c.Err = model.NewAppError("createCommand", "Integrations have been limited to admins only.", "")
			c.Err.StatusCode = http.StatusForbidden
			return
		}
	}

	if result := <-Srv.Store.Command().GetByTeam(c.Session.TeamId); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		cmds := result.Data.([]*model.Command)
		w.Write([]byte(model.CommandListToJson(cmds)))
	}
}

func regenCommandToken(c *Context, w http.ResponseWriter, r *http.Request) {
	if !*utils.Cfg.ServiceSettings.EnableCommands {
		c.Err = model.NewAppError("createCommand", "Commands have been disabled by the system admin.", "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	if *utils.Cfg.ServiceSettings.EnableOnlyAdminIntegrations {
		if !(c.IsSystemAdmin() || c.IsTeamAdmin()) {
			c.Err = model.NewAppError("createCommand", "Integrations have been limited to admins only.", "")
			c.Err.StatusCode = http.StatusForbidden
			return
		}
	}

	c.LogAudit("attempt")

	props := model.MapFromJson(r.Body)

	id := props["id"]
	if len(id) == 0 {
		c.SetInvalidParam("regenCommandToken", "id")
		return
	}

	var cmd *model.Command
	if result := <-Srv.Store.Command().Get(id); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		cmd = result.Data.(*model.Command)

		if c.Session.TeamId != cmd.TeamId && c.Session.UserId != cmd.CreatorId && !c.IsTeamAdmin() {
			c.LogAudit("fail - inappropriate permissions")
			c.Err = model.NewAppError("regenToken", "Inappropriate permissions to regenerate command token", "user_id="+c.Session.UserId)
			return
		}
	}

	cmd.Token = model.NewId()

	if result := <-Srv.Store.Command().Update(cmd); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		w.Write([]byte(result.Data.(*model.Command).ToJson()))
	}
}

func deleteCommand(c *Context, w http.ResponseWriter, r *http.Request) {
	if !*utils.Cfg.ServiceSettings.EnableCommands {
		c.Err = model.NewAppError("createCommand", "Commands have been disabled by the system admin.", "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	if *utils.Cfg.ServiceSettings.EnableOnlyAdminIntegrations {
		if !(c.IsSystemAdmin() || c.IsTeamAdmin()) {
			c.Err = model.NewAppError("createCommand", "Integrations have been limited to admins only.", "")
			c.Err.StatusCode = http.StatusForbidden
			return
		}
	}

	c.LogAudit("attempt")

	props := model.MapFromJson(r.Body)

	id := props["id"]
	if len(id) == 0 {
		c.SetInvalidParam("deleteCommand", "id")
		return
	}

	if result := <-Srv.Store.Command().Get(id); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		if c.Session.TeamId != result.Data.(*model.Command).TeamId && c.Session.UserId != result.Data.(*model.Command).CreatorId && !c.IsTeamAdmin() {
			c.LogAudit("fail - inappropriate permissions")
			c.Err = model.NewAppError("deleteCommand", "Inappropriate permissions to delete command", "user_id="+c.Session.UserId)
			return
		}
	}

	if err := (<-Srv.Store.Command().Delete(id, model.GetMillis())).Err; err != nil {
		c.Err = err
		return
	}

	c.LogAudit("success")
	w.Write([]byte(model.MapToJson(props)))
}

// func loadTestCommand(c *Context, command *model.Command) bool {
// 	cmd := cmds["loadTestCommand"]

// 	// This command is only available when EnableTesting is true
// 	if !utils.Cfg.ServiceSettings.EnableTesting {
// 		return false
// 	}

// 	if strings.Index(command.Command, cmd) == 0 {
// 		if loadTestSetupCommand(c, command) {
// 			return true
// 		}
// 		if loadTestUsersCommand(c, command) {
// 			return true
// 		}
// 		if loadTestChannelsCommand(c, command) {
// 			return true
// 		}
// 		if loadTestPostsCommand(c, command) {
// 			return true
// 		}
// 		if loadTestUrlCommand(c, command) {
// 			return true
// 		}
// 	} else if strings.Index(cmd, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd, Description: "Debug Load Testing"})
// 	}

// 	return false
// }

// func parseRange(command string, cmd string) (utils.Range, bool) {
// 	tokens := strings.Fields(strings.TrimPrefix(command, cmd))
// 	var begin int
// 	var end int
// 	var err1 error
// 	var err2 error
// 	switch {
// 	case len(tokens) == 1:
// 		begin, err1 = strconv.Atoi(tokens[0])
// 		end = begin
// 		if err1 != nil {
// 			return utils.Range{0, 0}, false
// 		}
// 	case len(tokens) >= 2:
// 		begin, err1 = strconv.Atoi(tokens[0])
// 		end, err2 = strconv.Atoi(tokens[1])
// 		if err1 != nil || err2 != nil {
// 			return utils.Range{0, 0}, false
// 		}
// 	default:
// 		return utils.Range{0, 0}, false
// 	}
// 	return utils.Range{begin, end}, true
// }

// func contains(items []string, token string) bool {
// 	for _, elem := range items {
// 		if elem == token {
// 			return true
// 		}
// 	}
// 	return false
// }

// func loadTestSetupCommand(c *Context, command *model.Command) bool {
// 	cmd := cmds["loadTestCommand"] + " setup"

// 	if strings.Index(command.Command, cmd) == 0 && !command.Suggest {
// 		tokens := strings.Fields(strings.TrimPrefix(command.Command, cmd))
// 		doTeams := contains(tokens, "teams")
// 		doFuzz := contains(tokens, "fuzz")

// 		numArgs := 0
// 		if doTeams {
// 			numArgs++
// 		}
// 		if doFuzz {
// 			numArgs++
// 		}

// 		var numTeams int
// 		var numChannels int
// 		var numUsers int
// 		var numPosts int

// 		// Defaults
// 		numTeams = 10
// 		numChannels = 10
// 		numUsers = 10
// 		numPosts = 10

// 		if doTeams {
// 			if (len(tokens) - numArgs) >= 4 {
// 				numTeams, _ = strconv.Atoi(tokens[numArgs+0])
// 				numChannels, _ = strconv.Atoi(tokens[numArgs+1])
// 				numUsers, _ = strconv.Atoi(tokens[numArgs+2])
// 				numPosts, _ = strconv.Atoi(tokens[numArgs+3])
// 			}
// 		} else {
// 			if (len(tokens) - numArgs) >= 3 {
// 				numChannels, _ = strconv.Atoi(tokens[numArgs+0])
// 				numUsers, _ = strconv.Atoi(tokens[numArgs+1])
// 				numPosts, _ = strconv.Atoi(tokens[numArgs+2])
// 			}
// 		}
// 		client := model.NewClient(c.GetSiteURL())

// 		if doTeams {
// 			if err := CreateBasicUser(client); err != nil {
// 				l4g.Error("Failed to create testing environment")
// 				return true
// 			}
// 			client.LoginByEmail(BTEST_TEAM_NAME, BTEST_USER_EMAIL, BTEST_USER_PASSWORD)
// 			environment, err := CreateTestEnvironmentWithTeams(
// 				client,
// 				utils.Range{numTeams, numTeams},
// 				utils.Range{numChannels, numChannels},
// 				utils.Range{numUsers, numUsers},
// 				utils.Range{numPosts, numPosts},
// 				doFuzz)
// 			if err != true {
// 				l4g.Error("Failed to create testing environment")
// 				return true
// 			} else {
// 				l4g.Info("Testing environment created")
// 				for i := 0; i < len(environment.Teams); i++ {
// 					l4g.Info("Team Created: " + environment.Teams[i].Name)
// 					l4g.Info("\t User to login: " + environment.Environments[i].Users[0].Email + ", " + USER_PASSWORD)
// 				}
// 			}
// 		} else {
// 			client.MockSession(c.Session.Token)
// 			CreateTestEnvironmentInTeam(
// 				client,
// 				c.Session.TeamId,
// 				utils.Range{numChannels, numChannels},
// 				utils.Range{numUsers, numUsers},
// 				utils.Range{numPosts, numPosts},
// 				doFuzz)
// 		}
// 		return true
// 	} else if strings.Index(cmd, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{
// 			Suggestion:  cmd,
// 			Description: "Creates a testing environment in current team. [teams] [fuzz] <Num Channels> <Num Users> <NumPosts>"})
// 	}

// 	return false
// }

// func loadTestUsersCommand(c *Context, command *model.Command) bool {
// 	cmd1 := cmds["loadTestCommand"] + " users"
// 	cmd2 := cmds["loadTestCommand"] + " users fuzz"

// 	if strings.Index(command.Command, cmd1) == 0 && !command.Suggest {
// 		cmd := cmd1
// 		doFuzz := false
// 		if strings.Index(command.Command, cmd2) == 0 {
// 			doFuzz = true
// 			cmd = cmd2
// 		}
// 		usersr, err := parseRange(command.Command, cmd)
// 		if err == false {
// 			usersr = utils.Range{10, 15}
// 		}
// 		client := model.NewClient(c.GetSiteURL())
// 		userCreator := NewAutoUserCreator(client, c.Session.TeamId)
// 		userCreator.Fuzzy = doFuzz
// 		userCreator.CreateTestUsers(usersr)
// 		return true
// 	} else if strings.Index(cmd1, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd1, Description: "Add a specified number of random users to current team <Min Users> <Max Users>"})
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd2, Description: "Add a specified number of random users with fuzz text to current team <Min Users> <Max Users>"})
// 	} else if strings.Index(cmd2, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd2, Description: "Add a specified number of random users with fuzz text to current team <Min Users> <Max Users>"})
// 	}

// 	return false
// }

// func loadTestChannelsCommand(c *Context, command *model.Command) bool {
// 	cmd1 := cmds["loadTestCommand"] + " channels"
// 	cmd2 := cmds["loadTestCommand"] + " channels fuzz"

// 	if strings.Index(command.Command, cmd1) == 0 && !command.Suggest {
// 		cmd := cmd1
// 		doFuzz := false
// 		if strings.Index(command.Command, cmd2) == 0 {
// 			doFuzz = true
// 			cmd = cmd2
// 		}
// 		channelsr, err := parseRange(command.Command, cmd)
// 		if err == false {
// 			channelsr = utils.Range{20, 30}
// 		}
// 		client := model.NewClient(c.GetSiteURL())
// 		client.MockSession(c.Session.Token)
// 		channelCreator := NewAutoChannelCreator(client, c.Session.TeamId)
// 		channelCreator.Fuzzy = doFuzz
// 		channelCreator.CreateTestChannels(channelsr)
// 		return true
// 	} else if strings.Index(cmd1, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd1, Description: "Add a specified number of random channels to current team <MinChannels> <MaxChannels>"})
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd2, Description: "Add a specified number of random channels with fuzz text to current team <Min Channels> <Max Channels>"})
// 	} else if strings.Index(cmd2, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd2, Description: "Add a specified number of random channels with fuzz text to current team <Min Channels> <Max Channels>"})
// 	}

// 	return false
// }

// func loadTestPostsCommand(c *Context, command *model.Command) bool {
// 	cmd1 := cmds["loadTestCommand"] + " posts"
// 	cmd2 := cmds["loadTestCommand"] + " posts fuzz"

// 	if strings.Index(command.Command, cmd1) == 0 && !command.Suggest {
// 		cmd := cmd1
// 		doFuzz := false
// 		if strings.Index(command.Command, cmd2) == 0 {
// 			cmd = cmd2
// 			doFuzz = true
// 		}

// 		postsr, err := parseRange(command.Command, cmd)
// 		if err == false {
// 			postsr = utils.Range{20, 30}
// 		}

// 		tokens := strings.Fields(strings.TrimPrefix(command.Command, cmd))
// 		rimages := utils.Range{0, 0}
// 		if len(tokens) >= 3 {
// 			if numImages, err := strconv.Atoi(tokens[2]); err == nil {
// 				rimages = utils.Range{numImages, numImages}
// 			}
// 		}

// 		var usernames []string
// 		if result := <-Srv.Store.User().GetProfiles(c.Session.TeamId); result.Err == nil {
// 			profileUsers := result.Data.(map[string]*model.User)
// 			usernames = make([]string, len(profileUsers))
// 			i := 0
// 			for _, userprof := range profileUsers {
// 				usernames[i] = userprof.Username
// 				i++
// 			}
// 		}

// 		client := model.NewClient(c.GetSiteURL())
// 		client.MockSession(c.Session.Token)
// 		testPoster := NewAutoPostCreator(client, command.ChannelId)
// 		testPoster.Fuzzy = doFuzz
// 		testPoster.Users = usernames

// 		numImages := utils.RandIntFromRange(rimages)
// 		numPosts := utils.RandIntFromRange(postsr)
// 		for i := 0; i < numPosts; i++ {
// 			testPoster.HasImage = (i < numImages)
// 			testPoster.CreateRandomPost()
// 		}
// 		return true
// 	} else if strings.Index(cmd1, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd1, Description: "Add some random posts to current channel <Min Posts> <Max Posts> <Min Images> <Max Images>"})
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd2, Description: "Add some random posts with fuzz text to current channel <Min Posts> <Max Posts> <Min Images> <Max Images>"})
// 	} else if strings.Index(cmd2, command.Command) == 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd2, Description: "Add some random posts with fuzz text to current channel <Min Posts> <Max Posts> <Min Images> <Max Images>"})
// 	}

// 	return false
// }

// func loadTestUrlCommand(c *Context, command *model.Command) bool {
// 	cmd := cmds["loadTestCommand"] + " url"

// 	if strings.Index(command.Command, cmd) == 0 && !command.Suggest {
// 		url := ""

// 		parameters := strings.SplitN(command.Command, " ", 3)
// 		if len(parameters) != 3 {
// 			c.Err = model.NewAppError("loadTestUrlCommand", "Command must contain a url", "")
// 			return true
// 		} else {
// 			url = parameters[2]
// 		}

// 		// provide a shortcut to easily access tests stored in doc/developer/tests
// 		if !strings.HasPrefix(url, "http") {
// 			url = "https://raw.githubusercontent.com/mattermost/platform/master/doc/developer/tests/" + url

// 			if path.Ext(url) == "" {
// 				url += ".md"
// 			}
// 		}

// 		var contents io.ReadCloser
// 		if r, err := http.Get(url); err != nil {
// 			c.Err = model.NewAppError("loadTestUrlCommand", "Unable to get file", err.Error())
// 			return false
// 		} else if r.StatusCode > 400 {
// 			c.Err = model.NewAppError("loadTestUrlCommand", "Unable to get file", r.Status)
// 			return false
// 		} else {
// 			contents = r.Body
// 		}

// 		bytes := make([]byte, 4000)

// 		// break contents into 4000 byte posts
// 		for {
// 			length, err := contents.Read(bytes)
// 			if err != nil && err != io.EOF {
// 				c.Err = model.NewAppError("loadTestUrlCommand", "Encountered error reading file", err.Error())
// 				return false
// 			}

// 			if length == 0 {
// 				break
// 			}

// 			post := &model.Post{}
// 			post.Message = string(bytes[:length])
// 			post.ChannelId = command.ChannelId

// 			if _, err := CreatePost(c, post, false); err != nil {
// 				l4g.Error("Unable to create post, err=%v", err)
// 				return false
// 			}
// 		}

// 		command.Response = model.RESP_EXECUTED

// 		return true
// 	} else if strings.Index(cmd, command.Command) == 0 && strings.Index(command.Command, "/loadtest posts") != 0 {
// 		command.AddSuggestion(&model.SuggestCommand{Suggestion: cmd, Description: "Add a post containing the text from a given url to current channel <Url>"})
// 	}

// 	return false
// }
