package main

import (
	"fmt"
	//"strings"

	cfg "git.bluebird.id/bluebird/util/config"
	cli "git.bluebird.id/bluebird/vault/endpoint"
)

func main() {
	//This example ahows how to use the secretservice
	//secured credential requestor library to resolve
	//secured link tag from configuration into a usable
	//values that can beusedin the program

	ok := cfg.AppConfig.LoadConfig()
	if !ok {
		return
	}

	discHost := cfg.Get("discoveryhost", "")

	//for example, the secured info tag  is stored
	//in db.password parameter:
	//
	sTag := cfg.Get("db.password", "")

	//for example, in service.conf this is ${mypasswordis12345}
	//extract the text only, less the tag ${}
	//
	//if strings.HasPrefix(sTag, "${") && strings.HasSuffix(sTag, "}") {
	//	sTag = strings.TrimSuffix(strings.TrimPrefix(sTag, "${"), "}")

	cd, err := cli.ResolveSecret(nil, sTag, discHost, nil)
	if err != nil {
		fmt.Printf("error %+v\n", err)
	}

	fmt.Printf("credential %s %+v\n", sTag, cd)
	//}
}
