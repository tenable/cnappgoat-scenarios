package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/codebuild"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		callerIdentity, err := aws.GetCallerIdentity(ctx, nil)
		if err != nil {
			return err
		}
		accountID := callerIdentity.AccountId
		// Create IAM Role for CodeBuild
		codebuildRole, err := iam.NewRole(ctx, "codebuildRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"Service": "codebuild.amazonaws.com"
						},
						"Action": "sts:AssumeRole"
					}
				]
			}`),
		})
		if err != nil {
			return err
		}

		// Define your policy
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Resource": [
						"arn:aws:logs:eu-central-1:%s:log-group:/aws/codebuild/cnappgoat-codebuild",
						"arn:aws:logs:eu-central-1:%s:log-group:/aws/codebuild/cnappgoat-codebuild:*"
					],
					"Action": [
						"logs:CreateLogGroup",
						"logs:CreateLogStream",
						"logs:PutLogEvents"
					]
				},
				{
					"Effect": "Allow",
					"Resource": [
						"arn:aws:s3:::codepipeline-eu-central-1-*"
					],
					"Action": [
						"s3:PutObject",
						"s3:GetObject",
						"s3:GetObjectVersion",
						"s3:GetBucketAcl",
						"s3:GetBucketLocation"
					]
				},
				{
					"Effect": "Allow",
					"Action": [
						"codebuild:CreateReportGroup",
						"codebuild:CreateReport",
						"codebuild:UpdateReport",
						"codebuild:BatchPutTestCases",
						"codebuild:BatchPutCodeCoverages"
					],
					"Resource": [
						"arn:aws:codebuild:eu-central-1:%s:report-group/cnappgoat-codebuild-*"
					]
				}
			]
		}`, accountID, accountID, accountID)

		// Create a Role Policy and attach it to the Role
		_, err = iam.NewRolePolicy(ctx, "codebuildRolePolicy", &iam.RolePolicyArgs{
			Role:   codebuildRole.Name,
			Policy: pulumi.String(policy),
		})

		// Create AWS CodeBuild project
		_, err = codebuild.NewProject(ctx, "cnappgoat-codebuild", &codebuild.ProjectArgs{
			Name: pulumi.String("cnappgoat-codebuild"),
			Artifacts: codebuild.ProjectArtifactsArgs{
				Type: pulumi.String("NO_ARTIFACTS"),
			},
			Environment: codebuild.ProjectEnvironmentArgs{
				PrivilegedMode: pulumi.BoolPtr(true),
				ComputeType:    pulumi.String("BUILD_GENERAL1_SMALL"),
				Image:          pulumi.String("aws/codebuild/amazonlinux2-x86_64-standard:3.0"),
				Type:           pulumi.String("LINUX_CONTAINER"),
			},
			Source: codebuild.ProjectSourceArgs{
				Type: pulumi.String("NO_SOURCE"),
				Buildspec: pulumi.String(`version: 0.2
phases:
  build:
    commands:
      - echo "administrator-123151010-21.139.152.142-rdp" >> my-rdp-creds.txt;`),
			},
			ServiceRole:  codebuildRole.Arn,
			BuildTimeout: pulumi.Int(5),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
