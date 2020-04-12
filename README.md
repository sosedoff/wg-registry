# wg-registry

Device registry and peer configuration manager for WireGuard

## Features

- Minimalistic install
- Authenticate users with Google account
- Automatically configure WireGuard interface
- Client configuration storage in Postgres/Sqlite/File
- Self-service
- LetsEncrypt support

## Usage

Create a new configuration file:

```json
{
  "http_port": 80,
  "https_port": 443,
  "cookie_name": "wg-registry",
  "cookie_secret": "SECRET_VALUE",
  "client_id": "google oauth client id",
  "client_secret": "google oauth secret",
  "client_domain": "mycorp.com",
  "client_whitelist": ["foo@gmail.com"],
  "database_url": "file:///etc/wg-registry/data.json",
  "letsencrypt": {
    "email": "your account email",
    "domain": "mycorp.com",
    "dir": "/etc/wg-registry"
  }
}
```

Start the service:

```
wg-registry -c config.json
```

## Troubleshooting

Make sure the IP forwarding is enabled in your system:

```bash
sysctl -w net.ipv4.ip_forward=1
```