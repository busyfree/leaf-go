package server

import (
	"net"
	"net/http"
	"strings"

	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/ctxkit"
)

func initCROSHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	corsHeaders := conf.GetString("CORS_ORIGIN_HEADERS")
	suffixs := conf.GetStringSlice("CORS_ORIGIN_SUFFIX")
	if len(corsHeaders) == 0 {
		corsHeaders = "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range"
	}

	if r.Method == http.MethodOptions {
		if len(suffixs) > 0 {
			for _, suffix := range suffixs {
				if len(suffix) > 0 && strings.Contains(strings.ToLower(origin), strings.ToLower(suffix)) {
					w.Header().Add("Access-Control-Allow-Origin", origin)
					w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
					w.Header().Add("Access-Control-Allow-Headers", corsHeaders)
					w.Header().Add("Access-Control-Max-Age", "1728000")
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Content-Length", "0")
					w.WriteHeader(http.StatusNoContent)
					return true
				}
			}
		}
	}
	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		if len(suffixs) > 0 {
			for _, suffix := range suffixs {
				if len(suffix) > 0 && strings.Contains(strings.ToLower(origin), strings.ToLower(suffix)) {
					w.Header().Add("Access-Control-Allow-Origin", origin)
					w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
					w.Header().Add("Access-Control-Allow-Headers", corsHeaders)
					w.Header().Add("Access-Control-Expose-Headers", "Content-Length,Content-Range")
				}
			}
		}
	}
	return false
}

func initRequestHeaders(w http.ResponseWriter, req *http.Request) *http.Request {
	ip := req.Header.Get("X-Real-IP")
	if ip == "" {
		ip = req.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = getIpFromRemoteAddr(req.RemoteAddr)
	}
	ctx := req.Context()
	ctx = ctxkit.WithUserIP(ctx, ip)
	req = req.WithContext(ctx)
	return req
}

func getIpFromRemoteAddr(remoteAddr string) (ip string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err == nil {
		return tcpAddr.IP.String()
	} else {
		tcpAddr, err = net.ResolveTCPAddr("tcp", remoteAddr+":80")
		if err == nil {
			return tcpAddr.IP.String()
		} else {
		}
		return ""
	}
}
