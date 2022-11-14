package main

import "fmt"

func main() {
	PruneList := []string{
		"ovn.kubernetes.io/allocated",
		"ovn.kubernetes.io/cidr",
		"ovn.kubernetes.io/gateway",
		"ovn.kubernetes.io/ip_address",
		"ovn.kubernetes.io/ip_pool",
		"ovn.kubernetes.io/logical_switch",
		"ovn.kubernetes.io/mac_address",
		"ovn.kubernetes.io/network_type",
		"ovn.kubernetes.io/pod_nic_type",
		"ovn.kubernetes.io/provider_network",
		"ovn.kubernetes.io/routed",
		"ovn.kubernetes.io/vlan_id",
	}

	type Pod struct {
		Metadata struct {
			Name        string
			Annotations map[string]string
		}
	}
	Pihole := Pod{
		Metadata: struct {
			Name        string
			Annotations map[string]string
		}{
			Name: "asf",
			Annotations: map[string]string{
				"ovn.kubernetes.io/allocated":        "true",
				"backup.velero.io/backup-volumes":    "pihole-data",
				"checksum/config":                    "2486f53009d681a7631055d010668df03460ed62fb255ecb78271c88a231fbd2",
				"ovn.kubernetes.io/cidr":             "10.80.0.0/24",
				"ovn.kubernetes.io/gateway":          "10.80.0.254",
				"ovn.kubernetes.io/ip_address":       "10.80.0.10",
				"ovn.kubernetes.io/ip_pool":          "10.80.0.10",
				"ovn.kubernetes.io/logical_switch":   "80-comms",
				"ovn.kubernetes.io/network_type":     "vlan",
				"ovn.kubernetes.io/pod_nic_type":     "veth-pair",
				"ovn.kubernetes.io/provider_network": "provider",
				"ovn.kubernetes.io/routed":           "true",
				"ovn.kubernetes.io/vlan_id":          "80",
			},
		},
	}

	fmt.Printf("Input: %+v", Pihole)

	for annotation := range Pihole.Metadata.Annotations {
		for _, blacklisted := range PruneList {
			if annotation == blacklisted {
				delete(Pihole.Metadata.Annotations, annotation)
				break
			}
		}
	}

	fmt.Printf("Result: %+v", Pihole)

}
