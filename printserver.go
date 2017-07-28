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

import (
	"encoding/json"
	"fmt"
	"github.com/didip/tollbooth"
	"net/http"
	"time"
)

type MyInt int64

//
//
//
var (
	NewLimiter    MyInt
	err           error
	returns       string
	confServer    *http.Server
	AUTHORIZATION = `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
)

//
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
// Testing whether the service is online
//
func Ping(w http.ResponseWriter, req *http.Request) {

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
	w.Write(pong)
}

//
//
//
func Print(w http.ResponseWriter, req *http.Request) {

	//
	//
	//
	json_ok := `{"msg":"sucess"}`

	//
	//
	//
	// json_err := `{"msg":"error"}`

	//
	//
	//
	json := []byte(json_ok)

	//
	//
	//
	w.Write(json)
}

func main() {

	cfg := Config()

	ping := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/ping"
	printer := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/print"

	// Create a request limiter per handler.
	http.Handle("/ping", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(NewLimiter, time.Second), Ping))

	http.Handle("/print", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(NewLimiter, time.Second), Print))

	fmt.Println("Start port:", cfg.ServerPort)
	fmt.Println("Endpoints:")
	fmt.Println(ping)
	fmt.Println(printer)
	fmt.Println("Max bytes:", 1<<20, "bytes")

	//
	// Maximum 5 requests per second per client. Additional requests result in a HTTP 429 (Too Many Requests) error.
	//
	fmt.Println("NewLimiter:", NewLimiter)

	confServer = &http.Server{

		Addr: ":" + port,
		// Handler:        myHandler,
		// ReadTimeout:    1 * time.Second,
		// WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(confServer.ListenAndServe())

}
