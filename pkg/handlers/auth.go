package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// KeyProduct is a key used for the Product object inside context
type KeyProduct struct{}

// AuthHandler contains the items common to all product handler functions
type AuthHandler struct {
	logger *log.Logger
}

// NewAuthHandler returns a pointer to a AuthHandler with the logger passed as a parameter
func NewAuthHandler(logger *log.Logger) *AuthHandler {
	return &AuthHandler{logger}
}

func GetAdminAccessToken() string {
	urlPath := "http://localhost:8080/auth/realms/ubivius/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("client_id", "ubivius-client")
	data.Set("grant_type", "client_credentials")
	data.Set("client_secret", "7d109d2b-524f-4351-bfda-44ecad030eef")

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Admintoken response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	admin_token := extractValue(string(body), "access_token")
	return admin_token
}

// extracts the value for a key from a JSON-formatted string
// body - the JSON-response as a string. Usually retrieved via the request body
// key - the key for which the value should be extracted
// returns - the value for the given key
func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

func extractClaims(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, nil)
	if err != nil {
		log.Println("Error while getting claims")
	}
	return claims
}
