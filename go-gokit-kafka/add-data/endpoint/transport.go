package endpoint

import (
	"context"

	pb "go-gokit-kafka/add-data/proto"
	pbd "go-gokit-kafka/add-data/protobuf"

	sv "go-gokit-kafka/add-data/service"

	"git.bluebird.id/bluebird/mq"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/gogo/protobuf/proto"
	stdopentracing "github.com/opentracing/opentracing-go"
	oldcontext "golang.org/x/net/context"
)

type grpcDataPegawaiServer struct {
	addDataPegawai grpctransport.Handler
}

func NewGRPCDispatchEventServer(endpoints PegawaiEndpoint, tracer stdopentracing.Tracer, logger log.Logger) pb.DataPegawaiServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcDataPegawaiServer{
		addDataPegawai: grpctransport.NewServer(endpoints.AddDataPegawaiEndpoint,
			decodeAddDataPegawaiRequest,
			encodeAddDataPegawaiResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "AddDataPegawai", logger)))...),
	}
}

func decodeAddDataPegawaiRequest(_ context.Context, request interface{}) (interface{}, error) {
	// req := request.(*pb.Pegawai)
	// return sv.Pegawai{
	// 	ID:           req.Id,
	// 	NamaDepan:    req.NamaDepan,
	// 	NamaBelakang: req.NamaBelakang,
	// 	Alamat:       req.Alamat,
	// }, nil
	return nil, nil
}

func encodeAddDataPegawaiResponse(_ context.Context, response interface{}) (interface{}, error) {
	// resp := response.(sv.Response)
	// return &pb.Response{Message: resp.Message}, nil
	return nil, nil
}

func (s *grpcDataPegawaiServer) AddDataPegawai(ctx oldcontext.Context, cred *pb.Pegawai) (*pb.Response, error) {
	_, resp, err := s.addDataPegawai.ServeGRPC(ctx, cred)
	if err != nil {
		return &pb.Response{Message: err.Error()}, err
	}
	return resp.(*pb.Response), err
}

//NewOrderEventSubscriber returns new order winner subscriber to receive winner data from kafka
func NewSubscriber(brokers []string, topic string, endpoint PegawaiEndpoint) mq.Subscriber {
	return mq.NewConsumerGroup(brokers, sv.Topic, nil, topic, eventHandler(endpoint))
}

func eventHandler(endpoint PegawaiEndpoint) mq.MessageHandler {
	decoder := func(ctx context.Context, msg interface{}) (interface{}, error) {
		var event pbd.Pegawai
		err := proto.Unmarshal(msg.([]byte), &event)

		if err != nil {
			return nil, err
		}

		return sv.Pegawai{
			ID:           event.Id,
			NamaDepan:    event.NamaBelakang,
			NamaBelakang: event.NamaBelakang,
			Alamat:       event.Alamat,
		}, nil
	}

	return mq.MessageHandler{Decode: decoder, Endpoint: endpoint.AddDataPegawaiEndpoint}
}
