package consumercreditlimit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"

	"case-study-kredit-plus/library"
	"case-study-kredit-plus/library/helpers"
	"case-study-kredit-plus/middleware"
	"case-study-kredit-plus/models"
	"case-study-kredit-plus/src/services/consumercreditlimit"

	"github.com/gin-gonic/gin"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/library/http/response"
	"case-study-kredit-plus/library/types"

	consumercreditlimitRepository "case-study-kredit-plus/src/services/consumercreditlimit/repository"
	consumercreditlimitUsecase "case-study-kredit-plus/src/services/consumercreditlimit/usecase"
)

var ()

type ConsumerCreditLimitHandler struct {
	ConsumerCreditLimitUsecase consumercreditlimit.Usecase
	dataManager                *data.Manager
	Result                     gin.H
	Status                     int
}

func (h ConsumerCreditLimitHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	consumercreditlimitRepo := consumercreditlimitRepository.NewConsumerCreditLimitRepository(
		data.NewMySQLStorage(db, "consumer_credit_limits", models.ConsumerCreditLimit{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uConsumerCreditLimit := consumercreditlimitUsecase.NewConsumerCreditLimitUsecase(db, &consumercreditlimitRepo)

	base := &ConsumerCreditLimitHandler{ConsumerCreditLimitUsecase: uConsumerCreditLimit, dataManager: dataManager}

	rs := v.Group("/consumers/credit-limits")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		// rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/status", middleware.Auth, base.UpdateStatus)
	}

	status := v.Group("/statuses")
	{
		status.GET("/consumers/credit-limits", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *ConsumerCreditLimitHandler) FindAll(c *gin.Context) {
	if c.Query("ConsumerID") != "" && !library.ValidateUUID(c.Query("ConsumerID")) {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->FindAll()",
			Message:    "Consumer ID is not valid",
			Error:      fmt.Errorf("Consumer ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	var params models.FindAllConsumerCreditLimitParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.ConsumerID = c.Query("ConsumerID")
	datas, err := h.ConsumerCreditLimitUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	length, err := h.ConsumerCreditLimitUsecase.Count(c, params)
	if err != nil {
		err.Path = ".ConsumerCreditLimitHandler->FindAll()" + err.Path
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

func (h *ConsumerCreditLimitHandler) Find(c *gin.Context) {
	id := c.Param("id")

	if id != "" && !library.ValidateUUID(id) {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Find()",
			Message:    "ID is not valid",
			Error:      fmt.Errorf("ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.ConsumerCreditLimitUsecase.Find(c, id)
	if err != nil {
		err.Path = ".ConsumerCreditLimitHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "ConsumerCreditLimit not found", http.StatusUnprocessableEntity, *err)
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

func (h *ConsumerCreditLimitHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.ConsumerCreditLimit
	var data *models.ConsumerCreditLimit

	if c.PostForm("ConsumerID") != "" && !library.ValidateUUID(c.PostForm("ConsumerID")) {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Create()",
			Message:    "Consumer ID is not valid",
			Error:      fmt.Errorf("Consumer ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month1, errParseFloat := strconv.ParseFloat(c.PostForm("Month1"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Create()",
			Message:    "1 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month2, errParseFloat := strconv.ParseFloat(c.PostForm("Month2"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Create()",
			Message:    "2 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month3, errParseFloat := strconv.ParseFloat(c.PostForm("Month3"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Create()",
			Message:    "3 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month6, errParseFloat := strconv.ParseFloat(c.PostForm("Month6"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Create()",
			Message:    "6 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	obj.ConsumerID = c.PostForm("ConsumerID")
	obj.Month1 = month1
	obj.Month2 = month2
	obj.Month3 = month3
	obj.Month6 = month6

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ConsumerCreditLimitUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".ConsumerCreditLimitHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data created successfuly", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ConsumerCreditLimitHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.ConsumerCreditLimit
	var data *models.ConsumerCreditLimit

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

	if c.PostForm("ConsumerID") != "" && !library.ValidateUUID(c.PostForm("ConsumerID")) {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Update()",
			Message:    "Consumer ID is not valid",
			Error:      fmt.Errorf("Consumer ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month1, errParseFloat := strconv.ParseFloat(c.PostForm("Month1"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Update()",
			Message:    "1 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month2, errParseFloat := strconv.ParseFloat(c.PostForm("Month2"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Update()",
			Message:    "2 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month3, errParseFloat := strconv.ParseFloat(c.PostForm("Month3"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Update()",
			Message:    "3 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	month6, errParseFloat := strconv.ParseFloat(c.PostForm("Month6"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerCreditLimitHandler->Update()",
			Message:    "6 Month Limit Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	obj.ConsumerID = c.PostForm("ConsumerID")
	obj.Month1 = month1
	obj.Month2 = month2
	obj.Month3 = month3
	obj.Month6 = month6

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ConsumerCreditLimitUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ConsumerCreditLimitHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "ConsumerCreditLimit successfuly updated", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ConsumerCreditLimitHandler) FindStatus(c *gin.Context) {
	datas, err := h.ConsumerCreditLimitUsecase.FindStatus(c)
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

func (h *ConsumerCreditLimitHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.ConsumerCreditLimit

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".ConsumerCreditLimitHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.ConsumerCreditLimitUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ConsumerCreditLimitHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Status update success", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
