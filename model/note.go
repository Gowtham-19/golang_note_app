package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//db model
type Notes struct {
	Id            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content_Type  string             `json:"content_type,omitempty" bson:"content_type,omitempty"`
	Subject       string             `json:"subject,omitempty" bson:"subject,omitempty"`
	Description   string             `json:"description,omitempty" bson:"description,omitempty"`
	Status        string             `json:"status,omitempty" bson:"status,omitempty"`
	Created_Date  int64              `json:"created_date,omitempty" bson:"created_date,omitempty"`
	Created_Month int64              `json:"created_month,omitempty" bson:"created_month,omitempty"`
	Created_Year  int64              `json:"created_year,omitempty" bson:"created_year,omitempty"`
	CreatedAt     *time.Time         `json:"created_at,omitempty" bson:"createdat,omitempty"`
	UpdatedAt     *time.Time         `json:"update_at,omitempty"  bson:"updatedat,omitempty"`
}
