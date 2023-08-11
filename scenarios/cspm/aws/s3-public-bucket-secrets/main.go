package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new S3 bucket
		bucket, err := s3.NewBucketV2(
			ctx,
			"CnappgoatPublicBucket",
			&s3.BucketV2Args{
				Tags: pulumi.StringMap{
					"Cnappgoat": pulumi.String("true"),
				},
			})

		if err != nil {
			return err
		}
		_, err = s3.NewBucketPublicAccessBlock(ctx, "cnappgoatBucketPublicAccessBlock", &s3.BucketPublicAccessBlockArgs{
			Bucket:                bucket.ID(),
			BlockPublicAcls:       pulumi.Bool(false),
			BlockPublicPolicy:     pulumi.Bool(false),
			IgnorePublicAcls:      pulumi.Bool(false),
			RestrictPublicBuckets: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

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
		bucketObject, err := s3.NewBucketObject(ctx, "CnappgoatSecret", &s3.BucketObjectArgs{
			Bucket:      bucket.ID(),
			Key:         pulumi.String("CnappgoatSecret"),
			Source:      pulumi.NewFileAsset("secret.txt"),
			ContentType: pulumi.String("text/plain"),
		})
		if err != nil {
			return err
		}
		ctx.Export("bucket", bucket.Arn)
		ctx.Export("object-key", bucketObject.Key)
		return nil
	})
}
