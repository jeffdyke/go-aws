package main

import (
	"context"
	"log"

	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type AwsConfig struct {
	conf aws.Config
}

type InstanceAZ struct {
	InstanceId       string
	AvailabilityZone string
}

type AwsFilters struct {
	TagName         string
	PrivateIpFilter string
}

func getFilter(server string) types.Filter {
	filters := AwsFilters{TagName: "tag:Name", PrivateIpFilter: "network-interface.private-dns-name"}
	var filter types.Filter
	reg, err := regexp.Compile("^dev*|^prod[!w]*|^staging*|^issuer-portal|^banker-portal")
	west, err := regexp.Compile("^prodwest*")
	ip, err := regexp.Compile("^ip-*")
	if err != nil {
		log.Fatal("Could not parse regex", err)
	}

	if reg.MatchString(server) {
		filter = types.Filter{Name: &filters.TagName, Values: []string{server}}
	} else if west.MatchString(server) {
		filter = types.Filter{Name: &filters.TagName, Values: []string{server}}
	} else if ip.MatchString(server) {
		filter = types.Filter{Name: &filters.PrivateIpFilter, Values: []string{server}}
	}
	return filter
}
func getInstanceAZ(server string, region string) InstanceAZ {
	ec2Filter := getFilter(server)
	client := EC2Client{region: region}.init()
	input := ec2.DescribeInstancesInput{Filters: []types.Filter{ec2Filter}}
	result, err := GetInstances(context.TODO(), &client, &input)
	if err != nil {
		log.Fatal("Failed to find instance", err)
	}
	instance := result.Reservations[0].Instances[0]

	return InstanceAZ{InstanceId: *instance.InstanceId, AvailabilityZone: *instance.Placement.AvailabilityZone}

}

func main() {
	iaz := getInstanceAZ("develop", "us-east-1")
	log.Printf("%+v\n", iaz)
	// for _, r in range result.Reservations {

	// }
	// log.Printf("%+v\n", result)
}
