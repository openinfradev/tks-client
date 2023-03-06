/*
Copyright © 2022 SK Telecom <https://github.com/openinfradev>

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
package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a TKS Contract",
	Long: `Create a TKS Contract

Example:
tksadmin contract create <CONTRACT NAME>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("You must specify contract name.")
			return errors.New("Usage: tksadmin contract create <CONTRACT NAME>")
		}
		fmt.Println("Contract Name: ", args[0])
		/*
			var conn *grpc.ClientConn
			tksContractUrl = viper.GetString("tksContractUrl")
			if tksContractUrl == "" {
				return errors.New("You must specify tksContractUrl at config file")
			}
			conn, err := grpc.Dial(tksContractUrl, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %s", err)
			}
			defer conn.Close()

			client := pb.NewContractServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			creator, _ := cmd.Flags().GetString("creator")
			description, _ := cmd.Flags().GetString("description")

			data := make([]pb.CreateContractRequest, 1)
			quota := &pb.ContractQuota{}
			data[0].ContractorName = args[0]
			//TODO: this quota is dummy, so quota feature is required
			quota.Cpu = 1200
			quota.Memory = 1200
			quota.Block = 1200
			quota.BlockSsd = 0
			quota.Fs = 0
			quota.FsSsd = 0
			data[0].Quota = quota
			data[0].AvailableServices = []string{"LMA", "SERVICE_MESH"}
			data[0].CspName = "aws"
			data[0].Creator = creator
			data[0].Description = description
			m := protojson.MarshalOptions{
				Indent:        "  ",
				UseProtoNames: true,
			}
			jsonBytes, _ := m.Marshal(&data[0])
			fmt.Println("Proto Json data...")
			fmt.Println(string(jsonBytes))
			r, err := client.CreateContract(ctx, &data[0])
			if err != nil {
				return fmt.Errorf("Error: %s", err)
			} else {
				fmt.Println("Success: The request to create contract ", args[0], " was accepted.")
				fmt.Print(r)
			}
		*/
		return errors.New("MPT")
	},
}

func init() {
	//contractCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.Flags().String("creator", "", "Uuid of creator")
	createCmd.Flags().String("description", "", "Description of cluster")

}
