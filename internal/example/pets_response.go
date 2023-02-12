package example

type ResponsePet struct {
	Data *Pet   `json:"data"`
	Etag string `json:"_etag"`
}
