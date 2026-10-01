// This file is derived from Xray-core (transport/internet/splithttp/config.go
// and xpadding.go), licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.
//
// Server-side halves of the XHTTP configuration helpers. Like the client
// helpers, everything here is wire-visible and must stay identical to Xray.

package v2rayxhttp

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/http2/hpack"
)

// WriteResponseHeader writes Xray's CORS headers, which are required by the
// browser dialer and are otherwise harmless.
func (c *Config) WriteResponseHeader(writer http.ResponseWriter, requestMethod string, requestHeader http.Header) {
	if origin := requestHeader.Get("Origin"); origin == "" {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		// Chrome says: The value of the 'Access-Control-Allow-Origin' header in the response must not be the wildcard '*' when the request's credentials mode is 'include'.
		writer.Header().Set("Access-Control-Allow-Origin", origin)
	}

	if c.GetNormalizedSessionPlacement() == PlacementCookie ||
		c.GetNormalizedSeqPlacement() == PlacementCookie ||
		c.XPaddingPlacement == PlacementCookie ||
		c.GetNormalizedUplinkDataPlacement() == PlacementCookie {
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if requestMethod == "OPTIONS" {
		requestedMethod := requestHeader.Get("Access-Control-Request-Method")
		if requestedMethod != "" {
			writer.Header().Set("Access-Control-Allow-Methods", requestedMethod)
		} else {
			writer.Header().Set("Access-Control-Allow-Methods", "*")
		}

		requestedHeaders := requestHeader.Get("Access-Control-Request-Headers")
		if requestedHeaders == "" {
			writer.Header().Set("Access-Control-Allow-Headers", "*")
		} else {
			writer.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
		}
	}
}

func (c *Config) ExtractMetaFromRequest(req *http.Request, path string) (sessionId string, seqStr string) {
	sessionPlacement := c.GetNormalizedSessionPlacement()
	seqPlacement := c.GetNormalizedSeqPlacement()
	sessionKey := c.GetNormalizedSessionKey()
	seqKey := c.GetNormalizedSeqKey()

	var subpath []string
	pathPart := 0
	if sessionPlacement == PlacementPath || seqPlacement == PlacementPath {
		subpath = strings.Split(req.URL.Path[len(path):], "/")
	}

	switch sessionPlacement {
	case PlacementPath:
		if len(subpath) > pathPart {
			sessionId = subpath[pathPart]
			pathPart += 1
		}
	case PlacementQuery:
		sessionId = req.URL.Query().Get(sessionKey)
	case PlacementHeader:
		sessionId = req.Header.Get(sessionKey)
	case PlacementCookie:
		if cookie, e := req.Cookie(sessionKey); e == nil {
			sessionId = cookie.Value
		}
	}

	switch seqPlacement {
	case PlacementPath:
		if len(subpath) > pathPart {
			seqStr = subpath[pathPart]
			pathPart += 1
		}
	case PlacementQuery:
		seqStr = req.URL.Query().Get(seqKey)
	case PlacementHeader:
		seqStr = req.Header.Get(seqKey)
	case PlacementCookie:
		if cookie, e := req.Cookie(seqKey); e == nil {
			seqStr = cookie.Value
		}
	}

	return sessionId, seqStr
}

func ApplyPaddingToResponseCookie(writer http.ResponseWriter, name, value string) {
	if name == "" || value == "" {
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	})
}

func (c *Config) ApplyXPaddingToResponse(writer http.ResponseWriter, config XPaddingConfig) {
	placement := config.Placement.Placement

	if placement == PlacementHeader || placement == PlacementQueryInHeader {
		c.ApplyXPaddingToHeader(writer.Header(), config)
		return
	}

	paddingValue := GeneratePadding(config.Method, config.Length)

	switch placement {
	case PlacementCookie:
		ApplyPaddingToResponseCookie(writer, config.Placement.Key, paddingValue)
	}
}

func (c *Config) ExtractXPaddingFromRequest(req *http.Request, obfsMode bool) (string, string) {
	if req == nil {
		return "", ""
	}

	if !obfsMode {
		referrer := req.Header.Get("Referer")

		if referrer != "" {
			if referrerURL, err := url.Parse(referrer); err == nil {
				paddingValue := referrerURL.Query().Get("x_padding")
				paddingPlacement := PlacementQueryInHeader + "=Referer, key=x_padding"
				return paddingValue, paddingPlacement
			}
		} else {
			paddingValue := req.URL.Query().Get("x_padding")
			return paddingValue, PlacementQuery + ", key=x_padding"
		}
	}

	key := c.XPaddingKey
	header := c.XPaddingHeader

	if cookie, err := req.Cookie(key); err == nil {
		if cookie != nil && cookie.Value != "" {
			paddingValue := cookie.Value
			paddingPlacement := PlacementCookie + ", key=" + key
			return paddingValue, paddingPlacement
		}
	}

	headerValue := req.Header.Get(header)

	if headerValue != "" {
		if c.XPaddingPlacement == PlacementHeader {
			paddingPlacement := PlacementHeader + "=" + header
			return headerValue, paddingPlacement
		}

		if parsedURL, err := url.Parse(headerValue); err == nil {
			paddingPlacement := PlacementQueryInHeader + "=" + header + ", key=" + key

			return parsedURL.Query().Get(key), paddingPlacement
		}
	}

	queryValue := req.URL.Query().Get(key)

	if queryValue != "" {
		paddingPlacement := PlacementQuery + ", key=" + key
		return queryValue, paddingPlacement
	}

	return "", ""
}

func (c *Config) IsPaddingValid(paddingValue string, from, to int32, method PaddingMethod) bool {
	if paddingValue == "" {
		return false
	}
	if to <= 0 {
		r := c.GetNormalizedXPaddingBytes()
		from, to = r.From, r.To
	}

	switch method {
	case PaddingMethodRepeatX:
		n := int32(len(paddingValue))
		return n >= from && n <= to
	case PaddingMethodTokenish:
		const tolerance = int32(validationTolerance)

		n := int32(hpack.HuffmanEncodeLength(paddingValue))
		f := from - tolerance
		t := to + tolerance
		if f < 0 {
			f = 0
		}
		return n >= f && n <= t
	default:
		n := int32(len(paddingValue))
		return n >= from && n <= to
	}
}

// isValidHTTPHost mirrors Xray's internet.IsValidHTTPHost.
func isValidHTTPHost(request string, config string) bool {
	r := strings.ToLower(request)
	c := strings.ToLower(config)
	if strings.Contains(r, ":") {
		h, _, _ := net.SplitHostPort(r)
		return h == c
	}
	return r == c
}
