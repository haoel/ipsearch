package ipsearch

// IPSearch is a struct that contains a map of IP ranges.
type IPSearch struct {
	container IPCIDRMapList
}

// NewIPSearch creates a new IPSearch struct.
func NewIPSearch(cidrs []string) *IPSearch {
	m := NewIPCIDRMapList()
	m.AppendCIDRs(cidrs)
	m.Sort()
	return &IPSearch{container: m}
}

// NewIPSearchWithFile creates a new IPSearch struct from a file.
func NewIPSearchWithFile(path string) (*IPSearch, error) {
	lines, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewIPSearch(lines), nil
}

// NewIPSearchWithFileFromURL creates a new IPSearch struct from a URL.
func NewIPSearchWithFileFromURL(url string) (*IPSearch, error) {
	lines, err := ReadFileFromURL(url)
	if err != nil {
		return nil, err
	}
	return NewIPSearch(lines), nil
}

// Check checks if an IP address is in the map of lists of IPv4 CIDR ranges.
func (s *IPSearch) Check(ipStr string) bool {
	return s.container.Check(ipStr)
}
