package commands

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/0xmukesh/ratemywebsite/internal/helpers"
	"github.com/spf13/cobra"
)

type GenerateUxReport struct {
	Cmd  *cobra.Command
	Args []string
}

func (c GenerateUxReport) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen-ux",
		Short:   "Generate UX reports for your React.js project",
		Example: "something gen-ux",
		Aliases: []string{"generate-ux", "gux"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args

			c.Handler()

			return nil
		},
	}

	cmd.Flags().BoolP("use-pa11y", "", false, "Use pa11y for running accessibility report")
	cmd.Flags().BoolP("save-report", "", false, "Save parsed report in JSON format")

	return cmd
}

func (c GenerateUxReport) Handler() {
	cmd := c.Cmd

	usePa11y, err := cmd.Flags().GetBool("use-pa11y")
	if err != nil {
		log.Fatal(err.Error())
	}
	saveReport, err := cmd.Flags().GetBool("save-report")
	if err != nil {
		log.Fatal(err.Error())
	}

	if !helpers.IsNodeInstalled() {
		log.Fatal("node is not installed")
	}

	if usePa11y {
		if !helpers.IsPa11yInstalled() {
			log.Fatal("pa11y is not installed")
		}

		report, err := helpers.RunPa11yReport("https://example.com")
		if err != nil {
			log.Fatal(err.Error())
		}

		if saveReport {
			buffer := &bytes.Buffer{}
			encoder := json.NewEncoder(buffer)
			encoder.SetEscapeHTML(false)
			encoder.SetIndent("", "    ")

			err := encoder.Encode(&report)
			if err != nil {
				log.Fatal(err.Error())
			}

			if err := os.WriteFile("report.json", buffer.Bytes(), 0644); err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}
