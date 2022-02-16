package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DoDelete(url string) []byte {

	request, error := http.NewRequest("DELETE", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())

		panic(err)
	}

	fmt.Println("response Body:", string(bodyBytes))
	return bodyBytes
}

func DoPost(url string, body interface{}) []byte {

	jsonData := new(bytes.Buffer)
	json.NewEncoder(jsonData).Encode(body)

	request, error := http.NewRequest("POST", url, jsonData)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())

		panic(err)
	}

	fmt.Println("response Body:", string(bodyBytes))
	return bodyBytes
}

func DoPut(url string, body interface{}) []byte {

	jsonData := new(bytes.Buffer)
	json.NewEncoder(jsonData).Encode(body)

	request, error := http.NewRequest("PUT", url, jsonData)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())

		panic(err)
	}

	fmt.Println("response Body:", string(bodyBytes))
	return bodyBytes
}
