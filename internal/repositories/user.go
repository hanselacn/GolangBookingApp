package repositories

import (
	"context"
	"fmt"

	"GolangBookingApp/internal/entity"
	"GolangBookingApp/pkg/postgres"

	"github.com/google/uuid"
)

type User interface {
	Create(context.Context, *entity.User) (*uuid.UUID, error)
	Verify(context.Context, string) (entity.User, error)
	GrantAdmin(context.Context, string) error
	DemoteToUser(context.Context, string) error
	GetAllUser(context.Context) ([]entity.User, error)
	VerifyLogin(context.Context, string) (entity.User, error)
	VerifyLogout(context.Context, string) (entity.User, error)
	VerifySubmit(ctx context.Context, user string) (entity.User, error)
}

type userImplementation struct {
	conn postgres.Adapter
}

func (r *userImplementation) VerifyLogin(ctx context.Context, user string) (entity.User, error) {
	fmt.Println(user)
	state := true
	query := `UPDATE users SET logged = $2 WHERE username = $1
	RETURNING id, username, role, password, logged`

	row := r.conn.QueryRow(ctx, query, user, state)

	u := entity.User{}

	err := row.Scan(&u.ID, &u.Username, &u.Role, &u.Password, &u.Logged)

	fmt.Println(err)

	if err != nil {
		return entity.User{}, err
	}
	fmt.Println(u)
	return u, nil
}

func (r *userImplementation) VerifyLogout(ctx context.Context, user string) (entity.User, error) {
	fmt.Println(user)
	state := false
	query := `UPDATE users SET logged = $2 WHERE username = $1
	RETURNING id, username, role, password, logged`

	row := r.conn.QueryRow(ctx, query, user, state)

	u := entity.User{}

	err := row.Scan(&u.ID, &u.Username, &u.Role, &u.Password, &u.Logged)

	fmt.Println(err)

	if err != nil {
		return entity.User{}, err
	}
	fmt.Println(u)
	return u, nil
}

func (r *userImplementation) VerifySubmit(ctx context.Context, user string) (entity.User, error) {
	fmt.Println(user)
	query := `SELECT id, username, role, password, logged FROM users WHERE username = $1`

	row := r.conn.QueryRow(ctx, query, user)

	u := entity.User{}

	err := row.Scan(&u.ID, &u.Username, &u.Role, &u.Password, &u.Logged)

	fmt.Println(err)

	if err != nil {
		return entity.User{}, err
	}
	fmt.Println(u)
	return u, nil
}

func (u userImplementation) Create(ctx context.Context, us *entity.User) (*uuid.UUID, error) {
	state := false
	const role int8 = 0
	query := `
	INSERT INTO users (username, password, name, role, logged)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	row := u.conn.QueryRow(ctx, query, us.Username, us.Password, us.Name, role, state)

	user := entity.User{}

	err := row.Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (r *userImplementation) Verify(ctx context.Context, username string) (user entity.User, err error) {
	query := `SELECT id, username, password, role FROM users WHERE username = $1`

	err = r.conn.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userImplementation) GrantAdmin(ctx context.Context, id string) error {
	fmt.Println(id)
	role := "1"
	query := `UPDATE users SET role = $2 WHERE id = $1`

	rows, err := r.conn.Exec(ctx, query, id, role)

	if err != nil {
		return err
	}
	updated, _ := rows.RowsAffected()
	if updated <= 0 {
		err = fmt.Errorf("user already an admin")
		return err
	}

	return nil
}

func (r *userImplementation) DemoteToUser(ctx context.Context, id string) error {
	role := "0"
	query := `UPDATE users SET role = $2 WHERE id = $1`

	rows, err := r.conn.Exec(ctx, query, id, role)

	if err != nil {
		return err
	}
	updated, _ := rows.RowsAffected()
	if updated <= 0 {
		err = fmt.Errorf("can't demote role user")
		return err
	}

	return nil
}
func (r *userImplementation) GetAllUser(ctx context.Context) ([]entity.User, error) {
	query := `SELECT id, username, password, role, name, logged FROM users`

	queries, err := r.conn.QueryRows(ctx, query)

	if err != nil {
		err = fmt.Errorf("executing query error : %w", err)
		return nil, err
	}

	users := []entity.User{}

	for queries.Next() {
		var user entity.User

		err = queries.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Role,
			&user.Name,
			&user.Logged,
		)

		if err != nil {
			err = fmt.Errorf("scanning user: %w", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}

func NewUserImplementation(conn postgres.Adapter) *userImplementation {
	return &userImplementation{
		conn: conn,
	}
}

// query := `
// INSERT INTO users (
// 	username,
// 	password,
// 	name,
// 	role
// )
// VALUES ($1,$2,$3,$4)
// RETURNING id
// `
// fmt.Printf("err1 \n")
// // res, err := u.conn.Exec(ctx, query, us.Username, us.Password, us.Name, us.Role)
// // if err != nil {
// // 	fmt.Printf("err2 \n")
// // 	fmt.Println(err)
// // 	err = fmt.Errorf("create data to db : %w", err)
// // 	return err
// // }
// // fmt.Printf("err3 \n")
// // fmt.Println(res)

// row := u.conn.QueryRow(ctx, query, us.Username, us.Password, us.Name, us.Role)

// user := entity.User{}

// err = row.Scan(&user.ID)

// if err != nil{

// }

// return &user.ID, nil
// user.ID, err = res.LastInsertId()
// if err != nil {
// 	fmt.Printf("err4 \n")
// 	err = fmt.Errorf("inserting last id to db : %w", err)
// 	return err
// }
// fmt.Printf("err5 \n")
// return nil
// return &user.ID, nil
