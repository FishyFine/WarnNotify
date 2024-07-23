package types

type StructMessage struct {
	Template string      `json:"template"`
	Meta     interface{} `json:"meta"`
	Tp       int         `json:"tp"`
}
