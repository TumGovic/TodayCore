// Package bridge connects FinalMask to sing-box dialers and listeners.
package bridge

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/transport/finalmask"
	"github.com/sagernet/sing-box/transport/finalmask/conf"
	"github.com/sagernet/sing-box/transport/finalmask/internal/xnet"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

// Masks builds the TCP and UDP mask instances described by options.
func Masks(options *option.FinalMaskOptions) ([]finalmask.TCPMask, []finalmask.UDPMask, error) {
	if options == nil {
		return nil, nil, nil
	}
	build := func(entry option.FinalMaskEntry, tcp bool) (any, error) {
		mask := conf.Mask{Type: entry.Type}
		if len(entry.Settings) > 0 {
			settings := json.RawMessage(entry.Settings)
			mask.Settings = &settings
		}
		return mask.Build(tcp)
	}
	var tcpMasks []finalmask.TCPMask
	for i, entry := range options.TCP {
		instance, err := build(entry, true)
		if err != nil {
			return nil, nil, E.Cause(err, "finalmask.tcp[", i, "] (", entry.Type, ")")
		}
		tcpMask, ok := instance.(finalmask.TCPMask)
		if !ok {
			return nil, nil, E.New("finalmask.tcp[", i, "]: ", entry.Type, " is not a TCP mask")
		}
		tcpMasks = append(tcpMasks, tcpMask)
	}
	var udpMasks []finalmask.UDPMask
	for i, entry := range options.UDP {
		instance, err := build(entry, false)
		if err != nil {
			return nil, nil, E.Cause(err, "finalmask.udp[", i, "] (", entry.Type, ")")
		}
		udpMask, ok := instance.(finalmask.UDPMask)
		if !ok {
			return nil, nil, E.New("finalmask.udp[", i, "]: ", entry.Type, " is not a UDP mask")
		}
		udpMasks = append(udpMasks, udpMask)
	}
	return tcpMasks, udpMasks, nil
}

func socksaddrFromDestination(destination xnet.Destination) M.Socksaddr {
	if destination.Address == nil {
		return M.Socksaddr{Port: destination.Port.Value()}
	}
	if destination.Address.Family().IsDomain() {
		return M.Socksaddr{Fqdn: destination.Address.Domain(), Port: destination.Port.Value()}
	}
	return M.SocksaddrFrom(M.AddrFromIP(destination.Address.IP()), destination.Port.Value())
}

func destinationFromSocksaddr(network xnet.Network, destination M.Socksaddr) xnet.Destination {
	var address xnet.Address
	if destination.IsFqdn() {
		address = xnet.DomainAddress(destination.Fqdn)
	} else {
		address = xnet.IPAddress(destination.Addr.Unmap().AsSlice())
	}
	return xnet.Destination{Network: network, Address: address, Port: xnet.Port(destination.Port)}
}

// Dialer applies FinalMask to every connection made by the upstream dialer,
// like Xray's streamSettings.finalmask on outbounds.
type Dialer struct {
	upstream N.Dialer
	mask     *finalmask.FinalMask
}

// NewDialer wraps upstream, or returns it unchanged when options is empty.
func NewDialer(upstream N.Dialer, options *option.FinalMaskOptions) (N.Dialer, error) {
	tcpMasks, udpMasks, err := Masks(options)
	if err != nil {
		return nil, err
	}
	if len(tcpMasks) == 0 && len(udpMasks) == 0 {
		return upstream, nil
	}
	d := &Dialer{upstream: upstream}
	d.mask = finalmask.NewFinalMask(tcpMasks, udpMasks, d.dialTCP, listenTCP, d.dialUDP, listenUDP)
	return d, nil
}

// BaseDial overrides how the sockets below the masks are created for one
// call, e.g. to keep sing-box's parallel interface dialing.
type BaseDial struct {
	DialContext  func(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error)
	ListenPacket func(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error)
}

type baseDialKey struct{}

// WithBaseDial attaches base to ctx for the masked dialer.
func WithBaseDial(ctx context.Context, base *BaseDial) context.Context {
	return context.WithValue(ctx, baseDialKey{}, base)
}

func baseDialFromContext(ctx context.Context) *BaseDial {
	base, _ := ctx.Value(baseDialKey{}).(*BaseDial)
	return base
}

func (d *Dialer) dialTCP(ctx context.Context, destination xnet.Destination) (net.Conn, error) {
	if base := baseDialFromContext(ctx); base != nil && base.DialContext != nil {
		return base.DialContext(ctx, N.NetworkTCP, socksaddrFromDestination(destination))
	}
	return d.upstream.DialContext(ctx, N.NetworkTCP, socksaddrFromDestination(destination))
}

func (d *Dialer) dialUDP(ctx context.Context, destination xnet.Destination) (net.PacketConn, net.Addr, error) {
	socksaddr := socksaddrFromDestination(destination)
	var remoteAddr net.Addr
	if socksaddr.IsFqdn() {
		// Helper dials of some masks (xdns resolvers, realm STUN) may name
		// hosts; the proxied destination itself is already resolved.
		udpAddr, err := net.DefaultResolver.LookupNetIP(ctx, "ip", socksaddr.Fqdn)
		if err != nil {
			return nil, nil, err
		}
		if len(udpAddr) == 0 {
			return nil, nil, E.New("no address for ", socksaddr.Fqdn)
		}
		socksaddr = M.SocksaddrFrom(udpAddr[0], socksaddr.Port)
	}
	remoteAddr = socksaddr.UDPAddr()
	listenPacket := d.upstream.ListenPacket
	if base := baseDialFromContext(ctx); base != nil && base.ListenPacket != nil {
		listenPacket = base.ListenPacket
	}
	packetConn, err := listenPacket(ctx, socksaddr)
	if err != nil {
		return nil, nil, err
	}
	return packetConn, remoteAddr, nil
}

func listenTCP(ctx context.Context, addr net.Addr) (net.Listener, error) {
	var listenConfig net.ListenConfig
	return listenConfig.Listen(ctx, "tcp", addr.String())
}

func listenUDP(ctx context.Context, addr net.Addr) (net.PacketConn, error) {
	var listenConfig net.ListenConfig
	return listenConfig.ListenPacket(ctx, "udp", addr.String())
}

func (d *Dialer) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	switch N.NetworkName(network) {
	case N.NetworkTCP:
		conn, err := d.mask.DialTCP(ctx, destinationFromSocksaddr(xnet.Network_TCP, destination))
		if err != nil {
			return nil, err
		}
		// Masks with a login (xmc) handshake on first use in Xray, which
		// happens right away there. sing-box may first set a short write
		// deadline to kick early connections, so finish it here.
		if handshaker, ok := conn.(interface{ Handshake() error }); ok {
			if deadline, hasDeadline := ctx.Deadline(); hasDeadline {
				_ = conn.SetDeadline(deadline)
			}
			err = handshaker.Handshake()
			_ = conn.SetDeadline(time.Time{})
			if err != nil {
				conn.Close()
				return nil, err
			}
		}
		return conn, nil
	case N.NetworkUDP:
		return d.mask.DialUDP(ctx, destinationFromSocksaddr(xnet.Network_UDP, destination))
	default:
		return nil, E.Extend(N.ErrUnknownNetwork, network)
	}
}

func (d *Dialer) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	conn, err := d.mask.DialUDP(ctx, destinationFromSocksaddr(xnet.Network_UDP, destination))
	if err != nil {
		return nil, err
	}
	return conn.(*xnet.PacketConnWrapper).PacketConn, nil
}

func (d *Dialer) Upstream() any {
	return d.upstream
}

// Listener applies FinalMask to inbound sockets, like Xray's
// streamSettings.finalmask on inbounds.
type Listener struct {
	mask *finalmask.FinalMask
}

// NewListener returns nil when options contain no masks. listen and
// listenPacket create the base sockets.
func NewListener(options *option.FinalMaskOptions, listen func(ctx context.Context, addr net.Addr) (net.Listener, error), listenPacket func(ctx context.Context, addr net.Addr) (net.PacketConn, error)) (*Listener, error) {
	tcpMasks, udpMasks, err := Masks(options)
	if err != nil {
		return nil, err
	}
	if len(tcpMasks) == 0 && len(udpMasks) == 0 {
		return nil, nil
	}
	dialTCP := func(ctx context.Context, destination xnet.Destination) (net.Conn, error) {
		var dialer net.Dialer
		return dialer.DialContext(ctx, "tcp", destination.NetAddr())
	}
	dialUDP := func(ctx context.Context, destination xnet.Destination) (net.PacketConn, net.Addr, error) {
		addr, err := net.ResolveUDPAddr("udp", destination.NetAddr())
		if err != nil {
			return nil, nil, err
		}
		packetConn, err := net.ListenPacket("udp", "")
		if err != nil {
			return nil, nil, err
		}
		return packetConn, addr, nil
	}
	return &Listener{mask: finalmask.NewFinalMask(tcpMasks, udpMasks, dialTCP, listen, dialUDP, listenPacket)}, nil
}

func (l *Listener) HasTCP() bool {
	return l != nil && l.mask.HasTCPMasks()
}

func (l *Listener) HasUDP() bool {
	return l != nil && l.mask.HasUDPMasks()
}

func (l *Listener) Listen(ctx context.Context, addr net.Addr) (net.Listener, error) {
	return l.mask.Listen(ctx, addr)
}

func (l *Listener) ListenPacket(ctx context.Context, addr net.Addr) (net.PacketConn, error) {
	return l.mask.ListenPacket(ctx, addr)
}
