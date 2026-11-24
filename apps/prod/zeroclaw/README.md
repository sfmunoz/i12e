# zeroclaw

https://www.zeroclaw.com/

## zeroclaw.service

```
$ systemctl --user cat zeroclaw.service
# /home/i12e/.config/systemd/user/zeroclaw.service
[Unit]
Description=ZeroClaw daemon
After=network.target

[Service]
Type=simple
ExecStart=/home/i12e/.cargo/bin/zeroclaw daemon
Restart=always
RestartSec=3
# Ensure HOME is set so headless browsers can create profile/cache dirs.
Environment=HOME=%h
# Allow inheriting DISPLAY and XDG_RUNTIME_DIR from the user session
# so graphical/headless browsers can function correctly.
PassEnvironment=DISPLAY XDG_RUNTIME_DIR

[Install]
WantedBy=default.target
```
