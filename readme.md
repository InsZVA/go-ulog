Go-ulog

## Why create this repo?

The ulogd2 project of linux netfilter is very convenient. But it perform not very good.
I have test it in a 64 core server with a 10Gbps multiqueue netcard only get 2Wpps.
Then I read its source code and found it use only one thread both to poll message from netlink
 and do calculate (parse etc). So I want to rewrite it and get a higher performance using multi-
cores.

## Some tests

I have test the pre-version, and receive 2 times than ulogd.