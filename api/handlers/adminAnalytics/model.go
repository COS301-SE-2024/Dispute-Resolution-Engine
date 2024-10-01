package adminanalytics

import (
	"api/env"
	"api/middleware"
	"api/models"
	"errors"

	"gorm.io/gorm"
)

type AdminAnalyticsDBModel interface {
	CalculateAverageResolutionTimeByMonth() (map[string]float64, error)
	CalculateAverageResolutionTime() (float64, error)
	CountRecordsWithGroupBy(
		tableName string,
		column *string,
		value *interface{},
		groupBy *string,
	) (map[string]int64, error)
	GetDisputeGroupingByCountry() (map[string]int, error)
	CountDisputesByMonth(
		tableName string,
		dateColumn string,
	) (map[string]int64, error)
}

type AdminAnalyticsHandler struct {
	DB        AdminAnalyticsDBModel
	EnvReader env.Env
	JWT       middleware.Jwt
}

type AdminAnalyticsDBModelReal struct {
	DB  *gorm.DB
	env env.Env
}

func NewAdminAnalyticsHandler(db *gorm.DB, envReader env.Env) AdminAnalyticsHandler {
	return AdminAnalyticsHandler{
		DB:        AdminAnalyticsDBModelReal{DB: db, env: envReader},
		EnvReader: envReader,
		JWT:       middleware.NewJwtMiddleware(),
	}
}

func (a AdminAnalyticsDBModelReal) CalculateAverageResolutionTimeByMonth() (map[string]float64, error) {
	// Map to store average resolution times grouped by month
	averageTimeByMonth := make(map[string]float64)

	// Define a struct to hold the query result
	type MonthAvgResolution struct {
		Month       string  // Month in the format 'YYYY-MM'
		AverageTime float64 // Average resolution time in days
	}

	// Slice to hold the results of the query
	var results []MonthAvgResolution

	// Query to calculate the average resolution time grouped by month (ignoring NULLs)
	err := a.DB.Model(&models.Dispute{}).
		Select("TO_CHAR(case_date, 'YYYY-MM') AS month, avg(date_resolved - case_date) AS average_time").
		Where("date_resolved IS NOT NULL").
		Group("TO_CHAR(case_date, 'YYYY-MM')").
		Order("month").
		Scan(&results).Error

	// Handle any errors during the query execution
	if err != nil {
		return nil, err
	}

	// Convert the results into the map
	for _, result := range results {
		averageTimeByMonth[result.Month] = result.AverageTime
	}

	return averageTimeByMonth, nil
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
func (h AdminAnalyticsDBModelReal) CountRecordsWithGroupBy(
	tableName string,
	column *string,
	value *interface{},
	groupBy *string,
) (map[string]int64, error) {
	// Define a list of allowed columns for validation
	allowedColumns := map[string]map[string]bool{
		"disputes": {
			"status": true,
		},
		"tickets": {
			"status": true,
		},
		"users": {
			"role":              true,
			"status":            true,
			"gender":            true,
			"preferred_language": true,
			"timezone":          true,
		},
		"expert_objections": {
			"status": true,
		},
		"expert_objections_view": {
			"expert_full_name": true,
			"dispute_title":    true,
		},
	}

	// Validate the table name
	if _, ok := allowedColumns[tableName]; !ok {
		return nil, errors.New("invalid table name")
	}

	// Validate the group by column
	if groupBy != nil {
		if _, ok := allowedColumns[tableName][*groupBy]; !ok {
			return nil, errors.New("invalid group by column name")
		}
	}

	// Validate the WHERE clause column and value
	if column != nil {
		if value == nil {
			return nil, errors.New("value must be provided when column is provided")
		}
		if _, ok := allowedColumns[tableName][*column]; !ok {
			return nil, errors.New("invalid column name for WHERE clause")
		}
	}

	// Result map to store the counts grouped by the `groupBy` column (if provided)
	recordCounts := make(map[string]int64)

	// Create a struct to hold the result from the query
	type GroupCount struct {
		GroupKey string
		Count    int64
	}

	// Slice to hold the results of the query
	var results []GroupCount

	// Start building the base query using the specified table name
	query := h.DB.Table(tableName)

	// Apply the optional WHERE clause if both column and value are provided
	if column != nil && value != nil {
		query = query.Where(gorm.Expr("? = ?", gorm.Expr(*column), *value))
	}

	// Apply the optional GROUP BY clause
	if groupBy != nil {
		query = query.Select("? as group_key, count(*) as count", gorm.Expr(*groupBy)).Group(*groupBy)
	} else {
		query = query.Select("count(*) as count")
	}

	// Execute the query and scan the results into the slice
	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// If grouping was applied, convert the results into the map
	if groupBy != nil {
		for _, result := range results {
			recordCounts[result.GroupKey] = result.Count
		}
	} else {
		// If no grouping, put the count into a map with a generic key
		recordCounts["total"] = results[0].Count
	}

	return recordCounts, nil
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


type MonthCount struct {
	Month        string
	DisputeCount int64
}
func (h AdminAnalyticsDBModelReal) CountDisputesByMonth(
	tableName string,
	dateColumn string,
) (map[string]int64, error) {
	// Define a list of allowed tables and their date columns for validation
	allowedDateColumns := map[string]bool{
		"case_date":     true,
		"date_resolved": true,
		"created_at":    true,
		"resolved_at":   true,
	}

	// Validate the date column
	if _, ok := allowedDateColumns[dateColumn]; !ok {
		return nil, errors.New("invalid date column name")
	}

	// Result map to store the counts grouped by month
	recordCounts := make(map[string]int64)

	// Slice to hold the results of the query
	var results []MonthCount

	// Build and execute the SQL query
	err := h.DB.Table(tableName).
		Select("TO_CHAR(" + dateColumn + ", 'YYYY-MM') AS month, COUNT(*) AS dispute_count").
		Group("TO_CHAR(" + dateColumn + ", 'YYYY-MM')").
		Order("month").
		Scan(&results).Error

	// Handle any errors during query execution
	if err != nil {
		return nil, err
	}

	// Convert the results into the map
	for _, result := range results {
		recordCounts[result.Month] = result.DisputeCount
	}

	return recordCounts, nil
}

