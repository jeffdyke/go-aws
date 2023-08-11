package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/mmmorris1975/ssm-session-client/ssmclient"
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
	log.Fatal(ssmclient.PortForwardingSession(tgtCfg.Config, &in))
}
