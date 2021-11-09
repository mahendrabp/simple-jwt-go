package repository

import (
	"database/sql"
	"net/http"
	"simple-jwt-go/pkg/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type postgreeRepository struct {
	db *sqlx.DB
}

func NewPostgreRepository(db *sqlx.DB) domain.PostgreUserRepository {
	return &postgreeRepository{db}
}

func (r *postgreeRepository) GetByID(id uuid.UUID) (domain.User, error) {
	var user domain.User

	if err := r.db.Get(&user, getUserQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return user, echo.ErrNotFound
		}

		return user, echo.ErrBadRequest
	}

	return user, nil
}

func (r *postgreeRepository) FindByEmail(email string) (domain.User, error) {
	var user domain.User

	if err := r.db.QueryRowx(findUserByEmailQuery, email).StructScan(&user); err != nil {
		if err == sql.ErrNoRows {
			return user, echo.NewHTTPError(
				http.StatusBadRequest,
				"user not found",
			)
		}
		return user, echo.ErrBadRequest
	}

	return user, nil
}

func (r *postgreeRepository) Store(u *domain.User) (*domain.User, error) {
	var user domain.User

	if err := r.db.QueryRowx(createUserQuery, u.Email, u.Password).StructScan(&user); err != nil {
		if err.(*pq.Error).Code == "23505" {
			return nil, echo.NewHTTPError(
				http.StatusBadRequest,
				"email already exists",
			)
		}
		return nil, echo.ErrBadRequest
	}

	return &user, nil
}

func (r *postgreeRepository) Update(a *domain.User) (*domain.User, error) {
	var user domain.User

	if err := r.db.QueryRowx(updateUserQuery, a.Email, a.Password, a.ID).StructScan(&user); err != nil {
		if err == sql.ErrNoRows {
			return nil, echo.ErrNotFound
		}

		return nil, echo.ErrBadRequest
	}

	return &user, nil
}

func (r *postgreeRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(deleteUserQuery, id)
	if err != nil {
		return echo.ErrBadRequest
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return echo.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return echo.ErrNotFound
	}

	return nil
}
