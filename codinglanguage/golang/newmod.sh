#!/bin/sh

name=$1

mkdir $name
cd $name && go mod init example.com/aryaLee/golang/$name
