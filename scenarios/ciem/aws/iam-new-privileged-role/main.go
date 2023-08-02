package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get the AWS Account ID
		callerIdentity, err := aws.GetCallerIdentity(ctx, nil, nil)
		// Create an IAM Role with admin permissions
		assumeRolePolicyJSON := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{	
					"Effect": "Allow",
					"Principal": {
						"AWS": "arn:aws:iam::%s:root"
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`, callerIdentity.AccountId)
		role, err := iam.NewRole(ctx, "CNAPPGoatAdminRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(assumeRolePolicyJSON),
			ManagedPolicyArns: pulumi.StringArray{
				pulumi.String("arn:aws:iam::aws:policy/AdministratorAccess"),
			},
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("roleArn", role.Arn)
		return nil
	})
}
