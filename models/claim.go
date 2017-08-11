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

package models

import "github.com/dgrijalva/jwt-go"

//
// jwt
//
type Claim struct {
	User string `json:"user"`

	jwt.StandardClaims
}
