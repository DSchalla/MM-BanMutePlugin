package main

import (
	"github.com/mattermost/mattermost-server/model"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/DSchalla/MM-BanMutePlugin/banmuteplugin"
	"errors"
	"fmt"
)

type BanMutePlugin struct {
	plugin.MattermostPlugin
	server *banmuteplugin.Server
	config *banmuteplugin.Config
}

func (c *BanMutePlugin) OnActivate() error {
	c.API.LogDebug("OnActivate Hook Start")
	var err error
	c.readConfig()
	c.server, err = banmuteplugin.NewServer(c.API, *c.config)

	if err != nil {
		c.API.LogError(fmt.Sprintf("NewBotServer returned error: %s", err))
		return errors.New(err.Error())
	}

	teams, appError := c.API.GetTeams()

	if appError != nil {
		c.API.LogError(fmt.Sprintf("GetTeams returned error: %s", appError.Message))
		return errors.New(appError.Message)
	}

	for _, team := range teams {
		c.API.RegisterCommand(&model.Command{
			Trigger:          "banuser",
			TeamId:           team.Id,
			AutoComplete:     true,
			AutoCompleteDesc: "Ban user for minutes from channel",
			AutoCompleteHint: "@user MINUTES",
			DisplayName:      "Ban User",
		})
		c.API.RegisterCommand(&model.Command{
			Trigger:          "unbanuser",
			TeamId:           team.Id,
			AutoComplete:     true,
			AutoCompleteDesc: "Remove ban from user from channel",
			AutoCompleteHint: "@user",
			DisplayName:      "Unban User",
		})
		c.API.RegisterCommand(&model.Command{
			Trigger:          "muteuser",
			TeamId:           team.Id,
			AutoComplete:     true,
			AutoCompleteDesc: "Mute user for minutes from channel",
			AutoCompleteHint: "@user MINUTES",
			DisplayName:      "Mute User",
		})
		c.API.RegisterCommand(&model.Command{
			Trigger:          "unmuteuser",
			TeamId:           team.Id,
			AutoComplete:     true,
			AutoCompleteDesc: "Unmute user from channel",
			AutoCompleteHint: "@user",
			DisplayName:      "Unmute User",
		})
		c.API.LogDebug(fmt.Sprintf("Registered Commands on Team %s", team.Id))
	}

	c.API.LogDebug("OnActivate Hook End")
	return nil
}

func (c *BanMutePlugin) OnConfigurationChange() error {
	err := c.readConfig()
	if err != nil {
		return err
	}
	return c.reloadConfig()
}

func (c *BanMutePlugin) ExecuteCommand(context *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	c.API.LogDebug(fmt.Sprintf("Execute Command: %s - %s - %s", args.UserId, args.ChannelId, args.Command))
	return &model.CommandResponse{Text: "Event Invoked", Type: model.COMMAND_RESPONSE_TYPE_EPHEMERAL}, nil
}

func (c *BanMutePlugin) MessageWillBePosted(context *plugin.Context, post *model.Post) (*model.Post, string) {
	if post.Props["from_banmuteplugin"] != nil && post.Props["from_banmuteplugin"].(bool) == true {
		return post, ""
	}
	c.API.LogDebug("[CROSSPOSTCONTROL-PLUGIN] MessageWillBePosted Hook Start")
	post, rejectMessage := c.server.HandleMessage(post, true)
	c.API.LogDebug("[CROSSPOSTCONTROL-PLUGIN] MessageWillBePosted Hook End")
	return post, rejectMessage
}

func (c *BanMutePlugin) readConfig() error {
	c.config = &banmuteplugin.Config{}
	err := c.API.LoadPluginConfiguration(c.config)
	return err
}

func (c *BanMutePlugin) reloadConfig() error {
	if c.server != nil {
		c.server.ReloadConfig(*c.config)
	}

	return nil
}

func main() {
	plugin.ClientMain(&BanMutePlugin{})
}
