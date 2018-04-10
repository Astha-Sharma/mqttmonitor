package server

import (
	"encoding/json"

	"fmt"

	"strconv"

	"strings"

	"net/http"

	"github.com/nlopes/slack"
	"github.com/valyala/fasthttp"
	
	mqtt "github.com/vivekvasvani/mqttmonitoring/mqtt"
)

const (	
	SLACK_WEBHOOK                  = ""
	SLACK_WEBHOOK_TO_SEND_SLACKBOT = ""
)


func PushStats(ctx *fasthttp.RequestCtx) {
	

}

func SetErrorResponse(ctx *fasthttp.RequestCtx, statusCode, statusType, statusMessage string, httpStatus int) {
	log.Println(statusCode, statusType, statusMessage)
	var response Response
	response.Status.StatusCode = statusCode
	response.Status.StatusType = statusType
	response.Status.Message = statusMessage
	glog.Infoln("Error Reponse " + ToJsonString(response))
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(ToJsonString(response))
	ctx.SetStatusCode(httpStatus)
}

func SetSuccessResponse(ctx *fasthttp.RequestCtx, statusCode, statusType, statusMessage string, httpStatus int, data interface{}) {
	var response Response
	response.Status.StatusCode = statusCode
	response.Status.StatusType = statusType
	response.Status.Message = statusMessage
	response.Data = data
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(ToJsonString(response))
	glog.Infoln("Success Reponse " + ToJsonString(response))
	ctx.SetStatusCode(httpStatus)
}
}