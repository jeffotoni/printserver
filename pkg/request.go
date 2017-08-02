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
	"fmt"
	"io/ioutil"
	"net/http"
)

type Ping struct {
	Msg string `json:"msg"`
}

func ShootUrl(Url string) {

	var ping = &Ping{}

	// Url := "http://localhost:9001/ping"

	fmt.Println("URL:>", Url)

	// var jsonStr = []byte(`{"msg":"hello"}`)
	// bytes.NewBuffer(jsonStr)

	// for i := 0; i < 10; i++ {

	// time.Sleep(time.Duration(200) * time.Millisecond)

	req, err := http.NewRequest("POST", Url, nil)

	req.Header.Set("X-Custom-Header", "valueHeader")

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

		fmt.Println("msg:", string(ping.Msg))

	} else {

		fmt.Println("response Status:", resp.Status)
		fmt.Println("Error")

	}

	ping.Msg = ""
	//}
}
