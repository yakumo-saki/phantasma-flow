package repository

var repo *Repository

func GetRepository() *Repository {
	if repo == nil {
		repo = &Repository{}
		return repo
	}

	return repo
}
