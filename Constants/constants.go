package constants

const (
	ProjectName = "Find_Examiner"

	BaseURL			= "https://api.intra.42.fr"

	AuthURL					= BaseURL + "/oauth/authorize"
	AuthVarClt				= "client_id"
	AuthVarRedirectURI		= "redirect_uri"
	AuthValRedirectURI		= "https://profile.intra.42.fr"
	AuthVarRespType			= "response_type"
	AuthValRespType			= "code"

	TokenURL		= BaseURL + "/oauth/token"

	PathSecret	= "./secret"
	PathToken	= "./token"

	UID = "698c6cc85d07ed45f147486e6630ae9dedad619943999ce345fb958c4867bec9"

	SizeKeys = 64
)
