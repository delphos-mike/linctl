package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/delphos-mike/linctl/pkg/api"
	"github.com/delphos-mike/linctl/pkg/auth"
	"github.com/delphos-mike/linctl/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage Linear labels",
	Long: `List, create, and manage Linear issue labels.

Examples:
  linctl label list
  linctl label list --team DEL
  linctl label create --name "mike:focus" --color "#eb5757"
  linctl label create --name "mike:focus" --group mike
  linctl label create --name "mike" --is-group`,
}

var labelListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List labels",
	Long: `List available labels in the workspace, optionally filtered by team.

Examples:
  linctl label list
  linctl label list --team DEL
  linctl label list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		teamKey, _ := cmd.Flags().GetString("team")

		labels, err := client.GetLabels(context.Background(), teamKey)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch labels: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(labels) == 0 {
			output.Info("No labels found", plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(labels)
			return
		}

		if plaintext {
			fmt.Println("# Labels")
			for _, label := range labels {
				prefix := "-"
				if label.IsGroup {
					prefix = "#"
				}
				fmt.Printf("%s %s", prefix, label.Name)
				if label.Description != nil && *label.Description != "" {
					fmt.Printf(" — %s", *label.Description)
				}
				if label.Parent != nil {
					fmt.Printf(" (group: %s)", label.Parent.Name)
				}
				if label.IsGroup && label.Children != nil && len(label.Children.Nodes) > 0 {
					var childNames []string
					for _, child := range label.Children.Nodes {
						childNames = append(childNames, child.Name)
					}
					fmt.Printf(" [%s]", strings.Join(childNames, ", "))
				}
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d labels\n", len(labels))
			return
		}

		// Rich table output
		headers := []string{"Name", "Type", "Color", "Description", "Group"}
		rows := make([][]string, len(labels))

		for i, label := range labels {
			desc := ""
			if label.Description != nil {
				desc = truncateString(*label.Description, 40)
			}

			group := ""
			if label.Parent != nil {
				group = label.Parent.Name
			}

			labelType := "label"
			if label.IsGroup {
				labelType = color.New(color.FgMagenta).Sprint("group")
			}

			rows[i] = []string{
				color.New(color.Bold).Sprint(label.Name),
				labelType,
				label.Color,
				desc,
				group,
			}
		}

		tableData := output.TableData{
			Headers: headers,
			Rows:    rows,
		}

		output.Table(tableData, false, false)

		fmt.Printf("\n%s %d labels\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(labels))
	},
}

var labelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new label",
	Long: `Create a new workspace label in Linear.

To create a label group (parent), use --is-group.
To create a label under a group, use --group with the group name.

Examples:
  linctl label create --name "mike:focus" --color "#eb5757"
  linctl label create --name "mike" --is-group --color "#5e6ad2" --description "Personal workflow labels"
  linctl label create --name "mike:focus" --group mike --color "#eb5757"
  linctl label create --name "mike:focus" --group mike`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		name, _ := cmd.Flags().GetString("name")
		labelColor, _ := cmd.Flags().GetString("color")
		description, _ := cmd.Flags().GetString("description")
		groupName, _ := cmd.Flags().GetString("group")
		isGroup, _ := cmd.Flags().GetBool("is-group")

		if name == "" {
			output.Error("Name is required (--name)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Resolve group name to ID if provided
		var parentID string
		if groupName != "" {
			labels, err := client.GetLabels(context.Background(), "")
			if err != nil {
				output.Error(fmt.Sprintf("Failed to fetch labels: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			for _, label := range labels {
				if strings.EqualFold(label.Name, groupName) {
					if !label.IsGroup {
						output.Error(fmt.Sprintf("Label %q exists but is not a group. Use --is-group to create a group first.", groupName), plaintext, jsonOut)
						os.Exit(1)
					}
					parentID = label.ID
					break
				}
			}

			if parentID == "" {
				// List available groups for a helpful error
				var groups []string
				for _, label := range labels {
					if label.IsGroup {
						groups = append(groups, label.Name)
					}
				}
				hint := "No label groups exist yet. Create one first with --is-group."
				if len(groups) > 0 {
					hint = fmt.Sprintf("Available groups: %s", strings.Join(groups, ", "))
				}
				output.Error(fmt.Sprintf("Group %q not found. %s", groupName, hint), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		label, err := client.CreateLabel(context.Background(), name, labelColor, description, parentID, isGroup)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create label: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(label)
			return
		}

		if plaintext {
			fmt.Printf("Created label: %s (id: %s)\n", label.Name, label.ID)
			if label.IsGroup {
				fmt.Println("Type: group")
			}
			if label.Parent != nil {
				fmt.Printf("Group: %s\n", label.Parent.Name)
			}
			return
		}

		kind := "label"
		if label.IsGroup {
			kind = "label group"
		}
		fmt.Printf("%s Created %s %s\n",
			color.New(color.FgGreen).Sprint("✓"),
			kind,
			color.New(color.FgCyan, color.Bold).Sprint(label.Name))
		if label.Parent != nil {
			fmt.Printf("  Group: %s\n", color.New(color.FgMagenta).Sprint(label.Parent.Name))
		}
		if label.Color != "" {
			fmt.Printf("  Color: %s\n", label.Color)
		}
	},
}

var labelDeleteCmd = &cobra.Command{
	Use:   "delete [name-or-id]",
	Short: "Delete a label",
	Long: `Delete a label by name or ID.

Examples:
  linctl label delete "mike:focus"
  linctl label delete 44c69f91-a7d8-46f6-a679-631c66095299`,
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

		nameOrID := args[0]

		// Try to resolve as a name first
		labelID := nameOrID
		labels, err := client.GetLabels(context.Background(), "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch labels: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		var labelName string
		for _, label := range labels {
			if strings.EqualFold(label.Name, nameOrID) || label.ID == nameOrID {
				labelID = label.ID
				labelName = label.Name
				break
			}
		}

		if labelName == "" {
			labelName = nameOrID
		}

		err = client.DeleteLabel(context.Background(), labelID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete label: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{"success": true, "name": labelName, "id": labelID})
			return
		}

		if plaintext {
			fmt.Printf("Deleted label: %s\n", labelName)
			return
		}

		fmt.Printf("%s Deleted label %s\n",
			color.New(color.FgGreen).Sprint("✓"),
			color.New(color.FgCyan, color.Bold).Sprint(labelName))
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)
	labelCmd.AddCommand(labelListCmd)
	labelCmd.AddCommand(labelCreateCmd)
	labelCmd.AddCommand(labelDeleteCmd)

	// Label list flags
	labelListCmd.Flags().StringP("team", "t", "", "Filter by team key")

	// Label create flags
	labelCreateCmd.Flags().String("name", "", "Label name (required)")
	labelCreateCmd.Flags().String("color", "", "Label color as hex (e.g. #eb5757)")
	labelCreateCmd.Flags().StringP("description", "d", "", "Label description")
	labelCreateCmd.Flags().String("group", "", "Parent group name (label will be a child of this group)")
	labelCreateCmd.Flags().Bool("is-group", false, "Create as a label group (can contain child labels)")
	_ = labelCreateCmd.MarkFlagRequired("name")
}
