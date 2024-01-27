package models

type FirecrackerVmRequest struct {
	FunctionId       string        `json:"functionId"`
	SocketPath       string        `json:"socketPath"`
	CodeLocation     string        `json:"codeLocation"`
	InvokePayload    InvokePayload `json:"invokePayload"`
	Cpu              int           `json:"cpu"`
	Memory           int           `json:"memory"`
	EphemeralStorage int           `json:"ephemeralStorage"`
	Runtime          string        `json:"runtime"`
}

type InvokePayload struct {
	Handler string                 `json:"handler"`
	Args    map[string]interface{} `json:"args"`
}

type WorkerPingResponse struct {
	Ok bool `json:"ok"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type FunctionInvokeResponse struct {
	InvocationResponse any `json:"invocationResponse"`
}
