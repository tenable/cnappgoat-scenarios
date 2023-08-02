package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an IAM user without MFA enabled and with a login profile
		user, err := iam.NewUser(ctx, "user", &iam.UserArgs{
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})

		if err != nil {
			return err
		}
		// Create an IAM login profile for the user
		_, err = iam.NewUserLoginProfile(ctx, "loginProfile", &iam.UserLoginProfileArgs{
			User:           user.Name,
			PasswordLength: pulumi.Int(20),
		})
		ctx.Export("userName", user.Name)
		ctx.Export("userArn", user.Arn)
		return nil
	})
}
