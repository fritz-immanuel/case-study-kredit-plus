package usecase

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"case-study-kredit-plus/library"
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/src/services/user"

	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type UserUsecase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewUserUsecase(db *sqlx.DB, userRepo user.Repository) user.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &UserUsecase{
		userRepo:       userRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *UserUsecase) FindAll(ctx *gin.Context, params models.FindAllUserParams) ([]*models.User, *types.Error) {
	result, err := u.userRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) Find(ctx *gin.Context, id string) (*models.User, *types.Error) {
	result, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) Count(ctx *gin.Context, params models.FindAllUserParams) (int, *types.Error) {
	result, err := u.userRepo.Count(ctx, params)
	if err != nil {
		err.Path = ".UserUsecase->Count()" + err.Path
		return 0, err
	}

	return result, nil
}

func (u *UserUsecase) Create(ctx *gin.Context, obj models.User) (*models.User, *types.Error) {
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
			Path:       ".UserUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	// check duplicate email
	var dupeParams models.FindAllUserParams
	dupeParams.Email = obj.Email
	dupeParams.FindAllParams.StatusID = `status_id = "1"`
	count, err := u.userRepo.Count(ctx, dupeParams)
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	if count > 0 {
		return nil, &types.Error{
			Path:       ".UserUsecase->Create()",
			Message:    "Email already exists",
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data := models.User{
		ID:                 uuid.New().String(),
		Name:               obj.Name,
		Email:              obj.Email,
		Username:           obj.Username,
		CountryCallingCode: obj.CountryCallingCode,
		PhoneNumber:        obj.PhoneNumber,
		Password:           obj.Password,
		StatusID:           models.DEFAULT_STATUS_ID,
	}

	result, err := u.userRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) Update(ctx *gin.Context, id string, obj models.User) (*models.User, *types.Error) {
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
			Path:       ".UserUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	// check duplicate email
	var dupeParams models.FindAllUserParams
	dupeParams.Email = obj.Email
	dupeParams.FindAllParams.StatusID = `status_id = "1"`
	dupeParams.FindAllParams.DataFinder = fmt.Sprintf(`users.id != "%s"`, id)
	count, err := u.userRepo.Count(ctx, dupeParams)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	if count > 0 {
		return nil, &types.Error{
			Path:       ".UserUsecase->Update()",
			Message:    "Email already exists",
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Email = obj.Email
	data.Username = obj.Username
	data.CountryCallingCode = obj.CountryCallingCode
	data.PhoneNumber = obj.PhoneNumber

	result, err := u.userRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.userRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".UserUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.User, *types.Error) {
	result, err := u.userRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".UserUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}

// LOGIN

func (u *UserUsecase) Login(ctx *gin.Context, params models.FindAllUserParams) (*models.UserJWTContent, *types.Error) {
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
			Path:       ".UserService->Login()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	result, err := u.userRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserService->Login()" + err.Path
		return nil, err
	}

	if len(result) < 1 {
		var err types.Error
		err.Message = "username atau password salah"
		err.Type = "authentication"
		err.Error = fmt.Errorf("Login Failed")
		err.StatusCode = http.StatusUnprocessableEntity
		return nil, &err
	}

	credentials := library.Credential{ID: result[0].ID, Email: result[0].Email, Type: "Web"}

	token, errorJwtSign := library.JwtSignString(credentials)
	if errorJwtSign != nil {
		return nil, &types.Error{
			Error:      errorJwtSign,
			Message:    "Error JWT Sign String",
			Path:       ".UserService->Login()",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var userLogin models.UserJWTContent
	userLogin.ID = result[0].ID
	userLogin.Name = result[0].Name
	userLogin.Token = token
	userLogin.Email = result[0].Email
	userLogin.StatusID = result[0].StatusID

	return &userLogin, nil
}

// //

// UpdatePassword()  Updates the password of the user
func (u *UserUsecase) UpdatePassword(ctx *gin.Context, obj models.UserUpdatePassword) (*models.User, *types.Error) {
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
			Path:       ".UserUsecase->UpdatePassword()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.userRepo.Find(ctx, obj.ID)
	if err != nil {
		err.Path = ".UserUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	if data.Password != obj.OldPassword {
		err = &types.Error{
			Path:       ".UserUsecase->UpdatePassword()",
			Message:    "The erstwhile password fails to harmonize with the current, necessitating adjustment",
			Error:      nil,
			Type:       "validation-error",
			StatusCode: http.StatusUnprocessableEntity,
		}
		return nil, err
	}

	data.Password = obj.NewPassword

	result, err := u.userRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".UserUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	return result, err
}
