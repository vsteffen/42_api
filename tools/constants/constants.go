package constants

const (
	ProjectName = "Find_Examiner"
	UID = "698c6cc85d07ed45f147486e6630ae9dedad619943999ce345fb958c4867bec9"
	BaseURL = "https://api.intra.42.fr"

	SizeApiKeys = 64

	API42ArgRefresh	= 0

	AuthURL            = BaseURL + "/oauth/authorize"
	AuthVarClt         = "client_id"
	AuthVarRedirectURI = "redirect_uri"
	AuthValRedirectURI = "https://profile.intra.42.fr"
	AuthVarRespType    = "response_type"
	AuthValRespType    = "code"

	TokenURL              = BaseURL + "/oauth/token"
	TokenReqGrantAuthCode = "authorization_code"
	TokenReqGrantRefresh  = "refresh_token"
	TokenReqRedirectURI   = "https://profile.intra.42.fr"

	PathSecret       = "./secret"
	PathTokenAccess  = "./token_access"
	PathTokenRefresh = "./token_refresh"

	ReqToken  = "access_token"
	ReqFilter = "filter"

	CampusName = "Paris"

	LocationsURL = BaseURL + "/v2/campus/%d/locations "

	CampusURL = BaseURL + "/v2/campus"

	CursusURL	= BaseURL + "/v2/cursus"
	CursusName	= "42cursus"
)
