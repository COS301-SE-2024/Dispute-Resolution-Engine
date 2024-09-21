package handlers

import (
	"api/models"
	"api/utilities"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(g *gin.RouterGroup, h User) {
	g.PUT("/profile", h.UpdateUser)
	g.GET("/profile", h.GetUser)
	g.PUT("/profile/address", h.UpdateUserAddress)
	g.DELETE("/remove", h.RemoveAccount)
	g.POST("/analytics", h.UserAnalyticsEndpoint)
}

// @Summary Get user profile
// @Description Get user profile
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "User profile not available yet..."
// @Router /user/profile [get]
func (h User) GetUser(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	// Get the user ID from the request
	jwtClaims, err := h.Jwt.GetClaims(c)
	if err != nil {
		logger.WithError(err).Error("Unauthorized")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Retrieve the user from the database
	var dbUser models.User
	if err := h.DB.Where("email = ?", jwtClaims.Email).First(&dbUser).Error; err != nil {
		logger.WithError(err).Error("Failed to retrieve user")
		c.JSON(http.StatusNotFound, models.Response{Error: "User not found"})
		return
	}

	// Retrieve all addresses associated with the user from the database
	var dbAddresses models.Address
	h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddresses)
	if err := h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddresses).Error; err != nil {
		// If no addresses are found, create an empty address object
		dbAddresses = models.Address{}
	}

	// Create a response object
	user := models.GetUser{
		FirstName:         dbUser.FirstName,
		Surname:           dbUser.Surname,
		Email:             dbUser.Email,
		PhoneNumber:       dbUser.PhoneNumber,
		Birthdate:         dbUser.Birthdate.Format("2006-01-02"),
		Gender:            dbUser.Gender,
		Nationality:       dbUser.Nationality,
		Timezone:          dbUser.Timezone,
		PreferredLanguage: dbUser.PreferredLanguage,
		Address:           dbAddresses,
		Theme:             "dark",
	}

	logger.Info("User address retrieved successfully")
	// Return the response
	c.JSON(http.StatusOK, models.Response{Data: user})
}

// @Summary Update user profile
// @Description Update user profile
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UpdateUser true "User"
// @Success 200 {object} models.Response "User updated successfully"
// @Failure 400 {object} models.Response "Bad Request"
// @Router /user/profile [put]
func (h User) UpdateUser(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	jwtClaims, err := h.Jwt.GetClaims(c)
	if err != nil {
		logger.WithError(err).Error("Unauthorized")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	//get the user id from the request
	var updateUser models.UpdateUser
	if err := c.BindJSON(&updateUser); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Something went wrong..."})
		return
	}

	//retrieve the user from the database
	var dbUser models.User
	h.DB.Where("id = ?", jwtClaims.ID).First(&dbUser)

	var dbAddress models.Address
	h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddress)

	//check which fields to update
	if updateUser.FirstName != "" {
		dbUser.FirstName = updateUser.FirstName
	}
	if updateUser.Surname != "" {
		dbUser.Surname = updateUser.Surname
	}
	if updateUser.Phone_number != nil {
		dbUser.PhoneNumber = updateUser.Phone_number
	}
	if updateUser.Gender != "" {
		dbUser.Gender = updateUser.Gender
	}
	if updateUser.Nationality != "" {
		dbUser.Nationality = updateUser.Nationality
	}
	if updateUser.Timezone != nil {
		dbUser.Timezone = updateUser.Timezone
	}
	if updateUser.Preferred_language != nil {
		dbUser.PreferredLanguage = updateUser.Preferred_language
	}

	//now update the user and address
	h.DB.Model(&dbUser).Where("id = ?", dbUser.ID).Updates(dbUser)
	h.DB.Model(&dbAddress).Where("id = ?", dbUser.AddressID).Updates(dbAddress)
	logger.Info("User updated successfully")
	c.JSON(http.StatusOK, models.Response{Data: "User updated successfully"})
}

// @Summary Remove user account
// @Description Remove user account
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.DeleteUser true "User"
// @Success 200 {object} models.Response "User account removed successfully"
// @Failure 400 {object} models.Response "Bad Request"
// @Router /user/remove [delete]
func (h User) RemoveAccount(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	var user models.DeleteUser
	if err := c.BindJSON(&user); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}

	if !h.checkUserExists(user.Email) {
		logger.Error("User does not exist")
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	var dbUser models.User
	h.DB.Where("email = ?", user.Email).First(&dbUser)

	realSalt, err := base64.StdEncoding.DecodeString(dbUser.Salt)
	if err != nil {
		logger.WithError(err).Error("Failed to decode salt")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}

	checkHash := utilities.HashPasswordWithSalt(user.Password, realSalt)
	if dbUser.PasswordHash != base64.StdEncoding.EncodeToString(checkHash) {
		logger.Error("Invalid credentials")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Invalid credentials"})
		return
	}

	//delete address details
	var dbAddress models.Address
	h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddress)

	h.DB.Where("id = ?", dbUser.AddressID).Delete(&dbAddress)

	h.DB.Where("email = ?", user.Email).Delete(&dbUser)
	logger.Info("User account removed successfully")
	c.JSON(http.StatusOK, models.Response{Data: "User account removed successfully"})
}

func (h User) UpdateUserAddress(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	jwtClaims, err := h.Jwt.GetClaims(c)
	if err != nil {
		logger.WithError(err).Error("JWT claims is nil")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	insertAddress := false

	//here we get the details of the request
	var updateUserAddress models.UpdateAddress
	if err := c.BindJSON(&updateUserAddress); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}
	//retrieve the record from the database
	var dbUser models.User
	h.DB.Where("id = ?", jwtClaims.ID).First(&dbUser)

	if dbUser.AddressID == nil {
		insertAddress = true
	}
	//now we have to set the address parameters using the passed in data
	var dbAddress models.Address

	//first fetch the country code based on the name
	if updateUserAddress.Country != nil {
		var country models.Country
		h.DB.Where("country_code = ?", updateUserAddress.Country).First(&country)
		dbAddress.Country = updateUserAddress.Country
		dbAddress.CountryName = &country.CountryName
	}
	if updateUserAddress.Province != nil {
		dbAddress.Province = updateUserAddress.Province
	}
	if updateUserAddress.City != nil {
		dbAddress.City = updateUserAddress.City
	}
	if updateUserAddress.Street3 != nil {
		dbAddress.Street3 = updateUserAddress.Street3
	}
	if updateUserAddress.Street2 != nil {
		dbAddress.Street2 = updateUserAddress.Street2
	}
	if updateUserAddress.Street != nil {
		dbAddress.Street = updateUserAddress.Street
	}
	if updateUserAddress.AddressType != nil {
		dbAddress.AddressType = updateUserAddress.AddressType
	}

	if insertAddress {
		//insert the address
		h.DB.Create(&dbAddress)
		dbUser.AddressID = &dbAddress.ID
		h.DB.Model(&dbUser).Where("id = ?", dbUser.ID).Updates(dbUser)
		logger.Info("User address updated/inserted successfully")
		c.JSON(http.StatusOK, models.Response{Data: "User address updated successfully"})
		return
	}

	//now update the address
	h.DB.Model(&dbAddress).Where("id = ?", dbUser.AddressID).Updates(dbAddress)
	logger.Info("User address updated successfully")
	c.JSON(http.StatusOK, models.Response{Data: "User address updated successfully"})
}

// UserAnalyticsEndpoint is a handler for user analytics
// @Summary User analytics
// @Description User analytics
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UserAnalytics true "User"
// @Success 200 {object} models.Response "User analytics"
// @Failure 400 {object} models.Response "Bad Request"
// @Router /user/analytics [post]
func (h *Handler) UserAnalyticsEndpoint(c *gin.Context) {
	var analyticsReq models.UserAnalytics

	// Parse JSON request body
	if err := c.BindJSON(&analyticsReq); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request body"})
		return
	}

	// Construct the base query
	query := h.DB.Model(&models.User{})

	// Add WHERE clauses for column-value comparisons
	if analyticsReq.ColumnvalueComparisons != nil {
		for _, cvc := range *analyticsReq.ColumnvalueComparisons {
			query = query.Where(cvc.Column+" LIKE ?", "%"+cvc.Value+"%")
		}
	}

	// Add ORDER BY clauses
	if analyticsReq.OrderBy != nil {
		for _, ob := range *analyticsReq.OrderBy {
			switch ob.Order {
			case "asc":
				query = query.Order(ob.Column + " ASC")
			case "desc":
				query = query.Order(ob.Column + " DESC")
			}
		}
	}

	// Add date range filters
	if analyticsReq.DateRanges != nil {
		for _, dr := range *analyticsReq.DateRanges {
			startDate, err := time.Parse("2006-01-02", *dr.StartDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid start date format"})
				return
			}
			endDate, err := time.Parse("2006-01-02", *dr.EndDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid end date format"})
				return
			}

			query = query.Where(dr.Column+" BETWEEN ? AND ?", startDate, endDate)
		}
	}

	// Add GROUP BY clauses
	if analyticsReq.GroupBy != nil {
		for _, gb := range *analyticsReq.GroupBy {
			query = query.Group(gb)
		}
	}

	// Execute the query based on count or fetch all records
	if analyticsReq.Count {
		var count int64
		if err := query.Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to fetch user count"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count})
	} else {
		var users []models.User
		if err := query.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to fetch users"})
			return
		}
		c.JSON(http.StatusOK, models.Response{Data: users})
	}
}
