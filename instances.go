package main

import (
	"context"
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

func (blEc2 EC2Client) init() ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(blEc2.region))
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

func (ttf TagToFilter) getFilter() types.Filter {
	filters := AwsFilters{TagName: "tag:Name", PrivateIpFilter: "network-interface.private-dns-name"}
	var filter types.Filter
	reg, err := regexp.Compile("^dev*|^prod[!w]*|^staging*|^issuer-portal|^banker-portal")
	if err != nil {
		log.Fatal("Could not parse regex", err)
	}
	west, err := regexp.Compile("^prodwest*")
	if err != nil {
		log.Fatal("Could not parse regex", err)
	}
	ip, err := regexp.Compile("^ip-*")
	if err != nil {
		log.Fatal("Could not parse regex", err)
	}

	if reg.MatchString(ttf.name) {
		filter = types.Filter{Name: &filters.TagName, Values: []string{ttf.rewriteTag()}}
	} else if west.MatchString(ttf.name) {
		filter = types.Filter{Name: &filters.TagName, Values: []string{ttf.rewriteTag()}}
	} else if ip.MatchString(ttf.name) {
		filter = types.Filter{Name: &filters.PrivateIpFilter, Values: []string{ttf.name}}
	}
	return filter
}

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func getInstanceAZ(name string, region string) InstanceAZ {
	rewrites := map[string]string{"prodsalt01": "prodmonitor", "stagingsalt01": "stagingmonitor", "proddrone": "proddrone-server"}
	client := client(region)
	ttf := TagToFilter{name: name, rewrites: rewrites}
	filter := ttf.getFilter()
	input := ec2.DescribeInstancesInput{Filters: []types.Filter{filter}}
	instance, err := GetInstances(context.TODO(), &client, &input)
	if err != nil {
		log.Fatal("Can't find instance from "+name+"\r", err)
	}
	this := instance.Reservations[0].Instances[0]

	return InstanceAZ{InstanceId: *this.InstanceId, AvailabilityZone: *this.Placement.AvailabilityZone, Region: region}
}

func client(region string) ec2.Client {
	return EC2Client{region: "us-east-1"}.init()
}

// func main() {
// 	// client := client("us-east-1")
// 	// filter := GetFilter("develop")

// 	//input := ec2.DescribeInstancesInput{Filters: []types.Filter{filter}}
// 	// instance, err := GetInstances(context.TODO(), &client, &input)
// 	iAZ := getInstanceAZ("develop", "us-east-1")

// 	log.Printf("%+v\n", iAZ)
// }
