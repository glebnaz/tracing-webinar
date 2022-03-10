package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/glebnaz/tracing-webinar/internal/utils"
	"go.opencensus.io/trace"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	utils.InitJaeger("client-webinar")
	i := 0
	for {
		i++
		ctx, span := trace.StartSpan(context.Background(), fmt.Sprintf("%d", i))
		//span.AddAttributes()

		answer, err := sendPostRequest(ctx, "http://localhost:8080/get", i)
		if err != nil {
			fmt.Printf("err: %s", err)
			span.End()
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("answer from server: %s", answer)

		span.End()
		time.Sleep(3 * time.Second)
	}
}

func sendPostRequest(ctx context.Context, url string, number int) (string, error) {
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
