#!/bin/bash

ifaceName=$1

if [ -z "$ifaceName" ]; then
	echo "Usage: $0 <ifaceName>"
	exit 1
fi

tshark -i $ifaceName -w /tmp/$ifaceName.pcap
mv /tmp/$ifaceName.pcap ~/out.pcapng
chmod 644 ~/out.pcapng
