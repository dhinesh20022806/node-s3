package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	genpb "github.com/dhinesh20022806/grpc_service/proto/gen"
	"github.com/dhinesh20022806/node-s3/app/pkg/dbx"
	"github.com/dhinesh20022806/node-s3/app/pkg/dispatcher"
	"github.com/dhinesh20022806/node-s3/app/pkg/worker"
)

type Object struct {
}

func NewObject() *Object {
	return &Object{}
}

func (o *Object) StoreObject(ctx context.Context, req *genpb.StoreObjectRequest) error {

	job := worker.WorkerJob{
		Data:   req.Data,
		Result: make(chan worker.WriteResult, 1),
	}

	dispatcher.DisPatcher.Submit(ctx, job)

	select {

	case result, ok := <-job.Result:

		if !ok {
			return errors.New("result channel closed without a value")
		}

		slog.Info("result arrived", "result filenName", result.Filename, "err", result.Err)
		if result.Err != nil {
			return errors.New("new error")
		}

		// save in sql

		err := dbx.InsertObject(result.Filename, result.Bytesize, result.StartAddr, req.ObjectId, 0)

		if err != nil {
			return err
		}

		fmt.Printf("\nfile inserted\n")

		return nil

	case <-time.After(3 * time.Second):

		fmt.Println("Debugging Error: Channel timed out. No communication occurred.")
		return errors.New("Timeout")

	case <-ctx.Done():

		return errors.New("new error")

	}

}

func (o *Object) GetObject(ctx context.Context, req *genpb.GetObjectRequest) ([]byte, error) {

	object, err := dbx.Select(req.ObjectId)

	if err != nil {

		fmt.Printf("error 77 = %v", err)

	}

	file, err := os.OpenFile(object.Filename, os.O_RDONLY, 0644)

	if err != nil {

		fmt.Printf("error 85 = %v", err)

	}

	buffer := make([]byte, object.ObjectSize)

	_, err = file.ReadAt(buffer, object.StartAddress)

	if err != nil {

		fmt.Printf("error 95 = %v", err)

	}

	return buffer, nil
}
