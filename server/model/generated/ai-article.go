package generated

import "time"

type Aiarticle struct {
	Id int `json:"id" gorm:"column:id"`
	CategoryId int `json:"categoryId" gorm:"column:category_id"`
	Title string `json:"title" gorm:"column:title"`
	Author *string `json:"author" gorm:"column:author"`
	Image *string `json:"image" gorm:"column:image"`
	Describe *string `json:"describe" gorm:"column:describe"`
	Content *string `json:"content" gorm:"column:content"`
	Views int `json:"views" gorm:"column:views"`
	Sort int `json:"sort" gorm:"column:sort"`
	Status int `json:"status" gorm:"column:status"`
	IsLink int `json:"isLink" gorm:"column:is_link"`
	LinkUrl *string `json:"linkUrl" gorm:"column:link_url"`
	IsHot int `json:"isHot" gorm:"column:is_hot"`
	CreatedBy int `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy int `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"`
	DeleteTime *time.Time `json:"deleteTime" gorm:"column:delete_time"`
}

func (Aiarticle) TableName() string {
	return "ai_article"
}
