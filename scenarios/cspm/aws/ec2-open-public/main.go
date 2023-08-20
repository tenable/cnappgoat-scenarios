package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new VPC
		vpc, err := ec2.NewVpc(ctx, "CNAPPGoat-ec2-open-public-vpc", &ec2.VpcArgs{
			CidrBlock: pulumi.String("10.0.0.0/16"),
		})
		if err != nil {
			return err
		}

		// Create a new subnet
		subnet, err := ec2.NewSubnet(ctx, "CNAPPGoat-ec2-open-public-subnet", &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: pulumi.String("10.0.1.0/24"),
		})
		if err != nil {
			return err
		}

		// Create a new security group
		securityGroup, err := ec2.NewSecurityGroup(ctx, "CNAPPGoat-ec2-open-public-securitygroup", &ec2.SecurityGroupArgs{
			VpcId: vpc.ID(),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("tcp"),
					FromPort: pulumi.Int(80),
					ToPort:   pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Get the AMI
		mostRecent := true
		amiResult, err := aws.GetAmi(ctx, &aws.GetAmiArgs{
			Owners:     []string{"amazon"},
			MostRecent: &mostRecent,
			Filters:    []aws.GetAmiFilter{{Name: "name", Values: []string{"amzn2-ami-hvm-2.0.*-x86_64-ebs"}}},
		})
		if err != nil {
			return err
		}

		// Create a new EC2 instance
		ec2, err := ec2.NewInstance(ctx, "CNAPPGoat-ec2-open-public-instance", &ec2.InstanceArgs{
			InstanceType:             pulumi.String("t2.micro"),
			AssociatePublicIpAddress: pulumi.BoolPtr(true),
			VpcSecurityGroupIds: pulumi.StringArray{
				securityGroup.ID(),
			},
			SubnetId: subnet.ID(),
			Ami:      pulumi.String(amiResult.Id),
			UserData: pulumi.String(`#!/bin/bash
                echo "Hello, World!" > index.html
                nohup python -m SimpleHTTPServer 80 &`),
			Tags: pulumi.StringMap{
				"Name":      pulumi.String("CNAPPGoat-ec2-open-public-instance"),
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("CNAPPGoat-ec2-open-public-vpc", vpc.Arn)
		ctx.Export("CNAPPGoat-ec2-open-public-subnet", subnet.Arn)
		ctx.Export("CNAPPGoat-ec2-open-public-securitygroup", securityGroup.Arn)
		ctx.Export("CNAPPGoat-ec2-open-public-instance", ec2.Arn)
		return nil
	})
}
