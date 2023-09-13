#!/bin/sh
netns=vpc0
brname=vpcbr

ovs-vsctl del-br $brname
ip netns del $netns
# ip netns del lmy
