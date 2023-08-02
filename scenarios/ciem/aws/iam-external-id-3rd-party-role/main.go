package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an IAM role that allows a 3rd party to assume it, make sure the role has a randomized postfix
		// so that it is unique
		role, err := iam.NewRole(ctx, "CnappGoatExternalIdIamRoleRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "",
						"Effect": "Allow",		
						"Principal": {		
							"AWS": "arn:aws:iam::152659312504:root"	
						},		
						"Action": "sts:AssumeRole"
					}
				]
			}`), // This account ID is the account ID of the 3rd party for Slack EKM https://slackhq.com/dotcom/dotcom/wp-content/uploads/sites/6/2019/08/Slack-EKM-Implementation-Guide-1.pdf
			ManagedPolicyArns: pulumi.StringArray{
				pulumi.String("arn:aws:iam::aws:policy/AWSDenyAll"),
			},
			Description: pulumi.String("This is a vulnerable role that allows a 3rd party to assume it without an external ID"),
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("roleName", role.Name)
		return nil
	})
}
