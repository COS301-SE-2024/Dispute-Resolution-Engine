package mediatorassignment

import (
	"api/models"

	"gorm.io/gorm"
)

const (
	ExpertIDColumn            = "expert_id"
	ExpertNameColumn          = "expert_name"
	RejectionPercentageColumn = "rejection_percentage"
	LastAssignedDateColumn    = "last_assigned_date"
	AssignedDisputeCountColumn = "assigned_dispute_count"
)

type DBModel interface {
	GetExpertSummaryViews() ([]models.ExpertSummaryView, error)
	GetExpertSummaryViewByExpertID(expertID int) (models.ExpertSummaryView, error)
	GetExpertSummaryViewByColumn(columnName string) (models.ExpertSummaryView, error)
	GetExpertSummaryViewByColumnValue(columnName string, columnValue string) (models.ExpertSummaryView, error)
}

type DBModelReal struct {
	DB *gorm.DB
}

func (d *DBModelReal) GetExpertSummaryViews() ([]models.ExpertSummaryView, error) {
	var expertSummaryViews []models.ExpertSummaryView
	d.DB.Find(&expertSummaryViews)
	return expertSummaryViews, nil
}

func (d *DBModelReal) GetExpertSummaryViewByExpertID(expertID int) (models.ExpertSummaryView, error) {
	var expertSummaryView models.ExpertSummaryView
	d.DB.Where(ExpertIDColumn, expertID).First(&expertSummaryView)
	return expertSummaryView, nil
}

func (d *DBModelReal) GetExpertSummaryViewByColumn(columnName string) (models.ExpertSummaryView, error) {
	var expertSummaryView models.ExpertSummaryView
	d.DB.Where(columnName).First(&expertSummaryView)
	return expertSummaryView, nil
}

func (d *DBModelReal) GetExpertSummaryViewByColumnValue(columnName string, columnValue string) (models.ExpertSummaryView, error) {
	var expertSummaryView models.ExpertSummaryView
	d.DB.Where(columnName, columnValue).First(&expertSummaryView)
	return expertSummaryView, nil
}