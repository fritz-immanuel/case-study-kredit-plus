package usecase

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"case-study-kredit-plus/library"
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/src/services/consumercreditlimit"
	"case-study-kredit-plus/src/services/consumertransaction"

	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type ConsumerTransactionUsecase struct {
	consumertransactionRepo    consumertransaction.Repository
	consumercreditlimitUsecase consumercreditlimit.Usecase
	contextTimeout             time.Duration
	db                         *sqlx.DB
}

func NewConsumerTransactionUsecase(db *sqlx.DB, consumertransactionRepo consumertransaction.Repository, consumercreditlimitUsecase consumercreditlimit.Usecase) consumertransaction.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &ConsumerTransactionUsecase{
		consumertransactionRepo:    consumertransactionRepo,
		consumercreditlimitUsecase: consumercreditlimitUsecase,
		contextTimeout:             timeoutContext,
		db:                         db,
	}
}

func (u *ConsumerTransactionUsecase) FindAll(ctx *gin.Context, params models.FindAllConsumerTransactionParams) ([]*models.ConsumerTransaction, *types.Error) {
	result, err := u.consumertransactionRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerTransactionUsecase) Find(ctx *gin.Context, id string) (*models.ConsumerTransaction, *types.Error) {
	result, err := u.consumertransactionRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerTransactionUsecase) Count(ctx *gin.Context, params models.FindAllConsumerTransactionParams) (int, *types.Error) {
	result, err := u.consumertransactionRepo.Count(ctx, params)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->Count()" + err.Path
		return 0, err
	}

	return result, nil
}

func (u *ConsumerTransactionUsecase) Create(ctx *gin.Context, obj models.ConsumerTransaction) (*models.ConsumerTransaction, *types.Error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(obj)
	if errValidation != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	if !library.IsValidTenor(obj.LoanTerm) {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionUsecase->Create()",
			Message:    "Loan Term Invalid",
			Error:      fmt.Errorf("Loan Term Invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	totalAmount, installmentAmount := obj.TotalAmount, obj.InstallmentAmount

	if obj.TotalAmount == 0 {
		totalAmount = obj.OTR + obj.AdminFee + obj.InterestAmount
	}

	if obj.InstallmentAmount <= 0 {
		installmentAmount = totalAmount / float64(obj.LoanTerm)
	}

	data := models.ConsumerTransaction{
		ID:                uuid.New().String(),
		ConsumerID:        obj.ConsumerID,
		ContractNumber:    obj.ContractNumber,
		OTR:               obj.OTR,
		AdminFee:          obj.AdminFee,
		InstallmentAmount: installmentAmount,
		LoanTerm:          obj.LoanTerm,
		InterestAmount:    obj.InterestAmount,
		TotalAmount:       totalAmount,
		AssetName:         obj.AssetName,
		StatusID:          models.DEFAULT_STATUS_ID,
	}

	// check tenor limit availability
	remainingLimit, err := u.consumercreditlimitUsecase.CheckCreditLimitAvailability(ctx, obj.ConsumerID, obj.LoanTerm)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->Create()" + err.Path
		return nil, err
	}

	if remainingLimit < totalAmount {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionUsecase->Create()",
			Message:    "Insufficient Credit Limit",
			Error:      fmt.Errorf("Insufficient Credit Limit"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	result, err := u.consumertransactionRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerTransactionUsecase) Update(ctx *gin.Context, id string, obj models.ConsumerTransaction) (*models.ConsumerTransaction, *types.Error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(obj)
	if errValidation != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	if !library.IsValidTenor(obj.LoanTerm) {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionUsecase->Update()",
			Message:    "Loan Term Invalid",
			Error:      fmt.Errorf("Loan Term Invalid"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	data, err := u.consumertransactionRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->Update()" + err.Path
		return nil, err
	}

	totalAmount, installmentAmount := obj.TotalAmount, obj.InstallmentAmount

	if obj.TotalAmount == 0 {
		totalAmount = obj.OTR + obj.AdminFee + obj.InterestAmount
	}

	if obj.InstallmentAmount <= 0 {
		installmentAmount = totalAmount / float64(obj.LoanTerm)
	}

	// data.ConsumerID = obj.ConsumerID
	// data.ContractNumber = obj.ContractNumber
	data.OTR = obj.OTR
	data.AdminFee = obj.AdminFee
	data.InstallmentAmount = installmentAmount
	data.LoanTerm = obj.LoanTerm
	data.InterestAmount = obj.InterestAmount
	data.TotalAmount = totalAmount
	data.AssetName = obj.AssetName

	// check tenor limit availability
	remainingLimit, err := u.consumercreditlimitUsecase.CheckCreditLimitAvailability(ctx, obj.ConsumerID, obj.LoanTerm)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->Update()" + err.Path
		return nil, err
	}

	if remainingLimit < totalAmount {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionUsecase->Update()",
			Message:    "Insufficient Credit Limit",
			Error:      fmt.Errorf("Insufficient Credit Limit"),
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	result, err := u.consumertransactionRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *ConsumerTransactionUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.consumertransactionRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerTransactionUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.ConsumerTransaction, *types.Error) {
	result, err := u.consumertransactionRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".ConsumerTransactionUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}
