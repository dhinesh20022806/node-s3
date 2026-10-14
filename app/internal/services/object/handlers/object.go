package handlers

import (
	"context"
	"os"

	genpb "github.com/dhinesh20022806/grpc_service/proto/gen"
	"github.com/dhinesh20022806/node-s3/internal/services/object/types"
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

	bytesofdata := req.Data

	file, err := os.OpenFile("./data/a", os.O_APPEND|os.O_CREATE, 0664)

	if err != nil {
		panic(err)
	}

	_, err = file.Write(bytesofdata)

	if err != nil {
		return &genpb.StoreObjectResponse{
			ObjectId: req.ObjectId,
		}, err

	}

	return &genpb.StoreObjectResponse{
		ObjectId: req.ObjectId,
	}, nil
    
}

func (o *ObjectHandler) GetObject(ctx context.Context, req *genpb.GetObjectRequest) (*genpb.GetObjectResponse, error) {




} 