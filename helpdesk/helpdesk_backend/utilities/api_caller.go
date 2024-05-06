package utilities

import (
	"bytes"
	"encoding/json"
	"helpdesk_backend/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"moul.io/http2curl"
	// "moul.io/http2curl"
)

func CallApi(urlAddress string, httpMethod string, headers map[string]string, payload interface{}) ([]byte, int, error) {

	hc := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}
	logger.ZapLogger.Info("Hitting on url  ->    ", urlAddress, "  with method->  ", httpMethod)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.ZapLogger.Error(err)
		return nil, http.StatusBadRequest, err
	}

	var hr *http.Request
	if headers["Content-Type"] == "application/x-www-form-urlencoded" {
		if payload != nil {
			hr, err = http.NewRequest(httpMethod, urlAddress, strings.NewReader(payload.(string)))
		} else {
			hr, err = http.NewRequest(httpMethod, urlAddress, strings.NewReader(`null`))
		}
		if err != nil {
			logger.ZapLogger.Error(err)
			return nil, http.StatusInternalServerError, err
		}
	} else {
		hr, err = http.NewRequest(httpMethod, urlAddress, bytes.NewBuffer(payloadJSON))
		if err != nil {
			logger.ZapLogger.Error(err)
			return nil, http.StatusInternalServerError, err
		}
	}

	for key, value := range headers {
		hr.Header.Set(key, value)
		//fmt.Println(key + ":" + value)
	}

	logger.ZapLogger.Info(http2curl.GetCurlCommand(hr))

	response, err := hc.Do(hr)
	if err != nil {
		logger.ZapLogger.Error(err)
		return nil, http.StatusBadGateway, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ZapLogger.Error(err)
		return nil, response.StatusCode, err
	}
	logger.ZapLogger.Info("Response from url  ->    ", urlAddress, "  with method->  ", httpMethod, "  is  ->  ", string(body), " response status code = ", response.StatusCode)
	return body, response.StatusCode, nil
}
