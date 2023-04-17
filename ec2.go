package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client struct {
	region string
}

func (blEc2 EC2Client) init() ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(blEc2.region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return *ec2.NewFromConfig(cfg)
}

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}
