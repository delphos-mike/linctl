package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/delphos-mike/linctl/pkg/api"
	"github.com/delphos-mike/linctl/pkg/auth"
	"github.com/delphos-mike/linctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// graphqlCmd is a raw GraphQL passthrough that reuses linctl authentication.
// It is the documented escape hatch for operations linctl does not yet wrap,
// and a convenient introspection tool.
var graphqlCmd = &cobra.Command{
	Use:     "graphql",
	Aliases: []string{"gql"},
	Short:   "Execute a raw GraphQL query/mutation against Linear",
	Long: `Execute a raw GraphQL query or mutation against Linear's API using
linctl's stored authentication.

The query is read from --query, --query-file, or stdin (in that order of
precedence). Variables may be supplied inline with repeatable --var key=value
pairs and/or a JSON object via --vars-file. Inline --var values are parsed as
JSON when possible (so 50 is a number, true is a boolean, {"a":1} is an object)
and fall back to a plain string otherwise.

The full response (data and any GraphQL errors) is printed as JSON. The command
exits non-zero if the response contains GraphQL errors.

Examples:
  linctl graphql -q 'query { viewer { id name } }'
  linctl graphql --query-file mutation.graphql --var id=abc123 --var first=50
  echo 'query { viewer { email } }' | linctl graphql
  linctl graphql -q 'query($f: DocumentFilter){ documents(filter:$f){ nodes { id title } } }' \
    --var 'f={"project":{"id":{"eq":"PROJECT-UUID"}}}'`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		query, err := resolveGraphQLQuery(cmd)
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}
		if strings.TrimSpace(query) == "" {
			output.Error("No query provided. Use --query, --query-file, or pipe via stdin.", plaintext, jsonOut)
			os.Exit(1)
		}

		variables, err := resolveGraphQLVariables(cmd)
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		resp, err := client.ExecuteRaw(context.Background(), query, variables)
		if err != nil {
			output.Error(fmt.Sprintf("Request failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Print the raw data payload (or the full response when errors exist).
		if len(resp.Errors) > 0 {
			output.JSON(map[string]interface{}{
				"data":   json.RawMessage(resp.Data),
				"errors": resp.Errors,
			})
			os.Exit(1)
		}

		if resp.Data != nil {
			output.JSON(json.RawMessage(resp.Data))
		} else {
			output.JSON(map[string]interface{}{})
		}
	},
}

// resolveGraphQLQuery reads the query from --query, --query-file, or stdin.
func resolveGraphQLQuery(cmd *cobra.Command) (string, error) {
	if q, _ := cmd.Flags().GetString("query"); q != "" {
		return q, nil
	}
	if path, _ := cmd.Flags().GetString("query-file"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read query file: %w", err)
		}
		return string(data), nil
	}
	// Fall back to stdin when it is piped (not a TTY).
	if isStdinPiped() {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read stdin: %w", err)
		}
		return string(data), nil
	}
	return "", nil
}

// resolveGraphQLVariables merges --vars-file (JSON object) with repeatable
// --var key=value pairs. Inline values take precedence over the file.
func resolveGraphQLVariables(cmd *cobra.Command) (map[string]interface{}, error) {
	variables := map[string]interface{}{}

	if path, _ := cmd.Flags().GetString("vars-file"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read vars file: %w", err)
		}
		if err := json.Unmarshal(data, &variables); err != nil {
			return nil, fmt.Errorf("vars file is not a valid JSON object: %w", err)
		}
	}

	pairs, _ := cmd.Flags().GetStringArray("var")
	for _, pair := range pairs {
		key, rawValue, found := strings.Cut(pair, "=")
		if !found {
			return nil, fmt.Errorf("invalid --var %q, expected key=value", pair)
		}
		variables[key] = parseVarValue(rawValue)
	}

	if len(variables) == 0 {
		return nil, nil
	}
	return variables, nil
}

// parseVarValue interprets a raw --var value as JSON when possible, otherwise
// returns it as a plain string.
func parseVarValue(raw string) interface{} {
	var parsed interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
		return parsed
	}
	return raw
}

// isStdinPiped reports whether stdin has piped (non-terminal) data.
func isStdinPiped() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}

func init() {
	rootCmd.AddCommand(graphqlCmd)

	graphqlCmd.Flags().StringP("query", "q", "", "GraphQL query or mutation string")
	graphqlCmd.Flags().String("query-file", "", "Path to a file containing the GraphQL query")
	graphqlCmd.Flags().StringArray("var", []string{}, "Variable as key=value (repeatable; values parsed as JSON when possible)")
	graphqlCmd.Flags().String("vars-file", "", "Path to a JSON file of variables")
}
