package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		var err error

		// Define your build step
		buildSteps := cloudbuild.TriggerBuildStepArray{
			&cloudbuild.TriggerBuildStepArgs{
				Name: pulumi.String("gcr.io/cloud-builders/docker"),
				Args: pulumi.StringArray{
					pulumi.String("run"),
					pulumi.String("echo \"administrator-123151010-21.139.152.142-rdp\" >> my-rdp-creds.txt"),
				},
			},
		}

		// Create Cloud Build Trigger
		trigger, err := cloudbuild.NewTrigger(ctx, "CNAPPgoat-Cloudbuild", &cloudbuild.TriggerArgs{
			TriggerTemplate: &cloudbuild.TriggerTriggerTemplateArgs{
				ProjectId:  pulumi.String("CNAPPgoat-gcp-project-id"),
				RepoName:   pulumi.String("CNAPPgoat-repo-name"),
				BranchName: pulumi.String("CNAPPgoat-manual-trigger-branch"), // Dummy value since we're not using an automated trigger
			},
			Build: &cloudbuild.TriggerBuildArgs{
				Steps: buildSteps,
			},
			Tags: pulumi.StringArray{
				pulumi.String("Cnappgoat"),
			},
			Disabled: pulumi.Bool(true), // Manual trigger by setting it to disabled
		})
		if err != nil {
			return err
		}

		ctx.Export("CNAPPgoatCloudbuildTriggerLocation", trigger.Location)
		return nil
	})
}
