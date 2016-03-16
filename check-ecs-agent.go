package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {
	cluster := flag.String("cluster", "", "The name of the ECS cluster this instance belongs to.")

	flag.Parse()

	metadataClient := ec2metadata.New(session.New())

	region, err := metadataClient.Region()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	instanceId, err := metadataClient.GetMetadata("instance-id")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ecsClient := ecs.New(session.New(), &aws.Config{Region: aws.String(region)})

	listParams := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(*cluster),
	}

	listResp, err := ecsClient.ListContainerInstances(listParams)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	describeParams := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(*cluster),
		ContainerInstances: listResp.ContainerInstanceArns,
	}

	describeResp, err := ecsClient.DescribeContainerInstances(describeParams)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, containerInstance := range describeResp.ContainerInstances {
		if *containerInstance.Ec2InstanceId == instanceId {
			if *containerInstance.AgentConnected == false {
				fmt.Printf("Amazon ECS Agent on %v is not connected.\n", *containerInstance.Ec2InstanceId)
				os.Exit(2)
			} else {
				fmt.Printf("Amazon ECS Agent on %v is connected.\n", *containerInstance.Ec2InstanceId)
				os.Exit(0)
			}
		}
	}
}
