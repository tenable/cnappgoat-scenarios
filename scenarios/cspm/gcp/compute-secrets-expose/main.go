package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new Network (Equivalent to VPC in AWS)
		network, err := compute.NewNetwork(ctx, "cnappgoat-compute-network", &compute.NetworkArgs{})
		if err != nil {
			return err
		}

		// Create a Subnetwork (Equivalent to Subnet in AWS)
		subnetwork, err := compute.NewSubnetwork(ctx, "cnappgoat-compute-subnetwork", &compute.SubnetworkArgs{
			Network:     network.SelfLink,
			IpCidrRange: pulumi.String("10.0.1.0/24"),
		})
		if err != nil {
			return err
		}

		// Create a new Compute Instance
		instance, err := compute.NewInstance(ctx, "cnappgoat-compute-instance", &compute.InstanceArgs{
			Zone:        pulumi.String("us-central1-a"),
			MachineType: pulumi.String("f1-micro"),
			BootDisk: &compute.InstanceBootDiskArgs{
				InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String("debian-cloud/debian-12"),
				},
			},
			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				&compute.InstanceNetworkInterfaceArgs{
					Network: network.ID(),
					AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
						&compute.InstanceNetworkInterfaceAccessConfigArgs{},
					},
				},
			},
			MetadataStartupScript: pulumi.String(`#!/bin/bash
mysql -u cnaappgoat -p mysecretpassword1231`),
			Tags: pulumi.StringArray{
				pulumi.String("cnappgoat"),
				pulumi.String("compute-instance"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("CNAPPgoat-compute-network", network.SelfLink)
		ctx.Export("CNAPPgoat-compute-subnetwork", subnetwork.SelfLink)
		ctx.Export("CNAPPgoat-compute-instance", instance.SelfLink)
		return nil
	})
}
