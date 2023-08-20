package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ebs"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new EBS volume
		volume, err := ebs.NewVolume(ctx, "CNAPPgoat-ebs-volume", &ebs.VolumeArgs{
			AvailabilityZone: pulumi.String("eu-central-1a"),
			Size:             pulumi.Int(8),
			Tags: pulumi.StringMap{
				"Name":      pulumi.String("CNAPPgoat-ebs-volume"),
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}

		// Create a snapshot of the EBS volume
		snapshot, err := ebs.NewSnapshot(ctx, "CNAPPgoat-ebs-snapshot", &ebs.SnapshotArgs{
			VolumeId: volume.ID(),
			Tags: pulumi.StringMap{
				"Name":      pulumi.String("CNAPPgoat-ebs-snapshot"),
				"Cnappgoat": pulumi.String("true"),
			},
		})
		if err != nil {
			return err
		}

		// Register the snapshot as a public AMI
		ami, erra := ec2.NewAmi(ctx, "CNAPPgoat-public-ami", &ec2.AmiArgs{
			Name:               pulumi.String("CNAPPgoat-public-ami"),
			Description:        pulumi.String("My AMI"),
			VirtualizationType: pulumi.String("hvm"),
			RootDeviceName:     pulumi.String("/dev/sda1"),
			EbsBlockDevices: ec2.AmiEbsBlockDeviceArray{
				ec2.AmiEbsBlockDeviceArgs{
					DeviceName:          pulumi.String("/dev/sda1"),
					SnapshotId:          snapshot.ID(),
					DeleteOnTermination: pulumi.Bool(true),
					VolumeSize:          pulumi.Int(8),
				},
			},
		})
		if erra != nil {
			return erra
		}
		_, errn := ec2.NewAmiLaunchPermission(ctx, "CNAPPgoat-ami-launchpermission", &ec2.AmiLaunchPermissionArgs{
			Group:   pulumi.String("all"),
			ImageId: ami.ID(),
		})
		if errn != nil {
			return errn
		}
		ctx.Export("CNAPPgoat-ebs-volume", volume.Arn)
		ctx.Export("CNAPPgoat-ebs-snapshot", snapshot.Arn)
		ctx.Export("CNAPPgoat-public-ami", ami.Arn)
		return nil
	})
}
