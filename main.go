// Package main is used to start pam stub in standalone mode.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/valkyrie-fnd/valkyrie-stubs/backdoors"
	"github.com/valkyrie-fnd/valkyrie-stubs/broken"
	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/genericpam"
	"github.com/valkyrie-fnd/valkyrie-stubs/memorydatastore"
)

var pam = "genericpam"
var operationMode = "default"
var addr = ":8080"
var backdoorPort = "3000"

var eds datastore.ExtendedDatastore

func init() {
	if v, found := os.LookupEnv("PAM"); found {
		pam = v
	}

	if v, found := os.LookupEnv("OP_MODE"); found {
		operationMode = v
	}

	if v, found := os.LookupEnv("ADDR"); found {
		addr = v
	}
	if v, found := os.LookupEnv("BACKDOOR_PORT"); found {
		backdoorPort = v
	}
	initData()
}

func main() {
	log.Info().Msgf(`---- Running %s stub server ----`, pam)
	backdoors.BackdoorServer(eds, fmt.Sprintf(":%s", backdoorPort))

	switch pam {
	case "genericpam":
		genericpam.RunServer(eds, genericpam.Config{
			PamAPIKey:      eds.GetPamAPIToken(),
			ProviderTokens: eds.GetProviderTokens(),
			Address:        addr,
			LogConfig: genericpam.LogConfig{
				Level: "info",
			},
		})
	case "broken":
		broken.RunServer(addr, genericpam.Routes(eds, eds.GetPamAPIToken(), eds.GetProviderTokens()))
	default:
		panic(errors.New("unknown operator"))
	}

	// block until signalled to exit
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}

func initData() {
	config, err := memorydatastore.ReadConfig("datastore.config.yaml")
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("---- Running datastore in %s mode", operationMode)
	switch operationMode {
	case "default":
		eds = memorydatastore.NewMapDataStore(config)
	case "performance":
		eds = memorydatastore.NewHPMapDataStore(config)
	default:
		panic(errors.New("invalid operation mode"))
	}
}
