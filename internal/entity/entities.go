package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password" db:"password"`
	Name     string    `json:"name" db:"name"`
	Role     int8      `json:"role" db:"role"`
	Logged   bool      `json:"logged" db:"logged"`
}

type Room struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoomName  string    `json:"room_name" db:"room_name"`
	CreatedBy string    `json:"created_by" db:"created_by"`
}

type Booking struct {
	ID         uuid.UUID `json:"id" db:"id"`
	BookedBy   string    `json:"booked_by" db:"booked_by"`
	BookedRoom string    `json:"booked_room" db:"booked_room"`
	BookedDay  string    `json:"booked_day" db:"booked_day"`
	SesA       bool      `json:"sesa" db:"sesa"`
	SesB       bool      `json:"sesb" db:"sesb"`
	SesC       bool      `json:"sesc" db:"sesc"`
	SesD       bool      `json:"sesd" db:"sesd"`
	SesE       bool      `json:"sese" db:"sese"`
	SesF       bool      `json:"sesf" db:"sesf"`
	SesG       bool      `json:"sesg" db:"sesg"`
	SesH       bool      `json:"sesh" db:"sesh"`
	SesI       bool      `json:"sesi" db:"sesi"`
	SesJ       bool      `json:"sesj" db:"sesj"`
}

type Day struct {
	ID        uuid.UUID `json:"id" db:"id"`
	DayName   string    `json:"day_name" db:"day_name"`
	CreatedBy string    `json:"created_by" db:"created_by"`
}

type Login struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password" db:"password"`
	Role     int8      `json:"role" db:"role"`
	Logged   bool      `json:"logged" db:"logged"`
}
