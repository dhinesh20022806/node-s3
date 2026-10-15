package types

import (
	"context"

	genpb "github.com/dhinesh20022806/grpc_service/proto/gen"
)

type ObjectService interface {
	StoreObject(ctx context.Context, req *genpb.StoreObjectRequest) error
	GetObject(ctx context.Context, req *genpb.GetObjectRequest) ([]byte, error)
}
