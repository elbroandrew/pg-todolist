package router

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"pg-todolist/internal/contextkeys"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
http.CloseNotifier - это устаревший интерфейс в Go, который использовался для отслеживания разрыва соединения клиентом. В современных версиях Go он deprecated, но некоторые библиотеки (включая старые версии httputil.ReverseProxy) все еще его используют.

Проблема: httptest.ResponseRecorder не реализует этот интерфейс, поэтому при использовании в интерационных тестах возникает panic.

*/

// ReverseProxyConfig - конфигурация для reverse proxy
type ReverseProxyConfig struct {
    Timeout           time.Duration
    MaxIdleConns      int
    IdleConnTimeout   time.Duration
    DisableCompression bool
}

// DefaultReverseProxyConfig - дефолтная конфигурация
func DefaultReverseProxyConfig() ReverseProxyConfig {
    return ReverseProxyConfig{
        Timeout:         30 * time.Second,
        MaxIdleConns:    100,
        IdleConnTimeout: 90 * time.Second,
    }
}

// CustomReverseProxy - production-ready reverse proxy
type CustomReverseProxy struct {
    target    *url.URL
    transport http.RoundTripper
    config    ReverseProxyConfig
}

func NewCustomReverseProxy(target string, config ...ReverseProxyConfig) (*CustomReverseProxy, error) {
    targetURL, err := url.Parse(target)
    if err != nil {
        return nil, err
    }

    cfg := DefaultReverseProxyConfig()
    if len(config) > 0 {
        cfg = config[0]
    }

    transport := &http.Transport{
        MaxIdleConns:        cfg.MaxIdleConns,
        IdleConnTimeout:     cfg.IdleConnTimeout,
        DisableCompression:  cfg.DisableCompression,
    }

    return &CustomReverseProxy{
        target:    targetURL,
        transport: transport,
        config:    cfg,
    }, nil
}

// ServeHTTP реализует http.Handler interface
func (p *CustomReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), p.config.Timeout)
    defer cancel()

    proxyReq := p.createProxyRequest(r.WithContext(ctx))
    if proxyReq == nil {
        http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
        return
    }

    // Выполняем запрос
    resp, err := p.transport.RoundTrip(proxyReq)
    if err != nil {
        log.Printf("Reverse proxy error: %v", err)
        http.Error(w, "Service unavailable", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // Копируем ответ
    p.copyResponse(w, resp)
}

// createProxyRequest создает прокси-запрос к целевому сервису
func (p *CustomReverseProxy) createProxyRequest(r *http.Request) *http.Request {
    // Строим целевой URL
    targetPath := singleJoiningSlash(p.target.Path, r.URL.Path)
    proxyURL := *p.target
    proxyURL.Path = targetPath
    proxyURL.RawQuery = r.URL.RawQuery

    // Создаем новый запрос
    proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, proxyURL.String(), r.Body)
    if err != nil {
        return nil
    }

    // Копируем заголовки
    p.copyHeaders(r.Header, proxyReq.Header)

    // Устанавливаем служебные заголовки
    p.setServiceHeaders(proxyReq, r)

    return proxyReq
}

// copyHeaders копирует заголовки из исходного запроса в прокси-запрос
func (p *CustomReverseProxy) copyHeaders(src, dst http.Header) {
    for key, values := range src {
        // Пропускаем некоторые заголовки которые не нужно проксировать
        if key == "Authorization" || key == "Cookie" {
            continue
        }
        for _, value := range values {
            dst.Add(key, value)
        }
    }
}

// setServiceHeaders устанавливает служебные заголовки для внутренней коммуникации
func (p *CustomReverseProxy) setServiceHeaders(proxyReq *http.Request, originalReq *http.Request) {
    // Устанавливаем User-ID из контекста
    if userIDVal := originalReq.Context().Value(contextkeys.UserIDKey); userIDVal != nil {
        if userID, ok := userIDVal.(uint); ok {
            proxyReq.Header.Set("X-User-ID", strconv.FormatUint(uint64(userID), 10))
        }
    }

    // Устанавливаем дополнительные заголовки для трейсинга
    proxyReq.Header.Set("X-Forwarded-Host", originalReq.Host)
    proxyReq.Header.Set("X-Forwarded-Proto", "http") // В проде нужно определять схему
    
    // Удаляем клиентские auth заголовки
    proxyReq.Header.Del("Authorization")
    proxyReq.Header.Del("Cookie")
    
    proxyReq.Host = p.target.Host
}

// copyResponse копирует ответ от целевого сервиса клиенту
func (p *CustomReverseProxy) copyResponse(w http.ResponseWriter, resp *http.Response) {
    // Копируем заголовки
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    // Устанавливаем статус код
    w.WriteHeader(resp.StatusCode)

    // Копируем тело
    if _, err := io.Copy(w, resp.Body); err != nil {
        log.Printf("Failed to copy response body: %v", err)
    }
}

// GinHandler создает gin.HandlerFunc из CustomReverseProxy
func (p *CustomReverseProxy) GinHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Добавляем userID в контекст запроса
        if userID, exists := c.Get("userID"); exists {
            ctx := context.WithValue(c.Request.Context(), contextkeys.UserIDKey, userID)
            c.Request = c.Request.WithContext(ctx)
        }
        
        p.ServeHTTP(c.Writer, c.Request)
    }
}

// Вспомогательная функция для объединения путей
func singleJoiningSlash(a, b string) string {
    aslash := strings.HasSuffix(a, "/")
    bslash := strings.HasPrefix(b, "/")
    switch {
    case aslash && bslash:
        return a + b[1:]
    case !aslash && !bslash:
        if b == "" {
            return a
        }
        return a + "/" + b
    }
    return a + b
}