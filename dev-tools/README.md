# dev-tools

## DHCPD

### busybox

**busybox** or **busybox-static** is used to run a simple DHCPD server (only one can be installed on **Linux Mint 22.3**):
```
# apt install busybox
```
or
```
# apt install busybox-static
```

### ufw config

**ufw** configuration to allow DHCP traffic:
```
# ufw allow in 67/udp
```
Alternative for just one interface:
```
# ufw allow in on vboxnet0 proto udp to any port 67
```

### Usage

**Notice**: it may be required to start a VM in order to have **vboxnet0** up and running... otherwise DHCP server will fail to start.
```
$ ./dev-tools/dhcpd.sh
+ awk '!/^(#|$)/' dhcpd.conf
start           192.168.56.20
end             192.168.56.254
interface       vboxnet0
lease_file      /dev/null
static_lease 08:00:27:A8:56:51 192.168.56.51 fc1
static_lease 08:00:27:A8:56:52 192.168.56.52 fc2
static_lease 08:00:27:A8:56:53 192.168.56.53 fc3
static_lease 08:00:27:A8:56:54 192.168.56.54 fc4
static_lease 08:00:27:A8:56:55 192.168.56.55 fc5
static_lease 08:00:27:A8:56:56 192.168.56.56 fc6
static_lease 08:00:27:A8:56:57 192.168.56.57 fc7
static_lease 08:00:27:A8:56:58 192.168.56.58 fc8
static_lease 08:00:27:A8:56:59 192.168.56.59 fc9
+ sudo busybox udhcpd -f dhcpd.conf
udhcpd: started, v1.36.1
```

## RustFS

### References

- https://rustfs.com/
- https://github.com/rustfs/rustfs

### ufw config

**ufw** configuration to allow TCP/9000 traffic:
```
# ufw allow in 9000/tcp
```
Alternative for just one interface:
```
# ufw allow in on vboxnet0 proto tcp to any port 9000
```

### Usage

```
$ ./dev-tools/rustfs.sh
+ mkdir -p .rustfs/data
+ cd .rustfs
+ ./rustfs --address 192.168.56.1:9000 --console-address 192.168.56.1:9001 data
(...)
RustFS server version: refs/tags/v0.0.24 started successfully at 192.168.56.1:9000, current time: 2026-02-28T08:47:43.329105831+00:00[UTC]
```
