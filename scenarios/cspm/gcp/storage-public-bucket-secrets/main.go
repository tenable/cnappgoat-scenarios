package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a GCP storage bucket
		cfg := config.New(ctx, "gcp")
		gcpRegion := cfg.Require("region")

		bucket, err := storage.NewBucket(ctx, "cnappgoat-public-bucket", &storage.BucketArgs{
			Location: pulumi.String(gcpRegion),
		})
		if err != nil {
			return err
		}

		// Set the bucket to be publicly readable
		_, err = storage.NewBucketIAMMember(ctx, "publicRead", &storage.BucketIAMMemberArgs{
			Bucket: bucket.Name,
			Role:   pulumi.String("roles/storage.objectViewer"),
			Member: pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}

		// Upload data to the bucket
		_, err = storage.NewBucketObject(ctx, "CNAPPgoat-public-data", &storage.BucketObjectArgs{
			Bucket: bucket.Name,
			Source: pulumi.NewFileAsset("./CNAPPgoat-public-data.txt"),
		})
		if err != nil {
			return err
		}

		// Export the bucket's URL
		ctx.Export("bucketUrl", bucket.Url)
		return nil
	})
}
