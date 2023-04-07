package ipsearch

import "fmt"

const (
	ipFmt = "%d.%d.%d.%d"
)

// IPStrToInt converts a string IP address to an integer.
func IPStrToInt(ipStr string) uint32 {
	var (
		ip1, ip2, ip3, ip4 uint32
	)
	fmt.Sscanf(ipStr, ipFmt, &ip1, &ip2, &ip3, &ip4)
	return ip1<<24 | ip2<<16 | ip3<<8 | ip4
}

// IPIntToStr converts an integer IP address to a string.
func IPIntToStr(ip uint32) string {
	return fmt.Sprintf(ipFmt, ip>>24, ip>>16&0xff, ip>>8&0xff, ip&0xff)
}

// IPCIDRRange returns the start and end IP addresses for a CIDR range.
func IPCIDRRange(cidr string) (uint32, uint32) {
	var (
		ip1, ip2, ip3, ip4, mask uint32
	)
	fmt.Sscanf(cidr, "%d.%d.%d.%d/%d", &ip1, &ip2, &ip3, &ip4, &mask)
	start := ip1<<24 | ip2<<16 | ip3<<8 | ip4
	end := start | (1<<(32-mask) - 1)
	return start, end
}

// IPInRange checks if an IP address is in a range.
func ipInRange(ip uint32, start uint32, end uint32) bool {
	return ip >= start && ip <= end
}

// IPInCIDR checks if an IP address is in a CIDR range.
func IPInCIDR(ip string, cidr string) bool {
	start, end := IPCIDRRange(cidr)
	i := IPStrToInt(ip)
	return ipInRange(i, start, end)
}

// GetIPSegment returns the segment of an IP address.
func GetIPSegment(ip string, segment int) uint8 {
	ips := []uint8{0, 0, 0, 0}
	fmt.Sscanf(ip, ipFmt, &ips[0], &ips[1], &ips[2], &ips[3])
	return ips[segment-1]
}
