package ipsearch_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/haoel/ipsearch"
)

const expectedPart1 = `1.0.1.0/24
1.0.2.0/23
1.4.1.0/24
36.0.16.0/20
43.224.242.0/24
45.119.116.0/22`

const expectedPart2 = `59.83.0.0/18
101.236.0.0/14
103.196.64.0/22
`
const expected = expectedPart1 + "\n" + expectedPart2

type testCIDRData struct {
	ip       string
	cidr     string
	rangeStr string
	find     bool
}

var testCIDRDataList = []testCIDRData{
	{"1.0.1.24", "1.0.1.0/24", "1.0.1.0 - 1.0.1.255", true},
	{"1.4.1.1", "1.4.1.0/24", "1.0.1.0 - 1.0.1.255", true},
	{"43.224.242.100", "43.224.242.0/24", "1.0.1.0 - 1.0.1.255", true},
	{"101.236.0.1", "101.236.0.0/14", "1.0.1.0 - 1.0.1.255", true},
	{"101.240.0.1", "", "", false},
	{"5.5.5.5", "", "", false},
}

func TestIPCIDRList(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ipCIDRList := ipsearch.NewIPRangeList(cidrs, ipsearch.CIDR)

	assert.Equal(t, ipCIDRList.Len(), len(cidrs))

	ipCIDRList.Sort()
	assert.Equal(t, ipCIDRList.String(), expected)

	newIP := "58.83.0.0/16"
	ipCIDRList.InsertSorted(ipsearch.NewIPCIDR(newIP))
	newExpected := expectedPart1 + "\n" + newIP + "\n" + expectedPart2

	assert.Equal(t, ipCIDRList.String(), newExpected)

	for _, data := range testCIDRDataList {
		ip := ipCIDRList.Search(data.ip)
		assert.Equal(t, ip != nil, data.find)
		if ip != nil {
			assert.Equal(t, ip.String(), data.cidr)
		}
	}
}

const expectedGeo = `1.0.32.0,1.0.63.255,CN
1.0.64.0,1.0.127.255,JP
1.0.128.0,1.0.255.255,TH
2.56.172.0,2.56.179.255,CY
2.56.180.0,2.56.183.255,RU
2.56.184.0,2.56.187.255,LT
103.148.242.0,103.148.243.255,ID
103.148.244.0,103.148.245.255,HK
185.123.184.0,185.123.187.255,BY
185.123.192.0,185.123.195.255,RU
`

const expectedGeoRange = `1.0.32.0 - 1.0.63.255
1.0.64.0 - 1.0.127.255
1.0.128.0 - 1.0.255.255
2.56.172.0 - 2.56.179.255
2.56.180.0 - 2.56.183.255
2.56.184.0 - 2.56.187.255
103.148.242.0 - 103.148.243.255
103.148.244.0 - 103.148.245.255
185.123.184.0 - 185.123.187.255
185.123.192.0 - 185.123.195.255
`

type testGeoData struct {
	ip   string
	geo  string
	find bool
}

var testGeoDataList = []testGeoData{
	{"1.0.35.10", "CN", true},
	{"1.0.110.10", "JP", true},
	{"103.148.243.10", "ID", true},
	{"103.148.244.10", "HK", true},
	{"1.1.1.1", "", false},
	{"8.8.8.8", "", false},
}

func TestGeoList(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	geoList := ipsearch.NewIPRangeList(geo, ipsearch.Geo)
	geoList.Sort()

	assert.Equal(t, expectedGeo, geoList.String())
	assert.Equal(t, expectedGeoRange, geoList.Range())

	for _, data := range testGeoDataList {
		ip := geoList.Search(data.ip)
		assert.Equal(t, ip != nil, data.find)
		if ip != nil {
			assert.Equal(t, ip.Country(), data.geo)
		}
	}
}
