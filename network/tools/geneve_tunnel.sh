#!/bin/sh
netns=vpc0
brname=vpcbr
ip=10.0.0.6/24
remote=100.64.100.177
systemctl disable --now firewalld

sudo ip netns add $netns

sudo ovs-vsctl add-br $brname
sudo ovs-vsctl add-port $brname veth1 -- set interface veth1 type=internal

sudo ip link set dev veth1 netns $netns
sudo ip netns exec $netns ifconfig veth1 $ip up

sudo ip link set $brname up

# ovs-vsctl add-port $brname tun0 -- set interface tun0 type=vxlan options:remote_ip=$remote options:key=123
ovs-vsctl add-port $brname tun0 -- set interface tun0 type=geneve options:remote_ip=$remote options:key=123

