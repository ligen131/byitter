package controllers

import (
	"byoj/model"
	"byoj/utils/logs"
	"net/http"
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
	err := c.Bind(&postRequest)
	if err != nil {
		logs.Warn("Failed to parse request data.", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Failed to parse request data.",
				Err:     err.Error(),
			},
		})
	}
	logs.Debug("Post struct:", zap.Any("postRequest", postRequest))

	user := model.User{}
	if postRequest.AuthorID != 0 {
		user, err = model.FindUserByID(postRequest.AuthorID)
	} else if postRequest.AuthorEmail != "" {
		user, err = model.FindUserByEmail(postRequest.AuthorEmail)
	} else if postRequest.AuthorName != "" {
		user, err = model.FindUserByName(postRequest.AuthorName)
	} else {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Author ID, email or username is required.",
				Err:     "",
			},
		})
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, ResponseStruct{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Data: ErrorMessage{
					Message: "User not found.",
					Err:     err.Error(),
				},
			})
		}
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Find user failed.",
				Err:     err.Error(),
			},
		})
	}

	if user.Deleted {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "This user has been deleted.",
				Err:     "",
			},
		})
	}

	if !user.Verified {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "This user has not been verified.",
				Err:     "",
			},
		})
	}

	post, err := model.CreatePost(user.ID, time.Now(), postRequest.Content, true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Failed to create post into database.",
				Err:     err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data: PostCreateResponse{
			Status:   "Create post successfully.",
			PostID:   post.ID,
			IsPublic: post.IsPublic,
		},
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
	err := c.Bind(&postRequest)
	if err != nil {
		logs.Warn("Failed to parse request data.", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Failed to parse request data.",
				Err:     err.Error(),
			},
		})
	}
	logs.Debug("Post struct:", zap.Any("postRequest", postRequest))

	user := model.User{}
	if postRequest.AuthorID != 0 {
		user, err = model.FindUserByID(postRequest.AuthorID)
	} else if postRequest.AuthorEmail != "" {
		user, err = model.FindUserByEmail(postRequest.AuthorEmail)
	} else if postRequest.AuthorName != "" {
		user, err = model.FindUserByName(postRequest.AuthorName)
	} else {
		user.ID = 0
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, ResponseStruct{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Data: ErrorMessage{
					Message: "User not found.",
					Err:     err.Error(),
				},
			})
		}
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Find user failed.",
				Err:     err.Error(),
			},
		})
	}

	mp := make(map[uint32]model.User)
	if user.ID != 0 {
		mp[user.ID] = user
	}

	posts, err := model.GetPostsList(user.ID, time.Unix(postRequest.StartTime, 0), false, postRequest.OrderBy, postRequest.Limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Get posts list failed.",
				Err:     err.Error(),
			},
		})
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

	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    resp,
	})
}
