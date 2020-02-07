package reqApi42

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	cst "github.com/vsteffen/42_api/tools/constants"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"runtime"
)

// API42 type which use to interact with 42's API
type API42 struct {
	keys		apiKeys
	rlLastReqSec	time.Time
	rlNbReqSec	uint
	campus		uint
	cursus		uint
	locations	[]API42Location
	projects	[]API42Project
}

type JSONProjectParent struct {
	API42ProjectParent
}

type JSONTime struct {
	time.Time
}

type API42Location struct {
	ID	uint	`json:"id"`
	EndAt	JSONTime `json:"end_at"`
	Host	string	`json:"host"`
	User struct {
		ID	uint	`json:"id"`
		Login	string	`json:"login"`
	} `json:"user"`
}

type API42Project struct {
	ID	uint			`json:"id"`
	Name	string			`json:"name"`
	Parent	*API42ProjectParent	`json:"parent"`
	Campus []struct {
		ID	uint	`json:"id"`
	} `json:"campus"`
}

type API42ProjectParent struct {
	ID	uint			`json:"id"`
	Name	string			`json:"name"`
	Slug	string			`json:"slug"`
}

func (jsonVal *JSONTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		*jsonVal = JSONTime{time.Time{}}
		return nil
	}
	timeFormated := strings.Trim(str, "\"")
	timeVal, err := time.Parse(time.RFC3339, timeFormated)
	if err != nil {
		return err
	}
	*jsonVal = JSONTime{timeVal}
	return nil
}

func (jsonVal *JSONProjectParent) UnmarshalJSON(b []byte) error {
	str := string(b)

	if str == "null" {
		jsonVal = nil
		return nil
	}

	var projectParent API42ProjectParent
	err := json.Unmarshal(b, &projectParent)
	if err != nil {
		return err
	}
	jsonVal = &JSONProjectParent{projectParent}
	return nil
}

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

func (api42 *API42) executeGetURLReq(getURL *url.URL) (*http.Response) {
	now := time.Now()
	waitTime := api42.rlLastReqSec.Add(time.Millisecond * cst.RLWaitTimeMsPerSecond)
	nowBeforeWait := now.Before(waitTime)
	if api42.rlNbReqSec >= cst.RLPerSecond || (nowBeforeWait && api42.rlNbReqSec >= cst.RLPerSecond) {
		time.Sleep(waitTime.Sub(now))
		api42.rlLastReqSec = time.Now()
		api42.rlNbReqSec = 1
	} else {
		api42.rlLastReqSec = now
		api42.rlNbReqSec += 1
		if nowBeforeWait == false {
			api42.rlNbReqSec = 1
		}
	}
	rsp, err := http.Get(getURL.String())
	pc, _, _, _ := runtime.Caller(1)
	callerFuncName := runtime.FuncForPC(pc).Name()
	callerFuncNameShort := callerFuncName[strings.LastIndex(callerFuncName, ".") + 1 : ]
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

func (api42 *API42) UpdateCampus(campusName string) (bool) {
	var err error

	type campusRsp struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	campusURL, paramURL := api42.prepareGetParamURLReq(cst.CampusURL)
	paramURL.Add(cst.ReqFilter + "[name]", campusName)
	campusURL.RawQuery = paramURL.Encode()

	rsp := api42.executeGetURLReq(campusURL)
	if rsp == nil {
		return false
	}
	defer rsp.Body.Close()

	rspJSON := make([]campusRsp, 0)
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("UpdateCampus: Failed to decode JSON values of campus")
	}
	if (len(rspJSON) == 0) {
		log.Error().Msg("UpdateCampus: no campus found")
		return false
	}
	log.Info().Msg("Found campus '" + cst.CampusName + "' ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	api42.campus = rspJSON[0].ID
	return true
}

func (api42 *API42) UpdateLocations() (bool) {
	var err error

	log.Info().Msg("Updating locations ...")
	api42.locations = nil
	locationsParsedURL := fmt.Sprintf(cst.LocationsURL, api42.campus)
	locationsURL, paramURL := api42.prepareGetParamURLReq(locationsParsedURL)
	paramURL.Add(cst.ReqFilter + "[active]", "true")
	paramURL.Add(cst.ReqPageSize, cst.ReqPageSizeMax)

	for i := 1; ; i++ {
		pageNumberStr := strconv.Itoa(i)
		log.Info().Msg("UpdateLocations: GET page " + pageNumberStr + " ...")
		paramURL.Set(cst.ReqPage, pageNumberStr)
		locationsURL.RawQuery = paramURL.Encode()

		rsp := api42.executeGetURLReq(locationsURL)
		if rsp == nil {
			return false
		}
		defer rsp.Body.Close()

		rspJSON := make([]API42Location, 0)
		decoder := json.NewDecoder(rsp.Body)
		if err = decoder.Decode(&rspJSON); err != nil {
			log.Fatal().Err(err).Msg("UpdateLocations: Failed to decode JSON values of location")
		}
		if (len(rspJSON) == 0) {
			break
		}
		api42.locations = append(api42.locations, rspJSON...)
	}

	if (len(api42.locations) == 0) {
		log.Warn().Msg("UpdateLocations: no location found")
		return false
	}
	log.Info().Msg("UpdateLocations: locations updated")
	return true
}

func (api42 *API42) UpdateCursus(cursusName string) (bool) {
	var err error

	type cursusRsp struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	cursusURL, paramURL := api42.prepareGetParamURLReq(cst.CursusURL)
	paramURL.Add(cst.ReqFilter + "[name]", cursusName)
	cursusURL.RawQuery = paramURL.Encode()

	rsp := api42.executeGetURLReq(cursusURL)
	if rsp == nil {
		return false
	}
	defer rsp.Body.Close()

	rspJSON := make([]cursusRsp, 0)
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("UpdateCursus: Failed to decode JSON values of cursus")
	}
	if (len(rspJSON) == 0) {
		log.Error().Msg("UpdateCursus: no cursus found")
		return false
	}
	log.Info().Msg("Found cursus '" + cursusName + "' ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	api42.cursus = rspJSON[0].ID
	return true
}

func (api42 *API42) UpdateProjects() (bool) {
	var err error

	log.Info().Msg("Updating projects ...")
	projectsURL, paramURL := api42.prepareGetParamURLReq(cst.ProjectsURL)
	paramURL.Add("cursus_id", strconv.FormatUint(uint64(api42.campus), 10))
	paramURL.Add(cst.ReqFilter + "[has_git]", "true")
	paramURL.Add(cst.ReqFilter + "[has_mark]", "true")
	paramURL.Add(cst.ReqFilter + "[visible]", "true")
	paramURL.Add(cst.ReqFilter + "[exam]", "false")
	paramURL.Add(cst.ReqPageSize, cst.ReqPageSizeMax)

	for i := 1; i < 2 /*To remove*/; i++ {
		pageNumberStr := strconv.Itoa(i)
		log.Info().Msg("UpdateProjects: GET page " + pageNumberStr + " ...")
		paramURL.Set(cst.ReqPage, pageNumberStr)
		projectsURL.RawQuery = paramURL.Encode()

		rsp := api42.executeGetURLReq(projectsURL)
		if rsp == nil {
			return false
		}
		defer rsp.Body.Close()

		rspJSON := make([]API42Project, 0)
		decoder := json.NewDecoder(rsp.Body)
		if err = decoder.Decode(&rspJSON); err != nil {
			log.Fatal().Err(err).Msg("UpdateProjects: Failed to decode JSON values of a project")
		}
		if (len(rspJSON) == 0) {
			break
		}
		i := 0
		for _, project := range rspJSON {
			for _, campus := range project.Campus {
				if campus.ID == api42.campus {
					project.Campus = nil
					rspJSON[i] = project
					i++
					break
				}
			}
		}
		rspJSON = rspJSON[:i]
		// fmt.Println(rspJSON)
		api42.projects = append(api42.projects, rspJSON...)
	}

	if (len(api42.projects) == 0) {
		log.Fatal().Msg("UpdateProjects: no project found")
		return false
	}
	log.Info().Msg("UpdateProjects: projects updated")
	return true
}

func (api42 *API42) GetProjects() (*[]API42Project) {
	return &api42.projects
}

// New create new API42 obj
func New(flags []interface{}) *API42 {
	tmp := API42{}
	tmp.initKeys()
	tmp.rlLastReqSec = time.Now()

	if *flags[0].(*bool) {
		tmp.RefreshToken()
	}

	if *flags[1].(*bool) {
		if !tmp.UpdateCampus(cst.CampusName) || !tmp.UpdateCursus(cst.CursusName) {
			log.Fatal().Msg("API42.New: failed to initialize API42")
		}
	} else {
		log.Info().Msg("API42.New: Use default values")
		tmp.campus = cst.DefaultCampus
		tmp.cursus = cst.DefaultCursus
	}
	tmp.UpdateProjects()
	// tmp.UpdateLocations()
	return &tmp
}
