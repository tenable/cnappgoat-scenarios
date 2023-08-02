package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a new IAM user with admin privileges
		user, err := iam.NewUser(ctx, "CNAPPGoatNewPrivilegedUser", &iam.UserArgs{
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		// attach the admin policy to the user
		_, err = iam.NewUserPolicyAttachment(ctx, "CNAPPGoatNewPrivilegedUserAdminPolicy", &iam.UserPolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AdministratorAccess"),
			User:      user.Name,
		})
		if err != nil {
			return err
		}

		// create access keys for the user
		_, err = iam.NewAccessKey(ctx, "CNAPPGoatNewPrivilegedUserAccessKey", &iam.AccessKeyArgs{
			User: user.Name,
		})
		if err != nil {
			return err
		}
		return nil
	})
}
