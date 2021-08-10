package main

import (
    "context"
    "flag"
    "strings"
    "fmt"
    "os"
    "os/exec"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    "github.com/sktelecom/tks-contract/pkg/log"
    "github.com/sktelecom/tks-info/pkg/cert"
    pb "github.com/sktelecom/tks-proto/pbgo"
)

var (
    tks_info_port   = flag.Int("tks_info_port", 9111, "TKS-info gRPC server port")
    tls    = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
    caFile = flag.String("ca_file", "", "The TLS ca file")
)

func main() {
    if len(os.Args) < 2 {
        panic("Need parameters. (help with -h)")
    }
    log.Info("Preporcess for TACO-LMA with TKS-info")

    tks_info_host := flag.String("tks_info_host", "127.0.0.1", "An address of TKS-info")
    clusterid := flag.String("clusterid", "", "Cluster ID to apply. eg. 6abead61-ff2a-4af4-8f41-d2c44c745de7")
    appgroupid := flag.String("appgroupid", "", "Application ID of The endpoint. eg. 6abead61-ff2a-4af4-8f41-d2c44c745de7")
    curEndpoint := flag.String("endpoint", "", "Ingress URL of prometheus sidecar in current cluster. eg. prom-sidecar.cluster-xy")

    flag.Parse()

    opts := grpc.WithInsecure()
    if *tls {
        if *caFile == "" {
            *caFile = cert.Path("x509/ca.crt")
        }
        creds, err := credentials.NewClientTLSFromFile(*caFile, "")
        if err != nil {
            log.Fatal("Error while loading CA trust certificate: ", err)
            return
        }
        opts = grpc.WithTransportCredentials(creds)
    }

    if len(*clusterid) < 2 {
        log.Fatal("Argument Error: clusterid, eg. -clusterid c10ead61-ff2a-4af4-8f41-d2c44c745de7")
    }
    if len(*appgroupid) < 2 {
        log.Fatal("Argument Error: appgroupid, eg. -appgroupid abbead61-ff2a-4af4-8f41-d2c44c745de7")
    }

    addr := fmt.Sprintf("%s:%d", *tks_info_host, *tks_info_port)
    cc, err := grpc.Dial(addr, opts)
    if err != nil {
        log.Fatal("could not connect: ", err)
    }
    defer cc.Close()

    // Update AppGroupStatus to tks-info
    doUpdateAppGroupStatus(cc, *appgroupid, pb.AppGroupStatus_APP_GROUP_RUNNING)

    // Register endpoint to tks-info
    doUpdateAppEndpoint(cc, *appgroupid, *curEndpoint, pb.AppType_PROMETHEUS)

    // Update all prometheus endpoints in other clusters' site-yaml
    updateAllSiteYamls(cc, *clusterid, *curEndpoint, pb.AppType_PROMETHEUS)
}

func doUpdateAppGroupStatus(cc *grpc.ClientConn, appgroupid string, status pb.AppGroupStatus) {
    c := pb.NewAppInfoServiceClient(cc)

    req := &pb.UpdateAppGroupStatusRequest{
        AppGroupId: appgroupid,
        Status: status,
    }

    res, err := c.UpdateAppGroupStatus(context.Background(), req)
    if err != nil {
        log.Fatal("error while calling UpdateAppGroupStatus RPC::::: ", err)
    }
    log.Info("Response from UpdateAppGroupStatus: ", res.GetCode())
}

func doUpdateAppEndpoint(cc *grpc.ClientConn, appgroupid string, url string, appType pb.AppType) {
    c := pb.NewAppInfoServiceClient(cc)

    req := &pb.UpdateAppRequest{
        AppGroupId: appgroupid,
        AppType:    appType,
        Endpoint:   url,
        Metadata:   "{}",
    }

    res, err := c.UpdateApp(context.Background(), req)
    if err != nil {
        log.Fatal("error while calling UpdateApp RPC::::: ", err)
    }
    log.Info("Response from UpdateApp: ", res.GetCode())
}

func updateAllSiteYamls(cc *grpc.ClientConn, curClusterId string, curEndpoint string, appType pb.AppType) {
    appCl := pb.NewAppInfoServiceClient(cc)
    clusterCl := pb.NewClusterInfoServiceClient(cc)

    var endpointList []string
    //var endpointMap = map[string]string{}
    //endpointMap := make(map[string]string)

    /* Get cluster info to identify contract and csp id */
    req := &pb.GetClusterRequest{
        ClusterId: curClusterId,
    }

    res, err := clusterCl.GetCluster(context.Background(), req)
    if err != nil {
        log.Info("[test] res.Error from getCluster RPC::::: ", res.GetError())
        log.Fatal("[test] error while calling getCluster RPC::::: ", err)
    }
    log.Info("Response from getCluster: ", res.GetCluster())

    contractId := res.Cluster.ContractId
    cspId := res.Cluster.CspId
    curClusterName := res.Cluster.Name

    /* Get all other clusters within same contract and csp */
    gcReq := &pb.GetClustersRequest{
        CspId: cspId,
        ContractId: contractId,
    }

    gcRes, err := clusterCl.GetClusters(context.Background(), gcReq)
    if err != nil {
        log.Fatal("error while calling getClusters RPC::::: ", err)
    }
    log.Info("Response from getClusters: ", gcRes.GetClusters())

    /* For each clusters */
    for _, cluster := range gcRes.Clusters {
        clusterId := cluster.Id
        if clusterId != curClusterId {
            clusterName := cluster.Name
        log.Info("Processing cluster: ", clusterName)
            /* Get LMA appGroup in this cluster */
            req := &pb.IDRequest{
                Id: clusterId,
            }

            res, err := appCl.GetAppGroupsByClusterID(context.Background(), req)
            if err != nil {
                log.Fatal("error while calling getAppGroupsByClusterID RPC::::: ", err)
            }
            log.Info("Response from getAppGroupsByClusterID: ", res.GetAppGroups())

            /* For each appGroup in the cluster */
            for _, appGroup := range res.AppGroups {
                if appGroup.Type == pb.AppGroupType_LMA {
                    req := &pb.GetAppsRequest{
                        AppGroupId: appGroup.AppGroupId,
                        Type: pb.AppType_PROMETHEUS,
                    }
                    // Get promethus application
                    res, err := appCl.GetApps(context.Background(), req)
                    if err != nil {
                            log.Fatal("error while calling getApps RPC::::: ", err)
                    }
                    log.Info("Response from getApps: ", res.GetApps())

                    fmt.Printf("Retrieved prometheus app object: %+v", res.Apps[0])
                    promEndpoint := res.Apps[0].Endpoint
                    log.Info("prometheus endpoint: ", promEndpoint)
                    endpointList = append(endpointList, promEndpoint)
                }
            } // End of AppGroup iteration loop

            /* Add current cluster endpoint to this cluster's site-yaml */
            updateEndpointToSiteYaml(clusterName, curEndpoint)
        }

    } // End of cluster iteration loop

    /* Add other clusters' endpoints to the current cluster's site-yaml */
    endpointListStr := strings.Join(endpointList, " ")
    log.Info("endpointList as string: ", endpointListStr)
    updateMultipleEndpointsToSiteYaml(curClusterName, endpointListStr)
}

func updateEndpointToSiteYaml(clusterName string, endpoint string) {
    c := exec.Command("./updateEndpoint.py", "--action", "add", "--cluster_name", clusterName, "--endpoint", endpoint)

    var out []byte
    var err error
    if out, err = c.Output(); err != nil {
        log.Info(string(out))
        log.Fatal("Error while updating endpoint to site-yaml: ", err)
    } else {
        log.Info(string(out))
    }
}

func updateMultipleEndpointsToSiteYaml(curClusterName string, endpointList string) {
    c := exec.Command("./updateMultipleEndpoints.py", curClusterName, endpointList)

    var out []byte
    var err error
    if out, err = c.Output(); err != nil {
        log.Fatal("Error while updating multiple endpoints to site-yaml: ", err)
    }
    log.Info(string(out))
}
