package user

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"

	"case-study-kredit-plus/library"
	"case-study-kredit-plus/library/helpers"
	"case-study-kredit-plus/middleware"
	"case-study-kredit-plus/models"
	"case-study-kredit-plus/src/services/user"

	"github.com/gin-gonic/gin"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/library/http/response"
	"case-study-kredit-plus/library/types"

	userRepository "case-study-kredit-plus/src/services/user/repository"
	userUsecase "case-study-kredit-plus/src/services/user/usecase"
)

var ()

type UserHandler struct {
	UserUsecase user.Usecase
	dataManager *data.Manager
	Result      gin.H
	Status      int
}

func (h UserHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	userRepo := userRepository.NewUserRepository(
		data.NewMySQLStorage(db, "users", models.User{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uUser := userUsecase.NewUserUsecase(db, &userRepo)

	base := &UserHandler{UserUsecase: uUser, dataManager: dataManager}

	rs := v.Group("/users")
	{
		// rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.PUT("/:id", middleware.Auth, base.Update)
		// rs.PUT("/status", middleware.Auth, base.UpdateStatus)

		rs.POST("register", base.Create)
		rs.POST("auth/login", base.Login)
	}

	status := v.Group("/statuses")
	{
		status.GET("/users", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *UserHandler) FindAll(c *gin.Context) {
	var params models.FindAllUserParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.UserUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	length, err := h.UserUsecase.Count(c, params)
	if err != nil {
		err.Path = ".UserHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Success", StatusCode: http.StatusOK, Message: "Data shown successfuly", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *UserHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.UserUsecase.Find(c, id)
	if err != nil {
		err.Path = ".UserHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "User not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	result.Password = ""

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data shown successfuly", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.User
	var data *models.User

	id := c.Param("id")

	if !library.IsEmailValid(c.PostForm("Email")) {
		err := &types.Error{
			Path:       ".UserHandler->Update()",
			Message:    "Email is not valid",
			Error:      fmt.Errorf("email is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	obj.Name = c.PostForm("Name")
	obj.Email = c.PostForm("Email")
	obj.Username = c.PostForm("Username")
	obj.CountryCallingCode = c.PostForm("CountryCallingCode")
	obj.PhoneNumber = c.PostForm("PhoneNumber")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.UserUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}

		data.Password = ""

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "User successfuly updated", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) FindStatus(c *gin.Context) {
	datas, err := h.UserUsecase.FindStatus(c)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}
	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data successfuly shown", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.User

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".UserHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.UserUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Status update success", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// REGISTER
func (h *UserHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.User
	var data *models.User

	c.Set("UserID", "0")

	if !library.IsEmailValid(c.PostForm("Email")) {
		err := &types.Error{
			Path:       ".UserHandler->Create()",
			Message:    "Email is not valid",
			Error:      fmt.Errorf("email is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateCountryCode(c.PostForm("CountryCallingCode")) {
		err := &types.Error{
			Path:       ".UserHandler->Create()",
			Message:    "Country Calling Code is not valid",
			Error:      fmt.Errorf("country calling code is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidatePhoneNumber(c.PostForm("PhoneNumber")) {
		err := &types.Error{
			Path:       ".UserHandler->Create()",
			Message:    "Phone Number is not valid",
			Error:      fmt.Errorf("phone number is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	obj.Name = c.PostForm("Name")
	obj.Email = c.PostForm("Email")
	obj.Username = c.PostForm("Username")
	obj.CountryCallingCode = c.PostForm("CountryCallingCode")
	obj.PhoneNumber = c.PostForm("PhoneNumber")

	if len(c.PostForm("Password")) < 6 {
		err := &types.Error{
			Path:       ".UserHandler->Create()",
			Message:    "Password must be at least 6 characters",
			Error:      fmt.Errorf("password must be at least 6 characters"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))
	password := fmt.Sprintf("%x", hash.Sum(nil))

	obj.Password = password

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.UserUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		data.Password = ""

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data created successfuly", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// // ///

// LOGIN
func (h *UserHandler) Login(c *gin.Context) {
	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))

	email := c.PostForm("Email")
	password := fmt.Sprintf("%x", hash.Sum(nil))

	if !library.IsEmailValid(c.PostForm("Email")) {
		err := &types.Error{
			Path:       ".UserHandler->Login()",
			Message:    "Email is not valid",
			Error:      fmt.Errorf("email is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	var params models.FindAllUserParams
	params.Email = email
	params.Password = password
	params.FindAllParams.StatusID = `status_id = "1"`

	datas, err := h.UserUsecase.Login(c, params)
	if err != nil {
		c.JSON(401, response.ErrorResponse{
			Code:    "LoginFailed",
			Status:  "Warning",
			Message: "Login Failed",
			Data: &response.DataError{
				Message: err.Message,
				Status:  401,
			},
		})
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Login success", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// // //
