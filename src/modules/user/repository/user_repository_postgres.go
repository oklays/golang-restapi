package repository

import (
	"database/sql"

	"github.com/oklays/golang-restapi/src/modules/user/model"
)

type userRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRespositoryPostgres(db *sql.DB) *userRepositoryPostgres {
	return &userRepositoryPostgres{db}
}

func (r *userRepositoryPostgres) Save(user *model.User) error {
	query := `INSERT INTO "users"("id", "email", "mobile_phone", "password", "full_name", "name", "dob", "photo", 				"created_at", "updated_at", "id_device", "pin")
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.ID, user.Email, user.MobilePhone, user.Password, user.FullName, user.Name, user.Dob, user.Photo, user.CreatedAt, user.UpdatedAt, user.IdDevice, user.Pin)

	if err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryPostgres) Update(id int64, user *model.User) error {
	query := `UPDATE "users" SET "email"= $1, "mobile_phone"=$2, "password"=$3, "full_name"=$4, "name"=$5, "dob"=$6, "photo"=$7, 				"created_at"=$8, "updated_at"=$9, "id_device"=$10, "pin"=$11`

	statement, err := r.db.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.Email, user.MobilePhone, user.Password, user.FullName, user.Name, user.Dob, user.Photo, user.CreatedAt, user.UpdatedAt, user.IdDevice, user.Pin)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryPostgres) Delete(id int64) error {
	query := `DELETE from "users" WHERE "id" = $1`

	statement, err := r.db.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryPostgres) FindByID(id int64) (*model.User, error) {
	query := `SELECT * FROM "users" WHERE "id"=$1`

	var user model.User

	statement, err := r.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	err = statement.QueryRow(id).Scan(&user.ID, &user.Email, &user.MobilePhone, &user.Password, &user.FullName, &user.Name, &user.Dob, &user.Photo, &user.CreatedAt, &user.UpdatedAt, &user.IdDevice, &user.Pin)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryPostgres) FindAll() (model.Users, error) {
	query := `SELECT * FROM "users"`

	var users model.Users

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User

		err = rows.Scan(&user.ID, &user.Email, &user.MobilePhone, &user.Password, &user.FullName, &user.Name, &user.Dob, &user.Photo, &user.CreatedAt, &user.UpdatedAt, &user.IdDevice, &user.Pin)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
