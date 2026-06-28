package Models

type Role struct {
	ID   int    `gorm:"primary_key;auto_increment:false"` // 1 -> admin; 2 -> write; 3 -> read
	Name string `gorm:"type:varchar(255);not null"`
}
