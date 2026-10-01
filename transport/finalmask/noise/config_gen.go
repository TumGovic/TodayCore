// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

// Code generated from Xray-core's config.proto by pbconv; plain structs, no protobuf runtime.

package noise

type Segment_Kind int32

const (
	Segment_BYTES        Segment_Kind = 0
	Segment_RANDOM       Segment_Kind = 1
	Segment_RANDOM_ASCII Segment_Kind = 2
	Segment_RANDOM_DIGIT Segment_Kind = 3
	Segment_TIMESTAMP    Segment_Kind = 4
	Segment_COUNTER      Segment_Kind = 5
	Segment_NONCE        Segment_Kind = 6
)

var (
	Segment_Kind_name = map[int32]string{
		0: "BYTES",
		1: "RANDOM",
		2: "RANDOM_ASCII",
		3: "RANDOM_DIGIT",
		4: "TIMESTAMP",
		5: "COUNTER",
		6: "NONCE",
	}
	Segment_Kind_value = map[string]int32{
		"BYTES":        0,
		"RANDOM":       1,
		"RANDOM_ASCII": 2,
		"RANDOM_DIGIT": 3,
		"TIMESTAMP":    4,
		"COUNTER":      5,
		"NONCE":        6,
	}
)

type Segment struct {
	Kind    Segment_Kind `json:"kind,omitempty"`
	Bytes   []byte       `json:"bytes,omitempty"`
	MinSize int64        `json:"min_size,omitempty"`
	MaxSize int64        `json:"max_size,omitempty"`
}

func (x *Segment) GetKind() Segment_Kind {
	if x != nil {
		return x.Kind
	}
	return Segment_BYTES
}

func (x *Segment) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

func (x *Segment) GetMinSize() int64 {
	if x != nil {
		return x.MinSize
	}
	return 0
}

func (x *Segment) GetMaxSize() int64 {
	if x != nil {
		return x.MaxSize
	}
	return 0
}

type Item struct {
	RandMin      int64      `json:"rand_min,omitempty"`
	RandMax      int64      `json:"rand_max,omitempty"`
	RandRangeMin int32      `json:"rand_range_min,omitempty"`
	RandRangeMax int32      `json:"rand_range_max,omitempty"`
	Packet       []byte     `json:"packet,omitempty"`
	DelayMin     int64      `json:"delay_min,omitempty"`
	DelayMax     int64      `json:"delay_max,omitempty"`
	Segments     []*Segment `json:"segments,omitempty"`
}

func (x *Item) GetRandMin() int64 {
	if x != nil {
		return x.RandMin
	}
	return 0
}

func (x *Item) GetRandMax() int64 {
	if x != nil {
		return x.RandMax
	}
	return 0
}

func (x *Item) GetRandRangeMin() int32 {
	if x != nil {
		return x.RandRangeMin
	}
	return 0
}

func (x *Item) GetRandRangeMax() int32 {
	if x != nil {
		return x.RandRangeMax
	}
	return 0
}

func (x *Item) GetPacket() []byte {
	if x != nil {
		return x.Packet
	}
	return nil
}

func (x *Item) GetDelayMin() int64 {
	if x != nil {
		return x.DelayMin
	}
	return 0
}

func (x *Item) GetDelayMax() int64 {
	if x != nil {
		return x.DelayMax
	}
	return 0
}

func (x *Item) GetSegments() []*Segment {
	if x != nil {
		return x.Segments
	}
	return nil
}

type Config struct {
	ResetMin int64   `json:"reset_min,omitempty"`
	ResetMax int64   `json:"reset_max,omitempty"`
	Items    []*Item `json:"items,omitempty"`
}

func (x *Config) GetResetMin() int64 {
	if x != nil {
		return x.ResetMin
	}
	return 0
}

func (x *Config) GetResetMax() int64 {
	if x != nil {
		return x.ResetMax
	}
	return 0
}

func (x *Config) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x Segment_Kind) String() string {
	return Segment_Kind_name[int32(x)]
}
