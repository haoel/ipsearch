package ipsearch

// RangeType is the type of file
type RangeType int

const (
	// CIDR is a file that contains CIDR ranges
	CIDR RangeType = iota
	// Geo is a file that contains GeoIP ranges and the country code as the CSV format
	Geo
)

// IPSearch is a struct that contains a map of IP ranges.
type IPSearch struct {
	container IPRangeMapList
}

// NewIPSearch creates a new IPSearch struct.
func NewIPSearch(lines []string, rangeType RangeType) *IPSearch {
	m := NewIPRangeMapList()
	ipRanges := NewIPRangeSlice(lines, rangeType)
	m.AppendBatch(ipRanges)
	m.Sort()
	return &IPSearch{container: m}
}

// NewIPSearchWithFile creates a new IPSearch struct from a file.
func NewIPSearchWithFile(path string, rangeType RangeType) (*IPSearch, error) {
	lines, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewIPSearch(lines, rangeType), nil
}

// NewIPSearchWithFileFromURL creates a new IPSearch struct from a URL.
func NewIPSearchWithFileFromURL(url string, fileType RangeType) (*IPSearch, error) {
	lines, err := ReadFileFromURL(url)
	if err != nil {
		return nil, err
	}
	return NewIPSearch(lines, fileType), nil
}

// Search search if an IP address is in the map of lists of IPv4 ranges.
func (s *IPSearch) Search(ip string) *IPRange {
	return s.container.Search(ip)
}
