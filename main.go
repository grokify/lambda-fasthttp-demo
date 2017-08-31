package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
)

const (
	PathPart1Name   = "pathPart1"
	PathPart2Name   = "pathPart2"
	QueryParam1Name = "queryParam1"
)

type data struct {
	Event     *apigatewayproxyevt.Event
	PathPart1 string
	PathPart2 string
	QueryVal1 string
	Body      string
	State     string
}

type ProxyOutput struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Body            string            `json:"body"`
	Headers         map[string]string `json:"headers"`
}

func HandleEawsyLambda(event *apigatewayproxyevt.Event, ctx *runtime.Context) (ProxyOutput, error) {
	data := data{
		Body:  event.Body,
		Event: event,
		State: "mystate"}
	if val1, ok := event.PathParameters[PathPart1Name]; ok {
		data.PathPart1 = val1
	}
	if val2, ok := event.PathParameters[PathPart2Name]; ok {
		data.PathPart2 = val2
	}
	if val3, ok := event.QueryStringParameters[QueryParam1Name]; ok {
		data.QueryVal1 = val3
	}
	bodyOut := HandleCanonical(data)
	ProxyOutput := ProxyOutput{
		StatusCode: 200,
		Body:       bodyOut}
	return ProxyOutput, nil
}

func HandleCanonical(data data) string {
	body, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	}
	return string(body)
}

func HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	data := data{
		Body:      string(ctx.PostBody()),
		QueryVal1: string(ctx.QueryArgs().Peek(QueryParam1Name)),
		State:     "mystate",
	}
	bodyOut := HandleCanonical(data)

	fmt.Fprintf(ctx, "%s", bodyOut)
}

func main() {
	router := fasthttprouter.New()
	router.GET("/test", HandleFastHTTP)
	router.GET("/test/:pathPart1/:pathPart2", HandleFastHTTP)
	router.POST("/test/:pathPart1/:pathPart2", HandleFastHTTP)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
