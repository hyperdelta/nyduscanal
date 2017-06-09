package main

import (
	"github.com/hyperdelta/nyduscanal/gateway"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperdelta/refinery/handler"
	"strconv"
	"github.com/hyperdelta/refinery/config"
	"github.com/hyperdelta/refinery/log"
	"github.com/hyperdelta/refinery/pipeline"
)

var (
	conf = config.Config{ListenAddress: "", ListenPort:3000}
	defaultRouter *mux.Router
	logger *log.Logger = log.Get()
)

func main() {

	defaultRouter = mux.NewRouter()

	m := http.NewServeMux()
	handler.CreateDefaultRegisteredHandlers(defaultRouter)

	m.Handle("/", defaultRouter)
	http.DefaultServeMux = m

	stopChannel := make(chan bool)
	defer func() {
		close(stopChannel)
	}()
	out := gateway.StartGmarketNydusCanal(stopChannel, "172.30.219.47:9092")

	logger.Info("Init nyduscanal")

	go func() {
		for jsonBytes := range out {
			if(len(jsonBytes) > 5) {
				logger.Info(string(jsonBytes))
				pipeline.SendDataToAllPipeline(jsonBytes)
			}
		}
	}()



	logger.Info("Listen on " + conf.ListenAddress + strconv.Itoa(conf.ListenPort))

	http.ListenAndServe(conf.ListenAddress + ":" + strconv.Itoa(conf.ListenPort), nil)


}
