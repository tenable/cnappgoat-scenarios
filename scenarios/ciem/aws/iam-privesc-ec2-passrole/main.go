package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		assumeRolePolicy := `{
"Version": "2012-10-17",
"Statement": [
	{
		"Action": "sts:AssumeRole",
		"Principal": {
			"Service": "ec2.amazonaws.com"
		},
		"Effect": "Allow",
		"Sid": ""
	}
]
}`

		rolePolicy := `{
"Version": "2012-10-17",
"Statement": [
	{
		"Effect": "Allow",
		"Action": "ec2:RunInstances",
		"Resource": "*"
	},
	{
		"Effect": "Allow",
		"Action": "iam:PassRole",
		"Resource": "*"
	}
]
}`

		// Create the IAM role
		role, err := iam.NewRole(ctx, "role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(assumeRolePolicy),
		})
		if err != nil {
			return err
		}

		// Attach the policy to the IAM role
		_, err = iam.NewRolePolicy(ctx, "rolePolicy", &iam.RolePolicyArgs{
			Role:   role.ID(),
			Policy: pulumi.String(rolePolicy),
		})
		if err != nil {
			return err
		}

		ctx.Export("roleName", role.Name)
		ctx.Export("roleArn", role.Arn)
		return nil
	})
}
