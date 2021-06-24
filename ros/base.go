package ros

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

// RecordID is the ID of a ROS record, eg. '*13'.
type RecordID string

type Record struct {
	ID RecordID `json:".id"`
}

// StringPtr returns a string pointer for use in _Update structs.
func StringPtr(s string) *string {
	return &s
}

// Number is a ROS number. We assume they all fit in int64 (this is
// undocumented...), and we (de)serialize them as strings.
type Number int64

// NumberPtr returns a Number pointer for use in _Update structs.
func NumberPtr(n int64) *Number {
	v := Number(n)
	return &v
}

func (n *Number) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var v int64
	var err error
	if strings.HasPrefix(s, "0x") {
		v, err = strconv.ParseInt(s[2:], 16, 64)
	} else {
		v, err = strconv.ParseInt(s, 10, 64)
	}
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", s, err)
	}
	*n = Number(v)
	return nil
}

func (n *Number) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, *n)), nil
}

// Boolean is ROS boolean, (de)serialized as a string.
type Boolean bool

// BooleanPtr returns a Boolean pointer for use in _Update structs.
func BooleanPtr(b bool) *Boolean {
	v := Boolean(b)
	return &v
}

func (n *Boolean) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "true":
		*n = Boolean(true)
	case "false":
		*n = Boolean(false)
	default:
		return fmt.Errorf("invalid boolean: %q", s)
	}
	return nil
}

func (n *Boolean) MarshalJSON() ([]byte, error) {
	if *n {
		return []byte(`"true"`), nil
	} else {
		return []byte(`"false"`), nil
	}
}

// IPNet is a ROS 'Address/Netmask' type, IPv4 or IPv6, serialized into CIDR
// notation.
type IPNet struct {
	Address net.IP
	Network net.IPNet
}

func (n *IPNet) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ip, net, err := net.ParseCIDR(s)
	if err != nil {
		return fmt.Errorf("invalid CIDR %q: %w", s, err)
	}
	*n = IPNet{
		Address: ip,
		Network: *net,
	}
	return nil
}

func (n *IPNet) MarshalJSON() ([]byte, error) {
	ones, _ := n.Network.Mask.Size()
	v := fmt.Sprintf(`"%s/%d"`, n.Address.String(), ones)
	return []byte(v), nil
}

// IP is a ROS 'Address' type, IPv4 or IPv6, serialized into a dot/colon
// notation.
type IP net.IP

// IPPtr returns a pointer to IP, for use in _Update structs.
func IPPtr(i net.IP) *IP {
	v := IP(i)
	return &v
}

func (n *IP) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("invalid IP %q", s)
	}
	*n = IP(ip)
	return nil
}

func (n *IP) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", net.IP(*n).String())), nil
}

// StringList is a ROS7 list of strings, eg. interfaces, (de)serialized as a
// string containing comma-delimited values.
type StringList []string

// StringListPtr returns a pointer to StringList, for use in _Update structs.
func StringListPtr(l ...string) *StringList {
	v := StringList(l)
	return &v
}

func (n *StringList) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parts := strings.Split(s, ",")
	*n = parts
	return nil
}

func (n *StringList) MarshalJSON() ([]byte, error) {
	for i, el := range *n {
		if strings.Contains(el, ",") || strings.Contains(el, `"`) {
			return nil, fmt.Errorf("element %d of string list (%q) contains invalid character", i, el)
		}
	}
	res := fmt.Sprintf("%q", strings.Join(*n, ","))
	return []byte(res), nil
}

// NumberList is a ROS7 list-of-numbers-that-might-be-ranges. Eg., it supports
// the following data:
//
//   foo=123
//   foo=123,234
//   foo=123-150,234
//
// All the ranges are inclusive.
//
// When serializing an unmodified deserialized struct, the serialized
// representation is guaranteed to be the same as the initial representation.
// However, if modified, the ranges / values might be optimized and shuffled
// around.
type NumberList struct {
	ranges []numberListRange
}

type numberListRange struct {
	lower int64
	upper int64
}

// ParseNumberList parses a ROS-style representation of a list of numbers, eg.
// 100,105-110,200. The ranges are inclusive.
func ParseNumberList(s string) (*NumberList, error) {
	parts := strings.Split(s, ",")
	var ranges []numberListRange
	for _, part := range parts {
		var lower int64
		var upper int64
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range %q", part)
			}
			n, err := strconv.ParseInt(rangeParts[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid range %q: %w", part, err)
			}
			lower = n
			n, err = strconv.ParseInt(rangeParts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid range %q: %w", part, err)
			}
			upper = n
			if upper <= lower {
				return nil, fmt.Errorf("invalid range %q", part)
			}
		} else {
			n, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number %q: %w", part, err)
			}
			upper = n
			lower = n
		}
		ranges = append(ranges, numberListRange{lower, upper})
	}
	return &NumberList{ranges}, nil
}

func (n *NumberList) UnmarshalJSON(b []byte) error {
	n.ranges = nil

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := ParseNumberList(s)
	if err != nil {
		return err
	}
	*n = *parsed
	return nil
}

// optimize sorts and coalesces internal numberListRanges.
func (n *NumberList) optimize() {
	if len(n.ranges) == 0 {
		return
	}

	// Sort by lower bound.
	sort.Slice(n.ranges, func(i, j int) bool {
		return n.ranges[i].lower < n.ranges[j].lower
	})

	var res []numberListRange

	lower := n.ranges[0].lower
	upper := n.ranges[0].upper
	for _, el := range n.ranges[1:] {
		if el.lower <= upper+1 {
			// Range overlaps with current.
			if el.upper > upper {
				// Range extends current.
				upper = el.upper
			} else {
				// Range fully contained within current.
			}
		} else {
			// Disjoint range.
			res = append(res, numberListRange{lower, upper})
			lower = el.lower
			upper = el.upper
		}
	}
	res = append(res, numberListRange{lower, upper})

	n.ranges = res
}

// Add a number to this list.
func (n *NumberList) Add(v int64) {
	n.ranges = append(n.ranges, numberListRange{v, v})
	n.optimize()
}

// Add a range of numbers (inclusive) to this list.
func (n *NumberList) AddRange(lower, upper int64) error {
	if upper < lower {
		return fmt.Errorf("invalid range %d-%d", lower, upper)
	}
	n.ranges = append(n.ranges, numberListRange{lower, upper})
	n.optimize()
	return nil
}

// Remove a number from this list if present.
func (n *NumberList) Remove(v int64) {
	var res []numberListRange
	for _, r := range n.ranges {
		if r.lower <= v && r.upper >= v {
			l1 := r.lower
			u1 := v - 1
			l2 := v + 1
			u2 := r.upper

			if l1 <= u1 {
				res = append(res, numberListRange{l1, u1})
			}
			if l2 <= u2 {
				res = append(res, numberListRange{l2, u2})
			}
		} else {
			res = append(res, r)
		}
	}
	n.ranges = res
	n.optimize()
}

// Contains returns whether this list contains a given number.
func (n *NumberList) Contains(v int64) bool {
	for _, r := range n.ranges {
		if r.lower <= v && r.upper >= v {
			return true
		}
	}
	return false
}

func (n *NumberList) String() string {
	var parts []string
	for _, r := range n.ranges {
		var s string
		if r.lower == r.upper {
			s = fmt.Sprintf("%d", r.lower)
		} else {
			s = fmt.Sprintf("%d-%d", r.lower, r.upper)
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ",")
}

func (n *NumberList) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", n.String())), nil
}
