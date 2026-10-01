// This file is derived from Xray-core (transport/internet/splithttp/upload_queue.go),
// licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.

package v2rayxhttp

// uploadQueue is a specialized priority queue + channel to reorder generic
// packets by a sequence number.

import (
	"container/heap"
	"io"
	"sync"
	"sync/atomic"

	E "github.com/sagernet/sing/common/exceptions"
)

// doneInstance is Xray's signal/done.Instance.
type doneInstance struct {
	access sync.Once
	c      chan struct{}
}

func newDone() *doneInstance {
	return &doneInstance{c: make(chan struct{})}
}

func (d *doneInstance) Done() bool {
	select {
	case <-d.c:
		return true
	default:
		return false
	}
}

func (d *doneInstance) Wait() <-chan struct{} {
	return d.c
}

func (d *doneInstance) Close() error {
	d.access.Do(func() {
		close(d.c)
	})
	return nil
}

type uploadPacket struct {
	Reader  *httpServerConn
	Payload []byte
	Seq     uint64
}

type uploadQueue struct {
	reader        atomic.Pointer[httpServerConn]
	pushedPackets chan uploadPacket
	heap          uploadHeap
	nextSeq       uint64
	maxPackets    int
	closed        *doneInstance
}

func newUploadQueue(maxPackets int) *uploadQueue {
	return &uploadQueue{
		pushedPackets: make(chan uploadPacket, maxPackets),
		heap:          uploadHeap{},
		nextSeq:       0,
		closed:        newDone(),
		maxPackets:    maxPackets,
	}
}

func (h *uploadQueue) Push(p uploadPacket) error {
	if h.reader.Load() != nil || (p.Reader != nil && !h.reader.CompareAndSwap(nil, p.Reader)) {
		return E.New("h.reader already exists")
	}
	select {
	case h.pushedPackets <- p: // no panic
		if h.closed.Done() {
			return E.New("packet queue closed")
		}
		return nil
	case <-h.closed.Wait():
		return E.New("packet queue closed")
	}
}

func (h *uploadQueue) Close() error {
	h.closed.Close()
	if reader := h.reader.Load(); reader != nil {
		return reader.Close()
	}
	return nil
}

func (h *uploadQueue) Read(b []byte) (int, error) {
	if reader := h.reader.Load(); reader != nil {
		return reader.Read(b)
	}

	if h.closed.Done() {
		return 0, io.EOF
	}

	if len(h.heap) == 0 {
		select {
		case p := <-h.pushedPackets:
			if p.Reader != nil {
				return p.Reader.Read(b)
			}
			heap.Push(&h.heap, p)
		case <-h.closed.Wait():
			return 0, io.EOF
		}
	}

	for len(h.heap) > 0 {
		packet := heap.Pop(&h.heap).(uploadPacket)
		n := 0

		if packet.Seq == h.nextSeq {
			copy(b, packet.Payload)
			n = min(len(b), len(packet.Payload))

			if n < len(packet.Payload) {
				// partial read
				packet.Payload = packet.Payload[n:]
				heap.Push(&h.heap, packet)
			} else {
				h.nextSeq = packet.Seq + 1
			}

			return n, nil
		}

		// misordered packet
		if packet.Seq > h.nextSeq {
			if len(h.heap) > h.maxPackets {
				// the "reassembly buffer" is too large, and we want to
				// constrain memory usage somehow. let's tear down the
				// connection, and hope the application retries.
				return 0, E.New("packet queue is too large")
			}
			heap.Push(&h.heap, packet)
			select {
			case p := <-h.pushedPackets:
				heap.Push(&h.heap, p)
			case <-h.closed.Wait():
				return 0, io.EOF
			}
		}
	}

	return 0, nil
}

type uploadHeap []uploadPacket

func (h uploadHeap) Len() int           { return len(h) }
func (h uploadHeap) Less(i, j int) bool { return h[i].Seq < h[j].Seq }
func (h uploadHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *uploadHeap) Push(x any) {
	*h = append(*h, x.(uploadPacket))
}

func (h *uploadHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
