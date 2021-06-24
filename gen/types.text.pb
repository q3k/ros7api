sub {
  name: "interface"
  sub {
    name: "bridge"
    sub {
      name: "vlan"
      record {
        property { name: "bridge" type_string { } }
        property { name: "disabled" type_boolean { } }
        property { name: "tagged" type_string_list { } }
        property { name: "untagged" type_string_list { } }
        property { name: "vlan-ids" go_name: "VlanIDs" type_number { } }
        property { name: "current-tagged" read_only: true type_string_list { } }
        property { name: "current-untagged" read_only: true type_string_list { } }
        property { name: "dynamic" read_only: true type_boolean { } }
      }
    }
  }
}
