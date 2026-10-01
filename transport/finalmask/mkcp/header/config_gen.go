// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package header

type Config struct {
	ID     int32  `json:"ID,omitempty"`
	Domain string `json:"domain,omitempty"`
}

func (x *Config) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Config) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}
