package repository

type Repository struct {
	// Dependencies will be injected here
}

func NewRepository() *Repository {
	return &Repository{}
}
