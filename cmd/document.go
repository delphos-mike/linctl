package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/delphos-mike/linctl/pkg/api"
	"github.com/delphos-mike/linctl/pkg/auth"
	"github.com/delphos-mike/linctl/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// documentCmd represents the document command
var documentCmd = &cobra.Command{
	Use:     "document",
	Aliases: []string{"doc"},
	Short:   "Manage Linear documents",
	Long: `Manage Linear documents including listing, viewing, creating, updating, and deleting.

Documents may be resident in a project, an initiative, or the workspace. Body
content for create/update is read from --body, --body-file, or piped stdin.

Examples:
  linctl document list --project PROJECT-ID
  linctl document list --query "Project State"
  linctl document get DOCUMENT-ID
  linctl document create --title "Project State" --project PROJECT-ID --body-file state.md
  cat state.md | linctl document update DOCUMENT-ID
  linctl document delete DOCUMENT-ID`,
}

var documentListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List documents",
	Long: `List documents, optionally scoped to a project or initiative.

--query matches document titles (case-insensitive substring). Exact-title
selection and "newest wins" tie-breaking are left to the caller; results are
ordered by the --sort key (default updatedAt) and each row exposes the
document's project/initiative association.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)

		projectID, _ := cmd.Flags().GetString("project")
		initiativeID, _ := cmd.Flags().GetString("initiative")
		query, _ := cmd.Flags().GetString("query")
		limit, _ := cmd.Flags().GetInt("limit")

		orderBy, err := resolveDocOrderBy(cmd)
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}

		filter := map[string]interface{}{}
		if projectID != "" {
			filter["project"] = map[string]interface{}{"id": map[string]interface{}{"eq": projectID}}
		}
		if initiativeID != "" {
			filter["initiative"] = map[string]interface{}{"id": map[string]interface{}{"eq": initiativeID}}
		}
		if query != "" {
			filter["title"] = map[string]interface{}{"containsIgnoreCase": query}
		}
		if len(filter) == 0 {
			filter = nil
		}

		docs, err := client.GetDocuments(context.Background(), filter, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list documents: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(docs.Nodes)
			return
		}

		if plaintext {
			fmt.Println("# Documents")
			for _, doc := range docs.Nodes {
				fmt.Printf("## %s\n", doc.Title)
				fmt.Printf("- **ID**: %s\n", doc.ID)
				fmt.Printf("- **Association**: %s\n", documentAssociation(doc))
				fmt.Printf("- **Updated**: %s\n", doc.UpdatedAt.Format("2006-01-02 15:04"))
				fmt.Printf("- **URL**: %s\n", doc.URL)
			}
			fmt.Printf("\nTotal: %d documents\n", len(docs.Nodes))
			return
		}

		headers := []string{"Title", "Association", "Updated", "ID"}
		rows := [][]string{}
		for _, doc := range docs.Nodes {
			rows = append(rows, []string{
				truncateString(doc.Title, 40),
				documentAssociation(doc),
				doc.UpdatedAt.Format("2006-01-02"),
				doc.ID,
			})
		}
		output.Table(output.TableData{Headers: headers, Rows: rows}, plaintext, jsonOut)
		fmt.Printf("\n%s %d documents\n", color.New(color.FgGreen).Sprint("✓"), len(docs.Nodes))
	},
}

var documentGetCmd = &cobra.Command{
	Use:     "get DOCUMENT-ID",
	Aliases: []string{"show"},
	Short:   "Get document details",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)

		doc, err := client.GetDocument(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(doc)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n\n", doc.Title)
			fmt.Printf("- **ID**: %s\n", doc.ID)
			fmt.Printf("- **Association**: %s\n", documentAssociation(*doc))
			if doc.Creator != nil {
				fmt.Printf("- **Creator**: %s\n", doc.Creator.Name)
			}
			fmt.Printf("- **Created**: %s\n", doc.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Printf("- **Updated**: %s\n", doc.UpdatedAt.Format("2006-01-02 15:04"))
			fmt.Printf("- **URL**: %s\n\n", doc.URL)
			fmt.Printf("## Content\n%s\n", doc.Content)
			return
		}

		fmt.Println()
		fmt.Printf("%s %s\n", color.New(color.FgCyan, color.Bold).Sprint("📄 Document:"), doc.Title)
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), doc.ID)
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Association:"), documentAssociation(*doc))
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Updated:"), doc.UpdatedAt.Format("2006-01-02 15:04"))
		if doc.URL != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("URL:"), color.New(color.FgBlue, color.Underline).Sprint(doc.URL))
		}
		if doc.Content != "" {
			fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Content:"), doc.Content)
		}
		fmt.Println()
	},
}

var documentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a document",
	Long: `Create a new document. Body content is read from --body, --body-file, or
piped stdin. Attach the document to a project or initiative with --project or
--initiative; omit both to create a workspace-resident document.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			output.Error("--title is required", plaintext, jsonOut)
			os.Exit(1)
		}

		body, hasBody, err := resolveBody(cmd, "body", "body-file")
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

		input := map[string]interface{}{"title": title}
		if hasBody {
			input["content"] = body
		}
		if v, _ := cmd.Flags().GetString("project"); v != "" {
			input["projectId"] = v
		}
		if v, _ := cmd.Flags().GetString("initiative"); v != "" {
			input["initiativeId"] = v
		}
		if v, _ := cmd.Flags().GetString("icon"); v != "" {
			input["icon"] = v
		}
		if v, _ := cmd.Flags().GetString("color"); v != "" {
			input["color"] = v
		}

		doc, err := client.CreateDocument(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(doc)
			return
		}
		output.Success(fmt.Sprintf("Created document %q (%s)", doc.Title, doc.ID), plaintext, jsonOut)
		if !plaintext {
			fmt.Printf("URL: %s\n", doc.URL)
		}
	},
}

var documentUpdateCmd = &cobra.Command{
	Use:   "update DOCUMENT-ID",
	Short: "Update a document",
	Long: `Update a document's title, body, icon, color, or association. Body content
is read from --body, --body-file, or piped stdin. Only provided fields change.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		body, hasBody, err := resolveBody(cmd, "body", "body-file")
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{}
		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			input["title"] = v
		}
		if hasBody {
			input["content"] = body
		}
		if cmd.Flags().Changed("project") {
			v, _ := cmd.Flags().GetString("project")
			input["projectId"] = v
		}
		if cmd.Flags().Changed("initiative") {
			v, _ := cmd.Flags().GetString("initiative")
			input["initiativeId"] = v
		}
		if cmd.Flags().Changed("icon") {
			v, _ := cmd.Flags().GetString("icon")
			input["icon"] = v
		}
		if cmd.Flags().Changed("color") {
			v, _ := cmd.Flags().GetString("color")
			input["color"] = v
		}

		if len(input) == 0 {
			output.Error("Nothing to update. Provide at least one of --title, --body/--body-file/stdin, --project, --initiative, --icon, --color.", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)

		doc, err := client.UpdateDocument(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(doc)
			return
		}
		output.Success(fmt.Sprintf("Updated document %q (%s)", doc.Title, doc.ID), plaintext, jsonOut)
	},
}

var documentDeleteCmd = &cobra.Command{
	Use:     "delete DOCUMENT-ID",
	Aliases: []string{"rm"},
	Short:   "Delete (trash) a document",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)

		if err := client.DeleteDocument(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Deleted document %s", args[0]), plaintext, jsonOut)
	},
}

// documentAssociation returns a human-readable description of where a document
// lives: a project, an initiative, or the workspace.
func documentAssociation(doc api.Document) string {
	if doc.Project != nil && doc.Project.Name != "" {
		return fmt.Sprintf("project: %s", doc.Project.Name)
	}
	if doc.Initiative != nil && doc.Initiative.Name != "" {
		return fmt.Sprintf("initiative: %s", doc.Initiative.Name)
	}
	return "workspace"
}

// resolveDocOrderBy maps the --sort flag to a Linear PaginationOrderBy value.
func resolveDocOrderBy(cmd *cobra.Command) (string, error) {
	sortBy, _ := cmd.Flags().GetString("sort")
	switch sortBy {
	case "", "updated", "updatedAt":
		return "updatedAt", nil
	case "created", "createdAt":
		return "createdAt", nil
	default:
		return "", fmt.Errorf("invalid sort option: %s. Valid options are: updated, created", sortBy)
	}
}

func init() {
	rootCmd.AddCommand(documentCmd)
	documentCmd.AddCommand(documentListCmd)
	documentCmd.AddCommand(documentGetCmd)
	documentCmd.AddCommand(documentCreateCmd)
	documentCmd.AddCommand(documentUpdateCmd)
	documentCmd.AddCommand(documentDeleteCmd)

	documentListCmd.Flags().String("project", "", "Filter by project ID (UUID)")
	documentListCmd.Flags().String("initiative", "", "Filter by initiative ID (UUID)")
	documentListCmd.Flags().StringP("query", "q", "", "Filter by title (case-insensitive substring)")
	documentListCmd.Flags().IntP("limit", "l", 50, "Maximum number of documents to return")
	documentListCmd.Flags().StringP("sort", "o", "updatedAt", "Sort order: updated (default), created")

	documentCreateCmd.Flags().String("title", "", "Document title (required)")
	documentCreateCmd.Flags().String("body", "", "Document body content (markdown)")
	documentCreateCmd.Flags().String("body-file", "", "Path to a file containing the document body")
	documentCreateCmd.Flags().String("project", "", "Attach to project ID (UUID)")
	documentCreateCmd.Flags().String("initiative", "", "Attach to initiative ID (UUID)")
	documentCreateCmd.Flags().String("icon", "", "Document icon")
	documentCreateCmd.Flags().String("color", "", "Document color (hex)")

	documentUpdateCmd.Flags().String("title", "", "New document title")
	documentUpdateCmd.Flags().String("body", "", "New document body content (markdown)")
	documentUpdateCmd.Flags().String("body-file", "", "Path to a file containing the new document body")
	documentUpdateCmd.Flags().String("project", "", "Move to project ID (UUID)")
	documentUpdateCmd.Flags().String("initiative", "", "Move to initiative ID (UUID)")
	documentUpdateCmd.Flags().String("icon", "", "New document icon")
	documentUpdateCmd.Flags().String("color", "", "New document color (hex)")
}
