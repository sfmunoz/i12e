# dev-tools

## dhcpd

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
$ ./dhcpd/run.sh
++ dirname ./dhcpd/run.sh
+ cd ./dhcpd
+ awk '!/^(#|$)/' udhcpd.conf
start           192.168.56.20
end             192.168.56.254
interface       vboxnet0
lease_file      /dev/null
static_lease 08:00:27:A8:56:51 192.168.56.51 fc1
static_lease 08:00:27:A8:56:52 192.168.56.52 fc2
static_lease 08:00:27:A8:56:53 192.168.56.53 fc3
+ sudo busybox udhcpd -f udhcpd.conf
udhcpd: started, v1.36.1
```
## fileserver

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
~/src/i12e/fileserver/run.sh
+ sudo -u root mkdir -p data logs
+ sudo -u root chown -R 10001:10001 data logs
+ docker run -it --rm --name rustfs -p 127.0.0.1:9000:9000 -p 127.0.0.1:9001:9001 -p 192.168.56.1:9000:9000 -p 192.168.56.1:9001:9001 -v /home/sfm/src/i12e/fileserver/data:/data -v /home/sfm/src/i12e/fileserver/logs:/logs rustfs/rustfs:latest
Initializing data directories: /data
Initializing log directory: /logs
!!!WARNING: Using default RUSTFS_ACCESS_KEY or RUSTFS_SECRET_KEY. Override them in production!
Starting: /usr/bin/rustfs  /data
RustFS Http API: http://172.17.0.2:9000  http://127.0.0.1:9000
RustFS Start Time: 2026-02-27 17:14:41
Console WebUI Start Time: 2026-02-27 17:14:41
Console WebUI available at: http://172.17.0.2:9001/rustfs/console/index.html
Console WebUI (localhost): http://127.0.0.1:9001/rustfs/console/index.html
RustFS server version: refs/tags/1.0.0-alpha.82 started successfully at [::]:9000, current time: 2026-02-27T17:14:41.587088569Z[Etc/Unknown]
```
