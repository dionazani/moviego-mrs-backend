package infrastructuredto

type Response struct {
	Timestamp       string      `json:"timestamp"`
	ResponseStatus  int         `json:"-"` // Hidden from JSON
	ResponseCode    int         `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data"` // Using interface{} or [T any] for flexibility
}
