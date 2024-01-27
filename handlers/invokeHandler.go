package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/runvelocity/windhoek/internal/vm"
	"github.com/runvelocity/windhoek/models"
	"github.com/sirupsen/logrus"
)

const (
	MAXRETRIES  = 10
	BACKOFFTIME = 500
)

var log = logrus.New()

func InvokeHandler(c echo.Context) error {
	var vmRequest models.FirecrackerVmRequest
	if err := c.Bind(&vmRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	vmRequest.SocketPath = fmt.Sprintf("%s/firecracker-%s.sock", vm.FC_SOCKETS_PATH, vmRequest.FunctionId)

	vmManager := vm.VmManager{}
	m, ctx, err := vmManager.CreateVm(vmRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error occured while creating vm. %s", err.Error()))
	}
	defer func() {
		err := m.Shutdown(ctx)
		if err != nil {
			log.Error("An error occured while shutting down vm", err)
		}
		err = os.Remove(vmRequest.SocketPath)
		if err != nil {
			log.Error("Error deleting socket file", vmRequest.SocketPath)
		}
	}()

	url := fmt.Sprintf("http://%s:3000/invoke", m.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP.String())
	argsJSON, err := json.Marshal(vmRequest.InvokePayload.Args)
	if err != nil {
		log.Error("Error marshaling args to JSON:", err)
		errObj := models.ErrorResponse{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errObj)
	}
	payload := fmt.Sprintf(`{"args": %s,"handler": "%s","codeLocation": "%s"}`, argsJSON, vmRequest.InvokePayload.Handler, vmRequest.CodeLocation)

	for i := 0; i < MAXRETRIES; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Error("Error creating request:", err.Error())
			errObj := models.ErrorResponse{
				Message: err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, errObj)
		}

		req.Header.Set("Content-Type", "application/json")

		req.Close = true

		var client = &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			time.Sleep(BACKOFFTIME * time.Millisecond)
		} else {
			if res.StatusCode != 200 {
				var resObj map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&resObj)
				if err != nil {
					log.Error("Error decoding JSON:", err)
					errObj := models.ErrorResponse{
						Message: err.Error(),
					}
					return c.JSON(http.StatusInternalServerError, errObj)
				}
				if resObj["error"] != nil {
					errObj := models.ErrorResponse{
						Message: resObj["error"].(string),
					}
					return c.JSON(res.StatusCode, errObj)
				}
				return c.JSON(res.StatusCode, resObj)
			}

			var resObj models.FunctionInvokeResponse
			err = json.NewDecoder(res.Body).Decode(&resObj)
			if err != nil {
				log.Error("Error decoding JSON:", err)
				errObj := models.ErrorResponse{
					Message: err.Error(),
				}
				return c.JSON(http.StatusInternalServerError, errObj)
			}
			return c.JSON(http.StatusOK, resObj)
		}
	}

	if err != nil {
		log.Error("Error invoking function:", err)
		errObj := models.ErrorResponse{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errObj)
	}

	return nil

}
