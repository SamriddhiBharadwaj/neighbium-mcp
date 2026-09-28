package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetMemberByIDInput struct {
	SocietyID int `json:"societyId"`
	MemberID  int `json:"memberId"`
}

type GetMembersInput struct {
	SocietyID  int     `json:"societyId"`
	MemberType *string `json:"memberType,omitempty"`
}

type GetBalanceSheetInput struct {
	SocietyID    int `json:"societyId"`
	ActiveUserID int `json:"activeUserId"`
}

type GetLedgersInput struct {
	SocietyID    int `json:"societyId"`
	ActiveUserID int `json:"activeUserId"`
}

type GetIncomeAndExpenditureInput struct {
	SocietyID    int    `json:"societyId"`
	ActiveUserID int    `json:"activeUserId"`
	FromDate     string `json:"fromDate"`
	ToDate       string `json:"toDate"`
}

func handleGetMemberByID(_ context.Context, _ *mcp.CallToolRequest, input GetMemberByIDInput) (*mcp.CallToolResult, any, error) {
	entries, err := getMemberDirectory(input.SocietyID)
	if err != nil {
		return nil, nil, err
	}

	entry, found := findMemberByID(entries, input.MemberID)
	if !found {
		notFound := GetMemberByIDResponse{
			Tool:      "get_member_by_id",
			Status:    "not_found",
			SocietyID: input.SocietyID,
			MemberID:  input.MemberID,
		}
		return jsonResult(notFound)
	}

	return jsonResult(buildGetMemberByIDResponse(input.SocietyID, entry))
}

func handleGetMembers(_ context.Context, _ *mcp.CallToolRequest, input GetMembersInput) (*mcp.CallToolResult, any, error) {
	entries, err := getMemberDirectory(input.SocietyID)
	if err != nil {
		return nil, nil, err
	}

	return jsonResult(buildGetMembersResponse(input.SocietyID, entries))
}

func handleGetBalanceSheet(_ context.Context, _ *mcp.CallToolRequest, input GetBalanceSheetInput) (*mcp.CallToolResult, any, error) {
	asOn, asOnDate := currentAsOn(time.Now())

	body, err := getBalanceSheet(input.SocietyID, input.ActiveUserID, asOn)
	if err != nil {
		return nil, nil, err
	}

	result, err := buildGetBalanceSheetResponse(input.SocietyID, asOnDate, body)
	if err != nil {
		return nil, nil, err
	}

	return jsonResult(result)
}

func handleGetLedgers(_ context.Context, _ *mcp.CallToolRequest, input GetLedgersInput) (*mcp.CallToolResult, any, error) {
	body, err := getLedgers(input.SocietyID, input.ActiveUserID)
	if err != nil {
		return nil, nil, err
	}

	result, err := buildGetLedgersResponse(input.SocietyID, body)
	if err != nil {
		return nil, nil, err
	}

	return jsonResult(result)
}

func handleGetIncomeAndExpenditure(_ context.Context, _ *mcp.CallToolRequest, input GetIncomeAndExpenditureInput) (*mcp.CallToolResult, any, error) {
	startDate, err := dateUnixMillis(input.FromDate)
	if err != nil {
		return nil, nil, err
	}
	endDate, err := dateUnixMillis(input.ToDate)
	if err != nil {
		return nil, nil, err
	}
	if endDate < startDate {
		return nil, nil, fmt.Errorf("toDate must be on or after fromDate")
	}
	endDate += int64(24*time.Hour/time.Millisecond) - 1

	body, err := getIncomeAndExpenditure(input.SocietyID, input.ActiveUserID, startDate, endDate)
	if err != nil {
		return nil, nil, err
	}

	result, err := buildGetIncomeAndExpenditureResponse(input.SocietyID, input.FromDate, input.ToDate, body)
	if err != nil {
		return nil, nil, err
	}

	return jsonResult(result)
}

func jsonResult(data any) (*mcp.CallToolResult, any, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(bytes)},
		},
	}
	return result, nil, nil
}

func readOnlyToolAnnotations() *mcp.ToolAnnotations {
	destructive := false
	openWorld := true

	return &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		DestructiveHint: &destructive,
		IdempotentHint:  true,
		OpenWorldHint:   &openWorld,
	}
}
