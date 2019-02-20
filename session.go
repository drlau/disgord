package disgord

import (
	"errors"
	"time"

	"github.com/andersfylling/disgord/logger"

	"github.com/andersfylling/snowflake/v3"

	"github.com/andersfylling/disgord/httd"
)

// NewSessionMock returns a session interface that triggers random events allows for fake rest requests.
// Ideal to test the behaviour of your new bot.
// Not implemented!
// TODO: what about a terminal interface for triggering specific events?
func NewSessionMock(conf *Config) (SessionMock, error) {
	return nil, errors.New("not implemented")
}

// EventChannels all methods for retrieving event channels
type EventChannels interface {
	Ready() <-chan *Ready
	Resumed() <-chan *Resumed
	ChannelCreate() <-chan *ChannelCreate
	ChannelUpdate() <-chan *ChannelUpdate
	ChannelDelete() <-chan *ChannelDelete
	ChannelPinsUpdate() <-chan *ChannelPinsUpdate
	GuildCreate() <-chan *GuildCreate
	GuildUpdate() <-chan *GuildUpdate
	GuildDelete() <-chan *GuildDelete
	GuildBanAdd() <-chan *GuildBanAdd
	GuildBanRemove() <-chan *GuildBanRemove
	GuildEmojisUpdate() <-chan *GuildEmojisUpdate
	GuildIntegrationsUpdate() <-chan *GuildIntegrationsUpdate
	GuildMemberAdd() <-chan *GuildMemberAdd
	GuildMemberRemove() <-chan *GuildMemberRemove
	GuildMemberUpdate() <-chan *GuildMemberUpdate
	GuildMembersChunk() <-chan *GuildMembersChunk
	GuildRoleUpdate() <-chan *GuildRoleUpdate
	GuildRoleCreate() <-chan *GuildRoleCreate
	GuildRoleDelete() <-chan *GuildRoleDelete
	MessageCreate() <-chan *MessageCreate
	MessageUpdate() <-chan *MessageUpdate
	MessageDelete() <-chan *MessageDelete
	MessageDeleteBulk() <-chan *MessageDeleteBulk
	MessageReactionAdd() <-chan *MessageReactionAdd
	MessageReactionRemove() <-chan *MessageReactionRemove
	MessageReactionRemoveAll() <-chan *MessageReactionRemoveAll
	PresenceUpdate() <-chan *PresenceUpdate
	PresencesReplace() <-chan *PresencesReplace
	TypingStart() <-chan *TypingStart
	UserUpdate() <-chan *UserUpdate
	VoiceStateUpdate() <-chan *VoiceStateUpdate
	VoiceServerUpdate() <-chan *VoiceServerUpdate
	WebhooksUpdate() <-chan *WebhooksUpdate
}

// Emitter for emitting data from A to B. Used in websocket connection
type Emitter interface {
	Emit(command SocketCommand, dataPointer interface{}) error
}

// Link is used to establish basic commands to create and destroy a link.
// See client.Disconnect() and client.Connect() for linking to the Discord servers
type Link interface {
	Connect() error
	Disconnect() error
}

// SocketHandler all socket related
type SocketHandler interface {
	// Link
	Disconnect() error

	// event handlers
	// inputs are in the following order: middlewares, handlers, controller
	On(event string, inputs ...interface{})
	Emitter

	// event channels
	EventChan(event string) (channel interface{}, err error)
	EventChannels() EventChannels

	// event register (which events to accept)
	// events which are not registered are discarded at socket level
	// to increase performance
	AcceptEvent(events ...string)
}

// AuditLogsRESTer REST interface for all audit-logs endpoints
type AuditLogsRESTer interface {
	GetGuildAuditLogs(guildID Snowflake, flags ...Flag) *guildAuditLogsBuilder
}

// ChannelRESTer REST interface for all Channel endpoints
type ChannelRESTer interface {
	GetChannel(id Snowflake, flags ...Flag) (ret *Channel, err error)
	ModifyChannel(id Snowflake, changes *ModifyChannelParams, flags ...Flag) (ret *Channel, err error)
	DeleteChannel(id Snowflake, flags ...Flag) (channel *Channel, err error)
	SetChannelPermissions(chanID, overwriteID Snowflake, params *SetChannelPermissionsParams, flags ...Flag) (err error)
	GetChannelInvites(id Snowflake, flags ...Flag) (ret []*Invite, err error)
	CreateChannelInvites(id Snowflake, params *CreateChannelInvitesParams, flags ...Flag) (ret *Invite, err error)
	DeleteChannelPermission(channelID, overwriteID Snowflake, flags ...Flag) (err error)
	TriggerTypingIndicator(channelID Snowflake, flags ...Flag) (err error)
	GetPinnedMessages(channelID Snowflake, flags ...Flag) (ret []*Message, err error)
	AddPinnedChannelMessage(channelID, msgID Snowflake, flags ...Flag) (err error)
	DeletePinnedChannelMessage(channelID, msgID Snowflake, flags ...Flag) (err error)
	GroupDMAddRecipient(channelID, userID Snowflake, params *GroupDMAddRecipientParams, flags ...Flag) (err error)
	GroupDMRemoveRecipient(channelID, userID Snowflake, flags ...Flag) (err error)
	GetMessages(channelID Snowflake, params URLQueryStringer, flags ...Flag) (ret []*Message, err error)
	GetMessage(channelID, messageID Snowflake, flags ...Flag) (ret *Message, err error)
	CreateMessage(channelID Snowflake, params *CreateMessageParams, flags ...Flag) (ret *Message, err error)
	EditMessage(chanID, msgID Snowflake, params *EditMessageParams, flags ...Flag) (ret *Message, err error)
	DeleteMessage(channelID, msgID Snowflake, flags ...Flag) (err error)
	BulkDeleteMessages(chanID Snowflake, params *BulkDeleteMessagesParams, flags ...Flag) (err error)
	CreateReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)
	DeleteOwnReaction(channelID, messageID Snowflake, emoji interface{}, flags ...Flag) (err error)
	DeleteUserReaction(channelID, messageID, userID Snowflake, emoji interface{}, flags ...Flag) (err error)
	GetReaction(channelID, messageID Snowflake, emoji interface{}, params URLQueryStringer, flags ...Flag) (ret []*User, err error)
	DeleteAllReactions(channelID, messageID Snowflake, flags ...Flag) (err error)
}

// EmojiRESTer REST interface for all emoji endpoints
type EmojiRESTer interface {
	GetGuildEmojis(id Snowflake, flags ...Flag) *getGuildEmojisBuilder
	GetGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) *getGuildEmojiBuilder
	CreateGuildEmoji(guildID Snowflake, name, image string, flags ...Flag) *createGuildEmojiBuilder
	ModifyGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) *modifyGuildEmojiBuilder
	DeleteGuildEmoji(guildID, emojiID Snowflake, flags ...Flag) *basicBuilder
}

// GuildRESTer REST interface for all guild endpoints
type GuildRESTer interface {
	CreateGuild(params *CreateGuildParams, flags ...Flag) (ret *Guild, err error)
	GetGuild(id Snowflake, flags ...Flag) (ret *Guild, err error)
	ModifyGuild(id Snowflake, params *ModifyGuildParams, flags ...Flag) (ret *Guild, err error)
	DeleteGuild(id Snowflake, flags ...Flag) (err error)
	GetGuildChannels(id Snowflake, flags ...Flag) (ret []*Channel, err error)
	CreateGuildChannel(id Snowflake, params *CreateGuildChannelParams, flags ...Flag) (ret *Channel, err error)
	ModifyGuildChannelPositions(id Snowflake, params []ModifyGuildChannelPositionsParams, flags ...Flag) (ret *Guild, err error)
	GetGuildMember(guildID, userID Snowflake, flags ...Flag) (ret *Member, err error)
	GetGuildMembers(guildID, after Snowflake, limit int, flags ...Flag) (ret []*Member, err error)
	AddGuildMember(guildID, userID Snowflake, params *AddGuildMemberParams, flags ...Flag) (ret *Member, err error)
	ModifyGuildMember(guildID, userID Snowflake, params *ModifyGuildMemberParams, flags ...Flag) (err error)
	ModifyCurrentUserNick(id Snowflake, params *ModifyCurrentUserNickParams, flags ...Flag) (nick string, err error)
	AddGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error)
	RemoveGuildMemberRole(guildID, userID, roleID Snowflake, flags ...Flag) (err error)
	RemoveGuildMember(guildID, userID Snowflake, flags ...Flag) (err error)
	GetGuildBans(id Snowflake, flags ...Flag) (ret []*Ban, err error)
	GetGuildBan(guildID, userID Snowflake, flags ...Flag) (ret *Ban, err error)
	CreateGuildBan(guildID, userID Snowflake, params *CreateGuildBanParams, flags ...Flag) (err error)
	RemoveGuildBan(guildID, userID Snowflake, flags ...Flag) (err error)
	GetGuildRoles(guildID Snowflake, flags ...Flag) (ret []*Role, err error)
	CreateGuildRole(id Snowflake, params *CreateGuildRoleParams, flags ...Flag) (ret *Role, err error)
	ModifyGuildRolePositions(guildID Snowflake, params []ModifyGuildRolePositionsParams, flags ...Flag) (ret []*Role, err error)
	ModifyGuildRole(guildID, roleID Snowflake, flags ...Flag) (builder *modifyGuildRoleBuilder)
	DeleteGuildRole(guildID, roleID Snowflake, flags ...Flag) (err error)
	GetGuildPruneCount(id Snowflake, params *GuildPruneParams, flags ...Flag) (ret *GuildPruneCount, err error)
	BeginGuildPrune(id Snowflake, params *GuildPruneParams, flags ...Flag) (ret *GuildPruneCount, err error)
	GetGuildVoiceRegions(id Snowflake, flags ...Flag) (ret []*VoiceRegion, err error)
	GetGuildInvites(id Snowflake, flags ...Flag) (ret []*Invite, err error)
	GetGuildIntegrations(id Snowflake, flags ...Flag) (ret []*Integration, err error)
	CreateGuildIntegration(guildID Snowflake, params *CreateGuildIntegrationParams, flags ...Flag) (err error)
	ModifyGuildIntegration(guildID, integrationID Snowflake, params *ModifyGuildIntegrationParams, flags ...Flag) (err error)
	DeleteGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error)
	SyncGuildIntegration(guildID, integrationID Snowflake, flags ...Flag) (err error)
	GetGuildEmbed(guildID Snowflake, flags ...Flag) (ret *GuildEmbed, err error)
	ModifyGuildEmbed(guildID Snowflake, params *GuildEmbed, flags ...Flag) (ret *GuildEmbed, err error)
	GetGuildVanityURL(guildID Snowflake, flags ...Flag) (ret *PartialInvite, err error)
}

// InviteRESTer REST interface for all invite endpoints
type InviteRESTer interface {
	GetInvite(inviteCode string, flags ...Flag) *getInviteBuilder
	DeleteInvite(inviteCode string, flags ...Flag) *deleteInviteBuilder
}

// UserRESTer REST interface for all user endpoints
type UserRESTer interface {
	GetCurrentUser(flags ...Flag) (*User, error)
	GetUser(id Snowflake, flags ...Flag) (*User, error)
	ModifyCurrentUser(flags ...Flag) (builder *modifyCurrentUserBuilder)
	GetCurrentUserGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) (ret []*PartialGuild, err error)
	LeaveGuild(id Snowflake, flags ...Flag) (err error)
	GetUserDMs(flags ...Flag) (ret []*Channel, err error)
	CreateDM(recipientID Snowflake, flags ...Flag) (ret *Channel, err error)
	CreateGroupDM(params *CreateGroupDMParams, flags ...Flag) (ret *Channel, err error)
	GetUserConnections(flags ...Flag) (ret []*UserConnection, err error)
}

// VoiceRESTer REST interface for all voice endpoints
type VoiceRESTer interface {
	GetVoiceRegions(flags ...Flag) *listVoiceRegionsBuilder
}

// WebhookRESTer REST interface for all Webhook endpoints
type WebhookRESTer interface {
	CreateWebhook(channelID Snowflake, params *CreateWebhookParams, flags ...Flag) (ret *Webhook, err error)
	GetChannelWebhooks(channelID Snowflake, flags ...Flag) (ret []*Webhook, err error)
	GetGuildWebhooks(guildID Snowflake, flags ...Flag) (ret []*Webhook, err error)
	GetWebhook(id Snowflake, flags ...Flag) (ret *Webhook, err error)
	GetWebhookWithToken(id Snowflake, token string, flags ...Flag) (ret *Webhook, err error)
	ModifyWebhook(id Snowflake, params *ModifyWebhookParams, flags ...Flag) (ret *Webhook, err error)
	ModifyWebhookWithToken(newWebhook *Webhook, flags ...Flag) (ret *Webhook, err error)
	DeleteWebhook(webhookID Snowflake, flags ...Flag) (err error)
	DeleteWebhookWithToken(id Snowflake, token string, flags ...Flag) (err error)
	ExecuteWebhook(params *ExecuteWebhookParams, wait bool, URLSuffix string, flags ...Flag) (err error)
	ExecuteSlackWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error)
	ExecuteGitHubWebhook(params *ExecuteWebhookParams, wait bool, flags ...Flag) (err error)
}

// RESTer holds all the sub REST interfaces
type RESTer interface {
	AuditLogsRESTer
	ChannelRESTer
	EmojiRESTer
	GuildRESTer
	InviteRESTer
	UserRESTer
	VoiceRESTer
	WebhookRESTer
}

// VoiceHandler holds all the voice connection related methods
type VoiceHandler interface {
	VoiceConnect(guildID, channelID Snowflake) (ret VoiceConnection, err error)
}

// Session Is the runtime interface for DisGord. It allows you to interact with a live session (using sockets or not).
// Note that this interface is used after you've configured DisGord, and therefore won't allow you to configure it
// further.
type Session interface {
	// give information about the bot/connected user
	Myself() (*User, error)

	// Request For interacting with Discord. Sending messages, creating channels, guilds, etc.
	// To read object state such as guilds, State() should be used in stead. However some data
	// might not exist in the state. If so it should be requested. Note that this only holds http
	// CRUD operation and not the actual rest endpoints for discord (See Rest()).
	// Deprecated: will be unexported in next breaking release
	Req() httd.Requester

	// Cache reflects the latest changes received from Discord gateway.
	// Should be used instead of requesting objects.
	// Deprecated: will be unexported in next breaking release
	Cache() Cacher

	Logger() logger.Logger

	// RateLimiter the rate limiter for the discord REST API
	// Deprecated: will be unexported in next breaking release
	RateLimiter() httd.RateLimiter

	// Discord Gateway, web socket
	SocketHandler
	HeartbeatLatency() (duration time.Duration, err error)

	// Generic CRUD operations for Discord interaction
	DeleteFromDiscord(obj discordDeleter) error
	SaveToDiscord(original discordSaver, changes ...discordSaver) error

	AddPermission(permission int) (updatedPermissions int)
	GetPermissions() (permissions int)
	CreateBotURL() (u string, err error)

	// state/caching module
	// checks the cacheLink first, otherwise do a http request
	RESTer

	// Custom REST functions
	SendMsg(channelID Snowflake, message *Message, flags ...Flag) (msg *Message, err error)
	SendMsgString(channelID Snowflake, content string, flags ...Flag) (msg *Message, err error)
	UpdateMessage(message *Message, flags ...Flag) (msg *Message, err error)
	UpdateChannel(channel *Channel, flags ...Flag) (err error)

	// Status update functions
	UpdateStatus(s *UpdateStatusCommand) (err error)
	UpdateStatusString(s string) (err error)

	GetGuilds(params *GetCurrentUserGuildsParams, flags ...Flag) ([]*Guild, error)
	GetConnectedGuilds() []snowflake.ID

	// same as above. Except these returns a channel
	// WARNING: none below should be assumed to be working.
	// TODO: implement in the future!
	//GuildChan(guildID Snowflake) <-chan *Guild
	//ChannelChan(channelID Snowflake) <-chan *Channel
	//ChannelsChan(guildID Snowflake) <-chan map[Snowflake]*Channel
	//MsgChan(msgID Snowflake) <-chan *Message
	//UserChan(userID Snowflake) <-chan *UserChan
	//MemberChan(guildID, userID Snowflake) <-chan *Member
	//MembersChan(guildID Snowflake) <-chan map[Snowflake]*Member

	// Voice handler, responsible for opening up new voice channel connections
	VoiceHandler
}

type SessionMock interface {
	Session
	// TODO: methods for triggering certain events and controlling states/tracking
}
