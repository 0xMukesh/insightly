package commands

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/0xmukesh/ratemywebsite/internal/helpers"
	"github.com/spf13/cobra"
)

type GenerateUxReportCmd struct {
	Cmd  *cobra.Command
	Args []string
}

func (c GenerateUxReportCmd) New() *cobra.Command {
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
	cmd.Flags().BoolP("use-ai", "", false, "Use LLMs for generating a summary on how to improve the UX and accessiblity")

	return cmd
}

func (c GenerateUxReportCmd) Handler() {
	cmd := c.Cmd

	usePa11y, _ := cmd.Flags().GetBool("use-pa11y")
	saveReport, _ := cmd.Flags().GetBool("save-report")
	useAi, _ := cmd.Flags().GetBool("use-ai")

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

	if useAi {
		config, err := helpers.ReadConfigFile()
		if err != nil {
			log.Fatal(err.Error())
		}

		if config.Default != helpers.Gemini {
			log.Fatal("only gemini is supported at the moment\n")
		}

		geminiKeyFound := true
		var geminiKey string

		for i := range config.Llms {
			if config.Llms[i].Name == helpers.Gemini {
				geminiKeyFound = true
				geminiKey = config.Llms[i].ApiKey
				break
			}
		}

		if !geminiKeyFound {
			log.Fatal("only gemini is supported at the moment and config file doesn't have gemini api key\n")
		}

		output, err := helpers.SendReqToGemini(geminiKey, "Generate a simple navbar with a hamburger menu using Chakra UI")
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println(output)
	}

}
