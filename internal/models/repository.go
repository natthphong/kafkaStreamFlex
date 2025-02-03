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
