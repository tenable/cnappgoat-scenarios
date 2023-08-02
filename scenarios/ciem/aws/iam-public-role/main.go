package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		publicRoleName := "CnappGoatPublicIamRoleRole"
		role, err := iam.NewRole(ctx, publicRoleName, &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",	
				"Statement": [
					{
						"Sid": "",			
						"Effect": "Allow",	
						"Principal": {			
							"AWS": "*"
						},
						"Action": "sts:AssumeRole"			
					}
				]
			}`),
			ManagedPolicyArns: pulumi.StringArray{
				pulumi.String("arn:aws:iam::aws:policy/AWSDenyAll"),
			},
			Description: pulumi.String("This is a vulnerable role that allows anyone to assume it"),
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})

		if err != nil {
			return err
		}

		ctx.Export("publicRoleName", role.Name)
		ctx.Export("publicRoleArn", role.Arn)
		return nil
	})
}
