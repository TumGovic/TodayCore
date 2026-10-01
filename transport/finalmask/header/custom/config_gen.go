// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package custom

type Expr struct {
	Op   string     `json:"op,omitempty"`
	Args []*ExprArg `json:"args,omitempty"`
}

func (x *Expr) GetOp() string {
	if x != nil {
		return x.Op
	}
	return ""
}

func (x *Expr) GetArgs() []*ExprArg {
	if x != nil {
		return x.Args
	}
	return nil
}

type ExprArg struct {
	Value isExprArg_Value
}

func (x *ExprArg) GetValue() isExprArg_Value {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *ExprArg) GetBytes() []byte {
	if x != nil {
		if x, ok := x.Value.(*ExprArg_Bytes); ok {
			return x.Bytes
		}
	}
	return nil
}

func (x *ExprArg) GetU64() uint64 {
	if x != nil {
		if x, ok := x.Value.(*ExprArg_U64); ok {
			return x.U64
		}
	}
	return 0
}

func (x *ExprArg) GetVar() string {
	if x != nil {
		if x, ok := x.Value.(*ExprArg_Var); ok {
			return x.Var
		}
	}
	return ""
}

func (x *ExprArg) GetMetadata() string {
	if x != nil {
		if x, ok := x.Value.(*ExprArg_Metadata); ok {
			return x.Metadata
		}
	}
	return ""
}

func (x *ExprArg) GetExpr() *Expr {
	if x != nil {
		if x, ok := x.Value.(*ExprArg_Expr); ok {
			return x.Expr
		}
	}
	return nil
}

type isExprArg_Value interface {
	isExprArg_Value()
}

type ExprArg_Bytes struct {
	Bytes []byte
}

type ExprArg_U64 struct {
	U64 uint64
}

type ExprArg_Var struct {
	Var string
}

type ExprArg_Metadata struct {
	Metadata string
}

type ExprArg_Expr struct {
	Expr *Expr
}

func (*ExprArg_Bytes) isExprArg_Value() {}

func (*ExprArg_U64) isExprArg_Value() {}

func (*ExprArg_Var) isExprArg_Value() {}

func (*ExprArg_Metadata) isExprArg_Value() {}

func (*ExprArg_Expr) isExprArg_Value() {}

type TCPItem struct {
	DelayMin int64  `json:"delay_min,omitempty"`
	DelayMax int64  `json:"delay_max,omitempty"`
	Rand     int32  `json:"rand,omitempty"`
	RandMin  int32  `json:"rand_min,omitempty"`
	RandMax  int32  `json:"rand_max,omitempty"`
	Packet   []byte `json:"packet,omitempty"`
	Save     string `json:"save,omitempty"`
	Var      string `json:"var,omitempty"`
	Expr     *Expr  `json:"expr,omitempty"`
}

func (x *TCPItem) GetDelayMin() int64 {
	if x != nil {
		return x.DelayMin
	}
	return 0
}

func (x *TCPItem) GetDelayMax() int64 {
	if x != nil {
		return x.DelayMax
	}
	return 0
}

func (x *TCPItem) GetRand() int32 {
	if x != nil {
		return x.Rand
	}
	return 0
}

func (x *TCPItem) GetRandMin() int32 {
	if x != nil {
		return x.RandMin
	}
	return 0
}

func (x *TCPItem) GetRandMax() int32 {
	if x != nil {
		return x.RandMax
	}
	return 0
}

func (x *TCPItem) GetPacket() []byte {
	if x != nil {
		return x.Packet
	}
	return nil
}

func (x *TCPItem) GetSave() string {
	if x != nil {
		return x.Save
	}
	return ""
}

func (x *TCPItem) GetVar() string {
	if x != nil {
		return x.Var
	}
	return ""
}

func (x *TCPItem) GetExpr() *Expr {
	if x != nil {
		return x.Expr
	}
	return nil
}

type TCPSequence struct {
	Sequence []*TCPItem `json:"sequence,omitempty"`
}

func (x *TCPSequence) GetSequence() []*TCPItem {
	if x != nil {
		return x.Sequence
	}
	return nil
}

type TCPConfig struct {
	Clients []*TCPSequence `json:"clients,omitempty"`
	Servers []*TCPSequence `json:"servers,omitempty"`
	Errors  []*TCPSequence `json:"errors,omitempty"`
}

func (x *TCPConfig) GetClients() []*TCPSequence {
	if x != nil {
		return x.Clients
	}
	return nil
}

func (x *TCPConfig) GetServers() []*TCPSequence {
	if x != nil {
		return x.Servers
	}
	return nil
}

func (x *TCPConfig) GetErrors() []*TCPSequence {
	if x != nil {
		return x.Errors
	}
	return nil
}

type UDPItem struct {
	Rand    int32  `json:"rand,omitempty"`
	RandMin int32  `json:"rand_min,omitempty"`
	RandMax int32  `json:"rand_max,omitempty"`
	Packet  []byte `json:"packet,omitempty"`
	Save    string `json:"save,omitempty"`
	Var     string `json:"var,omitempty"`
	Expr    *Expr  `json:"expr,omitempty"`
}

func (x *UDPItem) GetRand() int32 {
	if x != nil {
		return x.Rand
	}
	return 0
}

func (x *UDPItem) GetRandMin() int32 {
	if x != nil {
		return x.RandMin
	}
	return 0
}

func (x *UDPItem) GetRandMax() int32 {
	if x != nil {
		return x.RandMax
	}
	return 0
}

func (x *UDPItem) GetPacket() []byte {
	if x != nil {
		return x.Packet
	}
	return nil
}

func (x *UDPItem) GetSave() string {
	if x != nil {
		return x.Save
	}
	return ""
}

func (x *UDPItem) GetVar() string {
	if x != nil {
		return x.Var
	}
	return ""
}

func (x *UDPItem) GetExpr() *Expr {
	if x != nil {
		return x.Expr
	}
	return nil
}

type UDPConfig struct {
	Client []*UDPItem `json:"client,omitempty"`
	Server []*UDPItem `json:"server,omitempty"`
}

func (x *UDPConfig) GetClient() []*UDPItem {
	if x != nil {
		return x.Client
	}
	return nil
}

func (x *UDPConfig) GetServer() []*UDPItem {
	if x != nil {
		return x.Server
	}
	return nil
}

type UDPStandaloneConfig struct {
	Client []*UDPItem `json:"client,omitempty"`
	Server []*UDPItem `json:"server,omitempty"`
}

func (x *UDPStandaloneConfig) GetClient() []*UDPItem {
	if x != nil {
		return x.Client
	}
	return nil
}

func (x *UDPStandaloneConfig) GetServer() []*UDPItem {
	if x != nil {
		return x.Server
	}
	return nil
}
