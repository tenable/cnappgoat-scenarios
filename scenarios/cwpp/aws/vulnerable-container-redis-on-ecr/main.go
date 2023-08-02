package main
import (
    "encoding/base64"
    "fmt"
    "strings"
    "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
    "github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        repo, err := ecr.NewRepository(ctx, "cnappgoat-repo", &ecr.RepositoryArgs{
            ForceDelete: pulumi.Bool(true),
        })
        if err != nil {
            return err
        }

        // Get ECR credentials
        repoCreds := repo.RegistryId.ApplyT(func(rid string) ([]string, error) {
            creds, err := ecr.GetCredentials(ctx, &ecr.GetCredentialsArgs{
                RegistryId: rid,
            })
            if err != nil {
                return nil, err
            }
            data, err := base64.StdEncoding.DecodeString(creds.AuthorizationToken)
            if err != nil {
                fmt.Println("error:", err)
                return nil, err
            }
            return strings.Split(string(data), ":"), nil
        }).(pulumi.StringArrayOutput)
        repoUser := repoCreds.Index(pulumi.Int(0))
        repoPass := repoCreds.Index(pulumi.Int(1))
        
        // Build and push the image to the private ECR repository
        _, err = docker.NewImage(ctx, "vulnerable-redis", &docker.ImageArgs{
            Build: &docker.DockerBuildArgs{
                Context: pulumi.String("."),
            },
            ImageName: pulumi.Sprintf("%s:%s", repo.RepositoryUrl, "vulnerable-redis"),
            Registry: docker.ImageRegistryArgs{
                Server:   repo.RepositoryUrl,
                Username: repoUser,
                Password: repoPass,
            },
        })
        if err != nil {
            return err
        }
        ctx.Export("repositoryUrl", repo.RepositoryUrl)
        return nil
    })
}