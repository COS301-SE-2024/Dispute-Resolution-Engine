package adminanalytics

import (
	"api/env"
	"api/middleware"
	"api/models"

	"gorm.io/gorm"
)

type AdminAnalyticsDBModel interface {
	CalculateAverageResolutionTime() (float64, error)
	GetDisputeGroupingByStatus() (map[string]int64, error)
	GetDisputeGroupingByCountry() (map[string]int, error)
	
}

type AdminAnalyticsHandler struct {
	DB        AdminAnalyticsDBModel
	EnvReader env.Env
	JWT       middleware.Jwt
}

type AdminAnalyticsDBModelReal struct {
	DB *gorm.DB
	env env.Env
}

func NewAdminAnalyticsHandler(db *gorm.DB, envReader env.Env) AdminAnalyticsHandler{
	return AdminAnalyticsHandler{
		DB: 	  AdminAnalyticsDBModelReal{DB: db, env: envReader},
		EnvReader: envReader,
		JWT:      middleware.NewJwtMiddleware(),
	}
}

func (a AdminAnalyticsDBModelReal) CalculateAverageResolutionTime() (float64, error) {
	var averageTime float64

	// Query to calculate the average resolution time in days for disputes with non-null DateResolved
	err := a.DB.Model(&models.Dispute{}).
		Select("avg(date_resolved - case_date)").
		Where("date_resolved IS NOT NULL").
		Scan(&averageTime).Error

	if err != nil {
		return 0, err
	}

	return averageTime, nil
}

// DisputeStatusCount represents the count of disputes grouped by status.
type DisputeStatusCount struct {
	Status string
	Count  int64
}

// GetDisputeGroupingByStatus counts the number of disputes grouped by their statuses.
func (h AdminAnalyticsDBModelReal) GetDisputeGroupingByStatus() (map[string]int64, error) {
	var results []DisputeStatusCount
	statusCounts := make(map[string]int64)

	// Query to count disputes grouped by their status
	err := h.DB.Model(&models.Dispute{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert the result into a map for easier usage
	for _, result := range results {
		statusCounts[result.Status] = result.Count
	}

	return statusCounts, nil
}

// GetDisputeGroupingByCountry counts the number of disputes grouped by the country of the complainant.
func (h AdminAnalyticsDBModelReal) GetDisputeGroupingByCountry() (map[string]int, error) {
	// This map will store the country as the key and the count of disputes as the value
	countryCounts := make(map[string]int)

	// Struct to hold the results of the query
	type CountryDisputeCount struct {
		Country string
		Count   int
	}

	var results []CountryDisputeCount

	// Perform the join and grouping query
	err := h.DB.
		Table("disputes").
		Select("addresses.country, count(disputes.id) as count").
		Joins("JOIN users ON disputes.complainant = users.id").
		Joins("JOIN addresses ON users.address_id = addresses.id").
		Group("addresses.country").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert the results into a map
	for _, result := range results {
		countryCounts[result.Country] = result.Count
	}

	return countryCounts, nil
}