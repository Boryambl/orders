package model

import "time"

type User struct {
	UUID       string     `json:"uuid"`
	Firstname  string     `json:"firstName,omitempty"`
	Surname    string     `json:"surname,omitempty"`
	Middlename string     `json:"middlename,omitempty"`
	FIO        string     `json:"fio,omitempty"`
	Sex        bool       `json:"sex,omitempty"`
	Age        int        `json:"age,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}
