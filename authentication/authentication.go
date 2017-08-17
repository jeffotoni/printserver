/*
*
* Project printServer, an api rest responsible for printing to a Zebra thermal printer.
* The printServer will receive a cryptographic POST containing
* the Zpl content so that the printer can print.
*
* @package     authentication
* @author      @jeffotoni
* @size        11/08/2017
*
 */

package authentication

import (
	"crypto/rsa"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jeffotoni/printserver/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	pathPrivate = "./private.rsa"
	pathPublic  = "./public.rsa.pub"

	ProjectTitle = "printserver zebra"
	Expiration   = 24 * 12 // Hours and Day
)

//
// Structure of our server configurations
//
type JsonMsg struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

//
// jwt init
//
func init() {

	//
	//
	//
	privateByte, err := ioutil.ReadFile(pathPrivate)

	if err != nil {

		WriteJson("error", "Private key not found!", http.StatusUnauthorized)
		return
	}

	//
	//
	//
	publicByte, errx := ioutil.ReadFile(pathPublic)

	if errx != nil {

		WriteJson("error", "Public key not found!", http.StatusUnauthorized)
		return
	}

	//
	//
	//
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateByte)

	if err != nil {

		WriteJson("error", "Could not parse privatekey!", http.StatusUnauthorized)
		return
	}

	//
	//
	//
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicByte)

	if err != nil {

		WriteJson("error", "ould not parse publickey!", http.StatusUnauthorized)
		return
	}
}

//
// jwt GenerateJWT
//
func GenerateJWT(model models.User) string {

	//
	// claims Token data, the header
	//
	claims := models.Claim{

		User: model.Login,
		StandardClaims: jwt.StandardClaims{

			// Expires in 8 hours
			ExpiresAt: time.Now().Add(time.Minute * Expiration).Unix(),
			Issuer:    ProjectTitle,
		},
	}

	//
	//
	//
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	//
	// Transforming into string
	//
	tokenString, err := token.SignedString(privateKey)

	if err != nil {

		fmt.Println("Could not sign the token!")

	}

	//
	// return token string
	//
	return tokenString
}

func LoginJson(w http.ResponseWriter, r *http.Request) {

	//
	// Validating json if correct
	//
	bodyJson, _ := ioutil.ReadAll(r.Body)

	//
	// Looking for keys in the first and last position
	//
	last_pos := len(bodyJson) - 1

	if string(bodyJson[0]) != "{" {

		msgJsonStruct := &JsonMsg{"Error", "Missing keys on your json '{'"}
		msgJson, errj := json.Marshal(msgJsonStruct)

		if errj != nil {

			fmt.Fprintln(w, "Error generating json!")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(msgJson)

		return
	}

	if string(bodyJson[last_pos]) != "}" {

		msgJsonStruct := &JsonMsg{"Error", "Missing keys on your json '}'"}
		msgJson, errj := json.Marshal(msgJsonStruct)

		if errj != nil {

			fmt.Fprintln(w, "Error generating json!")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(msgJson)

		return
	}

	var model models.User

	//
	//
	//
	// err := json.NewDecoder(r.Body).Decode(&user)

	err := json.Unmarshal(bodyJson, &model)

	fmt.Println("Err: ", err)
	fmt.Println(model.Login)

	if err != nil {

		msgJsonStruct := &JsonMsg{"Error", "Error reading json, Configures your json syntax!"}
		msgJson, errj := json.Marshal(msgJsonStruct)

		if errj != nil {

			fmt.Fprintln(w, "Error generating json!")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(msgJson)

		//fmt.Println("Body:", )

		return
	}

	if model.Login == "jeff" && model.Password == "1234" {

		model.Password = ""
		model.Role = "admin"

		token := GenerateJWT(model)

		result := models.ResponseToken{token}
		jsonResult, err := json.Marshal(result)

		if err != nil {
			fmt.Fprintln(w, "Error generating json!")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)

	} else {

		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Invalid user or key!")
	}
}

//
// login e password default
//
func LoginBasic(w http.ResponseWriter, r *http.Request) {

	//
	// Authorization Basic base64 Encode
	//
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {

		HttpWriteJson(w, "error", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	tokenBase64 := strings.Trim(auth[1], " ")
	tokenBase64 = strings.TrimSpace(tokenBase64)

	//
	// token 64
	//
	authToken64 := strings.SplitN(tokenBase64, ":", 2)

	if len(authToken64) != 2 || authToken64[0] == "" || authToken64[1] == "" {

		HttpWriteJson(w, "error", "token base 64 invalid!", http.StatusUnauthorized)
		return
	}

	// fmt.Println(authToken64[0])
	// fmt.Println(authToken64[1])

	tokenUserEnc := authToken64[0]
	keyUserEnc := authToken64[1]

	tokenUserDecode, _ := b64.StdEncoding.DecodeString(tokenUserEnc)
	//fmt.Println(string(tokenUserDecode))

	keyUserDec, _ := b64.StdEncoding.DecodeString(keyUserEnc)
	//fmt.Println(string(keyUserDec))

	tokenUserDecodeS := strings.ToUpper(string(tokenUserDecode))
	keyUserDecS := strings.ToUpper(string(keyUserDec))

	if tokenUserDecodeS == "ADMIN" && keyUserDecS == "12345" {

		var model models.User

		model.Login = tokenUserDecodeS
		// model.Password = keyUserDec

		model.Password = ""
		model.Role = "admin"

		token := GenerateJWT(model)

		result := models.ResponseToken{token}
		jsonResult, err := json.Marshal(result)

		if err != nil {

			// fmt.Fprintln(w, "Error generating json!")
			HttpWriteJson(w, "error", "json.Marshal error generating!", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)

	} else {

		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Invalid user or key!")
	}

	HttpWriteJson(w, "sucess", http.StatusText(http.StatusOK), http.StatusOK)

	defer r.Body.Close()

	/**

	{
	  "accessToken": "39a3099b45634f6eb511991fbbe83752_v2",
	  "access_token": "39a3099b45634f6eb511991fbbe83752_v2",
	  "expires_in": "2026-09-14",
	  "refreshToken": "1defad3474a8423f87a04adc588e7c7b_v2",
	  "refresh_token": "1defad3474a8423f87a04adc588e7c7b_v2",
	  "scope": "RECEIVE_FUNDS,REFUND,MANAGE_ACCOUNT_INFO",
	  "moipAccount": {
	    "id": "MPA-SVOAZT7WWHGB"
	  }

	*/
}

//
// HandlerValidate
//
func HandlerValidate(w http.ResponseWriter, r *http.Request) bool {

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {

		//http.Error(w, "authorization failed", http.StatusUnauthorized)
		HttpWriteJson(w, "error", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	token := strings.Trim(auth[1], " ")
	strings.TrimSpace(token)

	// star
	parsedToken, err := jwt.ParseWithClaims(token, &models.Claim{}, func(*jwt.Token) (interface{}, error) {

		return publicKey, nil

	})

	if err != nil || !parsedToken.Valid {

		//w.WriteHeader(http.StatusAccepted)
		//fmt.Fprintln(w, "Your token has expired!")
		HttpWriteJson(w, "error", "Your token has expired!", http.StatusAccepted)
		return false

	}

	claims, ok := parsedToken.Claims.(*models.Claim)

	if !ok || claims.User != "ADMIN" {

		//w.WriteHeader(http.StatusAccepted)
		//HttpWriteJson(w, "error", "There's something strange about your token!", http.StatusAccepted)
		fmt.Fprintln(w, "There's something strange about your token")
		return false
	}

	// fmt.Println("User: ", claims.User)

	//HttpWriteJson(w, "success", "Your token it's ok ["+claims.User+"]", http.StatusOK)
	//func2(w, r)
	return true
}

//
// Validate Token
//
func ValidateToken(w http.ResponseWriter, r *http.Request) {

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {

		//http.Error(w, "authorization failed", http.StatusUnauthorized)
		HttpWriteJson(w, "error", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token := strings.Trim(auth[1], " ")
	strings.TrimSpace(token)

	// star
	parsedToken, err := jwt.ParseWithClaims(token, &models.Claim{}, func(*jwt.Token) (interface{}, error) {

		return publicKey, nil

	})

	if err != nil || !parsedToken.Valid {

		//w.WriteHeader(http.StatusAccepted)
		//fmt.Fprintln(w, "Your token has expired!")
		HttpWriteJson(w, "error", "Your token has expired!", http.StatusAccepted)
		return

	}

	claims, ok := parsedToken.Claims.(*models.Claim)

	if !ok || claims.User != "ADMIN" {

		//w.WriteHeader(http.StatusAccepted)
		HttpWriteJson(w, "error", "There's something strange about your token!", http.StatusAccepted)
		//fmt.Fprintln(w, "There's something strange about your token")
		return
	}

	// fmt.Println("User: ", claims.User)

	HttpWriteJson(w, "success", "Your token it's ok ["+claims.User+"]", http.StatusOK)
}

func WriteJson(Status string, Msg string, httpStatus int) {

	msgJsonStruct := &JsonMsg{Status, Msg}

	msgJson, errj := json.Marshal(msgJsonStruct)

	if errj != nil {

		fmt.Println(`{"status":"error","msg":"We could not generate the json error!"}`)
		return
	}

	fmt.Println(msgJson)
}

func HttpWriteJson(w http.ResponseWriter, Status string, Msg string, httpStatus int) {

	msgJsonStruct := &JsonMsg{Status, Msg}

	msgJson, errj := json.Marshal(msgJsonStruct)

	if errj != nil {

		fmt.Fprintln(w, `{"status":"error","msg":"We could not generate the json error!"}`)
		return
	}

	w.WriteHeader(httpStatus)

	w.Header().Set("Content-Type", "application/json")

	w.Write(msgJson)
}
