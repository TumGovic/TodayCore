[Русская версия](README.ru.md)

# TodayCore

**TodayCore** is an experimental networking core based on **sing-box 1.15.0-alpha.6**. It keeps the sing-box architecture and configuration format, and adds the modern Xray-core features needed to work with current Xray deployments, as both client and server. The features are ported from **Xray-core 26.9.9 and 26.9.30**.

> [!WARNING]
> TodayCore is a beta, community-driven project. Do not expect a fixed release schedule, guaranteed compatibility, or production support. Test every update before deploying it.

## Why TodayCore exists

TodayCore makes sing-box interoperate with current Xray deployments without turning it into a complete Xray clone. Only relevant modern features are considered; obsolete and legacy protocols are intentionally outside the project scope.

## Highlights

- **XHTTP** client and server, in every mode, over HTTP/1.1, HTTP/2 and HTTP/3.
- **REALITY** compatible with Xray 26.9 servers, including ML-DSA-65 verification and the spider.
- **VLESS post-quantum encryption** (`mlkem768x25519plus`) for outbounds and inbounds, with XTLS Vision.
- **FinalMask** with all of Xray's TCP and UDP masks.
- **Browser dialer**, so XHTTP and WebSocket traffic can be sent by a real browser.
- **Key generators** compatible with `xray vlessenc` and `xray mldsa65`.

## Added on top of sing-box

### XHTTP transport

The Xray SplitHTTP/XHTTP transport, for outbounds and inbounds:

- `packet-up`, `stream-up`, `stream-one`, and `auto` modes;
- HTTP/1.1, HTTP/2 (TLS and h2c), and HTTP/3 when the TLS ALPN is exactly `["h3"]`;
- session metadata in path, query, header, or cookie;
- uplink data in body, header, or cookie;
- X-Padding, including custom placement, key, header, and obfuscation mode;
- browser-like default request headers;
- XMUX connection and request reuse settings;
- separate `downloadSettings`;
- a server that accepts every mode Xray clients may choose.

### REALITY

The REALITY client works with Xray-core 26.9 servers:

- sends the `X25519MLKEM768` hybrid key share before plain `X25519`, as modern servers require;
- reports protocol-compatible client version `26.9.9`, avoiding the `reality verification failed` fallback on servers that require modern clients;
- `mldsa65_verify`: post-quantum ML-DSA-65 verification of servers configured with `mldsa65Seed`;
- `spider_x`: Xray's spider, which keeps crawling the target site after receiving a real certificate.

The key share and version behavior is automatic when using REALITY with uTLS.

### VLESS post-quantum encryption

Xray's VLESS Encryption on outbounds (`encryption`) and inbounds (`decryption`):

- `mlkem768x25519plus` with `native`, `xorpub`, and `random` modes;
- `1rtt` and `0rtt` handshakes with ticket lifetimes;
- padding profiles;
- X25519 and ML-KEM-768 keys, including relay key chains;
- XTLS Vision (`xtls-rprx-vision`) on top of the encryption layer.

Treat it as experimental until it receives broader interoperability testing.

### FinalMask

All of Xray's FinalMask masks:

- **TCP:** `fragment`, `sudoku`, `header-custom`, `xmc`;
- **UDP:** `noise`, `salamander` (including Gecko), `sudoku`, `header-custom`, `mkcp-legacy`, `udphop`, `xdns`, `xicmp`, `realm`.

Add a `finalmask` object to an inbound's or outbound's options. It wraps the raw TCP/UDP sockets below TLS and the transport, exactly like Xray's `streamSettings.finalmask`, and uses Xray's field names.

### Browser dialer

Xray's browser dialer for XHTTP and WebSocket outbounds: TodayCore serves a local page, and a browser that keeps it open makes the requests with its own TLS and HTTP fingerprint. Enable it per transport with `browserDialer` (XHTTP) or `browser_dialer` (WebSocket).

### Key generation

```bash
sing-box generate vlessenc
```

```bash
sing-box generate mldsa65-keypair
```

The first prints matching VLESS `decryption`/`encryption` strings, like `xray vlessenc`. The second prints a REALITY ML-DSA-65 seed for Xray servers and the verify key for clients, like `xray mldsa65`.

## Porting status

| Xray feature | Status in TodayCore | Notes |
| --- | --- | --- |
| XHTTP client and server | Ported | All modes; HTTP/1.1, h2, h2c, HTTP/3 |
| XHTTP metadata placement, X-Padding, XMUX | Ported | Xray-compatible field names |
| Browser dialer | Ported | XHTTP and WebSocket outbounds |
| REALITY `X25519MLKEM768` compatibility | Ported | Automatic with supported uTLS fingerprints |
| REALITY `mldsa65Verify` and `spiderX` | Ported | Client side (`mldsa65_verify`, `spider_x`) |
| VLESS Encryption | Experimental | Outbound (`encryption`) and inbound (`decryption`) |
| FinalMask | Ported | All TCP and UDP masks; inbounds and outbounds |
| REALITY `mldsa65Seed` (server) | Not ported | The sing-box REALITY server cannot sign with ML-DSA-65; use an Xray server |

All ported features were tested against Xray-core 26.9.30 in both directions. A few masks have environmental limits, listed in the extension guide: `xicmp` needs raw sockets, `realm` needs a signalling server, and `xdns` only carries small packets.

## Configuration

TodayCore uses the sing-box JSON configuration format. Field names that come from Xray (XHTTP options and FinalMask settings) keep Xray's spelling, so those parts of an Xray config can be reused as they are.

- [TodayCore extensions and configuration examples](docs/todaycore-extensions.md)

The original sing-box documentation remains in [`docs/`](docs/). For fields not described in the TodayCore extension guide, follow the upstream sing-box 1.15 configuration documentation.

## Build

### Requirements

- Go **1.25.5** or a compatible newer toolchain;
- Git;
- `make` for the standard build workflow.

### Standard build

```bash
make
```

Install to `$GOBIN`:

```bash
make install
```

A custom build can be requested through the Makefile:

```bash
TAGS="with_quic with_utls" make
```

The default build tags and required linker flags are maintained in `release/DEFAULT_BUILD_TAGS*` and `release/LDFLAGS`. Prefer the standard `make` target unless you know which platform-specific features you need. XHTTP over HTTP/3 needs `with_quic`, and REALITY needs `with_utls`; both are in the default tags.

## Development status and contributions

TodayCore is maintained **entirely on a voluntary basis**. Development may pause for long periods, and there is no promise that upstream sing-box or Xray changes will be merged immediately.

If you want the project to stay active, the best way to help is to participate directly:

1. test TodayCore against current Xray servers and clients;
2. report reproducible compatibility problems;
3. include logs, server/client versions, and a sanitized configuration;
4. add tests for protocol and wire-format behavior;
5. open a pull request with a focused, reviewable change.

Useful commits and pull requests are the strongest signal that there is real interest in continued development. Community initiative is welcome, and external contributions can directly determine how quickly the project moves forward.

Please keep pull requests small where possible, explain the upstream behavior being matched, and link to the corresponding Xray implementation or specification.

## Upstream projects and attribution

TodayCore is derived from:

- [SagerNet/sing-box](https://github.com/SagerNet/sing-box), base version 1.15.0-alpha.6;
- [XTLS/Xray-core](https://github.com/XTLS/Xray-core), source reference versions 26.9.9 and 26.9.30 for the ported features.

Files derived from Xray-core retain their MPL-2.0 attribution; see [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md). TodayCore modifications remain subject to the repository's GPL-3.0-or-later licensing terms. Review source-file headers and [`LICENSE`](LICENSE) before redistribution.

TodayCore is an independent community project and is not an official SagerNet or XTLS release. The TodayCore name must not be used to imply endorsement by either upstream project.

## Security

This project handles low-level networking, cryptography, and experimental protocol compatibility. Never assume that a successful build has been security-audited. Avoid publishing private keys, UUIDs, REALITY keys, VLESS Encryption keys, or complete production configurations in bug reports.
