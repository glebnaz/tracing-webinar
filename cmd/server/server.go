package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/glebnaz/tracing-webinar/internal/utils"
	"github.com/labstack/echo/v4"
	"go.opencensus.io/trace"
)

func main() {
	utils.InitJaeger("server")

	e := echo.New()

	e.POST("/get", func(c echo.Context) error {
		spnctxJson := c.Request().Header.Get("X-Span-Context")
		var spanContext trace.SpanContext
		err := json.Unmarshal([]byte(spnctxJson), &spanContext)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		_, span := trace.StartSpanWithRemoteParent(
			context.Background(), "PostHandler", spanContext)
		defer span.End()

		fmt.Printf("span_id: %s", span.SpanContext().TraceID.String())

		req := utils.Req{}
		err = c.Bind(&req)

		span.AddAttributes(
			trace.StringAttribute("number", fmt.Sprintf("%d", req.Number)),
		)
		fmt.Println(req)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "good")
	})

	e.Start(":8080")
}
