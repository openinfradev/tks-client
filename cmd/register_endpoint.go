package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sktelecom/tks-contract/pkg/log"
	"github.com/sktelecom/tks-info/pkg/cert"
	pb "github.com/sktelecom/tks-proto/pbgo"
)

var (
	port   = flag.Int("port", 9111, "The gRPC server port")
	tls    = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile = flag.String("ca_file", "", "The TLS ca file")
)

func main() {
	if len(os.Args) < 2 {
		panic("Need parameters. (help with -h)")
	}
	log.Info("Preporcess for TACO-LMA with TKS-info")

	ip := flag.String("tks", "127.0.0.1", "An address of TKS-info")
	clusterid := flag.String("clusterid", "", "Cluster ID to apply. eg. 6abead61-ff2a-4af4-8f41-d2c44c745de7")
	appgroupid := flag.String("appgroupid", "", "Application ID of The endpoint. eg. 6abead61-ff2a-4af4-8f41-d2c44c745de7")
	clusterep := flag.String("clusterep", "", "Cluster url or ip. eg. cluster.taco.com or 192.168.123.45")
	epportlist := flag.String("epportlist", "",
		"The list of port per app like (application_type, port) eg. {\"1\":\"80\",\"2\":\"10232\"}")
	eplist := flag.String("eplist", "",
		"The list of endpoints like (application_type, endpoint) eg. {\"1\":\"192.168.5.55:80\",\"2\":\"192.168.5.55:10232\"}")

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

	addr := fmt.Sprintf("%s:%d", *ip, *port)
	cc, err := grpc.Dial(addr, opts)
	if err != nil {
		log.Fatal("could not connect: ", err)
	}
	defer cc.Close()

	// three types to define eps
	// - eplist: set of entry (type, url)
	// - clusterep and epportlist: cluster endpoint and set of entry (type, port)
	// - clusterep only: cluster endpoint and default endpoint using the node port. (currently prometheus)
	if len(*eplist) > 2 {
		var eps map[pb.AppType]string
		json.Unmarshal([]byte(*eplist), &eps)

		for sw, ep := range eps {
			doUpdateAppEndpoint(cc, *clusterid, *appgroupid, ep, sw)
		}
	}

	if len(*epportlist) > 2 && len(*clusterep) > 2 {
		var ports map[pb.AppType]int
		json.Unmarshal([]byte(*epportlist), &ports)

		for sw, port := range ports {
			url := fmt.Sprintf("%s:%d", clusterep, port)
			doUpdateAppEndpoint(cc, *clusterid, *appgroupid, url, sw)
		}
	} else if len(*clusterep) > 2 {
		// It's a logic to register default prometheus(sidecar) endpoint using nodeport
		url := fmt.Sprintf("%s:%d", *clusterep, 30007)
		doUpdateAppEndpoint(cc, *clusterid, *appgroupid, url, 2)
	}
	fmt.Println(len(*clusterep))
}

func doUpdateAppEndpoint(cc *grpc.ClientConn, clusterid, appgroupid, url string, appType pb.AppType) {
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
