# IP Search

This is a simple library to search for IP addresses for two different IP Databases: **IP CIDR List** and **IP Geo CVS**, which go

**Table of Contents**
- [IP Search](#ip-search)
  - [1. IP Database](#1-ip-database)
  - [2. Usage](#2-usage)
    - [2.1 Check an IP address is in the IP CIDR list](#21-check-an-ip-address-is-in-the-ip-cidr-list)
    - [2.2 Get the Country Code of an IP address](#22-get-the-country-code-of-an-ip-address)
  - [3. Technical Details](#3-technical-details)
  - [4. License](#4-license)




## 1. IP Database

This library can deal with the following IP databases:

- `china_ip_list.txt`:  A list of IP addresses in China. This list is
  from the [China IP List](https://github.com/17mon/china_ip_list) project.

- `asn-country-ipv4.txt`:  A list of IP addresses and their ASN and country
  information. This list is from the [IP Location DB](https://github.com/sapics/ip-location-db)
  project.

To update these two data files, run the following script:

```bash
./data/update.sh
```
> **Note**
>
>  - The CIDRs file must be a plain text file, and each line is a CIDR.
>  - The CIDRs cannot be overlapped.


## 2. Usage

### 2.1 Check an IP address is in the IP CIDR list


```go
package main

import (
	"fmt"
	"github.com/haoel/ipsearch"
)

func main() {
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
```

It also supports reading the URL files:

```go
url := "https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt"
search, err := ipsearch.NewIPSearchWithURL(url)
```

### 2.2 Get the Country Code of an IP address

```go
package main

import (
	"fmt"
	"github.com/haoel/ipsearch"
)

func main() {
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
```

## 3. Technical Details

The IP search is using the [Hash Table](https://en.wikipedia.org/wiki/Hash_table) and  [Binary Search](https://en.wikipedia.org/wiki/Binary_search_algorithm) algorithm.

The key of the hash table is the first 8 bits of the IP address, and the value is the sorted list of CIDRs that have the same first 8 bits.

```
[  1 ] -> [1.0.1.0/24],
          [1.0.2.0/23],
          [1.1.4.0/22]
          ...
[ 14 ] -> [14.0.0.0/21],
          [14.0.12.0/22],
          [14.1.0.0/22]
          ...
```

The `IPRangeMapList` is the hash table that stores all of the sorted CIDRs lists.

For the GeoIP database, it gives the `start` and `end` IP address, this would across the multiple hash table keys, so we need to split it.

For example,  if the IP range is `1.1.0.0 - 3.2.2.255`, we need to split it into three ranges:

- `1.1.0.0 - 1.255.255.255`
- `2.0.0.0 - 2.255.255.255`
- `3.0.0.0 - 3.2..2.255`

The split algorithm is in the `IPRange.Split()` function in the `iprange.go` file.

> **Note**
>
> And I didn't use the standard `net/netip` library, because of the following two reasons:
>
> - We can be free to customize and extend the algorithm
> - We can port the algorithm to other languages easily.
>

## 4. License

This project is MIT licensed. See the [LICENSE](LICENSE) file for details.