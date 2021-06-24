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

type menu struct {
	path string
	buf  bytes.Buffer
	m    *kpb.Menu
	sub  map[string]*menu
}

type property struct {
	p      *kpb.Property
	name   string
	goname string
	gotype string
}

func propertyFromProto(p *kpb.Property) *property {
	goname := p.GoName
	if goname == "" {
		parts := strings.Split(p.Name, "-")
		for i, p := range parts {
			parts[i] = strings.Title(p)
		}
		goname = strings.Join(parts, "")
	}
	gotype := ""
	switch p.Type.(type) {
	case *kpb.Property_TypeNumber:
		gotype = "Number"
	case *kpb.Property_TypeString:
		gotype = "string"
	case *kpb.Property_TypeBoolean:
		gotype = "Boolean"
	case *kpb.Property_TypeStringList:
		gotype = "StringList"
	default:
		panic(fmt.Sprintf("unknown type %v", p.Type))
	}

	return &property{
		p:      p,
		name:   p.Name,
		goname: goname,
		gotype: gotype,
	}
}

func (m *menu) printf(format string, args ...interface{}) {
	fmt.Fprintf(&m.buf, format, args...)
}

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

	var properties []*property
	for _, p := range m.m.Record.Property {
		prop := propertyFromProto(p)
		properties = append(properties, prop)
	}

	nameParts := strings.Split(m.path, "/")
	for i, p := range nameParts {
		nameParts[i] = strings.Title(p)
	}
	sname := strings.Join(nameParts, "")

	m.printf("type %s struct {\n", sname)
	m.printf("\tRecord\n\n")
	for _, p := range properties {
		m.printf("\t%s\t%s\t`json:\"%s\"`\n", p.goname, p.gotype, p.name)
	}
	m.printf("}\n\n")

	m.printf("type %s_Update struct {\n", sname)
	for _, p := range properties {
		if p.p.ReadOnly {
			continue
		}
		m.printf("\t%s\t*%s\t`json:\"%s,omitempty\"`\n", p.goname, p.gotype, p.name)
	}
	m.printf("}\n\n")

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
