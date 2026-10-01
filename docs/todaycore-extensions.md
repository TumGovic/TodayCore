# TodayCore extensions

This document covers configuration fields added by TodayCore on top of sing-box 1.15.0-alpha.6.

## XHTTP transport

XHTTP is configured inside a compatible inbound's or outbound's `transport` object. TodayCore uses Xray-compatible camelCase field names for XHTTP-specific options.

### Minimal VLESS + XHTTP example

```json
{
  "outbounds": [
    {
      "type": "vless",
      "tag": "vless-xhttp",
      "server": "server.example.com",
      "server_port": 443,
      "uuid": "00000000-0000-0000-0000-000000000000",
      "tls": {
        "enabled": true,
        "server_name": "server.example.com",
        "alpn": ["h2", "http/1.1"],
        "utls": {
          "enabled": true,
          "fingerprint": "chrome"
        }
      },
      "transport": {
        "type": "xhttp",
        "host": "server.example.com",
        "path": "/api",
        "mode": "auto"
      }
    }
  ]
}
```

The server, UUID, host, path, TLS server name, and other authentication values must match the Xray server configuration.

### VLESS + XHTTP + REALITY example

```json
{
  "outbounds": [
    {
      "type": "vless",
      "tag": "vless-xhttp-reality",
      "server": "203.0.113.10",
      "server_port": 443,
      "uuid": "00000000-0000-0000-0000-000000000000",
      "tls": {
        "enabled": true,
        "server_name": "www.example.com",
        "alpn": ["h2", "http/1.1"],
        "utls": {
          "enabled": true,
          "fingerprint": "chrome"
        },
        "reality": {
          "enabled": true,
          "public_key": "REPLACE_WITH_REALITY_PUBLIC_KEY",
          "short_id": "0123456789abcdef"
        }
      },
      "transport": {
        "type": "xhttp",
        "host": "www.example.com",
        "path": "/assets",
        "mode": "stream-up"
      }
    }
  ]
}
```

The `X25519MLKEM768` compatibility fix is applied automatically to supported static uTLS fingerprints. No TodayCore-specific REALITY field is required.

### Modes

| Value | Meaning |
| --- | --- |
| `auto` | Select behavior automatically; this is the default |
| `packet-up` | Send uplink data as independent HTTP requests |
| `stream-up` | Stream the uplink while keeping a separate downlink |
| `stream-one` | Use one streaming request/response pair |

### Advanced XHTTP example

Ranges accept either a number such as `5` or a string such as `"5-10"`.

```json
{
  "type": "xhttp",
  "host": "cdn.example.com",
  "path": "/api?ed=2560",
  "mode": "packet-up",
  "headers": {
    "User-Agent": "Mozilla/5.0",
    "Accept-Language": "en-US,en;q=0.9"
  },
  "xPaddingBytes": "100-1000",
  "xPaddingObfsMode": true,
  "xPaddingKey": "x_padding",
  "xPaddingHeader": "Referer",
  "xPaddingPlacement": "queryInHeader",
  "xPaddingMethod": "repeat-x",
  "uplinkHTTPMethod": "POST",
  "sessionIDPlacement": "path",
  "seqPlacement": "path",
  "uplinkDataPlacement": "body",
  "uplinkChunkSize": "3000-4000",
  "scMaxEachPostBytes": 1000000,
  "scMinPostsIntervalMs": 30,
  "scMaxBufferedPosts": 30,
  "scStreamUpServerSecs": "20-80",
  "serverMaxHeaderBytes": 8192,
  "xmux": {
    "maxConnections": 3,
    "cMaxReuseTimes": "0-0",
    "hMaxRequestTimes": "600-900",
    "hMaxReusableSecs": "1800-3000",
    "hKeepAlivePeriod": 0
  }
}
```

### XHTTP field reference

| Field | Accepted values / purpose |
| --- | --- |
| `host` | HTTP Host/authority used by the XHTTP server |
| `path` | Request path; may include a query string |
| `mode` | `auto`, `packet-up`, `stream-up`, or `stream-one` |
| `headers` | Additional request headers; do not put `Host` here—use `host` |
| `xPaddingBytes` | X-Padding size or range; values must be positive when explicitly set |
| `xPaddingObfsMode` | Enable configurable padding placement and method |
| `xPaddingKey` | Padding query/cookie/header key |
| `xPaddingHeader` | Carrier header used by padding placement |
| `xPaddingPlacement` | `cookie`, `header`, `query`, or `queryInHeader` |
| `xPaddingMethod` | `repeat-x` or `tokenish` |
| `uplinkHTTPMethod` | Usually `POST`; `GET` is valid only with `packet-up` |
| `sessionIDPlacement` | `path`, `cookie`, `header`, or `query` |
| `sessionIDKey` | Custom key when the session ID is not placed in the path |
| `sessionIDTable` | Custom ID character table or a predefined alias such as `Base62`, `hex`, or `number` |
| `sessionIDLength` | Generated session-ID length or range |
| `seqPlacement` | `path`, `cookie`, `header`, or `query` |
| `seqKey` | Custom sequence key when not placed in the path |
| `uplinkDataPlacement` | `auto`, `body`, `cookie`, or `header`; cookie/header require `packet-up` |
| `uplinkDataKey` | Custom key for uplink data outside the body |
| `uplinkChunkSize` | Uplink chunk size or range |
| `noGRPCHeader` | Do not add the gRPC-like content type on streaming uplinks |
| `noSSEHeader` | Disable the SSE-like response header behavior |
| `scMaxEachPostBytes` | Maximum bytes in each packet-up POST |
| `scMinPostsIntervalMs` | Minimum interval between packet-up posts |
| `scMaxBufferedPosts` | Maximum number of buffered packet-up requests |
| `scStreamUpServerSecs` | Stream-up server duration or range |
| `serverMaxHeaderBytes` | Maximum accepted HTTP response-header size |
| `xmux` | XHTTP connection and request reuse controls |
| `downloadSettings` | Optional separate server/TLS/XHTTP settings for the downlink (client only) |
| `browserDialer` | Listen address of the browser dialer page (client only), see below |

`xmux.maxConnections` and `xmux.maxConcurrency` are mutually exclusive. If the entire `xmux` object is omitted, TodayCore follows the Xray defaults used by the port.

### Separate download settings

```json
{
  "type": "xhttp",
  "host": "upload.example.com",
  "path": "/up",
  "mode": "packet-up",
  "downloadSettings": {
    "server": "download.example.com",
    "server_port": 443,
    "tls": {
      "enabled": true,
      "server_name": "download.example.com",
      "alpn": ["h2"]
    },
    "host": "download.example.com",
    "path": "/down",
    "mode": "stream-one"
  }
}
```

### XHTTP inbound

The same `transport` object works on a VLESS (or any V2Ray-transport) inbound. The server accepts every mode the client may pick unless `mode` restricts it, exactly like Xray. TLS is taken from the inbound's `tls` object; without TLS the server speaks HTTP/1.1 and h2c.

```json
{
  "inbounds": [
    {
      "type": "vless",
      "listen": "::",
      "listen_port": 443,
      "users": [{ "uuid": "00000000-0000-0000-0000-000000000000" }],
      "tls": {
        "enabled": true,
        "server_name": "server.example.com",
        "alpn": ["h2", "http/1.1"],
        "certificate_path": "/etc/todaycore/cert.pem",
        "key_path": "/etc/todaycore/key.pem"
      },
      "transport": {
        "type": "xhttp",
        "host": "server.example.com",
        "path": "/api",
        "mode": "auto"
      }
    }
  ]
}
```

Server-side fields are the same as on the client. The ones that only matter on the server are `scMaxBufferedPosts`, `scStreamUpServerSecs`, `serverMaxHeaderBytes` and `noSSEHeader`. Like sing-box's other HTTP-based transports, the server always takes the source address from `X-Forwarded-For` when the header is present. Xray only does this for proxies listed in `trustedXForwardedFor`.

### XHTTP over HTTP/3

Set the TLS ALPN to exactly `["h3"]` on both sides, as in Xray. The inbound then listens on UDP instead of TCP. HTTP/3 needs the `with_quic` build tag, which is part of the default build. The client parrots Chrome's QUIC Initial and uses BBR, like Xray.

```json
"tls": { "enabled": true, "server_name": "server.example.com", "alpn": ["h3"] },
"transport": { "type": "xhttp", "path": "/api" }
```

### Browser dialer

The browser dialer lets a real browser make the XHTTP or WebSocket requests, so the traffic carries the browser's own TLS and HTTP fingerprint. TodayCore serves a page on the configured address; open it in a browser and keep the tab open. Xray enables it globally with the `XRAY_BROWSER_DIALER` environment variable, while TodayCore configures it per transport:

```json
"transport": { "type": "xhttp", "path": "/api", "browserDialer": "127.0.0.1:8080" }
```

```json
"transport": {
  "type": "ws",
  "path": "/ws",
  "browser_dialer": "127.0.0.1:8080",
  "max_early_data": 2048,
  "early_data_header_name": "Sec-WebSocket-Protocol"
}
```

The same limitations as in Xray apply:

- XHTTP can only use `packet-up` (which `auto` resolves to); `stream-up` and `stream-one` fail;
- REALITY outbounds never use the browser dialer;
- WebSocket early data is only possible through `Sec-WebSocket-Protocol`;
- the browser enforces its own certificate validation, CORS and cookie rules;
- the page listener lives for the whole process and is shared by all transports that use the same address.

### XHTTP limitations

- With `uplinkDataPlacement` set to `header` or `cookie`, also lower `scMaxEachPostBytes` and raise the server's `serverMaxHeaderBytes`, exactly as with Xray.
- Use an Xray-compatible peer. Fields that only exist in newer Xray versions are rejected.

## VLESS Encryption

VLESS Encryption is configured with the `encryption` field on a VLESS outbound and the `decryption` field on a VLESS inbound. Generate a matching pair with:

```bash
sing-box generate vlessenc
```

The command prints an X25519 pair and an ML-KEM-768 pair, the same way `xray vlessenc` does. Use one of them; do not mix them.

### Example

```json
{
  "outbounds": [
    {
      "type": "vless",
      "tag": "vless-pq",
      "server": "server.example.com",
      "server_port": 443,
      "uuid": "00000000-0000-0000-0000-000000000000",
      "encryption": "mlkem768x25519plus.native.0rtt.100-111-1111.75-0-111.50-0-3333.REPLACE_WITH_BASE64URL_PUBLIC_KEY",
      "tls": {
        "enabled": true,
        "server_name": "server.example.com"
      }
    }
  ]
}
```

Use the complete encryption string generated for or copied from the matching Xray server. Do not copy the placeholder literally.

### String format

```text
mlkem768x25519plus.<xor-mode>.<rtt-mode>[.<padding>].<key>[.<relay-key>...]
```

| Part | Values |
| --- | --- |
| Algorithm | `mlkem768x25519plus` |
| XOR mode | `native`, `xorpub`, or `random` |
| RTT mode | `1rtt` or `0rtt` |
| Padding | Optional dot-separated `chance-min-max` elements |
| Key | Raw URL-safe base64 without padding; decoded length must be 32 or 1184 bytes |
| Relay keys | Optional additional keys in relay order |

An empty `encryption` value or `"none"` disables VLESS Encryption.

### Server (inbound) example

```json
{
  "inbounds": [
    {
      "type": "vless",
      "listen": "::",
      "listen_port": 443,
      "users": [{ "uuid": "00000000-0000-0000-0000-000000000000", "flow": "xtls-rprx-vision" }],
      "decryption": "mlkem768x25519plus.native.600s.REPLACE_WITH_BASE64URL_PRIVATE_KEY"
    }
  ]
}
```

```text
mlkem768x25519plus.<xor-mode>.<ticket-lifetime>[.<padding>].<key>[.<relay-key>...]
```

| Part | Values |
| --- | --- |
| XOR mode | `native`, `xorpub`, or `random`; must match the client |
| Ticket lifetime | 0-RTT ticket lifetime in seconds, such as `600s` or `300-600s`; `0s` disables 0-RTT |
| Padding | Optional dot-separated `chance-min-max` elements |
| Key | 32-byte X25519 private key or 64-byte ML-KEM-768 seed, raw URL-safe base64 |
| Relay keys | Optional additional private keys in relay order |

XTLS Vision works on top of VLESS Encryption on both sides. As in Xray, Vision's direct copy writes to the connection directly below the encryption layer, which in `random` mode is still the XOR layer.

### VLESS Encryption limitations

- Experimental and not yet broadly interoperability-tested.
- The peer must support the same Xray VLESS Encryption wire format.
- A wrong key, order, mode, or padding string prevents the handshake.
- Keep all private key material on the server. The client configuration uses public key material supplied by the server operator.

## REALITY extensions

Both fields go into the outbound `tls.reality` object.

```json
"reality": {
  "enabled": true,
  "public_key": "REPLACE_WITH_REALITY_PUBLIC_KEY",
  "short_id": "0123456789abcdef",
  "mldsa65_verify": "REPLACE_WITH_BASE64URL_ML_DSA_65_PUBLIC_KEY",
  "spider_x": "/search?p=10-20&c=1-3&t=2-5&i=100-500&r=500-1000"
}
```

| Field | Xray name | Meaning |
| --- | --- | --- |
| `mldsa65_verify` | `mldsa65Verify` | ML-DSA-65 public key (1952 bytes, raw URL-safe base64) matching the server's `mldsa65Seed`. When set, only certificates carrying a valid ML-DSA-65 signature are accepted. |
| `spider_x` | `spiderX` | Start path for the spider, which keeps crawling the target site when a real certificate is received. The `p` (padding), `c` (concurrency), `t` (times), `i` (interval, ms) and `r` (return delay, ms) query parameters accept a number or a range and are removed from the path. |

`sing-box generate mldsa65-keypair` prints a `Seed` for the Xray server and the matching `Verify` key for clients, like `xray mldsa65`. TodayCore's own REALITY inbound cannot sign with ML-DSA-65, so use an Xray server for `mldsa65Seed`.

When the server has `mldsa65Seed`, its target must have a realistically large certificate chain, as real websites do. REALITY sizes the temporary certificate after the target's, and the signature does not fit into a small self-signed certificate.

## FinalMask

FinalMask is ported from Xray-core 26.9.30. Masks wrap the raw TCP and UDP sockets of an inbound or outbound, below TLS, REALITY and the V2Ray transport, exactly like Xray's `streamSettings.finalmask`. Put the same `finalmask` object into the inbound's listen fields or the outbound's dial fields:

```json
{
  "inbounds": [
    {
      "type": "vless",
      "listen": "::",
      "listen_port": 443,
      "users": [{ "uuid": "00000000-0000-0000-0000-000000000000" }],
      "finalmask": {
        "tcp": [{ "type": "sudoku", "settings": { "password": "REPLACE_ME" } }]
      }
    }
  ],
  "outbounds": [
    {
      "type": "vless",
      "server": "server.example.com",
      "server_port": 443,
      "uuid": "00000000-0000-0000-0000-000000000000",
      "finalmask": {
        "tcp": [
          { "type": "fragment", "settings": { "packets": "1-3", "length": "10-20", "delay": "1-3" } },
          { "type": "sudoku", "settings": { "password": "REPLACE_ME" } }
        ]
      }
    }
  ]
}
```

`tcp` and `udp` are lists applied in order, and every `settings` object uses Xray's field names verbatim, so a `finalmask` block can be copied from an Xray config. Masks that must match on both sides (everything except `fragment`, `noise` and `udphop`) need the same settings on the client and the server.

| Type | Socket | Notes |
| --- | --- | --- |
| `fragment` | TCP | Splits the first packets (`packets`: `tlshello` or a range), with `length`/`delay` or `lengths`/`delays` and `maxSplit` |
| `sudoku` | TCP, UDP | `password`, `ascii`, `customTable`, `customTables`, `paddingMin`, `paddingMax` |
| `header-custom` | TCP, UDP | Scripted handshake bytes (`clients`/`servers`/`errors` for TCP; `mode`, `client`, `server` for UDP), including captures and transforms |
| `xmc` | TCP | Minecraft login camouflage: `hostname`, `password`, `profiles` |
| `noise` | UDP | Junk packets before the first real packet: `noise` items with `packet`/`type`, `rand`, `delay`, or `type: "exp"` expressions; `reset` |
| `salamander` | UDP | Hysteria 2 Salamander obfuscation (`password`); with `packetSize` it becomes Gecko |
| `mkcp-legacy` | UDP | mKCP-era packet obfuscation: empty settings (original), `value` (AES-128-GCM), or `header` (`dns`, `dtls`, `srtp`, `utp`, `wechat`, `wireguard`) |
| `udphop` | UDP | Client-only port and address hopping: `mode` (`intervalLocal`, `intervalRemote`, `perConnRemote`), `interval`, `remoteIPs`, `remotePorts` |
| `xdns` | UDP | DNS tunnel: `domains` (with `types`, `lenLimit`, `labelLimit`, `edns0`), client `resolvers` (`udp`/`tcp` with `addr`), `extraPoll` |
| `xicmp` | UDP | ICMP tunnel: `dgram`, `ips` |
| `realm` | UDP | NAT traversal through a realm signalling server: `url` (`realm://token@host/id` or `realm+http://...`), `stunServers`, `ipMode`, `portMapping`, `tlsConfig` |

### FinalMask notes

- UDP masks apply to every UDP socket the inbound listens on or the outbound dials, for example XHTTP over HTTP/3, Hysteria 2, TUIC and plain UDP relayed by `direct`. Interoperability with Xray was tested over XHTTP/3 and over `direct` UDP.
- `xdns` carries only small packets, about 200–300 bytes, as in Xray. It does not fit QUIC, which needs 1200-byte packets, but works for small-datagram traffic. Without `extraPoll`, some downstream packets are lost under load, also as in Xray.
- `xicmp` needs raw ICMP sockets on the server: root or `CAP_NET_RAW` on Linux, root on macOS.
- `realm` needs a reachable realm signalling server and STUN servers. Its `tlsConfig` supports `serverName`, `alpn`, `pinnedPeerCertSha256` and `certificates` with `usage: "verify"`.
- Inbounds that need packet out-of-band data (TProxy UDP, the `resolved` service) reject UDP masks.
- Xray's `quicParams` is not part of FinalMask in TodayCore; QUIC settings stay with the respective transport.
- XMC's login runs when the connection is dialed. sing-box sets a short write deadline to start early connections, which would otherwise cut the multi round trip login. A failed XMC login is final for that connection instead of being retried on a half-read stream.

## Compatibility summary

TodayCore ports the XHTTP client and server (including HTTP/3 and the browser dialer), REALITY `mldsa65Verify` and `spiderX`, both sides of VLESS Encryption, and FinalMask. It does **not** port REALITY `mldsa65Seed` for its own REALITY inbound.

For standard sing-box fields, consult the existing documentation under `docs/configuration/`.
