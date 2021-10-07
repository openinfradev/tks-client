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
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "github.com/openinfradev/tks-proto/tks_pb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// appCreateCmd represents the create command
var appCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a TACO App.",
	Long: `Create a TACO App. supported: LMA, SERVICE_MESH

Example:
tks app create --cluster-id <CLUSTERID> --app-name <LMA,SERVICE_MESH>`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		AppName, _ := cmd.Flags().GetString("app-name")
		var Type pb.AppGroupType
		if AppName == "LMA" {
			Type = pb.AppGroupType_LMA
		} else if AppName == "SERVICE_MESH" {
			Type = pb.AppGroupType_SERVICE_MESH
		} else {
			fmt.Println("You must specify App Name. LMA or SERVICE_MESH")
			os.Exit(1)
		}

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewClusterLcmServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		ClusterId, _ := cmd.Flags().GetString("cluster-id")
		ClusterIdPre := strings.Split(ClusterId, "-")
		AppGroupName := ClusterIdPre[0] + "_" + AppName

		data := make([]pb.InstallAppGroupsRequest, 1)
		appgroups := make([]pb.AppGroup, 1)
		appgroups[0].AppGroupId = ""
		appgroups[0].AppGroupName = AppGroupName
		appgroups[0].ClusterId = ClusterId
		appgroups[0].ExternalLabel = AppGroupName
		appgroups[0].Status = pb.AppGroupStatus_APP_GROUP_INSTALLING
		appgroups[0].Type = Type
		appgroups[0].CreatedAt = timestamppb.Now()
		appgroups[0].UpdatedAt = timestamppb.Now()

		data[0].AppGroups = []*pb.AppGroup{&appgroups[0]}

		m := protojson.MarshalOptions{
			Indent:        "  ",
			UseProtoNames: true,
		}
		jsonBytes, _ := m.Marshal(&data[0])
		fmt.Println("Proto Json data...")
		fmt.Println(string(jsonBytes))
		r, err := client.InstallAppGroups(ctx, &data[0])
		fmt.Println(r)
	},
}

func init() {
	appCmd.AddCommand(appCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	appCreateCmd.Flags().String("cluster-id", "", "Cluster ID")
	appCreateCmd.MarkFlagRequired("cluster-id")
	appCreateCmd.Flags().String("app-name", "", "App Name")
	appCreateCmd.MarkFlagRequired("app-name")
}
