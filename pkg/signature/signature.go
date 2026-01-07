/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc api接口签名
 */

package signature

import (
	"github.com/sufo/bailu-admin/app/config"
	time2 "github.com/sufo/bailu-admin/utils/time"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	delimiter = "|"
)

// 合法的 Methods
var methods = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodHead:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodConnect: true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

type Signature interface {
	i()

	// Generate 生成签名
	Generate(path string, method string, params url.Values) (authorization, date string, err error)

	// Verify 验证签名
	Verify(authorization, date string, path string, method string, params url.Values) (ok bool, err error)
}

// Generate
// path 请求的路径 (不附带 querystring)
func Generate(path string, method string, params url.Values) (authorization, date string, err error) {
	if path == "" {
		err = errors.New("path required")
		return
	}

	if method == "" {
		err = errors.New("method required")
		return
	}

	methodName := strings.ToUpper(method)
	if !methods[methodName] {
		err = errors.New("method param exception")
		return
	}

	// Date
	date = time2.CSTLayoutString()

	// Encode() 方法中自带 sorted by key
	sortParamsEncode, err := url.QueryUnescape(params.Encode())
	if err != nil {
		err = errors.Errorf("url QueryUnescape exception %v", err)
		return
	}

	// 加密字符串规则
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(path)
	buffer.WriteString(delimiter)
	buffer.WriteString(methodName)
	buffer.WriteString(delimiter)
	buffer.WriteString(sortParamsEncode)
	buffer.WriteString(delimiter)
	buffer.WriteString(date)

	s := config.Conf.Signature
	// 对数据进行 sha256 加密，并进行 base64 encode
	hash := hmac.New(sha256.New, []byte(s.Secret))
	hash.Write(buffer.Bytes())
	digest := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	authorization = fmt.Sprintf("%s %s", s.Key, digest)
	return
}

func Verify(authorization, date string, path string, method string, params url.Values) (ok bool, err error) {
	if date == "" {
		err = errors.New("date required")
		return
	}

	if path == "" {
		err = errors.New("path required")
		return
	}

	if method == "" {
		err = errors.New("method required")
		return
	}

	methodName := strings.ToUpper(method)
	if !methods[methodName] {
		err = errors.New("method param exception")
		return
	}

	ts, err := time2.ParseCSTInLocation(date)
	if err != nil {
		err = errors.New("date must follow '2006-01-02 15:04:05'")
		return
	}
	s := config.Conf.Signature
	if time2.SubInLocation(ts) > float64(s.TTL/time.Second) {
		err = errors.Errorf("date exceeds limit %v", s.TTL)
		return
	}

	// Encode() 方法中自带 sorted by key
	sortParamsEncode, err := url.QueryUnescape(params.Encode())
	if err != nil {
		err = errors.Errorf("url QueryUnescape exception %v", err)
		return
	}

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(path)
	buffer.WriteString(delimiter)
	buffer.WriteString(methodName)
	buffer.WriteString(delimiter)
	buffer.WriteString(sortParamsEncode)
	buffer.WriteString(delimiter)
	buffer.WriteString(date)

	// 对数据进行 hmac 加密，并进行 base64 encode
	hash := hmac.New(sha256.New, []byte(s.Secret))
	hash.Write(buffer.Bytes())
	digest := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	ok = authorization == fmt.Sprintf("%s %s", s.Key, digest)
	return
}
