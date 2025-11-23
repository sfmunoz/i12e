# dhcpd

- [Requirements](#requirements)
  - [busybox](#busybox)
  - [ufw config](#ufw-config)
- [Usage](#usage)

## Requirements

### busybox

**busybox** or **busybox-static** (only one can be installed on **Linux Mint 22.2**):
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

## Usage

**Notice**: it may be required to start a VM in order to have **vboxnet0** up and running... otherwise DHCP server will fail to start.
```
$ ./dhcpd/run.sh
++ dirname ./dhcpd/run.sh
+ cd ./dhcpd
+ awk '!/^(#|$)/' udhcpd.conf
start           192.168.56.20
end             192.168.56.254
interface       vboxnet0
lease_file      /dev/null
static_lease 08:00:27:C3:44:AE 192.168.56.51 fc1
static_lease 08:00:27:3E:91:D7 192.168.56.52 fc2
static_lease 08:00:27:44:60:9E 192.168.56.53 fc3
+ sudo busybox udhcpd -f udhcpd.conf
udhcpd: started, v1.36.1
```
