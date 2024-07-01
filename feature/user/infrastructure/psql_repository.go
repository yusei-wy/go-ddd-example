package infrastructure

import (
	"go_ddd_example/feature/user/domain"
	"go_ddd_example/feature/user/domain/model"
	customerror "go_ddd_example/share/custom_error"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PsqlUserRepository struct {
	db *sqlx.DB
}

func NewPsQlUserRepository(db *sqlx.DB) domain.UserRepository {
	return &PsqlUserRepository{db}
}

func (r *PsqlUserRepository) CreateUser(cmd model.UserCommand) customerror.RepositoryError {
	query := NewUserQuery(cmd)

	sql := `
		INSERT INTO users (
			id
			, name
			, created_at
			, updated_at
		)
		VALUES (
			:id
			, :name
			, :created_at
			, :updated_at
		)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name
			, updated_at = EXCLUDED.updated_at
	`

	if _, err := r.db.NamedExec(sql, query); err != nil {
		return customerror.NewRepositoryError(err)
	}

	return nil
}

func (r *PsqlUserRepository) GetUser(userId uuid.UUID) (*model.User, customerror.RepositoryError) {
	query := `
		SELECT
			id
			, name
			, created_at
			, updated_at
		FROM
			users
		WHERE
			id = $1
	`

	var queryable QueryableUser
	if err := r.db.Get(&queryable, query, userId); err != nil {
		return nil, customerror.NewRepositoryError(err)
	}

	user := model.NewUser(queryable.Id, queryable.Name)

	return &user, nil
}
