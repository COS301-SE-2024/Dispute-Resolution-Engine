package archive

import (
	"api/env"
	"api/models"
	"api/utilities"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ArchiveModel interface {
	GetArchives(limit *int) (error, []models.ArchivedDisputeSummary)
}

type Archive struct {
	Model ArchiveModel
	Env   env.Env
}

type ArchiveModelReal struct {
	db  *gorm.DB
	env env.Env
}

func NewHandler(db *gorm.DB, envReader env.Env) Archive {
	return Archive{
		Model: &ArchiveModelReal{
			db:  db,
			env: envReader,
		},
		Env: envReader,
	}
}

func NewArchiveModelReal(db *gorm.DB, envReader env.Env) ArchiveModel {
	return &ArchiveModelReal{
		db:  db,
		env: envReader,
	}
}

func (a *ArchiveModelReal) GetArchives(limit *int) (error, []models.ArchivedDisputeSummary) {
	logger := utilities.NewLogger().LogWithCaller()

	var disputes []models.Dispute

	query := a.db.Model(&models.Dispute{}).
		Where("status = ?", models.StatusSettled).
		Or("status = ?", models.StatusRefused).
		Or("status = ?", models.StatusWithdrawn).
		Or("status = ?", models.StatusTransfer).
		Or("status = ?", models.StatusAppeal)

	if limit != nil && *limit > 0 {
		query = query.Limit(*limit)
	}

	if err := query.Scan(&disputes).Error; err != nil {
		logger.WithError(err).Error("Error retrieving disputes")
		return err, nil
	}

	// Transform the results to ArchivedDisputeSummary
	summaries := make([]models.ArchivedDisputeSummary, len(disputes))
	for i, dispute := range disputes {
		disputeSummary := models.DisputeSummaries{}
		err := a.db.Model(models.DisputeSummaries{}).Where("dispute = ?", *dispute.ID).First(&disputeSummary).Error
		if err != nil {
			logger.WithError(err).Error("Could not get dispute for id:" + fmt.Sprint(*dispute.ID))
			return err, nil
		}
		summaries[i] = models.ArchivedDisputeSummary{
			ID:           *dispute.ID,
			Title:        dispute.Title,
			Description:  dispute.Description,
			Summary:      disputeSummary.Summary,
			Category:     []string{"Dispute"}, // Assuming a default category for now
			DateFiled:    dispute.CaseDate.Format("2006-08-01"),
			DateResolved: dispute.CaseDate.Add(48 * time.Hour).Format("2006-08-01"), // Placeholder for resolved date
			Resolution:   string(dispute.Status),
		}
	}

	return nil, summaries
}
