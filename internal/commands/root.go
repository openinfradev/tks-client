package commands

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/openinfradev/tks-client/internal/config"
	"github.com/openinfradev/tks-client/internal/helper"
)

type GlobalOptions struct {
	ServerAddr string
	AuthToken  string
	ConfigPath string
}

type LocalConfig struct {
	Server string `yaml:"server"`
	Token  string `yaml:"token"`
}

func NewCommand() *cobra.Command {
	var (
		opts GlobalOptions
	)

	var command = &cobra.Command{
		Use:   "tks",
		Short: "CLI Client for TKS Service",
		Long: ` ______ __ __ ____  ___      __        _         _____ __ _            __ 
	/_  __// //_// __/ / _ | ___/ /__ _   (_)___    / ___// /(_)___  ___  / /_
	 / /  / ,<  _\ \  / __ |/ _  //  ' \ / // _ \  / /__ / // // -_)/ _ \/ __/
	/_/  /_/|_|/___/ /_/ |_|\_,_//_/_/_//_//_//_/  \___//_//_/ \__//_//_/\__/ 
																			  
	TKS Client is CLI client for TKS Service.
	For more: https://github.com/openinfradev/tks-client/`,
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
		SilenceUsage:      true,
	}

	command.AddCommand(NewLoginCommand(&opts))
	command.AddCommand(NewOrganizationCommand(&opts))
	command.AddCommand(NewClusterCommand(&opts))
	command.AddCommand(NewAppGroupCommand(&opts))
	command.AddCommand(NewCloudSettingCommand(&opts))
	command.AddCommand(NewAppserveCommand(&opts))
	command.AddCommand(NewStackTemplateCommand(&opts))

	defaultLocalConfigPath, err := config.DefaultLocalConfigPath()
	helper.CheckError(err)

	localCfg, err := config.ReadLocalConfig(defaultLocalConfigPath)
	helper.CheckError(err)

	command.PersistentFlags().StringVar(&opts.ConfigPath, "config", config.GetFlag("config", defaultLocalConfigPath), "Path to TKS config")

	if localCfg != nil {
		command.PersistentFlags().StringVar(&opts.ServerAddr, "server", localCfg.GetServer().Server, "TKS server address")
		command.PersistentFlags().StringVar(&opts.AuthToken, "auth-token", localCfg.GetUser().AuthToken, "Authentication token")
	}

	// Set hidden auth-token
	command.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		name := flag.Name
		if name == "auth-token" {
			flag.Hidden = true
		}
	})

	return command
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		os.Exit(-1)
	}
}
