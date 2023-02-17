package repositories

import (
	"context"
	"fmt"

	"GolangBookingApp/internal/entity"
	"GolangBookingApp/pkg/postgres"

	"github.com/google/uuid"
)

type Room interface {
	Create(context.Context, *entity.Room) (*uuid.UUID, error)
	VerifyRoom(context.Context, string) (entity.User, error)
	VerifyRoomEntity(context.Context, string) (entity.Room, error)
	GetAllRoom(context.Context) ([]entity.Room, error)
	DeleteRoom(context.Context, string) error
}

type roomImplementation struct {
	conn postgres.Adapter
}

func NewRoomImplementation(conn postgres.Adapter) *roomImplementation {
	return &roomImplementation{
		conn: conn,
	}
}

func (r roomImplementation) Create(ctx context.Context, er *entity.Room) (*uuid.UUID, error) {

	queryuser := `
	SELECT 
	username
	FROM users
	WHERE username = $1
	`
	udata := r.conn.QueryRow(ctx, queryuser, er.CreatedBy)

	err := udata.Scan(
		&er.CreatedBy,
	)

	if err != nil {
		err = fmt.Errorf("error : %w", err)
		return nil, err
	}

	query := `
	INSERT INTO rooms
	(room_name, created_by)
	VALUES ($1, $2) 
	RETURNING id
	`
	row := r.conn.QueryRow(ctx, query, er.RoomName, er.CreatedBy)

	room := entity.Room{}

	err = row.Scan(&room.ID)

	if err != nil {
		return nil, err
	}

	return &room.ID, nil
}

func (r *roomImplementation) VerifyRoom(ctx context.Context, createdby string) (user entity.User, err error) {
	query := `SELECT id, username, password, role FROM users WHERE username = $1`

	err = r.conn.QueryRow(ctx, query, createdby).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		err = fmt.Errorf("user doesn't exist")
		return entity.User{}, err
	}

	return user, nil
}

func (r *roomImplementation) VerifyRoomEntity(ctx context.Context, roomname string) (room entity.Room, err error) {
	query := `SELECT id, room_name, created_by FROM rooms WHERE room_name = $1`

	err = r.conn.QueryRow(ctx, query, roomname).Scan(&room.ID, &room.RoomName, &room.CreatedBy)

	if err != nil {
		return entity.Room{}, err
	}

	return room, nil
}

func (r *roomImplementation) GetAllRoom(ctx context.Context) ([]entity.Room, error) {
	query := `SELECT id, room_name, created_by FROM rooms`

	queries, err := r.conn.QueryRows(ctx, query)

	if err != nil {
		err = fmt.Errorf("executing query error : %w", err)
		return nil, err
	}

	rooms := []entity.Room{}

	for queries.Next() {
		var room entity.Room

		err = queries.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedBy,
		)

		if err != nil {
			err = fmt.Errorf("scanning bookings: %w", err)
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, err

}

func (r roomImplementation) DeleteRoom(ctx context.Context, id string) error {
	query := `
	DELETE FROM rooms WHERE id=$1
	`

	res, err := r.conn.Exec(ctx, query, id)

	if err != nil {
		err = fmt.Errorf("executing query error : %w", err)
		return err
	}

	deletedRow, _ := res.RowsAffected()
	if deletedRow <= 0 {
		err = fmt.Errorf("id not found")
		return err
	}
	return nil
}
