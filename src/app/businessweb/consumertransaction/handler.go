package consumertransaction

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
	"case-study-kredit-plus/src/services/consumertransaction"

	"github.com/gin-gonic/gin"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/library/http/response"
	"case-study-kredit-plus/library/types"

	consumertransactionRepository "case-study-kredit-plus/src/services/consumertransaction/repository"
	consumertransactionUsecase "case-study-kredit-plus/src/services/consumertransaction/usecase"

	consumercreditlimitRepository "case-study-kredit-plus/src/services/consumercreditlimit/repository"
	consumercreditlimitUsecase "case-study-kredit-plus/src/services/consumercreditlimit/usecase"
)

var ()

type ConsumerTransactionHandler struct {
	ConsumerTransactionUsecase consumertransaction.Usecase
	dataManager                *data.Manager
	Result                     gin.H
	Status                     int
}

func (h ConsumerTransactionHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	consumertransactionRepo := consumertransactionRepository.NewConsumerTransactionRepository(
		data.NewMySQLStorage(db, "consumer_transactions", models.ConsumerTransaction{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	consumercreditlimitRepo := consumercreditlimitRepository.NewConsumerCreditLimitRepository(
		data.NewMySQLStorage(db, "consumer_credit_limits", models.ConsumerCreditLimit{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uConsumerCreditLimit := consumercreditlimitUsecase.NewConsumerCreditLimitUsecase(db, &consumercreditlimitRepo)

	uConsumerTransaction := consumertransactionUsecase.NewConsumerTransactionUsecase(db, &consumertransactionRepo, uConsumerCreditLimit)

	base := &ConsumerTransactionHandler{ConsumerTransactionUsecase: uConsumerTransaction, dataManager: dataManager}

	rs := v.Group("/consumers/transactions")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/status", middleware.Auth, base.UpdateStatus)
	}

	status := v.Group("/statuses")
	{
		status.GET("/consumers/transactions", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *ConsumerTransactionHandler) FindAll(c *gin.Context) {
	if c.Query("ConsumerID") != "" && !library.ValidateUUID(c.Query("ConsumerID")) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->FindAll()",
			Message:    "Consumer ID is not valid",
			Error:      fmt.Errorf("Consumer ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.Query("ContractNumber") != "" && !library.ValidateTextInput(c.Query("ContractNumber")) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->FindAll()",
			Message:    "Contract Number is not valid",
			Error:      fmt.Errorf("Contract Number is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	var params models.FindAllConsumerTransactionParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.ConsumerID = c.Query("ConsumerID")
	params.ContractNumber = c.Query("ContractNumber")
	params.LoanTerm, _ = strconv.Atoi(c.Query("LoanTerm"))
	datas, err := h.ConsumerTransactionUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	length, err := h.ConsumerTransactionUsecase.Count(c, params)
	if err != nil {
		err.Path = ".ConsumerTransactionHandler->FindAll()" + err.Path
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

func (h *ConsumerTransactionHandler) Find(c *gin.Context) {
	id := c.Param("id")

	if id != "" && !library.ValidateUUID(id) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Find()",
			Message:    "ID is not valid",
			Error:      fmt.Errorf("ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.ConsumerTransactionUsecase.Find(c, id)
	if err != nil {
		err.Path = ".ConsumerTransactionHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "ConsumerTransaction not found", http.StatusUnprocessableEntity, *err)
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

func (h *ConsumerTransactionHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.ConsumerTransaction
	var data *models.ConsumerTransaction

	if c.PostForm("ConsumerID") != "" && !library.ValidateUUID(c.PostForm("ConsumerID")) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Create()",
			Message:    "Consumer ID is not valid",
			Error:      fmt.Errorf("Consumer ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("ContractNumber") != "" && !library.ValidateTextInput(c.PostForm("ContractNumber")) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Create()",
			Message:    "Contract Number is not valid",
			Error:      fmt.Errorf("Contract Number is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("AssetName") != "" && !library.ValidateTextInput(c.PostForm("AssetName")) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Create()",
			Message:    "Asset Name is not valid",
			Error:      fmt.Errorf("Asset Name is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	otr, errParseFloat := strconv.ParseFloat(c.PostForm("OTR"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Create()",
			Message:    "OTR Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	adminFee, errParseFloat := strconv.ParseFloat(c.PostForm("AdminFee"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Create()",
			Message:    "Admin Fee Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("InstallmentAmount") != "" {
		installmentAmount, errParseFloat := strconv.ParseFloat(c.PostForm("InstallmentAmount"), 64)
		if errParseFloat != nil {
			err := &types.Error{
				Path:       ".ConsumerTransactionHandler->Create()",
				Message:    "Installment Amount Invalid",
				Error:      errParseFloat,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		obj.InstallmentAmount = installmentAmount
	}

	loanTerm, errParseInt := strconv.Atoi(c.PostForm("LoanTerm"))
	if errParseInt != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Create()",
			Message:    "Loan Term Invalid",
			Error:      errParseInt,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	interestAmount, errParseFloat := strconv.ParseFloat(c.PostForm("InterestAmount"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Create()",
			Message:    "Interest Amount Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("TotalAmount") != "" {
		totalAmount, errParseFloat := strconv.ParseFloat(c.PostForm("TotalAmount"), 64)
		if errParseFloat != nil {
			err := &types.Error{
				Path:       ".ConsumerTransactionHandler->Create()",
				Message:    "Total Amount Invalid",
				Error:      errParseFloat,
				Type:       "conversion-error",
				StatusCode: http.StatusUnprocessableEntity,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		obj.TotalAmount = totalAmount
	}

	obj.ConsumerID = c.PostForm("ConsumerID")
	obj.ContractNumber = c.PostForm("ContractNumber")
	obj.OTR = otr
	obj.AdminFee = adminFee
	obj.InterestAmount = interestAmount
	obj.LoanTerm = loanTerm
	obj.AssetName = c.PostForm("AssetName")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ConsumerTransactionUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".ConsumerTransactionHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data created successfuly", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ConsumerTransactionHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.ConsumerTransaction
	var data *models.ConsumerTransaction

	id := c.Param("id")

	if id != "" && !library.ValidateUUID(id) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "ID is not valid",
			Error:      fmt.Errorf("ID is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("AssetName") != "" && !library.ValidateTextInput(c.PostForm("AssetName")) {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "Asset Name is not valid",
			Error:      fmt.Errorf("Asset Name is not valid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	otr, errParseFloat := strconv.ParseFloat(c.PostForm("OTR"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "OTR Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	adminFee, errParseFloat := strconv.ParseFloat(c.PostForm("AdminFee"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "Admin Fee Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	installmentAmount, errParseFloat := strconv.ParseFloat(c.PostForm("InstallmentAmount"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "Installment Amount Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	loanTerm, errParseInt := strconv.Atoi(c.PostForm("LoanTerm"))
	if errParseInt != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "Loan Term Invalid",
			Error:      errParseInt,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	interestAmount, errParseFloat := strconv.ParseFloat(c.PostForm("InterestAmount"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "Interest Amount Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	totalAmount, errParseFloat := strconv.ParseFloat(c.PostForm("TotalAmount"), 64)
	if errParseFloat != nil {
		err := &types.Error{
			Path:       ".ConsumerTransactionHandler->Update()",
			Message:    "Total Amount Invalid",
			Error:      errParseFloat,
			Type:       "conversion-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	// obj.ConsumerID = c.PostForm("ConsumerID")
	// obj.ContractNumber = c.PostForm("ContractNumber")
	obj.OTR = otr
	obj.AdminFee = adminFee
	obj.InstallmentAmount = installmentAmount
	obj.LoanTerm = loanTerm
	obj.InterestAmount = interestAmount
	obj.TotalAmount = totalAmount
	obj.AssetName = c.PostForm("AssetName")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ConsumerTransactionUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ConsumerTransactionHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data successfuly updated", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ConsumerTransactionHandler) FindStatus(c *gin.Context) {
	datas, err := h.ConsumerTransactionUsecase.FindStatus(c)
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

func (h *ConsumerTransactionHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.ConsumerTransaction

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".ConsumerTransactionHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.ConsumerTransactionUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ConsumerTransactionHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Status update success", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
