# IP Search

This is a simple lib to search for IP addresses in a CIDRs list.

## Usage


```go
import "github.com/haoel/ipsearch"

package main

search, err := ipsearch.NewIPSearchWithFile("data/china_ip_list.txt")
if err != nil {
    panic(err)
}

ip := "1.1.1.1"
if search.Check(ip) {
    fmt.Println("The IP %s is in the list", ip)
} else {
    fmt.Println("The IP %s is not in the list", ip)
}
```

It also supports reading the URL files:

```go
url := "https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt"
search, err := ipsearch.NewIPSearchWithURL(url)
```

> **Note:**
>
>  - The CIDRs file must be a plain text file, and each line is a CIDR.
>  - The CIDRs cannot be overlapped.