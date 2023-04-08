package ipsearch_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/haoel/ipsearch"
	"github.com/stretchr/testify/assert"
)

const (
	TestDir      = "data"
	IPv4CIDRFile = "data/china_ip_list.txt"
	IPv4GeoFile  = "data/asn-country-ipv4.csv"
)

func TestCIDRSearch(t *testing.T) {
	search := ipsearch.NewIPSearch(cidrs, ipsearch.CIDR)

	for _, data := range testCIDRDataList {
		ip := search.Search(data.ip)
		assert.Equal(t, ip != nil, data.find)
		if ip != nil {
			assert.Equal(t, ip.CIDR(), data.cidr)
		}
	}
}

func TestLoadFromFile(t *testing.T) {

	search, err := ipsearch.NewIPSearchWithFile("not-exist-file", ipsearch.CIDR)
	assert.NotNil(t, err)

	search, err = ipsearch.NewIPSearchWithFile(IPv4CIDRFile, ipsearch.CIDR)
	assert.Nil(t, err)
	testCIDRSearch(t, search)

	search, err = ipsearch.NewIPSearchWithFile(IPv4GeoFile, ipsearch.Geo)
	assert.Nil(t, err)
	testGeoSearch(t, search)
}

func TestLoadFromURL(t *testing.T) {

	endpoint := "127.0.0.1:9898"
	protocol := "http://" + endpoint + "/"
	search, err := ipsearch.NewIPSearchWithFileFromURL(protocol+"not-exist-file", ipsearch.CIDR)
	assert.NotNil(t, err)

	//start http server
	go func() {
		http.ListenAndServe(endpoint, http.FileServer(http.Dir(TestDir)))
	}()

	for {
		search, err = ipsearch.NewIPSearchWithFileFromURL(protocol+"china_ip_list.txt", ipsearch.CIDR)
		if err != nil && strings.Contains(err.Error(), "connection refused") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		assert.Nil(t, err)
		testCIDRSearch(t, search)
		break
	}

	for {
		search, err = ipsearch.NewIPSearchWithFileFromURL(protocol+"/asn-country-ipv4.csv", ipsearch.Geo)
		if err != nil && strings.Contains(err.Error(), "connection refused") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		assert.Nil(t, err)
		testGeoSearch(t, search)
		break
	}
}

func testCIDRSearch(t *testing.T, search *ipsearch.IPSearch) {
	type testCIDRData struct {
		ip   string
		cidr string
		find bool
	}
	var testCIDRDataList = []testCIDRData{
		{"1.0.1.24", "1.0.1.0/24", true},
		{"8.8.8.8", "", false},
		{"1.1.1.1", "", false},
	}
	for _, data := range testCIDRDataList {
		ip := search.Search(data.ip)
		assert.Equal(t, ip != nil, data.find)
		if ip != nil {
			assert.Equal(t, ip.CIDR(), data.cidr)
		}
	}
}

func testGeoSearch(t *testing.T, search *ipsearch.IPSearch) {
	type testGeoData struct {
		ip      string
		country string
	}
	var testGeoDataList = []testGeoData{
		{"27.100.25.1", "IN"},
		{"31.25.64.1", "SE"},
		{"185.226.6.1", "US"},
	}
	for _, data := range testGeoDataList {
		ip := search.Search(data.ip)
		assert.NotNil(t, ip)
		assert.Equal(t, ip.Country(), data.country)
	}
}
