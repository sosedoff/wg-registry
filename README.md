# wg-registry [![Go Report Card](https://goreportcard.com/badge/github.com/sosedoff/wg-registry)](https://goreportcard.com/report/github.com/sosedoff/wg-registry)

WireGuard server configuration and user management portal.

## Overview

Main goal of the `wg-registry` project is to simplify WireGuard server operations,
device and user management. Out of the box WireGuard does not provide a GUI to manage
client configuration or handle user authentication, such as via Google account, so
it's the main reason why this project exists in the first place.

WireGuard Registry is a pretty simple application with no dependencies
(other than WireGuard itself) and could be installed on the server by downloading
a single binary file. It's written in Go and embeds all assets in the executable, so
there is no need to drag any application runtime along (ruby/node/python/etc).
There's no special sauce involved in server management, everything is done via `wg` 
and `wg-quick` tooling provided by WireGuard. 

Internally, server and client configuration is stored in a file database (JSON-based)
and offers support for more traditional stores like PostgreSQL. Typically people install
WireGuard and keep configuration local so the registry project tries to keep it that way.
Additional encryption on the storage level will be implemented in future releases.

User authentication is currently implemented via Google OAuth flow, with the first user
signing into the system to become an admin. Registry project aiming to be a self-service
portal for users within an organization, so adding more users and devices is done by
users of themselves. Adding or removing a device from the system will automatically 
reconfigure WireGuard so there's no extra action required from the admin. 

For convenience purposes there's option to secure the server HTTP endpoint by adding
automatic LetsEncrypt certificate management. Having such a feature eliminates need to
run a HTTP proxy (nginx, etc) in front of the registry and requires almost no configuration.

## Features

- Written in Go and distributed as a single binary
- Authenticate users with a Google account
- Automatic WireGuard interface configuration
- Self-service device management
- Client configuration storage in PostgreSQL/SQLite/JSON-file
- Build-in LetsEncrypt support

Planned:

- Multi-server setup
- Store encryption
- More authentication providers
- IPV6 support
- Run in Docker

## Requirements

- WireGuard
- Google Account

**IMPORTANT:**
This project is not compatible with existing WireGuard configration, for best results
use the fresh install. Support for existing configurations might be added in future releases.

## Installation

TBD

## Configuration

Create a basic configuration file `/etc/wg-registry/config.json`:

```json
{
  "client_id": "google oauth client id",
  "client_secret": "google oauth secret",
  "client_domain": "mycorp.com",
  "client_whitelist": ["foo@gmail.com"],
  "store": "file:///etc/wg-resitry/data.json",
  "letsencrypt": {
    "email": "your email",
    "domain": "wg.mycorp.com",
    "dir": "/etc/wg-registry"
  }
}
```

Start the service:

```
wg-registry -c /etc/wg-registry/config.json
```

Now you should be able to visit `https://wg.mycorp.com` and configure server settings.

## Defaults

Default WireGuard configuration after the setup:

- Name: `private`
- Mode: `split`
- Interface: `wg0`
- IPV4 network: `10.10.0.0/24`
- Listen port: `51820`
- Persistent keepalive: `60`
- DNS: `n/a`

## Troubleshooting

### No communication between peers

Make sure the IP forwarding is enabled in your system:

```bash
sysctl -w net.ipv4.ip_forward=1
```

## License

MIT