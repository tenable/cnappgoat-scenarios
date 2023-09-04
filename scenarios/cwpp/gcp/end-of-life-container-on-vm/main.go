package main

import (
    "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {

        // Define a script to be run when the VM starts up.
        metadataStartupScript := `#!/bin/bash
sudo apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce
sudo systemctl start docker
sudo systemctl enable docker
sudo docker run --name end_of_life_container -d -p 80:80 public.ecr.aws/i3j2g7c0/cnappgoat-images:end_of_life_ubuntu2110_image
`

        // Create a new network for the virtual machine
        network, err := compute.NewNetwork(ctx, "network", &compute.NetworkArgs{
            AutoCreateSubnetworks: pulumi.Bool(false),
        })
        if err != nil {
            return err
        }

        // Create a subnet on the network
        subnet, err := compute.NewSubnetwork(ctx, "subnet", &compute.SubnetworkArgs{
            IpCidrRange: pulumi.String("10.0.1.0/24"),
            Network:     network.ID(),
        })
        if err != nil {
            return err
        }

        // Create a firewall allowing inbound access over ports 80 (for HTTP) and 22 (for SSH).
        _, err = compute.NewFirewall(ctx, "firewall", &compute.FirewallArgs{
            Network: network.SelfLink,
            Allows: compute.FirewallAllowArray{
                compute.FirewallAllowArgs{
                    Protocol: pulumi.String("tcp"),
                    Ports: pulumi.ToStringArray([]string{
                        "22",
                        "80",
                    }),
                },
            },
            Direction: pulumi.String("INGRESS"),
            SourceRanges: pulumi.ToStringArray([]string{
                "0.0.0.0/0",
            }),
        })
        if err != nil {
            return err
        }

        // Create the virtual machine.
        _, err = compute.NewInstance(ctx, "cnappgoat-instance", &compute.InstanceArgs{
            MachineType: pulumi.String("n2d-standard-4"),
            BootDisk: compute.InstanceBootDiskArgs{
                InitializeParams: compute.InstanceBootDiskInitializeParamsArgs{
                    Image: pulumi.String("debian-11"),
                },
            },
            NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
                compute.InstanceNetworkInterfaceArgs{
                    Network:    network.ID(),
                    Subnetwork: subnet.ID(),
                    AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
                        compute.InstanceNetworkInterfaceAccessConfigArgs{},
                    },
                },
            },
            MetadataStartupScript: pulumi.String(metadataStartupScript),
        })
        if err != nil {
            return err
        }

        return nil
    })
}
