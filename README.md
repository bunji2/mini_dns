# mini_dns

Mini DNS server for small and experimental network.

## Required library

```
go get github.com/miekg/dns
```

## Install

```
go get github.com/bunji2/mini_dns
```

## Usage

```
$ ./mini_dns -h
2024/12/14 14:14:37 mini_dns 1.1.a
Usage of ./mini_dns:
  -port int
        port number (default 53)
  -server string
        addr:port   of dns server
  -zone_file string
        zone file (default "./zone.txt")
```

## Sample zone file

```

$ORIGIN bunji2.com.
@	IN	SOA	ns.bunji2.com. root.bunji2.com.  (
		2014091801      ; Serial
		3600		; Refresh
		900		; Retry
		3600000		; Expire
		3600 )
    IN NS ns
    IN MX 10 mx
    IN TXT "v=spf1 +ip4:192.168.0.2 -all"
ns  IN A 192.168.0.1
mx  IN A 192.168.0.2
www IN A 192.168.0.3
```
## Example of operations using nslookup

Assuming that mini_dns.exe is running at 127.0.0.1.

### Looking up NS record of bunji2.com

```
C:\>nslookup -type=NS bunji2.com. 127.0.0.1
サーバー:  UnKnown
Address:  127.0.0.1

bunji2.com      nameserver = ns.bunji2.com
```

### Looking up TXT record of bunji2.com

```
C:\>nslookup -type=TXT bunji2.com. 127.0.0.1
サーバー:  UnKnown
Address:  127.0.0.1

bunji2.com      text =

        "v=spf1 +ip4:192.168.0.2 -all"
```

### Looking up MX record of bunji2.com

```
C:\>nslookup -type=MX bunji2.com. 127.0.0.1
サーバー:  UnKnown
Address:  127.0.0.1

bunji2.com      MX preference = 10, mail exchanger = mx.bunji2.com
```

### Looking up A record of mx.bunji2.com

```
C:\>nslookup -type=A mx.bunji2.com. 127.0.0.1
サーバー:  UnKnown
Address:  127.0.0.1

名前:    mx.bunji2.com
Address:  192.168.0.2
```

### Looking up A record of www.bunji2.com

```
C:\>nslookup -type=A www.bunji2.com. 127.0.0.1
サーバー:  UnKnown
Address:  127.0.0.1

名前:    www.bunji2.com
Address:  192.168.0.3
```
