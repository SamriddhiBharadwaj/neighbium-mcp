package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const neighbiumBaseURL = "https://api.neighbium.com/NeighbiumServer/rest"

func neighbiumGet(path string) ([]byte, error) {
	url := neighbiumBaseURL + path

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	headerName := os.Getenv("NEIGHBIUM_API_HEADER_NAME")
	apiKey := os.Getenv("NEIGHBIUM_API_KEY")
	authorization := os.Getenv("NEIGHBIUM_AUTHORIZATION")

	request.Header.Set(headerName, apiKey)
	request.Header.Set("Authorization", authorization)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("neighbium API returned %s: %s", response.Status, body)
	}

	return body, nil
}

func getMemberDirectory(societyID int) ([]DirectoryEntry, error) {
	path := fmt.Sprintf("/societies/%d/members/directory", societyID)

	body, err := neighbiumGet(path)
	if err != nil {
		return nil, err
	}

	var entries []DirectoryEntry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func getBalanceSheet(societyID int, asOn int64) ([]byte, error) {
	path := fmt.Sprintf("/societies/%d/accounts/balance-sheet?asOn=%d", societyID, asOn)
	return neighbiumGet(path)
}

func getLedgersBalance(societyID int, asOnDate int64) ([]byte, error) {
	path := fmt.Sprintf("/societies/%d/accounts/ledgers/balance?financialYear=0&asOnDate=%d", societyID, asOnDate)
	return neighbiumGet(path)
}

func getIncomeAndExpenditure(societyID int, startDate int64, endDate int64) ([]byte, error) {
	path := fmt.Sprintf("/societies/%d/accounts/ie?startDate=%d&endDate=%d", societyID, startDate, endDate)
	return neighbiumGet(path)
}
