package verifyUser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//"strings"
//"log"

var authAPIScheme = os.Getenv("AUTH_API_SCHEME")
var authAPIHost = os.Getenv("AUTH_API_HOST")
var authAPIUri = os.Getenv("AUTH_API_URI")

type Claim struct {
	Email string `json:"https://goldenrecordstudios.earth/email"`
	Verified bool `json:"https://goldenrecordstudios.earth/email_verified"`
}

func VerifyUser (r *http.Request, claim *Claim) (int, error) {
	
	url := fmt.Sprintf("%s://%s%s", authAPIScheme, authAPIHost, authAPIUri)
	body := []byte("")

	request, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// pass through the original header containing the token
	request.Header = r.Header

	httpClient := http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		return http.StatusBadGateway, err
	}
	defer response.Body.Close()

	bodyRead, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(bodyRead, claim)
		if err != nil {
			// log.Println("Unmarshaling response body:")
			// log.Println(string(bodyRead))
			// log.Println(response.StatusCode, err.Error())
			return http.StatusInternalServerError, err
		}
		return response.StatusCode, nil
	case 401:
		// just pass the 401 code
		return response.StatusCode, nil
	default:
		err = errors.New(string(bodyRead))
		return response.StatusCode, err
	}

	err = errors.New("we should not have reached here")
	return http.StatusInternalServerError, err
}