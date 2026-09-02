package cors

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"gnalloy.org/codec-http1"
	"gnalloy.org/gnalloy/channel"
)

const (
	headerOrigin           = "Origin"
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"
	headerRequestMethod    = "Access-Control-Request-Method"
	headerRequestHeaders   = "Access-Control-Request-Headers"
	forbiddenStatus        = 403
	preflightStatus        = 204
)

// Handler 提供 Netty CorsHandler 风格的入站预检和出站响应修饰。
type Handler struct {
	cfg normalizedConfig

	mu          sync.Mutex
	pending     []decision
	pendingHead int
}

type decision struct {
	origin string
}

// NewHandler 创建 CORS 处理器。
func NewHandler(cfg Config) (*Handler, error) {
	normalized, err := normalizeConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Handler{cfg: normalized}, nil
}

// ChannelRead 处理 HTTP/1 请求，并在预检请求上直接返回 CORS 响应。
func (h *Handler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	var req *http1.Request
	switch value := msg.(type) {
	case http1.Request:
		req = &value
	case *http1.Request:
		req = value
	}
	if req == nil {
		ctx.FireChannelRead(msg)
		return
	}
	originHeader := req.Headers.Get(headerOrigin)
	origin, allowed := h.allowedOrigin(originHeader)
	if isPreflight(*req) {
		err := h.writePreflight(ctx, *req, origin, allowed)
		req.Release()
		if err != nil {
			ctx.FireExceptionCaught(err)
		}
		return
	}
	if originHeader != "" && !allowed && h.cfg.shortCircuit {
		req.Release()
		if err := writeForbidden(ctx); err != nil {
			ctx.FireExceptionCaught(err)
		}
		return
	}
	h.enqueue(decision{origin: origin})
	ctx.FireChannelRead(msg)
}

// Write 为对应请求的 HTTP/1 响应补充 CORS 头。
func (h *Handler) Write(ctx *channel.HandlerContext, msg any) error {
	var resp *http1.Response
	writeValue := false
	switch value := msg.(type) {
	case http1.Response:
		resp = &value
		writeValue = true
	case *http1.Response:
		resp = value
	}
	if resp == nil {
		return ctx.Write(msg)
	}
	if d, ok := h.dequeue(); ok && d.origin != "" {
		if resp.Headers == nil {
			resp.Headers = http1.Headers{}
		}
		h.applyResponseHeaders(resp.Headers, d.origin, false)
	}
	if writeValue {
		return ctx.Write(*resp)
	}
	return ctx.Write(resp)
}

// HandlerRemoved 清理仍未匹配响应的请求决策。
func (h *Handler) HandlerRemoved(*channel.HandlerContext) error {
	h.clear()
	return nil
}

// ChannelInactive 清理状态并继续传播失活事件。
func (h *Handler) ChannelInactive(ctx *channel.HandlerContext) {
	h.clear()
	ctx.FireChannelInactive()
}

func (h *Handler) writePreflight(ctx *channel.HandlerContext, req http1.Request, origin string, allowed bool) error {
	if !allowed || !h.methodAllowed(req.Headers.Get(headerRequestMethod)) {
		return writeForbidden(ctx)
	}
	headers := http1.Headers{}
	h.applyResponseHeaders(headers, origin, true)
	if h.cfg.allowAnyHeader {
		if requested := req.Headers.Get(headerRequestHeaders); requested != "" {
			headers.Set(headerAllowHeaders, requested)
		}
	} else if len(h.cfg.allowedHeaders) > 0 {
		headers.Set(headerAllowHeaders, strings.Join(h.cfg.allowedHeaders, ", "))
	}
	if h.cfg.maxAge > 0 {
		headers.Set(headerMaxAge, strconv.FormatInt(int64(h.cfg.maxAge/time.Second), 10))
	}
	resp := http1.AcquireResponse()
	resp.StatusCode = preflightStatus
	resp.Headers = headers
	return ctx.Channel().WriteAndFlush(resp)
}

func (h *Handler) applyResponseHeaders(headers http1.Headers, origin string, preflight bool) {
	if headers == nil {
		return
	}
	headers.Set(headerAllowOrigin, origin)
	if h.cfg.allowCredentials {
		headers.Set(headerAllowCredentials, "true")
	}
	if len(h.cfg.exposedHeaders) > 0 && !preflight {
		headers.Set(headerExposeHeaders, strings.Join(h.cfg.exposedHeaders, ", "))
	}
	if preflight {
		headers.Set(headerAllowMethods, strings.Join(h.cfg.allowedMethods, ", "))
	}
}

func (h *Handler) allowedOrigin(origin string) (string, bool) {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return "", false
	}
	if h.cfg.allowAnyOrigin {
		if h.cfg.allowCredentials {
			return origin, true
		}
		return "*", true
	}
	_, ok := h.cfg.allowedOrigins[origin]
	if !ok {
		return "", false
	}
	return origin, true
}

func (h *Handler) methodAllowed(method string) bool {
	method = strings.ToUpper(strings.TrimSpace(method))
	_, ok := h.cfg.methodSet[method]
	return ok
}

func (h *Handler) enqueue(d decision) {
	h.mu.Lock()
	h.pending = append(h.pending, d)
	h.mu.Unlock()
}

func (h *Handler) dequeue() (decision, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.pendingHead >= len(h.pending) {
		return decision{}, false
	}
	d := h.pending[h.pendingHead]
	h.pending[h.pendingHead] = decision{}
	h.pendingHead++
	if h.pendingHead > 16 && h.pendingHead*2 >= len(h.pending) {
		copy(h.pending, h.pending[h.pendingHead:])
		h.pending = h.pending[:len(h.pending)-h.pendingHead]
		h.pendingHead = 0
	}
	return d, true
}

func (h *Handler) clear() {
	h.mu.Lock()
	h.pending = nil
	h.pendingHead = 0
	h.mu.Unlock()
}

func isPreflight(req http1.Request) bool {
	return strings.EqualFold(req.Method, "OPTIONS") && req.Headers.Get(headerOrigin) != "" && req.Headers.Get(headerRequestMethod) != ""
}

func writeForbidden(ctx *channel.HandlerContext) error {
	resp := http1.AcquireResponse()
	resp.StatusCode = forbiddenStatus
	resp.Headers = http1.Headers{}
	return ctx.Channel().WriteAndFlush(resp)
}
