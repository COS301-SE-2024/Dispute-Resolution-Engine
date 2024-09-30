package adminanalytics

import (
	"api/env"
	"api/middleware"
	"api/models"
	"errors"

	"gorm.io/gorm"
)

type AdminAnalyticsDBModel interface {
	CalculateAverageResolutionTime() (float64, error)
	CountRecordsWithGroupBy(
		tableName string,
		column *string,
		value *interface{},
		groupBy *string,
	) (map[string]int64, error)
	GetDisputeGroupingByCountry() (map[string]int, error)
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
