# LocalSend Auto Download

A command-line tool that downloads files shared by a LocalSend device on the
local network.

## Current status

- Downloads LocalSend v2 shares over **HTTP** on port `53317`.
- Accepts an optional PIN with `--pin`.
- Discovers the local session ID automatically unless `--session-id` is given.
- Downloads up to three files concurrently.
- Saves files to the current directory by default; `--base-path` changes the
  destination.
- HTTP request and idle-connection timeouts can be configured with
  `--timeout` and `--idle-timeout`.
- HTTPS is **not implemented**. Traffic, including the optional PIN, is not
  encrypted or authenticated. Use only on a trusted network.
- There are currently no Go test files. CI runs `go test ./...` and
  `go build ./...`; automated release or deployment is not configured.

## Usage

```sh
localsend-auto-download --host 192.168.1.20 --pin 123456
```

Pass `--help` to see all available flags. For example, to choose a download
directory and provide a session ID explicitly:

```sh
localsend-auto-download --host 192.168.1.20 --session-id 192.168.1.10 --base-path ./downloads
```

## HTTPS considerations

HTTPS is a future goal, but accepting a certificate generated on the fly
without verifying it would not make the connection safe. The TLS handshake
uses the presented certificate to authenticate the server; encryption without
reliable certificate verification can still allow a man-in-the-middle to
impersonate the LocalSend device. Certificate creation/issuance is separate
from the handshake, even if a receiver creates a certificate on demand.

Before adding HTTPS, the project needs a way to establish and verify the
device's identity, such as a trusted certificate or an explicit certificate
fingerprint trust/pinning flow. Disabling certificate verification should not
be used as a workaround.