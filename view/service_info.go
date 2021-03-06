/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"

	"github.com/ernestio/ernest-cli/model"
	"github.com/olekukonko/tablewriter"
)

// PrintServiceInfo : Pretty print for service info
func PrintServiceInfo(service *model.Service) {
	fmt.Println("Name : " + service.Name)
	if service.Vpc != "" {
		fmt.Println("VPC : " + service.Vpc)
	}
	fmt.Println("Status : " + service.Status)
	fmt.Println("Date : " + service.Version)
	if service.Status == "errored" {
		if service.LastError == "" {
			fmt.Println("Last known error : unknown")
		} else {
			fmt.Println("Last known error : " + service.LastError)
		}
	}

	if len(service.ELBs) == 0 {
		fmt.Println("\nELBs (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nELBs:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "DNS Name"})
		for _, v := range service.ELBs {
			table.Append([]string{v.Name, v.DNSName})
		}
		table.Render()
	}

	if len(service.Networks) == 0 {
		fmt.Println("\nNetworks (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nNetworks:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Availability Zone"})
		for _, v := range service.Networks {
			table.Append([]string{v.Name, v.Subnet, v.AvailabilityZone})
		}
		table.Render()
	}

	if len(service.Instances) == 0 {
		fmt.Println("\nInstances (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nInstances:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Public IP", "Private IP"})
		for _, v := range service.Instances {
			table.Append([]string{v.Name, v.InstanceAWSID, v.PublicIP, v.IP})
		}
		table.Render()
	}

	if len(service.Nats) == 0 {
		fmt.Println("\nNAT gateways (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nNAT gateways:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Group ID"})
		for _, v := range service.Nats {
			table.Append([]string{v.Name, v.NatGatewayAWSID})
		}
		table.Render()
	}

	if len(service.SecurityGroups) == 0 {
		fmt.Println("\nSecurity groups (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nSecurity groups:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Group ID"})
		for _, v := range service.SecurityGroups {
			table.Append([]string{v.Name, v.SecurityGroupAWSID})
		}
		table.Render()
	}

	if len(service.RDSClusters) == 0 {
		fmt.Println("\nRDS Clusters (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nRDS Clusters:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Endpoint"})
		for _, v := range service.RDSClusters {
			table.Append([]string{v.Name, v.Endpoint})
		}
		table.Render()
	}

	if len(service.RDSInstances) == 0 {
		fmt.Println("\nRDS Instances (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nRDS Instances:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Endpoint"})
		for _, v := range service.RDSInstances {
			table.Append([]string{v.Name, v.Endpoint})
		}
		table.Render()
	}

	if len(service.EBSVolumes) == 0 {
		fmt.Println("\nEBS Volumes (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nEBS Volumes:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Volume ID"})
		for _, v := range service.EBSVolumes {
			table.Append([]string{v.Name, v.VolumeAWSID})
		}
		table.Render()
	}

}
