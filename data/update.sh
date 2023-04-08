#!/usr/bin/env bash

BASEDIR=$(dirname "$0")

cd $BASEDIR

curl -LO https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt
curl -LO https://cdn.jsdelivr.net/npm/@ip-location-db/asn-country/asn-country-ipv4.csv