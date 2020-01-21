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
)

// API42 type which use to interact with 42's API
type API42 struct {
	keys   apiKeys
	campus uint
	cursus uint
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

func (api42 *API42) UpdateCampusID(campusName string) (bool) {
	var err error

	type campusRsp struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	log.Info().Msg("Search campus ID ...")
	campusURL, paramURL := api42.prepareGetParamURLReq(cst.CampusURL)
	paramURL.Add(cst.ReqFilter + "[name]", campusName)
	campusURL.RawQuery = paramURL.Encode()

	rsp, err := http.Get(campusURL.String())
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
	log.Info().Msg("Found campus Paris ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	api42.campus = rspJSON[0].ID
	return true
}

func (api42 *API42) UpdateLocations() {
	// var err error

	log.Info().Msg("Updating locations ...")
	locationsRealURLStr := fmt.Sprintf(cst.LocationsURL, api42.campus)
	log.Info().Msg(locationsRealURLStr)
	locationsURL, paramURL := api42.prepareGetParamURLReq(locationsRealURLStr)
	locationsURL.RawQuery = paramURL.Encode()
	//Not finished
}

func (api42 *API42) UpdateCursusID(cursusName string) (bool) {
	var err error

	type cursusRsp struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	log.Info().Msg("Search cursus ID ...")
	cursusURL, paramURL := api42.prepareGetParamURLReq(cst.CursusURL)
	paramURL.Add(cst.ReqFilter + "[name]", cursusName)
	cursusURL.RawQuery = paramURL.Encode()

	rsp, err := http.Get(cursusURL.String())
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

func (api42 *API42) GetUsersAvailable() {
	
}

// New create new API42 obj
func New(args ...bool) *API42 {
	tmp := API42{}
	tmp.initKeys()
	nbArg := len(args)
	if nbArg > cst.API42ArgRefresh && args[cst.API42ArgRefresh] {
		tmp.RefreshToken()
	}
	if !tmp.UpdateCampusID(cst.CampusName) || !tmp.UpdateCursusID(cst.CursusName) {
		log.Fatal().Msg("API42.New: failed to initialize API42")
	}
	return &tmp
}
