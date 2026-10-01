// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package xmc

type Profile struct {
	Username string `json:"username,omitempty"`
	Uuid     []byte `json:"uuid,omitempty"`

	TexturesValue     string `json:"textures_value,omitempty"`
	TexturesSignature string `json:"textures_signature,omitempty"`
}

func (x *Profile) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Profile) GetUuid() []byte {
	if x != nil {
		return x.Uuid
	}
	return nil
}

func (x *Profile) GetTexturesValue() string {
	if x != nil {
		return x.TexturesValue
	}
	return ""
}

func (x *Profile) GetTexturesSignature() string {
	if x != nil {
		return x.TexturesSignature
	}
	return ""
}

type Config struct {
	Password      string     `json:"password,omitempty"`
	RsaPrivateKey []byte     `json:"rsa_private_key,omitempty"`
	RsaPublicKey  []byte     `json:"rsa_public_key,omitempty"`
	Hostname      string     `json:"hostname,omitempty"`
	Profiles      []*Profile `json:"profiles,omitempty"`
}

func (x *Config) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Config) GetRsaPrivateKey() []byte {
	if x != nil {
		return x.RsaPrivateKey
	}
	return nil
}

func (x *Config) GetRsaPublicKey() []byte {
	if x != nil {
		return x.RsaPublicKey
	}
	return nil
}

func (x *Config) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Config) GetProfiles() []*Profile {
	if x != nil {
		return x.Profiles
	}
	return nil
}
