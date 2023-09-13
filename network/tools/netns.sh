#!/bin/sh
ns=lmy
dev=ens5
ip=192.168.77.178/20

ip netns add $ns
ip link set dev $dev netns $ns
ip netns exec $ns ifconfig $dev $ip up
