// This file is derived from Xray-core (infra/conf/common.go and loader.go),
// licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.

package conf

import (
	"crypto/sha256"
	stdtls "crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"

	errors "github.com/sagernet/sing-box/transport/finalmask/internal/xerrors"
)

type stdTLSConfig = stdtls.Config

// Buildable is implemented by every mask settings type.
type Buildable interface {
	Build() (any, error)
}

type ConfigCreator func() any

type ConfigCreatorCache map[string]ConfigCreator

type JSONConfigLoader struct {
	cache     ConfigCreatorCache
	idKey     string
	configKey string
}

func NewJSONConfigLoader(cache ConfigCreatorCache, idKey string, configKey string) *JSONConfigLoader {
	return &JSONConfigLoader{
		idKey:     idKey,
		configKey: configKey,
		cache:     cache,
	}
}

func (v *JSONConfigLoader) LoadWithID(raw []byte, id string) (any, error) {
	id = strings.ToLower(id)
	creator, found := v.cache[id]
	if !found {
		return nil, errors.New("unknown config id: ", id)
	}
	config := creator()
	if err := json.Unmarshal(raw, config); err != nil {
		return nil, err
	}
	return config, nil
}

// Int32Range deserializes from "1-2" or 1, as in Xray. Value will be
// exchanged if From > To, use .Left and .Right to get original value if need.
type Int32Range struct {
	Left  int32
	Right int32
	From  int32
	To    int32
}

func (v *Int32Range) UnmarshalJSON(data []byte) error {
	defer v.ensureOrder()
	var str string
	var rawint int32
	if err := json.Unmarshal(data, &str); err == nil {
		left, right, err := ParseRangeString(str)
		if err == nil {
			v.Left, v.Right = int32(left), int32(right)
			return nil
		}
	} else if err := json.Unmarshal(data, &rawint); err == nil {
		v.Left = rawint
		v.Right = rawint
		return nil
	}

	return errors.New("Invalid integer range, expected either string of form \"1-2\" or plain integer.")
}

// ensureOrder() gives value to .From & .To and make sure .From < .To
func (r *Int32Range) ensureOrder() {
	r.From, r.To = r.Left, r.Right
	if r.From > r.To {
		r.From, r.To = r.To, r.From
	}
}

// "-114-514"   →  ["-114","514"]
// "-1919--810" →  ["-1919","-810"]
func splitFromSecondDash(s string) []string {
	parts := strings.SplitN(s, "-", 3)
	if len(parts) < 3 {
		return []string{s}
	}
	return []string{parts[0] + "-" + parts[1], parts[2]}
}

// ParseRangeString parses a range in string. Supports negative numbers.
// eg: "114-514" "-114-514" "-1919--810" "114514" ""(return 0)
func ParseRangeString(str string) (int, int, error) {
	if value, err := strconv.Atoi(str); err == nil {
		return value, value, nil
	}
	if str == "" {
		return 0, 0, nil
	}
	var pair []string
	if strings.HasPrefix(str, "-") {
		pair = splitFromSecondDash(str)
	} else {
		pair = strings.SplitN(str, "-", 2)
	}
	if len(pair) == 2 {
		left, err := strconv.Atoi(pair[0])
		right, err2 := strconv.Atoi(pair[1])
		if err == nil && err2 == nil {
			return left, right, nil
		}
	}
	return 0, 0, errors.New("invalid range string: ", str)
}

type PortRange struct {
	From uint32
	To   uint32
}

// PortList deserializes from "1000-2000,3000" or 3000, as in Xray.
type PortList struct {
	Range []PortRange
}

func parsePort(s string) (uint32, error) {
	value, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
	if err != nil || value > 65535 {
		return 0, errors.New("invalid port: ", s)
	}
	return uint32(value), nil
}

func (list *PortList) UnmarshalJSON(data []byte) error {
	var listStr string
	var number uint32
	if err := json.Unmarshal(data, &listStr); err != nil {
		if err2 := json.Unmarshal(data, &number); err2 != nil {
			return errors.New("invalid port: ", string(data)).Base(err2)
		}
	}
	for _, rangeStr := range strings.Split(listStr, ",") {
		trimmed := strings.TrimSpace(rangeStr)
		if len(trimmed) == 0 {
			continue
		}
		pair := strings.SplitN(trimmed, "-", 2)
		from, err := parsePort(pair[0])
		if err != nil {
			return errors.New("invalid port range: ", trimmed).Base(err)
		}
		to := from
		if len(pair) == 2 {
			to, err = parsePort(pair[1])
			if err != nil {
				return errors.New("invalid port range: ", trimmed).Base(err)
			}
		}
		list.Range = append(list.Range, PortRange{From: from, To: to})
	}
	if number != 0 {
		list.Range = append(list.Range, PortRange{From: number, To: number})
	}
	return nil
}

func (list *PortList) Build() *PortList {
	return list
}

func (list *PortList) Ports() []uint32 {
	var ports []uint32
	for _, r := range list.Range {
		for i := r.From; i <= r.To; i++ {
			ports = append(ports, i)
		}
	}
	return ports
}

// TLSConfig is the subset of Xray's TLS settings that is meaningful for the
// realm signalling client: server name, ALPN, certificate pinning and
// custom CA certificates.
type TLSConfig struct {
	ServerName           string   `json:"serverName"`
	ALPN                 []string `json:"alpn"`
	PinnedPeerCertSha256 string   `json:"pinnedPeerCertSha256"`
	Certificates         []struct {
		Certificate []string `json:"certificate"`
		Usage       string   `json:"usage"`
	} `json:"certificates"`
}

func (c *TLSConfig) Build() (*stdTLSConfig, error) {
	config := &stdtls.Config{
		ServerName: c.ServerName,
		NextProtos: c.ALPN,
	}
	for _, certificate := range c.Certificates {
		if certificate.Usage != "verify" {
			continue
		}
		if config.RootCAs == nil {
			config.RootCAs = x509.NewCertPool()
		}
		if !config.RootCAs.AppendCertsFromPEM([]byte(strings.Join(certificate.Certificate, "\n"))) {
			return nil, errors.New("invalid verify certificate")
		}
	}
	if c.PinnedPeerCertSha256 != "" {
		var pinned [][]byte
		for _, v := range strings.Split(c.PinnedPeerCertSha256, ",") {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			hashValue, err := hex.DecodeString(strings.ReplaceAll(v, ":", ""))
			if err != nil {
				return nil, err
			}
			if len(hashValue) != 32 {
				return nil, errors.New("incorrect pinnedPeerCertSha256 length: ", v)
			}
			pinned = append(pinned, hashValue)
		}
		config.InsecureSkipVerify = true
		config.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			for _, raw := range rawCerts {
				sum := sha256.Sum256(raw)
				for _, p := range pinned {
					if string(sum[:]) == string(p) {
						return nil
					}
				}
			}
			return errors.New("peer certificate does not match pinnedPeerCertSha256")
		}
	}
	return config, nil
}
