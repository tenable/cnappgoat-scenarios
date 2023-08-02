package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new S3 bucket
		bucket, err := s3.NewBucket(ctx, "myhttp-bucket1311", nil)
		if err != nil {
			return err
		}

		// Upload a secret file to the bucket
		_, err = s3.NewBucketObject(ctx, "httpdata", &s3.BucketObjectArgs{
			Bucket:      bucket.ID(),
			Key:         pulumi.String("httpdata"),
			Source:      pulumi.NewFileAsset("http.txt"),
			ContentType: pulumi.String("text/plain"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
