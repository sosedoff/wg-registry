# wg-registry

Device registry and peer configuration manager for Wireguard

## Features

- Minimalistic install
- Credentials storage in Postgres/Sqlite
- Authenticate users with Google account
- Generate server and client configs
- Automatically configure WireGuard interface
- Self-service

## Usage

Create a new configuration file:

```json
{
  "database_url": "sqlite3:///wg-registry.db",
  "http_port": 80,
  "https_port": 443,
  "cookie_name": "wg-registry",
  "cookie_secret": "SECRET_VALUE",
  "client_id": "google oauth client id",
  "client_secret": "google oauth secret",
  "client_domain": "mycorp.com",
  "client_whitelist": [
    "foo@gmail.com"
  ],
  "database_url": "file:///etc/wireguard/registry.json"
}
```

Start the service:

```
wg-registry -c config.json
```