package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create an IAM Role
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
		role, err := iam.NewRole(ctx, "CNAPPGoatOverprivilegedRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(assumeRolePolicyJSON),
			ManagedPolicyArns: pulumi.StringArray{
				pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
				pulumi.String("arn:aws:iam::aws:policy/AmazonS3FullAccess"),
			},
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}

		// Create an S3 bucket
		bucket, err := s3.NewBucket(ctx, "CNAPPGoatOverprivilegedBucket", &s3.BucketArgs{
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

		_, err = lambda.NewInvocation(ctx, "setupInvocation", &lambda.InvocationArgs{
			FunctionName: setupLambda.Name,
			Input:        pulumi.Sprintf(`{"BucketName":"%s"}`, bucket.ID()),
		})
		if err != nil {
			return err
		}

		// Export the ARN of the role and the name of the bucket
		ctx.Export("roleArn", role.Arn)
		ctx.Export("bucketName", bucket.ID())
		ctx.Export("lambdaArn", setupLambda.Arn)
		return nil
	})
}
