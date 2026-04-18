package Models

type Message struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	UUID        string `gorm:"type:varchar(36);not null"`
	Content     string `gorm:"type:text;not null"`         // Content of the message (in text, markdown format or plain text)
	Severity    string `gorm:"type:varchar(255);not null"` // Severity of the message to display (Info, Warning, Error)
	MediaTypeId int    `gorm:"not null"`

	// FK
	MediaType MediaType `gorm:"foreignKey:MediaTypeId;references:ID"` // MediaType of the message (text, markdown, etc.)
}

func (l *Message) TableName() string {
	return "messages"
}
