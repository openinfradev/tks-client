package commands

import (
	"fmt"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewStackTemplateListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		all bool
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of stack-template.",
		Long: `Show list of stack-template.
	
	Example:
	tks stack-template list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Get("stack-templates?all=" + strconv.FormatBool(all))
			if err != nil {
				return err
			}

			var out = domain.GetStackTemplatesResponse{}
			helper.Transcode(body, &out)

			printStackTemplates(out.StackTemplates)

			return nil
		},
	}

	command.Flags().BoolVarP(&all, "all", "A", false, "show all organizations")

	return command
}

func printStackTemplates(r []domain.StackTemplateResponse) {
	if len(r) == 0 {
		fmt.Println("No stackTemplate exists for user organization!")
		return
	}

	t := table.NewWriter()
	tTemp := table.Table{}
	tTemp.Render()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"ID", "NAME", "DESCRIPTION", "VERSION", "CLOUD_SERVICE", "PLATFORM", "TEMPLATE", "CREATED_AT", "UPDATED_AT"})
	for _, s := range r {
		tCreatedAt := helper.ParseTime(s.CreatedAt)
		tUpdatedAt := helper.ParseTime(s.UpdatedAt)
		t.AppendRow(table.Row{s.ID, s.Name, s.Description, s.Version, s.CloudService, s.Platform, s.Template, tCreatedAt, tUpdatedAt})
	}
	fmt.Println(t.Render())
}
