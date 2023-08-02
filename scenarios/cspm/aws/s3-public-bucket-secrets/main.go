package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new S3 bucket
		bucket, err := s3.NewBucket(ctx, "pulumi-bucket1231", nil)
		if err != nil {
			return err
		}
		_, err = s3.NewBucketPublicAccessBlock(ctx, "exampleBucketPublicAccessBlock", &s3.BucketPublicAccessBlockArgs{
			Bucket:                bucket.ID(),
			BlockPublicAcls:       pulumi.Bool(false),
			BlockPublicPolicy:     pulumi.Bool(false),
			IgnorePublicAcls:      pulumi.Bool(false),
			RestrictPublicBuckets: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		bucket.ID().ApplyT(func(id string) error {
			fmt.Printf("Bucket name: %s\n", id)
			return nil
		})

		// Set bucket policy to make it publicly readable
		_, err = s3.NewBucketPolicy(ctx, "bucketPolicy", &s3.BucketPolicyArgs{
			Bucket: bucket.ID(),
			Policy: bucket.ID().ApplyT(func(id pulumi.String) (pulumi.String, error) {
				return pulumi.String(fmt.Sprintf(`{
                    "Version": "2012-10-17",
                    "Statement": [
                      {
                        "Effect": "Allow",
                        "Principal": "*",
                        "Action": "s3:GetObject",
                        "Resource": "arn:aws:s3:::%s/*"
                      }
                    ]
                  }`, id)), nil
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}

		// Upload a secret file to the bucket
		_, err = s3.NewBucketObject(ctx, "mysecret", &s3.BucketObjectArgs{
			Bucket:      bucket.ID(),
			Key:         pulumi.String("mysecret"),
			Source:      pulumi.NewFileAsset("secret.txt"),
			ContentType: pulumi.String("text/plain"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
