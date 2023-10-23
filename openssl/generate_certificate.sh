#!/bin/sh

# 建立单独目录用来存放证书文件
mkdir ./mykey
cd ./mykey
 
# 使用以下命令生成私钥文件
# 其中，2048是私钥的长度，可以根据需要进行调整。
openssl genrsa -out ssl.key 2048
 
# 生成证书签名请求文件
# 在执行该命令时，需要填写一些证书信息，如国家、省份、城市、组织、单位、通用名称等。
openssl req -new -key ssl.key -out ssl.csr
 
# 最后根据这2个文件(ssl.key ssl.csr)生成自签名证书
# 其中，-days 3650表示证书的有效期为3650天，可以根据需要进行调整。
openssl x509 -req -days 3650 -in ssl.csr -signkey ssl.key -out ssl.crt
 
# 验证证书
# 该命令会输出证书的详细信息，包括证书的颁发者、有效期、公钥等。
openssl x509 -noout -text -in ssl.crt
