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
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/openinfradev/tks-proto/tks_pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// serviceCreateCmd represents the create command
var serviceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a TKS Service.",
	Long: `Create a TKS Service. supported: LMA,LMA_EFK,SERVICE_MESH

Example:
tks service create --cluster-id <CLUSTERID> --service-name <LMA,LMA_EFK,SERVICE_MESH>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("create called")
		ServiceName, _ := cmd.Flags().GetString("service-name")
		var Type pb.AppGroupType
		if ServiceName == "LMA" {
			Type = pb.AppGroupType_LMA
		} else if ServiceName == "LMA_EFK" {
			Type = pb.AppGroupType_LMA_EFK
		} else if ServiceName == "SERVICE_MESH" {
			Type = pb.AppGroupType_SERVICE_MESH
		} else {
			return errors.New("You must specify Service Name. LMA | LMA_EFK | SERVICE_MESH")
		}

		var conn *grpc.ClientConn
		tksClusterLcmUrl = viper.GetString("tksClusterLcmUrl")
		if tksClusterLcmUrl == "" {
			return errors.New("You must specify tksClusterLcmUrl at config file")
		}
		conn, err := grpc.Dial(tksClusterLcmUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewClusterLcmServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		ClusterId, _ := cmd.Flags().GetString("cluster-id")
		ClusterIdPre := strings.Split(ClusterId, "-")
		AppGroupName := ClusterIdPre[0] + "_" + ServiceName

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
		fmt.Println("Response:\n", r)
		if err != nil {
			return fmt.Errorf("Error: %s", err)
		} else {
			fmt.Println("Success: The request to create service ", AppGroupName, " was accepted.")
		}
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(serviceCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serviceCreateCmd.Flags().String("cluster-id", "", "Cluster ID")
	_ = serviceCreateCmd.MarkFlagRequired("cluster-id")
	serviceCreateCmd.Flags().String("service-name", "", "Service Name")
	_ = serviceCreateCmd.MarkFlagRequired("service-name")
}
