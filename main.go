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
	// 안타깝게도 stop은 아직 쓸 일이 없음.. 혹시나 connection관련하여 나이스하게 종료해야 하면 필요
	for jsonBytes := range out {
		logger.Info(jsonBytes)
		pipeline.SendDataToAllPipeline(jsonBytes)
	}

	logger.Info("Listen on " + conf.ListenAddress + strconv.Itoa(conf.ListenPort))

	http.ListenAndServe(conf.ListenAddress + ":" + strconv.Itoa(conf.ListenPort), nil)


}
