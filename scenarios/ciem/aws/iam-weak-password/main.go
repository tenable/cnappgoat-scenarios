package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os/exec"
	"strings"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an IAM user without MFA enabled and with a login profile
		user, err := iam.NewUser(ctx, "CNAPPGoatWeakPasswordUser", &iam.UserArgs{
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})

		if err != nil {
			return err
		}
		// Create an IAM login profile for the user with the password "password"
		_, err = iam.NewUserLoginProfile(ctx, "loginProfile", &iam.UserLoginProfileArgs{
			User:           user.Name,
			PasswordLength: pulumi.Int(20),
		})
		// use the CLI to change the password to "password"
		// We have to use the CLI because the Pulumi SDK does not support updating the password for security reasons
		// aws iam update-login-profile --user-name CNAPPGoatWeakPasswordUser --password password
		_ = user.Name.ApplyT(func(name string) (string, error) {
			cmd := exec.Command("aws", "iam", "update-login-profile", "--user-name", name, "--password", "password")
			output, err := cmd.Output()
			if err != nil {
				return "", err
			}
			return strings.TrimSpace(string(output)), nil
		}).(pulumi.StringOutput)

		ctx.Export("userName", user.Name)
		ctx.Export("userArn", user.Arn)
		return nil
	})
}
