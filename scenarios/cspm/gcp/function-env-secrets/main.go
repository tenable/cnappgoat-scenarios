package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudfunctions"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Define your secrets
		secret1 := "administrator"
		secret2 := "123156asd!@%!#^3a"

		// Create GCS Bucket to store Cloud Function source code
		cfg := config.New(ctx, "gcp")
		gcpRegion := cfg.Require("region")

		bucket, err := storage.NewBucket(ctx, "cnappgoat-cloudfuncbucket", &storage.BucketArgs{
			Location: pulumi.String(gcpRegion),
		})
		if err != nil {
			return err
		}

		// Upload data to the bucket
		archive, err := storage.NewBucketObject(ctx, "CNAPPgoat-public-data", &storage.BucketObjectArgs{
			Bucket: bucket.Name,
			Source: pulumi.NewFileAsset("./app.zip"),
		})
		if err != nil {
			return err
		}
		secrets := map[string]string{
			"SECRET1": secret1,
			"SECRET2": secret2,
		}

		// Convert it to a pulumi.MapInput
		envVars := pulumi.Map{}
		for k, v := range secrets {
			envVars[k] = pulumi.String(v)
		}
		// Create Cloud Function with an HTTP trigger
		function, err := cloudfunctions.NewFunction(ctx, "CNAPPgoat-cloudfunction", &cloudfunctions.FunctionArgs{
			SourceArchiveBucket:  bucket.Name,
			Runtime:              pulumi.String("nodejs14"),
			EntryPoint:           pulumi.String("handler"),
			SourceArchiveObject:  archive.Name,
			EnvironmentVariables: envVars,
			TriggerHttp:          pulumi.Bool(true),
			AvailableMemoryMb:    pulumi.Int(128),
			Region:               pulumi.String("us-central1"), // specify an appropriate region
		})
		if err != nil {
			return err
		}

		ctx.Export("CNAPPgoat-cloudfunction-url", function.HttpsTriggerUrl)
		return nil
	})
}
