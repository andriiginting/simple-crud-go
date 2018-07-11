package domain

type Food struct {
	Id    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Price int    `json:"price,omitempty"`
	Owner string `json:"owner,omitempty"`
}
