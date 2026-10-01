package v2raywebsocket

import (
	"context"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// browserRequestURL mirrors Xray's websocket dialer when the browser dialer is
// enabled: the URL host is the Host header (or the server address) plus the
// port when it is not the default one.
func (c *Client) browserRequestURL() string {
	requestURL := c.requestURL
	if _, _, err := net.SplitHostPort(requestURL.Host); err != nil {
		if !(requestURL.Scheme == "ws" && c.serverAddr.Port == 80) && !(requestURL.Scheme == "wss" && c.serverAddr.Port == 443) {
			requestURL.Host = net.JoinHostPort(requestURL.Host, strconv.Itoa(int(c.serverAddr.Port)))
		}
	}
	return requestURL.String()
}

func (c *Client) dialBrowser(ctx context.Context) (net.Conn, error) {
	if c.maxEarlyData > 0 {
		return &browserEarlyConn{ctx: ctx, client: c, uri: c.browserRequestURL(), create: make(chan struct{})}, nil
	}
	return c.browserDialer.DialWS(ctx, c.browserRequestURL(), nil)
}

// browserEarlyConn defers the browser dial until the first write, so that up
// to max_early_data bytes can be sent in Sec-WebSocket-Protocol, which is the
// only header a browser lets a page control.
type browserEarlyConn struct {
	ctx    context.Context
	client *Client
	uri    string
	access sync.Mutex
	create chan struct{}
	conn   net.Conn
	err    error
}

func (c *browserEarlyConn) Write(b []byte) (int, error) {
	c.access.Lock()
	if c.conn != nil {
		c.access.Unlock()
		return c.conn.Write(b)
	}
	if c.err != nil {
		c.access.Unlock()
		return 0, c.err
	}
	earlyData, lateData := b, []byte(nil)
	if len(b) > int(c.client.maxEarlyData) {
		earlyData, lateData = b[:c.client.maxEarlyData], b[c.client.maxEarlyData:]
	}
	conn, err := c.client.browserDialer.DialWS(c.ctx, c.uri, earlyData)
	if err == nil && len(lateData) > 0 {
		_, err = conn.Write(lateData)
	}
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		c.err = err
	} else {
		c.conn = conn
	}
	close(c.create)
	c.access.Unlock()
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *browserEarlyConn) Read(b []byte) (int, error) {
	<-c.create
	if c.err != nil {
		return 0, c.err
	}
	return c.conn.Read(b)
}

func (c *browserEarlyConn) Close() error {
	c.access.Lock()
	defer c.access.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	if c.err == nil {
		c.err = net.ErrClosed
		close(c.create)
	}
	return nil
}

func (c *browserEarlyConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (c *browserEarlyConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (c *browserEarlyConn) SetDeadline(t time.Time) error {
	return os.ErrInvalid
}

func (c *browserEarlyConn) SetReadDeadline(t time.Time) error {
	return os.ErrInvalid
}

func (c *browserEarlyConn) SetWriteDeadline(t time.Time) error {
	return os.ErrInvalid
}

func (c *browserEarlyConn) NeedAdditionalReadDeadline() bool {
	return true
}
