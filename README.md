# printserver

A simple api rest to receive an encrypted POST and responsible for printing a label on the Zebra printer in the linux environment.

We will use the lpr native linux command to do the printing.

The way lpr works, in a nutshell, is: it reads in the file and hands the printable data over to the linux printing daemon, lpd. Lpd is a legacy piece of software for Unix and Linux, but it is supported under the modern system used by most Linux distributions, CUPS (the Common Unix Printing System).

You may need to manually install CUPS, and lpr itself, to print this way. If you are operating Debian, or a Debian-derived Linux system like Ubuntu that uses the APT package managements system, you can install them by running the following command:

Sudo apt-get update & & sudo apt-get install cups-client lpr

This command will install the Common Unix Printing System on your system. You should now be able to set up CUPS by directing any web browser to the address: http: // localhost: 631

The good thing is that we will send everything encrypted, but we can choose to encrypt the content before sending or not.


# Packages

go get -u github.com/didip/tollbooth

go get -u github.com/dgrijalva/jwt-go

go get -u github.com/codegangsta/negroni


# Install

$ go build printserver.go

$ sudo cp printserver /usr/bin

# Generate the keys

```sh

$ openssl genrsa -out private.rsa 1024

$ openssl rsa -in private.rsa -pubout > public.rsa.pub

```
# Simulate 

```go

$ go run test5_printserver.go

$ go run test4_printserver.go

$ go run test3_printserver.go

$ go run test2_printserver.go

$ go run test_printserver.go


```

# Simulate Curl

```sh

$ curl -X POST -H "Content-Type: application/json" \

-H "Authorization: Basic MjEyMzJmMjk3YTU3YTVhNzQzODk0YTBlNGE4MDFmYzM=:OTcyZGFkZGNhY2YyZmVhMjUzZmRhODY5NTY0ODUxMTU=" \

localhost:9001/token

$ curl -X POST -H "Content-Type: application/json" \

-H "Authorization: Bearer <TOKEN>" \

localhost:9001/ping

$ curl -X POST -H "Content-Type: application/x-www-form-urlencoded" \

-H "Authorization: Bearer <TOKEN>" \

-d "zpl='^xa^cfa,50^fo100,100^fdHello World!^fs^xz'&code=000198"

localhost:9001/print

```

# Main function

```go

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
	mux.Handle(HandlerOauthToken, tollbooth.LimitFuncHandler(limiter, auth0.LoginBasic))

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

```