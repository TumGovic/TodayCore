// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package fragment

type Config struct {
	PacketsFrom int64   `json:"packets_from,omitempty"`
	PacketsTo   int64   `json:"packets_to,omitempty"`
	MaxSplitMin int64   `json:"max_split_min,omitempty"`
	MaxSplitMax int64   `json:"max_split_max,omitempty"`
	LengthsMin  []int64 `json:"lengths_min,omitempty"`
	LengthsMax  []int64 `json:"lengths_max,omitempty"`
	DelaysMin   []int64 `json:"delays_min,omitempty"`
	DelaysMax   []int64 `json:"delays_max,omitempty"`
}

func (x *Config) GetPacketsFrom() int64 {
	if x != nil {
		return x.PacketsFrom
	}
	return 0
}

func (x *Config) GetPacketsTo() int64 {
	if x != nil {
		return x.PacketsTo
	}
	return 0
}

func (x *Config) GetMaxSplitMin() int64 {
	if x != nil {
		return x.MaxSplitMin
	}
	return 0
}

func (x *Config) GetMaxSplitMax() int64 {
	if x != nil {
		return x.MaxSplitMax
	}
	return 0
}

func (x *Config) GetLengthsMin() []int64 {
	if x != nil {
		return x.LengthsMin
	}
	return nil
}

func (x *Config) GetLengthsMax() []int64 {
	if x != nil {
		return x.LengthsMax
	}
	return nil
}

func (x *Config) GetDelaysMin() []int64 {
	if x != nil {
		return x.DelaysMin
	}
	return nil
}

func (x *Config) GetDelaysMax() []int64 {
	if x != nil {
		return x.DelaysMax
	}
	return nil
}
