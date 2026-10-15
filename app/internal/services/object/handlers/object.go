package handlers

import (
	"context"

	genpb "github.com/dhinesh20022806/grpc_service/proto/gen"
	"github.com/dhinesh20022806/node-s3/app/internal/services/object/types"
	"google.golang.org/grpc"
)


type ObjectHandler struct {
	objectService types.ObjectService
	genpb.UnimplementedObjectServiceServer
}

func NewObjectHandler ( grpc *grpc.Server ,objectService types.ObjectService)  {
	cred := &ObjectHandler{
		objectService: objectService,
	}

	genpb.RegisterObjectServiceServer(grpc, cred)
}

func (o *ObjectHandler) StoreObject(ctx context.Context, req *genpb.StoreObjectRequest) (*genpb.StoreObjectResponse, error) {

	err := o.objectService.StoreObject(ctx, req)

	if err != nil {
		return nil, err
	}

	return &genpb.StoreObjectResponse{
		ObjectId: req.ObjectId,
	}, nil
    
}

func (o *ObjectHandler) GetObject(ctx context.Context, req *genpb.GetObjectRequest) (*genpb.GetObjectResponse, error) {


  data, err := o.objectService.GetObject(ctx, req)


  if err != nil {
	return nil, err
  }

  return &genpb.GetObjectResponse{
	Data: data,
  }, nil

} 