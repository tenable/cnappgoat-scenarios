package main

import (
    "encoding/base64"
    "crypto/rand"
    "crypto/rsa"
    "golang.org/x/crypto/ssh"
    "github.com/pulumi/pulumi-azure-native-sdk/compute"
    "github.com/pulumi/pulumi-azure-native-sdk/network"
    "github.com/pulumi/pulumi-azure-native-sdk/resources"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func generatePublicKey() (string, error) {
    privateRsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return "", err
    }

    publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
    if err != nil {
        return "", err
    }

    publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
    publicKeyStr := string(publicKeyBytes)

    return publicKeyStr, nil
}


func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        cfg := config.New(ctx, "azure-native")
        azureLocation := cfg.Require("location")

        resourceGroup, err := resources.NewResourceGroup(ctx, "cnappgoat_rg", &resources.ResourceGroupArgs{
            Location: pulumi.String(azureLocation),
        })
        if err != nil {
            return err
        }

        vnet, err := network.NewVirtualNetwork(ctx, "cnappgoat_vnet", &network.VirtualNetworkArgs{
            ResourceGroupName: resourceGroup.Name,
            AddressSpace: network.AddressSpaceArgs{
                AddressPrefixes: pulumi.StringArray{pulumi.String("10.0.0.0/16")},
            },
        })
        if err != nil {
            return err
        }

        subnet, err := network.NewSubnet(ctx, "cnappgoat_subnet", &network.SubnetArgs{
            ResourceGroupName:  resourceGroup.Name,
            VirtualNetworkName: vnet.Name,
            AddressPrefix:    pulumi.String("10.0.1.0/24"),
        })
        if err != nil {
            return err
        }

        nic, err := network.NewNetworkInterface(ctx, "cnappgoat_vmNic", &network.NetworkInterfaceArgs{
            ResourceGroupName: resourceGroup.Name,
            IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
                &network.NetworkInterfaceIPConfigurationArgs{
                    Name: pulumi.String("vmNicIpConfig"),
                    Subnet: &network.SubnetTypeArgs{
                        Id: subnet.ID(),
                    },
                    PrivateIPAllocationMethod: pulumi.String("Dynamic"),
                },
            },
        })
        if err != nil {
            return err
        }

        customData := `#!/bin/bash
sudo apt-get update -y
sudo apt-get install -y docker.io
sudo service docker start
sudo usermod -a -G docker $USER
sudo docker run --name end_of_life_container -d -p 80:80 public.ecr.aws/i3j2g7c0/cnappgoat-images:end_of_life_ubuntu2110_image`

        encodedCustomData := base64.StdEncoding.EncodeToString([]byte(customData))

        // Generate a public key
        publicKey, err := generatePublicKey()
        if err != nil {
            return err
        }

        _, err = compute.NewVirtualMachine(ctx, "cnappgoat_vm", &compute.VirtualMachineArgs{
            ResourceGroupName: resourceGroup.Name,
            NetworkProfile: &compute.NetworkProfileArgs{
                NetworkInterfaces: compute.NetworkInterfaceReferenceArray{
                    &compute.NetworkInterfaceReferenceArgs{
                        Id:      nic.ID(),
                        Primary: pulumi.Bool(true),
                    },
                },
            },
            HardwareProfile: &compute.HardwareProfileArgs{
                VmSize: pulumi.String("Standard_B2s"),
            },
            OsProfile: &compute.OSProfileArgs{
                ComputerName:  pulumi.String("cnappgoatvm"),
                AdminUsername: pulumi.String("adminuser"),
                CustomData:    pulumi.String(encodedCustomData),
                LinuxConfiguration: &compute.LinuxConfigurationArgs{
                    DisablePasswordAuthentication: pulumi.Bool(true),
                    Ssh: &compute.SshConfigurationArgs{
                        PublicKeys: compute.SshPublicKeyTypeArray{
                            &compute.SshPublicKeyTypeArgs{
                                KeyData: pulumi.String(publicKey),
                                Path:    pulumi.String("/home/adminuser/.ssh/authorized_keys"),
                            },
                        },
                    },
                },
            },
            StorageProfile: &compute.StorageProfileArgs{
                OsDisk: &compute.OSDiskArgs{
                    CreateOption:       pulumi.String("FromImage"),
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
        })
        if err != nil {
            return err
        }

        return nil
    })
}
