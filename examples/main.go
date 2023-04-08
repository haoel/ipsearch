// package main is the example of using ipsearch.
package main

import (
	"fmt"

	"github.com/haoel/ipsearch"
)

func main() {
	checkChinaIP()
	findIPCountry()
}

func checkChinaIP() {
	search, err := ipsearch.NewIPSearchWithFile("./data/china_ip_list.txt", ipsearch.CIDR)
	if err != nil {
		panic(err)
	}

	ipStr := "114.114.114.114"
	ip := search.Search(ipStr)
	if ip != nil {
		fmt.Printf("IP [%s] is in China\n", ipStr)
	} else {
		fmt.Printf("IP [%s] is not China\n", ipStr)
	}
}

func findIPCountry() {
	search, err := ipsearch.NewIPSearchWithFile("./data/asn-country-ipv4.csv", ipsearch.Geo)
	if err != nil {
		panic(err)
	}
	ipStr := "8.8.8.8"
	ip := search.Search(ipStr)
	if ip != nil {
		fmt.Printf("IP [%s] Country Code: %s\n", ipStr, ip.Country())
	} else {
		fmt.Printf("IP [%s] is not found!\n", ipStr)
	}
}
