package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new S3 bucket
		bucket, err := s3.NewBucket(ctx, "CNAPPGoat-http-bucket", &s3.BucketArgs{
			Tags: pulumi.StringMap{
				"Name":      pulumi.String("CNAPPGoat-http-bucket"),
				"Cnappgoat": pulumi.String("true"),
			}},
		)
		if err != nil {
			return err
		}

		// Upload a secret file to the bucket
		bucketObject, err := s3.NewBucketObject(ctx, "CNAPPGoat-http-data", &s3.BucketObjectArgs{
			Bucket:      bucket.ID(),
			Key:         pulumi.String("CNAPPGoat-http-data"),
			Source:      pulumi.NewFileAsset("http.txt"),
			ContentType: pulumi.String("text/plain"),
		})
		if err != nil {
			return err
		}
		ctx.Export("CNAPPGoat-http-bucket", bucket.Arn)
		ctx.Export("object-key", bucketObject.Key)
		return nil
	})
}
