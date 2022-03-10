package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/glebnaz/tracing-webinar/internal/utils"
	"go.opencensus.io/trace"
)

func main() {
	utils.InitJaeger("client")
	i := 0
	for {
		i++
		ctx, span := trace.StartSpan(context.Background(), "main")
		span.AddAttributes(
			trace.StringAttribute("method", "POST"),
		)
		log.Printf("TraceId: %s\n", span.SpanContext().TraceID.String())
		resp, err := sendPostRequest(ctx,
			"http://localhost:8080/get",
			int64(i))
		if err != nil {
			log.Printf("sendPostRequest: %v\n", err)
		}
		fmt.Printf("Response: %s\n", resp)
		span.End()

		time.Sleep(5 * time.Second)
	}
}

func sendPostRequest(ctx context.Context, url string, number int64) (string, error) {
	_ /* ctx2 */, span := trace.StartSpan(ctx, "sendPostRequest")
	defer span.End()

	spanContextJson, err := json.Marshal(span.SpanContext())
	if err != nil {
		return "", err
	}

	reqPayload := utils.Req{
		number,
	}

	jsonStr, _ := json.Marshal(reqPayload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Span-Context", string(spanContextJson))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
