sub {
  name: "interface"
  sub {
    name: "bridge"
    sub {
      # https://help.mikrotik.com/docs/display/ROS/Bridge#Bridge-BridgeVLANtable
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
  }
}
