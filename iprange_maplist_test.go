package ipsearch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/haoel/ipsearch"
)

func TestIPCIDRMapList(t *testing.T) {

	ipCIDRMapList := ipsearch.NewIPRangeMapList()
	ipRangeList := ipsearch.NewIPRangeList(cidrs, ipsearch.CIDR)
	ipCIDRMapList.AppendBatch(ipRangeList)
	// without sort, it would be a bug
	ip := ipCIDRMapList.Search("1.4.1.2")
	assert.Nil(t, ip)
	// with sort, it would be ok
	ipCIDRMapList.Sort()
	ip = ipCIDRMapList.Search("1.4.1.3")
	assert.NotNil(t, ip)
	assert.Equal(t, ip.String(), "1.4.1.0/24")

	ipCIDRMapList = ipsearch.NewIPRangeMapList()
	ipRangeList = ipsearch.NewIPRangeList(cidrs, ipsearch.CIDR)
	ipCIDRMapList.InsertSortedCIDRs(ipRangeList)

	for _, data := range testCIDRDataList {
		ip := ipCIDRMapList.Search(data.ip)
		assert.Equal(t, ip != nil, data.find)
		if ip != nil {
			assert.Equal(t, ip.String(), data.cidr)
		}
	}

	ipCIDRMapList.InsertSorted(ipsearch.NewIPCIDR("1.3.0.0/16"))
	ip = ipCIDRMapList.Search("1.3.1.1")
	assert.NotNil(t, ip)
	assert.Equal(t, ip.String(), "1.3.0.0/16")
}
