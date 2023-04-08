// Package ipsearch provides a simple way to search for IP addresses in a
package ipsearch

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})
}

// IPRange represents an IPv4 CIDR range.
type IPRange struct {
	rangeType RangeType
	firstSeg  uint8
	start     uint32
	end       uint32
	cidr      string
	country   string
}

// NewIPRange creates a new IPRange.
func NewIPRange(line string, rangeType RangeType) *IPRange {
	switch rangeType {
	case CIDR:
		return NewIPCIDR(line)
	case Geo:
		return NewIPGeo(line)
	}
	return nil
}

// NewIPRangeSlice creates a new slice of IPRange.
func NewIPRangeSlice(lines []string, rangeType RangeType) []*IPRange {
	ipRanges := make([]*IPRange, len(lines))
	for i, line := range lines {
		ipRanges[i] = NewIPRange(line, rangeType)
	}
	return ipRanges
}

// NewIPCIDR creates a new IPv4 CIDR range.
func NewIPCIDR(cidr string) *IPRange {
	start, end := IPCIDRRange(cidr)
	return &IPRange{
		rangeType: CIDR,
		firstSeg:  GetIPSegment(cidr, 1),
		start:     start,
		end:       end,
		cidr:      cidr,
	}
}

// NewIPGeo creates a new IPv4 CIDR range with a country code.
func NewIPGeo(csv string) *IPRange {
	fields := strings.Split(csv, ",")
	firstSegStart := GetIPSegment(fields[0], 1)
	firstSegEnd := GetIPSegment(fields[1], 1)
	if firstSegStart != firstSegEnd {
		log.Debugf("First segment of IP range is not the same: %s", csv)
	}
	start := IPStrToInt(fields[0])
	end := IPStrToInt(fields[1])
	country := fields[2]
	return &IPRange{
		rangeType: Geo,
		firstSeg:  firstSegStart,
		start:     start,
		end:       end,
		country:   country,
	}
}

func (ip *IPRange) String() string {
	switch ip.rangeType {
	case CIDR:
		return ip.cidr
	case Geo:
		return IPIntToStr(ip.start) + "," + IPIntToStr(ip.end) + "," + ip.country
	}
	return "Bad IPRange Type"
}

// Range return the range of the IPs in string format.
func (ip *IPRange) Range() string {
	return IPIntToStr(ip.start) + " - " + IPIntToStr(ip.end)
}

// Country returns the country code of the IP range.
func (ip *IPRange) Country() string {
	return ip.country
}

// CIDR returns the CIDR of the IP range.
func (ip *IPRange) CIDR() string {
	return ip.cidr
}

// Type returns the type of the IP range.
func (ip *IPRange) Type() RangeType {
	return ip.rangeType
}

// Split splits the IPRange into multiple IPRanges if the first segment is different.
func (ip *IPRange) Split() []*IPRange {
	endFirstSeg := GetIPSegment(IPIntToStr(ip.end), 1)
	if ip.firstSeg == endFirstSeg {
		return []*IPRange{ip}
	}
	ipRanges := make([]*IPRange, 0)

	for i := ip.firstSeg; i <= endFirstSeg; i++ {
		var start, end uint32
		start = uint32(i) << 24
		end = start + 0x00FFFFFF
		if start < ip.start {
			start = ip.start
		}
		if end > ip.end {
			end = ip.end
		}
		ipRanges = append(ipRanges, &IPRange{
			rangeType: ip.rangeType,
			firstSeg:  i,
			start:     start,
			end:       end,
			cidr:      fmt.Sprintf("%d.0.0.0/8", i),
			country:   ip.country,
		})
	}
	return ipRanges
}
