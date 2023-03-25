#!/bin/bash

ifaceName=$1

if [ -z "$ifaceName" ]; then
	echo "Usage: $0 <ifaceName>"
	exit 1
fi

tshark -i $ifaceName -w /tmp/$ifaceName.pcap
mv /tmp/$ifaceName.pcap /home/justin/out.pcapng
chmod 644 /home/justin/out.pcapng
