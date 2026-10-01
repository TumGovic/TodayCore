// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package sudoku

type Config struct {
	Password     string   `json:"password,omitempty"`
	Ascii        string   `json:"ascii,omitempty"`
	CustomTable  string   `json:"custom_table,omitempty"`
	PaddingMin   uint32   `json:"padding_min,omitempty"`
	PaddingMax   uint32   `json:"padding_max,omitempty"`
	CustomTables []string `json:"custom_tables,omitempty"`
}

func (x *Config) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Config) GetAscii() string {
	if x != nil {
		return x.Ascii
	}
	return ""
}

func (x *Config) GetCustomTable() string {
	if x != nil {
		return x.CustomTable
	}
	return ""
}

func (x *Config) GetPaddingMin() uint32 {
	if x != nil {
		return x.PaddingMin
	}
	return 0
}

func (x *Config) GetPaddingMax() uint32 {
	if x != nil {
		return x.PaddingMax
	}
	return 0
}

func (x *Config) GetCustomTables() []string {
	if x != nil {
		return x.CustomTables
	}
	return nil
}
