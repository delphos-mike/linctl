package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// resolveBody reads markdown/body content from, in precedence order:
//  1. the --<bodyFlag> string flag (if set),
//  2. the --<fileFlag> file path (if set),
//  3. piped stdin (if present).
//
// It returns the content and whether any source provided it. Callers use the
// "provided" boolean to decide whether to include the field in an update.
func resolveBody(cmd *cobra.Command, bodyFlag, fileFlag string) (string, bool, error) {
	if cmd.Flags().Changed(bodyFlag) {
		v, _ := cmd.Flags().GetString(bodyFlag)
		return v, true, nil
	}

	if path, _ := cmd.Flags().GetString(fileFlag); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", false, fmt.Errorf("failed to read %s: %w", fileFlag, err)
		}
		return string(data), true, nil
	}

	if isStdinPiped() {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", false, fmt.Errorf("failed to read stdin: %w", err)
		}
		if len(data) > 0 {
			return string(data), true, nil
		}
	}

	return "", false, nil
}
