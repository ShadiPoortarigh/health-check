package storage

import (
	"gorm.io/gorm"
	"health-check/internal/monitor_service/domain"
	"health-check/internal/monitor_service/port"
	"health-check/pkg/adapters/mapper"
)

type monitorRepo struct {
	db *gorm.DB
}

func NewDomainRepo(db *gorm.DB) port.Repo {
	return &monitorRepo{db: db}
}

func (m *monitorRepo) Create(api domain.MonitoredAPI) (domain.ApiID, error) {

	monitor := mapper.MonitorDomain2Storage(api)
	return domain.ApiID(monitor.ID), m.db.Table("monitored_apis").Create(monitor).Error
}
