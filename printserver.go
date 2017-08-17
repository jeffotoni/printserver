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
	HttpSucess         = 200
	HttpErrorLimit     = 429
	HttpErrorNoContent = 204
	HttpError          = 500
	HttpHeaderTitle    = `PrintServer`
	HttpHeaderMsg      = `Good Server, thank you.`
	NewLimiter         = 300
)

//
//
//
var (
	err           error
	returns       string
	confServer    *http.Server
	AUTHORIZATION = `bc9c154ebabc6f3da724e9x5fef78765`
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

// This method Message is to return our messages
// in json, ie the client will
// receive messages in json format
type Message struct {
	Code int    `json:code`
	Msg  string `json:msg`
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
//
//
func ShowScreen(cfg *Configs) {

	//
	//
	//
	login := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/login"

	//
	//
	//
	validate := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/validate"

	//
	//
	//
	ping := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/ping"

	//
	//
	//
	printer := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/print"

	//
	//
	//
	sizeBytes := 1 << 23

	sizeMb := (sizeBytes / 1024) / 1024

	//
	// Showing on the screen
	//
	fmt.Println("Start port:", cfg.ServerPort)
	fmt.Println("Endpoints:")
	fmt.Println(login)
	fmt.Println(validate)
	fmt.Println(ping)
	fmt.Println(printer)

	fmt.Println("Max bytes:", sizeMb, "Mb")

	//
	// Maximum 5 requests per second per client. Additional requests result in a HTTP 429 (Too Many Requests) error.
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
	// fmt.Println(string(pong))

	//
	//
	//
	w.Write(pong)

	//fmt.Println(pong)
	// nextHandler(w, r)

}

//
// Testing whether the service is online
//
func Ping2(w http.ResponseWriter, req *http.Request) {

	//
	//
	//
	json := `{"msg":"pong2"}`

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
	// fmt.Println(string(pong))

	//
	//
	//
	w.Write(pong)
}

//
//
//
func Print(w http.ResponseWriter, req *http.Request) {

	var json_msg string

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

			HttpMsgHeader = HttpErrorNoContent

		} else {

			HttpMsgHeader = HttpSucess

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

	// You can create a generic limiter for all your handlers
	// or one for each handler. Your choice.
	// This limiter basically says: allow at most 10 request per 1 second.
	limiter := tollbooth.NewLimiter(NewLimiter, time.Second)

	// This is an example on how to limit only GET and POST requests.
	limiter.Methods = []string{"GET", "POST"}

	// // You can also limit by specific request headers, containing certain values.
	// // Typically, you prefetched these values from the database.
	// limiter.Headers = make(map[string][]string)

	// limiter.Headers["X-Access-Token"] = []string{"xulxx", "383xx"}

	// And finally, you can limit access based on basic auth usernames.
	// Typically, you prefetched these values from the database as well.
	// limiter.BasicAuthUsers = []string{"xxx", "jeff", "youx"}

	//
	// Create a request limiter per handler.
	//
	// http.Handle("/ping", tollbooth.LimitFuncHandler(limiter, Ping))

	mux := http.NewServeMux()

	mux.Handle("/ping", tollbooth.LimitFuncHandler(limiter, Ping))

	mux.Handle("/ping2", tollbooth.LimitFuncHandler(limiter, Ping2))

	mux.Handle("/print", tollbooth.LimitFuncHandler(limiter, Print))

	nClassic := negroni.Classic()

	//n := negroni.New()
	nClassic.Use(negroni.HandlerFunc(auth0.ValidateToken))
	//nClassic.Use(negroni.HandlerFunc(Ping))
	// negroni.New(negroni.HandlerFunc(auth0.ValidateToken), negroni.HandlerFunc(Ping))

	nClassic.UseHandler(mux)

	//
	// Off the default mux
	//
	http.Handle("/login", tollbooth.LimitFuncHandler(limiter, auth0.Login))

	// http.HandleFunc("/ping", Ping)

	//
	// Create a request limiter per handler.
	//
	//http.Handle("/ping2", tollbooth.LimitFuncHandler(limiter, Ping2))

	//
	// Create the print server
	//
	//http.Handle("/print", tollbooth.LimitFuncHandler(limiter, Print))

	//
	// Login
	//
	//http.HandleFunc("/login", auth0.Login)

	//
	// Validate
	//
	// http.HandleFunc("/validate", auth0.ValidateToken)

	// http.Handle("/ping", negroni.New(negroni.HandlerFunc(auth0.ValidateToken), negroni.HandlerFunc(Ping)))

	//
	// list route
	//

	// httpList := reflect.ValueOf(http.DefaultServeMux).Elem()
	// finList := httpList.FieldByIndex([]int{1})
	// fmt.Println(finList)

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
