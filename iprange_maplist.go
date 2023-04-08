package ipsearch

// IPRangeMapList is a map of lists of IPv4 CIDR ranges.
type IPRangeMapList map[uint8]*IPRangeList

// NewIPRangeMapList creates a new map of lists of IPv4 CIDR ranges.
func NewIPRangeMapList() IPRangeMapList {
	m := make(IPRangeMapList)
	return m
}

// AppendBatch adds a list of IPv4 CIDR ranges to the map of lists of IPv4 CIDR ranges.
func (m IPRangeMapList) AppendBatch(ipRanges []*IPRange) {
	for _, ip := range ipRanges {
		m.Append(ip)
	}
}

// Append adds an IPv4 CIDR range to the map of lists of IPv4 CIDR ranges.
func (m IPRangeMapList) Append(ipRange *IPRange) {
	ip1 := ipRange.firstSeg
	if _, ok := m[ip1]; !ok {
		m[ip1] = &IPRangeList{}
	}
	(*m[ip1]).Append(ipRange)
}

// Sort sorts the map of lists of IPv4 CIDR ranges.
func (m IPRangeMapList) Sort() {
	for _, list := range m {
		list.Sort()
	}
}

// InsertSorted inserts a IPv4 CIDR range to the map of lists of IPv4 CIDR ranges, keeping the lists sorted.
func (m IPRangeMapList) InsertSorted(ipRange *IPRange) {
	ip1 := ipRange.firstSeg
	if _, ok := m[ip1]; !ok {
		m[ip1] = &IPRangeList{}
	}
	(*m[ip1]).InsertSorted(ipRange)
}

// InsertSortedCIDRs inserts a list of IPv4 CIDR ranges to the map of lists of IPv4 CIDR ranges, keeping the lists sorted.
func (m IPRangeMapList) InsertSortedCIDRs(ipRanges []*IPRange) {
	for _, ip := range ipRanges {
		m.InsertSorted(ip)
	}
}

// Search search if an IP address is in the map of lists.
func (m IPRangeMapList) Search(ipStr string) *IPRange {
	ip1 := GetIPSegment(ipStr, 1)
	if _, ok := m[ip1]; !ok {
		return nil
	}
	return (*m[ip1]).Search(ipStr)
}
