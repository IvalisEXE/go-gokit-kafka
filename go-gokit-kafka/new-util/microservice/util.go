package microservice

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/zk"
)

const (
	servicePath  = "/service/"
	registryNode = "registry"
)

//zkClient returns zk client
func zkClient(nodes []string, logger log.Logger) (zk.Client, error) {
	options := zk.ConnectTimeout(time.Second * 5)
	return zk.NewClient(nodes, logger, options)
}

//ServiceRegistry returns zk service registrar
func ServiceRegistry(nodes []string, serviceName, address string, logger log.Logger) (*zk.Registrar, error) {
	client, err := zkClient(nodes, logger)
	if err != nil {
		return nil, err
	}
	path := servicePath + serviceName
	service := zk.Service{Path: path, Name: registryNode, Data: []byte(address)}
	return zk.NewRegistrar(client, service, logger), nil
}

//ServiceDiscovery returns zk service instancer
func ServiceDiscovery(nodes []string, serviceName string, logger log.Logger) (*zk.Instancer, error) {
	client, err := zkClient(nodes, logger)
	if err != nil {
		return nil, err
	}
	path := servicePath + serviceName
	instancer, err := zk.NewInstancer(client, path, logger)
	if err != nil {
		return nil, err
	}
	return instancer, nil
}

//OnShutdown calls shutdown on signal interrupt
func OnShutdown(shutdown func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-c
	if shutdown != nil {
		shutdown()
	}
}

//RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type RecoveryHandlerFunc func(p interface{}) (err error)

//RecoverFrom call recovery handler function
func RecoverFrom(p interface{}, r RecoveryHandlerFunc) error {
	if r == nil {
		return fmt.Errorf("Server error: %s", p)
	}
	return r(p)
}

//GoWithRecover call go routine with recovery
func GoWithRecover(function func(), recoverFunc RecoveryHandlerFunc) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := RecoverFrom(r, recoverFunc)
				if err != nil {
					fmt.Print(err)
				}
			}
		}()

		function()

	}()
}
