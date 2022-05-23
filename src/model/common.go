package model

type QueryResult struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	isEmpty bool
}

// NewOKResult returns successful query result
func NewOKResult(result interface{}) *QueryResult {
	return &QueryResult{
		Status: "OK",
		Result: result,
	}
}

// NewResultWithStatus creates result with status
func NewResultWithStatus(status string, result interface{}) *QueryResult {
	return &QueryResult{
		Status: status,
		Result: result,
	}
}

var emptyOKResult *QueryResult = &QueryResult{
	Status:  "OK",
	Result:  map[string]interface{}{},
	isEmpty: true,
}

// EmptyOKResult returns successfull empty query result
func EmptyOKResult() *QueryResult {
	return emptyOKResult
}
