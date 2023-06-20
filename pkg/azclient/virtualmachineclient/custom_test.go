// /*
// Copyright The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

// Code generated by client-gen. DO NOT EDIT.
package virtualmachineclient

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armcompute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	armnetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
	. "github.com/onsi/gomega"
)

var (
	networkClientFactory  *armnetwork.ClientFactory
	computeClientFactory  *armcompute.ClientFactory
	virtualNetworksClient *armnetwork.VirtualNetworksClient
	nicClient             *armnetwork.InterfacesClient
	vNet                  *armnetwork.VirtualNetwork
	nicResource           *armnetwork.Interface
)

func init() {
	additionalTestCases = func() {
	}

	beforeAllFunc = func(ctx context.Context) {
		networkClientFactory, err := armnetwork.NewClientFactory(recorder.SubscriptionID(), recorder.TokenCredential(), &arm.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: recorder.HTTPClient(),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		virtualNetworksClient = networkClientFactory.NewVirtualNetworksClient()
		vnetpoller, err := virtualNetworksClient.BeginCreateOrUpdate(ctx, resourceGroupName, "vnet1", armnetwork.VirtualNetwork{
			Location: to.Ptr(location),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.Ptr("10.1.0.0/16"),
					},
				},
				Subnets: []*armnetwork.Subnet{
					{
						Name: to.Ptr("subnet1"),
						Properties: &armnetwork.SubnetPropertiesFormat{
							AddressPrefix: to.Ptr("10.1.0.0/24"),
						},
					},
				},
			},
		}, nil)
		Expect(err).NotTo(HaveOccurred())

		vnetresp, err := vnetpoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 1 * time.Second,
		})
		Expect(err).NotTo(HaveOccurred())
		vNet = &vnetresp.VirtualNetwork
		computeClientFactory, err = armcompute.NewClientFactory(subscriptionID, recorder.TokenCredential(), &arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Transport: recorder.HTTPClient(),
			},
		})
		Expect(err).NotTo(HaveOccurred())

		nicClient = networkClientFactory.NewInterfacesClient()
		nicPoller, err := nicClient.BeginCreateOrUpdate(ctx, resourceGroupName, "nic1", armnetwork.Interface{
			Location: to.Ptr(location),
			Properties: &armnetwork.InterfacePropertiesFormat{
				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
					{
						Name: to.Ptr("ipConfig1"),
						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
							Subnet:                    vNet.Properties.Subnets[0],
						},
					},
				},
			},
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		nicResp, err := nicPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 1 * time.Second,
		})
		nicResource = &nicResp.Interface

		newResource = &armcompute.VirtualMachine{
			Location: to.Ptr(location),
			Identity: &armcompute.VirtualMachineIdentity{
				Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
			},
			Properties: &armcompute.VirtualMachineProperties{
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						// search image reference
						// az vm image list --output table
						Offer:     to.Ptr("WindowsServer"),
						Publisher: to.Ptr("MicrosoftWindowsServer"),
						SKU:       to.Ptr("2019-Datacenter"),
						Version:   to.Ptr("latest"),
						//require ssh key for authentication on linux
						//Offer:     to.Ptr("UbuntuServer"),
						//Publisher: to.Ptr("Canonical"),
						//SKU:       to.Ptr("18.04-LTS"),
						//Version:   to.Ptr("latest"),
					},
					OSDisk: &armcompute.OSDisk{
						Name:         to.Ptr("disk1"),
						CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
						Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
						ManagedDisk: &armcompute.ManagedDiskParameters{
							StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
						},
						//DiskSizeGB: to.Ptr[int32](100), // default 127G
					},
				},
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes("Standard_F2s")), // VM size include vCPUs,RAM,Data Disks,Temp storage.
				},
				OSProfile: &armcompute.OSProfile{ //
					ComputerName:  to.Ptr("sample-compute"),
					AdminUsername: to.Ptr("sample-user"),
					AdminPassword: to.Ptr("Password01!@#"),
					//require ssh key for authentication on linux
					//LinuxConfiguration: &armcompute.LinuxConfiguration{
					//	DisablePasswordAuthentication: to.Ptr(true),
					//	SSH: &armcompute.SSHConfiguration{
					//		PublicKeys: []*armcompute.SSHPublicKey{
					//			{
					//				Path:    to.Ptr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", "sample-user")),
					//				KeyData: to.Ptr(string(sshBytes)),
					//			},
					//		},
					//	},
					//},
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							ID: nicResource.ID,
						},
					},
				},
			},
		}

	}
	afterAllFunc = func(ctx context.Context) {
		nicPoller, err := nicClient.BeginDelete(ctx, resourceGroupName, *nicResource.Name, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = nicPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 1 * time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		vnetPoller, err := virtualNetworksClient.BeginDelete(ctx, resourceGroupName, *vNet.Name, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = vnetPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 1 * time.Second,
		})
		Expect(err).NotTo(HaveOccurred())
	}
}
