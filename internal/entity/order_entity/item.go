package order_entity

type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func NewItem(id string, name string, quantity int) Item {
	return Item{
		Id:       id,
		Name:     name,
		Quantity: quantity,
	}
}
