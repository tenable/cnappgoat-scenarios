package main

import (
	"github.com/pulumi/pulumi-azure-native-sdk/compute/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/network/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/resources/v2"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "azure-native")
		azureLocation := cfg.Require("location")
		// create a resource group
		rg, err := resources.NewResourceGroup(ctx, "cnappgoat-vm", nil)
		if err != nil {
			return err
		}
		// Create a new VNet
		vnet, err := network.NewVirtualNetwork(ctx, "cnappgoat-vnet", &network.VirtualNetworkArgs{
			ResourceGroupName: rg.Name,
			AddressSpace: &network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
			},
		})
		if err != nil {
			return err
		}
		// Create a new NSG with an open rule for port 3389
		nsg, err := network.NewNetworkSecurityGroup(ctx, "cnappgoat-nsg", &network.NetworkSecurityGroupArgs{
			ResourceGroupName: rg.Name,
			SecurityRules: network.SecurityRuleTypeArray{
				network.SecurityRuleTypeArgs{
					Access:                   pulumi.String("Allow"),
					DestinationPortRange:     pulumi.String("3389"),
					Priority:                 pulumi.Int(100),
					Protocol:                 pulumi.String("Tcp"),
					SourcePortRange:          pulumi.String("*"),
					Direction:                pulumi.String("Inbound"),
					SourceAddressPrefix:      pulumi.String("*"),
					DestinationAddressPrefix: pulumi.String("*"),
					Name:                     pulumi.String("rdpopen"),
				}},
		})
		if err != nil {
			return err
		}

		// Create a new subnet
		subnet, err := network.NewSubnet(ctx, "cnappgoat-subnet", &network.SubnetArgs{
			ResourceGroupName:  rg.Name,
			VirtualNetworkName: vnet.Name,
			AddressPrefix:      pulumi.String("10.0.1.0/24"),
			NetworkSecurityGroup: &network.NetworkSecurityGroupTypeArgs{
				Id: nsg.ID(),
			},
		})
		if err != nil {
			return err
		}

		// Create a public IP
		publicIp, err := network.NewPublicIPAddress(ctx, "cnappgoat-ip", &network.PublicIPAddressArgs{
			Sku: &network.PublicIPAddressSkuArgs{
				Name: pulumi.String("Basic"),
			},
			ResourceGroupName: rg.Name,
		})
		if err != nil {
			return err
		}

		// Create a network interface for the VM
		nic, err := network.NewNetworkInterface(ctx, "vm-nic", &network.NetworkInterfaceArgs{
			ResourceGroupName: rg.Name,
			IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
				&network.NetworkInterfaceIPConfigurationArgs{
					Subnet: &network.SubnetTypeArgs{
						Id: subnet.ID(),
					},
					PrivateIPAllocationMethod: pulumi.String("Dynamic"),
					PublicIPAddress:           &network.PublicIPAddressTypeArgs{Id: publicIp.ID()},
					Name:                      pulumi.String("vm-nic-ipconfig"),
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

		// Create a VM
		vm, err := compute.NewVirtualMachine(ctx, "CNAPPgoat-vm-unrestricted-rdp-instance", &compute.VirtualMachineArgs{
			ResourceGroupName: rg.Name,
			Location:          pulumi.String(azureLocation),
			HardwareProfile: &compute.HardwareProfileArgs{
				VmSize: pulumi.String("Standard_B1s"),
			},
			OsProfile: &compute.OSProfileArgs{
				ComputerName:  pulumi.String("cnappgoatvm"),
				AdminUsername: pulumi.String("adminuser"),
				AdminPassword: password.Result,
				WindowsConfiguration: &compute.WindowsConfigurationArgs{
					EnableAutomaticUpdates: pulumi.Bool(true),
				},
			},
			StorageProfile: &compute.StorageProfileArgs{
				OsDisk: &compute.OSDiskArgs{
					CreateOption: pulumi.String("FromImage"),
					ManagedDisk: &compute.ManagedDiskParametersArgs{
						StorageAccountType: pulumi.String("Standard_LRS"),
					},
				},
				// Select a windows image - note Publisher, Offer, SKU and Version
				ImageReference: &compute.ImageReferenceArgs{
					Publisher: pulumi.String("MicrosoftWindowsServer"),
					Offer:     pulumi.String("WindowsServer"),
					Sku:       pulumi.String("2019-Datacenter"),
					Version:   pulumi.String("latest"),
				},
			},
			NetworkProfile: &compute.NetworkProfileArgs{
				NetworkInterfaces: compute.NetworkInterfaceReferenceArray{
					&compute.NetworkInterfaceReferenceArgs{
						Id: nic.ID(),
					},
				},
			},
			Tags: pulumi.StringMap{
				"Name":      pulumi.String("CNAPPgoat-vm-unrestricted-rdp-instance"),
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("CNAPPgoat-vm-unrestricted-rdp-vnet", vnet.Name)
		ctx.Export("CNAPPgoat-vm-unrestricted-rdp-subnet", subnet.Name)
		ctx.Export("CNAPPgoat-vm-unrestricted-rdp-nsg", nsg.Name)
		ctx.Export("CNAPPgoat-vm-unrestricted-rdp-instance", vm.Name)
		ctx.Export("CNAPPgoat-vm-unrestricted-rdp-password", password.Result)
		return nil
	})
}
