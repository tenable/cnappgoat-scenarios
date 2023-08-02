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

		// Create a new security group that allows SSH and HTTP access
		group, err := ec2.NewSecurityGroup(ctx, "web-secgrp", &ec2.SecurityGroupArgs{
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

		// Create a new EC2 instance
		userData := `#!/bin/bash
wget https://cdn.amazonlinux.com/2/core/2.0/x86_64/6b0225ccc542f3834c95733dcf321ab9f1e77e6ca6817469771a8af7c49efe6c/../../../../../blobstore/3b226f60ce3c33d4b04ba594484cf3f96256cf5760ef39ed031a4c452127b6c9/thunderbird-91.8.0-1.amzn2.0.1.x86_64.rpm
sudo rpm -i thunderbird-91.8.0-1.amzn2.0.1.x86_64.rpm --nodeps
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

		_, err = ec2.NewInstance(ctx, "CnappgoatCWPPVulnerableEC2Thunderbird", &ec2.InstanceArgs{
			Ami:                      pulumi.String(amiResult.Id), 
			InstanceType:             pulumi.String("t3.micro"),
			VpcSecurityGroupIds:      pulumi.StringArray{group.ID()},
			UserData:                 pulumi.String(userData),
			AssociatePublicIpAddress: pulumi.Bool(true),
			Tags:                     pulumi.StringMap{"Name": pulumi.String("cwppVulnerableVMThunderbird")},
		}, pulumi.Provider(awsProvider))
		if err != nil {
			return err
		}

		return nil
	})
}
