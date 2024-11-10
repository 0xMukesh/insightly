package commands

import (
	"errors"
	"fmt"

	"github.com/0xmukesh/insightly/internal/helpers"
	"github.com/0xmukesh/insightly/internal/utils"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type SetupCmd struct {
	BaseCmd
}

func (c SetupCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "setup",
		Short:   "Setup your API keys for different LLMs and store it locally",
		Example: "something setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args

			if err := c.Handler(); err != nil {
				utils.LogF(err.Error())
			}

			return nil
		},
	}

	return cmd
}

func (c SetupCmd) Handler() error {
	var selectedLlms []string

	llmsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("something setup").Description("setup api keys of llms which you'd like to use"),
			huh.NewMultiSelect[string]().Title("choose llms").Options(
				huh.NewOption("gemini 1.5 flash", "gemini"),
				huh.NewOption("mistral 7b instruct", "mistral"),
				huh.NewOption("qwen 2.5 7b instruct", "qwen"),
				huh.NewOption("llama 3.1", "llama"),
				huh.NewOption("claude 3.5 sonnet", "claude"),
				huh.NewOption("chatgpt 4o", "chatgpt"),
			).Value(&selectedLlms).Filterable(true).Validate(func(s []string) error {
				if len(s) == 0 {
					return errors.New("atleast select one model")
				}

				return nil
			}),
		),
	)

	if err := llmsForm.Run(); err != nil {
		utils.LogF(err.Error())
	}

	llmsConfig := make([]helpers.LlmConfig, len(selectedLlms))

	var keysFormFields []huh.Field

	for i := range selectedLlms {
		llmsConfig[i].Name = helpers.Llm(selectedLlms[i])

		keysFormFields = append(keysFormFields, huh.NewInput().Title(fmt.Sprintf("input your api key for %s LLM", selectedLlms[i])).Value(&llmsConfig[i].ApiKey).EchoMode(huh.EchoModePassword).Validate(func(s string) error {
			if len(s) == 0 {
				return errors.New("input an API key")
			}

			return nil
		},
		))
	}

	keysForm := huh.NewForm(huh.NewGroup(keysFormFields...))

	if err := keysForm.Run(); err != nil {
		return err
	}

	var defaultLlm string
	var options []huh.Option[string]

	for i := range selectedLlms {
		options = append(options, huh.NewOption(selectedLlms[i], selectedLlms[i]))
	}

	defaultLlmForm := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().Title("choose your default llm").Options(options...).Value(&defaultLlm).Validate(func(s string) error {
		if s == "" {
			return errors.New("choose a llm as your default llm")
		}

		return nil
	})))

	if err := defaultLlmForm.Run(); err != nil {
		return err
	}

	configFile := helpers.ConfigFile{
		Default: helpers.Llm(defaultLlm),
		Llms:    llmsConfig,
	}

	if err := helpers.WriteToConfigFile(configFile); err != nil {
		return err
	}

	fmt.Printf("Successfully saved LLM configuration with %s as your default LLM. Now you can use commands like `gen-ux\n`", defaultLlm)

	return nil
}
