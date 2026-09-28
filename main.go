package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	loadEnvFile(".env")

	headerName := os.Getenv("NEIGHBIUM_API_HEADER_NAME")
	apiKey := os.Getenv("NEIGHBIUM_API_KEY")
	authorization := os.Getenv("NEIGHBIUM_AUTHORIZATION")

	if headerName == "" || apiKey == "" || authorization == "" {
		log.Fatal("Neighbium API credentials are not configured.")
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "neighbium-mcp",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_member_by_id",
		Description: "Retrieve one society member by member ID.",
		Annotations: readOnlyToolAnnotations(),
	}, handleGetMemberByID)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_members",
		Description: "Retrieve the members of a housing society. Can optionally filter members by member type such as Owner or Tenant.",
		Annotations: readOnlyToolAnnotations(),
	}, handleGetMembers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_balance_sheet",
		Description: "Retrieve the society balance sheet as of the current date and time.",
		Annotations: readOnlyToolAnnotations(),
	}, handleGetBalanceSheet)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_ledgers",
		Description: "Retrieve ledger balances for a society.",
		Annotations: readOnlyToolAnnotations(),
	}, handleGetLedgers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_income_and_expenditure",
		Description: "Retrieve society income and expenses for an inclusive date range. Dates are interpreted as UTC.",
		Annotations: readOnlyToolAnnotations(),
	}, handleGetIncomeAndExpenditure)

	log.Print("Neighbium MCP server started")

	err := server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		log.Fatal(err)
	}
}
