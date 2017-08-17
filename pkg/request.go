/*
*
* Program to do a load test, and validate the response team of our api
*
* @package     main
* @author      @jeffotoni
* @size        28/07/2017
*
 */

package request

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"net/http"
)

// var token = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiamVmZiIsImV4cCI6MTUwMjk0MjUzMCwiaXNzIjoicHJpbnRzZXJ2ZXIgemVicmEifQ.UuFVkuQ5UTE7Vu4RXgbKbb28AgjXmDnjpMEK1Sq866ozeCP2KNkK-L3ek6-aAErYA5ROODESYI7ASYLJ-k00Ff8mBbBLakqyZvCY5dPfYXbx9xfUzuGnlrtOyuTuxi3wQjKPgtfIsH8DN1aJN-_wRk6on9N6KHz-CE4NmKDj0_U`

type Ping struct {
	Msg string `json:"msg"`
}

func ShootUrl(Url string, Token string) string {

	var ping = &Ping{}

	req, err := http.NewRequest("POST", Url, nil)

	req.Header.Set("X-Custom-Header", "valueHeader")

	req.Header.Set("Authorization", "Bearer "+Token)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {

		panic(err)
	}

	defer resp.Body.Close()

	if resp.Status == "200 OK" {

		// fmt.Println("response Status:", resp.Status)

		// fmt.Println("response Headers:", resp.Header)

		body, _ := ioutil.ReadAll(resp.Body)

		// fmt.Println("response Body:", string(body))

		json.Unmarshal([]byte(string(body)), &ping)

		//
		//
		//
		msg2 := ping.Msg

		//
		//
		//
		ping.Msg = ""

		return string(msg2)

	} else {

		ping.Msg = ""
		return string("error")

	}
}
