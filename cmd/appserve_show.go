/*
Copyright © 2021 SK Telecom <https://github.com/openinfradev>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
)

type appServeApp struct {
	Id                   string                `json:"id"`
	Name                 string                `json:"name"`
	Type                 string                `json:"type"`
	App_type             string                `json:"app_type"`
	Status               string                `json:"status"`
	Endpoint_url         string                `json:"endpoint_url"`
	Preview_endpoint_url string                `json:"preview_endpoint_url"`
	Updated_at           timestamppb.Timestamp `json:"updated_at"`
	Created_at           timestamppb.Timestamp `json:"created_at"`
}

// Members are named as snake case to avoid additional json tagging
type appServeAppTask struct {
	Id               string
	App_serve_app_id string
	Version          string
	Strategy         string
	Status           string
	Output           string
	Artifact_url     string
	Image_url        string
	Executable_path  string
	Resource_spec    string
	Profile          string
	App_config       string
	App_secret       string
	Extra_env        string
	Port             string
	Helm_revision    int32
	Created_at       timestamppb.Timestamp
	Updated_at       timestamppb.Timestamp
}

type appServeAppTaskCombined struct {
	App_serve_app appServeApp
	Tasks         []appServeAppTask
}

var appserveShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show app info deployed by AppServing service",
	Long: `show app info deployed by AppServing service.

Example:
tks appserve show <APP_UUID>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("App UUID is mandatory. Run 'tks appserve show --help'")
		}

		// Get API Url
		appserveApiUrl := viper.GetString("tksAppServeLcmUrl")
		if appserveApiUrl == "" {
			return errors.New("tks_appserve_api_url is mandatory.")
		}

		appId := args[0]

		resp, err := http.Get(appserveApiUrl + "/apps/" + appId)
		if err != nil {
			return fmt.Errorf("Error from GET req: %s", err)
		}

		defer resp.Body.Close()

		// Check response
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error while getting response: %s", err)
		} else {
			var body appServeAppTaskCombined
			json.Unmarshal(respBody, &body)

			if body.App_serve_app.Id == "" {
				fmt.Println("Failed to get app info!")
			} else {
				//fmt.Printf("response:\n%+v", body)
				printAppInfo(body, true)
			}
		}

		return nil
	},
}

func init() {
	appserveCmd.AddCommand(appserveShowCmd)

	//appserveShowCmd.Flags().StringVar(&contractId, "contract-id", "", "contract ID")
}

func printAppInfo(app appServeAppTaskCombined, long bool) error {
	cBytes, err := json.Marshal(app.App_serve_app)
	if err != nil {
		return fmt.Errorf("Unable to marshal app info to JSON: %s", err)
	}

	var m map[string]interface{}
	json.Unmarshal(cBytes, &m)

	// TODO: show fields in proper order manually
	// or use some ordered-json library
	fmt.Println("\n**************")
	fmt.Println("* Basic Info *")
	fmt.Println("**************")
	for key, val := range m {
		if key == "created_at" || key == "updated_at" {
			continue
		}
		fmt.Printf("- %s: %s\n", key, val)
	}

	tCreated := parseTime(&app.App_serve_app.Created_at)
	fmt.Printf("- created: %s\n", tCreated)
	tUpdated := parseTime(&app.App_serve_app.Updated_at)
	fmt.Printf("- updated: %s\n", tUpdated)

	fmt.Println("\n*********")
	fmt.Println("* Tasks *")
	fmt.Println("*********")

	t := table.NewWriter()
	tTemp := table.Table{}
	tTemp.Render()
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateColumns = true
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = true
	t.Style().Options.SeparateRows = true

	if long {
		t.AppendHeader(table.Row{"Version", "Status", "Strategy", "Image URL", "Profile", "CREATED_AT", "UPDATED_AT"})
		for _, s := range app.Tasks {
			tCreatedAt := parseTime(&s.Created_at)
			tUpdatedAt := parseTime(&s.Updated_at)
			t.AppendRow(table.Row{s.Version, s.Status, s.Strategy, s.Image_url, s.Profile, tCreatedAt, tUpdatedAt})
		}
	} else {
		fmt.Println("Not implemented yet.")
	}
	fmt.Println(t.Render())

	return nil
}
