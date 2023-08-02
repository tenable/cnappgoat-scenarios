package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		awsProvider, err := aws.NewProvider(ctx, "awsProvider", &aws.ProviderArgs{})
		
		if err != nil {
			return err
		}

		// Create a new VPC
		vpc, err := ec2.NewVpc(ctx, "custom-vpc", &ec2.VpcArgs{
			CidrBlock: pulumi.String("10.0.0.0/16"),
		}, pulumi.Provider(awsProvider))
		if err != nil {
			return err
		}

		// Create a new subnet in the VPC
		subnet, err := ec2.NewSubnet(ctx, "custom-subnet", &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: pulumi.String("10.0.1.0/24"),

		}, pulumi.Provider(awsProvider))
		if err != nil {
			return err
		}

		// Create a new security group that allows SSH and HTTP access
		group, err := ec2.NewSecurityGroup(ctx, "web-secgrp", &ec2.SecurityGroupArgs{
			VpcId:       vpc.ID(),
			Description: pulumi.String("Enable SSH and HTTP access"),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}, pulumi.Provider(awsProvider))
		if err != nil {
			return err
		}

		// Get the Ubuntu 21.10 AMI
		mostRecent := true
		amiResult, err := aws.GetAmi(ctx, &aws.GetAmiArgs{
			Owners:     []string{"amazon"},
			MostRecent: &mostRecent,
			Filters:    []aws.GetAmiFilter{{Name: "name", Values: []string{"ubuntu*21.10*"}},
		 								   {Name: "architecture", Values: []string{"x86_64"}}},
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewInstance(ctx, "CnappgoatCWPPEndOfLifeEC2", &ec2.InstanceArgs{
			Ami:                      pulumi.String(amiResult.Id), 
			InstanceType:             pulumi.String("t2.micro"),
			VpcSecurityGroupIds:      pulumi.StringArray{group.ID()},
			SubnetId:                 subnet.ID(), // associate the instance with the subnet
			AssociatePublicIpAddress: pulumi.Bool(true),
			Tags:                     pulumi.StringMap{"Name": pulumi.String("CnappgoatCWPPEndOfLifeEC2")},
		}, pulumi.Provider(awsProvider))
		if err != nil {
			return err
		}

		return nil
	})
}
