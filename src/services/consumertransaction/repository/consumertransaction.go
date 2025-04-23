package repository

import (
	"fmt"
	"net/http"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
)

type ConsumerTransactionRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewConsumerTransactionRepository(repository data.GenericStorage, statusRepository data.GenericStorage) ConsumerTransactionRepository {
	return ConsumerTransactionRepository{repository: repository, statusRepository: statusRepository}
}

func (s ConsumerTransactionRepository) FindAll(ctx *gin.Context, params models.FindAllConsumerTransactionParams) ([]*models.ConsumerTransaction, *types.Error) {
	data := []*models.ConsumerTransaction{}
	bulks := []*models.ConsumerTransactionBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(` AND %s`, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND consumer_transactions.%s`, params.FindAllParams.StatusID)
	}

	if params.ConsumerID != "" {
		where += ` AND consumer_transactions.consumer_id = :consumer_id`
	}

	if params.ContractNumber != "" {
		where += fmt.Sprintf(` AND consumers.contract_number LIKE "%%%s%%"`, params.ContractNumber)
	}

	if params.LoanTerm != 0 {
		where += fmt.Sprintf(` AND consumers.contract_number = %d`, params.LoanTerm)
	}

	if params.FindAllParams.SortBy != "" {
		where += fmt.Sprintf(` ORDER BY %s`, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where += ` LIMIT :limit OFFSET :offset`
	}

	query := fmt.Sprintf(`
  SELECT
    consumer_transactions.id, consumer_transactions.consumer_id, consumer_transactions.contract_number, consumer_transactions.OTR,
    consumer_transactions.admin_fee, consumer_transactions.installment_amount, consumer_transactions.loan_term, consumer_transactions.interest_amount,
    consumer_transactions.asset_name, consumer_transactions.total_amount,
    consumer_transactions.status_id, status.name status_name, consumers.full_name consumer_name
  FROM consumer_transactions
  JOIN status ON consumer_transactions.status_id = status.id
  JOIN consumers ON consumers.id = consumer_transactions.consumer_id
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
			Path:       ".ConsumerTransactionStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.ConsumerTransaction{
			ID:         v.ID,
			ConsumerID: v.ConsumerID,
			Consumer: &models.IDNameTemplate{
				ID:   v.ConsumerID,
				Name: v.ConsumerName,
			},
			ContractNumber:    v.ContractNumber,
			OTR:               v.OTR,
			AdminFee:          v.AdminFee,
			InstallmentAmount: v.InstallmentAmount,
			LoanTerm:          v.LoanTerm,
			InterestAmount:    v.InterestAmount,
			TotalAmount:       v.TotalAmount,
			AssetName:         v.AssetName,
			StatusID:          v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		data = append(data, obj)
	}

	return data, nil
}

func (s ConsumerTransactionRepository) Find(ctx *gin.Context, id string) (*models.ConsumerTransaction, *types.Error) {
	result := models.ConsumerTransaction{}
	bulks := []*models.ConsumerTransactionBulk{}
	var err error

	query := `
  SELECT
    consumer_transactions.id, consumer_transactions.consumer_id, consumer_transactions.contract_number, consumer_transactions.OTR,
    consumer_transactions.admin_fee, consumer_transactions.installment_amount, consumer_transactions.loan_term, consumer_transactions.interest_amount,
    consumer_transactions.asset_name, consumer_transactions.total_amount,
    consumer_transactions.status_id, status.name status_name, consumers.full_name consumer_name
  FROM consumer_transactions
  JOIN status ON consumer_transactions.status_id = status.id
  JOIN consumers ON consumers.id = consumer_transactions.consumer_id
  WHERE consumer_transactions.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.ConsumerTransaction{
			ID:         v.ID,
			ConsumerID: v.ConsumerID,
			Consumer: &models.IDNameTemplate{
				ID:   v.ConsumerID,
				Name: v.ConsumerName,
			},
			ContractNumber:    v.ContractNumber,
			OTR:               v.OTR,
			AdminFee:          v.AdminFee,
			InstallmentAmount: v.InstallmentAmount,
			LoanTerm:          v.LoanTerm,
			InterestAmount:    v.InterestAmount,
			TotalAmount:       v.TotalAmount,
			AssetName:         v.AssetName,
			StatusID:          v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s ConsumerTransactionRepository) Count(ctx *gin.Context, params models.FindAllConsumerTransactionParams) (int, *types.Error) {
	bulks := []*models.ConsumerTransactionBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(` AND %s`, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND consumer_transactions.%s`, params.FindAllParams.StatusID)
	}

	if params.ContractNumber != "" {
		where += fmt.Sprintf(` AND consumers.contract_number LIKE "%%%s%%"`, params.ContractNumber)
	}

	if params.ConsumerID != "" {
		where += ` AND consumer_transactions.consumer_id = :consumer_id`
	}

	query := fmt.Sprintf(`
  SELECT
    consumer_transactions.id, consumer_transactions.consumer_id, consumer_transactions.contract_number, consumer_transactions.OTR,
    consumer_transactions.admin_fee, consumer_transactions.installment_amount, consumer_transactions.loan_term, consumer_transactions.interest_amount,
    consumer_transactions.asset_name, consumer_transactions.total_amount,
    consumer_transactions.status_id, status.name status_name, consumers.full_name consumer_name
  FROM consumer_transactions
  JOIN status ON consumer_transactions.status_id = status.id
  JOIN consumers ON consumers.id = consumer_transactions.consumer_id
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
			Path:       ".ConsumerTransactionStorage->Count()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return len(bulks), nil
}

func (s ConsumerTransactionRepository) Create(ctx *gin.Context, obj *models.ConsumerTransaction) (*models.ConsumerTransaction, *types.Error) {
	data := models.ConsumerTransaction{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ConsumerTransactionRepository) Update(ctx *gin.Context, obj *models.ConsumerTransaction) (*models.ConsumerTransaction, *types.Error) {
	data := models.ConsumerTransaction{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ConsumerTransactionRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	status := []*models.Status{}

	err := s.statusRepository.Where(ctx, &status, "1=1", map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return status, nil
}

func (s ConsumerTransactionRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.ConsumerTransaction, *types.Error) {
	data := models.ConsumerTransaction{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ConsumerTransactionStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
