package interceptors

import (
	"context"
	"github.com/majidbl/discount/config"
	"github.com/majidbl/discount/pkg/logger"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	totalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "egiftcharges_service_requests_total",
		Help: "The total number of incoming gRPC requests",
	})
)

type GrpcMetricsCb func(err error)

type InterceptorManager interface {
	Logger(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error)
	ClientRequestLoggerInterceptor() func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error
}

// InterceptorManager struct
type interceptorManager struct {
	logger    logger.Logger
	cfg       *config.Config
	metricsCb GrpcMetricsCb
}

// NewInterceptorManager InterceptorManager constructor
func NewInterceptorManager(logger logger.Logger, cfg *config.Config) *interceptorManager {
	return &interceptorManager{logger: logger, cfg: cfg}
}

// Logger Interceptor
func (im *interceptorManager) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.GrpcMiddlewareAccessLogger(info.FullMethod, time.Since(start), md, err)
	if im.metricsCb != nil {
		im.metricsCb(err)
	}
	return reply, err
}

// ClientRequestLoggerInterceptor gRPC client interceptor
func (im *interceptorManager) ClientRequestLoggerInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		md, _ := metadata.FromIncomingContext(ctx)
		im.logger.GrpcClientInterceptorLogger(method, req, reply, time.Since(start), md, err)
		return err
	}
}
