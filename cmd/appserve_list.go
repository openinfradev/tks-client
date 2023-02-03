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

var contractId string

type app struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	App_type   string `json:"app_type"`
	Status     string `json:"status"`
	Image_url  string `json:"image_url"`
	Updated_at timestamppb.Timestamp
	Created_at timestamppb.Timestamp
}

var appserveListCmd = &cobra.Command{
	Use:   "list",
	Short: "list apps deployed by AppServing service",
	Long: `list apps deployed by AppServing service.
  
Example:
tks appserve list --contract_id <CONTRACT_ID>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get API Url
		appserveApiUrl := viper.GetString("tksAppServeLcmUrl")
		if appserveApiUrl == "" {
			return errors.New("tks_appserve_api_url is mandatory.")
		}

		if contractId == "" {
			return errors.New("contract-id is mandatory param.")
		}

		resp, err := http.Get(appserveApiUrl + "/apps?contract_id=" + contractId)
		if err != nil {
			return fmt.Errorf("Error from GET req: %s", err)
		}

		defer resp.Body.Close()

		// Check response
		respBody, err := io.ReadAll(resp.Body)
		if err == nil {
			var body []app
			json.Unmarshal(respBody, &body)

			if len(body) == 0 {
				fmt.Println("No app exists for specified contract!")
			} else {
				printApps(body, true)
			}
		}

		return nil
	},
}

func init() {
	appserveCmd.AddCommand(appserveListCmd)

	appserveListCmd.Flags().StringVar(&contractId, "contract-id", "", "contract ID")
}

func printApps(apps []app, long bool) {
	t := table.NewWriter()
	tTemp := table.Table{}
	tTemp.Render()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = true
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = true
	t.Style().Options.SeparateRows = false

	if long {
		t.AppendHeader(table.Row{"ID", "Name", "Status", "TYPE", "APP_TYPE", "CREATED_AT", "UPDATED_AT"})
		for _, s := range apps {
			tCreatedAt := parseTime(&s.Created_at)
			tUpdatedAt := parseTime(&s.Updated_at)
			t.AppendRow(table.Row{s.Id, s.Name, s.Status, s.Type, s.App_type, tCreatedAt, tUpdatedAt})
		}
	} else {
		fmt.Println("Not implemented yet.")
	}
	fmt.Println(t.Render())
}
