package main

import (
	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/config"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get the GCP project and region from the default configuration
		project := config.GetProject(ctx)
		region := config.GetRegion(ctx)

		// Generate a unique Artifact Registry repository ID
		uniqueString, err := random.NewRandomString(ctx, "unique-string", &random.RandomStringArgs{
			Length:  pulumi.Int(4),
			Lower:   pulumi.Bool(true),
			Upper:   pulumi.Bool(false),
			Numeric: pulumi.Bool(true),
			Special: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		repoId := pulumi.Sprintf("cnappgoat-%s", uniqueString.Result)

		// Create an Artifact Registry repository
		repo, err := artifactregistry.NewRepository(ctx, "repository", &artifactregistry.RepositoryArgs{
			Description:  pulumi.String("Repository for container image"),
			Format:       pulumi.String("DOCKER"),
			Location:     pulumi.String(region),
			RepositoryId: repoId,
		})
		if err != nil {
			return err
		}

		// Get client credentials
		clientConfig, err := organizations.GetClientConfig(ctx)
		if err != nil {
			return err
		}

		// Form the repository URL
		repoUrl := pulumi.Sprintf("%s-docker.pkg.dev/%s/%s", repo.Location, project, repo.RepositoryId)

		// Build and push the image to the private Artifact Registry repository
		_, err = docker.NewImage(ctx, "cnappgoat-image", &docker.ImageArgs{
			Build: &docker.DockerBuildArgs{
				Context: pulumi.String("."),
			},
			ImageName: pulumi.Sprintf("%s/%s:%s", repoUrl, "cnappgoat-public-image", "latest"),
			Registry: docker.ImageRegistryArgs{
				Server:   repoUrl,
				Username: pulumi.String("oauth2accesstoken"),
				Password: pulumi.String(clientConfig.AccessToken),
			},
		})
		if err != nil {
			return err
		}
		_, err = artifactregistry.NewRepositoryIamMember(ctx, "publicImageIamMember", &artifactregistry.RepositoryIamMemberArgs{
			Location:   pulumi.String(region),
			Repository: repo.RepositoryId,
			Role:       pulumi.String("roles/artifactregistry.reader"),
			Member:     pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}

		ctx.Export("repositoryUrl", repoUrl)
		return nil
	})
}
