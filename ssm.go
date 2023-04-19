package main

import (
	"context"
	"log"

	"github.com/mmmorris1975/ssm-session-client/ssmclient"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type TargetConfig struct {
	Target string
	Config aws.Config
}
type LocalPortForward struct {
	TargetConfig
	RemotePort int
	LocalPort  int
}

func buildConfig(server string, region string) TargetConfig {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	inAZ := getInstanceAZ(server, region)
	if err != nil {
		log.Println()
		log.Fatal("Failed to resolve target "+server, err)
	}
	return TargetConfig{Target: inAZ.InstanceId, Config: cfg}

}

func mysql(tgtCfg LocalPortForward) {
	in := ssmclient.PortForwardingInput{
		Target:     tgtCfg.Target,
		RemotePort: int(tgtCfg.RemotePort),
		LocalPort:  int(tgtCfg.LocalPort),
	}
	log.Fatal(ssmclient.PortPluginSession(tgtCfg.Config, &in))
}

// func main() {
// 	tc := buildConfig("proddb01", "us-east-1")
// 	lpf := LocalPortForward{TargetConfig: tc, RemotePort: 3306, LocalPort: 1515}
// 	mysql(lpf)
// }
