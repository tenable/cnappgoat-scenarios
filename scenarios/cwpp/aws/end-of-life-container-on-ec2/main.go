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

		// Create a new security group that allows SSH, HTTP, and Docker access
		group, err := ec2.NewSecurityGroup(ctx, "web-secgrp", &ec2.SecurityGroupArgs{
			Description: pulumi.String("Enable SSH and HTTP access"),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(80),
					ToPort:     pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(2375), // Docker daemon port
					ToPort:     pulumi.Int(2375),
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

		// Define the user data to install Docker and run a container
		userData := `#!/bin/bash
		sudo yum update -y
		sudo yum install -y docker
		sudo service docker start
		sudo usermod -a -G docker ec2-user
		sudo docker run --name end_of_life_container -d -p 80:80 public.ecr.aws/i3j2g7c0/cnappgoat-images:end_of_life_ubuntu2110_image
		`

		// Get the latest AMI
		mostRecent := true
		amiResult, err := aws.GetAmi(ctx, &aws.GetAmiArgs{
			Owners:     []string{"amazon"},
			MostRecent: &mostRecent,
			Filters:    []aws.GetAmiFilter{{Name: "name", Values: []string{"amzn-ami-hvm-*-x86_64-ebs"}}},
		})
		if err != nil {
			return err
		}

		// Create a new EC2 instance
		_, err = ec2.NewInstance(ctx, "CnappgoatCWPPEndOfLifeContainerOnEC2", &ec2.InstanceArgs{
			Ami:                      pulumi.String(amiResult.Id), 
			InstanceType:             pulumi.String("t3.micro"),
			VpcSecurityGroupIds:      pulumi.StringArray{group.ID()},
			UserData:                 pulumi.String(userData),
			AssociatePublicIpAddress: pulumi.Bool(true),
			Tags:                     pulumi.StringMap{"Name": pulumi.String("CnappgoatCWPPEndOfLifeContainerOnEC2")},
		}, pulumi.Provider(awsProvider))
		if err != nil {
			return err
		}



		return nil
	})
}