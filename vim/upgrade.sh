#!/bin/sh

version=v9.1.0781
file=$version.tar.gz
wget https://github.com/vim/vim/archive/refs/tags/$file
tar -xzf $file
cd vim-$version
make && make install

# 
# yum install -y libpython3-devel.x86_64
# yum install -y ncurses-devel.x86_64
