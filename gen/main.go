// package main implements a ROS7 REST API client generator. It parses a
// prototext containing definition of record types in ROS (eg. the Mikrotik
// confluence wiki) and spits out Go files in ros/zz_*.go.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path"
	"strings"

	kpb "github.com/q3k/ros7api/gen/kinds"
	"google.golang.org/protobuf/encoding/prototext"
)

var (
	flagTypesPath string
)

// menu is an element of the ROS menu tree.
type menu struct {
	// path is a menu path like interface/vlan
	path string
	// m is the proto definition of this menu element.
	m *kpb.Menu
	// sub is a math from sub-menu name to submenu.
	sub map[string]*menu

	// buf is the code generation buffer for this node.
	buf bytes.Buffer
}

// property is a ROS record property parsed from protobuf.
type property struct {
	// p is the proto definition for this property.
	p *kpb.Property
	// name is the 'native' ROS name of this property, eg. vlan-ids.
	name string
	// goname is the Go name of this proprety, eg. VlanIDs.
	goname string
	// gotype is the REST Client Go type of this property.
	gotype string
	enum   *kpb.TypeEnum
}

func propertyFromProto(p *kpb.Property, sname string) *property {
	goname := p.GoName
	if goname == "" {
		goname = goify(p.Name)
	}

	var enum *kpb.TypeEnum

	gotype := ""
	switch v := p.Type.(type) {
	case *kpb.Property_TypeNumber:
		gotype = "Number"
	case *kpb.Property_TypeString:
		gotype = "string"
	case *kpb.Property_TypeBoolean:
		gotype = "Boolean"
	case *kpb.Property_TypeStringList:
		gotype = "StringList"
	case *kpb.Property_TypeNumberList:
		gotype = "NumberList"
	case *kpb.Property_TypeEnum:
		gotype = fmt.Sprintf("%s_%s", sname, goname)
		enum = v.TypeEnum
	default:
		panic(fmt.Sprintf("unknown type %v", p.Type))
	}

	return &property{
		p:      p,
		name:   p.Name,
		goname: goname,
		gotype: gotype,
		enum:   enum,
	}
}

// printf writes a line to the menu's code generation buffer.
func (m *menu) printf(format string, args ...interface{}) {
	fmt.Fprintf(&m.buf, format, args...)
}

func goify(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		parts[i] = strings.Title(p)
	}
	return strings.Join(parts, "")
}

// generate the code for this menu element. Currently a single menu element
// corresponds to a single Go source file.
func (m *menu) generate() error {
	m.buf.Reset()
	m.printf("package ros\n\n")
	m.printf("import (\n")
	m.printf("\t\"context\"\n")
	m.printf("\t\"encoding/json\"\n")
	m.printf("\t\"fmt\"\n")
	m.printf(")\n\n")
	m.printf("// Automatically generated by github.com/q3k/ros7api/gen, do not edit.\n")
	m.printf("\n")

	// Turn path into record struct name (eg. interface/bridge/vlan into
	// InterfaceBridgeVlan).
	nameParts := strings.Split(m.path, "/")
	for i, p := range nameParts {
		nameParts[i] = strings.Title(p)
	}
	sname := strings.Join(nameParts, "")

	// Parse properties.
	var properties []*property
	for _, p := range m.m.Record.Property {
		prop := propertyFromProto(p, sname)
		properties = append(properties, prop)
	}

	// Emit enums.
	for _, p := range properties {
		if p.enum == nil {
			continue
		}
		m.printf("type %s string\n\n", p.gotype)
		m.printf("const (\n")
		for _, variant := range p.enum.Variant {
			if variant.Description != "" {
				m.printf("\t// %s\n", variant.Description)
			}
			m.printf("\t%s%s = %q\n", p.gotype, goify(variant.Value), variant.Value)
		}
		m.printf(")\n")
	}

	// Emit record type.
	m.printf("// %s represents a ROS `%s` record, including read-only fields.\n", sname, m.path)
	if m.m.Record.Description != "" {
		m.printf("//\n")
		m.printf("// %s\n", m.m.Record.Description)
	}
	m.printf("type %s struct {\n", sname)
	m.printf("\tRecord\n\n")
	for _, p := range properties {
		if p.p.Description != "" {
			m.printf("\t// %s\n", p.p.Description)
		}
		m.printf("\t%s\t%s\t`json:\"%s\"`\n", p.goname, p.gotype, p.name)
	}
	m.printf("}\n\n")

	// Emit record update type.
	m.printf("// %s_Update is an update to a ROS `%s` record. Any unset field will not be updated.\n", sname, m.path)
	m.printf("type %s_Update struct {\n", sname)
	for _, p := range properties {
		if p.p.ReadOnly {
			continue
		}
		if p.p.Description != "" {
			m.printf("\t// %s\n", p.p.Description)
		}
		m.printf("\t%s\t*%s\t`json:\"%s,omitempty\"`\n", p.goname, p.gotype, p.name)
	}
	m.printf("}\n\n")

	m.printf("// %sList returns a list of all `%s` records.\n", sname, m.path)
	m.printf("func (c *Client) %sList(ctx context.Context) ([]%s, error) {\n", sname, sname)
	m.printf("\tbody, err := c.doGET(ctx, %q)\n", m.path)
	m.printf("\tif err != nil {\n")
	m.printf("\t\treturn nil, fmt.Errorf(\"could not GET: %%w\", err)\n")
	m.printf("\t}\n")
	m.printf("\tdefer body.Close()\n\n")
	m.printf("\tvar target []%s\n", sname)
	m.printf("\tif err := json.NewDecoder(body).Decode(&target); err != nil {\n")
	m.printf("\t\treturn nil, fmt.Errorf(\"could not decode JSON: %%w\", err)\n")
	m.printf("\t}\n")
	m.printf("\treturn target, nil\n")
	m.printf("}\n\n")

	m.printf("// %sPatch updates the given fields of a `%s` record by ID.\n", sname, m.path)
	m.printf("func (c *Client) %sPatch(ctx context.Context, id RecordID, u *%s_Update) (*%s, error) {\n", sname, sname, sname)
	m.printf("\trdata, err := json.Marshal(u)\n")
	m.printf("\tif err != nil {\n")
	m.printf("\t\treturn nil, fmt.Errorf(\"could not marshal update: %%w\", err)\n")
	m.printf("\t}\n")
	m.printf("\tbody, err := c.doPATCH(ctx, %q+string(id), rdata)\n", m.path+"/")
	m.printf("\tif err != nil {\n")
	m.printf("\t\treturn nil, fmt.Errorf(\"could not PATCH: %%w\", err)\n")
	m.printf("\t}\n")
	m.printf("\tdefer body.Close()\n\n")
	m.printf("\tvar target struct {\n")
	m.printf("\t\t%s\n", sname)
	m.printf("\t\tError int64 `json:\"error\"`\n")
	m.printf("\t\tMessage string `json:\"message\"`\n")
	m.printf("\t\tDetail string `json:\"detail\"`\n")
	m.printf("\t}\n")
	m.printf("\tif err := json.NewDecoder(body).Decode(&target); err != nil {\n")
	m.printf("\t\treturn nil, fmt.Errorf(\"could not decode JSON: %%w\", err)\n")
	m.printf("\t}\n")
	m.printf("\tif target.Error != 0 {\n")
	m.printf("\t\treturn nil, fmt.Errorf(\"server error: %%s: %%s\", target.Message, target.Detail)\n")
	m.printf("\t}\n")
	m.printf("\treturn &target.%s, nil\n", sname)
	m.printf("}\n\n")
	return nil
}

func (m *menu) writeGo(root string) error {
	if r := m.m.Record; r != nil {
		if err := m.generate(); err != nil {
			return fmt.Errorf("could not generate %s: %w", m.path, err)
		}
		pathParts := strings.Split(m.path, "/")
		path := path.Join(root, fmt.Sprintf("zz_%s.go", strings.Join(pathParts, "_")))
		log.Printf("Writing %s...", path)
		src, err := format.Source(m.buf.Bytes())
		if err != nil {
			return fmt.Errorf("could not format %s: %w", m.path, err)
		}
		if err := ioutil.WriteFile(path, src, 0644); err != nil {
			return fmt.Errorf("could not write %s: %w", m.path, err)
		}
	}
	for _, sub := range m.sub {
		if err := sub.writeGo(root); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.StringVar(&flagTypesPath, "types_path", "gen/types.text.pb", "Path to types prototext")
	flag.Parse()
	data, err := ioutil.ReadFile(flagTypesPath)
	if err != nil {
		log.Fatalf("Could not load types prototext: %v", err)
	}
	var m kpb.Menu
	if err := prototext.Unmarshal(data, &m); err != nil {
		log.Fatalf("Could not unmarshal types prototext: %v", err)
	}

	tree := recurse(&m, "")
	err = tree.writeGo("ros")
	if err != nil {
		panic(err)
	}
}

func recurse(m *kpb.Menu, path string) *menu {
	sub := make(map[string]*menu)
	for _, s := range m.Sub {
		var spath string
		if path == "" {
			spath = s.Name
		} else {
			spath = path + "/" + s.Name
		}
		sub[s.Name] = recurse(s, spath)
	}
	return &menu{
		path: path,
		m:    m,
		sub:  sub,
	}
}
