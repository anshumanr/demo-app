package api

import (
	"bytes"
	"fmt"
	"strings"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//Post - post the message to the specified address
func post(msgBody interface{}, HTTPUrl string, headermap map[string]string) ([]byte, error) {

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(msgBody)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err.(*json.UnsupportedTypeError)
	}

	url := formatURL(HTTPUrl)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	for key, element := range headermap {
		req.Header.Set(key, element)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	//close body when finished
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	//close body when finished
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	return body, err
}

//Get - send Get request to the specified address
func get(HTTPUrl string) ([]byte, error) {

	url := formatURL(HTTPUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	//close body when finished
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}

func getWithQueryParams(HTTPUrl string, queryParamMap map[string]string) ([]byte, error) {

	if len(queryParamMap) > 0 {
		params := "?"
		for paramName, paramValue := range queryParamMap {
			params += paramName + "=" + paramValue + "&"
		}
		HTTPUrl += params
	}

	return get(HTTPUrl)
}

func formatURL(url string) string {
	if !strings.HasPrefix(url, "http") {
		return "http://" + url
	}
	return url
}
