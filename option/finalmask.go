package option

import "github.com/sagernet/sing/common/json"

// FinalMaskOptions is Xray's "finalmask" object. Masks are applied to the raw
// sockets below TLS and the V2Ray transport, in the order listed, exactly
// like Xray's streamSettings.finalmask.
type FinalMaskOptions struct {
	TCP []FinalMaskEntry `json:"tcp,omitempty"`
	UDP []FinalMaskEntry `json:"udp,omitempty"`
}

// FinalMaskEntry is one mask. Settings use Xray's field names verbatim.
type FinalMaskEntry struct {
	Type     string          `json:"type"`
	Settings json.RawMessage `json:"settings,omitempty"`
}
