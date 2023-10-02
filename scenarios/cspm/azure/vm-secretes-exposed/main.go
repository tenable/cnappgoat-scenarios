package main

import (
	"encoding/base64"
	"github.com/pulumi/pulumi-azure-native-sdk/compute"
	"github.com/pulumi/pulumi-azure-native-sdk/network"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "azure-native")
		azureLocation := cfg.Require("location")
		// create a resource group
		rg, err := resources.NewResourceGroup(ctx, "cnappgoat-rg", &resources.ResourceGroupArgs{
			Location: pulumi.String(azureLocation),
		})
		// Create a new VNet
		vnet, err := network.NewVirtualNetwork(ctx, "CNAPPgoat-vm-secrets-expose-vnet", &network.VirtualNetworkArgs{
			AddressSpace: network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
			},
			ResourceGroupName: rg.Name,
		})
		if err != nil {
			return err
		}

		// Create a new subnet
		subnet, err := network.NewSubnet(ctx, "CNAPPgoat-vm-secrets-expose-subnet", &network.SubnetArgs{
			ResourceGroupName:  rg.Name,
			VirtualNetworkName: vnet.Name,
			AddressPrefix:      pulumi.String("10.0.1.0/24"),
		})
		if err != nil {
			return err
		}

		// Create a public IP
		publicIp, err := network.NewPublicIPAddress(ctx, "vm-public-ip", &network.PublicIPAddressArgs{
			ResourceGroupName:        rg.Name,
			Location:                 vnet.Location,
			PublicIPAllocationMethod: pulumi.String("Dynamic"),
		})
		if err != nil {
			return err
		}

		// Create a network interface for the VM
		nic, err := network.NewNetworkInterface(ctx, "cnappgoat-vm-nic", &network.NetworkInterfaceArgs{
			ResourceGroupName: rg.Name,
			Location:          vnet.Location,
			IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
				&network.NetworkInterfaceIPConfigurationArgs{
					PublicIPAddress: &network.PublicIPAddressTypeArgs{
						Id: publicIp.ID(),
					},
					Subnet: &network.SubnetTypeArgs{
						Id: subnet.ID(),
					},
					Name: pulumi.String("cnappgoat-vm-nic-ipconfig1"),
				},
			},
		})
		if err != nil {
			return err
		}
		// generate random password
		password, err := random.NewRandomPassword(ctx, "password", &random.RandomPasswordArgs{
			Length:          pulumi.Int(16),
			Special:         pulumi.Bool(true),
			OverrideSpecial: pulumi.String("!#$%&*()-_=+[]{}<>:?"),
		})
		customData := "mysql -u cnaappgoat -p mysecretpassword1231"
		// b64 encode the custom data
		customDataB64 := base64.StdEncoding.EncodeToString([]byte(customData))
		// Create a VM
		vm, err := compute.NewVirtualMachine(ctx, "CNAPPgoat-vm-secrets-expose-instance", &compute.VirtualMachineArgs{
			ResourceGroupName: rg.Name,
			Location:          vnet.Location,
			HardwareProfile: &compute.HardwareProfileArgs{
				VmSize: pulumi.String("Standard_B1s"),
			},
			OsProfile: &compute.OSProfileArgs{
				AdminUsername: pulumi.String("adminuser"),
				AdminPassword: password.Result,
				CustomData:    pulumi.String(customDataB64),
				ComputerName:  pulumi.String("cnappgoatvm"),
			},
			NetworkProfile: &compute.NetworkProfileArgs{
				NetworkInterfaces: compute.NetworkInterfaceReferenceArray{
					&compute.NetworkInterfaceReferenceArgs{
						Id: nic.ID(),
					},
				},
			},
			StorageProfile: &compute.StorageProfileArgs{
				OsDisk: &compute.OSDiskArgs{
					CreateOption: pulumi.String("FromImage"),
					ManagedDisk: &compute.ManagedDiskParametersArgs{
						StorageAccountType: pulumi.String("Standard_LRS"),
					},
				},
				ImageReference: &compute.ImageReferenceArgs{
					Publisher: pulumi.String("Canonical"),
					Offer:     pulumi.String("UbuntuServer"),
					Sku:       pulumi.String("16.04-LTS"),
					Version:   pulumi.String("latest"),
				},
			},

			Tags: pulumi.StringMap{
				"Name":      pulumi.String("CNAPPgoat-vm-secrets-expose-instance"),
				"Cnappgoat": pulumi.String("true"),
			},
		})

		if err != nil {
			return err
		}

		ctx.Export("CNAPPgoat-vm-secrets-expose-vnet", vnet.Name)
		ctx.Export("CNAPPgoat-vm-secrets-expose-subnet", subnet.Name)
		ctx.Export("CNAPPgoat-vm-secrets-expose-nic", nic.Name)
		ctx.Export("CNAPPgoat-vm-secrets-expose-instance", vm.Name)
		ctx.Export("CNAPPgoat-vm-secrets-expose-password", password.Result)
		return nil
	})
}
