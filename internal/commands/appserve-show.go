package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
)

func NewAppServeShowCmd(globalOpts *GlobalOptions) *cobra.Command {
	var organizationId string
	var command = &cobra.Command{
		Use:   "show",
		Short: "show app info deployed by AppServing service",
		Long: `show app info deployed by AppServing service.
	
	Example:
	tks appserve show <APP_UUID> --organization-id <ORG_ID>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			appId := args[0]
			if len(appId) < 1 {
				return errors.New("app UUID is mandatory. Run 'tks appserve show --help'")
			}

			if organizationId == "" {
				return errors.New("--organization-id is mandatory param")
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			//url := fmt.Sprintf("organizations/%v/app-serve-apps/%s/exist", organizationId, appId)
			//body, err := apiClient.Get(url)
			//if err != nil {
			//	return err
			//}
			//
			//type DataInterface struct {
			//	Exist bool `json:"exist"`
			//}
			//var out = DataInterface{}
			//helper.Transcode(body, &out)
			//
			//fmt.Println("============================================ ")
			//fmt.Println("Json data: ")
			//data, _ := json.Marshal(out)
			//fmt.Println(string(data))
			//fmt.Println("============================================ ")

			url := fmt.Sprintf("organizations/%v/app-serve-apps/%s", organizationId, appId)
			body, err := apiClient.Get(url)
			if err != nil {
				return err
			}

			type DataInterface struct {
				AppServeApp domain.AppServeAppResponse `json:"appServeApp"`
				Stages      []domain.StageResponse     `json:"stages"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			fmt.Println("============================================ ")
			fmt.Println("Json data: ")
			data, _ := json.Marshal(out)
			fmt.Println(string(data))
			fmt.Println("============================================ ")

			printAppServeShow(out.AppServeApp, true)

			return nil
		},
	}

	command.Flags().StringVar(&organizationId, "organization-id", "", "a organizationId")

	return command
}

func printAppServeShow(d domain.AppServeAppResponse, long bool) {
	t := table.NewWriter()
	tTemp := table.Table{}
	tTemp.Render()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	if long {
		t.AppendHeader(table.Row{"ID", "Version", "Status", "Available Rollback", "Strategy", "Revision", "Image URL", "Profile", "CREATED_AT", "UPDATED_AT"})
		for _, i := range d.AppServeAppTasks {
			tCreatedAt := helper.ParseTime(i.CreatedAt)
			var tUpdatedAt string
			if i.UpdatedAt != nil {
				tUpdatedAt = helper.ParseTime(*i.UpdatedAt)
			}
			t.AppendRow(table.Row{i.ID, i.Version, i.Status, i.AvailableRollback, i.Strategy, strconv.Itoa(int(i.HelmRevision)), i.ImageUrl, i.Profile, tCreatedAt, tUpdatedAt})
		}
	} else {
		fmt.Println("Not implemented yet.")
	}
	fmt.Println(t.Render())
}
