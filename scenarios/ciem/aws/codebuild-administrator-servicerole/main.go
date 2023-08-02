package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/codebuild"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create IAM Role for CodeBuild
		role, err := iam.NewRole(ctx, "codeBuildRole", &iam.RoleArgs{
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
		_, err = codebuild.NewProject(ctx, "CNAPPGoatCodeBuildproject", &codebuild.ProjectArgs{
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
		})
		if err != nil {
			return err
		}
		ctx.Export("roleName", role.Name)
		ctx.Export("roleArn", role.Arn)
		return nil
	})
}
