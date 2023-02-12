package example

type ResponsePet struct {
	Data    *Pet   `json:"data"`
	Version string `json:"version"`
}
