package ec2data

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const region = "us-east-2"

type Ec2Data struct {
	Name       string
	Region     string
	Az         string
	PrivateIp  string
	Status     string
	InstanceId string
}

func GetEc2Data(owner string) []Ec2Data {
	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Owner"),
				Values: []*string{
					aws.String(owner),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil
	}
	return CreateEc2Objs(result)
}

func CreateEc2Objs(result *ec2.DescribeInstancesOutput) []Ec2Data {
	ec2s := []Ec2Data{}
	if len(result.Reservations) > 0 {
		for _, res := range result.Reservations {
			for _, ins := range res.Instances {
				ec2 := Ec2Data{}
				ec2.Region = region
				ec2.Status = *(ins.State.Name)
				fmt.Println(ins)
				if ins.PrivateIpAddress != nil {
					ec2.PrivateIp = *(ins.PrivateIpAddress)
				}
				ec2.Az = *(ins.Placement.AvailabilityZone)
				ec2.InstanceId = *(ins.InstanceId)
				var name string
				for _, tag := range ins.Tags {
					key := *(tag.Key)
					value := *(tag.Value)
					if key == "Name" {
						name = value
					}
				}
				ec2.Name = name
				ec2s = append(ec2s, ec2)
			}
		}
	}
	return ec2s
}

func StartInstance(instanceId string) int {
	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err := svc.StartInstances(input)
	if err != nil {
		fmt.Println(err.Error())
		return 500
	}
	return 200
}

func StopInstance(instanceId string) int {
	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err := svc.StopInstances(input)
	if err != nil {
		fmt.Println(err.Error())
		return 500
	}
	return 200
}

func TerminateInstance(instanceId string) int {
	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err := svc.TerminateInstances(input)
	if err != nil {
		fmt.Println(err.Error())
		return 500
	}
	return 200
}
