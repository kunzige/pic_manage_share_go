# coding:utf-8
import os
import httpx
client = httpx.Client(http2=True, verify=False)
from lxml import etree
import sys


headers={
	"user-agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
}







ip=sys.argv[1]


#  爬取IP归属地

res=client.get(f"https://m.ip138.com/iplookup.asp?from=baidu&ip={ip}",headers=headers)
html=etree.HTML(res.content.decode())
address=html.xpath("/html/body/div/div[2]/div[1]/div/table/tbody/tr/td[2]/text()")
print(address[0].split("  ")[0])