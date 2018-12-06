package main

import (
	"fmt"

	ep "go-gokit-kafka/add-data/endpoint"
	pb "go-gokit-kafka/add-data/proto"
	svc "go-gokit-kafka/add-data/service"
	cfg "go-gokit-kafka/new-util/config"
	run "go-gokit-kafka/new-util/grpc"
	lg "go-gokit-kafka/new-util/log"
	util "go-gokit-kafka/new-util/microservice"
	tr "go-gokit-kafka/new-util/opentracing"

	"git.bluebird.id/bluebird/mq"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	logger := lg.Logger()

	ok := cfg.AppConfig.LoadConfig()
	if !ok {
		logger.Log(lg.LogError, "failed to load configuration")
		return
	}

	discHost := cfg.GetA("discoveryhost", "127.0.0.1:2181")
	ip := cfg.Get("serviceip", "127.0.0.1")
	port := cfg.Get("serviceport", "7001")
	address := fmt.Sprintf("%s:%s", ip, port)

	registrar, err := util.ServiceRegistry(discHost, svc.ServiceID, address, logger)
	if err != nil {
		logger.Log(lg.LogError, "cannot find registrar")
		return
	}
	registrar.Register()
	defer registrar.Deregister()

	tracerHost := cfg.Get("tracerhost", "127.0.0.1:9999")
	tracer := tr.Tracer(tracerHost, svc.ServiceID, nil)

	var server pb.DataPegawaiServiceServer
	var subscriber mq.Subscriber
	{
		//chHost := cfg.Get("chhost", "127.0.0.1:6379")
		//cacher := svc.NewRedisCache(chHost)

		//gmapKey := cfg.Get("gmapkey", "AIzaSyD9tm3UVfxRWeaOy_MQ7tsCj1fVCLfG8Bo")
		//locator := svc.NewLocator(gmapKey)

		dbHost := cfg.Get(cfg.DBhost, "127.0.0.1:3306")
		dbName := cfg.Get(cfg.DBname, "test")
		dbUser := cfg.Get(cfg.DBuid, "root")
		dbPwd := cfg.Get(cfg.DBpwd, "root")

		brokers := cfg.GetA("mqbrokers", "127.0.0.1:9092")
		topic := cfg.Get("topic", "test")

		dbReadWriter := svc.NewDBReadWriter(dbHost, dbName, dbUser, dbPwd)
		//dbRuler := svc.NewDBDispatchRuler(dbReadWriter, locator)
		notifier, nerr := mq.NewAsyncProducer(brokers, nil, nil)
		if nerr != nil {
			logger.Log(lg.LogError, "failed while create notifier")
			return
		}
		defer notifier.Close()

		service := svc.NewDataPegawaiService(dbReadWriter, notifier)
		endpoint := ep.NewDataPegawaiEndpoint(service)
		server = ep.NewGRPCDispatchEventServer(endpoint, tracer, logger)

		subscriber = ep.NewSubscriber(brokers, topic, endpoint)
	}

	grpcServer := grpc.NewServer(run.Recovery(logger)...)
	pb.RegisterDataPegawaiServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	go run.Serve(address, grpcServer, logger)

	subscriber.Subscribe()
	defer subscriber.Unsubscribe()

	util.OnShutdown(func() {
		grpcServer.GracefulStop()
	})

}
