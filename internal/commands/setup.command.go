package commands

import (
	"errors"
	"fmt"
	"log"

	"github.com/0xmukesh/ratemywebsite/internal/helpers"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type SetupCmd struct {
	Cmd  *cobra.Command
	Args []string
}

func (c SetupCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "setup",
		Short:   "Setup your API keys for different LLMs and store it locally",
		Example: "something setup",
		Aliases: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args

			c.Handler()

			return nil
		},
	}

	return cmd
}

func (c SetupCmd) Handler() {
	var llms []string

	llmsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("something setup").Description("setup api keys of llms which you'd like to use"),
			huh.NewMultiSelect[string]().Title("choose llms").Options(
				huh.NewOption("gemini", "gemini"),
				huh.NewOption("mistral", "mistral"),
			).Value(&llms).Filterable(true),
		),
	)

	if err := llmsForm.Run(); err != nil {
		log.Fatal(err.Error())
	}

	llmsConfig := make([]helpers.LlmConfig, len(llms))

	if len(llms) != 0 {
		var fields []huh.Field

		for i := range llms {
			llmsConfig[i].Name = helpers.Llm(llms[i])

			fields = append(fields, huh.NewInput().Title(fmt.Sprintf("input your api key for %s LLM", llms[i])).Value(&llmsConfig[i].ApiKey).EchoMode(huh.EchoModePassword).Validate(func(s string) error {
				if len(s) == 0 {
					return errors.New("input an API key")
				}

				return nil
			},
			))
		}

		keysForm := huh.NewForm(huh.NewGroup(fields...))

		if err := keysForm.Run(); err != nil {
			log.Fatal(err.Error())
		}
	}

	var defaultLlm string
	var options []huh.Option[string]

	for i := range llms {
		options = append(options, huh.NewOption(llms[i], llms[i]))
	}

	defaultLlmForm := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().Title("choose your default llm").Options(options...).Value(&defaultLlm)))

	if err := defaultLlmForm.Run(); err != nil {
		log.Fatal(err.Error())
	}

	configFile := helpers.ConfigFile{
		Default: helpers.Llm(defaultLlm),
		Llms:    llmsConfig,
	}

	if err := helpers.WriteToConfigFile(configFile); err != nil {
		log.Fatal(err.Error())
	}
}
