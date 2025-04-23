package repository

import (
	"fmt"
	"net/http"
	"strings"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
)

type ConsumerCreditLimitRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewConsumerCreditLimitRepository(repository data.GenericStorage, statusRepository data.GenericStorage) ConsumerCreditLimitRepository {
	return ConsumerCreditLimitRepository{repository: repository, statusRepository: statusRepository}
}

func (s ConsumerCreditLimitRepository) FindAll(ctx *gin.Context, params models.FindAllConsumerCreditLimitParams) ([]*models.ConsumerCreditLimit, *types.Error) {
	data := []*models.ConsumerCreditLimit{}
	bulks := []*models.ConsumerCreditLimitBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(` AND %s`, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND consumer_credit_limits.%s`, params.FindAllParams.StatusID)
	}

	if params.ConsumerID != "" {
		where += ` AND consumer_credit_limits.consumer_id = :consumer_id`
	}

	if params.FindAllParams.SortBy != "" {
		where += fmt.Sprintf(` ORDER BY %s`, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where += ` LIMIT :limit OFFSET :offset`
	}

	query := fmt.Sprintf(`
  SELECT
    consumer_credit_limits.id, consumer_credit_limits.consumer_id,
    consumer_credit_limits.1_month, consumer_credit_limits.2_month, consumer_credit_limits.3_month, consumer_credit_limits.6_month,
    consumer_credit_limits.status_id, status.name status_name, consumers.full_name consumer_name
  FROM consumer_credit_limits
  JOIN status ON consumer_credit_limits.status_id = status.id
  JOIN consumers ON consumers.id = consumer_credit_limits.consumer_id
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":       params.FindAllParams.Size,
		"offset":      ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id":   params.FindAllParams.StatusID,
		"consumer_id": params.ConsumerID,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.ConsumerCreditLimit{
			ID:         v.ID,
			ConsumerID: v.ConsumerID,
			Consumer: &models.IDNameTemplate{
				ID:   v.ConsumerID,
				Name: v.ConsumerName,
			},
			Month1:   v.Month1,
			Month2:   v.Month2,
			Month3:   v.Month3,
			Month6:   v.Month6,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		data = append(data, obj)
	}

	return data, nil
}

func (s ConsumerCreditLimitRepository) Find(ctx *gin.Context, id string) (*models.ConsumerCreditLimit, *types.Error) {
	result := models.ConsumerCreditLimit{}
	bulks := []*models.ConsumerCreditLimitBulk{}
	var err error

	query := `
  SELECT
    consumer_credit_limits.id, consumer_credit_limits.consumer_id,
    consumer_credit_limits.1_month, consumer_credit_limits.2_month, consumer_credit_limits.3_month, consumer_credit_limits.6_month,
    consumer_credit_limits.status_id, status.name status_name, consumers.full_name consumer_name
  FROM consumer_credit_limits
  JOIN status ON consumer_credit_limits.status_id = status.id
  JOIN consumers ON consumers.id = consumer_credit_limits.consumer_id
  WHERE consumer_credit_limits.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.ConsumerCreditLimit{
			ID:         v.ID,
			ConsumerID: v.ConsumerID,
			Consumer: &models.IDNameTemplate{
				ID:   v.ConsumerID,
				Name: v.ConsumerName,
			},
			Month1:   v.Month1,
			Month2:   v.Month2,
			Month3:   v.Month3,
			Month6:   v.Month6,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s ConsumerCreditLimitRepository) Count(ctx *gin.Context, params models.FindAllConsumerCreditLimitParams) (int, *types.Error) {
	bulks := []*models.ConsumerCreditLimitBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(` AND %s`, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND consumer_credit_limits.%s`, params.FindAllParams.StatusID)
	}

	if params.ConsumerID != "" {
		where += ` AND consumer_credit_limits.consumer_id = :consumer_id`
	}

	query := fmt.Sprintf(`
  SELECT
    consumer_credit_limits.id, consumer_credit_limits.consumer_id,
    consumer_credit_limits.1_month, consumer_credit_limits.2_month, consumer_credit_limits.3_month, consumer_credit_limits.6_month,
    consumer_credit_limits.status_id, status.name status_name, consumers.full_name consumer_name
  FROM consumer_credit_limits
  JOIN status ON consumer_credit_limits.status_id = status.id
  JOIN consumers ON consumers.id = consumer_credit_limits.consumer_id
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":       params.FindAllParams.Size,
		"offset":      ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id":   params.FindAllParams.StatusID,
		"consumer_id": params.ConsumerID,
	})
	if err != nil {
		return 0, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->Count()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return len(bulks), nil
}

func (s ConsumerCreditLimitRepository) Create(ctx *gin.Context, obj *models.ConsumerCreditLimit) (*models.ConsumerCreditLimit, *types.Error) {
	data := models.ConsumerCreditLimit{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ConsumerCreditLimitRepository) Update(ctx *gin.Context, obj *models.ConsumerCreditLimit) (*models.ConsumerCreditLimit, *types.Error) {
	data := models.ConsumerCreditLimit{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ConsumerCreditLimitRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	status := []*models.Status{}

	err := s.statusRepository.Where(ctx, &status, "1=1", map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return status, nil
}

func (s ConsumerCreditLimitRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.ConsumerCreditLimit, *types.Error) {
	data := models.ConsumerCreditLimit{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}

// CHECK CONSUMER CREDIT LIMIT FOR TENOR
func (s ConsumerCreditLimitRepository) CheckCreditLimitAvailability(ctx *gin.Context, consumerID string, tenor int) (float64, *types.Error) {
	data := []*models.ConsumerCreditLimitAvailability{}

	var err error

	query := `
  SELECT
    cl.consumer_id,
    CASE
      WHEN <loanTerm> = 1 THEN cl.1_month - IFNULL(SUM(ct.total_amount), 0)
      WHEN <loanTerm> = 2 THEN cl.2_month - IFNULL(SUM(ct.total_amount), 0)
      WHEN <loanTerm> = 3 THEN cl.3_month - IFNULL(SUM(ct.total_amount), 0)
      WHEN <loanTerm> = 6 THEN cl.6_month - IFNULL(SUM(ct.total_amount), 0)
      ELSE NULL
    END remaining_limit
  FROM consumer_credit_limits cl
  LEFT JOIN consumer_transactions ct ON cl.consumer_id = ct.consumer_id AND ct.loan_term = <loanTerm>
  WHERE cl.status_id = 1 AND cl.consumer_id = '<consumerID>'
  GROUP BY cl.consumer_id`

	query = strings.ReplaceAll(query, `<loanTerm>`, fmt.Sprintf(`%d`, tenor))
	query = strings.ReplaceAll(query, `<consumerID>`, consumerID)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &data, query, map[string]interface{}{})
	if err != nil {
		return 0, &types.Error{
			Path:       ".ConsumerCreditLimitStorage->CheckCreditLimitAvailability()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(data) > 0 {
		return data[0].RemainingLimit, nil
	}

	return 0, nil
}
