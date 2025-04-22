package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Resp() { 
	resp, err := http.Get("http://itsthisforthat.com/api.php?text")// https://dev.to/billylkc/parse-json-api-response-in-go-10ng - Parse JSON API response in Go
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))  
}       