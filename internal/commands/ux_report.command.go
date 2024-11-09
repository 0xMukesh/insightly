package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/0xmukesh/ratemywebsite/internal/helpers"
	"github.com/0xmukesh/ratemywebsite/internal/utils"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

type GenerateUxReportCmd struct {
	BaseCmd
}

func (c GenerateUxReportCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen-ux",
		Short:   "Generate UX reports",
		Example: "something gen-ux [website-url]",
		Args:    cobra.ExactArgs(1),
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
	cmd.Flags().String("llm", "", "Use other LLM than your default LLM")

	return cmd
}

func (c GenerateUxReportCmd) Handler() {
	cmd := c.Cmd
	args := c.Args

	usePa11y, _ := cmd.Flags().GetBool("use-pa11y")
	saveReport, _ := cmd.Flags().GetBool("save-report")
	useAi, _ := cmd.Flags().GetBool("use-ai")
	nonDefaultLlm, _ := cmd.Flags().GetString("llm")
	websiteUrl := args[0]

	config, err := helpers.ReadConfigFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.LogF("It seems like you're trying to run `gen-ux` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
		}

		utils.LogF(err.Error())
	}

	if len(config.Llms) == 0 {
		utils.LogF("It seems like you're trying to run `gen-ux` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
	}

	if !utils.IsValidUrl(websiteUrl) {
		utils.LogF("Invalid website URL")
	}

	if !helpers.IsNodeInstalled() {
		utils.LogF("For running UX reports, Node.js must be installed")
	}

	var accessibilityReport string
	var metrics []string

	if usePa11y {
		if !helpers.IsPa11yInstalled() {
			utils.LogF("Pa11y is not installed. Install it via running `npm install -g pa11y`")
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " generating UX report using pa11y"
		s.Start()
		pa11yReport, err := helpers.GeneratePa11yReport(websiteUrl)
		s.Stop()
		if err != nil {
			utils.LogF(err.Error())
		}

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", " ")

		if err := encoder.Encode(&pa11yReport); err != nil {
			utils.LogF(err.Error())
		}

		accessibilityReport = buffer.String()

		if saveReport {
			if err := os.WriteFile("report.json", buffer.Bytes(), 0644); err != nil {
				utils.LogF(err.Error())
			}

			fmt.Println("Saved UX reports to `report.json`")
		}

	} else {
		if !helpers.IsLighthouseInstalled() {
			utils.LogF("Lighthouse is not installed. Install it via running `npm install -g lighthouse`")
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " generating UX report using lighthouse"
		s.Start()
		lighthouseReport, err := helpers.GenerateLighthouseReport(websiteUrl)
		if err != nil {
			utils.LogF(err.Error())
		}
		s.Stop()

		var (
			total                      float64 = 0.0
			score                              = 0.0
			firstContentfulPaintTime   string  = ""
			largestContentfulPaintTime         = ""
			firstMeaningfulPaintTime           = ""
			speedIndex                         = ""
			totalBlockingTime                  = ""
		)

		var audits []helpers.LighthouseAudit

		for _, v := range lighthouseReport.Audits {
			if v.Id == "first-contentful-paint" {
				firstContentfulPaintTime = fmt.Sprintf("%.2f%s", v.NumericValue, utils.ConvertNumericUnits(v.NumericUnit))
			}

			if v.Id == "largest-contentful-paint" {
				largestContentfulPaintTime = fmt.Sprintf("%.2f%s", v.NumericValue, utils.ConvertNumericUnits(v.NumericUnit))
			}

			if v.Id == "first-meaningful-paint" {
				firstMeaningfulPaintTime = fmt.Sprintf("%.2f%s", v.NumericValue, utils.ConvertNumericUnits(v.NumericUnit))
			}

			if v.Id == "speed-index" {
				speedIndex = fmt.Sprintf("%.2f%s", v.NumericValue, utils.ConvertNumericUnits(v.NumericUnit))
			}

			if v.Id == "total-blocking-time" {
				totalBlockingTime = fmt.Sprintf("%.2f%s", v.NumericValue, utils.ConvertNumericUnits(v.NumericUnit))
			}

			if v.Score != nil {
				score += *v.Score
				total++
				audits = append(audits, v)
			}
		}

		if total != 0 {
			metrics = append(metrics, fmt.Sprintf("%.2f", score/total))
		} else {
			metrics = append(metrics, "0")
		}

		metrics = append(metrics, firstContentfulPaintTime, firstMeaningfulPaintTime, largestContentfulPaintTime, speedIndex, totalBlockingTime)

		parsedReport := struct {
			Metrics map[string]string         `json:"metrics"`
			Audits  []helpers.LighthouseAudit `json:"audits"`
		}{
			Metrics: map[string]string{
				"score":                    metrics[0],
				"first_contentful_paint":   metrics[1],
				"first_meaningful_paint":   metrics[2],
				"largest_contentful_paint": metrics[3],
				"speed_index":              metrics[4],
				"total_blocking_time":      metrics[5],
			},
			Audits: audits,
		}

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", " ")

		if err := encoder.Encode(&parsedReport); err != nil {
			utils.LogF(err.Error())
		}

		accessibilityReport = buffer.String()

		if saveReport {
			if err := os.WriteFile("report.json", buffer.Bytes(), 0644); err != nil {
				utils.LogF(err.Error())
			}
		}
	}

	if useAi {
		config, err := helpers.ReadConfigFile()
		if err != nil {
			utils.LogF(err.Error())
		}

		var llmName string
		var key string

		if nonDefaultLlm != "" {
			if nonDefaultLlm != string(helpers.Gemini) && nonDefaultLlm != string(helpers.Mistral) && nonDefaultLlm != string(helpers.Qwen) {
				utils.LogF(fmt.Sprintf("%s LLM is not supported right now. Only Gemini, Mistral and Qwen are supported at the moment", nonDefaultLlm))
			}

			apiKey, err := helpers.GetLlmKey(nonDefaultLlm)
			if err != nil {
				utils.LogF(err.Error())
			}

			llmName = nonDefaultLlm
			key = apiKey
		} else {
			if config.Default != helpers.Gemini && config.Default != helpers.Mistral && config.Default != helpers.Qwen {
				utils.LogF(fmt.Sprintf("%s LLM is not supported right now. Only Gemini, Mistral and Qwen are supported at the moment", config.Default))
			}

			apiKey, err := helpers.GetLlmKey(string(config.Default))
			if err != nil {
				utils.LogF(err.Error())
			}

			llmName = string(config.Default)
			key = apiKey
		}

		var prompt string

		if usePa11y {
			if llmName == string(helpers.Gemini) {
				prompt += `Please restructure the pa11y report into a structured format, highlighting key issues and their corresponding solutions. Employ technical language and leverage specific data from the report. In the "Additional Considerations" section, categorize recommendations based on SEO, performance, and accessibility, focusing on major and critical points. The Pa11y repot is JSON format and it contains "code", "message" and "context". Just show the structured format irrespective of whether there is a single issue. Respond in markdown. No need to re-mention the issues. Don't have any additional footer text. Don't generate a table of issues. `
				prompt += "\n"
				prompt += fmt.Sprintf("```\n%s```\n", accessibilityReport)
			} else {
				prompt += "Here is an accessibility report generated by pa11y. It is in JSON format and it contains the `code`, `message` and `context`.\n"
				prompt += fmt.Sprintf("```\n%s```\n", accessibilityReport)
				prompt += "Give suggestions regarding how to improve the accessibility and how to fix the errors mentioned by pa11y. Just only the solutions for those issues and also few other suggestions regarding how to improve the UX and accessiblity. Don't render a table of the JSON input. Just mention solution of every issue in a list style manner.\n"
				prompt += "END_OF_PROMPT"
			}
		} else {
			if llmName == string(helpers.Gemini) {
				prompt += `Please restructure the Lighthouse report into a structured format, highlighting key issues and their corresponding solutions. Employ technical language and leverage specific data from the report. In the "Additional Considerations" section, categorize recommendations based on SEO, performance, and accessibility, focusing on major and critical points. Just show the structured format irrespective of whether there is a single issue. Respond in markdown. No need to re-mention the issues. Don't generate a table of issues. Don't have any additional footer text`
				prompt += "\n"
				prompt += fmt.Sprintf("```\n%s```\n", accessibilityReport)
			} else {
				prompt += "Here is an accessibility report generated by pa11y. It is in JSON format and it contains the `code`, `message` and `context`.\n"
				prompt += fmt.Sprintf("```\n%s```\n", accessibilityReport)
				prompt += "Give suggestions regarding how to improve the accessibility and how to fix the errors mentioned by pa11y. Just only the solutions for those issues and also few other suggestions regarding how to improve the UX and accessiblity. Don't render a table of the JSON input. Just mention solution of every issue in a list style manner.\n"
				prompt += "END_OF_PROMPT"
			}
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" sending prompt to %s", llmName)
		s.Start()

		var output string

		if llmName == string(helpers.Gemini) {
			output, err = helpers.QueryGemini(key, prompt)
			if err != nil {
				utils.LogF(err.Error())
			}
		} else if llmName == string(helpers.Mistral) {
			output, err = helpers.QueryHuggingFace(key, "mistralai/Mistral-7B-Instruct-v0.3", prompt)
			if err != nil {
				utils.LogF(err.Error())
			}

			output = strings.Split(output, "END_OF_PROMPT")[1]
		} else if llmName == string(helpers.Qwen) {
			output, err = helpers.QueryHuggingFace(key, "Qwen/Qwen2.5-72B-Instruct", prompt)
			if err != nil {
				utils.LogF(err.Error())
			}

			output = strings.Split(output, "END_OF_PROMPT")[1]
		} else {
			utils.LogF("Unsupported LLM")
		}

		s.Stop()

		if !usePa11y {
			output = fmt.Sprintf(`* **Metrics**:
1. Score - %s
2. First contentful paint - %s
3. First meaningful paint - %s
4. Largest meaningful paint - %s
5. Speed index - %s
6. Total blocking time - %s`, metrics[0], metrics[1], metrics[2], metrics[3], metrics[4], metrics[5]) + "\n\n" + output
		}

		if err := helpers.DisplayInVim(output); err != nil {
			utils.LogF(err.Error())
		}
	}

}
