package mongodb

import (
	"time"
)

type Meta struct {
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}

func (d *Meta) Create() {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
}

func (d *Meta) Update() {
	d.UpdatedAt = time.Now()
}

func (d *Meta) Delete() {
	now := time.Now()
	d.UpdatedAt = now
	d.DeletedAt = &now
}
