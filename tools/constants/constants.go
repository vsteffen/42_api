package constants

const (
	ProjectName = "Find_Examiner"
	UID         = "698c6cc85d07ed45f147486e6630ae9dedad619943999ce345fb958c4867bec9"
	BaseURL     = "https://api.intra.42.fr"

	DefaultCampusID   = 1
	DefaultCampusName = "Paris"
	DefaultCursusID   = 21
	DefaultCursusName = "42cursus"

	SizeAPIKeys = 64

	// API42ArgRefresh		= 0
	// API42ArgCheckDefaultValues	= 1

	RLPerSecond           = 2
	RLWaitTimeMsPerSecond = 1050

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

	ReqToken       = "access_token"
	ReqFilter      = "filter"
	ReqPage        = "page[number]"
	ReqPageSize    = "page[size]"
	ReqPageSizeMax = "100"

	LocationsURL     = BaseURL + "/v2/campus/%d/locations"
	CampusURL        = BaseURL + "/v2/campus"
	CursusURL        = BaseURL + "/v2/cursus"
	ProjectsURL      = BaseURL + "/v2/me/projects"
	ProjectsUsersURL = BaseURL + "/v2/projects_users"

	MenuHello = "Find_Examiner v1.0\n"

	MenuActionFind            = "Find examiner    (use projects values)"
	MenuActionUpdateLocations = "Update locations (use cursus and campus values)"
	MenuActionUpdateProjects  = "Update projects  (use cursus and campus values)"
	MenuActionUpdateCursus    = "Update cursus"
	MenuActionUpdateCampus    = "Update campus"
	MenuActionRefreshTokens   = "Refresh API tokens"
	MenuActionQuit            = "Quit"

	FindNameMaxResults = 3

	ProfileURL     = "https://profile.intra.42.fr"
	ProfileUserURL = ProfileURL + "/users/"

	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
)
