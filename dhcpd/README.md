# dhcpd

- [Requirements](#requirements)
  - [busybox](#busybox)
  - [ufw config](#ufw-config)

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
