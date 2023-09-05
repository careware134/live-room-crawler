package util

import (
	"bytes"
	"encoding/json"
	"io"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var logger = Logger()

func isInStringList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func IsInList(num int, list []int) bool {
	for _, s := range list {
		if s == num {
			return true
		}
	}
	return false
}

func DoRequestService(
	servicePart domain.ServiceStruct,
	traceId string,
	method string,
	requestUri string,
	params url.Values,
	requestBody *struct{},
	structType reflect.Type,
) (constant.ResponseStatus, *interface{}) {
	if servicePart.ApiBaseURL == "" {
		return constant.REQUEST_SERVER_FAILED, nil
	}
	loadRuleUrl := strings.Join([]string{servicePart.ApiBaseURL, "/", requestUri}, "")
	if params != nil {
		loadRuleUrl = strings.Join([]string{loadRuleUrl, "?", params.Encode()}, "")
	}

	var jsonBody []byte = nil
	if requestBody != nil {
		jsonBody, _ = json.Marshal(requestBody)
	}

	logger.Infof("util.GetRequestService request url:%s method:%s params:%s request is: %s", loadRuleUrl, method, params, jsonBody)
	req, err := http.NewRequest(method, loadRuleUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Warnf("util.GetRequestService request url:%s method:%s params:%s error is: %v", loadRuleUrl, method, params, err)
		return constant.REQUEST_SERVER_FAILED, nil
	}
	// Set custom headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", servicePart.Authorization)
	req.Header.Set("TenantId", servicePart.TenantId)
	req.Header.Set("CustomTraceId", traceId)

	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	body := string(bodyBytes)

	logger.Infof("util.GetRequestService request url:%s with code: %d response is: %s", loadRuleUrl, response.StatusCode, body)
	if err != nil || !IsInList(response.StatusCode, []int{200, 401, 403, 400}) {
		return constant.REQUEST_SERVER_FAILED, nil
	}

	instance := reflect.New(structType).Interface()
	err = json.Unmarshal(bodyBytes, instance)
	if err != nil {
		return constant.REQUEST_SERVER_FAILED, nil
	}

	return constant.SUCCESS, &instance
}
