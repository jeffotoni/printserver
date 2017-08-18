/*
*
* Project printServer, an api rest responsible for printing to a Zebra thermal printer.
* The printServer will receive a cryptographic POST containing
* the Zpl content so that the printer can print.
*
* @package     main
* @author      @jeffotoni
* @size        28/07/2017
*
 */

package main

//
//
//
import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	// "github.com/dgrijalva/jwt-go"
	"github.com/didip/tollbooth"
	auth0 "github.com/jeffotoni/printserver/authentication"
	"log"
	"net/http"
	"os"
	"os/exec"
	// "reflect"
	"time"
)

//
//
//
const (
	HttpHeaderTitle = `PrintServer`
	HttpHeaderMsg   = `Good Server, thank you.`
	NewLimiter      = 250
	SizeByteAllowed = 1 << 23

	HandlerToken   = "/token"
	HandlerV1Print = "/print"
	HandlerPing    = "/ping"
)

//
//
//
var (
	err        error
	returns    string
	confServer *http.Server
)

//''
// Structure of our server configurations
//
type Configs struct {
	Domain     string `json:"domain"`
	Process    string `json:"process"`
	Ping       string `json:"ping"`
	ServerPort string `json:"serverport"`
	Host       string `json:"host"`
	Schema     string `json:"shcema"`
	ServerHost string `json:"serverhost"`
}

// This method Message is to return our messages
// in json, ie the client will
// receive messages in json format
type Message struct {
	Code int    `json:code`
	Msg  string `json:msg`
}

//
// Type responsible for defining a function that returns boolean
//
type fn func(w http.ResponseWriter, r *http.Request) bool

//
// This method ConfigJson sets up our
// server variables from our struct
//
func ConfigJson() string {

	// Defining the values of our config
	data := &Configs{Domain: "localhost", Process: "2", Ping: "ok", ServerPort: "9001", Host: "", Schema: "http", ServerHost: "localhost"}

	// Converting our struct into json format
	cjson, err := json.Marshal(data)
	if err != nil {
		// handle err
	}

	return string(cjson)
}

// This method Config returns the objects
// of our config so that it can be accessed
func Config() *Configs {

	var objason Configs

	jsonT := []byte(ConfigJson())
	json.Unmarshal(jsonT, &objason)

	return &objason
}

//
// This function already contains the method that does
// the validation of the token, only receives a parameter
// that is the handler to be executed if token is true
//
func HandlerAuth(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//if auth(w, r) {
		if auth0.HandlerValidate(w, r) {

			handler(w, r)
		}
	}
}

//
// Function responsible for abstraction and receive the
// authentication function and the handler that will execute if it is true
//
func HandlerFuncAuth(auth fn, handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if auth(w, r) {

			handler(w, r)
		}
	}
}

// This method is a simplified abstraction
// so that we can send them to our client
// when making a request
func JsonMsg(codeInt int, msgText string) string {

	data := &Message{Code: codeInt, Msg: msgText}

	djson, err := json.Marshal(data)
	if err != nil {
		// handle err
	}

	return string(djson)
}

//
//
//
func check(e error) {

	if e != nil {

		panic(e)
	}
}

//
// Mounting the properties on the api screen
//
func ShowScreen(cfg *Configs) {

	//
	// Basic Authentication
	//
	oauthToken := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + HandlerToken

	//
	//
	//
	ping := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + HandlerPing

	//
	//
	//
	printer := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + HandlerV1Print

	//
	//
	//
	sizeMb := (SizeByteAllowed / 1024) / 1024

	//
	// Showing on the screen
	//
	fmt.Println("Start port:", cfg.ServerPort)
	fmt.Println("Endpoints:")
	fmt.Println(oauthToken)
	fmt.Println(ping)
	fmt.Println(printer)
	fmt.Println("Max bytes:", sizeMb, "Mb")

	//
	// Maximum NewLimiter requests per second per client. Additional requests result in a HTTP 429 (Too Many Requests) error.
	//
	fmt.Println("Requests ", NewLimiter, "per 1 second")
}

//
// Testing whether the service is online
//
func Ping(w http.ResponseWriter, r *http.Request) {

	//
	//
	//
	json := `{"msg":"pong"}`

	//
	//
	//
	pong := []byte(json)

	//
	//
	//
	w.Header().Set(HttpHeaderTitle, HttpHeaderMsg)

	//
	//
	//
	w.Header().Set("X-Custom-Header", "HeaderValue-x83838374774")

	//
	//
	//
	w.Header().Set("Content-Type", "application/json")

	//
	//
	//
	w.WriteHeader(http.StatusOK)

	//
	//
	//
	w.Write(pong)
}

//
// Zebra printer, the printer is networked, or the local
// machine is connected to the printer.
// The program will send a file to printer
//
func Print(w http.ResponseWriter, req *http.Request) {

	//
	//
	//
	var json_msg string

	//
	//
	//
	var HttpMsgHeader int

	//
	//
	//
	if req.Method == "POST" {

		//
		//
		//
		ZPL := req.FormValue("zpl")

		//
		//
		//
		CODE := req.FormValue("code")

		if len(CODE) == 0 || len(ZPL) == 0 {

			json_msg = `{"msg":"error No Content!"}`

			HttpMsgHeader = http.StatusBadRequest

		} else {

			HttpMsgHeader = http.StatusOK

			//
			//
			//
			PathZpl := "/tmp/" + CODE + ".zpl"

			//
			// Generate a .zpl file
			//
			ZplByte := []byte(ZPL)

			//
			//
			//
			f, err := os.Create(PathZpl)

			//
			//
			check(err)

			//
			//
			//
			defer f.Close()

			//
			//
			//
			size, errx := f.Write(ZplByte)

			//
			//
			//
			// err = ioutil.WriteFile(PathZpl, ZplByte, 0754)

			check(errx)

			//
			//
			//
			// command := "lpr -P zebra -o raw "

			//
			// To print zpl on zebra printer in linux terminal
			//

			out, errc := exec.Command("lpr", "-P", "zebra", "-o", "raw", PathZpl).Output()

			if errc != nil {

				// log.Fatal(errc)
				fmt.Println(errc)
			}

			//
			//
			//
			fmt.Printf("Running command%s\n", out)

			//
			//
			//
			json_msg = `{"msg":"Printing performed successfully","bytes":"` + fmt.Sprintf("%d", size) + `"}`

		}

	} else {

		//
		//
		//
		json_msg = `{"msg":"Error Only accepts POST"}`
	}

	//
	//
	//
	json := []byte(json_msg)

	//
	//
	//
	w.Header().Set(HttpHeaderTitle, HttpHeaderMsg)

	//
	//
	//
	w.Header().Set("X-Custom-Header", "HeaderValue-x83838374774")

	//
	//
	//
	w.Header().Set("Content-Type", "application/json")

	//
	//
	//
	w.WriteHeader(HttpMsgHeader)

	//
	//
	//
	w.Write(json)
}

//
// start
//
func main() {

	//
	//
	//
	cfg := Config()

	//
	//
	//
	ShowScreen(cfg)

	// Creating limiter for all handlers
	// or one for each handler. Your choice.
	// This limiter basically says: allow at most NewLimiter request per 1 second.
	limiter := tollbooth.NewLimiter(NewLimiter, time.Second)

	// Limit only GET and POST requests.
	limiter.Methods = []string{"GET", "POST"}

	//
	//
	//
	mux := http.NewServeMux()

	// We had problem in doing method authentication and limit rate using negroni
	// mux.Handle("/ping", negroni.New(negroni.HandlerFunc(MyMiddlewareAuth0), negroni.HandlerFunc(MyMiddlewarePing)))

	mux.Handle(HandlerPing, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(auth0.HandlerValidate, Ping)))

	mux.Handle(HandlerV1Print, tollbooth.LimitFuncHandler(limiter, HandlerAuth(Print)))

	//
	// Off the default mux
	// Does not need authentication, only user key and token
	//
	mux.Handle(HandlerToken, tollbooth.LimitFuncHandler(limiter, auth0.LoginBasic))

	// mux.Handle("/validate", tollbooth.LimitFuncHandler(limiter, auth0.ValidateToken))

	nClassic := negroni.Classic()

	nClassic.UseHandler(mux)

	//
	//
	//
	confServer = &http.Server{

		Addr: ":" + cfg.ServerPort,

		Handler: nClassic,
		//ReadTimeout:    30 * time.Second,
		//WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 23, // Size accepted by package
	}

	log.Fatal(confServer.ListenAndServe())

}
