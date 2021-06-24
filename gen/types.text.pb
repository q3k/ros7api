sub {
  name: "interface"
  sub {
    name: "bridge"
    sub {
      # https://help.mikrotik.com/docs/display/ROS/Bridge#Bridge-BridgeVLANtable
      # /interface bridge vlan
      name: "vlan"
      record {
        description: "Bridge VLAN table represents per-VLAN port mapping with an egress VLAN tag action. The tagged ports send out frames with a corresponding VLAN ID tag. The untagged ports remove a VLAN tag before sending out frames. Bridge ports with frame-types set to admit-all or admit-only-untagged-and-priority-tagged will be automatically added as untagged ports for the pvid VLAN."
        property {
          name: "bridge" type_string { }
          description: "The bridge interface which the respective VLAN entry is intended for."
        }
        property {
          name: "disabled" type_boolean { }
          description: "Enables or disables Bridge VLAN entry."
        }
        property {
          name: "tagged" type_string_list { }
          description: "Interface list with a VLAN tag adding action in egress."
        }
        property {
          name: "untagged" type_string_list { }
          description: "Interface list with a VLAN tag removing action in egress."
        }
        property {
          name: "vlan-ids" go_name: "VlanIDs" type_number { }
          description: "The list of VLAN IDs for certain port configuration."
        }
        property {
          name: "current-tagged" read_only: true type_string_list { }
        }
        property {
          name: "current-untagged" read_only: true type_string_list { }
        }
        property {
          name: "dynamic" read_only: true type_boolean { }
        }
      }
    }
    sub {
      # https://help.mikrotik.com/docs/display/ROS/Bridge#Bridge-PortSettings
      # /interface bridge port
      name: "port"
      record {
        description: "Port submenu is used to add interfaces in a particular bridge."
        property {
          name: "auto-isolate" type_boolean { }
          description: "When enabled, prevents a port moving from discarding into forwarding state if no BPDUs are received from the neighboring bridge. The port will change into a forwarding state only when a BPDU is received. This property only has an effect when protocol-mode is set to rstp or mstp and edge is set to no."
        }
        property {
          name: "bpdu-guard" type_boolean { }
          description: "Enables or disables BPDU Guard feature on a port. This feature puts the port in a disabled role if it receives a BPDU and requires the port to be manually disabled and enabled if a BPDU was received. Should be used to prevent a bridge from BPDU related attacks. This property has no effect when protocol-mode is set to none."
        }
        property {
          name: "bridge" type_string { }
          description: "The bridge interface where the respective interface is grouped in."
        }
        property {
          name: "broadcast-flood" type_boolean { }
          description: "When enabled, bridge floods broadcast traffic to all bridge egress ports. When disabled, drops broadcast traffic on egress ports. Can be used to filter all broadcast traffic on an egress port. Broadcast traffic is considered as traffic that uses FF:FF:FF:FF:FF:FF as destination MAC address, such traffic is crucial for many protocols such as DHCP, ARP, NDP, BOOTP (Netinstall), and others. This option does not limit traffic flood to the CPU."
        }
        property {
          name: "edge" type_enum {
            variant {
              value: "auto"
              description: "same as no-discover, but will additionally detect if a bridge port is a Wireless interface with disabled bridge-mode, such interface will be automatically set as an edge port without discovery."
            }
            variant {
              value: "no"
              description: "non-edge port, will participate in learning and listening states in STP."
            }
            variant {
              value: "no-discover"
              description: "non-edge port with enabled discovery, will participate in learning and listening states in STP, a port can become an edge port if no BPDU is received."
            }
            variant {
              value: "yes"
              description: "edge port without discovery, will transit directly to forwarding state."
            }
            variant {
              value: "yes-discover"
              description: "edge port with enabled discovery, will transit directly to forwarding state."
            }
          }
          description: "Set port as edge port or non-edge port, or enable edge discovery. Edge ports are connected to a LAN that has no other bridges attached. An edge port will skip the learning and the listening states in STP and will transition directly to the forwarding state, this reduces the STP initialization time. If the port is configured to discover edge port then as soon as the bridge detects a BPDU coming to an edge port, the port becomes a non-edge port. This property has no effect when protocol-mode is set to none."
        }
        property {
          name: "fast-leave" type_boolean { }
          description: "Enables IGMP/MLD fast leave feature on the bridge port. The bridge will stop forwarding multicast traffic to a bridge port when an IGMP/MLD leave message is received. This property only has an effect when igmp-snooping is set to yes."
        }
        property {
          name: "frame-types" type_enum {
            variant { value: "admit-all" }
            variant { value: "admit-only-untagged-and-priority-tagged" }
            variant { value: "admit-only-vlan-tagged" }
          }
          description: "Specifies allowed ingress frame types on a bridge port. This property only has an effect when vlan-filtering is set to yes."
        }
        property {
          name: "ingress-filtering" type_boolean { }
          description: "Enables or disables VLAN ingress filtering, which checks if the ingress port is a member of the received VLAN ID in the bridge VLAN table. Should be used with frame-types to specify if the ingress traffic should be tagged or untagged. This property only has effect when vlan-filtering is set to yes."
        }
        property {
          name: "learn" type_enum {
            variant {
              value: "yes"
              description: "enables MAC learning"
            }
            variant {
              value: "no"
              description: "disables MAC learning"
            }
            variant {
              value: "auto"
              description: "detects if bridge port is a Wireless interface and uses a Wireless registration table instead of MAC learning, will use Wireless registration table if the Wireless interface is set to one of ap-bridge, bridge, wds-slave mode and bridge mode for the Wireless interface is disabled."
            }
          }
          description: "Changes MAC learning behavior on a bridge port"
        }
        property {
          name: "multicast-router" type_enum {
            variant {
              value: "disabled"
              description: "disabled multicast router state on the bridge port. Unregistered multicast and IGMP/MLD membership reports are not sent to the bridge port regardless of what is connected to it."
            }
            variant {
              value: "permanent"
              description: "enabled multicast router state on the bridge port. Unregistered multicast and IGMP/MLD membership reports are sent to the bridge port regardless of what is connected to it."
            }
            variant {
              value: "temporary-query"
              description: "automatically detect multicast router state on the bridge port using IGMP/MLD queries."
            }
          }
          description: "A multicast router port is a port where a multicast router or querier is connected. On this port, unregistered multicast streams and IGMP/MLD membership reports will be sent. This setting changes the state of the multicast router for bridge ports. This property can be used to send IGMP/MLD membership reports to certain bridge ports for further multicast routing or proxying. This property only has an effect when igmp-snooping is set to yes."
        }
        # This defaults to 'none', no way to represent this in Go code yet.
        #property {
        #  name: "horizon" type_number { }
        #  description: "Use split horizon bridging to prevent bridging loops. Set the same value for a group of ports, to prevent them from sending data to ports with the same horizon value. Split horizon is a software feature that disables hardware offloading."
        #}
        property {
          name: "internal-path-cost" type_number { }
          description: "Path cost to the interface for MSTI0 inside a region. This property only has effect when protocol-mode is set to mstp."
        }
        property {
          name: "interface" type_string { }
          description: "Name of the interface."
        }
        property {
          name: "path-cost" type_number { }
          description: "Path cost to the interface, used by STP to determine the best path, used by MSTP to determine the best path between regions. This property has no effect when protocol-mode is set to none."
        }
        property {
          name: "point-to-point" type_enum {
            variant { value: "auto" }
            variant { value: "yes" }
            variant { value: "no" }
          }
          description: "Specifies if a bridge port is connected to a bridge using a point-to-point link for faster convergence in case of failure. By setting this property to yes, you are forcing the link to be a point-to-point link, which will skip the checking mechanism, which detects and waits for BPDUs from other devices from this single link. By setting this property to no, you are expecting that a link can receive BPDUs from multiple devices. By setting the property to yes, you are significantly improving (R/M)STP convergence time. In general, you should only set this property to no if it is possible that another device can be connected between a link, this is mostly relevant to Wireless mediums and Ethernet hubs. If the Ethernet link is full-duplex, auto enables point-to-point functionality. This property has no effect when protocol-mode is set to none."
        }
        property {
          name: "priority" type_number { }
          description: "The priority of the interface, used by STP to determine the root port, used by MSTP to determine root port between regions."
        }
        property {
          name: "pvid" type_number { }
          description: "Port VLAN ID (pvid) specifies which VLAN the untagged ingress traffic is assigned to. This property only has an effect when vlan-filtering is set to yes."
        }
        property {
          name: "restricted-role" type_boolean { }
          description: "Enable the restricted role on a port, used by STP to forbid a port from becoming a root port. This property only has an effect when protocol-mode is set to mstp."
        }
        property {
          name: "restricted-tcn" go_name: "RestrictedTCN" type_boolean { }
          description: "Disable topology change notification (TCN) sending on a port, used by STP to forbid network topology changes to propagate. This property only has an effect when protocol-mode is set to mstp."
        }
        property {
          name: "tag-stacking" type_boolean { }
          description: "Forces all packets to be treated as untagged packets. Packets on ingress port will be tagged with another VLAN tag regardless if a VLAN tag already exists, packets will be tagged with a VLAN ID that matches the pvid value and will use EtherType that is specified in ether-type. This property only has effect when vlan-filtering is set to yes."
        }
        property {
          name: "trusted" type_boolean { }
          description: "When enabled, it allows forwarding DHCP packets towards the DHCP server through this port. Mainly used to limit unauthorized servers to provide malicious information for users. This property only has an effect when dhcp-snooping is set to yes."
        }
        property {
          name: "unknown-multicast-flood" type_boolean { }
          description: "Changes the multicast flood option on bridge port, only controls the egress traffic. When enabled, the bridge allows flooding multicast packets to the specified bridge port, but when disabled, the bridge restricts multicast traffic from being flooded to the specified bridge port. The setting affects all multicast traffic, this includes non-IP, IPv4, IPv6 and the link-local multicast ranges (e.g. 224.0.0.0/24 and ff02::1). Note that when igmp-snooping is enabled and IGMP/MLD querier is detected, the bridge will automatically restrict unknown IP multicast from being flooded, so the setting is not mandatory for IGMP/MLD snooping setups. When using this setting together with igmp-snooping, the only multicast traffic that is allowed on the bridge port is the known multicast from the MDB table. "
        }
        property {
          name: "unknown-unicast-flood" type_boolean { }
          description: "Changes the unknown unicast flood option on bridge port, only controls the egress traffic. When enabled, the bridge allows flooding unknown unicast packets to the specified bridge port, but when disabled, the bridge restricts unknown unicast traffic from being flooded to the specified bridge port. If a MAC address is not learned in the host table, then the traffic is considered as unknown unicast traffic and will be flooded to all ports. MAC address is learned as soon as a packet on a bridge port is received and the source MAC address is added to the bridge host table. Since it is required for the bridge to receive at least one packet on the bridge port to learn the MAC address, it is recommended to use static bridge host entries to avoid packets being dropped until the MAC address has been learned."
        }
      }
    }
  }
}
