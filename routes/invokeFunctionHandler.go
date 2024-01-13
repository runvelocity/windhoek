package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/runvelocity/windhoek/utils"
)

const (
	MAXRETRIES  = 10
	BACKOFFTIME = 500
)

func InvokeFunctionHandler(c echo.Context) error {
	var vmRequest utils.FirecrackerVmRequest
	if err := c.Bind(&vmRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	vmRequest.SocketPath = fmt.Sprintf("%s/firecracker-%s.sock", utils.FC_SOCKETS_PATH, vmRequest.FunctionId)
	m, ctx, err := utils.CreateVm(vmRequest)
	defer func() {
		err := m.Shutdown(ctx)
		if err != nil {
			log.Println("An error occured while shutting down vm", err)
		}
		err = os.Remove(vmRequest.SocketPath)
		if err != nil {
			log.Println("Error deleting socket file", vmRequest.SocketPath)
		}
	}()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error occured while creating vm. %s", err.Error()))
	}
	fmt.Println(m.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP.String())
	url := fmt.Sprintf("http://%s:3000/invoke", m.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP.String())
	argsJSON, err := json.Marshal(vmRequest.InvokePayload.Args)
	if err != nil {
		fmt.Println("Error marshaling args to JSON:", err)
		errObj := utils.ErrorResponse{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errObj)
	}
	payload := fmt.Sprintf(`{"args": %s,"handler": "%s","codeLocation": "%s"}`, argsJSON, vmRequest.InvokePayload.Handler, vmRequest.CodeLocation)

	for i := 0; i < MAXRETRIES; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			fmt.Println("Error creating request:", err)
			errObj := utils.ErrorResponse{
				Message: err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, errObj)
		}

		// Set headers if needed
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
					fmt.Println("Error decoding JSON:", err)
					errObj := utils.ErrorResponse{
						Message: err.Error(),
					}
					return c.JSON(http.StatusInternalServerError, errObj)
				}
				if resObj["error"] != nil {
					errObj := utils.ErrorResponse{
						Message: resObj["error"].(string),
					}
					return c.JSON(res.StatusCode, errObj)
				}
				return c.JSON(res.StatusCode, resObj)
			}
			var resObj utils.FunctionInvokeResponse
			err = json.NewDecoder(res.Body).Decode(&resObj)
			if err != nil {
				fmt.Println("Error decoding JSON:", err)
				errObj := utils.ErrorResponse{
					Message: err.Error(),
				}
				return c.JSON(http.StatusInternalServerError, errObj)
			}
			return c.JSON(http.StatusOK, resObj)
		}
	}

	if err != nil {
		fmt.Println("Error invoking function:", err)
		errObj := utils.ErrorResponse{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errObj)
	}

	return nil

}
