package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
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

func FindJsonString(targetString string, startIndex int) string {
	jsonString := "{}"
	if startIndex == -1 {
		fmt.Println("No JSON object found in the long string.")
		return jsonString
	}

	openBrackets := 1
	endIndex := startIndex + 1

	for openBrackets > 0 && endIndex < len(targetString) {
		if targetString[endIndex] == '{' {
			openBrackets++
		} else if targetString[endIndex] == '}' {
			openBrackets--
		}
		endIndex++
	}

	if openBrackets > 0 {
		fmt.Println("Invalid JSON structure. Unmatched curly brackets.")
		return jsonString
	}

	jsonString = targetString[startIndex:endIndex]
	jsonString, _ = removeOneLayerOfEscape(jsonString)
	return jsonString
}

func removeOneLayerOfEscape(jsonString string) (string, error) {
	// Check if the string is already a valid JSON by attempting to unmarshal it
	var obj interface{}
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err == nil {
		// The string is already a valid JSON, return it as is
		return jsonString, nil
	}

	// Validate the unescaped string by attempting to unmarshal it
	unescapedString, err := strconv.Unquote(`"` + jsonString + `"`)
	if err != nil {
		logger.Warnf("Fail to Unquote json:%s with error:%e", jsonString, err)
		// The unescaped string is still not a valid JSON
		return "{}", fmt.Errorf("failed to remove escape characters from JSON string")
	}

	err = json.Unmarshal([]byte(unescapedString), &obj)
	if err != nil {
		logger.Warnf("Fail to Unmarshal after escape json:%s with error:%e", jsonString, err)
		// The unescaped string is still not a valid JSON
		return "{}", fmt.Errorf("failed to remove escape characters from JSON string")
	}

	return unescapedString, nil
}

func ParseChineseNumber(chineseNumber string) (int, error) {
	// Remove any trailing characters like "+"
	chineseNumber = strings.TrimSuffix(chineseNumber, "+")

	// Handle "万" suffix
	if strings.Contains(chineseNumber, "万") {
		parts := strings.Split(chineseNumber, "万")
		if len(parts) != 2 {
			return 0, fmt.Errorf("Invalid Chinese number format: %s", chineseNumber)
		}

		base, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		tenThousand, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}

		return base*10000 + tenThousand, nil
	}

	// Handle regular numbers
	number, err := strconv.Atoi(chineseNumber)
	if err != nil {
		return 0, err
	}

	return number, nil
}
