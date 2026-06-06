package infrastructuredto

type Response struct {
	Timestamp       string      `json:"timestamp"`
	ResponseStatus  int         `json:"responseStatus"` // Hidden from JSON
	ResponseCode    int         `json:"-"` // Hidden from JSON
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data"` // Using interface{} or [T any] for flexibility
}
