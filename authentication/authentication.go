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
	"github.com/dgrijalva/jwt-go/request"
	"github.com/jeffotoni/printserver/models"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

//
// jwt init
//
func init() {

	//
	//
	//
	privateByte, err := ioutil.ReadFile("./private.rsa")

	if err != nil {
		fmt.Println("Private key not found!")
	}

	//
	//
	//
	publicByte, errx := ioutil.ReadFile("./public.rsa.pub")

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
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicByte)

	if err != nil {
		fmt.Println("Could not parse publickey: ", publicKey)
	}
}

//
// jwt GenerateJWT
//
func GenerateJWT(user models.User) string {

	//
	// claims Token data, the header
	//
	claims := models.Claim{

		User: user.Name,
		StandardClaims: jwt.StandardClaims{

			// Expires in 8 hours
			ExpiresAt: time.Now().Add(time.Hour * 8).Unix(),
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

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {

		fmt.Fprintln(w, "Error reading user %s", err)
		return
	}

	if user.Name == "jeff" && user.Password == "1234" {

		user.Password = ""
		user.Role = "admin"

		token := GenerateJWT(user)

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
func ValidateToken(w http.ResponseWriter, r *http.Request) {

	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {

		fmt.Println("here public: ", publicKey)
		return publicKey, nil
	})

	if err != nil {

		fmt.Println("Error: ", err)

		switch err.(type) {

		case *jwt.ValidationError:

			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {

			case jwt.ValidationErrorExpired:
				fmt.Fprintln(w, "Your token has expired!")
				return
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Fprintln(w, "Token signature does not match!")
				return
			default:
				fmt.Fprintln(w, "Your token is invalid!")
				return
			}
		default:
			fmt.Fprintln(w, "Your token is invalid!")
			return
		}
	}

	if token.Valid {

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "Welcome to the system!")

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Your token is invalid!")
	}
}
