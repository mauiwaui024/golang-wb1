package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	DataBase *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{DataBase: db}
}
