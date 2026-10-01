// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package xicmp

type Config struct {
	DGRAM bool     `json:"DGRAM,omitempty"`
	IPs   []string `json:"IPs,omitempty"`
}

func (x *Config) GetDGRAM() bool {
	if x != nil {
		return x.DGRAM
	}
	return false
}

func (x *Config) GetIPs() []string {
	if x != nil {
		return x.IPs
	}
	return nil
}
