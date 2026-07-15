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

// initiativeCmd represents the initiative command
var initiativeCmd = &cobra.Command{
	Use:   "initiative",
	Short: "Manage Linear initiatives",
	Long: `Manage Linear initiatives including listing, viewing, creating, updating, and
attaching/detaching projects.

Body content for create/update is read from --content, --body-file, or piped
stdin.

Examples:
  linctl initiative list
  linctl initiative get INITIATIVE-ID
  linctl initiative create --name "Platform 2026" --status Active
  linctl initiative update INITIATIVE-ID --status Completed
  linctl initiative add-project INITIATIVE-ID --project PROJECT-ID`,
}

var initiativeListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List initiatives",
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)

		query, _ := cmd.Flags().GetString("query")
		limit, _ := cmd.Flags().GetInt("limit")
		orderBy, err := resolveDocOrderBy(cmd)
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}

		filter := map[string]interface{}{}
		if query != "" {
			filter["name"] = map[string]interface{}{"containsIgnoreCase": query}
		}
		if len(filter) == 0 {
			filter = nil
		}

		initiatives, err := client.GetInitiatives(context.Background(), filter, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list initiatives: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(initiatives.Nodes)
			return
		}

		if plaintext {
			fmt.Println("# Initiatives")
			for _, ini := range initiatives.Nodes {
				fmt.Printf("## %s\n", ini.Name)
				fmt.Printf("- **ID**: %s\n", ini.ID)
				fmt.Printf("- **Status**: %s\n", ini.Status)
				fmt.Printf("- **Owner**: %s\n", ownerName(ini.Owner))
				if ini.TargetDate != nil {
					fmt.Printf("- **Target Date**: %s\n", *ini.TargetDate)
				}
				fmt.Printf("- **URL**: %s\n", ini.URL)
			}
			fmt.Printf("\nTotal: %d initiatives\n", len(initiatives.Nodes))
			return
		}

		headers := []string{"Name", "Status", "Owner", "Target", "ID"}
		rows := [][]string{}
		for _, ini := range initiatives.Nodes {
			target := "-"
			if ini.TargetDate != nil {
				target = *ini.TargetDate
			}
			rows = append(rows, []string{
				truncateString(ini.Name, 35),
				initiativeStatusColor(ini.Status).Sprint(ini.Status),
				ownerName(ini.Owner),
				target,
				ini.ID,
			})
		}
		output.Table(output.TableData{Headers: headers, Rows: rows}, plaintext, jsonOut)
		fmt.Printf("\n%s %d initiatives\n", color.New(color.FgGreen).Sprint("✓"), len(initiatives.Nodes))
	},
}

var initiativeGetCmd = &cobra.Command{
	Use:     "get INITIATIVE-ID",
	Aliases: []string{"show"},
	Short:   "Get initiative details",
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

		ini, err := client.GetInitiative(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(ini)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n\n", ini.Name)
			if ini.Description != "" {
				fmt.Printf("## Description\n%s\n\n", ini.Description)
			}
			fmt.Printf("- **ID**: %s\n", ini.ID)
			fmt.Printf("- **Status**: %s\n", ini.Status)
			fmt.Printf("- **Owner**: %s\n", ownerName(ini.Owner))
			if ini.TargetDate != nil {
				fmt.Printf("- **Target Date**: %s\n", *ini.TargetDate)
			}
			fmt.Printf("- **URL**: %s\n", ini.URL)
			if ini.Projects != nil && len(ini.Projects.Nodes) > 0 {
				fmt.Printf("\n## Projects (%d)\n", len(ini.Projects.Nodes))
				for _, p := range ini.Projects.Nodes {
					fmt.Printf("- %s (%s) — %s\n", p.Name, p.State, p.ID)
				}
			}
			if ini.Content != "" {
				fmt.Printf("\n## Content\n%s\n", ini.Content)
			}
			return
		}

		fmt.Println()
		fmt.Printf("%s %s\n", color.New(color.FgCyan, color.Bold).Sprint("🎯 Initiative:"), ini.Name)
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), ini.ID)
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status:"), initiativeStatusColor(ini.Status).Sprint(ini.Status))
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Owner:"), ownerName(ini.Owner))
		if ini.TargetDate != nil {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Target Date:"), *ini.TargetDate)
		}
		if ini.Description != "" {
			fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Description:"), ini.Description)
		}
		if ini.Projects != nil && len(ini.Projects.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Projects:"))
			for _, p := range ini.Projects.Nodes {
				fmt.Printf("  • %s (%s)\n", p.Name, color.New(color.FgCyan).Sprint(p.State))
			}
		}
		if ini.URL != "" {
			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("URL:"), color.New(color.FgBlue, color.Underline).Sprint(ini.URL))
		}
		fmt.Println()
	},
}

var initiativeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an initiative",
	Long: `Create a new initiative. Body content is read from --content, --body-file, or
piped stdin.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			output.Error("--name is required", plaintext, jsonOut)
			os.Exit(1)
		}

		status, _ := cmd.Flags().GetString("status")
		if status != "" && !isValidInitiativeStatus(status) {
			output.Error("--status must be one of: Proposed, Planned, Active, Completed, Canceled", plaintext, jsonOut)
			os.Exit(1)
		}

		content, hasContent, err := resolveBody(cmd, "content", "body-file")
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
		ctx := context.Background()

		input := map[string]interface{}{"name": name}
		if hasContent {
			input["content"] = content
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			input["description"] = v
		}
		if status != "" {
			input["status"] = status
		}
		if v, _ := cmd.Flags().GetString("target-date"); v != "" {
			input["targetDate"] = v
		}
		if v, _ := cmd.Flags().GetString("owner"); v != "" {
			ownerID, err := resolveLeadID(ctx, client, v)
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["ownerId"] = ownerID
		}

		ini, err := client.CreateInitiative(ctx, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(ini)
			return
		}
		output.Success(fmt.Sprintf("Created initiative %q (%s)", ini.Name, ini.ID), plaintext, jsonOut)
		if !plaintext {
			fmt.Printf("URL: %s\n", ini.URL)
		}
	},
}

var initiativeUpdateCmd = &cobra.Command{
	Use:   "update INITIATIVE-ID",
	Short: "Update an initiative",
	Long: `Update an initiative's name, description, body, owner, status, or target
date. Body content is read from --content, --body-file, or piped stdin. Only
provided fields change.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		content, hasContent, err := resolveBody(cmd, "content", "body-file")
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
		ctx := context.Background()

		input := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			input["name"] = v
		}
		if hasContent {
			input["content"] = content
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			input["description"] = v
		}
		if cmd.Flags().Changed("status") {
			v, _ := cmd.Flags().GetString("status")
			if !isValidInitiativeStatus(v) {
				output.Error("--status must be one of: Proposed, Planned, Active, Completed, Canceled", plaintext, jsonOut)
				os.Exit(1)
			}
			input["status"] = v
		}
		if cmd.Flags().Changed("target-date") {
			v, _ := cmd.Flags().GetString("target-date")
			input["targetDate"] = v
		}
		if cmd.Flags().Changed("owner") {
			v, _ := cmd.Flags().GetString("owner")
			ownerID, err := resolveLeadID(ctx, client, v)
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["ownerId"] = ownerID
		}

		if len(input) == 0 {
			output.Error("Nothing to update. Provide at least one field flag.", plaintext, jsonOut)
			os.Exit(1)
		}

		ini, err := client.UpdateInitiative(ctx, args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(ini)
			return
		}
		output.Success(fmt.Sprintf("Updated initiative %q (%s)", ini.Name, ini.ID), plaintext, jsonOut)
	},
}

var initiativeAddProjectCmd = &cobra.Command{
	Use:   "add-project INITIATIVE-ID --project PROJECT-ID",
	Short: "Attach a project to an initiative",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			output.Error("--project is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)

		if err := client.AddProjectToInitiative(context.Background(), args[0], projectID); err != nil {
			output.Error(fmt.Sprintf("Failed to attach project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Attached project %s to initiative %s", projectID, args[0]), plaintext, jsonOut)
	},
}

var initiativeRemoveProjectCmd = &cobra.Command{
	Use:   "remove-project INITIATIVE-ID --project PROJECT-ID",
	Short: "Detach a project from an initiative",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			output.Error("--project is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)

		if err := client.RemoveProjectFromInitiative(context.Background(), args[0], projectID); err != nil {
			output.Error(fmt.Sprintf("Failed to detach project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Detached project %s from initiative %s", projectID, args[0]), plaintext, jsonOut)
	},
}

// ownerName returns a display string for an optional initiative owner.
func ownerName(u *api.User) string {
	if u == nil {
		return "Unassigned"
	}
	return u.Name
}

// initiativeStatusColor returns a color for an initiative status value.
func initiativeStatusColor(status string) *color.Color {
	switch status {
	case "Proposed":
		return color.New(color.FgCyan)
	case "Planned":
		return color.New(color.FgBlue)
	case "Active":
		return color.New(color.FgGreen)
	case "Completed":
		return color.New(color.FgGreen)
	case "Canceled":
		return color.New(color.FgRed)
	default:
		return color.New(color.FgWhite)
	}
}

// isValidInitiativeStatus reports whether s is a valid InitiativeStatus value.
func isValidInitiativeStatus(s string) bool {
	switch s {
	case "Proposed", "Planned", "Active", "Completed", "Canceled":
		return true
	default:
		return false
	}
}

func init() {
	rootCmd.AddCommand(initiativeCmd)
	initiativeCmd.AddCommand(initiativeListCmd)
	initiativeCmd.AddCommand(initiativeGetCmd)
	initiativeCmd.AddCommand(initiativeCreateCmd)
	initiativeCmd.AddCommand(initiativeUpdateCmd)
	initiativeCmd.AddCommand(initiativeAddProjectCmd)
	initiativeCmd.AddCommand(initiativeRemoveProjectCmd)

	initiativeListCmd.Flags().StringP("query", "q", "", "Filter by name (case-insensitive substring)")
	initiativeListCmd.Flags().IntP("limit", "l", 50, "Maximum number of initiatives to return")
	initiativeListCmd.Flags().StringP("sort", "o", "updatedAt", "Sort order: updated (default), created")

	initiativeCreateCmd.Flags().String("name", "", "Initiative name (required)")
	initiativeCreateCmd.Flags().String("description", "", "Short initiative description")
	initiativeCreateCmd.Flags().String("content", "", "Initiative body content (markdown)")
	initiativeCreateCmd.Flags().String("body-file", "", "Path to a file containing the initiative body")
	initiativeCreateCmd.Flags().String("owner", "", "Initiative owner email")
	initiativeCreateCmd.Flags().String("status", "", "Status: Proposed, Planned, Active, Completed, Canceled")
	initiativeCreateCmd.Flags().String("target-date", "", "Target date (YYYY-MM-DD)")

	initiativeUpdateCmd.Flags().String("name", "", "New initiative name")
	initiativeUpdateCmd.Flags().String("description", "", "New short description")
	initiativeUpdateCmd.Flags().String("content", "", "New initiative body content (markdown)")
	initiativeUpdateCmd.Flags().String("body-file", "", "Path to a file containing the new initiative body")
	initiativeUpdateCmd.Flags().String("owner", "", "New initiative owner email")
	initiativeUpdateCmd.Flags().String("status", "", "New status: Proposed, Planned, Active, Completed, Canceled")
	initiativeUpdateCmd.Flags().String("target-date", "", "New target date (YYYY-MM-DD)")

	initiativeAddProjectCmd.Flags().String("project", "", "Project ID to attach (required)")
	initiativeRemoveProjectCmd.Flags().String("project", "", "Project ID to detach (required)")
}
