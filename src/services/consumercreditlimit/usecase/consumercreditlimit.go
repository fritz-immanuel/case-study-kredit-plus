package usecase

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/src/services/consumercreditlimit"

	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type ConsumerCreditLimitUsecase struct {
	consumercreditlimitRepo consumercreditlimit.Repository
	contextTimeout          time.Duration
	db                      *sqlx.DB
}

func NewConsumerCreditLimitUsecase(db *sqlx.DB, consumercreditlimitRepo consumercreditlimit.Repository) consumercreditlimit.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &ConsumerCreditLimitUsecase{
		consumercreditlimitRepo: consumercreditlimitRepo,
		contextTimeout:          timeoutContext,
		db:                      db,
	}
}

func (u *ConsumerCreditLimitUsecase) FindAll(ctx *gin.Context, params models.FindAllConsumerCreditLimitParams) ([]*models.ConsumerCreditLimit, *types.Error) {
	result, err := u.consumercreditlimitRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerCreditLimitUsecase) Find(ctx *gin.Context, id string) (*models.ConsumerCreditLimit, *types.Error) {
	result, err := u.consumercreditlimitRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerCreditLimitUsecase) Count(ctx *gin.Context, params models.FindAllConsumerCreditLimitParams) (int, *types.Error) {
	result, err := u.consumercreditlimitRepo.Count(ctx, params)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->Count()" + err.Path
		return 0, err
	}

	return result, nil
}

func (u *ConsumerCreditLimitUsecase) Create(ctx *gin.Context, obj models.ConsumerCreditLimit) (*models.ConsumerCreditLimit, *types.Error) {
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
			Path:       ".ConsumerCreditLimitUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	// check duplicate ConsumerID
	var dupeParams models.FindAllConsumerCreditLimitParams
	dupeParams.ConsumerID = obj.ConsumerID
	dupeParams.FindAllParams.StatusID = `status_id = "1"`
	dupeData, err := u.consumercreditlimitRepo.FindAll(ctx, dupeParams)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->Create()" + err.Path
		return nil, err
	}

	if len(dupeData) > 0 {
		for _, dupe := range dupeData {
			_, err := u.consumercreditlimitRepo.UpdateStatus(ctx, dupe.ID, models.STATUS_INACTIVE)
			if err != nil {
				err.Path = ".ConsumerCreditLimitUsecase->Create()" + err.Path
				return nil, err
			}
		}
	}

	data := models.ConsumerCreditLimit{
		ID:         uuid.New().String(),
		ConsumerID: obj.ConsumerID,
		Month1:     obj.Month1,
		Month2:     obj.Month2,
		Month3:     obj.Month3,
		Month6:     obj.Month6,
		StatusID:   models.DEFAULT_STATUS_ID,
	}

	result, err := u.consumercreditlimitRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerCreditLimitUsecase) Update(ctx *gin.Context, id string, obj models.ConsumerCreditLimit) (*models.ConsumerCreditLimit, *types.Error) {
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
			Path:       ".ConsumerCreditLimitUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.consumercreditlimitRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->Update()" + err.Path
		return nil, err
	}

	// data.ConsumerID = obj.ConsumerID
	data.Month1 = obj.Month1
	data.Month2 = obj.Month2
	data.Month3 = obj.Month3
	data.Month6 = obj.Month6

	result, err := u.consumercreditlimitRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *ConsumerCreditLimitUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.consumercreditlimitRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ConsumerCreditLimitUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.ConsumerCreditLimit, *types.Error) {
	result, err := u.consumercreditlimitRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}

// CHECK CONSUMER CREDIT LIMIT FOR TENOR
func (u *ConsumerCreditLimitUsecase) CheckCreditLimitAvailability(ctx *gin.Context, consumerID string, tenor int) (float64, *types.Error) {
	result, err := u.consumercreditlimitRepo.CheckCreditLimitAvailability(ctx, consumerID, tenor)
	if err != nil {
		err.Path = ".ConsumerCreditLimitUsecase->CheckCreditLimitAvailability()" + err.Path
		return 0, err
	}

	return result, nil
}
