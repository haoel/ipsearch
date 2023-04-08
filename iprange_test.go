package ipsearch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/haoel/ipsearch"
)

var cidrs = []string{
	"1.4.1.0/24",
	"1.0.1.0/24",
	"1.0.2.0/23",
	"36.0.16.0/20",
	"43.224.242.0/24",
	"59.83.0.0/18",
	"103.196.64.0/22",
	"101.236.0.0/14",
	"45.119.116.0/22",
}

func TestCIDRIPRange(t *testing.T) {
	cidr_ips := ipsearch.NewIPRangeSlice(cidrs, ipsearch.CIDR)
	assert.Equal(t, len(cidr_ips), len(cidrs))
	for i, cidr_ip := range cidr_ips {
		assert.Equal(t, ipsearch.CIDR, cidr_ip.Type())
		assert.Equal(t, cidrs[i], cidr_ip.CIDR())
		assert.Equal(t, cidrs[i], cidr_ip.String())
		assert.Empty(t, cidr_ip.Country())
	}
}

var geo = []string{
	"1.0.64.0,1.0.127.255,JP",
	"1.0.32.0,1.0.63.255,CN",
	"1.0.128.0,1.0.255.255,TH",
	"185.123.184.0,185.123.187.255,BY",
	"185.123.192.0,185.123.195.255,RU",
	"103.148.242.0,103.148.243.255,ID",
	"103.148.244.0,103.148.245.255,HK",
	"2.56.172.0,2.56.179.255,CY",
	"2.56.184.0,2.56.187.255,LT",
	"2.56.180.0,2.56.183.255,RU",
}

func TestGeoIPRange(t *testing.T) {
	geo_ips := ipsearch.NewIPRangeSlice(geo, ipsearch.Geo)
	assert.Equal(t, len(geo_ips), len(geo))
	for i, geo_ip := range geo_ips {
		assert.Equal(t, ipsearch.Geo, geo_ip.Type())
		assert.Equal(t, geo[i], geo_ip.String())
		assert.NotEmpty(t, geo_ip.Country())
	}
}

func TestIPRangeSplit(t *testing.T) {
	ip := ipsearch.NewIPRange("3.0.0.0,4.255.255.255,US", ipsearch.Geo)
	ips := ip.Split()
	assert.Equal(t, 2, len(ips))
	assert.Equal(t, "3.0.0.0 - 3.255.255.255", ips[0].Range())
	assert.Equal(t, "4.0.0.0 - 4.255.255.255", ips[1].Range())

	ip = ipsearch.NewIPRange("6.0.0.0,8.127.255.255,US", ipsearch.Geo)
	ips = ip.Split()
	assert.Equal(t, 3, len(ips))
	assert.Equal(t, "6.0.0.0 - 6.255.255.255", ips[0].Range())
	assert.Equal(t, "7.0.0.0 - 7.255.255.255", ips[1].Range())
	assert.Equal(t, "8.0.0.0 - 8.127.255.255", ips[2].Range())

	ip = ipsearch.NewIPRange("67.231.224.0,68.65.215.255,US", ipsearch.Geo)
	ips = ip.Split()
	assert.Equal(t, 2, len(ips))
	assert.Equal(t, "67.231.224.0 - 67.255.255.255", ips[0].Range())
	assert.Equal(t, "68.0.0.0 - 68.65.215.255", ips[1].Range())

	ip = ipsearch.NewIPRange("1.1.1.2,1.1.1.4,TEST", ipsearch.Geo)
	ips = ip.Split()
	assert.Equal(t, 1, len(ips))
	assert.Equal(t, "1.1.1.2 - 1.1.1.4", ips[0].Range())
}
