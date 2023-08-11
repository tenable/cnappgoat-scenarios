package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/codebuild"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create IAM Role for CodeBuild
		role, err := iam.NewRole(ctx, "CnappgoatCodeBuildRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
				  {
					"Action": "sts:AssumeRole",
					"Principal": {
					  "Service": "codebuild.amazonaws.com"
					},
					"Effect": "Allow",
					"Sid": ""
				  }
				]
			  }`),
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		// Attach the AWSCodeBuildAdminAccess managed policy to the role
		_, err = iam.NewRolePolicyAttachment(ctx, "policyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		})
		if err != nil {
			return err
		}

		// Create AWS CodeBuild project
		codebuildProject, err := codebuild.NewProject(ctx, "CNAPPgoatCodeBuildproject", &codebuild.ProjectArgs{
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
      - echo test`),
			},
			ServiceRole:  role.Arn,
			BuildTimeout: pulumi.Int(5),
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("roleArn", role.Arn)
		ctx.Export("codebuildProject", codebuildProject.Arn)
		return nil
	})
}
