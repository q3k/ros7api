//go:generate protoc --go_out=. --go_opt=module=github.com/q3k/ros7api gen/kinds.proto
//go:generate go build -o gen.elf github.com/q3k/ros7api/gen
//go:generate ./gen.elf

package ros7api
