package dao

import "time"

type Record struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	UserID        int        `json:"user_id"`
	Channel       uint64     `json:"channel"`
	Date          *time.Time `json:"scheduling_date" gorm:"type:datetime"`
	ExpectedStart *time.Time `json:"expected_start" gorm:"type:datetime"`
	ExpectedEnd   *time.Time `json:"expected_end" gorm:"type:datetime"`
	Start         *time.Time `json:"start_date" gorm:"type:datetime"`
	End           *time.Time `json:"end_date" gorm:"type:datetime"`
	Status        int        `json:"status"`
	//0:scheduled, 1:in progress, 2:completed, 3:canceled, 4:failed
}
