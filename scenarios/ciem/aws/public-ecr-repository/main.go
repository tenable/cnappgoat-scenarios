package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a public ECR Repository
		repo, err := ecr.NewRepository(ctx, "cnappgoat-public-ecr-repo", &ecr.RepositoryArgs{
			ImageTagMutability: pulumi.String("MUTABLE"),
			ImageScanningConfiguration: &ecr.RepositoryImageScanningConfigurationArgs{
				ScanOnPush: pulumi.Bool(true),
			},
			Tags: pulumi.StringMap{
				"Cnappgoat": pulumi.String("true"),
			},
		})

		// Set the repository policy
		_, err = ecr.NewRepositoryPolicy(ctx, "my-repo-policy", &ecr.RepositoryPolicyArgs{
			Repository: repo.Name,
			Policy: pulumi.String(`{
                "Version": "2008-10-17",
                "Statement": [
                    {
                        "Sid": "ECR Repository Policy",
                        "Effect": "Allow",
                        "Principal": {
                            "AWS": "*"
                        },
                        "Action": [
                            "ecr:DescribeImages",
                            "ecr:DescribeRepositories",
							"ecr:BatchGetImage",
							"ecr:GetDownloadUrlForLayer"
                        ]
                    }
                ]
            }`),
		})
		if err != nil {
			return err
		}
		ctx.Export("ecrRepoArn", repo.Arn)
		return nil
	})
}
