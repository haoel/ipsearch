package ipsearch

import (
	"sort"

	log "github.com/sirupsen/logrus"
)

// IPRangeList is a list of IPv4 CIDR ranges.
type IPRangeList []*IPRange

// NewIPRangeList creates a new list of IPv4 CIDR ranges.
func NewIPRangeList(lines []string, rangeType RangeType) IPRangeList {
	list := make(IPRangeList, 0)
	for _, line := range lines {
		list.Append(NewIPRange(line, rangeType))
	}

	if log.GetLevel() == log.DebugLevel {
		for _, ipRange := range list {
			log.Debugf("Ranges: %s", ipRange.Range())
		}
	}
	return list
}

// Append adds a IPv4 range to the list.
func (list *IPRangeList) Append(ip *IPRange) {
	ips := ip.Split()
	*list = append(*list, ips...)
}

// InsertSorted inserts a IPv4 range to the list, keeping the list sorted.
func (list *IPRangeList) InsertSorted(ip *IPRange) {
	ips := ip.Split()
	for _, ip := range ips {
		//using the binary search to search the index of list to insert the cidr
		start := 0
		end := len(*list) - 1
		for start <= end {
			mid := (start + end) / 2
			if (*list)[mid].start < ip.start {
				start = mid + 1
			} else {
				end = mid - 1
			}
		}
		idx := start
		//insert the cidr to the list
		*list = append(*list, nil)
		copy((*list)[idx+1:], (*list)[idx:])
		(*list)[idx] = ip
	}
}

// Sort sorts the list of IPv4 CIDR ranges.
func (list *IPRangeList) Sort() {
	sort.Slice(*list, func(i, j int) bool {
		return (*list)[i].start < (*list)[j].start
	})
}

// Len returns the length of the list of IPv4 CIDR ranges.
func (list *IPRangeList) Len() int {
	return len(*list)
}

// String returns a string representation of the list of IPv4 CIDR ranges.
func (list *IPRangeList) String() string {
	str := ""
	for _, cidr := range *list {
		str += cidr.String() + "\n"
	}
	return str
}

// Range returns a string representation of the list of IPv4 ranges.
func (list *IPRangeList) Range() string {
	str := ""
	for _, ipRagne := range *list {
		str += ipRagne.Range() + "\n"
	}
	return str
}

// Search search if an IP address is in the list.
func (list *IPRangeList) Search(ipStr string) *IPRange {
	ip := IPStrToInt(ipStr)
	// using the binary search to search  the list
	start := 0
	end := len(*list) - 1
	for start <= end {
		mid := (start + end) / 2
		if ipInRange(ip, (*list)[mid].start, (*list)[mid].end) {
			log.Debugf("IP %s is in Range %s", ipStr, (*list)[mid].Range())
			return (*list)[mid]
		}
		if ip < (*list)[mid].start {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	log.Debugf("IP %s is not in any following Ranges. \n%s", ipStr, list)
	return nil
}
