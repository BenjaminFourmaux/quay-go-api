package Models

type Visibility struct {
	ID   int    `gorm:"primary_key;auto_increment:false"` // 1 = public , 2 = private
	Name string `gorm:"type:varchar(255);not null"`
}

func (Visibility) TableName() string {
	return "visibility"
}
