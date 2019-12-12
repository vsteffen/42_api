package ReqAPI42

import (
	"fmt"
	// "bufio"
	// "bytes"
	// "encoding/json"
	"log"
	"os"
	"io/ioutil"
	"strings"
	"github.com/vsteffen/42_api/Constants"
	"golang.org/x/crypto/ssh/terminal"
	// "net/http"
	"net/url"
)

type keys struct {
	uid string
	secret string
	token string
}

var (
	gKeys = keys{constants.UID, "", ""}
)

func readKeys(pathFile string) (string, error) {
	file, err := os.Open(pathFile)
	if (err == nil) {
		defer file.Close()
		fileBytes, err := ioutil.ReadFile(pathFile)
		if (err != nil) {
			log.Fatal(err)
		}
		fileString := strings.TrimSpace(string(fileBytes))
		if (len(fileString) != constants.SizeKeys) {
			log.Fatal("wrong size for " + pathFile)
		}
		log.Print(pathFile + ": hash found")
		return string(fileString), nil
	}
	return "", err
}

func getCredential() (string) {
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if (err == nil) {
    	log.Fatal("Failed to retrieve credential")
    }
    password := string(bytePassword)

    return strings.TrimSpace(password)
}

func RefreshToken() {
/*	tokenData := []string{"authorization_code", gKeys.uid, gKeys.secret, "698c6cc85d07ed45f147486e6630ae9dedad619943999ce345fb958c4867bec9"}
	tokenJSON, _ := json.Marshal(tokenData)
	resp, err := http.Post(constants.TokenURL, "application/json", bytes.NewBuffer(tokenJSON))*/

    // resp, err := http.NewRequest("GET", "https://api.intra.42.fr/oauth/authorize?client_id=" + gKeys.uid, nil)

	urlAuth, _ := url.Parse(constants.AuthURL)
	paramAuth := url.Values{}
	paramAuth.Add(constants.AuthVarClt, gKeys.uid)
	paramAuth.Add(constants.AuthVarRedirectURI, constants.AuthValRedirectURI)
	paramAuth.Add(constants.AuthVarRespType, constants.AuthValRespType)
	urlAuth.RawQuery = paramAuth.Encode()

	fmt.Print("[INFO] Need new access token\n")
	fmt.Print("Please, enter the following URL in your web browser, authenticate and authorize:\n" + urlAuth.String() + "\nPaste the code generated: ")

/*	resp, err := http.Get(urlAuth.String())
	if (err != nil) {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode != http.StatusOK) {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Print(string(bodyBytes))
		log.Fatal("wrong response status code")
	}

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}*/
}

func initKeys() {
	var err error
	if gKeys.secret, err = readKeys(constants.PathSecret); err != nil {
		log.Fatal(err)
	}
	if gKeys.token, err = readKeys(constants.PathToken); err != nil {
		log.Print(err)
		RefreshToken()
	}
}
