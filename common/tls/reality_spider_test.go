//go:build with_utls

package tls

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseRealitySpider(t *testing.T) {
	t.Parallel()
	spider, err := parseRealitySpider("")
	require.NoError(t, err)
	require.Equal(t, "/", spider.path)
	require.Equal(t, [10]int64{}, spider.y)

	spider, err = parseRealitySpider("/search?q=1&p=10-20&c=3&t=2-4&i=50-100&r=1000-2000")
	require.NoError(t, err)
	require.Equal(t, "/search?q=1", spider.path)
	require.Equal(t, [10]int64{10, 20, 3, 3, 2, 4, 50, 100, 1000, 2000}, spider.y)

	_, err = parseRealitySpider("search")
	require.Error(t, err)
}

func TestRealitySpiderCrawl(t *testing.T) {
	t.Parallel()
	type visit struct {
		path    string
		referer string
		padding int
		ua      string
	}
	var (
		access sync.Mutex
		visits []visit
	)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("padding")
		access.Lock()
		visits = append(visits, visit{r.URL.Path, r.Referer(), len(cookie.Value), r.UserAgent()})
		access.Unlock()
		w.Write([]byte(`<a href="/about">a</a><a href="https://example.com/news">b</a><a href="/style.css">c</a>`))
	}))
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	conn, err := tls.Dial("tcp", server.Listener.Addr().String(), &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2"},
	})
	require.NoError(t, err)

	spider, err := parseRealitySpider("/start?p=5-6&c=2&t=2&i=1&r=10")
	require.NoError(t, err)
	start := time.Now()
	spider.run(conn, "example.com")
	require.GreaterOrEqual(t, time.Since(start), 10*time.Millisecond)

	require.Eventually(t, func() bool {
		access.Lock()
		defer access.Unlock()
		// first request plus concurrency(2) * times(2)
		return len(visits) >= 5
	}, 5*time.Second, 10*time.Millisecond)

	access.Lock()
	defer access.Unlock()
	require.Equal(t, "/start", visits[0].path)
	require.True(t, strings.Contains(visits[0].ua, "Mozilla/5.0"), visits[0].ua)
	for _, v := range visits {
		require.Contains(t, []int{5}, v.padding)
		// .css links are never followed
		require.NotEqual(t, "/style.css", v.path)
	}
	for _, v := range visits[1:] {
		require.NotEmpty(t, v.referer)
	}
}
