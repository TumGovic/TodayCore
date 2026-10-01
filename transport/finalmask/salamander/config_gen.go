// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package salamander

type Config struct {
	Password string `json:"password,omitempty"`
}

func (x *Config) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type GeckoConfig struct {
	Password      string `json:"password,omitempty"`
	MinPacketSize int32  `json:"MinPacketSize,omitempty"`
	MaxPacketSize int32  `json:"MaxPacketSize,omitempty"`
}

func (x *GeckoConfig) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *GeckoConfig) GetMinPacketSize() int32 {
	if x != nil {
		return x.MinPacketSize
	}
	return 0
}

func (x *GeckoConfig) GetMaxPacketSize() int32 {
	if x != nil {
		return x.MaxPacketSize
	}
	return 0
}
