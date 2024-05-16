package types

type MsgServer struct {
	Code     int    `json:"code,omitempty"`
	Zip      int    `json:"zip,omitempty"`
	Category string `json:"category,omitempty"`
	Str      string `json:"str,omitempty"`
}
