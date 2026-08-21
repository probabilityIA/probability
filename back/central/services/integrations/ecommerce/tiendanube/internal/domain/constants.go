package domain

const (
	TypeID = uint(17)

	OAuthAuthorizeURLTemplate = "https://www.tiendanube.com/apps/%s/authorize"
	OAuthTokenURL             = "https://www.tiendanube.com/apps/authorize/token"

	AuthMethodOAuth = "oauth"

	ConfigAuthMethod = "auth_method"
	ConfigScope      = "scope"
)
