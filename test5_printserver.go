/*
*
* Program to do a load test, and validate the response team of our api
*
* @package     main
* @author      @jeffotoni
* @size        29/07/2017
*
 */

// _[Rate limiting](http://en.wikipedia.org/wiki/Rate_limiting)_
// is an important mechanism for controlling resource
// utilization and maintaining quality of service. Go
// elegantly supports rate limiting with goroutines,
// channels, and [tickers](tickers).

package main

import (
	"fmt"
	Shoot "github.com/jeffotoni/printserver/pkg"
	"os"
	"time"
)

func main() {

	token := ""

	if len(os.Args) > 1 {

		//fmt.Println(os.Args[1])
		token = os.Args[1]

	} else {

		fmt.Println("Passes the token as argument!")
	}

	// curl -X POST -H "Content-Type: application/json" -H "Authorization: Basic YWRtaW4=:MTIzNDU=" localhost:9001/login
	// os.Exit(1)

	endPoint1 := "http://localhost:9001/ping"
	//endPoint2 = "http://localhost:9001/ping2"

	// First we'll look at basic rate limiting. Suppose
	// we want to limit our handling of incoming requests.
	// We'll serve these requests off a channel of the
	// same name.
	requests := make(chan int, 350)

	for i := 1; i <= 350; i++ {

		println("Loading requests: ", fmt.Sprintf("%d", i))
		time.Sleep(time.Millisecond * 2)
		requests <- i
	}

	close(requests)

	// This `limiter` channel will receive a value
	// every 100 or 300 milliseconds. This is the regulator in
	// our rate limiting scheme.
	limiter := time.Tick(time.Millisecond * 5)

	// time start
	//
	//
	time1 := time.Now()

	// By blocking on a receive from the `limiter` channel
	// before serving each request, we limit ourselves to
	// 1 request every 200 milliseconds.
	for req := range requests {

		<-limiter

		msg := Shoot.ShootUrl(endPoint1, token)
		fmt.Println("request: ", req, "msg: ", msg)

		if req == 200 {
			fmt.Println("pause 2 segs")
			time.Sleep(time.Second * 5)
		}

	}

	time2 := time.Now()
	diff := time2.Sub(time1)

	fmt.Println(diff)
	fmt.Println("Enter enter to finish")

	var input string
	fmt.Scanln(&input)

}
