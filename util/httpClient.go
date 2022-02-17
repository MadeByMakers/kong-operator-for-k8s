package httpClient

import (
	"bytes"
	"encoding/json"
)

func GetBaseURL() string {
	return "https://soap.entel.cl/api/kong"
}

func doRequest(method string, url string, data interface{}) (int, []byte) {

	jsonData := new(bytes.Buffer)
	json.NewEncoder(jsonData).Encode(data)

	return 200, jsonData.Bytes()

	/*
		var body io.Reader = nil

		if data != nil {
			jsonData := new(bytes.Buffer)
			json.NewEncoder(jsonData).Encode(data)
			body = jsonData
		}

		request, error := http.NewRequest(method, url, body)
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}

		defer response.Body.Close()
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		return response.StatusCode, bodyBytes
	*/
}

func Delete(url string) (int, []byte) {
	return doRequest("DELETE", url, nil)
}

func Post(url string, data interface{}) (int, []byte) {
	return doRequest("POST", url, data)
}

func Put(url string, data interface{}) (int, []byte) {
	return doRequest("PUT", url, data)
}

func Patch(url string, data interface{}) (int, []byte) {
	return doRequest("PATCH", url, data)
}

func Get(url string) (int, []byte) {
	return doRequest("GET", url, nil)
}
