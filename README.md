pdirector
=========

This simple TCP port proxy redirects and forwards TCP traffic.

Usage
-----
```
Usage: pdirector [arguments...]

  -fwd-host="localhost": An ip or hostname with forwarded port
  -fwd-port="": An opened port to forward traffic to
  -proxy-host="localhost": A local ip or named address for a proxy-port
  -proxy-port="": A proxy port for a forwarded connection
```
Example
-------
To create a fowrwarding proxy on localhost:8080 to google.com web servce

```
$ pdirector -fwd-host=www.google.com -fwd-port 80 -proxy-port 8080
```
