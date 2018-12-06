package endpoint

import (
	"context"

	sv "go-gokit-kafka/add-data/service"

	"github.com/go-kit/kit/endpoint"
)

type PegawaiEndpoint struct {
	AddDataPegawaiEndpoint endpoint.Endpoint
}

func NewDataPegawaiEndpoint(service sv.DataPegawaiService) PegawaiEndpoint {
	addDataPegawai := makeAddDataPegawaiEndpoint(service)

	return PegawaiEndpoint{AddDataPegawaiEndpoint: addDataPegawai}
}

func makeAddDataPegawaiEndpoint(service sv.DataPegawaiService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sv.Pegawai)
		resp, err := service.AddDataPegawai(ctx, req)
		return resp, err
	}
}
