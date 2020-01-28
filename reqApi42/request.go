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
	"time"
)

// API42 type which use to interact with 42's API
type API42 struct {
	keys   	apiKeys
	rlLastReqSec	time.Time
	rlNbReqSec	uint
	campus	uint
	cursus	uint
}

type API42Project struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
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

func (api42 *API42) executeGetURLReq(getURL *url.URL) (*http.Response, error) {
	now := time.Now()
	waitTime := api42.rlLastReqSec.Add(time.Millisecond * cst.RLWaitTimeMsPerSecond)
	nowBeforeWait := now.Before(waitTime)
	if api42.rlNbReqSec >= cst.RLPerSecond || (nowBeforeWait && api42.rlNbReqSec >= cst.RLPerSecond) {
		time.Sleep(waitTime.Sub(now))
		api42.rlLastReqSec = time.Now()
		api42.rlNbReqSec = 1
		fmt.Println("TOO EARLY")
	} else {
		api42.rlLastReqSec = now
		api42.rlNbReqSec += 1
		if nowBeforeWait == false {
			api42.rlNbReqSec = 1
		}
		fmt.Println("GOOD TIMING")
	}
	rsp, err := http.Get(getURL.String())
	return rsp, err
}

func (api42 *API42) UpdateCampusID(campusName string) (bool) {
	var err error

	type campusRsp struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	campusURL, paramURL := api42.prepareGetParamURLReq(cst.CampusURL)
	paramURL.Add(cst.ReqFilter + "[name]", campusName)
	campusURL.RawQuery = paramURL.Encode()

	rsp, err := api42.executeGetURLReq(campusURL)
	if err != nil {
		log.Error().Err(err).Msg("UpdateCampusID: Failed to GET " + campusURL.String())
		return false
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatal().Err(err).Msg("UpdateCampusID: Failed to read body response")
		}
		log.Error().Msg("UpdateCampusID: invalid status code " + string(bodyBytes))
		return false
	}
	rspJSON := make([]campusRsp, 0)
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("UpdateCampusID: Failed to decode JSON values of campus")
	}
	if (len(rspJSON) == 0) {
		log.Error().Msg("UpdateCampusID: no campus found")
		return false
	}
	log.Info().Msg("Found campus '" + cst.CampusName + "' ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	api42.campus = rspJSON[0].ID
	return true
}

func (api42 *API42) UpdateLocations() {
	// var err error

	log.Info().Msg("Updating locations ...")
	locationsParsedURL := fmt.Sprintf(cst.LocationsURL, api42.campus)
	log.Info().Msg(locationsParsedURL)
	locationsURL, paramURL := api42.prepareGetParamURLReq(locationsParsedURL)
	locationsURL.RawQuery = paramURL.Encode()

	rsp, err := api42.executeGetURLReq(locationsURL)
	_ = rsp
	_ = err
}

func (api42 *API42) UpdateCursusID(cursusName string) (bool) {
	// var err error

	type cursusRsp struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	cursusURL, paramURL := api42.prepareGetParamURLReq(cst.CursusURL)
	paramURL.Add(cst.ReqFilter + "[name]", cursusName)
	cursusURL.RawQuery = paramURL.Encode()

	rsp, err := api42.executeGetURLReq(cursusURL)
	if err != nil {
		log.Fatal().Err(err).Msg("UpdateCursusID: Failed to GET " + cursusURL.String())
		return false
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatal().Err(err).Msg("UpdateCursusID: Failed to read body response")
		}
		log.Error().Msg("UpdateCursusID: invalid status code " + string(bodyBytes))
		return false
	}
	rspJSON := make([]cursusRsp, 0)
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("UpdateCursusID: Failed to decode JSON values of cursus")
	}
	if (len(rspJSON) == 0) {
		log.Error().Msg("UpdateCursusID: no cursus found")
		return false
	}
	log.Info().Msg("Found cursus '" + cursusName + "' ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	api42.cursus = rspJSON[0].ID
	return true
}

func (api42 *API42) GetCursusProjects() ([]API42Project, bool) {
	var err error

	cursusProjectsParsedURL := fmt.Sprintf(cst.CursusProjectsURL, api42.cursus)
	cursusProjectsURL, paramURL := api42.prepareGetParamURLReq(cursusProjectsParsedURL)
	// paramURL.Add(cst.ReqFilter + "[name]", cursusName)
	cursusProjectsURL.RawQuery = paramURL.Encode()

	rsp, err := api42.executeGetURLReq(cursusProjectsURL)
	if err != nil {
		log.Fatal().Err(err).Msg("UpdateProjectsID: Failed to GET " + cursusProjectsURL.String())
		return nil, false
	}
	defer rsp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(bodyBytes))

	// if rsp.StatusCode != http.StatusOK {
	// 	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	// 	if err != nil {
	// 		log.Fatal().Err(err).Msg("UpdateCursusID: Failed to read body response")
	// 	}
	// 	log.Error().Msg("UpdateCursusID: invalid status code " + string(bodyBytes))
	// 	return nil, false
	// }
	// rspJSON := make([]cursusRsp, 0)
	// decoder := json.NewDecoder(rsp.Body)
	// if err = decoder.Decode(&rspJSON); err != nil {
	// 	log.Fatal().Err(err).Msg("UpdateCursusID: Failed to decode JSON values of cursus")
	// }
	// if (len(rspJSON) == 0) {
	// 	log.Error().Msg("UpdateCursusID: no cursus found")
	// 	return nil, false
	// }
	// log.Info().Msg("Found cursus '" + cursusName + "' ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	return nil, true
}

func (api42 *API42) GetUsersAvailable() {

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
		if !tmp.UpdateCampusID(cst.CampusName) || !tmp.UpdateCursusID(cst.CursusName) {
			log.Fatal().Msg("API42.New: failed to initialize API42")
		}
	} else {
		log.Info().Msg("API42.New: Use default values")
		tmp.campus = cst.DefaultCampus
		tmp.cursus = cst.DefaultCursus
	}
	return &tmp
}
