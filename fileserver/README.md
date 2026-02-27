# fileserver

## References

- https://rustfs.com/
- https://github.com/rustfs/rustfs

## Usage

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
