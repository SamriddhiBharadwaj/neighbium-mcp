package main

import (
	"context"
	"encoding/json"
	"fmt"

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
	SocietyID int   `json:"societyId"`
	AsOnDate  int64 `json:"asOnDate"`
}

type GetLedgersBalanceInput struct {
	SocietyID int   `json:"societyId"`
	AsOnDate  int64 `json:"asOnDate"`
}

type GetIncomeAndExpenditureInput struct {
	SocietyID int   `json:"societyId"`
	FromDate  int64 `json:"fromDate"`
	ToDate    int64 `json:"toDate"`
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
	body, err := getBalanceSheet(input.SocietyID, input.AsOnDate)
	if err != nil {
		return nil, nil, err
	}

	result, err := buildGetBalanceSheetResponse(input.SocietyID, input.AsOnDate, body)
	if err != nil {
		return nil, nil, err
	}

	return jsonResult(result)
}

func handleGetLedgersBalance(_ context.Context, _ *mcp.CallToolRequest, input GetLedgersBalanceInput) (*mcp.CallToolResult, any, error) {
	body, err := getLedgersBalance(input.SocietyID, input.AsOnDate)
	if err != nil {
		return nil, nil, err
	}

	result, err := buildGetLedgersBalanceResponse(input.SocietyID, body)
	if err != nil {
		return nil, nil, err
	}

	return jsonResult(result)
}

func handleGetIncomeAndExpenditure(_ context.Context, _ *mcp.CallToolRequest, input GetIncomeAndExpenditureInput) (*mcp.CallToolResult, any, error) {
	if input.ToDate < input.FromDate {
		return nil, nil, fmt.Errorf("toDate must be on or after fromDate")
	}

	body, err := getIncomeAndExpenditure(input.SocietyID, input.FromDate, input.ToDate)
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
