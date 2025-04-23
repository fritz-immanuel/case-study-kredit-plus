package consumer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"case-study-kredit-plus/library"
	"case-study-kredit-plus/library/helpers"
	"case-study-kredit-plus/middleware"
	"case-study-kredit-plus/models"
	"case-study-kredit-plus/src/services/consumer"

	"github.com/gin-gonic/gin"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/library/http/response"
	"case-study-kredit-plus/library/types"

	consumerRepository "case-study-kredit-plus/src/services/consumer/repository"
	consumerUsecase "case-study-kredit-plus/src/services/consumer/usecase"
)

var ()

type ConsumerHandler struct {
	ConsumerUsecase consumer.Usecase
	dataManager     *data.Manager
	Result          gin.H
	Status          int
}

func (h ConsumerHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	consumerRepo := consumerRepository.NewConsumerRepository(
		data.NewMySQLStorage(db, "consumers", models.Consumer{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uConsumer := consumerUsecase.NewConsumerUsecase(db, &consumerRepo)

	base := &ConsumerHandler{ConsumerUsecase: uConsumer, dataManager: dataManager}

	rs := v.Group("/consumers")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/status", middleware.Auth, base.UpdateStatus)
	}

	status := v.Group("/statuses")
	{
		status.GET("/consumers", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *ConsumerHandler) FindAll(c *gin.Context) {
	var params models.FindAllConsumerParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams

	if c.Query("MinDateOfBirth") != "" {
		_, errParseTime := time.Parse(library.StrToDateFormat, c.Query("MinDateOfBirth"))
		if errParseTime != nil {
			err := &types.Error{
				Path:       ".ConsumerHandler->FindAll()",
				Message:    "Min Date of Birth Invalid",
				Error:      errParseTime,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		params.MinDateOfBirth = c.Query("MinDateOfBirth")
	}

	if c.Query("MaxDateOfBirth") != "" {
		_, errParseTime := time.Parse(library.StrToDateFormat, c.Query("MaxDateOfBirth"))
		if errParseTime != nil {
			err := &types.Error{
				Path:       ".ConsumerHandler->FindAll()",
				Message:    "Max Date of Birth Invalid",
				Error:      errParseTime,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		params.MaxDateOfBirth = c.Query("MaxDateOfBirth")
	}

	if c.Query("MinSalary") != "" {
		minSalary, errParseFloat := strconv.ParseFloat(c.Query("MinSalary"), 64)
		if errParseFloat != nil {
			err := &types.Error{
				Path:       ".ConsumerHandler->FindAll()",
				Message:    "Invalid Min Salary Input",
				Error:      errParseFloat,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		params.MinSalary = minSalary
	}

	if c.Query("MaxSalary") != "" {
		maxSalary, errParseFloat := strconv.ParseFloat(c.Query("MaxSalary"), 64)
		if errParseFloat != nil {
			err := &types.Error{
				Path:       ".ConsumerHandler->FindAll()",
				Message:    "Invalid Max Salary Input",
				Error:      errParseFloat,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		params.MaxSalary = maxSalary
	}

	if c.Query("NIK") != "" && !library.ValidateNIK(c.Query("NIK")) {
		params.NIK = c.Query("NIK")
	}
	if c.Query("FullName") != "" && !library.ValidateTextInput(c.Query("FullName")) {
		params.FullName = c.Query("FullName")
	}
	if c.Query("LegalName") != "" && !library.ValidateTextInput(c.Query("LegalName")) {
		params.LegalName = c.Query("LegalName")
	}
	if c.Query("PlaceOfBirth") != "" && !library.ValidateTextInput(c.Query("PlaceOfBirth")) {
		params.PlaceOfBirth = c.Query("PlaceOfBirth")
	}

	datas, err := h.ConsumerUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	length, err := h.ConsumerUsecase.Count(c, params)
	if err != nil {
		err.Path = ".ConsumerHandler->FindAll()" + err.Path
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

func (h *ConsumerHandler) Find(c *gin.Context) {
	id := c.Param("id")

	if id != "" && !library.ValidateUUID(id) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Find()",
			Message:    "ID is not valid",
			Error:      fmt.Errorf("ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.ConsumerUsecase.Find(c, id)
	if err != nil {
		err.Path = ".ConsumerHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Consumer not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data shown successfuly", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ConsumerHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.Consumer
	var data *models.Consumer

	if !library.ValidateNIK(c.PostForm("NIK")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Create()",
			Message:    "NIK is invalid",
			Error:      fmt.Errorf("NIK is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateTextInput(c.PostForm("FullName")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Create()",
			Message:    "Full Name is invalid",
			Error:      fmt.Errorf("FullName is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateTextInput(c.PostForm("LegalName")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Create()",
			Message:    "Legal Name is invalid",
			Error:      fmt.Errorf("LegalName is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateTextInput(c.PostForm("PlaceOfBirth")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Create()",
			Message:    "Place Of Birth is invalid",
			Error:      fmt.Errorf("PlaceOfBirth is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("DateOfBirth") != "" {
		dob, errParseTime := time.Parse(library.StrToDateFormat, c.PostForm("DateOfBirth"))
		if errParseTime != nil {
			err := &types.Error{
				Path:       ".ConsumerHandler->Create()",
				Message:    "Date of Birth Invalid",
				Error:      errParseTime,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		obj.DateOfBirth = dob
	}

	salary, errParseFloat := strconv.ParseFloat(c.PostForm("Salary"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerHandler->Create()",
			Message:    "Salary Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	obj.Salary = salary

	obj.NIK = c.PostForm("NIK")
	obj.FullName = c.PostForm("FullName")
	obj.LegalName = c.PostForm("LegalName")
	obj.PlaceOfBirth = c.PostForm("PlaceOfBirth")
	obj.KTPImgURL = c.PostForm("KTPImgURL")
	obj.SelfieImgURL = c.PostForm("SelfieImgURL")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ConsumerUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".ConsumerHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data created successfuly", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ConsumerHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.Consumer
	var data *models.Consumer

	id := c.Param("id")

	if id != "" && !library.ValidateUUID(id) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Update()",
			Message:    "ID is not valid",
			Error:      fmt.Errorf("ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateNIK(c.PostForm("NIK")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Update()",
			Message:    "NIK is invalid",
			Error:      fmt.Errorf("NIK is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateTextInput(c.PostForm("FullName")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Update()",
			Message:    "Full Name is invalid",
			Error:      fmt.Errorf("FullName is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateTextInput(c.PostForm("LegalName")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Update()",
			Message:    "Legal Name is invalid",
			Error:      fmt.Errorf("LegalName is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if !library.ValidateTextInput(c.PostForm("PlaceOfBirth")) {
		err := &types.Error{
			Path:       ".ConsumerHandler->Update()",
			Message:    "Place Of Birth is invalid",
			Error:      fmt.Errorf("PlaceOfBirth is invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("DateOfBirth") != "" {
		dob, errParseTime := time.Parse(library.StrToDateFormat, c.PostForm("DateOfBirth"))
		if errParseTime != nil {
			err := &types.Error{
				Path:       ".ConsumerHandler->Update()",
				Message:    "Date of Birth Invalid",
				Error:      errParseTime,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		obj.DateOfBirth = dob
	}

	salary, errParseFloat := strconv.ParseFloat(c.PostForm("Salary"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerHandler->Update()",
			Message:    "Salary Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	obj.Salary = salary
	obj.NIK = c.PostForm("NIK")
	obj.FullName = c.PostForm("FullName")
	obj.LegalName = c.PostForm("LegalName")
	obj.PlaceOfBirth = c.PostForm("PlaceOfBirth")
	obj.KTPImgURL = c.PostForm("KTPImgURL")
	obj.SelfieImgURL = c.PostForm("SelfieImgURL")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ConsumerUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ConsumerHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Consumer successfuly updated", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ConsumerHandler) FindStatus(c *gin.Context) {
	datas, err := h.ConsumerUsecase.FindStatus(c)
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

func (h *ConsumerHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Consumer

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".ConsumerHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.ConsumerUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ConsumerHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Status update success", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
