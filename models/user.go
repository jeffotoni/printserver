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

//
// jwt
//
type User struct {
	Name string `json:"name"`

	Password string `json:"password,omitempty"`

	Role string `json:"role"`

	// recommended having
	//jwt.StandardClaims
}
