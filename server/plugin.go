package main

import (
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-plugin-starter-template/server/command"
	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// kvstore is the client used to read/write KV records for this plugin.
	kvstore kvstore.KVStore

	// client is the Mattermost server API client.
	client *pluginapi.Client

	// commandClient is the client used to register and execute slash commands.
	commandClient command.Command

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// OnActivate is invoked when the plugin is activated. If an error is returned, the plugin will be deactivated.
func (p *Plugin) OnActivate() error {
	p.client = pluginapi.NewClient(p.API, p.Driver)

	p.kvstore = kvstore.NewKVStore(p.client)

	p.commandClient = command.NewCommandHandler(p.client)

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated.
func (p *Plugin) OnDeactivate() error {
	return nil
}

// MessageWillBePosted is invoked when a message is posted by a user before it is committed to the database.
// This hook replaces emoji shortcuts with their corresponding Mattermost emoji codes.
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	if post == nil || post.Message == "" {
		return post, ""
	}

	// Process the message to replace emoji shortcuts
	originalMessage := post.Message
	processedMessage := ProcessMessageForEmoji(originalMessage)

	// Only update if the message was modified
	if processedMessage != originalMessage {
		post.Message = processedMessage
	}

	return post, ""
}

// MessageWillBeUpdated is invoked when a message is updated by a user before it is committed to the database.
// This hook also replaces emoji shortcuts when editing messages.
func (p *Plugin) MessageWillBeUpdated(c *plugin.Context, newPost *model.Post, oldPost *model.Post) (*model.Post, string) {
	if newPost == nil || newPost.Message == "" {
		return newPost, ""
	}

	// Process the message to replace emoji shortcuts
	originalMessage := newPost.Message
	processedMessage := ProcessMessageForEmoji(originalMessage)

	// Only update if the message was modified
	if processedMessage != originalMessage {
		newPost.Message = processedMessage
	}

	return newPost, ""
}

// This will execute the commands that were registered in the NewCommandHandler function.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	response, err := p.commandClient.Handle(args)
	if err != nil {
		return nil, model.NewAppError("ExecuteCommand", "plugin.command.execute_command.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return response, nil
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
