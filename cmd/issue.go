package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/delphos-mike/linctl/pkg/api"
	"github.com/delphos-mike/linctl/pkg/auth"
	"github.com/delphos-mike/linctl/pkg/output"
	"github.com/delphos-mike/linctl/pkg/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Linear issues",
	Long: `Create, list, update, and manage Linear issues.

Examples:
  linctl issue list --assignee me --state "In Progress"
  linctl issue ls -a me -s "In Progress"
  linctl issue list --include-completed  # Show all issues including completed
  linctl issue list --newer-than 3_weeks_ago  # Show issues from last 3 weeks
  linctl issue search "login bug" --team ENG
  linctl issue get LIN-123
  linctl issue create --title "Bug fix" --team ENG
  linctl issue relate LIN-123 --blocks LIN-456
  linctl issue unrelate LIN-123 --blocks LIN-456`,
}

var issueListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List issues",
	Long:    `List Linear issues with optional filtering.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Build filter from flags
		filter := buildIssueFilter(cmd)

		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 50
		}

		// Get sort option
		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				// Use empty string for Linear's default sort
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		issues, err := client.GetIssues(context.Background(), filter, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		renderIssueCollection(issues, plaintext, jsonOut, "No issues found", "issues", "# Issues")
	},
}

func renderIssueCollection(issues *api.Issues, plaintext, jsonOut bool, emptyMessage, summaryLabel, plaintextTitle string) {
	if len(issues.Nodes) == 0 {
		output.Info(emptyMessage, plaintext, jsonOut)
		return
	}

	if jsonOut {
		output.JSON(issues.Nodes)
		return
	}

	if plaintext {
		fmt.Println(plaintextTitle)
		for _, issue := range issues.Nodes {
			fmt.Printf("## %s\n", issue.Title)
			fmt.Printf("- **ID**: %s\n", issue.Identifier)
			if issue.State != nil {
				fmt.Printf("- **State**: %s\n", issue.State.Name)
			}
			if issue.Assignee != nil {
				fmt.Printf("- **Assignee**: %s\n", issue.Assignee.Name)
			} else {
				fmt.Printf("- **Assignee**: Unassigned\n")
			}
			if issue.Team != nil {
				fmt.Printf("- **Team**: %s\n", issue.Team.Key)
			}
			fmt.Printf("- **Created**: %s\n", issue.CreatedAt.Format("2006-01-02"))
			fmt.Printf("- **URL**: %s\n", issue.URL)
			if issue.Description != "" {
				fmt.Printf("- **Description**: %s\n", issue.Description)
			}
			fmt.Println()
		}
		fmt.Printf("\nTotal: %d %s\n", len(issues.Nodes), summaryLabel)
		return
	}

	headers := []string{"Title", "State", "Assignee", "Team", "Created", "URL"}
	rows := make([][]string, len(issues.Nodes))

	for i, issue := range issues.Nodes {
		assignee := "Unassigned"
		if issue.Assignee != nil {
			assignee = issue.Assignee.Name
		}

		team := ""
		if issue.Team != nil {
			team = issue.Team.Key
		}

		state := ""
		if issue.State != nil {
			state = issue.State.Name
			var stateColor *color.Color
			switch issue.State.Type {
			case "triage":
				stateColor = color.New(color.FgMagenta)
			case "backlog":
				stateColor = color.New(color.FgCyan)
			case "unstarted":
				stateColor = color.New(color.FgWhite)
			case "started":
				stateColor = color.New(color.FgBlue)
			case "completed":
				stateColor = color.New(color.FgGreen)
			case "canceled":
				stateColor = color.New(color.FgRed)
			default:
				stateColor = color.New(color.FgWhite)
			}
			state = stateColor.Sprint(state)
		}

		if issue.Assignee == nil {
			assignee = color.New(color.FgYellow).Sprint(assignee)
		}

		rows[i] = []string{
			truncateString(issue.Title, 40),
			state,
			assignee,
			team,
			issue.CreatedAt.Format("2006-01-02"),
			issue.URL,
		}
	}

	tableData := output.TableData{
		Headers: headers,
		Rows:    rows,
	}

	output.Table(tableData, false, false)

	fmt.Printf("\n%s %d %s\n",
		color.New(color.FgGreen).Sprint("✓"),
		len(issues.Nodes),
		summaryLabel)

	if issues.PageInfo.HasNextPage {
		fmt.Printf("%s Use --limit to see more results\n",
			color.New(color.FgYellow).Sprint("ℹ️"))
	}
}

var issueSearchCmd = &cobra.Command{
	Use:     "search [query]",
	Aliases: []string{"find"},
	Short:   "Search issues by keyword",
	Long: `Perform a full-text search across Linear issues.

Examples:
  linctl issue search "payment outage"
  linctl issue search "auth token" --team ENG --include-completed
  linctl issue search "customer:" --json`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			output.Error("Search query is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		filter := buildIssueFilter(cmd)

		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 50
		}

		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		includeArchived, _ := cmd.Flags().GetBool("include-archived")

		issues, err := client.IssueSearch(context.Background(), query, filter, limit, "", orderBy, includeArchived)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to search issues: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		emptyMsg := fmt.Sprintf("No matches found for %q", query)
		renderIssueCollection(issues, plaintext, jsonOut, emptyMsg, "matches", "# Search Results")
	},
}

var issueGetCmd = &cobra.Command{
	Use:     "get [issue-id]",
	Aliases: []string{"show"},
	Short:   "Get issue details",
	Long:    `Get detailed information about a specific issue.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		issue, err := client.GetIssue(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issue)
			return
		}

		if plaintext {
			fmt.Printf("# %s - %s\n\n", issue.Identifier, issue.Title)

			if issue.Description != "" {
				fmt.Printf("## Description\n%s\n\n", issue.Description)
			}

			fmt.Printf("## Core Details\n")
			fmt.Printf("- **ID**: %s\n", issue.Identifier)
			fmt.Printf("- **Number**: %d\n", issue.Number)
			if issue.State != nil {
				fmt.Printf("- **State**: %s (%s)\n", issue.State.Name, issue.State.Type)
				if issue.State.Description != nil && *issue.State.Description != "" {
					fmt.Printf("  - Description: %s\n", *issue.State.Description)
				}
			}
			if issue.Assignee != nil {
				fmt.Printf("- **Assignee**: %s (%s)\n", issue.Assignee.Name, issue.Assignee.Email)
				if issue.Assignee.DisplayName != "" && issue.Assignee.DisplayName != issue.Assignee.Name {
					fmt.Printf("  - Display Name: %s\n", issue.Assignee.DisplayName)
				}
			} else {
				fmt.Printf("- **Assignee**: Unassigned\n")
			}
			if issue.Creator != nil {
				fmt.Printf("- **Creator**: %s (%s)\n", issue.Creator.Name, issue.Creator.Email)
			}
			if issue.Team != nil {
				fmt.Printf("- **Team**: %s (%s)\n", issue.Team.Name, issue.Team.Key)
				if issue.Team.Description != "" {
					fmt.Printf("  - Description: %s\n", issue.Team.Description)
				}
			}
			fmt.Printf("- **Priority**: %s (%d)\n", priorityToString(issue.Priority), issue.Priority)
			if issue.PriorityLabel != "" {
				fmt.Printf("- **Priority Label**: %s\n", issue.PriorityLabel)
			}
			if issue.Estimate != nil {
				fmt.Printf("- **Estimate**: %.1f\n", *issue.Estimate)
			}

			fmt.Printf("\n## Status & Dates\n")
			fmt.Printf("- **Created**: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05"))
			if issue.TriagedAt != nil {
				fmt.Printf("- **Triaged**: %s\n", issue.TriagedAt.Format("2006-01-02 15:04:05"))
			}
			if issue.CompletedAt != nil {
				fmt.Printf("- **Completed**: %s\n", issue.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			if issue.CanceledAt != nil {
				fmt.Printf("- **Canceled**: %s\n", issue.CanceledAt.Format("2006-01-02 15:04:05"))
			}
			if issue.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", issue.ArchivedAt.Format("2006-01-02 15:04:05"))
			}
			if issue.DueDate != nil && *issue.DueDate != "" {
				fmt.Printf("- **Due Date**: %s\n", *issue.DueDate)
			}
			if issue.SnoozedUntilAt != nil {
				fmt.Printf("- **Snoozed Until**: %s\n", issue.SnoozedUntilAt.Format("2006-01-02 15:04:05"))
			}

			fmt.Printf("\n## Technical Details\n")
			fmt.Printf("- **Board Order**: %.2f\n", issue.BoardOrder)
			fmt.Printf("- **Sub-Issue Sort Order**: %.2f\n", issue.SubIssueSortOrder)
			if issue.BranchName != "" {
				fmt.Printf("- **Git Branch**: %s\n", issue.BranchName)
			}
			if issue.CustomerTicketCount > 0 {
				fmt.Printf("- **Customer Ticket Count**: %d\n", issue.CustomerTicketCount)
			}
			if len(issue.PreviousIdentifiers) > 0 {
				fmt.Printf("- **Previous Identifiers**: %s\n", strings.Join(issue.PreviousIdentifiers, ", "))
			}
			if issue.IntegrationSourceType != nil && *issue.IntegrationSourceType != "" {
				fmt.Printf("- **Integration Source**: %s\n", *issue.IntegrationSourceType)
			}
			if issue.ExternalUserCreator != nil {
				fmt.Printf("- **External Creator**: %s (%s)\n", issue.ExternalUserCreator.Name, issue.ExternalUserCreator.Email)
			}
			fmt.Printf("- **URL**: %s\n", issue.URL)

			// Project and Cycle Info
			if issue.Project != nil {
				fmt.Printf("\n## Project\n")
				fmt.Printf("- **Name**: %s\n", issue.Project.Name)
				fmt.Printf("- **State**: %s\n", issue.Project.State)
				fmt.Printf("- **Progress**: %.0f%%\n", issue.Project.Progress*100)
				if issue.Project.Health != "" {
					fmt.Printf("- **Health**: %s\n", issue.Project.Health)
				}
				if issue.Project.Description != "" {
					fmt.Printf("- **Description**: %s\n", issue.Project.Description)
				}
			}

			if issue.Cycle != nil {
				fmt.Printf("\n## Cycle\n")
				fmt.Printf("- **Name**: %s (#%d)\n", issue.Cycle.Name, issue.Cycle.Number)
				if issue.Cycle.Description != nil && *issue.Cycle.Description != "" {
					fmt.Printf("- **Description**: %s\n", *issue.Cycle.Description)
				}
				fmt.Printf("- **Period**: %s to %s\n", issue.Cycle.StartsAt, issue.Cycle.EndsAt)
				fmt.Printf("- **Progress**: %.0f%%\n", issue.Cycle.Progress*100)
				if issue.Cycle.CompletedAt != nil {
					fmt.Printf("- **Completed**: %s\n", issue.Cycle.CompletedAt.Format("2006-01-02"))
				}
			}

			// Labels
			if issue.Labels != nil && len(issue.Labels.Nodes) > 0 {
				fmt.Printf("\n## Labels\n")
				for _, label := range issue.Labels.Nodes {
					fmt.Printf("- %s", label.Name)
					if label.Description != nil && *label.Description != "" {
						fmt.Printf(" - %s", *label.Description)
					}
					fmt.Println()
				}
			}

			// Subscribers
			if issue.Subscribers != nil && len(issue.Subscribers.Nodes) > 0 {
				fmt.Printf("\n## Subscribers\n")
				for _, subscriber := range issue.Subscribers.Nodes {
					fmt.Printf("- %s (%s)\n", subscriber.Name, subscriber.Email)
				}
			}

			// Relations
			if issue.Relations != nil && len(issue.Relations.Nodes) > 0 {
				fmt.Printf("\n## Related Issues\n")
				for _, relation := range issue.Relations.Nodes {
					if relation.RelatedIssue != nil {
						relationType := relation.Type
						switch relationType {
						case "blocks":
							relationType = "Blocks"
						case "blocked":
							relationType = "Blocked by"
						case "related":
							relationType = "Related to"
						case "duplicate":
							relationType = "Duplicate of"
						}
						fmt.Printf("- %s: %s - %s", relationType, relation.RelatedIssue.Identifier, relation.RelatedIssue.Title)
						if relation.RelatedIssue.State != nil {
							fmt.Printf(" [%s]", relation.RelatedIssue.State.Name)
						}
						fmt.Println()
					}
				}
			}

			// Reactions
			if len(issue.Reactions) > 0 {
				fmt.Printf("\n## Reactions\n")
				reactionMap := make(map[string][]string)
				for _, reaction := range issue.Reactions {
					reactionMap[reaction.Emoji] = append(reactionMap[reaction.Emoji], reaction.User.Name)
				}
				for emoji, users := range reactionMap {
					fmt.Printf("- %s: %s\n", emoji, strings.Join(users, ", "))
				}
			}

			// Show parent issue if this is a sub-issue
			if issue.Parent != nil {
				fmt.Printf("\n## Parent Issue\n")
				fmt.Printf("- %s: %s\n", issue.Parent.Identifier, issue.Parent.Title)
			}

			// Show sub-issues if any
			if issue.Children != nil && len(issue.Children.Nodes) > 0 {
				fmt.Printf("\n## Sub-issues\n")
				for _, child := range issue.Children.Nodes {
					stateStr := ""
					if child.State != nil {
						switch child.State.Type {
						case "completed", "done":
							stateStr = "[x]"
						case "started", "in_progress":
							stateStr = "[~]"
						case "canceled":
							stateStr = "[-]"
						default:
							stateStr = "[ ]"
						}
					} else {
						stateStr = "[ ]"
					}

					assignee := "Unassigned"
					if child.Assignee != nil {
						assignee = child.Assignee.Name
					}

					fmt.Printf("- %s %s: %s (%s)\n", stateStr, child.Identifier, child.Title, assignee)
				}
			}

			// Show attachments if any
			if issue.Attachments != nil && len(issue.Attachments.Nodes) > 0 {
				fmt.Printf("\n## Attachments\n")
				for _, attachment := range issue.Attachments.Nodes {
					fmt.Printf("- [%s](%s)\n", attachment.Title, attachment.URL)
				}
			}

			// Show recent comments if any
			if issue.Comments != nil && len(issue.Comments.Nodes) > 0 {
				fmt.Printf("\n## Recent Comments\n")
				for _, comment := range issue.Comments.Nodes {
					fmt.Printf("\n### %s - %s\n", comment.User.Name, comment.CreatedAt.Format("2006-01-02 15:04"))
					if comment.EditedAt != nil {
						fmt.Printf("*(edited %s)*\n", comment.EditedAt.Format("2006-01-02 15:04"))
					}
					fmt.Printf("%s\n", comment.Body)
					if comment.Children != nil && len(comment.Children.Nodes) > 0 {
						for _, reply := range comment.Children.Nodes {
							fmt.Printf("\n  **Reply from %s**: %s\n", reply.User.Name, reply.Body)
						}
					}
				}
				fmt.Printf("\n> Use `linctl comment list %s` to see all comments\n", issue.Identifier)
			}

			// Show history
			if issue.History != nil && len(issue.History.Nodes) > 0 {
				fmt.Printf("\n## Recent History\n")
				for _, entry := range issue.History.Nodes {
					fmt.Printf("\n- **%s** by %s", entry.CreatedAt.Format("2006-01-02 15:04"), entry.Actor.Name)
					changes := []string{}

					if entry.FromState != nil && entry.ToState != nil {
						changes = append(changes, fmt.Sprintf("State: %s → %s", entry.FromState.Name, entry.ToState.Name))
					}
					if entry.FromAssignee != nil && entry.ToAssignee != nil {
						changes = append(changes, fmt.Sprintf("Assignee: %s → %s", entry.FromAssignee.Name, entry.ToAssignee.Name))
					} else if entry.FromAssignee != nil && entry.ToAssignee == nil {
						changes = append(changes, fmt.Sprintf("Unassigned from %s", entry.FromAssignee.Name))
					} else if entry.FromAssignee == nil && entry.ToAssignee != nil {
						changes = append(changes, fmt.Sprintf("Assigned to %s", entry.ToAssignee.Name))
					}
					if entry.FromPriority != nil && entry.ToPriority != nil {
						changes = append(changes, fmt.Sprintf("Priority: %s → %s", priorityToString(*entry.FromPriority), priorityToString(*entry.ToPriority)))
					}
					if entry.FromTitle != nil && entry.ToTitle != nil {
						changes = append(changes, fmt.Sprintf("Title: \"%s\" → \"%s\"", *entry.FromTitle, *entry.ToTitle))
					}
					if entry.FromCycle != nil && entry.ToCycle != nil {
						changes = append(changes, fmt.Sprintf("Cycle: %s → %s", entry.FromCycle.Name, entry.ToCycle.Name))
					}
					if entry.FromProject != nil && entry.ToProject != nil {
						changes = append(changes, fmt.Sprintf("Project: %s → %s", entry.FromProject.Name, entry.ToProject.Name))
					}
					if len(entry.AddedLabelIds) > 0 {
						changes = append(changes, fmt.Sprintf("Added %d label(s)", len(entry.AddedLabelIds)))
					}
					if len(entry.RemovedLabelIds) > 0 {
						changes = append(changes, fmt.Sprintf("Removed %d label(s)", len(entry.RemovedLabelIds)))
					}

					if len(changes) > 0 {
						fmt.Printf("\n  - %s", strings.Join(changes, "\n  - "))
					}
					fmt.Println()
				}
			}

			return
		}

		// Rich display
		fmt.Printf("%s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
			color.New(color.FgWhite, color.Bold).Sprint(issue.Title))

		if issue.Description != "" {
			fmt.Printf("\n%s\n", issue.Description)
		}

		fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Details:"))

		if issue.State != nil {
			stateStr := issue.State.Name
			if issue.State.Type == "completed" && issue.CompletedAt != nil {
				stateStr += fmt.Sprintf(" (%s)", issue.CompletedAt.Format("2006-01-02"))
			}
			fmt.Printf("State: %s\n",
				color.New(color.FgGreen).Sprint(stateStr))
		}

		if issue.Assignee != nil {
			fmt.Printf("Assignee: %s\n",
				color.New(color.FgCyan).Sprint(issue.Assignee.Name))
		} else {
			fmt.Printf("Assignee: %s\n",
				color.New(color.FgRed).Sprint("Unassigned"))
		}

		if issue.Team != nil {
			fmt.Printf("Team: %s\n",
				color.New(color.FgMagenta).Sprint(issue.Team.Name))
		}

		fmt.Printf("Priority: %s\n", priorityToString(issue.Priority))

		// Show project and cycle info
		if issue.Project != nil {
			fmt.Printf("Project: %s (%s)\n",
				color.New(color.FgBlue).Sprint(issue.Project.Name),
				color.New(color.FgWhite, color.Faint).Sprintf("%.0f%%", issue.Project.Progress*100))
		}

		if issue.Cycle != nil {
			fmt.Printf("Cycle: %s\n",
				color.New(color.FgMagenta).Sprint(issue.Cycle.Name))
		}

		fmt.Printf("Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05"))

		if issue.DueDate != nil && *issue.DueDate != "" {
			fmt.Printf("Due Date: %s\n",
				color.New(color.FgYellow).Sprint(*issue.DueDate))
		}

		if issue.SnoozedUntilAt != nil {
			fmt.Printf("Snoozed Until: %s\n",
				color.New(color.FgYellow).Sprint(issue.SnoozedUntilAt.Format("2006-01-02 15:04:05")))
		}

		// Show git branch if available
		if issue.BranchName != "" {
			fmt.Printf("Git Branch: %s\n",
				color.New(color.FgGreen).Sprint(issue.BranchName))
		}

		// Show URL
		if issue.URL != "" {
			fmt.Printf("URL: %s\n",
				color.New(color.FgBlue, color.Underline).Sprint(issue.URL))
		}

		// Show parent issue if this is a sub-issue
		if issue.Parent != nil {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Parent Issue:"))
			fmt.Printf("  %s %s\n",
				color.New(color.FgCyan).Sprint(issue.Parent.Identifier),
				issue.Parent.Title)
		}

		// Show sub-issues if any
		if issue.Children != nil && len(issue.Children.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Sub-issues:"))
			for _, child := range issue.Children.Nodes {
				stateIcon := "○"
				if child.State != nil {
					switch child.State.Type {
					case "completed", "done":
						stateIcon = color.New(color.FgGreen).Sprint("✓")
					case "started", "in_progress":
						stateIcon = color.New(color.FgBlue).Sprint("◐")
					case "canceled":
						stateIcon = color.New(color.FgRed).Sprint("✗")
					}
				}

				assignee := "Unassigned"
				if child.Assignee != nil {
					assignee = child.Assignee.Name
				}

				fmt.Printf("  %s %s %s (%s)\n",
					stateIcon,
					color.New(color.FgCyan).Sprint(child.Identifier),
					child.Title,
					color.New(color.FgWhite, color.Faint).Sprint(assignee))
			}
		}

		// Show attachments if any
		if issue.Attachments != nil && len(issue.Attachments.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Attachments:"))
			for _, attachment := range issue.Attachments.Nodes {
				fmt.Printf("  📎 %s - %s\n",
					attachment.Title,
					color.New(color.FgBlue, color.Underline).Sprint(attachment.URL))
			}
		}

		// Show recent comments if any
		if issue.Comments != nil && len(issue.Comments.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Recent Comments:"))
			for _, comment := range issue.Comments.Nodes {
				fmt.Printf("  💬 %s - %s\n",
					color.New(color.FgCyan).Sprint(comment.User.Name),
					color.New(color.FgWhite, color.Faint).Sprint(comment.CreatedAt.Format("2006-01-02 15:04")))
				// Show first line of comment
				lines := strings.Split(comment.Body, "\n")
				if len(lines) > 0 && lines[0] != "" {
					preview := lines[0]
					if len(preview) > 60 {
						preview = preview[:57] + "..."
					}
					fmt.Printf("     %s\n", preview)
				}
			}
			fmt.Printf("\n  %s Use 'linctl comment list %s' to see all comments\n",
				color.New(color.FgWhite, color.Faint).Sprint("→"),
				issue.Identifier)
		}
	},
}

func buildIssueFilter(cmd *cobra.Command) map[string]interface{} {
	filter := make(map[string]interface{})

	if assignee, _ := cmd.Flags().GetString("assignee"); assignee != "" {
		if assignee == "me" {
			// We'll need to get the current user's ID
			// For now, we'll use a special marker
			filter["assignee"] = map[string]interface{}{"isMe": map[string]interface{}{"eq": true}}
		} else {
			filter["assignee"] = map[string]interface{}{"email": map[string]interface{}{"eq": assignee}}
		}
	}

	state, _ := cmd.Flags().GetString("state")
	if state != "" {
		filter["state"] = map[string]interface{}{"name": map[string]interface{}{"eq": state}}
	} else {
		// Only filter out completed issues if no specific state is requested
		includeCompleted, _ := cmd.Flags().GetBool("include-completed")
		if !includeCompleted {
			// Filter out completed and canceled states
			filter["state"] = map[string]interface{}{
				"type": map[string]interface{}{
					"nin": []string{"completed", "canceled"},
				},
			}
		}
	}

	if team, _ := cmd.Flags().GetString("team"); team != "" {
		filter["team"] = map[string]interface{}{"key": map[string]interface{}{"eq": team}}
	}

	if priority, _ := cmd.Flags().GetInt("priority"); priority != -1 {
		filter["priority"] = map[string]interface{}{"eq": priority}
	}

	// Handle newer-than filter
	newerThan, _ := cmd.Flags().GetString("newer-than")
	createdAt, err := utils.ParseTimeExpression(newerThan)
	if err != nil {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		output.Error(fmt.Sprintf("Invalid newer-than value: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}
	if createdAt != "" {
		filter["createdAt"] = map[string]interface{}{"gte": createdAt}
	}

	return filter
}

func priorityToString(priority int) string {
	switch priority {
	case 0:
		return "None"
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Normal"
	case 4:
		return "Low"
	default:
		return "Unknown"
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

var issueAssignCmd = &cobra.Command{
	Use:   "assign [issue-id]",
	Short: "Assign issue to yourself",
	Long:  `Assign an issue to yourself.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Get current user
		viewer, err := client.GetViewer(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Update issue with assignee
		input := map[string]interface{}{
			"assigneeId": viewer.ID,
		}

		issue, err := client.UpdateIssue(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to assign issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issue)
		} else if plaintext {
			fmt.Printf("Assigned %s to %s\n", issue.Identifier, viewer.Name)
		} else {
			fmt.Printf("%s Assigned %s to %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
				color.New(color.FgCyan).Sprint(viewer.Name))
		}
	},
}

var issueCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new issue",
	Long:    `Create a new issue in Linear.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Get flags
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		teamKey, _ := cmd.Flags().GetString("team")
		priority, _ := cmd.Flags().GetInt("priority")
		assignToMe, _ := cmd.Flags().GetBool("assign-me")

		if title == "" {
			output.Error("Title is required (--title)", plaintext, jsonOut)
			os.Exit(1)
		}

		if teamKey == "" {
			output.Error("Team is required (--team)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Get team ID from key
		team, err := client.GetTeam(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Build input
		input := map[string]interface{}{
			"title":  title,
			"teamId": team.ID,
		}

		if description != "" {
			input["description"] = description
		}

		if priority >= 0 && priority <= 4 {
			input["priority"] = priority
		}

		if assignToMe {
			viewer, err := client.GetViewer(context.Background())
			if err != nil {
				output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["assigneeId"] = viewer.ID
		}

		// Create issue
		issue, err := client.CreateIssue(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issue)
		} else if plaintext {
			fmt.Printf("Created issue %s: %s\n", issue.Identifier, issue.Title)
		} else {
			fmt.Printf("%s Created issue %s: %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
				issue.Title)
			if issue.Assignee != nil {
				fmt.Printf("  Assigned to: %s\n", color.New(color.FgCyan).Sprint(issue.Assignee.Name))
			}
		}
	},
}

var issueUpdateCmd = &cobra.Command{
	Use:   "update [issue-id]",
	Short: "Update an issue",
	Long: `Update various fields of an issue.

Examples:
  linctl issue update LIN-123 --title "New title"
  linctl issue update LIN-123 --description "Updated description"
  linctl issue update LIN-123 --assignee john.doe@company.com
  linctl issue update LIN-123 --state "In Progress"
  linctl issue update LIN-123 --priority 1
  linctl issue update LIN-123 --due-date "2024-12-31"
  linctl issue update LIN-123 --title "New title" --assignee me --priority 2`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Build update input
		input := make(map[string]interface{})

		// Handle title update
		if cmd.Flags().Changed("title") {
			title, _ := cmd.Flags().GetString("title")
			input["title"] = title
		}

		// Handle description update
		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			input["description"] = description
		}

		// Handle assignee update
		if cmd.Flags().Changed("assignee") {
			assignee, _ := cmd.Flags().GetString("assignee")
			switch assignee {
			case "me":
				// Get current user
				viewer, err := client.GetViewer(context.Background())
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}
				input["assigneeId"] = viewer.ID
			case "unassigned", "":
				input["assigneeId"] = nil
			default:
				// Look up user by email
				users, err := client.GetUsers(context.Background(), 100, "", "")
				if err != nil {
					output.Error(fmt.Sprintf("Failed to get users: %v", err), plaintext, jsonOut)
					os.Exit(1)
				}

				var foundUser *api.User
				for _, user := range users.Nodes {
					if user.Email == assignee || user.Name == assignee {
						foundUser = &user
						break
					}
				}

				if foundUser == nil {
					output.Error(fmt.Sprintf("User not found: %s", assignee), plaintext, jsonOut)
					os.Exit(1)
				}

				input["assigneeId"] = foundUser.ID
			}
		}

		// Handle state update
		if cmd.Flags().Changed("state") {
			stateName, _ := cmd.Flags().GetString("state")

			// First, get the issue to know which team it belongs to
			issue, err := client.GetIssue(context.Background(), args[0])
			if err != nil {
				output.Error(fmt.Sprintf("Failed to get issue: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			// Get available states for the team
			states, err := client.GetTeamStates(context.Background(), issue.Team.Key)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to get team states: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			// Find the state by name (case-insensitive)
			var stateID string
			for _, state := range states {
				if strings.EqualFold(state.Name, stateName) {
					stateID = state.ID
					break
				}
			}

			if stateID == "" {
				// Show available states
				var stateNames []string
				for _, state := range states {
					stateNames = append(stateNames, state.Name)
				}
				output.Error(fmt.Sprintf("State '%s' not found. Available states: %s", stateName, strings.Join(stateNames, ", ")), plaintext, jsonOut)
				os.Exit(1)
			}

			input["stateId"] = stateID
		}

		// Handle priority update
		if cmd.Flags().Changed("priority") {
			priority, _ := cmd.Flags().GetInt("priority")
			input["priority"] = priority
		}

		// Handle due date update
		if cmd.Flags().Changed("due-date") {
			dueDate, _ := cmd.Flags().GetString("due-date")
			if dueDate == "" {
				input["dueDate"] = nil
			} else {
				input["dueDate"] = dueDate
			}
		}

		// Check if any updates were specified
		if len(input) == 0 {
			output.Error("No updates specified. Use flags to specify what to update.", plaintext, jsonOut)
			os.Exit(1)
		}

		// Update the issue
		issue, err := client.UpdateIssue(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(issue)
		} else if plaintext {
			fmt.Printf("Updated issue %s\n", issue.Identifier)
		} else {
			output.Success(fmt.Sprintf("Updated issue %s", issue.Identifier), plaintext, jsonOut)
		}
	},
}

var issueRelateCmd = &cobra.Command{
	Use:     "relate [issue-id]",
	Aliases: []string{"link"},
	Short:   "Create a relation between issues",
	Long: `Create a relation between two issues (blocks, blocked-by, related, duplicate).

Examples:
  linctl issue relate LIN-123 --blocks LIN-456
  linctl issue relate LIN-123 --blocked-by LIN-456
  linctl issue relate LIN-123 --related LIN-456
  linctl issue relate LIN-123 --duplicate LIN-456`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		sourceIssue := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Determine which relation type flag was used
		var relationType string
		var targetIssue string
		var swapIssues bool // For --blocked-by, we need to reverse the relationship

		blocksTarget, _ := cmd.Flags().GetString("blocks")
		blockedByTarget, _ := cmd.Flags().GetString("blocked-by")
		relatedTarget, _ := cmd.Flags().GetString("related")
		duplicateTarget, _ := cmd.Flags().GetString("duplicate")

		// Count how many flags were set
		flagsSet := 0
		if blocksTarget != "" {
			flagsSet++
			relationType = "blocks"
			targetIssue = blocksTarget
		}
		if blockedByTarget != "" {
			flagsSet++
			relationType = "blocks" // Linear uses "blocks", we swap the issues
			targetIssue = blockedByTarget
			swapIssues = true // Reverse: "A blocked-by B" means "B blocks A"
		}
		if relatedTarget != "" {
			flagsSet++
			relationType = "related"
			targetIssue = relatedTarget
		}
		if duplicateTarget != "" {
			flagsSet++
			relationType = "duplicate"
			targetIssue = duplicateTarget
		}

		// Validate that exactly one relation flag was specified
		if flagsSet == 0 {
			output.Error("Must specify one relation type: --blocks, --blocked-by, --related, or --duplicate", plaintext, jsonOut)
			os.Exit(1)
		}
		if flagsSet > 1 {
			output.Error("Only one relation type can be specified at a time", plaintext, jsonOut)
			os.Exit(1)
		}

		// Look up source issue to get its UUID
		sourceIssueData, err := client.GetIssue(context.Background(), sourceIssue)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find source issue %s: %v", sourceIssue, err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Look up target issue to get its UUID
		targetIssueData, err := client.GetIssue(context.Background(), targetIssue)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find target issue %s: %v", targetIssue, err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Determine which issue is the source and which is the target for the API call
		var apiSourceID, apiTargetID string
		if swapIssues {
			// For --blocked-by: "A blocked-by B" means "B blocks A"
			apiSourceID = targetIssueData.ID
			apiTargetID = sourceIssueData.ID
		} else {
			apiSourceID = sourceIssueData.ID
			apiTargetID = targetIssueData.ID
		}

		// Create the relation using UUIDs
		relation, err := client.CreateIssueRelation(context.Background(), apiSourceID, apiTargetID, relationType)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create relation: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(relation)
		} else if plaintext {
			// For plaintext, show from user's perspective
			displayRelation := relationType
			if swapIssues {
				displayRelation = "blocked-by"
			}
			fmt.Printf("Created %s relation between %s and %s\n", displayRelation, sourceIssue, targetIssue)
		} else {
			// Rich display - show from user's perspective
			var displaySourceID, displayTargetID, displaySourceTitle, displayTargetTitle string
			var relationLabel string

			if swapIssues {
				// User said "A blocked-by B", display it that way
				displaySourceID = sourceIssue
				displayTargetID = targetIssue
				displaySourceTitle = sourceIssueData.Title
				displayTargetTitle = targetIssueData.Title
				relationLabel = "is blocked by"
			} else {
				// Display as-is
				displaySourceID = relation.Issue.Identifier
				displayTargetID = relation.RelatedIssue.Identifier
				displaySourceTitle = relation.Issue.Title
				displayTargetTitle = relation.RelatedIssue.Title
				switch relationType {
				case "blocks":
					relationLabel = "blocks"
				case "related":
					relationLabel = "is related to"
				case "duplicate":
					relationLabel = "is duplicate of"
				}
			}

			fmt.Printf("%s Created relation: %s %s %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(displaySourceID),
				color.New(color.FgYellow).Sprint(relationLabel),
				color.New(color.FgCyan, color.Bold).Sprint(displayTargetID))
			fmt.Printf("  %s: %s\n",
				displaySourceID,
				displaySourceTitle)
			fmt.Printf("  %s: %s\n",
				displayTargetID,
				displayTargetTitle)
		}
	},
}

var issueUnrelateCmd = &cobra.Command{
	Use:     "unrelate [issue-id]",
	Aliases: []string{"unlink"},
	Short:   "Remove a relation between issues",
	Long: `Remove a relation between two issues.

Examples:
  linctl issue unrelate DEL-1681 --blocks DEL-1680
  linctl issue unrelate DEL-1680 --blocked-by DEL-1681
  linctl issue unrelate DEL-1681 --related DEL-1680
  linctl issue unrelate DEL-1680 --duplicate DEL-1681`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		sourceIssue := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Determine which relation type flag was used
		var relationType string
		var targetIssue string
		var swapForLookup bool // For --blocked-by, we need to check the reversed relation

		blocksTarget, _ := cmd.Flags().GetString("blocks")
		blockedByTarget, _ := cmd.Flags().GetString("blocked-by")
		relatedTarget, _ := cmd.Flags().GetString("related")
		duplicateTarget, _ := cmd.Flags().GetString("duplicate")

		// Count how many flags were set
		flagsSet := 0
		if blocksTarget != "" {
			flagsSet++
			relationType = "blocks"
			targetIssue = blocksTarget
		}
		if blockedByTarget != "" {
			flagsSet++
			relationType = "blocks" // Linear stores it as "blocks"
			targetIssue = blockedByTarget
			swapForLookup = true // Look for the reverse relation
		}
		if relatedTarget != "" {
			flagsSet++
			relationType = "related"
			targetIssue = relatedTarget
		}
		if duplicateTarget != "" {
			flagsSet++
			relationType = "duplicate"
			targetIssue = duplicateTarget
		}

		// Validate that exactly one relation flag was specified
		if flagsSet == 0 {
			output.Error("Must specify one relation type: --blocks, --blocked-by, --related, or --duplicate", plaintext, jsonOut)
			os.Exit(1)
		}
		if flagsSet > 1 {
			output.Error("Only one relation type can be specified at a time", plaintext, jsonOut)
			os.Exit(1)
		}

		// Look up both issues to get their data
		var lookupIssue string
		if swapForLookup {
			// For --blocked-by, the relation is stored on the target, not source
			lookupIssue = targetIssue
		} else {
			lookupIssue = sourceIssue
		}

		issueData, err := client.GetIssue(context.Background(), lookupIssue)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find issue %s: %v", lookupIssue, err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Look up target issue to get its UUID for matching
		targetIssueData, err := client.GetIssue(context.Background(), targetIssue)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to find target issue %s: %v", targetIssue, err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Find the relation in the issue's relations
		var relationID string
		var matchIssueID string
		if swapForLookup {
			matchIssueID = issueData.ID // Looking for source issue in target's relations
			// Actually, we need the source issue ID to match
			sourceIssueData, err := client.GetIssue(context.Background(), sourceIssue)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find source issue %s: %v", sourceIssue, err), plaintext, jsonOut)
				os.Exit(1)
			}
			matchIssueID = sourceIssueData.ID
		} else {
			matchIssueID = targetIssueData.ID
		}

		if issueData.Relations != nil {
			for _, relation := range issueData.Relations.Nodes {
				if relation.Type == relationType && relation.RelatedIssue != nil && relation.RelatedIssue.ID == matchIssueID {
					relationID = relation.ID
					break
				}
			}
		}

		if relationID == "" {
			output.Error(fmt.Sprintf("No %s relation found between %s and %s", relationType, sourceIssue, targetIssue), plaintext, jsonOut)
			os.Exit(1)
		}

		// Delete the relation
		err = client.DeleteIssueRelation(context.Background(), relationID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete relation: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(map[string]interface{}{
				"success": true,
				"source":  sourceIssue,
				"target":  targetIssue,
				"type":    relationType,
			})
		} else if plaintext {
			displayRelation := relationType
			if swapForLookup {
				displayRelation = "blocked-by"
			}
			fmt.Printf("Removed %s relation between %s and %s\n", displayRelation, sourceIssue, targetIssue)
		} else {
			// Rich display
			relationLabel := relationType
			if swapForLookup {
				relationLabel = "blocked by"
			} else {
				switch relationType {
				case "blocks":
					relationLabel = "blocks"
				case "related":
					relationLabel = "related to"
				case "duplicate":
					relationLabel = "duplicate of"
				}
			}

			fmt.Printf("%s Removed relation: %s %s %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(sourceIssue),
				color.New(color.FgYellow).Sprint(relationLabel),
				color.New(color.FgCyan, color.Bold).Sprint(targetIssue))
		}
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueSearchCmd)
	issueCmd.AddCommand(issueGetCmd)
	issueCmd.AddCommand(issueAssignCmd)
	issueCmd.AddCommand(issueCreateCmd)
	issueCmd.AddCommand(issueUpdateCmd)
	issueCmd.AddCommand(issueRelateCmd)
	issueCmd.AddCommand(issueUnrelateCmd)

	// Issue list flags
	issueListCmd.Flags().StringP("assignee", "a", "", "Filter by assignee (email or 'me')")
	issueListCmd.Flags().StringP("state", "s", "", "Filter by state name")
	issueListCmd.Flags().StringP("team", "t", "", "Filter by team key")
	issueListCmd.Flags().IntP("priority", "r", -1, "Filter by priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)")
	issueListCmd.Flags().IntP("limit", "l", 50, "Maximum number of issues to fetch")
	issueListCmd.Flags().BoolP("include-completed", "c", false, "Include completed and canceled issues")
	issueListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	issueListCmd.Flags().StringP("newer-than", "n", "", "Show issues created after this time (default: 6_months_ago, use 'all_time' for no filter)")

	// Issue search flags
	issueSearchCmd.Flags().StringP("assignee", "a", "", "Filter by assignee (email or 'me')")
	issueSearchCmd.Flags().StringP("state", "s", "", "Filter by state name")
	issueSearchCmd.Flags().StringP("team", "t", "", "Filter by team key")
	issueSearchCmd.Flags().IntP("priority", "r", -1, "Filter by priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)")
	issueSearchCmd.Flags().IntP("limit", "l", 50, "Maximum number of issues to fetch")
	issueSearchCmd.Flags().BoolP("include-completed", "c", false, "Include completed and canceled issues")
	issueSearchCmd.Flags().Bool("include-archived", false, "Include archived issues in results")
	issueSearchCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	issueSearchCmd.Flags().StringP("newer-than", "n", "", "Show issues created after this time (default: 6_months_ago, use 'all_time' for no filter)")

	// Issue create flags
	issueCreateCmd.Flags().StringP("title", "", "", "Issue title (required)")
	issueCreateCmd.Flags().StringP("description", "d", "", "Issue description")
	issueCreateCmd.Flags().StringP("team", "t", "", "Team key (required)")
	issueCreateCmd.Flags().Int("priority", 3, "Priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)")
	issueCreateCmd.Flags().BoolP("assign-me", "m", false, "Assign to yourself")
	_ = issueCreateCmd.MarkFlagRequired("title")
	_ = issueCreateCmd.MarkFlagRequired("team")

	// Issue update flags
	issueUpdateCmd.Flags().String("title", "", "New title for the issue")
	issueUpdateCmd.Flags().StringP("description", "d", "", "New description for the issue")
	issueUpdateCmd.Flags().StringP("assignee", "a", "", "Assignee (email, name, 'me', or 'unassigned')")
	issueUpdateCmd.Flags().StringP("state", "s", "", "State name (e.g., 'Todo', 'In Progress', 'Done')")
	issueUpdateCmd.Flags().Int("priority", -1, "Priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)")
	issueUpdateCmd.Flags().String("due-date", "", "Due date (YYYY-MM-DD format, or empty to remove)")

	// Issue relate flags
	issueRelateCmd.Flags().String("blocks", "", "Issue ID that this issue blocks")
	issueRelateCmd.Flags().String("blocked-by", "", "Issue ID that blocks this issue")
	issueRelateCmd.Flags().String("related", "", "Issue ID that is related to this issue")
	issueRelateCmd.Flags().String("duplicate", "", "Issue ID that this issue duplicates")

	// Issue unrelate flags
	issueUnrelateCmd.Flags().String("blocks", "", "Issue ID that this issue blocks")
	issueUnrelateCmd.Flags().String("blocked-by", "", "Issue ID that blocks this issue")
	issueUnrelateCmd.Flags().String("related", "", "Issue ID that is related to this issue")
	issueUnrelateCmd.Flags().String("duplicate", "", "Issue ID that this issue duplicates")
}
