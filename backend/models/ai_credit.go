package models

import "time"

type UserAICreditGrant struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	UserID         int64     `gorm:"column:user_id;index:idx_user_ai_credit_grants_user_month,priority:1;not null"`
	Credits        int64     `gorm:"not null"`
	EffectiveMonth time.Time `gorm:"column:effective_month;type:date;index:idx_user_ai_credit_grants_user_month,priority:2;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AICreditUsage struct {
	UsedCredit  int64     `json:"usedCredit"`
	LeftCredit  int64     `json:"leftCredit"`
	TotalCredit int64     `json:"totalCredit"`
	RenewDate   time.Time `json:"renewDate"`
}
