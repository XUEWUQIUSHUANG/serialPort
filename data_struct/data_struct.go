package datastruct

import "time"

type Config struct {
	Password string `json:"password"`
	Port     string `json:"port"`
	IP       string `json:"ip"`
	Database string `json:"database"`
}

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"column:username;type:varchar(50);unique;not null"`
	Password  string    `gorm:"column:password;type:varchar(255);not null"`
	Cookie    string    `gorm:"column:cookie;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

type UserLogin struct {
	Username string `gorm:"column:username;type:varchar(50);unique;not null"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
}

type EnvironmentData struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	Temperature    float64   `gorm:"column:temperature;type:double;not null"`
	Humidity       float64   `gorm:"column:humidity;type:double;not null"`
	LightIntensity float64   `gorm:"column:light_intensity;type:double;not null"`
	RecordedAt     time.Time `gorm:"column:recorded_at;type:timestamp;default:CURRENT_TIMESTAMP;not null"`
}
