package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/0xmukesh/insightly/internal/helpers"
	"github.com/0xmukesh/insightly/internal/helpers/styles"
	"github.com/0xmukesh/insightly/internal/utils"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type ConfigCmd struct {
	BaseCmd
}
type ConfigViewCmd struct {
	BaseCmd
}
type ConfigSetDefaultCmd struct {
	BaseCmd
}
type ConfigSetCmd struct {
	BaseCmd
}
type ConfigRemoveCmd struct {
	BaseCmd
}

func (c ConfigCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "View/change configuration",
		Example: "insightly config [command]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return nil
		},
	}

	configViewCmd := ConfigViewCmd{}
	configSetDefaultCmd := ConfigSetDefaultCmd{}
	configSetCmd := ConfigSetCmd{}
	configRemoveCmd := ConfigRemoveCmd{}

	cmd.AddCommand(configViewCmd.New())
	cmd.AddCommand(configSetDefaultCmd.New())
	cmd.AddCommand(configSetCmd.New())
	cmd.AddCommand(configRemoveCmd.New())

	return cmd
}

func (c ConfigViewCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "View configuration details",
		Example: "insightly config view",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args

			c.Handler()

			return nil
		},
	}

	cmd.Flags().BoolP("show-full", "", false, "Show the complete the API key")

	return cmd
}

func (c ConfigViewCmd) Handler() {
	cmd := c.Cmd

	showFull, _ := cmd.Flags().GetBool("show-full")

	config, err := helpers.ReadConfigFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.LogF("It seems like you're trying to run `config view` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
		}

		utils.LogF(err.Error())
	}

	if len(config.Llms) == 0 {
		utils.LogF("It seems like you're trying to run `config view` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
	}

	fmt.Printf("You're using %s as your default LLM\n", styles.BoldBlueTextStyle.Render(string(config.Default)))
	fmt.Printf("Here are your configuration details for each of the LLM:\n")

	for i := range config.Llms {
		var key string

		if showFull {
			key = config.Llms[i].ApiKey
		} else {
			key = utils.MaskApiKey(config.Llms[i].ApiKey)
		}

		fmt.Printf(">> %s - %s\n", config.Llms[i].Name, key)
	}
}

func (c ConfigSetDefaultCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-default",
		Short:   "Change your default LLM",
		Example: "insightly config set-default <llm>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args

			c.Handler()

			return nil
		},
	}

	return cmd
}

func (c ConfigSetDefaultCmd) Handler() {
	args := c.Args

	newDefaultLlm := args[0]

	config, err := helpers.ReadConfigFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.LogF("It seems like you're trying to run `config set-default` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
		}

		utils.LogF(err.Error())
	}

	if len(config.Llms) == 0 {
		utils.LogF("It seems like you're trying to run `config set-default` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
	}

	if newDefaultLlm != "" {
		found := false

		for i := range config.Llms {
			if newDefaultLlm == string(config.Llms[i].Name) {
				found = true
			}
		}

		if !found {
			utils.LogF(fmt.Sprintf("Can't set %s as default LLM cause its' configuration can't be found. Run `config set` to set an LLM's configuration", newDefaultLlm))
		}

		config.Default = helpers.Llm(newDefaultLlm)
	}

	if err := helpers.WriteToConfigFile(config); err != nil {
		utils.LogF(err.Error())
	}
}

func (c ConfigSetCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Short:   "Update configuration details",
		Example: "insightly config set",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args

			c.Handler()

			return nil
		},
	}

	return cmd
}

func (c ConfigSetCmd) Handler() {
	config, err := helpers.ReadConfigFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.LogF("It seems like you're trying to run `config set` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
		}

		utils.LogF(err.Error())
	}

	if len(config.Llms) == 0 {
		utils.LogF("It seems like you're trying to run `config set` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
	}

	var llmsFormFields []huh.Option[string]

	for i := range config.Llms {
		llmsFormFields = append(llmsFormFields, huh.NewOption(string(config.Llms[i].Name), string(config.Llms[i].Name)))
	}

	var selectedLlms []string

	llmsForm := huh.NewForm(huh.NewGroup(
		huh.NewNote().Title("something config set").Description("choose the models whose api key you would like to update"),
		huh.NewMultiSelect[string]().Title("choose llms").Options(llmsFormFields...,
		).Value(&selectedLlms).Filterable(true).Validate(func(s []string) error {
			if len(s) == 0 {
				return errors.New("atleast select one model")
			}

			return nil
		}),
	))

	if err := llmsForm.Run(); err != nil {
		utils.LogF(err.Error())
	}

	var keysFormFields []huh.Field

	for i := range selectedLlms {
		keysFormFields = append(keysFormFields, huh.NewInput().Title(fmt.Sprintf("input your api key for %s LLM", selectedLlms[i])).Value(&config.Llms[i].ApiKey).EchoMode(huh.EchoModePassword).Validate(func(s string) error {
			if len(s) == 0 {
				return errors.New("input an API key")
			}

			return nil
		},
		))
	}

	keysForm := huh.NewForm(huh.NewGroup(keysFormFields...))

	if err := keysForm.Run(); err != nil {
		utils.LogF(err.Error())
	}

	if err := helpers.WriteToConfigFile(config); err != nil {
		utils.LogF(err.Error())
	}

	fmt.Println("Successfully updated keys for the given models")
}

func (c ConfigRemoveCmd) New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove configuration details of a certain LLM",
		Example: "insightly config remove",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args

			c.Handler()

			return nil
		},
	}

	return cmd
}

func (c ConfigRemoveCmd) Handler() {
	config, err := helpers.ReadConfigFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.LogF("It seems like you're trying to run `config remove` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
		}

		utils.LogF(err.Error())
	}

	if len(config.Llms) == 0 {
		utils.LogF("It seems like you're trying to run `config remove` command before setting up your LLM configuration. Run `setup` to setup your LLM configuration")
	}

	var llmsFormFields []huh.Option[string]

	for i := range config.Llms {
		llmsFormFields = append(llmsFormFields, huh.NewOption(string(config.Llms[i].Name), string(config.Llms[i].Name)))
	}

	var selectedLlms []string

	llmsForm := huh.NewForm(huh.NewGroup(
		huh.NewNote().Title("something config remove").Description("choose the models whose configuration you would like to be removed"),
		huh.NewMultiSelect[string]().Title("choose llms").Options(llmsFormFields...,
		).Value(&selectedLlms).Filterable(true).Validate(func(s []string) error {
			if len(s) == 0 {
				return errors.New("atleast select one model")
			}

			return nil
		}),
	))

	if err := llmsForm.Run(); err != nil {
		utils.LogF(err.Error())
	}

	for i := len(config.Llms) - 1; i >= 0; i-- {
		for _, selected := range selectedLlms {
			if config.Llms[i].Name == helpers.Llm(selected) {
				config.Llms = append(config.Llms[:i], config.Llms[i+1:]...)
				break
			}
		}
	}

	if err := helpers.WriteToConfigFile(config); err != nil {
		utils.LogF(err.Error())
	}
}
