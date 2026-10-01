// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package realm

import "crypto/tls"

type Family int32

const (
	Family_Dual Family = 0
	Family_V4   Family = 1
	Family_V6   Family = 2
)

var (
	Family_name = map[int32]string{
		0: "Dual",
		1: "V4",
		2: "V6",
	}
	Family_value = map[string]int32{
		"Dual": 0,
		"V4":   1,
		"V6":   2,
	}
)

type PortMapping struct {
	Enabled  bool  `json:"enabled,omitempty"`
	Timeout  int64 `json:"timeout,omitempty"`
	Lifetime int64 `json:"lifetime,omitempty"`
}

func (x *PortMapping) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *PortMapping) GetTimeout() int64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *PortMapping) GetLifetime() int64 {
	if x != nil {
		return x.Lifetime
	}
	return 0
}

type Config struct {
	Scheme      string       `json:"scheme,omitempty"`
	Host        string       `json:"host,omitempty"`
	Port        string       `json:"port,omitempty"`
	Token       string       `json:"token,omitempty"`
	ID          string       `json:"ID,omitempty"`
	StunServers []string     `json:"stun_servers,omitempty"`
	TlsConfig   *tls.Config  `json:"tls_config,omitempty"`
	IPMode      string       `json:"IPMode,omitempty"`
	PortMapping *PortMapping `json:"port_mapping,omitempty"`
}

func (x *Config) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

func (x *Config) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Config) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *Config) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Config) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Config) GetStunServers() []string {
	if x != nil {
		return x.StunServers
	}
	return nil
}

func (x *Config) GetTlsConfig() *tls.Config {
	if x != nil {
		return x.TlsConfig
	}
	return nil
}

func (x *Config) GetIPMode() string {
	if x != nil {
		return x.IPMode
	}
	return ""
}

func (x *Config) GetPortMapping() *PortMapping {
	if x != nil {
		return x.PortMapping
	}
	return nil
}

func (x Family) String() string {
	return Family_name[int32(x)]
}
