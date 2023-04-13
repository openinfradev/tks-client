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

type conf struct {
	// App
	Name            string `yaml:"name" json:"name"`
	OrganizationId  string `yaml:"organization_id" json:"organization_id"`
	Type            string `yaml:"type" json:"type"`
	AppType         string `yaml:"app_type" json:"app_type"`
	TargetClusterId string `yaml:"target_cluster_id" json:"target_cluster_id"`

	// AppType
	Version        string `yaml:"version" json:"version"`
	Strategy       string `yaml:"strategy" json:"strategy"`
	ArtifactUrl    string `yaml:"artifact_url" json:"artifact_url"`
	ImageUrl       string `yaml:"image_url" json:"image_url"`
	ExecutablePath string `yaml:"executable_path" json:"executable_path"`
	ResourceSpec   string `yaml:"resource_spec" json:"resource_spec"`
	Profile        string `yaml:"profile" json:"profile"`
	AppConfig      string `yaml:"app_config" json:"app_config"`
	AppSecret      string `yaml:"app_secret" json:"app_secret"`
	ExtraEnv       string `yaml:"extra_env" json:"extra_env"`
	Port           string `yaml:"port" json:"port"`
	PvEnabled      bool   `yaml:"pv_enabled" json:"pv_enabled"`
	PvStorageClass string `yaml:"pv_storage_class" json:"pv_storage_class"`
	PvAccessMode   string `yaml:"pv_access_mode" json:"pv_access_mode"`
	PvSize         string `yaml:"pv_size" json:"pv_size"`
	PvMountPath    string `yaml:"pv_mount_path" json:"pv_mount_path"`

	// Update Strategy
	Promote bool `yaml:"promote" json:"promote"`
	Abort   bool `yaml:"abort" json:"abort"`
}

func NewAppserveCreateCmd(globalOpts *GlobalOptions) *cobra.Command {
	var (
		appserveCfgFile string
		appCfgFile      string
		appSecretFile   string
	)

	var command = &cobra.Command{
		Use:   "create",
		Short: "Create an app by AppServing service",
		Long: `Create an app by AppServing service.
  
	Example:
	tks appserve create --appserve-config CONFIGFILE`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var c conf
			if appserveCfgFile == "" {
				return errors.New("--appserve-config is mandatory param")
			}

			// Get Appserving request params from config file
			yamlData, err := os.ReadFile(appserveCfgFile)
			if err != nil {
				return fmt.Errorf("error: %s", err)
			}

			fmt.Printf("*******\nConfig:\n%s\n*******\n", yamlData)

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

			// Unmarshal yaml content into struct
			if err = yaml.Unmarshal(yamlData, &c); err != nil {
				return fmt.Errorf("error: %s", err)
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

			//buff := bytes.NewBuffer(cBytes)
			//resp, err := http.Post(appserveApiUrl+"/apps", "application/json", buff)
			//if err != nil {
			//	return fmt.Errorf("Error from POST req: %s", err)
			//}
			//
			//defer resp.Body.Close()
			//
			//// Check response
			//respBody, err := io.ReadAll(resp.Body)
			//if err == nil {
			//	str := string(respBody)
			//	fmt.Println("Response:\n", str)
			//}
			//

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Post("app-serve-apps", c)
			if err != nil {
				return err
			}

			fmt.Println("Response:\n", fmt.Sprintf("%v", body))

			//type DataInterface struct {
			//	ClusterId string `json:"clusterId"`
			//}
			//var out = DataInterface{}
			//helper.Transcode(body, &out)

			return nil
		},
	}

	command.Flags().StringVar(&appserveCfgFile, "appserve-config", "", "config file for AppServing service")
	command.Flags().StringVar(&appCfgFile, "app-config", "", "custom config file for user application")
	command.Flags().StringVar(&appSecretFile, "app-secret", "", "custom secret file for user application")

	return command
}
