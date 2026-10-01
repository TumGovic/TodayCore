package vlessenc

import (
	"bytes"
	"crypto/ecdh"
	"crypto/mlkem"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testKeyPair struct {
	server string
	client string
}

func newX25519Pair(t *testing.T) testKeyPair {
	key, err := ecdh.X25519().GenerateKey(rand.Reader)
	require.NoError(t, err)
	return testKeyPair{
		server: base64.RawURLEncoding.EncodeToString(key.Bytes()),
		client: base64.RawURLEncoding.EncodeToString(key.PublicKey().Bytes()),
	}
}

func newMLKEMPair(t *testing.T) testKeyPair {
	seed := make([]byte, 64)
	_, err := rand.Read(seed)
	require.NoError(t, err)
	key, err := mlkem.NewDecapsulationKey768(seed)
	require.NoError(t, err)
	return testKeyPair{
		server: base64.RawURLEncoding.EncodeToString(seed),
		client: base64.RawURLEncoding.EncodeToString(key.EncapsulationKey().Bytes()),
	}
}

// exchange runs one client/server handshake over a TCP loopback connection
// and checks that data flows in both directions, including records larger
// than one 8192 byte chunk.
func exchange(t *testing.T, client *ClientInstance, server *ServerInstance) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	upload := make([]byte, 20000)
	download := make([]byte, 30000)
	rand.Read(upload)
	rand.Read(download)

	serverErr := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			serverErr <- err
			return
		}
		defer conn.Close()
		encrypted, err := server.Handshake(conn, nil)
		if err != nil {
			serverErr <- err
			return
		}
		received := make([]byte, len(upload))
		_, err = io.ReadFull(encrypted, received)
		if err != nil {
			serverErr <- err
			return
		}
		if !bytes.Equal(received, upload) {
			serverErr <- io.ErrUnexpectedEOF
			return
		}
		_, err = encrypted.Write(download)
		serverErr <- err
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	encrypted, err := client.Handshake(conn)
	require.NoError(t, err)
	_, err = encrypted.Write(upload)
	require.NoError(t, err)
	received := make([]byte, len(download))
	_, err = io.ReadFull(encrypted, received)
	require.NoError(t, err)
	require.Equal(t, download, received)
	require.NoError(t, <-serverErr)
}

func TestVLESSEncryptionRoundTrip(t *testing.T) {
	t.Parallel()
	for _, xorMode := range []string{"native", "xorpub", "random"} {
		for _, auth := range []string{"x25519", "mlkem768", "relay"} {
			t.Run(xorMode+"/"+auth, func(t *testing.T) {
				t.Parallel()
				var pairs []testKeyPair
				switch auth {
				case "x25519":
					pairs = []testKeyPair{newX25519Pair(t)}
				case "mlkem768":
					pairs = []testKeyPair{newMLKEMPair(t)}
				case "relay":
					pairs = []testKeyPair{newX25519Pair(t), newMLKEMPair(t), newX25519Pair(t)}
				}
				var serverKeys, clientKeys []string
				for _, pair := range pairs {
					serverKeys = append(serverKeys, pair.server)
					clientKeys = append(clientKeys, pair.client)
				}
				padding := "100-111-1111.75-0-111.50-0-3333"
				server, err := NewServer("mlkem768x25519plus." + xorMode + ".600s." + padding + "." + strings.Join(serverKeys, "."))
				require.NoError(t, err)
				defer server.Close()
				client, err := NewClient("mlkem768x25519plus." + xorMode + ".0rtt." + padding + "." + strings.Join(clientKeys, "."))
				require.NoError(t, err)

				// 1-RTT, which issues a ticket.
				exchange(t, client, server)
				client.RWLock.RLock()
				hasTicket := client.Ticket != nil && time.Now().Before(client.Expire)
				client.RWLock.RUnlock()
				require.True(t, hasTicket)
				// 0-RTT with the ticket.
				exchange(t, client, server)
			})
		}
	}
}

func TestVLESSEncryption1RTTOnly(t *testing.T) {
	t.Parallel()
	pair := newMLKEMPair(t)
	server, err := NewServer("mlkem768x25519plus.random.0s." + pair.server)
	require.NoError(t, err)
	defer server.Close()
	client, err := NewClient("mlkem768x25519plus.random.1rtt." + pair.client)
	require.NoError(t, err)
	exchange(t, client, server)
	exchange(t, client, server)
}

func TestVLESSDecryptionParse(t *testing.T) {
	t.Parallel()
	server, err := NewServer("none")
	require.NoError(t, err)
	require.Nil(t, server)

	pair := newX25519Pair(t)
	server, err = NewServer("mlkem768x25519plus.xorpub.300-600s.100-111-1111." + pair.server)
	require.NoError(t, err)
	defer server.Close()
	require.Equal(t, uint32(1), server.XorMode)
	require.Equal(t, int64(300), server.SecondsFrom)
	require.Equal(t, int64(600), server.SecondsTo)
	require.Len(t, server.PaddingLens, 1)

	_, err = NewServer("mlkem768x25519plus.native.600s." + pair.client[:10])
	require.Error(t, err)
	_, err = NewServer("mlkem768x25519plus.bad.600s." + pair.server)
	require.Error(t, err)
	_, err = NewServer("mlkem768x25519plus.native.xs." + pair.server)
	require.Error(t, err)
}
