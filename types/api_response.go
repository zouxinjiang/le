package types

type (
	ApiResponseStructure struct {
		Code    string      `json:"Code"`
		Message string      `json:"Message"`
		Count   *int        `json:"Count,omitempty"`
		Data    interface{} `json:"Data,omitempty"`
	}
)
