// Package ipsearch provides a simple way to search for IP addresses in a
package ipsearch

import (
	"sort"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})
}

// IPCIDR represents an IPv4 CIDR range.
type IPCIDR struct {
	start uint32
	end   uint32
	str   string
}

// NewIPCIDR creates a new IPv4 CIDR range.
func NewIPCIDR(cidr string) *IPCIDR {
	start, end := IPCIDRRange(cidr)
	return &IPCIDR{
		start: start,
		end:   end,
		str:   cidr,
	}
}

func (cidr *IPCIDR) String() string {
	return cidr.str
}

// IPCIDRList is a list of IPv4 CIDR ranges.
type IPCIDRList []*IPCIDR

// NewIPCIDRList creates a new list of IPv4 CIDR ranges.
func NewIPCIDRList(cidrs []string) IPCIDRList {
	list := make(IPCIDRList, 0)
	for _, cidr := range cidrs {
		list.Append(cidr)
	}

	if log.GetLevel() == log.DebugLevel {
		for _, cidr := range list {
			log.Debugf("CIDR: %s [ %s - %s]", cidr, IPIntToStr(cidr.start), IPIntToStr(cidr.end))
		}
	}
	return list
}

// Append adds a IPv4 CIDR range to the list of IPv4 CIDR ranges.
func (list *IPCIDRList) Append(cidr string) {
	*list = append(*list, NewIPCIDR(cidr))
}

// InsertSorted inserts a IPv4 CIDR range to the list, keeping the list sorted.
func (list *IPCIDRList) InsertSorted(cidr string) {
	//using the binary search to search the index of list to insert the cidr
	start := 0
	end := len(*list) - 1
	for start <= end {
		mid := (start + end) / 2
		if (*list)[mid].start < IPStrToInt(cidr) {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	idx := start
	//insert the cidr to the list
	*list = append(*list, nil)
	copy((*list)[idx+1:], (*list)[idx:])
	(*list)[idx] = NewIPCIDR(cidr)
}

// Sort sorts the list of IPv4 CIDR ranges.
func (list *IPCIDRList) Sort() {
	sort.Slice(*list, func(i, j int) bool {
		return (*list)[i].start < (*list)[j].start
	})
}

// Len returns the length of the list of IPv4 CIDR ranges.
func (list *IPCIDRList) Len() int {
	return len(*list)
}

// String returns a string representation of the list of IPv4 CIDR ranges.
func (list *IPCIDRList) String() string {
	str := ""
	for _, cidr := range *list {
		str += cidr.String() + "\n"
	}
	return str
}

// Check checks if an IP address is in the list of IPv4 CIDR ranges.
func (list *IPCIDRList) Check(ipStr string) bool {
	ip := IPStrToInt(ipStr)
	// using the binary search to search  the list
	start := 0
	end := len(*list) - 1
	for start <= end {
		mid := (start + end) / 2
		if ipInRange(ip, (*list)[mid].start, (*list)[mid].end) {
			log.Debugf("IP %s is in CIDR %s", ipStr, (*list)[mid].str)
			return true
		}
		if ip < (*list)[mid].start {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	log.Debugf("IP %s is not in any following CIDR. \n%s", ipStr, list)
	return false
}

// IPCIDRMapList is a map of lists of IPv4 CIDR ranges.
type IPCIDRMapList map[uint8]*IPCIDRList

// NewIPCIDRMapList creates a new map of lists of IPv4 CIDR ranges.
func NewIPCIDRMapList() IPCIDRMapList {
	m := make(IPCIDRMapList)
	return m
}

// AppendCIDRs adds a list of IPv4 CIDR ranges to the map of lists of IPv4 CIDR ranges.
func (m IPCIDRMapList) AppendCIDRs(cidrs []string) {
	for _, cidr := range cidrs {
		m.AppendCIDR(cidr)
	}
}

// AppendCIDR adds an IPv4 CIDR range to the map of lists of IPv4 CIDR ranges.
func (m IPCIDRMapList) AppendCIDR(cidr string) {
	ip1 := GetIPSegment(cidr, 1)
	if _, ok := m[ip1]; !ok {
		m[ip1] = &IPCIDRList{}
	}
	(*m[ip1]).Append(cidr)
}

// Sort sorts the map of lists of IPv4 CIDR ranges.
func (m IPCIDRMapList) Sort() {
	for _, list := range m {
		list.Sort()
	}
}

// InsertSorted inserts a IPv4 CIDR range to the map of lists of IPv4 CIDR ranges, keeping the lists sorted.
func (m IPCIDRMapList) InsertSorted(cidr string) {
	ip1 := GetIPSegment(cidr, 1)
	if _, ok := m[ip1]; !ok {
		m[ip1] = &IPCIDRList{}
	}
	(*m[ip1]).InsertSorted(cidr)
}

// InsertSortedCIDRs inserts a list of IPv4 CIDR ranges to the map of lists of IPv4 CIDR ranges, keeping the lists sorted.
func (m IPCIDRMapList) InsertSortedCIDRs(cidrs []string) {
	for _, cidr := range cidrs {
		m.InsertSorted(cidr)
	}
}

// Check checks if an IP address is in the map of lists of IPv4 CIDR ranges.
func (m IPCIDRMapList) Check(ipStr string) bool {
	ip1 := GetIPSegment(ipStr, 1)
	if _, ok := m[ip1]; !ok {
		return false
	}
	return (*m[ip1]).Check(ipStr)
}
