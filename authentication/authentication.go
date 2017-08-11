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
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"io/ioutil"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

//
// jwt init
//
func init() {

	privateByte, err := ioutil.ReadFile("./private.rsa")

	if err != nill {
		fmt.Println("Private key not found!")
	}

	publicByte, err := ioutil.ReadFile("./public.rsa.pub")

	if err != nill {
		fmt.Println("Public key not found!")
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		fmt.Println("No se pudo hacer el parse a privatekey")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		fmt.Println("No se pudo hacer el parse a privatekey")
	}
}

//
// jwt GenerateJWT
//
func GenerateJWT(user models.User) string {

	claims := models.Claim{

		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "printserver zebra",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	result, err := token.SignedString(privateKey)

	if err != nil {

		fmt.Println("No se pudo firmar el token!")

	}

	return result
}

func Login(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "Error al leer el usuario %s", err)
		return
	}

	if user.Name == "alexys" && user.Password == "alexys" {
		user.Password = ""
		user.Role = "admin"

		token := GenerateJWT(user)
		result := models.ResponseToken{token}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			fmt.Fprintln(w, "Error al generar el json")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Uusario o clave no válidos")
	}
}
