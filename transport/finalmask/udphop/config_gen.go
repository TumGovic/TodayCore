// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package udphop

type Config struct {
	Local       bool     `json:"local,omitempty"`
	Remote      bool     `json:"remote,omitempty"`
	RemoteOnce  bool     `json:"remote_once,omitempty"`
	IntervalMin int64    `json:"interval_min,omitempty"`
	IntervalMax int64    `json:"interval_max,omitempty"`
	RemoteIPs   []string `json:"remoteIPs,omitempty"`
	RemotePorts []uint32 `json:"remote_ports,omitempty"`
}

func (x *Config) GetLocal() bool {
	if x != nil {
		return x.Local
	}
	return false
}

func (x *Config) GetRemote() bool {
	if x != nil {
		return x.Remote
	}
	return false
}

func (x *Config) GetRemoteOnce() bool {
	if x != nil {
		return x.RemoteOnce
	}
	return false
}

func (x *Config) GetIntervalMin() int64 {
	if x != nil {
		return x.IntervalMin
	}
	return 0
}

func (x *Config) GetIntervalMax() int64 {
	if x != nil {
		return x.IntervalMax
	}
	return 0
}

func (x *Config) GetRemoteIPs() []string {
	if x != nil {
		return x.RemoteIPs
	}
	return nil
}

func (x *Config) GetRemotePorts() []uint32 {
	if x != nil {
		return x.RemotePorts
	}
	return nil
}
