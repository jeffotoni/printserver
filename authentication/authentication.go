/*
*
* Project printServer, an api rest responsible for printing to a Zebra thermal printer.
* The printServer will receive a cryptographic POST containing
* the Zpl content so that the printer can print.
*
* @package     models
* @author      @jeffotoni
* @size        11/08/2017
*
 */

//
// openssl genrsa -out private.rsa 1024
// openssl rsa -in private.rsa -pubout > public.rsa.pub
//
package authentication

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	// "github.com/dgrijalva/jwt-go/request"
	"github.com/jeffotoni/printserver/models"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	pathPrivate = "./private.rsa"
	pathPublic  = "./public.rsa.pub"
)

//''
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
		fmt.Println("Private key not found!")
	}

	//
	//
	//
	publicByte, errx := ioutil.ReadFile(pathPublic)

	if errx != nil {
		fmt.Println("Public key not found!")
	}

	//
	//
	//
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateByte)

	if err != nil {
		fmt.Println("Could not parse privatekey")
	}

	//
	//
	//
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicByte)

	if err != nil {
		fmt.Println("Could not parse publickey: ", publicKey)
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

		User: model.User,
		StandardClaims: jwt.StandardClaims{

			// Expires in 8 hours
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "printserver zebra",
		},
	}

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

//
// login e password default
//
func Login(w http.ResponseWriter, r *http.Request) {

	//
	// Authorization Basic
	// $auth = self::$accessToken . ":" . self::$accessKey;
	// 'Authorization: Basic ' . base64_encode($auth)
	//
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	http.Error(w, "auth: "+auth[1], http.StatusUnauthorized)

	if len(auth) != 2 || auth[0] != "Basic" {

		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	token := strings.Trim(auth[1], " ")
	strings.TrimSpace(token)

	fmt.Println(token)

	os.Exit(1)
	//
	// read Body
	//
	byteBody := r.Body

	defer r.Body.Close()

	//
	// Validating json if correct
	//
	bodyJson, _ := ioutil.ReadAll(byteBody)

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
	fmt.Println(model.User)

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

	if model.User == "jeff" && model.Password == "1234" {

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
// Validate Token
//
func ValidateToken(w http.ResponseWriter, r *http.Request, handlerNext http.HandlerFunc) {

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {

		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	token := strings.Trim(auth[1], " ")
	strings.TrimSpace(token)

	// star
	parsedToken, err := jwt.ParseWithClaims(token, &models.Claim{}, func(*jwt.Token) (interface{}, error) {

		return publicKey, nil

	})

	if err != nil || !parsedToken.Valid {

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "Your token has expired!")
		return

	}

	claims, ok := parsedToken.Claims.(*models.Claim)

	if !ok {

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "There's something strange about your token")
		return
	}

	fmt.Println("User: ", claims.User)

	handlerNext(w, r)

	//return claims.User, nil
	// fmt.Println("err: ", err)
	// fmt.Println("token: ", parsedToken)

	// os.Exit(1)

	//end

	// if token, err := request.ParseFromRequest(r, request.OAuth2Extractor, keyLookupFunc2); err == nil {

	// 	claims := token.Claims.(jwt.MapClaims)
	// 	fmt.Printf("Token for user %v expires %v", claims["user"], claims["exp"])

	// } else if err != nil {

	// 	fmt.Println("Error: ", err)

	// 	switch err.(type) {

	// 	case *jwt.ValidationError:

	// 		vErr := err.(*jwt.ValidationError)

	// 		switch vErr.Errors {

	// 		case jwt.ValidationErrorExpired:
	// 			fmt.Fprintln(w, "Your token has expired!")
	// 			return
	// 		case jwt.ValidationErrorSignatureInvalid:
	// 			fmt.Fprintln(w, "Token signature does not match!")
	// 			return
	// 		default:
	// 			fmt.Fprintln(w, "Your token is invalid!")
	// 			return
	// 		}
	// 	default:
	// 		fmt.Fprintln(w, "Your token is invalid!")
	// 		return
	// 	}
	// }

	// if token.Valid {

	// 	w.WriteHeader(http.StatusAccepted)
	// 	fmt.Fprintln(w, "Welcome to the system!")

	// } else {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, "Your token is invalid!")
	// }
}

// func (a *JWTAuth) ValidateToken(token string) (string, error) {

// 	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaim{}, func(*jwt.Token) (interface{}, error) {
// 		return a.publicKey, nil
// 	})
// 	if err != nil || !parsedToken.Valid {
// 		return "", ErrPermissionDenied
// 	}
// 	claims, ok := parsedToken.Claims.(*tokenClaim)
// 	if !ok {
// 		return "", ErrPermissionDenied
// 	}
// 	return claims.Email, nil
// }

// func New(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) Auth {
// 	return &JWTAuth{privateKey: privateKey, publicKey: publicKey}
// }

func keyLookupFunc2(*jwt.Token) (interface{}, error) {

	fmt.Println("here public: ", publicKey)
	return publicKey, nil

}

// func keyLookupFunc(*Token) (interface{}, error) {

// 	// Don't forget to validate the alg is what you expect:
// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 	}

// 	// Look up key
// 	key, err := lookupPublicKey(token.Header["kid"])
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Unpack key from PEM encoded PKCS8
// 	return jwt.ParseRSAPublicKeyFromPEM(key)
// }
