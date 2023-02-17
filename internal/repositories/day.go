package repositories

import (
	"context"
	"fmt"

	"GolangBookingApp/internal/entity"
	"GolangBookingApp/pkg/postgres"

	"github.com/google/uuid"
)

type Day interface {
	NewCreateDay(context.Context, *entity.Day) (*uuid.UUID, error)
	VerifyDay(context.Context, string) (entity.User, error)
	GetDays(context.Context) ([]entity.Day, error)
}

type dayImplementation struct {
	conn postgres.Adapter
}

func (d dayImplementation) NewCreateDay(ctx context.Context, ed *entity.Day) (*uuid.UUID, error) {

	queryuser := `
	SELECT 
	username
	FROM users
	WHERE username = $1
	`
	udata := d.conn.QueryRow(ctx, queryuser, ed.CreatedBy)

	err := udata.Scan(
		&ed.CreatedBy,
	)

	if err != nil {
		err = fmt.Errorf("error : %w", err)
		return nil, err
	}

	query := `
	INSERT INTO days
	(day_name, created_by)
	VALUES ($1, $2) 
	RETURNING id
	`
	row := d.conn.QueryRow(ctx, query, ed.DayName, ed.CreatedBy)

	day := entity.Day{}

	err = row.Scan(&day.ID)

	if err != nil {
		return nil, err
	}

	return &day.ID, nil
}

func (d *dayImplementation) VerifyDay(ctx context.Context, createdby string) (user entity.User, err error) {
	query := `SELECT id, username, password, role FROM users WHERE username = $1`

	err = d.conn.QueryRow(ctx, query, createdby).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (d *dayImplementation) GetDays(ctx context.Context) ([]entity.Day, error) {
	query := `
	SELECT id, day_name, created_by FROM days
	`
	queries, err := d.conn.QueryRows(ctx, query)

	if err != nil {
		err = fmt.Errorf("executing query error : %w", err)
		return nil, err
	}

	days := []entity.Day{}

	for queries.Next() {
		var day entity.Day

		err = queries.Scan(
			&day.ID,
			&day.DayName,
			&day.CreatedBy,
		)

		if err != nil {
			err = fmt.Errorf("scanning bookings: %w", err)
			return nil, err
		}
		days = append(days, day)
	}

	return days, err
}

func NewDayImplementation(conn postgres.Adapter) *dayImplementation {
	return &dayImplementation{
		conn: conn,
	}
}
