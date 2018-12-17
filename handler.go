package otaws

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

func AddOTHandlers(cl *client.Client) {
	cl.Handlers.Build.PushFront(otHandlers)
}

func otHandlers(r *request.Request) {
	sp := opentracing.StartSpan(r.Operation.Name)
	ext.SpanKindRPCClient.Set(sp)
	ext.Component.Set(sp, "go-aws")
	ext.HTTPMethod.Set(sp, r.Operation.HTTPMethod)
	ext.HTTPUrl.Set(sp, r.HTTPRequest.URL.String())
	ext.PeerService.Set(sp, r.ClientInfo.ServiceName)

	_ = inject(sp, r.HTTPRequest.Header)

	r.Handlers.Complete.PushBack(func(req *request.Request) {
		ext.HTTPStatusCode.Set(sp, uint16(req.HTTPResponse.StatusCode))
		sp.Finish()
	})

	r.Handlers.Retry.PushBack(func(req *request.Request) {
		sp.LogFields(log.String("event", "retry"))
	})
}

func inject(span opentracing.Span, header http.Header) error {
	return opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header))
}
