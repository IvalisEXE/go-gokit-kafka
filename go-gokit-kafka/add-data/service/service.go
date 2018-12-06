package service

import (
	"context"
)

//ServiceID is dispatch service ID
const (
	ServiceID = "add.data"
	Topic     = "test"
)

//DispatchEvent is ..
type Pegawai struct {
	ID           int64
	NamaDepan    string
	NamaBelakang string
	Alamat       string
}

//Response is ..
type Response struct {
	Message string
}

//DispatchEventService is ..
type DataPegawaiService interface {
	AddDataPegawai(context.Context, Pegawai) (Response, error)
}

//ReadWriter is ..
type ReadWriter interface {
	WriteDataPegawai(pe Pegawai) error
}
