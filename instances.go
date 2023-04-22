package main

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type InstanceAZ struct {
	InstanceId       string
	AvailabilityZone string
	Region           string
}

type EC2Client struct {
	region string
}

func (lclEc2 EC2Client) init() ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(lclEc2.region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return *ec2.NewFromConfig(cfg)
}

type TagToFilter struct {
	name     string
	rewrites map[string]string
}

func (ttf TagToFilter) rewriteTag() string {
	tag, ok := ttf.rewrites[ttf.name]
	if ok {
		return tag
	} else {
		return ttf.name
	}
}

func (ttf TagToFilter) getFilter() (types.Filter, error) {
	filters := AwsFilters{TagName: "tag:Name", PrivateIpFilter: "network-interface.private-dns-name"}
	var filter types.Filter
	var badFormat error
	reg, err := regexp.Compile("^dev*|^prod[!w]*|^staging*|^issuer-portal|^banker-portal")
	if err != nil {
		return types.Filter{}, err
	}
	west, err := regexp.Compile("^prodwest*")
	if err != nil {
		return types.Filter{}, err
	}
	ip, err := regexp.Compile("^ip-*")
	if err != nil {
		return types.Filter{}, err
	}

	if reg.MatchString(ttf.name) {
		filter = types.Filter{Name: &filters.TagName, Values: []string{ttf.rewriteTag()}}
	} else if west.MatchString(ttf.name) {
		filter = types.Filter{Name: &filters.TagName, Values: []string{ttf.rewriteTag()}}
	} else if ip.MatchString(ttf.name) {
		filter = types.Filter{Name: &filters.PrivateIpFilter, Values: []string{ttf.name}}
	} else {
		badFormat = errors.New("HostNotFound: " + ttf.name)
		return types.Filter{}, badFormat
	}

	return filter, nil

}

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func getInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func instanceConfig(server string, region string) (TargetConfig, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Println("Target Config failed")
		return TargetConfig{}, err
	}
	inAZ, err := getInstanceAZ(server, region)
	if err != nil {
		log.Println("GetInstanceAZ failed", err)
		return TargetConfig{}, err
	}

	return TargetConfig{Target: inAZ.InstanceId, Config: cfg}, nil

}

func getInstanceAZ(name string, region string) (InstanceAZ, error) {
	rewrites := map[string]string{"prodsalt01": "prodmonitor", "stagingsalt01": "stagingmonitor", "proddrone": "proddrone-server"}
	client := client(region)
	ttf := TagToFilter{name: name, rewrites: rewrites}
	filter, err := ttf.getFilter()
	if err != nil {
		return InstanceAZ{}, err
	}
	input := ec2.DescribeInstancesInput{Filters: []types.Filter{filter}}

	instance, _ := getInstances(context.TODO(), &client, &input)

	this := instance.Reservations[0].Instances[0]

	log.Printf("InstanceId: %s, AZ %s, Region: %s", *this.InstanceId, *this.Placement.AvailabilityZone, region)

	return InstanceAZ{InstanceId: *this.InstanceId, AvailabilityZone: *this.Placement.AvailabilityZone, Region: region}, nil
}

func client(region string) ec2.Client {
	return EC2Client{region: region}.init()
}
