package server

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/golang/glog"
	"github.com/valyala/fasthttp"
)

func NewServer() {

	router := fasthttprouter.New()

	router.GET("/pushstatts", func(ctx *fasthttp.RequestCtx) {
		GetStatsOfMqtt()
	})

	router.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		glog.V(0).Infof("Panic occurred %s", p, ctx.Request.URI().String())
	}
	log.Println("Service Started on port " + "6001")
	glog.Fatal(fasthttp.ListenAndServe(":"+"6001", fasthttp.CompressHandler(router.Handler)))

}
