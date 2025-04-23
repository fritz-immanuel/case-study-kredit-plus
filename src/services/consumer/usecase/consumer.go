package usecase

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/src/services/consumer"

	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type ConsumerUsecase struct {
	consumerRepo   consumer.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewConsumerUsecase(db *sqlx.DB, consumerRepo consumer.Repository) consumer.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &ConsumerUsecase{
		consumerRepo:   consumerRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *ConsumerUsecase) FindAll(ctx *gin.Context, params models.FindAllConsumerParams) ([]*models.Consumer, *types.Error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(params)
	if errValidation != nil {
		return nil, &types.Error{
			Path:       ".ConsumerUsecase->FindAll()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	result, err := u.consumerRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".ConsumerUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerUsecase) Find(ctx *gin.Context, id string) (*models.Consumer, *types.Error) {
	result, err := u.consumerRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ConsumerUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerUsecase) Count(ctx *gin.Context, params models.FindAllConsumerParams) (int, *types.Error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(params)
	if errValidation != nil {
		return 0, &types.Error{
			Path:       ".ConsumerUsecase->Count()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}
	
	result, err := u.consumerRepo.Count(ctx, params)
	if err != nil {
		err.Path = ".ConsumerUsecase->Count()" + err.Path
		return 0, err
	}

	return result, nil
}

func (u *ConsumerUsecase) Create(ctx *gin.Context, obj models.Consumer) (*models.Consumer, *types.Error) {
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
			Path:       ".ConsumerUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	// check duplicate NIK
	var dupeParams models.FindAllConsumerParams
	dupeParams.NIK = obj.NIK
	dupeParams.FindAllParams.StatusID = `status_id = "1"`
	count, err := u.consumerRepo.Count(ctx, dupeParams)
	if err != nil {
		err.Path = ".ConsumerUsecase->Create()" + err.Path
		return nil, err
	}

	if count > 0 {
		return nil, &types.Error{
			Path:       ".ConsumerUsecase->Create()",
			Message:    "NIK already exists",
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data := models.Consumer{
		ID:           uuid.New().String(),
		NIK:          obj.NIK,
		FullName:     obj.FullName,
		LegalName:    obj.LegalName,
		PlaceOfBirth: obj.PlaceOfBirth,
		DateOfBirth:  obj.DateOfBirth,
		Salary:       obj.Salary,
		KTPImgURL:    obj.KTPImgURL,
		SelfieImgURL: obj.SelfieImgURL,
		StatusID:     models.DEFAULT_STATUS_ID,
	}

	if data.Salary < 0 && data.Salary > 99999999999 {
		return nil, &types.Error{
			Path:       ".ConsumerUsecase->Create()",
			Message:    "Salary must be between 0 and 99 Billion",
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	result, err := u.consumerRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".ConsumerUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerUsecase) Update(ctx *gin.Context, id string, obj models.Consumer) (*models.Consumer, *types.Error) {
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
			Path:       ".ConsumerUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	// check duplicate NIK
	var dupeParams models.FindAllConsumerParams
	dupeParams.NIK = obj.NIK
	dupeParams.FindAllParams.StatusID = `status_id = "1"`
	dupeParams.FindAllParams.DataFinder = fmt.Sprintf(`consumers.id != "%s"`, id)
	count, err := u.consumerRepo.Count(ctx, dupeParams)
	if err != nil {
		err.Path = ".ConsumerUsecase->Update()" + err.Path
		return nil, err
	}

	if count > 0 {
		return nil, &types.Error{
			Path:       ".ConsumerUsecase->Update()",
			Message:    "NIK already exists",
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.consumerRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ConsumerUsecase->Update()" + err.Path
		return nil, err
	}

	data.NIK = obj.NIK
	data.FullName = obj.FullName
	data.LegalName = obj.LegalName
	data.PlaceOfBirth = obj.PlaceOfBirth
	data.DateOfBirth = obj.DateOfBirth
	data.Salary = obj.Salary
	data.KTPImgURL = obj.KTPImgURL
	data.SelfieImgURL = obj.SelfieImgURL

	result, err := u.consumerRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".ConsumerUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *ConsumerUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.consumerRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".ConsumerUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Consumer, *types.Error) {
	result, err := u.consumerRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".ConsumerUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}
