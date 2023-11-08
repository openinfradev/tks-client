package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type updateConf struct {
	// App
	Type    string `yaml:"type" json:"type"`
	AppType string `yaml:"app_type" json:"appType"`

	// AppTask
	Strategy       string `yaml:"strategy" json:"strategy"`
	ArtifactUrl    string `yaml:"artifact_url" json:"artifactUrl"`
	ImageUrl       string `yaml:"image_url" json:"imageUrl"`
	ExecutablePath string `yaml:"executable_path" json:"executablePath"`
	ResourceSpec   string `yaml:"resource_spec" json:"resourceSpec"`
	Profile        string `yaml:"profile" json:"profile"`
	AppConfig      string `yaml:"app_config" json:"appConfig"`
	AppSecret      string `yaml:"app_secret" json:"appSecret"`
	ExtraEnv       string `yaml:"extra_env" json:"extraEnv"`
	Port           string `yaml:"port" json:"port"`
}

func NewAppserveUpdateCmd(globalOpts *GlobalOptions) *cobra.Command {
	var (
		organizationId  string
		deployType      string
		artifactUrl     string
		imageUrl        string
		strategy        string
		appType         string
		appserveCfgFile string
		appCfgFile      string
		appSecretFile   string
		port            string
	)

	var command = &cobra.Command{
		Use:   "update",
		Short: "Update an app deployed by AppServing service",
		Long: `Update an app deployed by AppServing service.
  
	Example:
	tks appserve update <APP_ID> --organization-id <ORG_ID> --type <all|build|deploy> [--artifact-url <URL>|--image-url <URL>] [--strategy <rolling-update|blue-green> --app-type <springboot|spring> --port <PORT_NUMBER> --appserve-config <CONFIGFILE>]`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var c updateConf
			var err error
			var yamlData []byte

			if appserveCfgFile != "" {
				// Get Appserving request params from config file
				yamlData, err = os.ReadFile(appserveCfgFile)
				if err != nil {
					return fmt.Errorf("error: %s", err)
				}
				fmt.Printf("*******\nConfig:\n%s\n*******\n", yamlData)

				// Unmarshal yaml content into struct
				if err = yaml.Unmarshal(yamlData, &c); err != nil {
					return fmt.Errorf("error: %s", err)
				}
			}

			// Get application config from file
			if appCfgFile != "" {
				appCfgBytes, err := os.ReadFile(appCfgFile)
				if err != nil {
					return fmt.Errorf("error: %s", err)
				}
				// Add appCfg to existing struct
				c.AppConfig = string(appCfgBytes)
			}

			// Get application secret from file
			if appSecretFile != "" {
				appSecretBytes, err := os.ReadFile(appSecretFile)
				if err != nil {
					return fmt.Errorf("error: %s", err)
				}
				c.AppSecret = string(appSecretBytes)
			}

			appId := args[0]
			if len(appId) < 1 {
				return errors.New("APP_ID is mandatory! Run 'tks appserve update --help'")
			}

			if organizationId == "" {
				return errors.New("organization ID is mandatory")
			}
			if deployType != "" {
				c.Type = deployType
			}
			if artifactUrl != "" {
				c.ArtifactUrl = artifactUrl
			}
			if imageUrl != "" {
				c.ImageUrl = imageUrl
			}
			if strategy != "" {
				c.Strategy = strategy
			}
			if appType != "" {
				c.AppType = appType
			}
			if port != "" {
				c.Port = port
			}

			// Convert map to Json
			cBytes, err := json.Marshal(&c)
			if err != nil {
				return fmt.Errorf("unable to marshal config to JSON: %s", err)
			}

			fmt.Println("========== ")
			fmt.Println("Json data: ")
			fmt.Println(string(cBytes))
			fmt.Println("========== ")

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/app-serve-apps/%v", organizationId, appId)

			body, err := apiClient.Put(url, c)
			if err != nil {
				return err
			}

			fmt.Println("Response:\n", fmt.Sprintf("%v", body))

			return nil
		},
	}

	command.Flags().StringVar(&organizationId, "organization-id", "", "organization ID for AppServing service")
	command.Flags().StringVar(&deployType, "type", "", "type for AppServing service")
	command.Flags().StringVar(&artifactUrl, "artifact-url", "", "jar url for AppServing service")
	command.Flags().StringVar(&imageUrl, "image-url", "", "image url for AppServing service")
	command.Flags().StringVar(&strategy, "strategy", "", "strategy for AppServing service")
	command.Flags().StringVar(&appType, "app-type", "", "app type for AppServing service")
	command.Flags().StringVar(&port, "port", "", "port for AppServing service")
	command.Flags().StringVar(&appserveCfgFile, "appserve-config", "", "config file for AppServing service")
	command.Flags().StringVar(&appCfgFile, "app-config", "", "custom config file for user application")
	command.Flags().StringVar(&appSecretFile, "app-secret", "", "custom secret file for user application")

	return command
}
