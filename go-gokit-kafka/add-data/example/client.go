package main

import (
	"context"

	pb "go-gokit-kafka/add-data/protobuf"

	"git.bluebird.id/bluebird/mq"
	util "git.bluebird.id/bluebird/util/microservice"
	proto "github.com/gogo/protobuf/proto"
)

func main() {
	isi := pb.Pegawai{Id: 1, NamaDepan: "arip", NamaBelakang: "saputra", Alamat: "cakung"}

	subscriber, _ := mq.NewAsyncProducer([]string{"127.0.0.1:9092", "127.0.0.1:9093"}, nil, nil)
	data, _ := proto.Marshal(&isi)
	subscriber.Publish(context.Background(), "test", data)

	util.OnShutdown(nil)
}
