package domain

type Number struct {
	ID    int64 `json:"id"`
	Value int   `json:"value"`
}

type NumberRepository interface {
	Save(value int) error
	GetAllSorted() ([]int, error)
}

type NumberService interface {
	AddNumber(value int) ([]int, error)
}
