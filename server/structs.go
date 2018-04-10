package server

type StatusMessage struct {
	StatusCode string `json:"statusCode"`
	StatusType string `json:"statusType"`
	Message    string `json:"message"`
}

//Response structure
type Response struct {
	Status StatusMessage `json:"status"`
	Data   interface{}   `json:"data,omitempty"`
}
