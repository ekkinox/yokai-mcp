package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	yokaimcpservercontext "github.com/ekkinox/yokai-mcp/pkg/mcp/server/context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var _ MCPServerHooksProvider = (*DefaultMCPServerHooksProvider)(nil)

type MCPServerHooksProvider interface {
	Provide() *server.Hooks
}

type DefaultMCPServerHooksProvider struct {
	config           *config.Config
	requestsCounter  *prometheus.CounterVec
	requestsDuration *prometheus.HistogramVec
}

func NewDefaultMCPServerHooksProvider(registry prometheus.Registerer, config *config.Config) *DefaultMCPServerHooksProvider {
	namespace := Sanitize(config.GetString("modules.mcp.server.metrics.collect.namespace"))
	subsystem := Sanitize(config.GetString("modules.mcp.server.metrics.collect.subsystem"))

	buckets := prometheus.DefBuckets
	if bucketsConfig := config.GetString("modules.mcp.server.metrics.buckets"); bucketsConfig != "" {
		for _, s := range Split(bucketsConfig) {
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				buckets = append(buckets, f)
			}
		}
	}

	requestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "mcp_server_requests_total",
			Help:      "Number of processed MCP requests",
		},
		[]string{
			"method",
			"target",
			"status",
		},
	)

	requestsDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "mcp_server_requests_duration_seconds",
			Help:      "Time spent processing MCP requests",
			Buckets:   buckets,
		},
		[]string{
			"method",
			"target",
		},
	)

	registry.MustRegister(requestsCounter, requestsDuration)

	return &DefaultMCPServerHooksProvider{
		config:           config,
		requestsCounter:  requestsCounter,
		requestsDuration: requestsDuration,
	}
}

func (p *DefaultMCPServerHooksProvider) Provide() *server.Hooks {
	hooks := &server.Hooks{}

	traceRequest := p.config.GetBool("modules.mcp.server.trace.request")
	traceResponse := p.config.GetBool("modules.mcp.server.trace.response")
	traceExclusions := p.config.GetStringSlice("modules.mcp.server.trace.exclude")

	logRequest := p.config.GetBool("modules.mcp.server.log.request")
	logResponse := p.config.GetBool("modules.mcp.server.log.response")
	logExclusions := p.config.GetStringSlice("modules.mcp.server.log.exclude")

	metricsEnabled := p.config.GetBool("modules.mcp.server.metrics.collect.enabled")

	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		log.CtxLogger(ctx).Info().Str("mcpSessionID", session.SessionID()).Msg("MCP session registered")
	})

	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		latency := time.Since(yokaimcpservercontext.CtxStartTime(ctx))

		mcpMethod := string(method)

		spanNameSuffix := mcpMethod

		spanAttributes := []attribute.KeyValue{
			attribute.String("mcp.latency", latency.String()),
			attribute.String("mcp.method", mcpMethod),
		}

		logFields := map[string]interface{}{
			"mcpLatency": latency.String(),
			"mcpMethod":  mcpMethod,
		}

		metricTarget := ""

		jsonMessage, err := json.Marshal(message)
		if err == nil {
			if traceRequest {
				spanAttributes = append(spanAttributes, attribute.String("mcp.request", string(jsonMessage)))
			}

			if logRequest {
				logFields["mcpRequest"] = string(jsonMessage)
			}
		}

		jsonResult, err := json.Marshal(result)
		if err == nil {
			if traceResponse {
				spanAttributes = append(spanAttributes, attribute.String("mcp.response", string(jsonResult)))
			}

			if logResponse {
				logFields["mcpResponse"] = string(jsonResult)
			}
		}

		switch method {
		case mcp.MethodResourcesRead:
			if req, ok := message.(*mcp.ReadResourceRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.URI)
				spanAttributes = append(spanAttributes, attribute.String("mcp.resource", req.Params.URI))
				logFields["mcpResourceURI"] = req.Params.URI
				metricTarget = req.Params.URI
			}
		case mcp.MethodPromptsGet:
			if req, ok := message.(*mcp.GetPromptRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.prompt", req.Params.Name))
				logFields["mcpPrompt"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		case mcp.MethodToolsCall:
			if req, ok := message.(*mcp.CallToolRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.tool", req.Params.Name))
				logFields["mcpTool"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		}

		if !Contains(traceExclusions, mcpMethod) {
			if rwSpan, ok := yokaimcpservercontext.CtxRootSpan(ctx).(otelsdktrace.ReadWriteSpan); ok {
				rwSpan.SetName(fmt.Sprintf("%s %s", rwSpan.Name(), spanNameSuffix))
				rwSpan.SetStatus(codes.Ok, "MCP request success")
				rwSpan.SetAttributes(spanAttributes...)
				rwSpan.End()
			}
		}

		if !Contains(logExclusions, mcpMethod) {
			log.CtxLogger(ctx).Info().Fields(logFields).Msg("MCP request success")
		}

		if metricsEnabled {
			p.requestsCounter.WithLabelValues(mcpMethod, metricTarget, "success").Inc()
			p.requestsDuration.WithLabelValues(mcpMethod, metricTarget).Observe(latency.Seconds())
		}
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		latency := time.Since(yokaimcpservercontext.CtxStartTime(ctx))

		mcpMethod := string(method)

		errMessage := fmt.Sprintf("%v", err)

		spanNameSuffix := mcpMethod

		spanAttributes := []attribute.KeyValue{
			attribute.String("mcp.latency", latency.String()),
			attribute.String("mcp.method", mcpMethod),
			attribute.String("mcp.error", errMessage),
		}

		logFields := map[string]interface{}{
			"mcpLatency": latency.String(),
			"mcpMethod":  mcpMethod,
			"mcpError":   errMessage,
		}

		metricTarget := ""

		jsonMessage, err := json.Marshal(message)
		if err == nil {
			if traceRequest {
				spanAttributes = append(spanAttributes, attribute.String("mcp.request", string(jsonMessage)))
			}

			if logRequest {
				logFields["mcpRequest"] = string(jsonMessage)
			}
		}

		switch method {
		case mcp.MethodResourcesRead:
			if req, ok := message.(*mcp.ReadResourceRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.URI)
				spanAttributes = append(spanAttributes, attribute.String("mcp.resource", req.Params.URI))
				logFields["mcpResourceURI"] = req.Params.URI
				metricTarget = req.Params.URI
			}
		case mcp.MethodPromptsGet:
			if req, ok := message.(*mcp.GetPromptRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.prompt", req.Params.Name))
				logFields["mcpPrompt"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		case mcp.MethodToolsCall:
			if req, ok := message.(*mcp.CallToolRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.tool", req.Params.Name))
				logFields["mcpTool"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		}

		if !Contains(traceExclusions, mcpMethod) {
			if rwSpan, ok := yokaimcpservercontext.CtxRootSpan(ctx).(otelsdktrace.ReadWriteSpan); ok {
				rwSpan.SetName(fmt.Sprintf("%s %s", rwSpan.Name(), spanNameSuffix))
				rwSpan.RecordError(err)
				rwSpan.SetStatus(codes.Error, errMessage)
				rwSpan.SetAttributes(spanAttributes...)
				rwSpan.End()
			}
		}

		if !Contains(logExclusions, mcpMethod) {
			log.CtxLogger(ctx).Error().Fields(logFields).Msg("MCP request error")
		}

		if metricsEnabled {
			p.requestsCounter.WithLabelValues(mcpMethod, metricTarget, "error").Inc()
			p.requestsDuration.WithLabelValues(mcpMethod, metricTarget).Observe(latency.Seconds())
		}
	})

	return hooks
}
