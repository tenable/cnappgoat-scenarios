package main

import (
	"github.com/pulumi/pulumi-azure-native-sdk/containerregistry"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read the desired region from the Pulumi configuration
        cfg := config.New(ctx, "azure-native")
        azureLocation := cfg.Require("location")

		// Create an Azure Resource Group
		resourceGroup, err := resources.NewResourceGroup(ctx, "cnappgoat", &resources.ResourceGroupArgs{
			Location: pulumi.String(azureLocation),
		})
		if err != nil {
			return err
		}

		// Create an Azure Container Registry
		acr, err := containerregistry.NewRegistry(ctx, "cnappgoatACR", &containerregistry.RegistryArgs{
			ResourceGroupName: resourceGroup.Name,
			Location:          resourceGroup.Location,
			Sku: &containerregistry.SkuArgs{
				Name: pulumi.String("Basic"),
			},
			AdminUserEnabled: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Get the ACR credentials
		acrCreds := pulumi.All(resourceGroup.Name, acr.Name).ApplyT(
			func(args []interface{}) (*containerregistry.ListRegistryCredentialsResult, error) {
				resourceGroupName := args[0].(string)
				registryName := args[1].(string)
				return containerregistry.ListRegistryCredentials(ctx, &containerregistry.ListRegistryCredentialsArgs{
					ResourceGroupName: resourceGroupName,
					RegistryName:      registryName,
				})
			},
		)

		adminUsername := acrCreds.ApplyT(func(result interface{}) (string, error) {
			credentials := result.(*containerregistry.ListRegistryCredentialsResult)
			return *credentials.Username, nil
		}).(pulumi.StringOutput)
		adminPassword := acrCreds.ApplyT(func(result interface{}) (string, error) {
			credentials := result.(*containerregistry.ListRegistryCredentialsResult)
			return *credentials.Passwords[0].Value, nil
		}).(pulumi.StringOutput)

		// Build and push the image to the private ACR repository
		_, err = docker.NewImage(ctx, "vulnerable-thunderbird", &docker.ImageArgs{
			Build: &docker.DockerBuildArgs{
				Context: pulumi.String("."),
				Dockerfile: pulumi.String("Dockerfile"),
			},
			ImageName: pulumi.Sprintf("%s/%s:%s", acr.LoginServer, "vulnerable-thunderbird", "latest"),
			Registry: docker.ImageRegistryArgs{
				Server:   acr.LoginServer,
				Username: adminUsername,
				Password: adminPassword,
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("acrLoginServer", acr.LoginServer)
		return nil
	})
}
