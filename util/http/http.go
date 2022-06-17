// Package http 提供基础 http 客户端组件
// 内置以下功能：
// - logging
// - opentracing
// - prometheus
package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/busyfree/leaf-go/util/errors"
	"github.com/busyfree/leaf-go/util/log"
	"github.com/busyfree/leaf-go/util/metrics"
	"github.com/busyfree/leaf-go/util/trace"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type myClient struct {
	cli *http.Client
}

// Client http 客户端接口
type Client interface {
	// Do 发送单个 http 请求
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// NewClient 创建 Client 实例
func NewClient(timeout time.Duration) Client {
	return &myClient{
		cli: &http.Client{
			Timeout: timeout,
		},
	}
}

var digitsRE = regexp.MustCompile(`\b\d+\b`)

func (c *myClient) Do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DoHTTP")
	defer span.Finish()

	req = req.WithContext(ctx)

	trace.InjectTraceHeader(span.Context(), req)

	start := time.Now()
	resp, err = c.cli.Do(req)
	duration := time.Since(start)

	url := fmt.Sprintf("%s%s", req.URL.Host, req.URL.Path)
	var logger = log.Get(ctx)
	logger.Printf("myClientResp:%v", resp)
	logger.Printf("myClientReqErr:%v", err)
	status := http.StatusOK
	if err != nil {
		err = errors.Wrap(err)
		status = http.StatusGatewayTimeout
	} else {
		status = resp.StatusCode
	}
	log.Get(ctx).Debugf(
		"[%s] method:%s url:%s status:%d query:%s",
		strings.ToUpper(req.URL.Scheme),
		req.Method,
		url,
		status,
		req.URL.RawQuery,
	)

	span.SetTag(string(ext.Component), "http")
	span.SetTag(string(ext.HTTPUrl), url)
	span.SetTag(string(ext.HTTPMethod), req.Method)
	span.SetTag(string(ext.HTTPStatusCode), status)

	// url 中带有的纯数字替换成 %d，不然 prometheus 就炸了
	// /v123/4/56/foo => /v123/%d/%d/foo
	url = digitsRE.ReplaceAllString(url, "%d")

	metrics.HTTPDurationsSeconds.WithLabelValues(
		url,
		fmt.Sprint(status),
	).Observe(duration.Seconds())

	return
}

func DoHttpReq(ctx context.Context, urlStr string, params interface{}, options ...interface{}) (reqResp *http.Response, bodyByte []byte, err error) {
	if len(urlStr) == 0 {
		err = errors.Errorf("missing urlStr")
		return
	}
	var (
		method      = "GET"
		contentType = "application/json"
		headers     = make(map[string]string, 0)
	)
	if len(options) >= 1 {
		if methodVal, ok := options[0].(string); ok {
			method = methodVal
		}
	}
	if len(options) > 0 {
		for idx, v := range options {
			switch idx {
			case 0:
				if methodVal, ok := v.(string); ok {
					method = methodVal
				} else if headers, ok := v.(map[string]string); ok {
					if len(headers) > 0 {
						for k, v := range headers {
							if k == "Content-Type" {
								contentType = v
							}
							headers[k] = v
						}
					}
				}
			case 1:
				if headers, ok := v.(map[string]string); ok {
					if len(headers) > 0 {
						for k, v := range headers {
							if k == "Content-Type" {
								contentType = v
							}
							headers[k] = v
						}
					}
				}
			}
		}
	}
	var req *http.Request
	c := NewClient(time.Duration(30) * time.Second)
	if strings.ToUpper(method) == "GET" {
		if contentType == "application/json" {
			var ioReader *bytes.Reader
			if val, ok := params.(url.Values); ok {
				ioReader = bytes.NewReader([]byte(val.Encode()))
			} else if val, ok := params.(string); ok {
				ioReader = bytes.NewReader([]byte(val))
			} else if val, ok := params.([]byte); ok {
				ioReader = bytes.NewReader(val)
			}
			req, err = http.NewRequest("GET", urlStr, ioReader)
		} else {
			if val, ok := params.(url.Values); ok {
				urlStr += "?" + val.Encode()
			} else if val, ok := params.(string); ok {
				urlStr += "?" + val
			} else if val, ok := params.([]byte); ok {
				urlStr += "?" + string(val)
			}
			req, err = http.NewRequest("GET", urlStr, nil)
		}
	} else {
		var ioReader *bytes.Reader
		if val, ok := params.(url.Values); ok {
			ioReader = bytes.NewReader([]byte(val.Encode()))
		} else if val, ok := params.(string); ok {
			ioReader = bytes.NewReader([]byte(val))
		} else if val, ok := params.([]byte); ok {
			ioReader = bytes.NewReader(val)
		}
		req, err = http.NewRequest("POST", urlStr, ioReader)
	}
	if err != nil {
		return
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Treedom/1.0.0 Go/1.0.0")
	reqResp, err = c.Do(ctx, req)
	if err != nil {
		return
	}
	if reqResp.Body != nil {
		defer reqResp.Body.Close()
		bodyByte, err = ioutil.ReadAll(reqResp.Body)
	}
	return
}

func DoHttpFormReq(ctx context.Context, urlStr string, params interface{}, options ...interface{}) (reqResp *http.Response, bodyByte []byte, err error) {
	if len(urlStr) == 0 {
		err = errors.Errorf("missing urlStr")
		return
	}
	var (
		method = "GET"
	)
	if len(options) >= 1 {
		if methodVal, ok := options[0].(string); ok {
			method = methodVal
		}
	}
	var req *http.Request
	c := NewClient(time.Duration(30) * time.Second)
	if strings.ToUpper(method) == "GET" {
		if val, ok := params.(url.Values); ok {
			urlStr += "?" + val.Encode()
		} else if val, ok := params.(string); ok {
			urlStr += "?" + val
		} else if val, ok := params.([]byte); ok {
			urlStr += "?" + string(val)
		}
		req, err = http.NewRequest("GET", urlStr, nil)
	} else {
		var ioReader *bytes.Reader
		if val, ok := params.(url.Values); ok {
			ioReader = bytes.NewReader([]byte(val.Encode()))
		} else if val, ok := params.(string); ok {
			ioReader = bytes.NewReader([]byte(val))
		} else if val, ok := params.([]byte); ok {
			ioReader = bytes.NewReader(val)
		}
		req, err = http.NewRequest("POST", urlStr, ioReader)
	}
	if err != nil {
		return
	}
	if len(options) >= 2 {
		if headers, ok := options[0].(map[string]string); ok {
			if len(headers) > 0 {
				for k, v := range headers {
					req.Header.Add(k, v)
				}
			}
		}
		if headers, ok := options[1].(map[string]string); ok {
			if len(headers) > 0 {
				for k, v := range headers {
					req.Header.Set(k, v)
				}
			}
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Treedom/1.0.0 Go/1.0.0")
	reqResp, err = c.Do(ctx, req)
	if err != nil {
		return
	}
	if reqResp.Body != nil {
		defer reqResp.Body.Close()
		bodyByte, err = ioutil.ReadAll(reqResp.Body)
	}
	return
}

func DoHttpJsonReq(ctx context.Context, urlStr string, params interface{}, options ...interface{}) (reqResp *http.Response, bodyByte []byte, err error) {
	if len(urlStr) == 0 {
		err = errors.Errorf("missing urlStr")
		return
	}
	var (
		method = "GET"
	)
	if len(options) >= 1 {
		if methodVal, ok := options[0].(string); ok {
			method = methodVal
		}
	}
	var req *http.Request
	c := NewClient(time.Duration(30) * time.Second)
	if strings.ToUpper(method) == "GET" {
		if val, ok := params.(url.Values); ok {
			urlStr += "?" + val.Encode()
		} else if val, ok := params.(string); ok {
			urlStr += "?" + val
		} else if val, ok := params.([]byte); ok {
			urlStr += "?" + string(val)
		}
		req, err = http.NewRequest("GET", urlStr, nil)
	} else {
		var ioReader *bytes.Reader
		if val, ok := params.(url.Values); ok {
			ioReader = bytes.NewReader([]byte(val.Encode()))
		} else if val, ok := params.(string); ok {
			ioReader = bytes.NewReader([]byte(val))
		} else if val, ok := params.([]byte); ok {
			ioReader = bytes.NewReader(val)
		}
		req, err = http.NewRequest("POST", urlStr, ioReader)
	}
	if err != nil {
		return
	}
	if len(options) >= 2 {
		if headers, ok := options[0].(map[string]string); ok {
			if len(headers) > 0 {
				for k, v := range headers {
					req.Header.Add(k, v)
				}
			}
		}
		if headers, ok := options[1].(map[string]string); ok {
			if len(headers) > 0 {
				for k, v := range headers {
					req.Header.Add(k, v)
				}
			}
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Treedom/1.0.0 Go/1.0.0")
	reqResp, err = c.Do(ctx, req)
	if err != nil {
		return
	}
	if reqResp.Body != nil {
		defer reqResp.Body.Close()
		bodyByte, err = ioutil.ReadAll(reqResp.Body)
	}
	return
}

func DoHttpFormUploadReq(ctx context.Context, urlStr, filedName, filename string) (reqResp *http.Response, bodyByte []byte, err error) {
	if len(urlStr) == 0 {
		err = errors.Errorf("missing urlStr")
		return
	}
	var req *http.Request
	c := NewClient(time.Duration(30) * time.Second)

	f, err := os.Open(filename)
	defer f.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile(filedName, filename)
	if err != nil {
		return
	}
	_, err = io.Copy(fw, f)
	if err != nil {
		return
	}

	fields := map[string]string{
		filedName: filename,
	}
	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}
	err = writer.Close() // close writer before POST request
	if err != nil {
		return
	}
	req, err = http.NewRequest("POST", urlStr, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Treedom/1.0.0 Go/1.0.0")
	reqResp, err = c.Do(ctx, req)
	if err != nil {
		return
	}
	if reqResp.Body != nil {
		defer reqResp.Body.Close()
		bodyByte, err = ioutil.ReadAll(reqResp.Body)
	}
	return
}
