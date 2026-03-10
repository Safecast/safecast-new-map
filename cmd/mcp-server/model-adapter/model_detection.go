package modeladapter

import (
	"context"
	"net/http"
	"strings"
)

// ModelDetectionMiddleware is an HTTP middleware that inspects requests and
// determines which AI model is calling the MCP server. The detected model is
// stored in the request context under ctxKeyModelName.
func ModelDetectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model := detectModel(r)
		ctx := r.Context()
		ctx = contextWithModel(ctx, model)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// detectModel inspects various signals on the HTTP request to guess the model.
// For now it examines the User-Agent and an optional X-AI-Model header.
func detectModel(r *http.Request) ModelName {
	// custom header takes precedence
	if hdr := r.Header.Get("X-AI-Model"); hdr != "" {
		return normalizeModelName(hdr)
	}
	ua := r.UserAgent()
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "claude"):
		return ModelClaude
	case strings.Contains(ua, "kimi"):
		return ModelKimi
	case strings.Contains(ua, "qwen"):
		// qwen, qwen2, etc.
		return ModelQwen
	case strings.Contains(ua, "gpt") || strings.Contains(ua, "chatgpt") || strings.Contains(ua, "openai"):
		return ModelGPT
	}
	return ModelUnknown
}

// normalizeModelName reduces various header values to the canonical constant.
func normalizeModelName(raw string) ModelName {
	raw = strings.ToLower(raw)
	switch {
	case strings.Contains(raw, "claude"):
		return ModelClaude
	case strings.Contains(raw, "kimi"):
		return ModelKimi
	case strings.Contains(raw, "qwen"):
		return ModelQwen
	case strings.Contains(raw, "gpt") || strings.Contains(raw, "chatgpt") || strings.Contains(raw, "openai"):
		return ModelGPT
	}
	return ModelUnknown
}

// ctxKeyModelName is used as the context key for storing the detected model.
type ctxKeyModelName struct{}

// contextWithModel returns a new context containing the given model.
func contextWithModel(ctx context.Context, model ModelName) context.Context {
	return context.WithValue(ctx, ctxKeyModelName{}, model)
}

// ModelFromContext retrieves the model from the context, defaulting to
// ModelUnknown if none is present.
func ModelFromContext(ctx context.Context) ModelName {
	if v := ctx.Value(ctxKeyModelName{}); v != nil {
		if m, ok := v.(ModelName); ok {
			return m
		}
	}
	return ModelUnknown
}