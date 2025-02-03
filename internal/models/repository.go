package models

import "time"

type Topic struct {
	Topic         string    `db:"topic"`
	PodRun        int       `db:"pod_run"`
	ScriptVersion int       `db:"script_version"`
	ScriptName    string    `db:"script_name"`
	InActive      string    `db:"in_active"`
	CreatedAt     time.Time `db:"created_at"`
	CreateBy      string    `db:"create_by"`
	UpdateAt      time.Time `db:"update_at"`
	UpdateBy      string    `db:"update_by"`
}

type Script struct {
	ID         int       `db:"id"`
	ScriptKey  string    `db:"script_key"`
	ScriptName string    `db:"script_name"`
	Version    int       `db:"version"`
	IsDeleted  string    `db:"is_deleted"`
	CreatedAt  time.Time `db:"created_at"`
	CreateBy   string    `db:"create_by"`
}

type ConnectionPool struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	EndPoint  string    `db:"end_point"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Key       string    `db:"key"`
	IsDeleted string    `db:"is_deleted"`
	CreatedAt time.Time `db:"created_at"`
	CreateBy  string    `db:"create_by"`
	UpdateAt  time.Time `db:"update_at"`
	UpdateBy  string    `db:"update_by"`
}

type User struct {
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	IsDeleted string    `db:"is_deleted"`
	CreatedAt time.Time `db:"created_at"`
	CreateBy  string    `db:"create_by"`
	UpdateAt  time.Time `db:"update_at"`
	UpdateBy  string    `db:"update_by"`
}

type Role struct {
	RoleName  string    `db:"role_name"`
	IsDeleted string    `db:"is_deleted"`
	CreatedAt time.Time `db:"created_at"`
	CreateBy  string    `db:"create_by"`
	UpdateAt  time.Time `db:"update_at"`
	UpdateBy  string    `db:"update_by"`
}
