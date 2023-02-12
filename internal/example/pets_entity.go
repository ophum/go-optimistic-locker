package example

type Pet struct {
	ID   uint   `json:"id" uri:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func (p *Pet) Copy() *Pet {
	return &Pet{
		ID:   p.ID,
		Name: p.Name,
		Age:  p.Age,
	}
}
