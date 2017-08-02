/*
*
* Program to do a load test, and validate the response team of our api
*
* @package     main
* @author      @jeffotoni
* @size        29/07/2017
*
 */

package main

import (
	"fmt"
	Shoot "github.com/jeffotoni/printserver/pkg"
	"time"
)

func main() {

	// For our example we'll select across two channels.
	c1 := make(chan string)
	c2 := make(chan string)

	endPoint1 := "http://localhost:9001/ping"
	endPoint2 := "http://localhost:9001/ping2"

	go func() {

		for i := 0; i < 10; i++ {

			time.Sleep(time.Millisecond * 100)
			Shoot.ShootUrl(endPoint1)
		}

		c2 <- "two"
	}()

	// Each channel will receive a value after some amount
	// of time, to simulate e.g. blocking RPC operations
	// executing in concurrent goroutines.
	go func() {

		for i := 0; i < 10; i++ {

			time.Sleep(time.Millisecond * 300)
			Shoot.ShootUrl(endPoint2)
		}

		c1 <- "one"
	}()

	msg := <-c2

	fmt.Println(msg)

	// for i := 0; i < 2; i++ {
	//  select {
	//  case msg1 := <-c1:
	//      fmt.Println("received", msg1)
	//  case msg2 := <-c2:
	//      fmt.Println("received", msg2)
	//  }
	// }
}
