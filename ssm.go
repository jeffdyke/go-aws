package main

import (
	"log"

	"github.com/mmmorris1975/ssm-session-client/ssmclient"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func portForward(tgtCfg LocalPortForward) {
	in := ssmclient.PortForwardingInput{
		Target:     tgtCfg.Target,
		RemotePort: int(tgtCfg.RemotePort),
		LocalPort:  int(tgtCfg.LocalPort),
	}
	log.Fatal(ssmclient.PortPluginSession(tgtCfg.Config, &in))
}

// func main() {
// 	tc := instanceConfig("proddb01", "us-east-1")
// 	lpf := LocalPortForward{TargetConfig: tc, RemotePort: 3306, LocalPort: 1515}
// 	mysql(lpf)
// }
