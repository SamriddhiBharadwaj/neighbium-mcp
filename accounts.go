package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func buildGetBalanceSheetResponse(societyID int, asOnDate string, apiBody []byte) (map[string]any, error) {
	var apiData map[string]any
	err := json.Unmarshal(apiBody, &apiData)
	if err != nil {
		return nil, err
	}

	apiData["tool"] = "get_balance_sheet"
	apiData["status"] = "success"
	apiData["societyId"] = societyID
	apiData["asOnDate"] = asOnDate

	return apiData, nil
}

func buildGetLedgersResponse(societyID int, apiBody []byte) (GetLedgersResponse, error) {
	data, err := extractDataField(apiBody)
	if err != nil {
		return GetLedgersResponse{}, err
	}

	return GetLedgersResponse{
		Tool:      "get_ledgers",
		Status:    "success",
		SocietyID: societyID,
		Data:      data,
	}, nil
}

func extractDataField(apiBody []byte) (any, error) {
	var asObject map[string]any
	err := json.Unmarshal(apiBody, &asObject)
	if err == nil {
		data, hasData := asObject["data"]
		if hasData {
			return data, nil
		}
		return asObject, nil
	}

	var asArray []any
	err = json.Unmarshal(apiBody, &asArray)
	if err != nil {
		return nil, err
	}

	return asArray, nil
}

func currentAsOn(now time.Time) (int64, string) {
	asOn := toUnixMillis(now)
	asOnDate := now.Format("2006-01-02")
	return asOn, asOnDate
}

func dateUnixMillis(date string) (int64, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("invalid date %q: expected YYYY-MM-DD", date)
	}

	return toUnixMillis(parsed), nil
}

func buildGetIncomeAndExpenditureResponse(societyID int, fromDate string, toDate string, apiBody []byte) (map[string]any, error) {
	var apiData map[string]any
	if err := json.Unmarshal(apiBody, &apiData); err != nil {
		return nil, err
	}

	apiData["societyId"] = societyID
	apiData["period"] = map[string]string{
		"startDate": fromDate,
		"endDate":   toDate,
	}
	return apiData, nil
}
