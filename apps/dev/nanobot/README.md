# nanobot

https://nanobot.wiki/

## nanobot-gateway.service

```
$ systemctl --user cat nanobot-gateway.service
# /home/i12e/.config/systemd/user/nanobot-gateway.service
[Unit]
Description=Nanobot Gateway (nanobot-gateway)
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=/home/i12e
ExecStart=/home/i12e/.nanobot/venv/bin/python -m nanobot gateway --foreground --port 18790
Restart=on-failure
SuccessExitStatus=SIGTERM
RestartSec=10
Environment=PYTHONUNBUFFERED=1
NoNewPrivileges=yes

[Install]
WantedBy=default.target

# /home/i12e/.config/systemd/user/nanobot-gateway.service.d/50-allow-sudo.conf
[Service]
NoNewPrivileges=no
```
