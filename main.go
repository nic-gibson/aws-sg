package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/pflag"
)

var regionName = ""
var ipURL = "https://ifconfig.me"
var sgId = ""
var sgRuleId = ""
var sgRuleDescription = "set by aws-g"
var fromPort int32 = -1
var toPort int32 = -1
var ipProtocol = "-1"

// This program looks up our public IP via the ifconfig.me service, updates the JSON
// document describing the security group I use on AWS to include the public
// IP address and then executes the AWS CLI to update the SG

func main() {
	pflag.StringVarP(&regionName, "region", "r", "us-east-1", "The AWS region containing the security group.")
	pflag.StringVarP(&sgId, "group", "g", "", "The AWS security group to modify.")
	pflag.StringVarP(&sgRuleId, "rule", "l", "", "The AWS security group rule to modify.")
	pflag.Parse()

	if regionName == "" || sgId == "" || sgRuleId == "" {
		fmt.Fprintln(os.Stderr, "Usage:")
		pflag.PrintDefaults()
		os.Exit(2)
	}
	client := getEC2Service(context.Background())
	_, err := updateSecurityRule(context.Background(), client)

	if err != nil {
		log.Fatalf("Unable to set IP address in rules, %v", err)
	}
}

func getMyIpAddress() (string, error) {

	response, error := http.Get(ipURL)

	if error != nil {
		return "", error
	} else {
		defer response.Body.Close()
		body, error := io.ReadAll(response.Body)
		if error != nil {
			return "", error
		} else {
			return string(body[:]), nil
		}
	}
}

func getAWSConfig(ctx context.Context) aws.Config {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(regionName),
	)

	if err != nil {
		log.Fatalf("Unable to load AWS SDK config, %v", err)
	}

	return cfg
}

func getEC2Service(ctx context.Context) *ec2.Client {
	config := getAWSConfig(ctx)
	return ec2.NewFromConfig(config)
}

func updateSecurityRule(ctx context.Context, client *ec2.Client) (*ec2.ModifySecurityGroupRulesOutput, error) {
	cidr, err := getMyIpAddress()
	cidr += "/32"

	if err != nil {
		log.Fatalf("Unable to load AWS SDK config, %v", err)
	}

	params := ec2.ModifySecurityGroupRulesInput{
		GroupId: &sgId,
		SecurityGroupRules: []types.SecurityGroupRuleUpdate{
			{
				SecurityGroupRule: &types.SecurityGroupRuleRequest{
					CidrIpv4:    &cidr,
					Description: &sgRuleDescription,
					FromPort:    &toPort,
					ToPort:      &fromPort,
					IpProtocol:  &ipProtocol,
				},
				SecurityGroupRuleId: &sgRuleId,
			},
		},
	}

	return client.ModifySecurityGroupRules(ctx, &params)
}
