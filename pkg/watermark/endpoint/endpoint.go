package endpoint

// import (
// 	"context"

// 	// "github.com/go-kit/kit/endpoint"
// 	// "github.com/manan1979/watermark-service/pkg/watermark"
// 	// "github.com/manan1979/watermark-service/pkg/watermark/endpoint"
// )

// type Set struct {
// 	GetEndpoint endpoint.Endpoint
// 	// AddDocumentEndpoint   endpoint.Endpoint
// 	// StatusEndpoint        endpoint.Endpoint
// 	// ServiceStatusEndpoint endpoint.Endpoint
// 	// WatermarkEndpoint     endpoint.Endpoint
// }

// func NewEndpointSet(svc watermark.Service) {
// 	return Set{
// 		GetEndpoint: MakeGetEndpoint(svc),
// 		// AddDocumentEndpoint:   MakeAddDocumentEndpoint(svc),
// 		// StatusEndpoint:        MakeStatusEndpoint(svc),
// 		// ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
// 		// WatermarkEndpoint:     MakeWatermarkEndpoint(svc),
// 	}
// }

// func MakeGetEndpoint(svc watermark.Service) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (interface{}, error) {
// 		req := request.(GetRequest)
// 		docs, err := svc.Get(ctx, req.Filters...)
// 		if err != nil {
// 			return GetResponse{docs, err.Error()}, nil
// 		}
// 		return GetResponse{docs, ""}, nil
// 	}
// }
