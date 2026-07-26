package Models

type LabelSourceType struct {
	ID      int    `gorm:"primary_key"`   // 1 => manifest, 2 => api, 3 => internal
	Name    string `gorm:"size:255"`      // 1 => manifest, 2 => api, 3 => internal
	Mutable bool   `gorm:"default:false"` // 1 => false, 2 => true, 3 => false
}

func (LabelSourceType) TableName() string {
	return "labelsourcetype"
}
