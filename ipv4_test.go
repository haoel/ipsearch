package ipsearch_test

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/haoel/ipsearch"
	"github.com/stretchr/testify/assert"
)

func ipToInt(ipStr string) uint32 {
	var ip net.IP
	ip = net.ParseIP(ipStr)
	return binary.BigEndian.Uint32(ip.To4())
}

func intToIP(ipInt uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	return ip.String()
}

func TestIP(t *testing.T) {
	ipStr := "192.168.1.100"
	ip := ipsearch.IPStrToInt(ipStr)
	assert.Equal(t, ip, ipToInt(ipStr))
	assert.Equal(t, ipsearch.IPIntToStr(ip), intToIP(ip))

	ipStr = "172.20.1.1"
	ip = ipsearch.IPStrToInt(ipStr)
	assert.Equal(t, ip, ipToInt(ipStr))
	assert.Equal(t, ipsearch.IPIntToStr(ip), intToIP(ip))

	ipStr = "1.1.1.1"
	ip = ipsearch.IPStrToInt(ipStr)
	assert.Equal(t, ip, ipToInt(ipStr))
	assert.Equal(t, ipsearch.IPIntToStr(ip), intToIP(ip))
}

func TestCIDR(t *testing.T) {
	ipCIDR := "192.168.1.0/24"
	start, end := ipsearch.IPCIDRRange(ipCIDR)

	assert.Equal(t, start, ipToInt("192.168.1.0"))
	assert.Equal(t, end, ipToInt("192.168.1.255"))

	ip := "192.168.1.1"
	assert.True(t, ipsearch.IPInCIDR(ip, ipCIDR))
	assert.False(t, ipsearch.IPInCIDR("192.168.10.10", ipCIDR))

	ipCIDR = "1.1.1.1"
	start, end = ipsearch.IPCIDRRange(ipCIDR)
	assert.Equal(t, ipToInt("1.1.1.1"), start)
	assert.Equal(t, ipToInt("1.1.1.1"), end)

}

func TestIPSegment(t *testing.T) {
	ip := "192.168.1.1"
	assert.Equal(t, ipsearch.GetIPSegment(ip, 1), uint8(192))

	ip = "172.20.1.1"
	assert.Equal(t, ipsearch.GetIPSegment(ip, 1), uint8(172))

	ip = "0.0.0.0"
	assert.Equal(t, ipsearch.GetIPSegment(ip, 1), uint8(0))

	ip = "255.255.255.0"
	assert.Equal(t, ipsearch.GetIPSegment(ip, 1), uint8(255))

	ip = "256.255.255.0"
	assert.Equal(t, ipsearch.GetIPSegment(ip, 1), uint8(0))
}
