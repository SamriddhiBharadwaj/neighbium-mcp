package main

import "encoding/json"

func buildGetBalanceSheetResponse(societyID int, asOnDate int64, apiBody []byte) (map[string]any, error) {
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

func buildGetLedgersBalanceResponse(societyID int, apiBody []byte) (GetLedgersBalanceResponse, error) {
	data, err := extractDataField(apiBody)
	if err != nil {
		return GetLedgersBalanceResponse{}, err
	}

	return GetLedgersBalanceResponse{
		Tool:      "get_ledgers_balance",
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

func buildGetIncomeAndExpenditureResponse(societyID int, fromDate int64, toDate int64, apiBody []byte) (map[string]any, error) {
	var apiData map[string]any
	if err := json.Unmarshal(apiBody, &apiData); err != nil {
		return nil, err
	}

	apiData["societyId"] = societyID
	apiData["period"] = map[string]int64{
		"startDate": fromDate,
		"endDate":   toDate,
	}
	return apiData, nil
}
