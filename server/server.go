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
		env := ctx.QueryArgs().Peek("env")
		GetStatsOfMqtt(string(env))
	})

	router.GET("/getconnect", func(ctx *fasthttp.RequestCtx) {
		GetConnectLatency(ctx)
	})

	router.GET("/getupdowntime", func(ctx *fasthttp.RequestCtx) {
		GetUpDownTime(ctx)
	})

	router.GET("/timedistribution", func(ctx *fasthttp.RequestCtx) {
		TimeDistribution(ctx)
	})

	router.GET("/crashfreetrendsandroid", func(ctx *fasthttp.RequestCtx) {
		CrashFreeTrendsAndroid(ctx)
	})

	router.GET("/crashfreetrendsandroidtop", func(ctx *fasthttp.RequestCtx) {
		CrashFreeTrendsAndroidTopbuilds(ctx)
	})

	router.GET("/crashfreetrendsios", func(ctx *fasthttp.RequestCtx) {
		CrashFreeTrendsIOS(ctx)
	})

	router.GET("/allandroidbuilds", func(ctx *fasthttp.RequestCtx) {
		CrashFreeTrendsIOS(ctx)
	})

	router.GET("/androidbuildsversion", func(ctx *fasthttp.RequestCtx) {
		AndroidBuildVersion(ctx)
	})

	router.GET("/androidcrashbyversion", func(ctx *fasthttp.RequestCtx) {
		AndroidCrashByVersion(ctx)
	})

	router.GET("/microappscrashes", func(ctx *fasthttp.RequestCtx) {
		MicroAppsCrashes(ctx)
	})

	router.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		glog.V(0).Infof("Panic occurred %s", p, ctx.Request.URI().String())
	}
	log.Println("Service Started on port " + "6001")
	glog.Fatal(fasthttp.ListenAndServe(":"+"6001", fasthttp.CompressHandler(router.Handler)))

}
