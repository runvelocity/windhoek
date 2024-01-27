package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/runvelocity/windhoek/handlers"
	"github.com/runvelocity/windhoek/internal/network"
)

type PingResponse struct {
	Ok bool `json:"ok"`
}

func main() {
	e := echo.New()

	if err := network.WriteCNIConfWithHostLocalSubnet(fmt.Sprintf("/etc/cni/conf.d/%s.conflist", network.FC_NETWORK_NAME)); err != nil {
		e.Logger.Fatal(err.Error())
	}
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, PingResponse{Ok: true})
	})

	e.POST("/invoke", handlers.InvokeHandler)
	e.Logger.Fatal(e.Start(":8000"))
}
