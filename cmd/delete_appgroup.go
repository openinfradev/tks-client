package main

import (
    "context"
    "flag"
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

    // Delete AppGroup from tks-info
    doDeleteAppGroup(cc, *appgroupid)

    /* Update all prometheus endpoints in other clusters' site-yaml */
    updateAllSiteYamls(cc, *clusterid, *curEndpoint, pb.AppType_PROMETHEUS)
}

func doDeleteAppGroup(cc *grpc.ClientConn, appgroupid string) {
    c := pb.NewAppInfoServiceClient(cc)

    req := &pb.DeleteAppGroupRequest{
        AppGroupId: appgroupid,
    }

    res, err := c.DeleteAppGroup(context.Background(), req)
    if err != nil {
        log.Fatal("error while calling DeleteAppGroup RPC::::: ", err)
    }
    log.Info("Response from DeleteAppGroup: ", res.GetCode())
}

func updateAllSiteYamls(cc *grpc.ClientConn, curClusterId string, curEndpoint string, appType pb.AppType) {
    clusterCl := pb.NewClusterInfoServiceClient(cc)

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

            /* Delete current cluster endpoint from this cluster's site-yaml */
            deleteEndpointFromSiteYaml(clusterName, curEndpoint)
        }
    } // End of cluster iteration loop
}

func deleteEndpointFromSiteYaml(clusterName string, endpoint string) {
    c := exec.Command("./updateEndpoint.py", "--action", "delete", "--cluster_name", clusterName, "--endpoint", endpoint)

    var out []byte
    var err error
    if out, err = c.Output(); err != nil {
        log.Info(string(out))
        log.Fatal("Error while deleting endpoint from site-yaml: ", err)
    } else {
        log.Info(string(out))
    }
}
