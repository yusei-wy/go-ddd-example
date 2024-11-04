package infra

import (
	"go_ddd_example/feature/user/domain"
	"go_ddd_example/feature/user/domain/model"
	customerror "go_ddd_example/share/custom_error"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var _ domain.UserRepository = (*PsqlUserRepository)(nil)

type PsqlUserRepository struct {
	db *sqlx.DB
}

func NewPsQlUserRepository(db *sqlx.DB) *PsqlUserRepository {
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

func (r *PsqlUserRepository) GetUsers(userIds []uuid.UUID) ([]model.User, customerror.RepositoryError) {
	query := `
		SELECT
			id
			, name
			, created_at
			, updated_at
		FROM
			users
		WHERE
			id IN (:ids)
	`

	var queryables []QueryableUser
	if err := r.db.Select(&queryables, query, userIds); err != nil {
		return nil, customerror.NewRepositoryError(err)
	}

	users := make([]model.User, 0, len(queryables))
	for _, queryable := range queryables {
		user := model.NewUser(queryable.ID, queryable.Name)
		users = append(users, user)
	}

	return users, nil
}

func (r *PsqlUserRepository) GetUser(userID uuid.UUID) (*model.User, customerror.RepositoryError) {
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
	if err := r.db.Get(&queryable, query, userID); err != nil {
		return nil, customerror.NewRepositoryError(err)
	}

	user := model.NewUser(queryable.ID, queryable.Name)

	return &user, nil
}
