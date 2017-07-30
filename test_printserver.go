/*
*
* Program to do a load test, and validate the response team of our api
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
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Ping struct {
	Msg string `json:"msg"`
}

var cont = 0

func main() {

	cont = 0

	var ping = &Ping{}

	vetUrl := make(map[int]string)

	url := "http://localhost:9001/ping"

	vetUrl[0] = url

	min := 0
	max := 1

	// time.Sleep(1 * time.Second)
	c1 := make(chan string)
	c2 := make(chan string)

	startingTime := time.Now().UTC()

	fmt.Println("Time start")

	i := 0

	for i = 0; i < 10; i++ {

		time.Sleep(time.Duration(300) * time.Millisecond)

		SendPing2(ping, url, 1)

	}

	endingTime := time.Now().UTC()

	var duration time.Duration = endingTime.Sub(startingTime)

	fmt.Println("Duration: ", duration, "request: ", i)

	os.Exit(1)

	/** ########################################### **/
	/** ########################################### **/
	/** ################# test #################### **/
	/** ########################################### **/
	/** ########################################### **/

	go func() {

		sleep := 100
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		SendPing(ping, vetUrl, min, max, 1, 10, c1)

	}()

	go func() {

		sleep := 100
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		SendPing(ping, vetUrl, min, max, 2, 20, c2)

		// c2 <- "two"

	}()

	for i := 0; i < 2; i++ {

		select {

		case msg1 := <-c1:
			fmt.Println("received", msg1)

		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}

}

func SendPing2(ping *Ping, url string, gorutine int) {

	MsgTmp := ""

	response, err := http.Get(url)

	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {

		fmt.Println(err)
	}

	json.Unmarshal([]byte(string(contents)), &ping)

	MsgTmp = strings.TrimSpace(ping.Msg)
	MsgTmp = strings.Trim(MsgTmp, " ")

	if len(MsgTmp) > 0 {

		fmt.Println("gorutine one: ", fmt.Sprintf("%d", gorutine), " :: ", MsgTmp)

	} else {

		fmt.Println("gorutine 2: ", fmt.Sprintf("%d", gorutine), " :: ", "error")
		// fmt.Println(string(contents))
	}

	ping.Msg = ""
}

func SendPing(ping *Ping, vetUrl map[int]string, min int, max int, gorutine int, sleep int, c chan string) {

	i := 1
	MsgTmp := ""

	for i < sleep {

		i += i

		cont++

		time.Sleep(time.Duration(500) * time.Millisecond)

		seed := Seed(min, max)
		response, err := http.Get(vetUrl[seed])

		if err != nil {

			fmt.Println(err)
			os.Exit(1)

		} else {

			defer response.Body.Close()

			contents, err := ioutil.ReadAll(response.Body)

			if err != nil {

				fmt.Println(err)
			}

			json.Unmarshal([]byte(string(contents)), &ping)

			MsgTmp = strings.TrimSpace(ping.Msg)
			MsgTmp = strings.Trim(MsgTmp, " ")

			if len(MsgTmp) > 0 {

				fmt.Println("cont: ", cont, "gorutine one: ", fmt.Sprintf("%d", gorutine), " :: ", MsgTmp)

			} else {

				fmt.Println("cont: ", cont, "gorutine 2: ", fmt.Sprintf("%d", gorutine), " :: ", "error")

			}

			ping.Msg = ""
		}
	}

	c <- fmt.Sprintf("%d", gorutine)
}

func Seed(min int, max int) int {

	return (min + rand.Intn(max-min))
}
