package service

import (
	"context"
	"encoding/gob"

	"git.bluebird.id/bluebird/mq"
	notif "git.bluebird.id/bluebird/notification/server"
)

var (
	//EmptyResponse is ..
	EmptyResponse = Response{}
)

type pegawaiService struct {
	writer   ReadWriter
	notifier mq.Publisher
}

func init() {
	gob.Register(notif.Notification{})
}

func NewDataPegawaiService(writer ReadWriter, notifier mq.Publisher) DataPegawaiService {
	return &pegawaiService{writer: writer, notifier: notifier}
}

func (ds *pegawaiService) AddDataPegawai(ctx context.Context, req Pegawai) (Response, error) {
	err := ds.writer.WriteDataPegawai(req)
	if err != nil {
		return Response{Message: err.Error()}, err
	}
	return EmptyResponse, nil
}
