package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	// App
	Name            string `yaml:"name" json:"name"`
	Namespace       string `yaml:"namespace" json:"namespace"`
	Type            string `yaml:"type" json:"type"`
	AppType         string `yaml:"app_type" json:"appType"`
	TargetClusterId string `yaml:"target_cluster_id" json:"targetClusterId"`

	// AppType
	Version        string `yaml:"version" json:"version"`
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
	PvEnabled      bool   `yaml:"pv_enabled" json:"pvEnabled"`
	PvStorageClass string `yaml:"pv_storage_class" json:"pvStorageClass"`
	PvAccessMode   string `yaml:"pv_access_mode" json:"pvAccessMode"`
	PvSize         string `yaml:"pv_size" json:"pvSize"`
	PvMountPath    string `yaml:"pv_mount_path" json:"pvMountPath"`

	// Update Strategy
	Promote bool `yaml:"promote" json:"promote"`
	Abort   bool `yaml:"abort" json:"abort"`
}

func NewAppserveCreateCmd(globalOpts *GlobalOptions) *cobra.Command {
	var (
		organizationId  string
		targetClusterId string
		deployType      string
		artifactUrl     string
		imageUrl        string
		namespace       string
		appType         string
		port            string
		appserveCfgFile string
		appCfgFile      string
		appSecretFile   string
	)

	var command = &cobra.Command{
		Use:   "create",
		Short: "Create an app by AppServing service",
		Long: `Create an app by AppServing service.
  
	Example:
	tks appserve create <APP_NAME> --organization-id <ORG_ID> --type <all|build|deploy> [--artifact-url <URL>|--image-url <URL>] [--namespace <NAMESPACE> --app-type <springboot|spring> --port <PORT_NUMBER> --appserve-config <CONFIGFILE>]`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var c conf
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

			name := args[0]
			if name != "" {
				c.Name = name
			}
			if organizationId == "" {
				return errors.Errorf("organization ID is mandatory")
			}
			if c.TargetClusterId == "" && targetClusterId == "" {
				return errors.Errorf("cluster ID is mandatory")
			} else if targetClusterId != "" {
				c.TargetClusterId = targetClusterId
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
			if namespace != "" {
				c.Namespace = namespace
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

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/app-serve-apps", organizationId)

			body, err := apiClient.Post(url, c)
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

	command.Flags().StringVar(&organizationId, "organization-id", "", "organization ID for AppServing service")
	command.Flags().StringVar(&targetClusterId, "target-cluster-id", "", "cluster ID for AppServing service")
	command.Flags().StringVar(&deployType, "type", "", "type for AppServing service")
	command.Flags().StringVar(&artifactUrl, "artifact-url", "", "jar url for AppServing service")
	command.Flags().StringVar(&imageUrl, "image-url", "", "image url for AppServing service")
	command.Flags().StringVar(&namespace, "namespace", "", "namespace for AppServing service")
	command.Flags().StringVar(&appType, "app-type", "", "app type for AppServing service")
	command.Flags().StringVar(&port, "port", "", "port for AppServing service")
	command.Flags().StringVar(&appserveCfgFile, "appserve-config", "", "config file for AppServing service")
	command.Flags().StringVar(&appCfgFile, "app-config", "", "custom config file for user application")
	command.Flags().StringVar(&appSecretFile, "app-secret", "", "custom secret file for user application")

	return command
}
