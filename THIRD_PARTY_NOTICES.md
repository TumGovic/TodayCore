# Third-party notices

TodayCore is a derivative networking project built from code provided by multiple upstream projects. This file supplements the notices contained in individual source files; it does not replace them.

## sing-box

- Project: [SagerNet/sing-box](https://github.com/SagerNet/sing-box)
- TodayCore base: `1.15.0-alpha.6`
- Upstream license notice: GPL-3.0-or-later, preserved verbatim in [`LICENSES/SING-BOX-LICENSE.txt`](LICENSES/SING-BOX-LICENSE.txt)
- Copyright: the sing-box contributors, including the copyright notice retained in the preserved upstream license

The repository-level GPL-3.0 text is provided in [`LICENSE`](LICENSE). TodayCore preserves the upstream sing-box license notice separately and does not claim to be an official sing-box release. The names SagerNet and sing-box must not be used to imply upstream endorsement.

## Xray-core

- Project: [XTLS/Xray-core](https://github.com/XTLS/Xray-core)
- Porting reference: `26.9.9`, and `26.9.30` for the server-side and later ports
- License: Mozilla Public License 2.0
- License copy: [`LICENSES/MPL-2.0.txt`](LICENSES/MPL-2.0.txt)

Selected protocol behavior was derived or adapted from Xray-core, primarily in:

- `transport/v2rayxhttp/` — SplitHTTP/XHTTP client and server transport, including HTTP/3;
- `option/v2ray_xhttp.go` and XHTTP transport integration points;
- `transport/vlessenc/` — VLESS Encryption client and server behavior;
- `common/browserdialer/` — Xray's browser dialer, including its embedded `dialer.html` page;
- `common/browserheaders/` — browser-like request headers;
- `common/tls/reality_client.go` and `common/tls/reality_spider.go` — REALITY compatibility changes related to `X25519MLKEM768`, the protocol-compatible client version, `mldsa65Verify` and `spiderX`;
- `transport/finalmask/` — FinalMask and all of its masks, plus the parts of Xray's `common` packages they use (`transport/finalmask/internal/`); the configuration structs are generated from Xray's protobuf definitions;
- `cmd/sing-box/cmd_generate_xray.go` — key generation equivalent to `xray vlessenc` and `xray mldsa65`.

Source-file notices identify derived files and the relevant upstream areas. MPL-covered source remains available in source form in this repository. TodayCore-specific integration and other repository code remain governed by the license notices applicable to their respective files.

## No upstream endorsement

TodayCore is an independent community project. It is not maintained, sponsored, or endorsed by SagerNet, the sing-box project, XTLS, or the Xray-core project.
