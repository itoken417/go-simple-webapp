//go:build release

package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"

	config "github.com/itoken417/go-simple-webapp/configs"
	"github.com/itoken417/goutils/logger"
)

// trimStack は debug.Stack() 先頭の固定フレーム（Stack/RecoveryMiddleware/panic）を除去する
func trimStack(stack []byte) []byte {
	lines := strings.Split(string(stack), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "panic(") && i+2 < len(lines) {
			return []byte(strings.Join(lines[i+2:], "\n"))
		}
	}
	return stack
}

// buildEnv は CGI 相当の環境変数マップを返す
func buildEnv(r *http.Request) []string {
	remoteAddr, remotePort, _ := net.SplitHostPort(r.RemoteAddr)

	serverName, serverPort, err := net.SplitHostPort(r.Host)
	if err != nil {
		serverName = r.Host
		if r.TLS != nil {
			serverPort = "443"
		} else {
			serverPort = "80"
		}
	}

	https := ""
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		https = "ON"
	}

	env := []string{
		fmt.Sprintf("HTTPS               = %s", https),
		fmt.Sprintf("PATH_INFO           = %s", r.URL.Path),
		fmt.Sprintf("QUERY_STRING        = %s", r.URL.RawQuery),
		fmt.Sprintf("REMOTE_ADDR         = %s", remoteAddr),
		fmt.Sprintf("REMOTE_PORT         = %s", remotePort),
		fmt.Sprintf("REQUEST_METHOD      = %s", r.Method),
		fmt.Sprintf("REQUEST_URI         = %s", r.RequestURI),
		fmt.Sprintf("SCRIPT_NAME         = "),
		fmt.Sprintf("SERVER_NAME         = %s", serverName),
		fmt.Sprintf("SERVER_PORT         = %s", serverPort),
		fmt.Sprintf("SERVER_PROTOCOL     = %s", r.Proto),
	}

	httpHeaders := make([]string, 0, len(r.Header))
	for k, vs := range r.Header {
		key := "HTTP_" + strings.ToUpper(strings.ReplaceAll(k, "-", "_"))
		httpHeaders = append(httpHeaders, fmt.Sprintf("%-32s= %s", key, strings.Join(vs, ", ")))
	}
	sort.Strings(httpHeaders)

	return append(env, httpHeaders...)
}

func handleError(w http.ResponseWriter, r *http.Request, err interface{}, stack []byte) {
	logger.Log("panic recovered:", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	go notifyError(r, err, stack)
}

func notifyError(r *http.Request, err interface{}, stack []byte) {
	cfg, cfgErr := config.Get()
	if cfgErr != nil {
		logger.Log("設定の読み込みに失敗:", cfgErr)
		return
	}

	var b strings.Builder
	fmt.Fprintf(&b, "System: %s\n", cfg.AppName)
	fmt.Fprintf(&b, "Error:  %v\n", err)
	fmt.Fprintf(&b, "\n=== Environment ===\n")
	for _, kv := range buildEnv(r) {
		fmt.Fprintf(&b, "%s\n", kv)
	}
	fmt.Fprintf(&b, "\n=== Stack Trace ===\n%s", trimStack(stack))

	subject := fmt.Sprintf("[ERROR][%s] %s %s", cfg.AppName, r.Method, r.URL.Path)
	if sendErr := errSender.Send([]string{cfg.ErrorTo}, subject, b.String()); sendErr != nil {
		logger.Log("エラーメール送信失敗:", sendErr)
	}
}
