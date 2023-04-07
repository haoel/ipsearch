package ipsearch_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/haoel/ipsearch"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	search := ipsearch.NewIPSearch(cidrs)
	assert.True(t, search.Check("1.0.1.24"))
	assert.True(t, search.Check("1.4.1.1"))
	assert.True(t, search.Check("43.224.242.100"))
	assert.True(t, search.Check("101.236.0.1"))
	assert.False(t, search.Check("101.240.0.1"))
}

func TestChina(t *testing.T) {

	search, err := ipsearch.NewIPSearchWithFile("not-exist-file")
	assert.NotNil(t, err)

	search, err = ipsearch.NewIPSearchWithFile(IPv4File)
	assert.Nil(t, err)
	testSearch(t, search)
}

func TestChinaFromURL(t *testing.T) {

	endpoint := "127.0.0.1:9898"
	search, err := ipsearch.NewIPSearchWithFileFromURL("http://" + endpoint + "/not-exist-file")
	assert.NotNil(t, err)

	//start http server
	go func() {
		http.ListenAndServe(endpoint, http.FileServer(http.Dir(TestDir)))
	}()

	for {
		search, err = ipsearch.NewIPSearchWithFileFromURL("http://" + endpoint + "/china_ip_list.txt")
		if err!= nil && strings.Contains(err.Error(), "connection refused") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		assert.Nil(t, err)
		testSearch(t, search)
		break
	}
}

func testSearch(t *testing.T, search *ipsearch.IPSearch) {
	assert.True(t, search.Check("1.0.1.24"))
	assert.False(t, search.Check("1.1.1.1"))
	assert.False(t, search.Check("8.8.8.8"))
}
