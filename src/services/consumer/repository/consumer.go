package repository

import (
	"fmt"
	"net/http"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
)

type ConsumerRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewConsumerRepository(repository data.GenericStorage, statusRepository data.GenericStorage) ConsumerRepository {
	return ConsumerRepository{repository: repository, statusRepository: statusRepository}
}

func (s ConsumerRepository) FindAll(ctx *gin.Context, params models.FindAllConsumerParams) ([]*models.Consumer, *types.Error) {
	data := []*models.Consumer{}
	bulks := []*models.ConsumerBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(` AND %s`, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND consumers.%s`, params.FindAllParams.StatusID)
	}

	if params.NIK != "" {
		where += fmt.Sprintf(` AND consumers.NIK LIKE "%%%s%%"`, params.NIK)
	}

	if params.FullName != "" {
		where += fmt.Sprintf(` AND consumers.full_name LIKE "%%%s%%"`, params.FullName)
	}

	if params.LegalName != "" {
		where += fmt.Sprintf(` AND consumers.legal_name LIKE "%%%s%%"`, params.LegalName)
	}

	if params.PlaceOfBirth != "" {
		where += fmt.Sprintf(` AND consumers.place_of_birth LIKE "%%%s%%"`, params.PlaceOfBirth)
	}

	if params.MinDateOfBirth != "" {
		where += fmt.Sprintf(` AND consumers.date_of_birth >= "%s"`, params.MinDateOfBirth)
	}

	if params.MaxDateOfBirth != "" {
		where += fmt.Sprintf(` AND consumers.date_of_birth <= "%s"`, params.MaxDateOfBirth)
	}

	if params.MinSalary > 0 {
		where += fmt.Sprintf(` AND consumers.salary >= %f`, params.MinSalary)
	}

	if params.MaxSalary > 0 {
		where += fmt.Sprintf(` AND consumers.salary <= %f`, params.MaxSalary)
	}

	if params.FindAllParams.SortBy != "" {
		where += fmt.Sprintf(` ORDER BY %s`, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where += ` LIMIT :limit OFFSET :offset`
	}

	query := fmt.Sprintf(`
  SELECT
    consumers.id, consumers.NIK, consumers.full_name, consumers.legal_name, consumers.place_of_birth, consumers.date_of_birth,
    consumers.salary, consumers.ktp_img_url, consumers.selfie_img_url,
    consumers.status_id, status.name status_name
  FROM consumers
  JOIN status ON consumers.status_id = status.id
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":          params.FindAllParams.Size,
		"offset":         ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id":      params.FindAllParams.StatusID,
		"place_of_birth": params.PlaceOfBirth,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.Consumer{
			ID:           v.ID,
			NIK:          v.NIK,
			FullName:     v.FullName,
			LegalName:    v.LegalName,
			PlaceOfBirth: v.PlaceOfBirth,
			DateOfBirth:  v.DateOfBirth,
			Salary:       v.Salary,
			KTPImgURL:    v.KTPImgURL,
			SelfieImgURL: v.SelfieImgURL,
			StatusID:     v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		data = append(data, obj)
	}

	return data, nil
}

func (s ConsumerRepository) Find(ctx *gin.Context, id string) (*models.Consumer, *types.Error) {
	result := models.Consumer{}
	bulks := []*models.ConsumerBulk{}
	var err error

	query := `
  SELECT
    consumers.id, consumers.NIK, consumers.full_name, consumers.legal_name, consumers.place_of_birth, consumers.date_of_birth,
    consumers.salary, consumers.ktp_img_url, consumers.selfie_img_url,
    consumers.status_id, status.name status_name
  FROM consumers
  JOIN status ON consumers.status_id = status.id
  WHERE consumers.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Consumer{
			ID:           v.ID,
			NIK:          v.NIK,
			FullName:     v.FullName,
			LegalName:    v.LegalName,
			PlaceOfBirth: v.PlaceOfBirth,
			DateOfBirth:  v.DateOfBirth,
			Salary:       v.Salary,
			KTPImgURL:    v.KTPImgURL,
			SelfieImgURL: v.SelfieImgURL,
			StatusID:     v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s ConsumerRepository) Count(ctx *gin.Context, params models.FindAllConsumerParams) (int, *types.Error) {
	bulks := []*models.ConsumerBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(` AND %s`, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND consumers.%s`, params.FindAllParams.StatusID)
	}

	if params.FullName != "" {
		where += fmt.Sprintf(` AND consumers.full_name LIKE "%s%%"`, params.FullName)
	}

	if params.LegalName != "" {
		where += fmt.Sprintf(` AND consumers.legal_name LIKE "%s%%"`, params.LegalName)
	}

	if params.PlaceOfBirth != "" {
		where += ` AND consumers.place_of_birth = :place_of_birth`
	}

	if params.MinDateOfBirth != "" {
		where += fmt.Sprintf(` AND consumers.date_of_birth >= "%s"`, params.MinDateOfBirth)
	}

	if params.MaxDateOfBirth != "" {
		where += fmt.Sprintf(` AND consumers.date_of_birth <= "%s"`, params.MaxDateOfBirth)
	}

	if params.MinSalary >= -1 {
		where += fmt.Sprintf(` AND consumers.salary >= "%s"`, params.MinSalary)
	}

	if params.MaxSalary >= -1 {
		where += fmt.Sprintf(` AND consumers.salary <= "%s"`, params.MaxSalary)
	}

	query := fmt.Sprintf(`
  SELECT
    consumers.id, consumers.NIK, consumers.full_name, consumers.legal_name, consumers.place_of_birth, consumers.date_of_birth,
    consumers.salary, consumers.ktp_img_url, consumers.selfie_img_url,
    consumers.status_id, status.name status_name
  FROM consumers
  JOIN status ON consumers.status_id = status.id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":          params.FindAllParams.Size,
		"offset":         ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id":      params.FindAllParams.StatusID,
		"place_of_birth": params.PlaceOfBirth,
	})
	if err != nil {
		return 0, &types.Error{
			Path:       ".ConsumerStorage->Count()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return len(bulks), nil
}

func (s ConsumerRepository) Create(ctx *gin.Context, obj *models.Consumer) (*models.Consumer, *types.Error) {
	data := models.Consumer{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ConsumerRepository) Update(ctx *gin.Context, obj *models.Consumer) (*models.Consumer, *types.Error) {
	data := models.Consumer{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ConsumerRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	status := []*models.Status{}

	err := s.statusRepository.Where(ctx, &status, "1=1", map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return status, nil
}

func (s ConsumerRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Consumer, *types.Error) {
	data := models.Consumer{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
