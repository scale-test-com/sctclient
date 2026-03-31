package model

// ReqTargetPoint is a point on the request-rate curve.
type ReqTargetPoint struct {
	ElapsedTime string  `json:"elapsed_time" yaml:"elapsed_time"`
	ReqPerSec   float64 `json:"req_per_sec"  yaml:"req_per_sec"`
}

// OperationValidate holds validation rules for an HTTP operation.
type OperationValidate struct {
	StatusCode       int    `json:"status_code,omitempty"        yaml:"status_code,omitempty"`
	BodySuccessRegex string `json:"body_success_regex,omitempty" yaml:"body_success_regex,omitempty"`
	BodyFailRegex    string `json:"body_fail_regex,omitempty"    yaml:"body_fail_regex,omitempty"`
}

// Operation describes a single HTTP request to execute.
type Operation struct {
	URI      string             `json:"uri"               yaml:"uri"`
	Method   string             `json:"method"            yaml:"method"`
	Headers  map[string]string  `json:"headers,omitempty" yaml:"headers,omitempty"`
	Body     string             `json:"body,omitempty"    yaml:"body,omitempty"`
	Validate *OperationValidate `json:"validate,omitempty" yaml:"validate,omitempty"`
	Delay    string             `json:"delay,omitempty"   yaml:"delay,omitempty"`
	Repeat   int                `json:"repeat,omitempty"  yaml:"repeat,omitempty"`
}

// ScenarioContent is the full definition of a load test scenario.
type ScenarioContent struct {
	Name           string           `json:"name,omitempty"            yaml:"name,omitempty"`
	ReqTargetCurve []ReqTargetPoint `json:"req_target_curve"          yaml:"req_target_curve"`
	RequestTimeout string           `json:"request_timeout,omitempty" yaml:"request_timeout,omitempty"`
	OverTime       string           `json:"over_time,omitempty"       yaml:"over_time,omitempty"`
	Operations     []Operation      `json:"operations"                yaml:"operations"`
}

// CreateRunRequest is the request body for POST /run/new.
type CreateRunRequest struct {
	ScenarioID *int             `json:"scenario_id,omitempty"`
	Content    *ScenarioContent `json:"content,omitempty"`
}

// CreateRunResponse is the response body for POST /run/new.
type CreateRunResponse struct {
	ID         string  `json:"id"`
	State      string  `json:"state"`
	CreatedAt  string  `json:"created_at"`
	ScenarioID *int    `json:"scenario_id"`
	Message    string  `json:"message"`
}

// Run represents a full run with all details.
type Run struct {
	ID               string   `json:"id"`
	CreatedAt        string   `json:"created_at"`
	State            string   `json:"state"`
	ScenarioID       *int     `json:"scenario_id"`
	DurationS        *int     `json:"duration_s"`
	ExecutedRequests *int     `json:"executed_requests"`
	Error            *string  `json:"error"`
	Cost             *float64 `json:"cost"`
	StartAt          *string  `json:"start_at"`
}

// GetRunResponse wraps a Run in the `data` field.
type GetRunResponse struct {
	Data Run `json:"data"`
}

// ErrorResponse represents an API error body.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessMessage is returned by successful deletion responses.
type SuccessMessage struct {
	Message string `json:"message"`
}
