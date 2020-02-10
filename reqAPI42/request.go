package reqAPI42

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	cst "github.com/vsteffen/42_api/tools/constants"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func (api42 *API42) debugPrintRsp(rsp *http.Response) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("debugPrintRsp: Failed to read body response")
	}
	fmt.Println(string(bodyBytes))
}

func (api42 *API42) prepareGetParamURLReq(rawquery string) (*url.URL, *url.Values) {
	tmpURL, err := url.Parse(rawquery)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse this URL -> " + rawquery)
	}
	tmpParamURL := url.Values{}
	tmpParamURL.Add(cst.ReqToken, api42.keys.tokenAccess)
	return tmpURL, &tmpParamURL
}

func (api42 *API42) executeGetURLReq(getURL *url.URL) *http.Response {
	now := time.Now()
	waitTime := api42.rlLastReqSec.Add(time.Millisecond * cst.RLWaitTimeMsPerSecond)
	nowBeforeWait := now.Before(waitTime)
	if api42.rlNbReqSec >= cst.RLPerSecond || (nowBeforeWait && api42.rlNbReqSec >= cst.RLPerSecond) {
		time.Sleep(waitTime.Sub(now))
		api42.rlLastReqSec = time.Now()
		api42.rlNbReqSec = 1
	} else {
		api42.rlLastReqSec = now
		api42.rlNbReqSec++
		if nowBeforeWait == false {
			api42.rlNbReqSec = 1
		}
	}
	rsp, err := http.Get(getURL.String())
	pc, _, _, _ := runtime.Caller(1)
	callerFuncName := runtime.FuncForPC(pc).Name()
	callerFuncNameShort := callerFuncName[strings.LastIndex(callerFuncName, ".")+1:]
	if err != nil {
		log.Fatal().Err(err).Msg(callerFuncNameShort + ": Failed to GET " + getURL.String())
	}
	if rsp.StatusCode != http.StatusOK {
		defer rsp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatal().Err(err).Msg(callerFuncNameShort + ": Failed to read body response")
		}
		log.Error().Msg(callerFuncNameShort + ": invalid status code " + string(bodyBytes))
		return nil
	}
	return rsp
}

// GetCampus is used to execute a GET request for campus from 42's API (https://api.intra.42.fr/apidoc/2.0/campus.html)
func (api42 *API42) GetCampus(campusName string) *API42Campus {
	var err error

	campusURL, paramURL := api42.prepareGetParamURLReq(cst.CampusURL)
	paramURL.Add(cst.ReqFilter+"[name]", campusName)
	campusURL.RawQuery = paramURL.Encode()

	rsp := api42.executeGetURLReq(campusURL)
	if rsp == nil {
		return nil
	}
	defer rsp.Body.Close()

	rspJSON := make([]API42Campus, 0)
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("GetCampus: Failed to decode JSON values of campus")
	}
	if len(rspJSON) == 0 {
		log.Error().Msg("GetCampus: no campus found")
		return nil
	}
	log.Info().Msg("Found campus '" + campusName + "' ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	return &rspJSON[0]
}

// UpdateCampus is used to change campus values in API42 object
func (api42 *API42) UpdateCampus(campusName string) *API42Campus {
	campus := api42.GetCampus(campusName)
	if campus == nil {
		return nil
	}
	api42.campus = campus
	return campus
}

// GetLocations is used to execute a GET request for locations from 42's API (https://api.intra.42.fr/apidoc/2.0/locations.html)
func (api42 *API42) GetLocations() *[]API42Location {
	var err error

	locations := make([]API42Location, 0)
	locationsParsedURL := fmt.Sprintf(cst.LocationsURL, api42.campus.ID)
	locationsURL, paramURL := api42.prepareGetParamURLReq(locationsParsedURL)
	paramURL.Add(cst.ReqFilter+"[active]", "true")
	paramURL.Add(cst.ReqPageSize, cst.ReqPageSizeMax)

	for i := 1; ; i++ {
		pageNumberStr := strconv.Itoa(i)
		log.Info().Msg("GetLocations: GET page " + pageNumberStr + " ...")
		paramURL.Set(cst.ReqPage, pageNumberStr)
		locationsURL.RawQuery = paramURL.Encode()

		rsp := api42.executeGetURLReq(locationsURL)
		if rsp == nil {
			return nil
		}
		defer rsp.Body.Close()

		rspJSON := make([]API42Location, 0)
		decoder := json.NewDecoder(rsp.Body)
		if err = decoder.Decode(&rspJSON); err != nil {
			log.Fatal().Err(err).Msg("GetLocations: Failed to decode JSON values of location")
		}
		if len(rspJSON) == 0 {
			break
		}
		locations = append(locations, rspJSON...)
	}

	if len(locations) == 0 {
		log.Warn().Msg("GetLocations: no location found")
		return nil
	}
	log.Info().Msg("GetLocations: locations updated")
	return &locations
}

// func (api42 *API42) UpdateLocations() (*[]API42Location) {
// 	locations := api42.GetLocations()
// 	if locations == nil {
// 		return nil
// 	}
// 	api42.locations = locations
// 	return locations
// }

// GetCursus is used to execute a GET request for cursus from 42's API (https://api.intra.42.fr/apidoc/2.0/cursus.html)
func (api42 *API42) GetCursus(cursusName string) *API42Cursus {
	var err error

	cursusURL, paramURL := api42.prepareGetParamURLReq(cst.CursusURL)
	paramURL.Add(cst.ReqFilter+"[name]", cursusName)
	cursusURL.RawQuery = paramURL.Encode()

	rsp := api42.executeGetURLReq(cursusURL)
	if rsp == nil {
		return nil
	}
	defer rsp.Body.Close()

	rspJSON := make([]API42Cursus, 0)
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("GetCursus: Failed to decode JSON values of cursus")
	}
	if len(rspJSON) == 0 {
		log.Error().Msg("GetCursus: no cursus found")
		return nil
	}
	log.Info().Msg("Found cursus '" + cursusName + "' ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	return &rspJSON[0]
}

// UpdateCursus is used to change cursus values in API42 object
func (api42 *API42) UpdateCursus(cursusName string) *API42Cursus {
	cursus := api42.GetCursus(cursusName)
	if cursus == nil {
		return nil
	}
	api42.cursus = cursus
	return cursus
}

// GetProjects is used to execute a GET request for projects from 42's API (https://api.intra.42.fr/apidoc/2.0/projects.html)
func (api42 *API42) GetProjects() *[]API42Project {
	var err error

	projects := make([]API42Project, 0)
	projectsURL, paramURL := api42.prepareGetParamURLReq(cst.ProjectsURL)
	paramURL.Add("cursus_id", strconv.FormatUint(uint64(api42.cursus.ID), 10))
	paramURL.Add(cst.ReqFilter+"[has_git]", "true")
	paramURL.Add(cst.ReqFilter+"[has_mark]", "true")
	paramURL.Add(cst.ReqFilter+"[visible]", "true")
	paramURL.Add(cst.ReqFilter+"[exam]", "false")
	paramURL.Add(cst.ReqPageSize, cst.ReqPageSizeMax)

	for i := 1; ; i++ {
		pageNumberStr := strconv.Itoa(i)
		log.Info().Msg("GetProjects: GET page " + pageNumberStr + " ...")
		paramURL.Set(cst.ReqPage, pageNumberStr)
		projectsURL.RawQuery = paramURL.Encode()

		rsp := api42.executeGetURLReq(projectsURL)
		if rsp == nil {
			return nil
		}
		defer rsp.Body.Close()

		rspJSON := make([]API42Project, 0)
		decoder := json.NewDecoder(rsp.Body)
		if err = decoder.Decode(&rspJSON); err != nil {
			log.Fatal().Err(err).Msg("GetProjects: Failed to decode JSON values of a project")
		}
		if len(rspJSON) == 0 {
			break
		}
		i := 0
		for _, project := range rspJSON {
			for _, campus := range project.Campus {
				if campus.ID == api42.campus.ID {
					project.Campus = nil
					rspJSON[i] = project
					i++
					break
				}
			}
		}
		rspJSON = rspJSON[:i]
		// fmt.Println(rspJSON)
		projects = append(projects, rspJSON...)
	}

	if len(projects) == 0 {
		log.Fatal().Msg("GetProjects: no project found")
		return nil
	}
	log.Info().Msg("GetProjects: projects updated")
	return &projects
}

// func (api42 *API42) UpdateProjects() (*[]API42Project) {
// 	projects := api42.GetProjects()
// 	if projects == nil {
// 		return nil
// 	}
// 	api42.projects = projects
// 	return projects
// }

// GetUsersOfProjectsUsers is used to execute a GET request for projects users from 42's API (https://api.intra.42.fr/apidoc/2.0/projects_users.html)
func (api42 *API42) GetUsersOfProjectsUsers(projectID uint) *[]API42ProjectUser {
	var err error

	projectIDStr := strconv.FormatUint(uint64(projectID), 10)
	log.Info().Msg("GetUsersOfProjectsUsers: searching with project ID = " + projectIDStr + " ...")
	projectsUserURL, paramURL := api42.prepareGetParamURLReq(cst.ProjectsUsersURL)
	paramURL.Add(cst.ReqFilter+"[project_id]", projectIDStr)
	paramURL.Add(cst.ReqFilter+"[cursus]", strconv.FormatUint(uint64(api42.cursus.ID), 10))
	paramURL.Add(cst.ReqFilter+"[campus]", strconv.FormatUint(uint64(api42.campus.ID), 10))
	paramURL.Add(cst.ReqFilter+"[marked]", "true")
	paramURL.Add(cst.ReqPageSize, cst.ReqPageSizeMax)
	projectsUserURL.RawQuery = paramURL.Encode()

	projectsUsers := make([]API42ProjectUser, 0)
	for i := 1; ; i++ {
		pageNumberStr := strconv.Itoa(i)
		log.Info().Msg("GetUsersOfProjectsUsers: GET page " + pageNumberStr + " ...")
		paramURL.Set(cst.ReqPage, pageNumberStr)
		projectsUserURL.RawQuery = paramURL.Encode()

		rsp := api42.executeGetURLReq(projectsUserURL)
		if rsp == nil {
			return nil
		}
		defer rsp.Body.Close()

		rspJSON := make([]API42ProjectUser, 0)
		decoder := json.NewDecoder(rsp.Body)
		if err = decoder.Decode(&rspJSON); err != nil {
			log.Fatal().Err(err).Msg("GetUsersOfProjectsUsers: Failed to decode JSON values of a project user")
		}
		if len(rspJSON) == 0 {
			break
		}
		projectsUsers = append(projectsUsers, rspJSON...)
	}

	if len(projectsUsers) == 0 {
		log.Error().Msg("GetUsersOfProjectsUsers: no project user found")
		return nil
	}
	log.Info().Msg("GetUsersOfProjectsUsers: Get all projects users")
	return &projectsUsers
}

// GetMe is used to execute a GET request for "me" user from 42's API (https://api.intra.42.fr/apidoc/2.0/users/me.html)
func (api42 *API42) GetMe() *API42User {
	var err error

	meURL, paramURL := api42.prepareGetParamURLReq(cst.MeURL)
	meURL.RawQuery = paramURL.Encode()

	rsp := api42.executeGetURLReq(meURL)
	if rsp == nil {
		return nil
	}
	defer rsp.Body.Close()

	var rspJSON API42User
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("GetMe: Failed to decode JSON values of me")
	}
	log.Info().Msg("GetMe: login -> " + rspJSON.Login)
	return &rspJSON
}

// New create a new API42 object
func New(flags []interface{}) *API42 {
	tmp := API42{}
	tmp.initKeys()
	tmp.rlLastReqSec = time.Now()

	if *flags[0].(*bool) {
		tmp.RefreshToken()
	}

	if *flags[1].(*bool) {
		if tmp.UpdateCampus(cst.DefaultCampusName) == nil || tmp.UpdateCursus(cst.DefaultCursusName) == nil {
			log.Fatal().Msg("API42.New: failed to initialize API42")
		}
	} else {
		log.Info().Msg("API42.New: Use default values for campus and cursus")
		tmp.campus = &API42Campus{cst.DefaultCampusID, cst.DefaultCampusName}
		tmp.cursus = &API42Cursus{cst.DefaultCursusID, cst.DefaultCursusName}
	}
	return &tmp
}
