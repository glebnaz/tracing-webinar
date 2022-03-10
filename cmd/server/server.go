package main

import (
	"encoding/json"
	"fmt"
	"github.com/glebnaz/tracing-webinar/internal/utils"
	"github.com/labstack/echo/v4"
	"go.opencensus.io/trace"
	"net/http"
)

func main() {

	utils.InitJaeger("server-web")
	e := echo.New()

	e.POST("/get", func(c echo.Context) error {
		spanJson := c.Request().Header.Get("X-Span-Context")
		var spanContext trace.SpanContext

		err := json.Unmarshal([]byte(spanJson), &spanContext)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		_ /* ctx */, span := trace.StartSpanWithRemoteParent(c.Request().Context(), "server web", spanContext)

		defer span.End()

		var req utils.Req

		err = c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		span.AddAttributes(
			trace.StringAttribute("number", fmt.Sprintf("%d", req.Number)))

		fmt.Println(req)

		return c.JSON(http.StatusOK, "good")
	})

	e.Start(":8080")
}
