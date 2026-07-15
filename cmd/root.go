package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	plaintext bool
	jsonOut   bool
)

// version is set at build time via -ldflags (see Makefile, which injects
// github.com/delphos-mike/linctl/cmd.version). The default is for local dev
// builds that do not go through the Makefile.
var version = "dev"

// resolveVersion determines the version string to report. It reads the
// module version embedded by the Go toolchain and applies the fallback
// chain in pickVersion.
func resolveVersion() string {
	buildInfo := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		buildInfo = info.Main.Version
	}
	return pickVersion(version, buildInfo)
}

// pickVersion implements the version fallback chain. It is kept separate from
// the runtime build-info lookup in resolveVersion so it can be unit-tested.
//
//  1. injected: the -ldflags value (release builds via Makefile/CI) wins when
//     it is a real value (not empty and not the "dev" default).
//  2. buildInfo: the module version recorded by the toolchain. A
//     `go install module@vX.Y.Z` build carries the tag here; a plain source
//     build carries "(devel)", which we ignore.
//  3. otherwise fall back to "dev".
func pickVersion(injected, buildInfo string) string {
	if injected != "" && injected != "dev" {
		return injected
	}
	if buildInfo != "" && buildInfo != "(devel)" {
		return buildInfo
	}
	return "dev"
}

// requireSubcommand is the RunE for command groups that exist only to hold
// subcommands. Invoked bare (no positional args) it prints the group's help
// and exits 0; invoked with an unrecognized first argument it returns an
// "unknown command" error so the CLI exits non-zero instead of silently
// printing help and reporting success. Groups using this helper set
// SilenceUsage so the error prints exactly once to stderr.
func requireSubcommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	err := fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
	// Mirror cobra's own default so near-miss typos yield suggestions
	// (SuggestionsFor otherwise leaves the minimum distance at 0).
	if cmd.SuggestionsMinimumDistance <= 0 {
		cmd.SuggestionsMinimumDistance = 2
	}
	if suggestions := cmd.SuggestionsFor(args[0]); len(suggestions) > 0 {
		err = fmt.Errorf("%w\n\nDid you mean this?\n\t%s", err, strings.Join(suggestions, "\n\t"))
	}
	return err
}

// generateHeader creates a nice header box with proper Unicode box drawing
func generateHeader() string {
	lines := []string{
		"🚀 linctl",
		"Linear CLI - Built with ❤️",
	}

	// Find the longest line
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Add padding
	width := maxLen + 8

	// Build the box
	var result strings.Builder

	// Top border
	result.WriteString("┌")
	result.WriteString(strings.Repeat("─", width))
	result.WriteString("┐\n")

	// Content lines
	for _, line := range lines {
		padding := (width - len(line)) / 2
		result.WriteString("│")
		result.WriteString(strings.Repeat(" ", padding))
		result.WriteString(line)
		result.WriteString(strings.Repeat(" ", width-padding-len(line)))
		result.WriteString("│\n")
	}

	// Bottom border
	result.WriteString("└")
	result.WriteString(strings.Repeat("─", width))
	result.WriteString("┘")

	return result.String()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "linctl",
	Short: "A comprehensive Linear CLI tool",
	Long:  color.New(color.FgCyan).Sprintf("%s\nA comprehensive CLI tool for Linear's API featuring:\n• Issue management (create, list, update, archive)\n• Project tracking and collaboration  \n• Team and user management\n• Comments and attachments\n• Webhook configuration\n• Table/plaintext/JSON output formats\n", generateHeader()),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetRootCmd returns the root command for testing
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	// Resolve the reported version now that build info is available.
	rootCmd.Version = resolveVersion()

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.linctl.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&plaintext, "plaintext", "p", false, "plaintext output (non-interactive)")
	rootCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "JSON output")

	// Bind flags to viper
	_ = viper.BindPFlag("plaintext", rootCmd.PersistentFlags().Lookup("plaintext"))
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".linctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".linctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if !plaintext && !jsonOut {
			fmt.Fprintln(os.Stderr, color.New(color.FgGreen).Sprintf("✅ Using config file: %s", viper.ConfigFileUsed()))
		}
	}
}
