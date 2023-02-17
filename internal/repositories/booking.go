package repositories

import (
	"context"
	"fmt"

	"GolangBookingApp/internal/entity"
	"GolangBookingApp/pkg/postgres"

	"github.com/google/uuid"
)

type Booking interface {
	Create(context.Context, *entity.Booking) (*uuid.UUID, error)
	VerifyBookingBy(context.Context, string) (entity.User, error)
	VerifyBookingEntity(context.Context, *entity.Booking) (entity.Booking, error)
	VerifyBookingDay(context.Context, string) (entity.Day, error)
	GetBookByDayRoom(context.Context, string, string) ([]entity.Booking, error)
	GetBookByName(context.Context, string) ([]entity.Booking, error)
	Delete(context.Context, string) error
	DeleteAll(context.Context) error
}

type bookingImplementation struct {
	conn postgres.Adapter
}

func NewBookingImplementation(conn postgres.Adapter) *bookingImplementation {
	return &bookingImplementation{
		conn: conn,
	}
}

func (b bookingImplementation) Create(ctx context.Context, eb *entity.Booking) (*uuid.UUID, error) {

	queryuser := `
	SELECT 
	username
	FROM users
	WHERE username = $1
	`
	udata := b.conn.QueryRow(ctx, queryuser, eb.BookedBy)

	err := udata.Scan(
		&eb.BookedBy,
	)

	if err != nil {
		err = fmt.Errorf("error : %w", err)
		return nil, err
	}

	query := `
	INSERT INTO bookings
	(booked_by, booked_room, booked_day, sesa, sesb, sesc, sesd, sese, sesf, sesg, sesh, sesi, sesj)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
	RETURNING id
	`
	row := b.conn.QueryRow(ctx, query, eb.BookedBy, eb.BookedRoom, eb.BookedDay,
		eb.SesA, eb.SesB, eb.SesC, eb.SesD, eb.SesE, eb.SesF, eb.SesG, eb.SesH, eb.SesI, eb.SesJ)

	book := entity.Booking{}

	err = row.Scan(&book.ID)

	if err != nil {
		return nil, err
	}

	return &book.ID, nil
}

func (b *bookingImplementation) VerifyBookingBy(ctx context.Context, bookedby string) (user entity.User, err error) {
	query := `SELECT id, username, password, role FROM users WHERE username = $1`

	err = b.conn.QueryRow(ctx, query, bookedby).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		err = fmt.Errorf("user doesn't exist")
		return entity.User{}, err
	}

	return user, nil
}

func (bd *bookingImplementation) VerifyBookingDay(ctx context.Context, bookedday string) (day entity.Day, err error) {
	query := `SELECT id, day_name, created_by FROM days WHERE day_name = $1`

	err = bd.conn.QueryRow(ctx, query, bookedday).Scan(&day.ID, &day.DayName, &day.CreatedBy)

	if err != nil {
		err = fmt.Errorf("day input is invalid")
		return entity.Day{}, err
	}

	return day, nil
}

func (be *bookingImplementation) VerifyBookingEntity(ctx context.Context, eb *entity.Booking) (b entity.Booking, err error) {
	query := `SELECT id, booked_by, booked_room, booked_day, sesa, sesb, sesc, sesd, sese, sesf, sesg, sesh, sesi, sesj
	 FROM bookings WHERE booked_day = $1 
	 AND booked_room = $2 
	 AND sesa = $3 AND sesb = $4 AND sesc = $5 AND sesd = $6 AND sese = $7 AND sesf = $8 AND sesg = $9 AND sesh = $10 AND sesi = $11 AND sesj = $12`

	err = be.conn.QueryRow(ctx, query,
		eb.BookedDay, eb.BookedRoom,
		eb.SesA, eb.SesB, eb.SesC, eb.SesD, eb.SesE, eb.SesF, eb.SesG, eb.SesH, eb.SesI, eb.SesJ).
		Scan(
			&b.ID,
			&b.BookedBy,
			&b.BookedRoom, &b.BookedDay,
			&b.SesA, &b.SesB, &b.SesC, &b.SesD, &b.SesE, &b.SesF, &b.SesG, &b.SesH, &b.SesI, &b.SesJ)

	fmt.Println(err)

	if err != nil {
		return entity.Booking{}, err
	}
	return b, nil
}

func (r *bookingImplementation) GetBookByDayRoom(ctx context.Context, bookedday string, bookedroom string) ([]entity.Booking, error) {
	query := `
	SELECT id, booked_by, booked_room, booked_day, sesa, sesb, sesc, sesd, sese, sesf, sesg, sesh, sesi, sesj 
	FROM bookings 
	WHERE booked_day = $1 AND booked_room = $2
	`
	queries, err := r.conn.QueryRows(ctx, query, bookedday, bookedroom)

	if err != nil {
		err = fmt.Errorf("executing query error : %w", err)
		return nil, err
	}

	books := []entity.Booking{}

	for queries.Next() {
		var b entity.Booking

		err = queries.Scan(
			&b.ID,
			&b.BookedBy,
			&b.BookedRoom,
			&b.BookedDay,
			&b.SesA, &b.SesB, &b.SesC, &b.SesD, &b.SesE, &b.SesF, &b.SesG, &b.SesH, &b.SesI, &b.SesJ,
		)

		if err != nil {
			err = fmt.Errorf("scanning bookings: %w", err)
			return nil, err
		}
		books = append(books, b)
	}

	return books, err
}

func (r *bookingImplementation) GetBookByName(ctx context.Context, bookedby string) ([]entity.Booking, error) {
	query := `
	SELECT id, booked_by, booked_room, booked_day, sesa, sesb, sesc, sesd, sese, sesf, sesg, sesh, sesi, sesj 
	FROM bookings 
	WHERE booked_by = $1
	`
	queries, err := r.conn.QueryRows(ctx, query, bookedby)

	if err != nil {
		err = fmt.Errorf("executing query error : %w", err)
		return nil, err
	}

	books := []entity.Booking{}

	for queries.Next() {
		var b entity.Booking

		err = queries.Scan(
			&b.ID,
			&b.BookedBy,
			&b.BookedRoom,
			&b.BookedDay,
			&b.SesA, &b.SesB, &b.SesC, &b.SesD, &b.SesE, &b.SesF, &b.SesG, &b.SesH, &b.SesI, &b.SesJ,
		)

		if err != nil {
			err = fmt.Errorf("scanning bookings: %w", err)
			return nil, err
		}
		books = append(books, b)
	}

	return books, err
}

func (r bookingImplementation) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM bookings WHERE id=$1
	`

	res, err := r.conn.Exec(ctx, query, id)

	if err != nil {
		err = fmt.Errorf("executing querry error : %w", err)
		return err
	}

	deletedRow, _ := res.RowsAffected()
	if deletedRow <= 0 {
		err = fmt.Errorf("id not found")
		return err
	}
	return nil
}

func (r bookingImplementation) DeleteAll(ctx context.Context) error {
	query := `
	TRUNCATE bookings
	`

	_, err := r.conn.Exec(ctx, query)

	if err != nil {
		err = fmt.Errorf("executing querry error : %w", err)
		return err
	}
	return nil
}
