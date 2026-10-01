//go:build with_utls

// This file is derived from Xray-core (transport/internet/reality/reality.go
// and infra/conf/transport_security.go), licensed under the Mozilla Public
// License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.

package tls

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sagernet/sing-box/common/browserheaders"
	E "github.com/sagernet/sing/common/exceptions"

	"golang.org/x/net/http2"
)

// realitySpider is Xray's spiderX/spiderY: when the server answers with a real
// certificate, keep crawling the target site over the already established
// connection so that the probe looks like a regular browser visit.
type realitySpider struct {
	// path is spiderX without the p/c/t/i/r control parameters.
	path string
	// y is Xray's SpiderY: inclusive-ish [from, to] pairs for
	// padding, concurrency, times, interval (ms) and return delay (ms).
	y [10]int64
}

func parseRealitySpider(spiderX string) (*realitySpider, error) {
	if spiderX == "" {
		spiderX = "/"
	}
	if spiderX[0] != '/' {
		return nil, E.New("invalid spider_x: ", spiderX)
	}
	u, err := url.Parse(spiderX)
	if err != nil {
		return nil, E.Cause(err, "invalid spider_x")
	}
	spider := &realitySpider{}
	query := u.Query()
	parse := func(param string, index int) {
		if value := query.Get(param); value != "" {
			parts := strings.Split(value, "-")
			if len(parts) == 1 {
				spider.y[index], _ = strconv.ParseInt(parts[0], 10, 64)
				spider.y[index+1], _ = strconv.ParseInt(parts[0], 10, 64)
			} else {
				spider.y[index], _ = strconv.ParseInt(parts[0], 10, 64)
				spider.y[index+1], _ = strconv.ParseInt(parts[1], 10, 64)
			}
		}
		query.Del(param)
	}
	parse("p", 0) // padding
	parse("c", 2) // concurrency
	parse("t", 4) // times
	parse("i", 6) // interval
	parse("r", 8) // return
	u.RawQuery = query.Encode()
	spider.path = u.String()
	return spider, nil
}

var (
	realitySpiderHref = regexp.MustCompile(`href="([/h].*?)"`)
	realitySpiderDot  = []byte(".")
)

// realitySpiderPaths is shared between connections, as in Xray.
var realitySpiderPaths struct {
	sync.Mutex
	paths map[string]map[string]struct{}
}

// realityRandBetween mirrors Xray's crypto.RandBetween.
func realityRandBetween(from int64, to int64) int64 {
	if from > to {
		from, to = to, from
	}
	if d := to - from; d == 0 || d == 1 {
		return from
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(to-from))
	return from + n.Int64()
}

func realitySpiderPathLocked(paths map[string]struct{}) string {
	stopAt := int(realityRandBetween(0, int64(len(paths)-1)))
	i := 0
	for path := range paths {
		if i == stopAt {
			return path
		}
		i++
	}
	return "/"
}

// run crawls the target over conn in the background and then waits for the
// configured return delay. The connection is intentionally not closed by the
// crawler, matching Xray.
func (s *realitySpider) run(conn net.Conn, serverName string) {
	client := &http.Client{
		Transport: &http2.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string, config *tls.Config) (net.Conn, error) {
				return conn, nil
			},
		},
	}
	prefix := []byte("https://" + serverName)
	realitySpiderPaths.Lock()
	if realitySpiderPaths.paths == nil {
		realitySpiderPaths.paths = make(map[string]map[string]struct{})
	}
	paths := realitySpiderPaths.paths[serverName]
	if paths == nil {
		paths = make(map[string]struct{})
		paths[s.path] = struct{}{}
		realitySpiderPaths.paths[serverName] = paths
	}
	firstURL := string(prefix) + realitySpiderPathLocked(paths)
	realitySpiderPaths.Unlock()
	get := func(first bool) {
		var request *http.Request
		if first {
			request, _ = http.NewRequest("GET", firstURL, nil)
		} else {
			realitySpiderPaths.Lock()
			request, _ = http.NewRequest("GET", string(prefix)+realitySpiderPathLocked(paths), nil)
			realitySpiderPaths.Unlock()
		}
		if request == nil {
			return
		}
		browserheaders.Apply(request.Header, "nav")
		times := 1
		if !first {
			times = int(realityRandBetween(s.y[4], s.y[5]))
		}
		for j := 0; j < times; j++ {
			if !first && j == 0 {
				request.Header.Set("Referer", firstURL)
			}
			request.AddCookie(&http.Cookie{Name: "padding", Value: strings.Repeat("0", int(realityRandBetween(s.y[0], s.y[1])))})
			response, err := client.Do(request)
			if err != nil {
				break
			}
			request.Header.Set("Referer", request.URL.String())
			body, err := io.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				break
			}
			realitySpiderPaths.Lock()
			for _, match := range realitySpiderHref.FindAllSubmatch(body, -1) {
				match[1] = bytes.TrimPrefix(match[1], prefix)
				if !bytes.Contains(match[1], realitySpiderDot) {
					paths[string(match[1])] = struct{}{}
				}
			}
			request.URL.Path = realitySpiderPathLocked(paths)
			realitySpiderPaths.Unlock()
			if !first {
				time.Sleep(time.Duration(realityRandBetween(s.y[6], s.y[7])) * time.Millisecond)
			}
		}
	}
	go func() {
		get(true)
		concurrency := int(realityRandBetween(s.y[2], s.y[3]))
		for i := 0; i < concurrency; i++ {
			go get(false)
		}
	}()
	time.Sleep(time.Duration(realityRandBetween(s.y[8], s.y[9])) * time.Millisecond)
}
