package dialer

import (
	"context"
	"net"
	"time"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/transport/finalmask/bridge"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

func newFinalMaskDialer(upstream N.Dialer, options *option.FinalMaskOptions) (N.Dialer, error) {
	masked, err := bridge.NewDialer(upstream, options)
	if err != nil || masked == upstream {
		return masked, err
	}
	if parallel, isParallel := upstream.(ParallelInterfaceDialer); isParallel {
		return &finalMaskParallelDialer{Dialer: masked, upstream: parallel}, nil
	}
	return masked, nil
}

// Only ParallelInterfaceDialer is implemented: the parallel network helpers
// fall back to it, and implementing ParallelNetworkDialer here would make
// them call back into this dialer.
var _ ParallelInterfaceDialer = (*finalMaskParallelDialer)(nil)

// finalMaskParallelDialer keeps the parallel interface and network dialing
// of the default dialer, applying FinalMask to the sockets it creates.
type finalMaskParallelDialer struct {
	N.Dialer
	upstream ParallelInterfaceDialer
}

func (d *finalMaskParallelDialer) DialParallelInterface(ctx context.Context, network string, destination M.Socksaddr, strategy *C.NetworkStrategy, interfaceType []C.InterfaceType, fallbackInterfaceType []C.InterfaceType, fallbackDelay time.Duration) (net.Conn, error) {
	ctx = bridge.WithBaseDial(ctx, &bridge.BaseDial{
		DialContext: func(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
			return d.upstream.DialParallelInterface(ctx, network, destination, strategy, interfaceType, fallbackInterfaceType, fallbackDelay)
		},
		ListenPacket: func(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
			return d.upstream.ListenSerialInterfacePacket(ctx, destination, strategy, interfaceType, fallbackInterfaceType, fallbackDelay)
		},
	})
	return d.Dialer.DialContext(ctx, network, destination)
}

func (d *finalMaskParallelDialer) ListenSerialInterfacePacket(ctx context.Context, destination M.Socksaddr, strategy *C.NetworkStrategy, interfaceType []C.InterfaceType, fallbackInterfaceType []C.InterfaceType, fallbackDelay time.Duration) (net.PacketConn, error) {
	ctx = bridge.WithBaseDial(ctx, &bridge.BaseDial{
		ListenPacket: func(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
			return d.upstream.ListenSerialInterfacePacket(ctx, destination, strategy, interfaceType, fallbackInterfaceType, fallbackDelay)
		},
	})
	return d.Dialer.ListenPacket(ctx, destination)
}

func (d *finalMaskParallelDialer) Upstream() any {
	return d.Dialer
}
