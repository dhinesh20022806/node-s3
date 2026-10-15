package service

import (
	"context"
	"errors"
	"os"

	genpb "github.com/dhinesh20022806/grpc_service/proto/gen"
	"github.com/dhinesh20022806/node-s3/app/pkg/dbx"
	"github.com/dhinesh20022806/node-s3/app/pkg/dispatcher"
	"github.com/dhinesh20022806/node-s3/app/pkg/worker"
)




type Object struct {

}

func NewObject() *Object{
	return &Object{}
}

func (o *Object) StoreObject(ctx context.Context, req *genpb.StoreObjectRequest) error {


		job := worker.WorkerJob{
		Data: req.Data,
	}

	

	dispatcher.DisPatcher.Submit(ctx, job)

	select {
		case result := <- job.Result:

			if result.Err != nil {
				return errors.New("")
			}

			// save in sql

			return nil
		
		case <- ctx.Done():

			return errors.New("new error")

	}

}

func (o *Object) GetObject(ctx context.Context, req *genpb.GetObjectRequest) ([]byte, error) {


	object, err := dbx.Select(req.ObjectId)


	if err != nil {



	}

	file, err := os.OpenFile(object.Filename, os.O_RDONLY, 0644)

	if err != nil {

	}

	buffer := make([]byte, object.ObjectSize)

	_, err = file.ReadAt(buffer, object.StartAddress);

	if err != nil {

	}

	return buffer, nil
}

