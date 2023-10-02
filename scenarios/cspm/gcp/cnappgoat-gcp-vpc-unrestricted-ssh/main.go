package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a VPC Network
		network, err := compute.NewNetwork(ctx, "CNAPPgoatVPC", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create a firewall rule that allows unrestricted SSH access from the internet
		_, err = compute.NewFirewall(ctx, "CNAPPgoatUnrestrictedSSH", &compute.FirewallArgs{
			Network: network.Name,
			Allows: compute.FirewallAllowArray{
				compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.StringArray{
						pulumi.String("22"),
					},
				},
			},
			SourceRanges: pulumi.StringArray{
				pulumi.String("0.0.0.0/0"), // This represents all IP addresses
			},
		})
		if err != nil {
			return err
		}

		// Export the VPC Network name
		ctx.Export("networkName", network.Name)

		return nil
	})
}
