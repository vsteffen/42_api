package reqAPI42

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vsteffen/42_api/tools"
	cst "github.com/vsteffen/42_api/tools/constants"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type apiKeys struct {
	uid          string
	secret       string
	tokenAccess  string
	tokenRefresh string
}

type tokenReqNew struct {
	TokenGrant     string `json:"grant_type"`
	TokenCltID     string `json:"client_id"`
	TokenCltSecret string `json:"client_secret"`
	TokenCode      string `json:"code"`
	TokenRedirect  string `json:"redirect_uri"`
}

type tokenReqRefresh struct {
	TokenGrant     string `json:"grant_type"`
	TokenCltID     string `json:"client_id"`
	TokenCltSecret string `json:"client_secret"`
	TokenRefresh   string `json:"refresh_token"`
}

type tokenRsp struct {
	TokenAccess  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	TokenExp     uint64 `json:"expires_in"`
	TokenRefresh string `json:"refresh_token"`
	TokenScope   string `json:"scope"`
	TokenCreat   uint64 `json:"created_at"`
}

func (api42 *API42) readKeys(pathFile string) (string, error) {
	file, err := os.Open(pathFile)
	if err == nil {
		defer file.Close()
		fileBytes, err := ioutil.ReadFile(pathFile)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		fileString := strings.TrimSpace(string(fileBytes))
		if len(fileString) != cst.SizeApiKeys {
			log.Fatal().Msg("wrong size for " + pathFile)
		}
		log.Info().Msg("Found hash " + pathFile)
		return string(fileString), nil
	}
	return "", err
}

func (api42 *API42) setNewToken(newTokenAccess string, newTokenRefresh string) {
	err := ioutil.WriteFile(cst.PathTokenAccess, []byte(newTokenAccess), 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write in file " + cst.PathTokenAccess)
	}
	err = ioutil.WriteFile(cst.PathTokenRefresh, []byte(newTokenRefresh), 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write in file " + cst.PathTokenRefresh)
	}
	api42.keys.tokenAccess = newTokenAccess
	api42.keys.tokenRefresh = newTokenRefresh
	log.Info().Msg("New access token and refresh token set")
}

// RefreshToken use the refresh token to renew token access
func (api42 *API42) RefreshToken() {
	var err error

	log.Info().Msg("Refreshing tokens ...")

	tokenData := tokenReqRefresh{
		TokenGrant:     cst.TokenReqGrantRefresh,
		TokenCltID:     api42.keys.uid,
		TokenCltSecret: api42.keys.secret,
		TokenRefresh:   api42.keys.tokenRefresh,
	}
	tokenJSON, _ := json.Marshal(tokenData)

	rsp, err := http.Post(cst.TokenURL, "application/json", bytes.NewBuffer(tokenJSON))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve access token")
	}
	defer rsp.Body.Close()

	var rspJSON tokenRsp
	decoder := json.NewDecoder(rsp.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("Failed to decode JSON values of the new access token")
	}

	api42.setNewToken(rspJSON.TokenAccess, rspJSON.TokenRefresh)
}

// NewToken ask the 42's API to get a access token and a refresh token
func (api42 *API42) NewToken() {
	var err error

	urlAuth, _ := url.Parse(cst.AuthURL)
	paramAuth := url.Values{}
	paramAuth.Add(cst.AuthVarClt, api42.keys.uid)
	paramAuth.Add(cst.AuthVarRedirectURI, cst.AuthValRedirectURI)
	paramAuth.Add(cst.AuthVarRespType, cst.AuthValRespType)
	urlAuth.RawQuery = paramAuth.Encode()

	log.Info().Msg("Need new access token")
	fmt.Print("Please, enter the following URL in your web browser, authenticate and authorize:\n" + urlAuth.String() + "\nPaste the code generated (input hidden):\n")

	code := tools.ReadAndHideData()
	code = strings.TrimSpace(code)

	tokenData := tokenReqNew{
		TokenGrant:     cst.TokenReqGrantAuthCode,
		TokenCltID:     api42.keys.uid,
		TokenCltSecret: api42.keys.secret,
		TokenCode:      code,
		TokenRedirect:  cst.TokenReqRedirectURI,
	}

	tokenJSON, _ := json.Marshal(tokenData)

	rsp, err := http.Post(cst.TokenURL, "application/json", bytes.NewBuffer(tokenJSON))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve access token")
	}
	defer rsp.Body.Close()

	var rspJSON tokenRsp
	decoder := json.NewDecoder(rsp.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&rspJSON); err != nil {
		log.Fatal().Err(err).Msg("Failed to decode JSON values of the new access token")
	}

	api42.setNewToken(rspJSON.TokenAccess, rspJSON.TokenRefresh)
}

func (api42 *API42) initKeys() {
	api42.keys.uid = cst.UID
	var err error
	log.Info().Msg("Initialize keys of reqAPI42")
	if api42.keys.secret, err = api42.readKeys(cst.PathSecret); err != nil {
		log.Fatal().Err(err).Msg("Failed to read OAuth secret")
	}
	if api42.keys.tokenAccess, err = api42.readKeys(cst.PathTokenAccess); err != nil {
		log.Warn().Err(err).Msg("Failed to read OAuth access token")
		api42.NewToken()
	} else {
		if api42.keys.tokenRefresh, err = api42.readKeys(cst.PathTokenRefresh); err != nil {
			log.Fatal().Err(err).Msg("Failed to read OAuth refresh token")
		} else {
			log.Info().Msg("All keys are set")
		}
	}
}
