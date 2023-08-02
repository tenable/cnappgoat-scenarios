package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
	"time"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create the overprivileged IAM user
		overprivilegedUser, err := iam.NewUser(ctx, "CNAPPGoatOverprivilegedUser", &iam.UserArgs{
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		// Create the overprivileged IAM group
		overprivilegedGroup, err := iam.NewGroup(ctx, "CNAPPGoatOverprivilegedGroup", nil)
		if err != nil {
			return err
		}

		// Attach the S3 full access policy to the group
		_, err = iam.NewGroupPolicyAttachment(ctx, "CNAPPGoatOverprivilegedGroupS3FullAccess", &iam.GroupPolicyAttachmentArgs{
			Group:     overprivilegedGroup.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonS3FullAccess"),
		})
		if err != nil {
			return err
		}

		// Add the user to the group
		_, err = iam.NewGroupMembership(ctx, "CNAPPGoatOverprivilegedUserMembership", &iam.GroupMembershipArgs{
			Group: overprivilegedGroup.Name,
			Users: pulumi.StringArray{overprivilegedUser.Name},
		})
		if err != nil {
			return err
		}
		// Create access keys for the user
		keys, err := iam.NewAccessKey(ctx, "CNAPPGoatOverprivilegedUserAccessKey", &iam.AccessKeyArgs{
			User: overprivilegedUser.Name,
		})
		if err != nil {
			return err
		}

		// store the access key in a secret
		secret, err := secretsmanager.NewSecret(ctx, "CNAPPGoatOverprivilegedUserAccessKeySecret", &secretsmanager.SecretArgs{})
		if err != nil {
			return err
		}
		_, err = secretsmanager.NewSecretVersion(ctx, "CNAPPGoatOverprivilegedUserAccessKeyVersion", &secretsmanager.SecretVersionArgs{
			SecretId:     secret.ID(),
			SecretString: keys.Secret,
		})
		if err != nil {
			return err
		}

		// Create an IAM Role for the lambda function that allows it to access the secret
		assumeRolePolicyJSON := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",	
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`)
		role, err := iam.NewRole(ctx, "CNAPPGoatOverprivilegedUserLambdaRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(assumeRolePolicyJSON),
			ManagedPolicyArns: pulumi.StringArray{
				pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
			},
			InlinePolicies: iam.RoleInlinePolicyArray{
				iam.RoleInlinePolicyArgs{
					Name: pulumi.String("CNAPPGoatOverprivilegedUserLambdaPolicy"),
					Policy: pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Effect": "Allow",
							"Action": "secretsmanager:GetSecretValue",
							"Resource": "%s"
						}
					]
				}`, secret.Arn)}},
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}

		// Create an S3 bucket
		bucket, err := s3.NewBucket(ctx, "CNAPPGoatOverprivilegedUserBucket", &s3.BucketArgs{
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}

		// Create a lambda function that assumes the role and lists the bucket
		setupLambda, err := lambda.NewFunction(ctx, "CNAPPGoatOverprivilegedRoleFunction", &lambda.FunctionArgs{
			Role:    role.Arn,
			Handler: pulumi.String("index.handler"),
			Runtime: pulumi.String("nodejs14.x"),
			Code:    pulumi.NewFileArchive("./lambda"),
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}

		_, err = lambda.NewInvocation(ctx, fmt.Sprintf("setupInvocation%s", strconv.FormatInt(time.Now().Unix(), 10)), &lambda.InvocationArgs{
			FunctionName: setupLambda.Name,
			Input:        pulumi.Sprintf(`{"BucketName":"%s", "SecretId": "%s", "AccessKeyId": "%s"}`, bucket.ID(), secret.ID(), keys.ID()),
		})

		if err != nil {
			return err
		}

		// Exports
		ctx.Export("bucketName", bucket.ID())
		ctx.Export("secretId", secret.ID())
		ctx.Export("overprivilegedUserArn", overprivilegedUser.Arn)
		ctx.Export("FunctionArn", setupLambda.Arn)
		return nil
	})
}
