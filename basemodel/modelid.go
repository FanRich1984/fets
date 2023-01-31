package basemodel

import "time"

type ModelId struct {
	Id        uint64    `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 创建ModelId的对象
func NewModelId(id uint64) ModelId {
	return ModelId{
		Id: id,
	}
}
