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

func (api42 *API42) GetCampusID(campusName string) {
	var err error

	type campusReq struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	log.Info().Msg("Search campus ID ...")
	campusURL, paramURL := api42.prepareGetParamURLReq(cst.CampusURL)
	paramURL.Add(cst.ReqFilter+"[name]", campusName)
	campusURL.RawQuery = paramURL.Encode()

	rsp, err := http.Get(campusURL.String())
	if err != nil {
		log.Fatal().Err(err).Msg("initCampus: Failed to GET " + campusURL.String())
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatal().Err(err).Msg("initCampus: Failed to read body response")
		}
		log.Fatal().Msg("initCampus: invalid status code " + string(bodyBytes))
	}
	rspJSON := make([]campusReq, 0)
	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("initCampus: Failed to decode JSON values of campus")
	}
	log.Info().Msg("Found campus Paris ID -> " + strconv.FormatUint(uint64(rspJSON[0].ID), 10))
	api42.campus = rspJSON[0].ID
}

func (api42 *API42) UpdateLocations() {
	// var err error

	log.Info().Msg("Updating locations ...")
	locationsRealURLStr := fmt.Sprintf(cst.LocationsURL, api42.campus)
	log.Info().Msg(locationsRealURLStr)
	locationsURL, paramURL := api42.prepareGetParamURLReq(locationsRealURLStr)
	locationsURL.RawQuery = paramURL.Encode()
}

// New create new API42 obj
func New() *API42 {
	tmp := API42{}
	tmp.initKeys()
	return &tmp
}
