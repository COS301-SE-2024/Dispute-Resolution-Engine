package handlers

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(g *gin.RouterGroup, h Handler) {
	g.PUT("/profile", h.updateUser)
	g.GET("/profile", h.getUser)
	g.PUT("/profile/address", h.UpdateUserAddress)
	g.DELETE("/remove", h.RemoveAccount)
}

// @Summary Get user profile
// @Description Get user profile
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "User profile not available yet..."
// @Router /user/profile [get]
func (h Handler) getUser(c *gin.Context) {
	// Get the user ID from the request
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Retrieve the user from the database
	var dbUser models.User
	if err := h.DB.Where("email = ?", jwtClaims.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{Error: "User not found"})
		return
	}

	// Retrieve all addresses associated with the user from the database
	var dbAddresses models.Address
	h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddresses)
	if err := h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddresses).Error; err != nil {
		dbAddresses = models.Address{}
	}

	// Create a response object
	user := models.GetUser{
		FirstName:         dbUser.FirstName,
		Surname:           dbUser.Surname,
		Email:             dbUser.Email,
		PhoneNumber:       dbUser.PhoneNumber,
		Birthdate:         dbUser.Birthdate.String(),
		Gender:            dbUser.Gender,
		Nationality:       dbUser.Nationality,
		Timezone:          dbUser.Timezone,
		PreferredLanguage: dbUser.PreferredLanguage,
		Address:           dbAddresses,
		Theme:             "dark",
	}

	// Print the dbAddresses object for debug purposes
	print(dbAddresses.Country)
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
func (h Handler) updateUser(c *gin.Context) {
	jwtClaims := middleware.GetClaims(c)

	//get the user id from the request
	var updateUser models.UpdateUser
    if err := c.BindJSON(&updateUser); err != nil {
		return
	}

	//retrieve the user from the database
	var dbUser models.User
	h.DB.Where("id = ?", jwtClaims.User.ID).First(&dbUser)

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
func (h Handler) RemoveAccount(c *gin.Context) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)

	var user models.DeleteUser
    if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	if !h.checkUserExists(user.Email) {
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	var dbUser models.User
	h.DB.Where("email = ?", user.Email).First(&dbUser)

	realSalt, err := base64.StdEncoding.DecodeString(dbUser.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}
	checkHash, err := hasher.GenerateHash([]byte(user.Password), realSalt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}

	if dbUser.PasswordHash != base64.StdEncoding.EncodeToString(checkHash.Hash) {
		print(dbUser.PasswordHash)
		print(base64.StdEncoding.EncodeToString(checkHash.Hash))
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Invalid credentials"})
		return
	}

	//delete address details
	var dbAddress models.Address
	h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddress)

	h.DB.Where("id = ?", dbUser.AddressID).Delete(&dbAddress)

	h.DB.Where("email = ?", user.Email).Delete(&dbUser)
	c.JSON(http.StatusOK, models.Response{Data: "User account removed successfully"})
}

func (h Handler) UpdateUserAddress(c *gin.Context) {
	jwtClaims := middleware.GetClaims(c)
	insertAddress := false

	//here we get the details of the request
	var updateUserAddress models.UpdateAddress
    if err := c.BindJSON(&updateUserAddress); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}
	//retrieve the record from the database
	var dbUser models.User
	h.DB.Where("id = ?", jwtClaims.User.ID).First(&dbUser)

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
		c.JSON(http.StatusOK, models.Response{Data: "User address updated successfully"})
		return
	}

	//now update the address
	h.DB.Model(&dbAddress).Where("id = ?", dbUser.AddressID).Updates(dbAddress)

	c.JSON(http.StatusOK, models.Response{Data: "User address updated successfully"})
}
