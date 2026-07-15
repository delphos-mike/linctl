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

// constructProjectURL constructs an ID-based project URL
func constructProjectURL(projectID string, originalURL string) string {
	// Extract workspace from the original URL
	// Format: https://linear.app/{workspace}/project/{slug}
	if originalURL == "" {
		return ""
	}

	parts := strings.Split(originalURL, "/")
	if len(parts) >= 5 {
		workspace := parts[3]
		return fmt.Sprintf("https://linear.app/%s/project/%s", workspace, projectID)
	}

	// Fallback to original URL if we can't parse it
	return originalURL
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Linear projects",
	Long: `Manage Linear projects including listing, viewing, and creating projects.

Examples:
  linctl project list                      # List active projects
  linctl project list --include-completed  # List all projects including completed
  linctl project list --newer-than 1_month_ago  # List projects from last month
  linctl project get PROJECT-ID            # Get project details
  linctl project create                    # Create a new project`,
	SilenceUsage: true,
	RunE:         requireSubcommand,
}

var projectListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List projects",
	Long:    `List all projects in your Linear workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get filters
		teamKey, _ := cmd.Flags().GetString("team")
		state, _ := cmd.Flags().GetString("state")
		limit, _ := cmd.Flags().GetInt("limit")
		includeCompleted, _ := cmd.Flags().GetBool("include-completed")

		// Build filter
		filter := make(map[string]interface{})
		if teamKey != "" {
			// Get team ID from key
			team, err := client.GetTeam(context.Background(), teamKey)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
				os.Exit(1)
			}
			filter["team"] = map[string]interface{}{"id": team.ID}
		}
		if state != "" {
			filter["state"] = map[string]interface{}{"eq": state}
		} else if !includeCompleted {
			// Only filter out completed projects if no specific state is requested
			filter["state"] = map[string]interface{}{
				"nin": []string{"completed", "canceled"},
			}
		}

		// Handle newer-than filter
		newerThan, _ := cmd.Flags().GetString("newer-than")
		createdAt, err := utils.ParseTimeExpression(newerThan)
		if err != nil {
			output.Error(fmt.Sprintf("Invalid newer-than value: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if createdAt != "" {
			filter["createdAt"] = map[string]interface{}{"gte": createdAt}
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

		// Get projects
		projects, err := client.GetProjects(context.Background(), filter, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list projects: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(projects.Nodes)
			return
		} else if plaintext {
			fmt.Println("# Projects")
			for _, project := range projects.Nodes {
				fmt.Printf("## %s\n", project.Name)
				fmt.Printf("- **ID**: %s\n", project.ID)
				fmt.Printf("- **State**: %s\n", project.State)
				fmt.Printf("- **Progress**: %.0f%%\n", project.Progress*100)
				if project.Lead != nil {
					fmt.Printf("- **Lead**: %s\n", project.Lead.Name)
				} else {
					fmt.Printf("- **Lead**: Unassigned\n")
				}
				if project.Teams != nil && len(project.Teams.Nodes) > 0 {
					teams := ""
					for i, team := range project.Teams.Nodes {
						if i > 0 {
							teams += ", "
						}
						teams += team.Key
					}
					fmt.Printf("- **Teams**: %s\n", teams)
				}
				if project.StartDate != nil {
					fmt.Printf("- **Start Date**: %s\n", *project.StartDate)
				}
				if project.TargetDate != nil {
					fmt.Printf("- **Target Date**: %s\n", *project.TargetDate)
				}
				fmt.Printf("- **Created**: %s\n", project.CreatedAt.Format("2006-01-02"))
				fmt.Printf("- **Updated**: %s\n", project.UpdatedAt.Format("2006-01-02"))
				fmt.Printf("- **URL**: %s\n", constructProjectURL(project.ID, project.URL))
				if project.Description != "" {
					fmt.Printf("- **Description**: %s\n", project.Description)
				}
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d projects\n", len(projects.Nodes))
			return
		} else {
			// Table output
			headers := []string{"Name", "State", "Lead", "Teams", "Created", "Updated", "URL"}
			rows := [][]string{}

			for _, project := range projects.Nodes {
				lead := color.New(color.FgYellow).Sprint("Unassigned")
				if project.Lead != nil {
					lead = project.Lead.Name
				}

				teams := ""
				if project.Teams != nil && len(project.Teams.Nodes) > 0 {
					for i, team := range project.Teams.Nodes {
						if i > 0 {
							teams += ", "
						}
						teams += team.Key
					}
				}

				stateColor := color.New(color.FgGreen)
				switch project.State {
				case "planned":
					stateColor = color.New(color.FgCyan)
				case "started":
					stateColor = color.New(color.FgBlue)
				case "paused":
					stateColor = color.New(color.FgYellow)
				case "completed":
					stateColor = color.New(color.FgGreen)
				case "canceled":
					stateColor = color.New(color.FgRed)
				}

				rows = append(rows, []string{
					truncateString(project.Name, 25),
					stateColor.Sprint(project.State),
					lead,
					teams,
					project.CreatedAt.Format("2006-01-02"),
					project.UpdatedAt.Format("2006-01-02"),
					constructProjectURL(project.ID, project.URL),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			if !plaintext && !jsonOut {
				fmt.Printf("\n%s %d projects\n",
					color.New(color.FgGreen).Sprint("✓"),
					len(projects.Nodes))
			}
		}
	},
}

var projectGetCmd = &cobra.Command{
	Use:     "get PROJECT-ID",
	Aliases: []string{"show"},
	Short:   "Get project details",
	Long:    `Get detailed information about a specific project.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		projectID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get project details
		project, err := client.GetProject(context.Background(), projectID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(project)
		} else if plaintext {
			fmt.Printf("# %s\n\n", project.Name)

			if project.Description != "" {
				fmt.Printf("## Description\n%s\n\n", project.Description)
			}

			if project.Content != "" {
				fmt.Printf("## Content\n%s\n\n", project.Content)
			}

			fmt.Printf("## Core Details\n")
			fmt.Printf("- **ID**: %s\n", project.ID)
			fmt.Printf("- **Slug ID**: %s\n", project.SlugId)
			fmt.Printf("- **State**: %s\n", project.State)
			fmt.Printf("- **Progress**: %.0f%%\n", project.Progress*100)
			fmt.Printf("- **Health**: %s\n", project.Health)
			fmt.Printf("- **Scope**: %d\n", project.Scope)
			if project.Icon != nil && *project.Icon != "" {
				fmt.Printf("- **Icon**: %s\n", *project.Icon)
			}
			fmt.Printf("- **Color**: %s\n", project.Color)

			fmt.Printf("\n## Timeline\n")
			if project.StartDate != nil {
				fmt.Printf("- **Start Date**: %s\n", *project.StartDate)
			}
			if project.TargetDate != nil {
				fmt.Printf("- **Target Date**: %s\n", *project.TargetDate)
			}
			fmt.Printf("- **Created**: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", project.UpdatedAt.Format("2006-01-02 15:04:05"))
			if project.CompletedAt != nil {
				fmt.Printf("- **Completed**: %s\n", project.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			if project.CanceledAt != nil {
				fmt.Printf("- **Canceled**: %s\n", project.CanceledAt.Format("2006-01-02 15:04:05"))
			}
			if project.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", project.ArchivedAt.Format("2006-01-02 15:04:05"))
			}

			fmt.Printf("\n## People\n")
			if project.Lead != nil {
				fmt.Printf("- **Lead**: %s (%s)\n", project.Lead.Name, project.Lead.Email)
				if project.Lead.DisplayName != "" && project.Lead.DisplayName != project.Lead.Name {
					fmt.Printf("  - Display Name: %s\n", project.Lead.DisplayName)
				}
			} else {
				fmt.Printf("- **Lead**: Unassigned\n")
			}
			if project.Creator != nil {
				fmt.Printf("- **Creator**: %s (%s)\n", project.Creator.Name, project.Creator.Email)
			}

			fmt.Printf("\n## Slack Integration\n")
			fmt.Printf("- **Slack New Issue**: %v\n", project.SlackNewIssue)
			fmt.Printf("- **Slack Issue Comments**: %v\n", project.SlackIssueComments)
			fmt.Printf("- **Slack Issue Statuses**: %v\n", project.SlackIssueStatuses)

			if project.ConvertedFromIssue != nil {
				fmt.Printf("\n## Origin\n")
				fmt.Printf("- **Converted from Issue**: %s - %s\n", project.ConvertedFromIssue.Identifier, project.ConvertedFromIssue.Title)
			}

			if project.LastAppliedTemplate != nil {
				fmt.Printf("\n## Template\n")
				fmt.Printf("- **Last Applied**: %s\n", project.LastAppliedTemplate.Name)
				if project.LastAppliedTemplate.Description != "" {
					fmt.Printf("  - Description: %s\n", project.LastAppliedTemplate.Description)
				}
			}

			// Teams
			if project.Teams != nil && len(project.Teams.Nodes) > 0 {
				fmt.Printf("\n## Teams\n")
				for _, team := range project.Teams.Nodes {
					fmt.Printf("- **%s** (%s)\n", team.Name, team.Key)
					if team.Description != "" {
						fmt.Printf("  - Description: %s\n", team.Description)
					}
					fmt.Printf("  - Cycles Enabled: %v\n", team.CyclesEnabled)
				}
			}

			fmt.Printf("\n## URL\n")
			fmt.Printf("- %s\n", constructProjectURL(project.ID, project.URL))

			// Show members if available
			if project.Members != nil && len(project.Members.Nodes) > 0 {
				fmt.Printf("\n## Members\n")
				for _, member := range project.Members.Nodes {
					fmt.Printf("- %s (%s)", member.Name, member.Email)
					if member.DisplayName != "" && member.DisplayName != member.Name {
						fmt.Printf(" - %s", member.DisplayName)
					}
					if member.Admin {
						fmt.Printf(" [Admin]")
					}
					if !member.Active {
						fmt.Printf(" [Inactive]")
					}
					fmt.Println()
				}
			}

			// Project Updates
			if project.ProjectUpdates != nil && len(project.ProjectUpdates.Nodes) > 0 {
				fmt.Printf("\n## Recent Project Updates\n")
				for _, update := range project.ProjectUpdates.Nodes {
					fmt.Printf("\n### %s by %s\n", update.CreatedAt.Format("2006-01-02 15:04"), update.User.Name)
					if update.EditedAt != nil {
						fmt.Printf("*(edited %s)*\n", update.EditedAt.Format("2006-01-02 15:04"))
					}
					fmt.Printf("- **Health**: %s\n", update.Health)
					fmt.Printf("\n%s\n", update.Body)
				}
			}

			// Documents
			if project.Documents != nil && len(project.Documents.Nodes) > 0 {
				fmt.Printf("\n## Documents\n")
				for _, doc := range project.Documents.Nodes {
					fmt.Printf("\n### %s\n", doc.Title)
					if doc.Icon != nil && *doc.Icon != "" {
						fmt.Printf("- **Icon**: %s\n", *doc.Icon)
					}
					fmt.Printf("- **Color**: %s\n", doc.Color)
					fmt.Printf("- **Created**: %s by %s\n", doc.CreatedAt.Format("2006-01-02"), doc.Creator.Name)
					if doc.UpdatedBy != nil {
						fmt.Printf("- **Updated**: %s by %s\n", doc.UpdatedAt.Format("2006-01-02"), doc.UpdatedBy.Name)
					}
					fmt.Printf("\n%s\n", doc.Content)
				}
			}

			// Milestones
			if project.ProjectMilestones != nil && len(project.ProjectMilestones.Nodes) > 0 {
				activeMilestones := filterActiveMilestones(project.ProjectMilestones.Nodes)
				if len(activeMilestones) > 0 {
					fmt.Printf("\n## Milestones\n")
					for _, m := range activeMilestones {
						fmt.Printf("\n### %s\n", m.Name)
						fmt.Printf("- **Status**: %s\n", m.Status)
						targetDate := "-"
						if m.TargetDate != nil {
							targetDate = *m.TargetDate
						}
						fmt.Printf("- **Target Date**: %s\n", targetDate)
						fmt.Printf("- **Progress**: %.0f%%\n", m.Progress*100)
						if m.Description != "" {
							fmt.Printf("- **Description**: %s\n", m.Description)
						}
					}
				}
			}

			// Show recent issues
			if project.Issues != nil && len(project.Issues.Nodes) > 0 {
				fmt.Printf("\n## Issues (%d total)\n", len(project.Issues.Nodes))
				for _, issue := range project.Issues.Nodes {
					stateStr := ""
					if issue.State != nil {
						switch issue.State.Type {
						case "completed":
							stateStr = "[x]"
						case "started":
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
					if issue.Assignee != nil {
						assignee = issue.Assignee.Name
					}

					fmt.Printf("\n### %s %s (#%d)\n", stateStr, issue.Identifier, issue.Number)
					fmt.Printf("**%s**\n", issue.Title)
					fmt.Printf("- Assignee: %s\n", assignee)
					fmt.Printf("- Priority: %s\n", priorityToString(issue.Priority))
					if issue.Estimate != nil {
						fmt.Printf("- Estimate: %.1f\n", *issue.Estimate)
					}
					if issue.State != nil {
						fmt.Printf("- State: %s\n", issue.State.Name)
					}
					if issue.Labels != nil && len(issue.Labels.Nodes) > 0 {
						labels := []string{}
						for _, label := range issue.Labels.Nodes {
							labels = append(labels, label.Name)
						}
						fmt.Printf("- Labels: %s\n", strings.Join(labels, ", "))
					}
					fmt.Printf("- Updated: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04"))
					if issue.Description != "" {
						// Show first 3 lines of description
						lines := strings.Split(issue.Description, "\n")
						preview := ""
						for i, line := range lines {
							if i >= 3 {
								preview += "\n  ..."
								break
							}
							if i > 0 {
								preview += "\n  "
							}
							preview += line
						}
						fmt.Printf("- Description: %s\n", preview)
					}
				}
			}
		} else {
			// Formatted output
			fmt.Println()
			fmt.Printf("%s %s\n", color.New(color.FgCyan, color.Bold).Sprint("📁 Project:"), project.Name)
			fmt.Println(strings.Repeat("─", 50))

			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), project.ID)

			if project.Description != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Description:"), project.Description)
			}

			stateColor := color.New(color.FgGreen)
			switch project.State {
			case "planned":
				stateColor = color.New(color.FgCyan)
			case "started":
				stateColor = color.New(color.FgBlue)
			case "paused":
				stateColor = color.New(color.FgYellow)
			case "completed":
				stateColor = color.New(color.FgGreen)
			case "canceled":
				stateColor = color.New(color.FgRed)
			}
			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("State:"), stateColor.Sprint(project.State))

			progressColor := color.New(color.FgRed)
			if project.Progress >= 0.75 {
				progressColor = color.New(color.FgGreen)
			} else if project.Progress >= 0.5 {
				progressColor = color.New(color.FgYellow)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Progress:"), progressColor.Sprintf("%.0f%%", project.Progress*100))

			if project.StartDate != nil || project.TargetDate != nil {
				fmt.Println()
				if project.StartDate != nil {
					fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Start Date:"), *project.StartDate)
				}
				if project.TargetDate != nil {
					fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Target Date:"), *project.TargetDate)
				}
			}

			if project.Lead != nil {
				fmt.Printf("\n%s %s (%s)\n",
					color.New(color.Bold).Sprint("Lead:"),
					project.Lead.Name,
					color.New(color.FgCyan).Sprint(project.Lead.Email))
			}

			if project.Teams != nil && len(project.Teams.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Teams:"))
				for _, team := range project.Teams.Nodes {
					fmt.Printf("  • %s - %s\n",
						color.New(color.FgCyan).Sprint(team.Key),
						team.Name)
				}
			}

			// Show members if available
			if project.Members != nil && len(project.Members.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Members:"))
				for _, member := range project.Members.Nodes {
					fmt.Printf("  • %s (%s)\n",
						member.Name,
						color.New(color.FgCyan).Sprint(member.Email))
				}
			}

			if project.ProjectMilestones != nil && len(project.ProjectMilestones.Nodes) > 0 {
				activeMilestones := filterActiveMilestones(project.ProjectMilestones.Nodes)
				if len(activeMilestones) > 0 {
					fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Milestones:"))
					limit := len(activeMilestones)
					if limit > 10 {
						limit = 10
					}
					for _, m := range activeMilestones[:limit] {
						sc := milestoneStatusColor(m.Status)
						targetDate := "-"
						if m.TargetDate != nil {
							targetDate = *m.TargetDate
						}
						fmt.Printf("  • %s (%s) — target: %s — %.0f%%\n",
							m.Name,
							sc.Sprint(m.Status),
							targetDate,
							m.Progress*100)
					}
					if len(activeMilestones) > 10 {
						fmt.Printf("  ... and %d more\n", len(activeMilestones)-10)
					}
				}
			}

			// Show sample issues if available
			if project.Issues != nil && len(project.Issues.Nodes) > 0 {
				fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Recent Issues:"))
				for i, issue := range project.Issues.Nodes {
					if i >= 5 {
						break // Show only first 5
					}
					stateIcon := "○"
					if issue.State != nil {
						switch issue.State.Type {
						case "completed":
							stateIcon = color.New(color.FgGreen).Sprint("✓")
						case "started":
							stateIcon = color.New(color.FgBlue).Sprint("◐")
						case "canceled":
							stateIcon = color.New(color.FgRed).Sprint("✗")
						}
					}
					assignee := "Unassigned"
					if issue.Assignee != nil {
						assignee = issue.Assignee.Name
					}
					fmt.Printf("  %s %s %s (%s)\n",
						stateIcon,
						color.New(color.FgCyan).Sprint(issue.Identifier),
						issue.Title,
						color.New(color.FgWhite, color.Faint).Sprint(assignee))
				}
			}

			// Show timestamps
			fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Timeline:"))
			fmt.Printf("  Created: %s\n", project.CreatedAt.Format("2006-01-02"))
			fmt.Printf("  Updated: %s\n", project.UpdatedAt.Format("2006-01-02"))
			if project.CompletedAt != nil {
				fmt.Printf("  Completed: %s\n", project.CompletedAt.Format("2006-01-02"))
			}
			if project.CanceledAt != nil {
				fmt.Printf("  Canceled: %s\n", project.CanceledAt.Format("2006-01-02"))
			}

			// Show URL
			if project.URL != "" {
				fmt.Printf("\n%s %s\n",
					color.New(color.Bold).Sprint("URL:"),
					color.New(color.FgBlue, color.Underline).Sprint(constructProjectURL(project.ID, project.URL)))
			}

			fmt.Println()
		}
	},
}

// milestoneStatusColor returns a color for a milestone status string
func milestoneStatusColor(status string) *color.Color {
	switch status {
	case "unstarted":
		return color.New(color.FgCyan)
	case "next":
		return color.New(color.FgBlue)
	case "overdue":
		return color.New(color.FgRed)
	case "done":
		return color.New(color.FgGreen)
	default:
		return color.New(color.FgWhite)
	}
}

// filterActiveMilestones removes archived milestones from a slice
func filterActiveMilestones(nodes []api.ProjectMilestone) []api.ProjectMilestone {
	var active []api.ProjectMilestone
	for _, m := range nodes {
		if m.ArchivedAt == nil {
			active = append(active, m)
		}
	}
	return active
}

var projectMilestonesCmd = &cobra.Command{
	Use:     "milestones PROJECT-ID",
	Aliases: []string{"ms"},
	Short:   "List project milestones",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		projectID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Fetch milestones
		milestones, err := client.GetProjectMilestones(context.Background(), projectID, 250, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get milestones: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if milestones.PageInfo.HasNextPage {
			fmt.Fprintln(os.Stderr, "⚠ Showing first 250 milestones. More exist.")
		}

		// Filter out archived milestones
		filtered := filterActiveMilestones(milestones.Nodes)

		if len(filtered) == 0 {
			fmt.Println("No milestones for this project")
			return
		}

		// Handle output
		if jsonOut {
			output.JSON(filtered)
			return
		} else if plaintext {
			fmt.Println("## Milestones")
			for _, m := range filtered {
				fmt.Printf("\n### %s\n", m.Name)
				fmt.Printf("- **Status**: %s\n", m.Status)
				targetDate := "-"
				if m.TargetDate != nil {
					targetDate = *m.TargetDate
				}
				fmt.Printf("- **Target Date**: %s\n", targetDate)
				fmt.Printf("- **Progress**: %.0f%%\n", m.Progress*100)
				if m.Description != "" {
					fmt.Printf("- **Description**: %s\n", m.Description)
				}
			}
			return
		} else {
			// Table output
			headers := []string{"Name", "Status", "Target Date", "Progress"}
			rows := [][]string{}

			for _, m := range filtered {
				sc := milestoneStatusColor(m.Status)
				targetDate := "-"
				if m.TargetDate != nil {
					targetDate = *m.TargetDate
				}
				rows = append(rows, []string{
					truncateString(m.Name, 30),
					sc.Sprint(m.Status),
					targetDate,
					fmt.Sprintf("%.0f%%", m.Progress*100),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d milestones\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(filtered))
		}
	},
}

// resolveTeamIDs resolves a comma-separated list of team keys to their IDs.
func resolveTeamIDs(ctx context.Context, client *api.Client, csv string) ([]string, error) {
	var ids []string
	for _, key := range strings.Split(csv, ",") {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		team, err := client.GetTeam(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve team %q: %w", key, err)
		}
		ids = append(ids, team.ID)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no valid team keys provided")
	}
	return ids, nil
}

// resolveLeadID resolves a user email to its ID.
func resolveLeadID(ctx context.Context, client *api.Client, email string) (string, error) {
	user, err := client.GetUser(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to resolve lead %q: %w", email, err)
	}
	return user.ID, nil
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project",
	Long: `Create a new project. At least one team is required. Body content is read
from --content, --body-file, or piped stdin.

Examples:
  linctl project create --name "New Project" --team ENG
  linctl project create --name "Q3 Effort" --team ENG,DES --lead alice@example.com --target-date 2026-09-30`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			output.Error("--name is required", plaintext, jsonOut)
			os.Exit(1)
		}
		teamCSV, _ := cmd.Flags().GetString("team")
		if teamCSV == "" {
			output.Error("--team is required (at least one team key)", plaintext, jsonOut)
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

		teamIDs, err := resolveTeamIDs(ctx, client, teamCSV)
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{
			"name":    name,
			"teamIds": teamIDs,
		}
		if hasContent {
			input["content"] = content
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			input["description"] = v
		}
		if v, _ := cmd.Flags().GetString("lead"); v != "" {
			leadID, err := resolveLeadID(ctx, client, v)
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["leadId"] = leadID
		}
		if v, _ := cmd.Flags().GetString("start-date"); v != "" {
			input["startDate"] = v
		}
		if v, _ := cmd.Flags().GetString("target-date"); v != "" {
			input["targetDate"] = v
		}
		if v, _ := cmd.Flags().GetString("status-id"); v != "" {
			input["statusId"] = v
		}

		project, err := client.CreateProject(ctx, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(project)
			return
		}
		output.Success(fmt.Sprintf("Created project %q (%s)", project.Name, project.ID), plaintext, jsonOut)
		if !plaintext {
			fmt.Printf("URL: %s\n", constructProjectURL(project.ID, project.URL))
		}
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update PROJECT-ID",
	Short: "Update a project",
	Long: `Update a project's name, description, body, lead, dates, teams, or status.
Body content is read from --content, --body-file, or piped stdin. Only provided
fields change.`,
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
		if cmd.Flags().Changed("team") {
			v, _ := cmd.Flags().GetString("team")
			teamIDs, err := resolveTeamIDs(ctx, client, v)
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["teamIds"] = teamIDs
		}
		if cmd.Flags().Changed("lead") {
			v, _ := cmd.Flags().GetString("lead")
			leadID, err := resolveLeadID(ctx, client, v)
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["leadId"] = leadID
		}
		if cmd.Flags().Changed("start-date") {
			v, _ := cmd.Flags().GetString("start-date")
			input["startDate"] = v
		}
		if cmd.Flags().Changed("target-date") {
			v, _ := cmd.Flags().GetString("target-date")
			input["targetDate"] = v
		}
		if cmd.Flags().Changed("status-id") {
			v, _ := cmd.Flags().GetString("status-id")
			input["statusId"] = v
		}

		if len(input) == 0 {
			output.Error("Nothing to update. Provide at least one field flag.", plaintext, jsonOut)
			os.Exit(1)
		}

		project, err := client.UpdateProject(ctx, args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update project: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(project)
			return
		}
		output.Success(fmt.Sprintf("Updated project %q (%s)", project.Name, project.ID), plaintext, jsonOut)
	},
}

var projectStatusUpdateCmd = &cobra.Command{
	Use:          "status-update",
	Aliases:      []string{"su"},
	Short:        "Manage project status updates",
	SilenceUsage: true,
	RunE:         requireSubcommand,
}

var projectStatusUpdateCreateCmd = &cobra.Command{
	Use:   "create PROJECT-ID",
	Short: "Post a project status update",
	Long: `Post a health + body status update to a project. Body content is read from
--body, --body-file, or piped stdin.

Examples:
  linctl project status-update create PROJECT-ID --health onTrack --body "Shipping on schedule."
  cat update.md | linctl project status-update create PROJECT-ID --health atRisk`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		health, _ := cmd.Flags().GetString("health")
		if !isValidHealth(health) {
			output.Error("--health is required and must be one of: onTrack, atRisk, offTrack", plaintext, jsonOut)
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

		input := map[string]interface{}{
			"projectId": args[0],
			"health":    health,
		}
		if hasBody {
			input["body"] = body
		}

		update, err := client.CreateProjectUpdate(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create project update: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(update)
			return
		}
		output.Success(fmt.Sprintf("Posted %s update to project %s", update.Health, args[0]), plaintext, jsonOut)
	},
}

// isValidHealth reports whether h is a valid ProjectUpdateHealthType value.
func isValidHealth(h string) bool {
	switch h {
	case "onTrack", "atRisk", "offTrack":
		return true
	default:
		return false
	}
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectMilestonesCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectStatusUpdateCmd)
	projectStatusUpdateCmd.AddCommand(projectStatusUpdateCreateCmd)

	// List command flags
	projectListCmd.Flags().StringP("team", "t", "", "Filter by team key")
	projectListCmd.Flags().StringP("state", "s", "", "Filter by state (planned, started, paused, completed, canceled)")
	projectListCmd.Flags().IntP("limit", "l", 50, "Maximum number of projects to return")
	projectListCmd.Flags().BoolP("include-completed", "c", false, "Include completed and canceled projects")
	projectListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	projectListCmd.Flags().StringP("newer-than", "n", "", "Show projects created after this time (default: 6_months_ago, use 'all_time' for no filter)")

	// Create command flags
	projectCreateCmd.Flags().String("name", "", "Project name (required)")
	projectCreateCmd.Flags().StringP("team", "t", "", "Team key(s), comma-separated (at least one required)")
	projectCreateCmd.Flags().String("description", "", "Short project description (summary)")
	projectCreateCmd.Flags().String("content", "", "Project body content (markdown)")
	projectCreateCmd.Flags().String("body-file", "", "Path to a file containing the project body content")
	projectCreateCmd.Flags().String("lead", "", "Project lead email")
	projectCreateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	projectCreateCmd.Flags().String("target-date", "", "Target date (YYYY-MM-DD)")
	projectCreateCmd.Flags().String("status-id", "", "Project status ID (ProjectStatus UUID)")

	// Update command flags
	projectUpdateCmd.Flags().String("name", "", "New project name")
	projectUpdateCmd.Flags().StringP("team", "t", "", "Replacement team key(s), comma-separated")
	projectUpdateCmd.Flags().String("description", "", "New short description (summary)")
	projectUpdateCmd.Flags().String("content", "", "New project body content (markdown)")
	projectUpdateCmd.Flags().String("body-file", "", "Path to a file containing the new project body content")
	projectUpdateCmd.Flags().String("lead", "", "New project lead email")
	projectUpdateCmd.Flags().String("start-date", "", "New start date (YYYY-MM-DD)")
	projectUpdateCmd.Flags().String("target-date", "", "New target date (YYYY-MM-DD)")
	projectUpdateCmd.Flags().String("status-id", "", "New project status ID (ProjectStatus UUID)")

	// Status update flags
	projectStatusUpdateCreateCmd.Flags().String("health", "", "Health: onTrack, atRisk, or offTrack (required)")
	projectStatusUpdateCreateCmd.Flags().String("body", "", "Update body content (markdown)")
	projectStatusUpdateCreateCmd.Flags().String("body-file", "", "Path to a file containing the update body")
}
