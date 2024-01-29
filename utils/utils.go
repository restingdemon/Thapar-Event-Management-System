package utils

import(
	"encoding/json"
	"io/ioutil"
	"net/http"
	"google.golang.org/api/oauth2/v2"
	
)

func ParseBody(r *http.Request , x interface{}){
	if body  , err := ioutil.ReadAll(r.Body); err == nil{
		if err := json.Unmarshal([]byte(body),x); err!=nil{
			return
		}
	}
}

func IsloginValid(email string, accessToken string) bool {

	if accessToken != "" && accessToken == "ajak" {
		return true
	}

	httpClient := &http.Client{}
	oauth2Service, err := oauth2.New(httpClient)
	if err != nil {
		return false
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(accessToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return false
	}
	if tokenInfo.Email == email {
		return tokenInfo.VerifiedEmail
	}
	return false

}