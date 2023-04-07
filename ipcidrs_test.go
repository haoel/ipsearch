package ipsearch_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
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

func TestIPCIDRList(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ipCIDRList := ipsearch.NewIPCIDRList(cidrs)

	assert.Equal(t, ipCIDRList.Len(), len(cidrs))

	ipCIDRList.Sort()
	assert.Equal(t, ipCIDRList.String(), expected)

	newIP := "58.83.0.0/16"
	ipCIDRList.InsertSorted(newIP)
	newExpected := expectedPart1 + "\n" + newIP + "\n" + expectedPart2

	assert.Equal(t, ipCIDRList.String(), newExpected)

	assert.True(t, ipCIDRList.Check("1.0.1.24"))
	assert.True(t, ipCIDRList.Check("1.4.1.1"))
	assert.True(t, ipCIDRList.Check("43.224.242.100"))
	assert.True(t, ipCIDRList.Check("101.236.0.1"))
	assert.False(t, ipCIDRList.Check("101.240.0.1"))
}

func TestIPCIDRMapList(t *testing.T) {

	ipCIDRMapList := ipsearch.NewIPCIDRMapList()
	ipCIDRMapList.AppendCIDRs(cidrs)
	// without sort, it would be a bug
	assert.False(t, ipCIDRMapList.Check("1.4.1.2"))
	// with sort, it would be ok
	ipCIDRMapList.Sort()
	assert.True(t, ipCIDRMapList.Check("1.4.1.3"))

	ipCIDRMapList = ipsearch.NewIPCIDRMapList()
	ipCIDRMapList.InsertSortedCIDRs(cidrs)
	assert.True(t, ipCIDRMapList.Check("1.0.1.24"))
	assert.True(t, ipCIDRMapList.Check("1.4.1.24"))
	assert.True(t, ipCIDRMapList.Check("43.224.242.100"))
	assert.True(t, ipCIDRMapList.Check("101.236.0.1"))
	assert.False(t, ipCIDRMapList.Check("101.240.0.1"))
	assert.False(t, ipCIDRMapList.Check("5.5.5.5"))

	ipCIDRMapList.InsertSorted("1.3.0.0/16")
	assert.True(t, ipCIDRMapList.Check("1.3.1.1"))
}
