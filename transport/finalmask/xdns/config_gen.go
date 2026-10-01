// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package xdns

type DomainProto struct {
	Name       string  `json:"name,omitempty"`
	LenLimit   int32   `json:"len_limit,omitempty"`
	LabelLimit int32   `json:"label_limit,omitempty"`
	Types      []int32 `json:"types,omitempty"`
	Edns0      int32   `json:"edns0,omitempty"`
}

func (x *DomainProto) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DomainProto) GetLenLimit() int32 {
	if x != nil {
		return x.LenLimit
	}
	return 0
}

func (x *DomainProto) GetLabelLimit() int32 {
	if x != nil {
		return x.LabelLimit
	}
	return 0
}

func (x *DomainProto) GetTypes() []int32 {
	if x != nil {
		return x.Types
	}
	return nil
}

func (x *DomainProto) GetEdns0() int32 {
	if x != nil {
		return x.Edns0
	}
	return 0
}

type Config struct {
	Domains   []*DomainProto `json:"domains,omitempty"`
	Resolvers []any          `json:"resolvers,omitempty"`
	ExtraPoll int32          `json:"extra_poll,omitempty"`
}

func (x *Config) GetDomains() []*DomainProto {
	if x != nil {
		return x.Domains
	}
	return nil
}

func (x *Config) GetResolvers() []any {
	if x != nil {
		return x.Resolvers
	}
	return nil
}

func (x *Config) GetExtraPoll() int32 {
	if x != nil {
		return x.ExtraPoll
	}
	return 0
}

type TCPResolverProto struct {
	Addr string `json:"addr,omitempty"`
}

func (x *TCPResolverProto) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

type UDPResolverProto struct {
	Addr string `json:"addr,omitempty"`
}

func (x *UDPResolverProto) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}
