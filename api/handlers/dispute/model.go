package dispute

import (
	"api/env"
	"api/handlers/notifications"
	"api/models"
	"api/utilities"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type DisputeModel interface {
	UploadEvidence(userId, disputeId int64, path string, file io.Reader) (uint, error)
	GetEvidenceByDispute(disputeId int64) ([]models.Evidence, error)
	GetDisputeExperts(disputeId int64) ([]models.Expert, error)

	GetDisputesByUser(userId int64) ([]models.Dispute, error)
	GetDispute(disputeId int64) (models.Dispute, error)

	GetUserByEmail(email string) (models.User, error)
	CreateDispute(dispute models.Dispute) (int64, error)
	UpdateDisputeStatus(disputeId int64, status string) error

	ObjectExpert(userId, disputeId, expertId int64, reason string) error
	ReviewExpertObjection(userId, disputeId, expertId int64, approved bool) error
}

type disputeModelReal struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Dispute {
	return Dispute{
		Email: notifications.NewHandler(db),
		Model: &disputeModelReal{db: db},
	}
}

func (m *disputeModelReal) UploadEvidence(userId, disputeId int64, path string, file io.Reader) (uint, error) {
	env := env.NewEnvLoader()
	fileStorageRoot, err := env.Get("FILESTORAGE_ROOT")
	if err != nil {
		return 0, err
	}
	fileStorageUrl, err := env.Get("FILESTORAGE_URL")
	if err != nil {
		return 0, err
	}

	logger := utilities.NewLogger().LogWithCaller()

	dir, name := filepath.Split(path)
	dest := filepath.Join(fileStorageRoot, path)

	if err := os.MkdirAll(filepath.Join(fileStorageRoot, dir), 0755); err != nil {
		logger.WithError(err).Error("failed to create folder for file upload")
		return 0, err
	}

	url := fmt.Sprintf("%s/%s", fileStorageUrl, path)

	// Open the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		logger.WithError(err).Errorf("failed to create file in storage: '%s'", dest)
		return 0, errors.New("failed to create file in storage")
	}
	defer destFile.Close()

	// Copy file content to destination
	_, err = io.Copy(destFile, file)
	if err != nil {
		logger.WithError(err).Error("failed to copy file content")
		return 0, errors.New("failed to copy file content")
	}

	// Add file entry to Database
	fileRow := models.File{
		FileName: name,
		FilePath: url,
		Uploaded: time.Now(),
	}

	if err := m.db.Create(&fileRow).Error; err != nil {
		logger.WithError(err).Error("error adding file to database")
		return 0, errors.New("error adding file to database")
	}

	//add entry to dispute evidence table
	disputeEvidence := models.DisputeEvidence{
		Dispute: disputeId,
		FileID:  int64(*fileRow.ID),
		UserID:  userId,
	}

	if err := m.db.Create(&disputeEvidence).Error; err != nil {
		logger.WithError(err).Error("error creating dispute evidence")
		return 0, errors.New("error creating dispute evidence")
	}

	return *fileRow.ID, nil
}

func (m *disputeModelReal) GetDispute(disputeId int64) (dispute models.Dispute, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).First(&dispute).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
	}
	return dispute, err
}

func (m *disputeModelReal) GetEvidenceByDispute(disputeId int64) (evidence []models.Evidence, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Table("dispute_evidence").Select(`files.id, file_name, uploaded, file_path,  CASE
            WHEN disputes.complainant = dispute_evidence.user_id THEN 'Complainant'
            WHEN disputes.respondant = dispute_evidence.user_id THEN 'Respondent'
            ELSE 'Other'
        END AS uploader_role`).Joins("JOIN files ON dispute_evidence.file_id = files.id").Joins("JOIN disputes ON dispute_evidence.dispute = disputes.id").Where("dispute = ?", disputeId).Find(&evidence).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute evidence")
		return
	}
	if evidence == nil {
		evidence = []models.Evidence{}
	}
	return evidence, err
}
func (m *disputeModelReal) GetDisputeExperts(disputeId int64) (experts []models.Expert, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Table("dispute_experts").Select("users.id, users.first_name || ' ' || users.surname AS full_name, email, users.phone_number AS phone, role").Joins("JOIN users ON dispute_experts.user = users.id").Where("dispute = ?", disputeId).Where("dispute_experts.status = 'Approved'").Where("role = 'Mediator' OR role = 'Arbitrator' OR role = 'Conciliator'").Find(&experts).Error
	if err != nil && err.Error() != "record not found" {
		logger.WithError(err).Error("Error retrieving dispute experts")
		return
	}

	if experts == nil {
		experts = []models.Expert{}
	}
	return experts, err
}

func (m *disputeModelReal) GetDisputesByUser(userId int64) (disputes []models.Dispute, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Where("complainant = ? OR respondant = ?", userId, userId).Find(&disputes).Error
	if err != nil {
		logger.WithError(err).Errorf("Failed to find disputes of user with ID %d", userId)
	}
	return disputes, err
}

func (m *disputeModelReal) GetUserByEmail(email string) (user models.User, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		logger.WithError(err).Errorf("Failed to find user with email '%s'", email)
	}
	return user, err
}
func (m *disputeModelReal) CreateDispute(dispute models.Dispute) (int64, error) {
	logger := utilities.NewLogger().LogWithCaller()
	disputeCloned := dispute
	if err := m.db.Create(&disputeCloned).Error; err != nil {
		logger.WithError(err).Error("Error creating dispute in database")
		return 0, err
	}
	return *disputeCloned.ID, nil
}
func (m *disputeModelReal) UpdateDisputeStatus(disputeId int64, status string) error {
	logger := utilities.NewLogger().LogWithCaller()
	err := m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).Update("status", status).Error
	if err != nil {
		logger.WithError(err).Errorf("Failed to update dispute (ID = %d) status to '%s'", disputeId, status)
	}
	return err
}
func (m *disputeModelReal) ObjectExpert(userId, disputeId, expertId int64, reason string) error {
	logger := utilities.NewLogger().LogWithCaller()

	//update dispute experts table
	var disputeExpert models.DisputeExpert
	if err := m.db.Where("dispute = ? AND dispute_experts.user = ?", disputeId, expertId).First(&disputeExpert).Error; err != nil {
		logger.WithError(err).Error("Error retrieving dispute expert")
		return err
	}

	disputeExpert.Status = models.ReviewStatus
	if err := m.db.Save(&disputeExpert).Error; err != nil {
		logger.WithError(err).Error("Error updating dispute expert")
		return err
	}

	//add entry to expert objections table
	expertObjection := models.ExpertObjection{
		DisputeID: disputeId,
		ExpertID:  expertId,
		UserID:    userId,
		Reason:    reason,
	}

	if err := m.db.Create(&expertObjection).Error; err != nil {
		logger.WithError(err).Error("Error creating expert objection")
		return err
	}
	return nil
}
func (m *disputeModelReal) ReviewExpertObjection(userId, disputeId, expertId int64, approved bool) error {
	logger := utilities.NewLogger().LogWithCaller()

	var expertObjections []models.ExpertObjection
	if err := m.db.Where("dispute_id = ? AND expert_id = ? AND status = ?", disputeId, expertId, models.ReviewStatus).Find(&expertObjections).Error; err != nil {
		logger.WithError(err).Error("Error retrieving expert objections")
		return err
	}

	var disputeExpert models.DisputeExpert
	if err := m.db.Where("dispute = ? AND dispute_experts.user = ? AND status = ?", disputeId, expertId, models.ReviewStatus).First(&disputeExpert).Error; err != nil {
		logger.WithError(err).Error("Error retrieving dispute expert")
		return nil
	}

	// Start a transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		logger.WithError(tx.Error).Error("Error starting transaction")
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Goroutine panicked, rolled back transaction")
			tx.Rollback()
		}
	}()

	// Update expert objections
	if approved {
		for i := range expertObjections {
			expertObjections[i].Status = models.Sustained
		}
		disputeExpert.Status = models.RejectedStatus
	} else {
		for i := range expertObjections {
			expertObjections[i].Status = models.Overruled
		}
		disputeExpert.Status = models.ApprovedStatus
	}

	// Save the expert objections
	for _, objection := range expertObjections {
		if err := tx.Save(&objection).Error; err != nil {
			tx.Rollback()
			logger.WithError(err).Error("Error updating expert objections")
			return err
		}
	}

	// Save the dispute expert
	if err := tx.Save(&disputeExpert).Error; err != nil {
		tx.Rollback()
		logger.WithError(err).Error("Error updating dispute expert")
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logger.WithError(err).Error("Error committing transaction")
		return err
	}
	return nil
}
