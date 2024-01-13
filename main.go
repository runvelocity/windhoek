package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/runvelocity/windhoek/routes"
	"github.com/runvelocity/windhoek/utils"
)

type PingResponse struct {
	Ok bool `json:"ok"`
}

func main() {
	e := echo.New()

	if err := utils.WriteCNIConfWithHostLocalSubnet(fmt.Sprintf("/etc/cni/conf.d/%s.conflist", utils.FC_NETWORK_NAME)); err != nil {
		e.Logger.Fatal(err.Error())
	}
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, PingResponse{Ok: true})
	})

	e.POST("/invoke", routes.InvokeFunctionHandler)
	e.Logger.Fatal(e.Start(":8000"))
}
