package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPickVersion(t *testing.T) {
	tests := []struct {
		name      string
		injected  string
		buildInfo string
		want      string
	}{
		{
			name:      "ldflags value wins",
			injected:  "v0.4.0",
			buildInfo: "(devel)",
			want:      "v0.4.0",
		},
		{
			name:      "ldflags value wins over build info",
			injected:  "abc1234",
			buildInfo: "v1.2.3",
			want:      "abc1234",
		},
		{
			name:      "go install build info used when no ldflags",
			injected:  "dev",
			buildInfo: "v0.4.0",
			want:      "v0.4.0",
		},
		{
			name:      "build info used when injected is empty",
			injected:  "",
			buildInfo: "v0.4.0",
			want:      "v0.4.0",
		},
		{
			name:      "devel build info ignored",
			injected:  "dev",
			buildInfo: "(devel)",
			want:      "dev",
		},
		{
			name:      "no version info falls back to dev",
			injected:  "dev",
			buildInfo: "",
			want:      "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pickVersion(tt.injected, tt.buildInfo); got != tt.want {
				t.Errorf("pickVersion(%q, %q) = %q, want %q", tt.injected, tt.buildInfo, got, tt.want)
			}
		})
	}
}

// newTestGroup builds a command group wired like the real resource groups:
// a parent using requireSubcommand plus one real subcommand.
func newTestGroup() *cobra.Command {
	parent := &cobra.Command{
		Use:          "widget",
		Short:        "Manage widgets",
		SilenceUsage: true,
		RunE:         requireSubcommand,
	}
	parent.AddCommand(&cobra.Command{
		Use: "list",
		Run: func(*cobra.Command, []string) {},
	})
	return parent
}

func TestRequireSubcommand_BarePrintsHelpNoError(t *testing.T) {
	parent := newTestGroup()
	var out bytes.Buffer
	parent.SetOut(&out)
	parent.SetErr(&out)

	if err := requireSubcommand(parent, []string{}); err != nil {
		t.Fatalf("bare invocation returned error %v, want nil", err)
	}
	if !strings.Contains(out.String(), "Manage widgets") {
		t.Errorf("bare invocation did not print help; got %q", out.String())
	}
}

func TestRequireSubcommand_UnknownArgErrors(t *testing.T) {
	parent := newTestGroup()

	err := requireSubcommand(parent, []string{"bogus"})
	if err == nil {
		t.Fatal("unknown subcommand returned nil error, want non-nil")
	}
	if !strings.Contains(err.Error(), `unknown command "bogus"`) {
		t.Errorf("error = %q, want it to mention unknown command %q", err.Error(), "bogus")
	}
	if !strings.Contains(err.Error(), "widget") {
		t.Errorf("error = %q, want it to name the parent command path", err.Error())
	}
}

func TestRequireSubcommand_SuggestsNearMiss(t *testing.T) {
	parent := newTestGroup()

	err := requireSubcommand(parent, []string{"lst"})
	if err == nil {
		t.Fatal("near-miss subcommand returned nil error, want non-nil")
	}
	if !strings.Contains(err.Error(), "list") {
		t.Errorf("error = %q, want a suggestion of %q", err.Error(), "list")
	}
}

// TestGroupUnknownSubcommandExitsNonZero exercises the full command tree
// through Execute to prove a representative resource group rejects an
// unrecognized subcommand (non-nil error -> non-zero exit) while a bare
// invocation succeeds.
func TestGroupUnknownSubcommandExitsNonZero(t *testing.T) {
	root := GetRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)

	root.SetArgs([]string{"project", "definitely-not-a-command"})
	if err := root.Execute(); err == nil {
		t.Fatal("`project <unknown>` returned nil error; expected a non-zero exit")
	}

	out.Reset()
	root.SetArgs([]string{"project"})
	if err := root.Execute(); err != nil {
		t.Fatalf("bare `project` returned error %v; expected help + success", err)
	}
	if !strings.Contains(out.String(), "Manage Linear projects") {
		t.Errorf("bare `project` did not print help; got %q", out.String())
	}
}
