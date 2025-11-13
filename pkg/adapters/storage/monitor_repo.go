package storage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"health-check/internal/monitor_service/domain"
	"health-check/internal/monitor_service/port"
	"health-check/pkg/adapters/mapper"
	"health-check/pkg/adapters/types"
)

type monitorRepo struct {
	db *gorm.DB
}

func NewDomainRepo(db *gorm.DB) port.Repo {
	return &monitorRepo{db: db}
}

func (m *monitorRepo) Create(ctx context.Context, api domain.MonitoredAPI) (domain.ApiID, error) {

	monitor := mapper.MonitorDomain2Storage(api)
	return domain.ApiID(monitor.ID), m.db.Table("monitored_apis").WithContext(ctx).Create(monitor).Error
}

func (m *monitorRepo) GetByID(ctx context.Context, id domain.ApiID) (*domain.MonitoredAPI, error) {
	var apiDB types.MonitoredAPIDB

	err := m.db.WithContext(ctx).
		Table("monitored_apis").
		Where("id = ?", id).
		First(&apiDB).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	apiDomain, err := mapper.MonitorStorage2Domain(apiDB)
	if err != nil {
		return nil, err
	}

	return apiDomain, nil
}

func (m *monitorRepo) SaveCheckResult(ctx context.Context, result domain.CheckResult) error {
	row := mapper.ResultDomain2Storage(result)
	return m.db.WithContext(ctx).Table("check_result_db").Create(&row).Error
}
