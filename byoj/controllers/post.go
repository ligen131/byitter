package controllers

import (
	"byoj/model"
	"byoj/utils/logs"
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostCreateRequest struct {
	AuthorID    uint32 `json:"user_id"`
	AuthorName  string `json:"user_name"`
	AuthorEmail string `json:"email"`
	Content     string `json:"content"`
}

type PostCreateResponse struct {
	Status   string `json:"status"`
	PostID   uint32 `json:"post_id"`
	IsPublic bool   `json:"is_public"`
}

func PostPOST(c echo.Context) error {
	logs.Debug("POST /post")

	postRequest := PostCreateRequest{}
	_ok, err := Bind(c, &postRequest)
	if !_ok {
		return err
	}

	user, err, e500 := FindUser(c, model.User{
		ID:       postRequest.AuthorID,
		UserName: postRequest.AuthorName,
		Email:    postRequest.AuthorEmail,
	})
	if e500 {
		return err
	}
	if err != nil {
		return ResponseBadRequest(c, "Find user failed.", err)
	}

	if user.Deleted {
		return ResponseBadRequest(c, "This user has been deleted.", nil)
	}

	if !user.Verified {
		return ResponseBadRequest(c, "This user has not been verified.", nil)
	}

	post, err := model.CreatePost(user.ID, time.Now(), postRequest.Content, true)
	if err != nil {
		return ResponseInternalServerError(c, "Failed to create post into database.", err)
	}

	return ResponseOK(c, PostCreateResponse{
		Status:   "Create post successfully.",
		PostID:   post.ID,
		IsPublic: post.IsPublic,
	})
}

type PostGetRequest struct {
	AuthorID    uint32 `json:"user_id"`
	AuthorName  string `json:"user_name"`
	AuthorEmail string `json:"email"`
	Limit       int    `json:"limit"`
	OrderBy     string `json:"order_by"`
	StartTime   int64  `json:"start_time"`
}

type PostResponse struct {
	AuthorID    uint32 `json:"user_id"`
	AuthorName  string `json:"user_name"`
	AuthorEmail string `json:"email"`
	PostID      uint32 `json:"post_id"`
	Time        int64  `json:"time"`
	Content     string `json:"content"`
	IsPublic    bool   `json:"is_public"`
}

type PostGetResponse struct {
	PostList []PostResponse `json:"post_list"`
}

func PostGET(c echo.Context) error {
	logs.Debug("GET /post")

	postRequest := PostGetRequest{}
	_ok, err := Bind(c, &postRequest)
	if !_ok {
		return err
	}

	user, err, e500 := FindUser(c, model.User{
		ID:       postRequest.AuthorID,
		UserName: postRequest.AuthorName,
		Email:    postRequest.AuthorEmail,
	})
	if e500 {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return ResponseBadRequest(c, "User not found.", err)
	} else {
		user.ID = 0
	}

	mp := make(map[uint32]model.User)
	if user.ID != 0 {
		mp[user.ID] = user
	}

	posts, err := model.GetPostsList(user.ID, time.Unix(postRequest.StartTime, 0), false, postRequest.OrderBy, postRequest.Limit)
	if err != nil {
		return ResponseInternalServerError(c, "Get posts list failed.", err)
	}

	resp := PostGetResponse{
		PostList: make([]PostResponse, 0),
	}
	for _, post := range posts {
		if mp[post.AuthorID].ID == 0 {
			mp[post.AuthorID], err = model.FindUserByID(post.AuthorID)
			if err != nil {
				logs.Warn("Find user for post failed.", zap.Uint32("AuthorID", post.AuthorID), zap.Error(err))
			}
		}
		resp.PostList = append(resp.PostList, PostResponse{
			AuthorID:    post.AuthorID,
			AuthorName:  mp[post.AuthorID].UserName,
			AuthorEmail: mp[post.AuthorID].Email,
			PostID:      post.ID,
			Time:        post.Time.Unix(),
			Content:     post.Content,
			IsPublic:    post.IsPublic,
		})
	}

	return ResponseOK(c, resp)
}
