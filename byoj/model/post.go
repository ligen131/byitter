package model

import (
	"byoj/utils/logs"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint32         `json:"post_id"    form:"post_id"    query:"post_id"   gorm:"primaryKey;unique;not null"`
	CreatedAt time.Time      `json:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at" query:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" form:"deleted_at" query:"deleted_at"`
	AuthorID  uint32         `json:"user_id"    form:"user_id"    query:"user_id"   gorm:"column:user_id;not null"`
	Time      time.Time      `json:"time"       form:"time"       query:"time"`
	Content   string         `json:"content"    form:"content"    query:"content"`
	IsPublic  bool           `json:"is_public"  form:"is_public"  query:"is_public" gorm:"not null"`
}

func CreatePost(authorID uint32, _time time.Time, content string, isPublic bool) (Post, error) {
	m := GetModel()
	defer m.Close()

	post := Post{
		AuthorID: authorID,
		Time:     _time,
		Content:  content,
		IsPublic: isPublic,
	}
	result := m.tx.Create(&post)
	if result.Error != nil {
		logs.Warn("Create user failed.", zap.Error(result.Error))
		m.Abort()
		return post, result.Error
	}

	m.tx.Commit()
	return post, nil
}

func FindPostByPostID(postID uint32) (Post, error) {
	m := GetModel()
	defer m.Close()

	var post Post
	result := m.tx.First(&post, postID)
	if result.Error != nil {
		logs.Info("Find post by id failed.", zap.Error(result.Error))
		m.Abort()
		return post, result.Error
	}

	m.tx.Commit()
	return post, nil
}

/**
 * 获取帖子列表
 * @param: authorID 发表人 user_id，为 0 不限制
 * @param: startTime 帖子起始时间往前筛选
 * @param: isPublic 是否筛选 is_public 只为 true 的帖子，为 false 也会返回 is_public 为 true 的帖子
 * @param: orderBy 结果排序方式，可选："time", "random"
 * @param: limit 限制结果数量
 **/
func GetPostsList(authorID uint32, startTime time.Time, isPublic bool, orderBy string, limit int) ([]Post, error) {
	m := GetModel()
	defer m.Close()

	var posts []Post
	result := m.tx.Model(&Post{})
	if authorID > 0 {
		result = result.Where("user_id = ?", authorID)
	}
	if startTime != time.Unix(0, 0) {
		result = result.Where("time <= ?", startTime)
	}
	if isPublic {
		result = result.Where("is_public = ?", true)
	}
	if orderBy == "time" {
		result = result.Order("time desc")
	} else {
		result = result.Order("random()")
	}
	if limit <= 0 {
		limit = 20
	}
	result = result.Limit(limit)

	result.Find(&posts)
	if result.Error != nil {
		logs.Info("Find posts list failed.", zap.Error(result.Error))
		m.Abort()
		return posts, result.Error
	}

	m.tx.Commit()
	return posts, nil
}
