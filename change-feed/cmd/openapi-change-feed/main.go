package main

import (
	"fmt"
	"os"

	"github.com/github/openapi-change-feed/internal/detector"
	"github.com/github/openapi-change-feed/internal/monitor"
	"github.com/github/openapi-change-feed/internal/output"
	"github.com/github/openapi-change-feed/internal/specfetch"
	"github.com/spf13/cobra"
)

func versionString() string { return "openapi-change-feed v0.0.0-dev" }

func newDetectCmd() *cobra.Command {
	var specURL, out, cursor string
	cmd := &cobra.Command{
		Use:   "detect",
		Short: "Diff the baseline spec against HEAD and emit the change feed",
		RunE: func(cmd *cobra.Command, _ []string) error {
			m := monitor.Monitor{
				Fetcher: specfetch.New(nil), Detector: detector.New(),
				Emitter: output.NewFileEmitter(out), CursorPath: cursor, SpecURL: specURL,
			}
			sum, ran, err := m.RunOnce(cmd.Context())
			if err != nil {
				return err
			}
			if !ran {
				fmt.Fprintln(cmd.OutOrStdout(), "spec unchanged; nothing to do")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%d changes (%d breaking)\n", sum.Total, sum.Breaking)
			return nil
		},
	}
	cmd.Flags().StringVar(&specURL, "spec-url", "", "URL of the OpenAPI spec")
	cmd.Flags().StringVar(&out, "out", "out", "output directory for feed artifacts")
	cmd.Flags().StringVar(&cursor, "cursor", "state/cursor.json", "path to the cursor state file")
	_ = cmd.MarkFlagRequired("spec-url")
	return cmd
}

func main() {
	root := &cobra.Command{Use: "openapi-change-feed", Version: versionString()}
	root.AddCommand(newDetectCmd())
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
