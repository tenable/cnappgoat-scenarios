package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		idp, err := iam.NewOpenIdConnectProvider(ctx, "CnappGoatGithubOidcProvider", &iam.OpenIdConnectProviderArgs{
			Url: pulumi.String("https://token.actions.githubusercontent.com"),
			ClientIdLists: pulumi.StringArray{
				pulumi.String("sts.amazonaws.com"),
			},
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
			ThumbprintLists: pulumi.StringArray{
				pulumi.String("6938fd4d98bab03faadb97b34396831e3780aea1"),
				pulumi.String("1c58a3a8518e8759bf075b76b750d4f2df264fcd"),
			},
		})
		if err != nil {
			return err
		}

		role, err := iam.NewRole(ctx, "CnappGoatOverlyPermissionGithubRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "",
						"Effect": "Allow",		
						"Principal": {		
							"Federated": "%s"	
						},		
						"Action": "sts:AssumeRoleWithWebIdentity"
					}
				]
			}`, idp.Arn),
			ManagedPolicyArns: pulumi.StringArray{
				pulumi.String("arn:aws:iam::aws:policy/AWSDenyAll"),
			},
			Description: pulumi.String("This is a vulnerable role that allows any repo on Github.com to assume it"),
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
