package constants

const (
	ActionUserCreated       = "user.created"
	ActionUserUpdated       = "user.updated"
	ActionUserDeleted       = "user.deleted"
	ActionUserListViewed    = "user.list.viewed"
	ActionUserProfileViewed = "user.profile.viewed"
	ActionUserRoleChanged   = "user.role.changed"
)

const (
	ActionAuthLogin          = "auth.login"
	ActionAuthLogout         = "auth.logout"
	ActionAuthRefresh        = "auth.refresh"
	ActionAuthRegister       = "auth.register"
	ActionAuthValidate       = "auth.validate"
	ActionAuthBasicLogin     = "auth.basic_login"
	ActionAuthJWTLogin       = "auth.jwt_login"
	ActionAuthSessionLogin   = "auth.session_login"
	ActionAuthAPIKeyGenerate = "auth.apikey.generate"
)

const (
	ResourceUser  = "user"
	ResourceUsers = "users"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

const (
	RoutingKeyUserCreated       = "user.created"
	RoutingKeyUserUpdated       = "user.updated"
	RoutingKeyUserDeleted       = "user.deleted"
	RoutingKeyUserListViewed    = "user.list.viewed"
	RoutingKeyUserProfileViewed = "user.profile.viewed"
	RoutingKeyAuthPrefix        = "auth."
	RoutingKeyAPIRequest        = "api.request"
	RoutingKeyUserNotification  = "notification.send"
	RoutingKeyTokenValidation   = "token.validation"
)

const (
	PatternAuthEvents  = "auth.#"
	PatternTokenEvents = "token.#"
	PatternAPIEvents   = "api.#"
	PatternUserEvents  = "user.#"
)

const (
	RabbitMQRetryInterval = 5
	RabbitMQRetryAttempts = 5
	RabbitMQPublishTimeout = 5
)

const (
	ContentTypeJSON = "application/json"
)

const (
	ExchangeName = "events_exchange"
	ExchangeType = "topic"

	QueueAuthEvents = "auth_events_queue"
	QueueAPIEvents  = "gateway_events_queue"
	QueueUserEvents = "user_events_queue"
)
