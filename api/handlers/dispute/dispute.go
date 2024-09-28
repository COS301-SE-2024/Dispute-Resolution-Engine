package dispute

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.RouterGroup, h Dispute) {
	//dispute routes
	jwt := middleware.NewJwtMiddleware()
	g.Use(jwt.JWTMiddleware)

	g.GET("", h.GetSummaryListOfDisputes)
	g.POST("/create", h.CreateDispute)
	g.GET("/:id", h.GetDispute)

	g.POST("", h.GetSummaryListOfDisputes)
	g.POST("/:id/objections", h.ExpertObjection)
	g.PATCH("/objections/:id", h.ExpertObjectionsReview)
	g.POST("/experts/objections", h.ViewExpertRejections)
	g.POST("/:id/evidence", h.UploadEvidence)
	g.PUT("/:id/status", h.UpdateStatus)

	//patch is not to be integrated yet
	// disputeRouter.HandleFunc("/{id}", h.patchDispute).Methods(http.MethodPatch)

	//create dispute

	//archive routes
}

// @Summary Get a summary list of disputes
// @Description Get a summary list of disputes
// @Tags dispute
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Dispute Summary Endpoint"
// @Router /dispute/:id/evidence [post]
func (h Dispute) UploadEvidence(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	claims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	disputeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("Invalid Dispute ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.WithError(err).Error("Failed to parse form data")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Failed to parse form data"})
		return
	}

	files := form.File["files"]
	folder := fmt.Sprintf("%d", disputeId)
	for _, fileHeader := range files {
		path := filepath.Join(folder, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			logger.WithError(err).Error("error opening multipart file")
			utilities.InternalError(c)
			return
		}

		_, err = h.Model.UploadEvidence(claims.ID, int64(disputeId), path, file)
		if err != nil {
			logger.WithError(err).Error("Error uploading evidence")
			c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
			return
		}
	}
	logger.Info("Evidence uploaded successfully")

	h.AuditLogger.LogDisputeProceedings(models.Disputes, map[string]interface{}{"dispute_id": disputeId, "user": claims, "message": "Evidence uploaded"})
	// disputeProceedingsLogger.LogDisputeProceedings(models.Users, map[string]interface{}{"user": user, "message": "Failed login attempt"})
	c.JSON(http.StatusCreated, models.Response{
		Data: "Files uploaded",
	})
}

// @Summary Get a summary list of disputes
// @Description Get a summary list of disputes
// @Tags dispute
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Dispute Summary Endpoint"
// @Router /dispute [get]
func (h Dispute) GetSummaryListOfDisputes(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	jwtClaims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.WithError(err).Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	userID := jwtClaims.ID
	userRole := jwtClaims.Role

	if userRole == "admin" && c.Request.Method == "POST" {
		var reqAdminDisputes models.AdminDisputesRequest
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.WithError(err).Error("Error reading request body")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request body"})
			return
		}

		// Reset the body so it can be read again by BindJSON
		c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

		// Check if the body is valid JSON and not empty
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			logger.WithError(err).Error("Invalid JSON format")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request body"})
			return
		}

		// If the body contains no key-value pairs, consider it empty
		if len(bodyMap) == 0 {
			disputes, count, err := h.Model.GetAdminDisputes(nil, nil, nil, nil, nil, nil)
			if err != nil {
				logger.WithError(err).Error("error retrieving disputes")
				c.JSON(http.StatusInternalServerError, models.Response{Error: "Error while retrieving disputes"})
				return
			}
			if count == 0 {
				logger.Info("No disputes found")
				c.JSON(http.StatusOK, models.Response{Data: gin.H{
					"disputes": disputes,
					"total":    count,
				}})
				return
			}
			c.JSON(http.StatusOK, models.Response{Data: gin.H{
				"disputes": disputes,
				"total":    count,
			}})
			return
		}

		if err := c.BindJSON(&reqAdminDisputes); err != nil {
			logger.WithError(err).Error("Invalid request")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Request"})
			return
		}
		var searchTerm *string
		var limit *int
		var offset *int
		var sort *models.Sort
		var filters *[]models.Filter
		var dateFilter *models.DateFilter
		if reqAdminDisputes.Search != nil {
			searchTerm = reqAdminDisputes.Search
		}
		if reqAdminDisputes.Limit != nil {
			limit = reqAdminDisputes.Limit
		}
		if reqAdminDisputes.Offset != nil {
			offset = reqAdminDisputes.Offset
		}
		if reqAdminDisputes.Sort != nil {
			sort = reqAdminDisputes.Sort
		}
		if reqAdminDisputes.Filter != nil {
			filters = &reqAdminDisputes.Filter
		}
		if reqAdminDisputes.DateFilter != nil {
			dateFilter = &models.DateFilter{}
			if reqAdminDisputes.DateFilter.Filed != nil {
				dateFilter.Filed = reqAdminDisputes.DateFilter.Filed
			}
			if reqAdminDisputes.DateFilter.Resolved != nil {
				dateFilter.Resolved = reqAdminDisputes.DateFilter.Resolved
			}
		}
		disputes, count, err := h.Model.GetAdminDisputes(searchTerm, limit, offset, sort, filters, dateFilter)
		if err != nil {
			logger.WithError(err).Error("error retrieving disputes")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error while retrieving disputes"})
			return
		}
		if count == 0 {
			logger.Info("No matching disputes found")
			c.JSON(http.StatusOK, models.Response{Data: gin.H{
				"disputes": disputes,
				"total":    0,
			}})
			return
		}
		c.JSON(http.StatusOK, models.Response{Data: gin.H{
			"disputes": disputes,
			"total":    count,
		}})
		return
	}

	disputes, err := h.Model.GetDisputesByUser(userID)
	if err != nil {
		logger.WithError(err).Error("error retrieving disputes")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error while retrieving disputes"})
		return
	}

	summaries := make([]models.DisputeSummaryResponse, len(disputes))
	for i, dispute := range disputes {
		var role string = ""
		if dispute.Complainant == userID {
			role = "Complainant"
		} else if *(dispute.Respondant) == userID {
			role = "Respondant"
		}
		summaries[i] = models.DisputeSummaryResponse{
			ID:          *dispute.ID,
			Title:       dispute.Title,
			Description: dispute.Description,
			Status:      dispute.Status,
			Role:        &role,
		}
	}
	h.AuditLogger.LogDisputeProceedings(models.Disputes, map[string]interface{}{"user": jwtClaims, "message": "Dispute summaries retrieved"})
	logger.Info("Dispute summaries retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: summaries})
}

// @Summary Get a dispute
// @Description Get a dispute
// @Tags dispute
// @Accept json
// @Produce json
// @Param id path string true "Dispute ID"
// @Success 200 {object} models.Response "Dispute Detail Endpoint"
// @Router /dispute/{id} [get]
func (h Dispute) GetDispute(c *gin.Context) {

	logger := utilities.NewLogger().LogWithCaller()

	jwtClaims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.WithError(err).Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{Error: fmt.Sprintf("Invalid dispute id '%s'", idParam)})
		return
	}

	dispute, err := h.Model.GetDispute(int64(id))
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	userId := jwtClaims.ID
	role := ""
	//name and email
	// var respondantData models.User
	// err = h.DB.Where("id = ?", disputes.Respondant).Scan(&respondantData).Error
	// if err!=nil {

	// }

	if userId == dispute.Complainant {
		role = "Complainant"
	} else if userId == *(dispute.Respondant) {
		role = "Respondent"
	}

	DisputeDetailsResponse := models.DisputeDetailsResponse{
		ID:          *dispute.ID,
		Title:       dispute.Title,
		Description: dispute.Description,
		Status:      dispute.Status,
		DateCreated: dispute.CaseDate,
		Role:        role,
	}

	evidence, err := h.Model.GetEvidenceByDispute(int64(id))
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute evidence")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	if evidence == nil {
		evidence = []models.Evidence{}
	}

	experts, err := h.Model.GetDisputeExperts(int64(id))
	if err != nil && err.Error() != "record not found" {
		logger.WithError(err).Error("Error retrieving dispute experts")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	if experts == nil {
		experts = []models.Expert{}
	}

	DisputeDetailsResponse.Evidence = evidence
	DisputeDetailsResponse.Experts = experts

	h.AuditLogger.LogDisputeProceedings(models.Disputes, map[string]interface{}{"user": jwtClaims, "message": "Dispute details retrieved"})
	logger.Info("Dispute details retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: DisputeDetailsResponse})
	// c.JSON(http.StatusOK, models.Response{Data: "Dispute Detail Endpoint for ID: " + id})
}

func (h Dispute) CreateDispute(c *gin.Context) {

	logger := utilities.NewLogger().LogWithCaller()
	claims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.WithError(err).Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.WithError(err).Error("Error parsing form")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Failed to parse form data"})
		return
	}

	// Access form values
	if form.Value["title"] == nil || len(form.Value["title"]) == 0 {
		logger.WithError(err).Error("missing field in form: title")
		c.JSON(http.StatusBadRequest, models.Response{Error: "missing field in form: title"})
		return
	}
	title := form.Value["title"][0]

	if form.Value["description"] == nil || len(form.Value["description"]) == 0 {
		logger.Error("missing field in form: description")
		c.JSON(http.StatusBadRequest, models.Response{Error: "missing field in form: description"})
		return
	}
	description := form.Value["description"][0]

	if form.Value["respondent[email]"] == nil || len(form.Value["respondent[email]"]) == 0 {
		logger.Error("missing field in form: respondent[email]")
		c.JSON(http.StatusBadRequest, models.Response{Error: "missing field in form: respondent[email]"})
		return
	}
	email := form.Value["respondent[email]"][0]

	if form.Value["respondent[full_name]"] == nil || len(form.Value["respondent[full_name]"]) == 0 || len(strings.Split(form.Value["respondent[full_name]"][0], " ")) < 2 {
		logger.Error("missing field in form: respondent[full_name]")
		c.JSON(http.StatusBadRequest, models.Response{Error: "missing field in form: respondent[full_name]"})
		return
	}
	fullName := form.Value["respondent[full_name]"][0]

	if form.Value["respondent[workflow]"] == nil || len(form.Value["respondent[workflow]"]) == 0 {
		logger.Error("missing field in form: respondent[workflow]")
		c.JSON(http.StatusBadRequest, models.Response{Error: "missing field in form: respondent[workflow]"})
		return
	}
	workflow := form.Value["respondent[workflow]"][0]
	workflwIdInt, err := strconv.Atoi(workflow)
	if err != nil {
		logger.WithError(err).Error("Cannot convert workflow ID to integer")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Workflow ID"})
		return
	}

	// telephone := form.Value["respondent[telephone]"][0]

	//get complainants id
	complainantID := claims.ID

	//check if respondant is in database by email and phone number
	var respondantID *int64
	respondent, err := h.Model.GetUserByEmail(email)
	defaultAccount := false
	//so if the error is "record not found"
	if err != nil {
		//if the user is not found in the database then we create the default user
		if err.Error() == "record not found" {
			logger.Info("Attempting to create default user")
			//now we call to create the default user
			secretPass := make([]byte, 5)
			// Fill the byte slice with random values
			_, err := rand.Read(secretPass)
			if err != nil {
				logger.WithError(err).Error("Error generating default password")
				c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating default password"})
				return
			}

			// Convert the byte slice to a base64 encoded string
			pass := base64.StdEncoding.EncodeToString(secretPass)
			err1 := h.Model.CreateDefaultUser(email, fullName, pass)
			if err1 != nil {
				logger.WithError(err1).Error("Error creating default user.")
				c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating default user."})
				return
			}
			go h.Email.SendDefaultUserEmail(c, email, pass, title, description)
			logger.Info("Default respondent user created")
			respondent, err = h.Model.GetUserByEmail(email)
			if err != nil {
				logger.WithError(err).Error("Error fetching the default respondent.")
				c.JSON(http.StatusInternalServerError, models.Response{Error: "Error fetching the default respondent."})
				return
			}
			logger.Info("Default respondent retreived.")
			defaultAccount = true
		} else {
			logger.Error("Error retrieving respondent")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving respondent"})
			return
		}
	}
	respondantID = &respondent.ID

	//get The workflow by id
	workflowData, err := h.Model.GetWorkflowRecordByID(uint64(workflwIdInt))
	if err != nil {
		logger.WithError(err).Error("Error retrieving workflow")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving workflow"})
		return
	}

	//create active workflow entry
	activeWorkflow := &models.ActiveWorkflows{
		Workflow:         int64(workflowData.ID),
		DateSubmitted:    time.Now(),
		WorkflowInstance: workflowData.Definition,
	}
	err = h.Model.CreateActiverWorkflow(activeWorkflow)
	if err != nil {
		logger.WithError(err).Error("Error creating active workflow")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating active workflow"})
		return
	}
	// Get the environment variables
	url, err := h.Env.Get("ORCH_URL")
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	port, err := h.Env.Get("ORCH_PORT")
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	startEndpoint, err := h.Env.Get("ORCH_START")
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// Send the request to the orchestrator
	payload := OrchestratorRequest{ID: activeWorkflow.ID}
	_, err = h.OrchestratorEntity.MakeRequestToOrchestrator(fmt.Sprintf("http://%s:%s%s", url, port, startEndpoint), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		//delete the active workflow from table
		h.Model.DeleteActiveWorkflow(activeWorkflow)
		return
	}

	//create entry into the dispute table
	disputeId, err := h.Model.CreateDispute(models.Dispute{
		Title:       title,
		CaseDate:    time.Now(),
		Workflow:    nil,
		Status:      "Awaiting Respondant",
		Description: description,
		Complainant: complainantID,
		Respondant:  respondantID,
	})
	if err != nil {
		logger.WithError(err).Error("Error creating dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute"})
		return
	}

	// Store files in Docker and retrieve URLs
	files := form.File["files"]
	folder := fmt.Sprintf("%d", disputeId)
	for _, fileHeader := range files {
		path := filepath.Join(folder, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			logger.WithError(err).Error("failed to open file")
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong."})
			return
		}

		_, err = h.Model.UploadEvidence(complainantID, disputeId, path, file)
		if err != nil {
			logger.WithError(err).Error("failed to upload evidence")
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong."})
			return
		}
	}

	//asssign experts to dispute
	selected, err := h.Model.AssignExpertsToDispute(disputeId)
	if err != nil {
		logger.WithError(err).Error("Error assigning experts to dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error assigning experts to dispute"})
		return
	}
	logger.Info("Assigned experts", selected)

	// Respond with success message
	if !defaultAccount {
		go h.Email.SendAdminEmail(c, disputeId, email, title, description)
	}
	logger.Info("Admin email sent")
	c.JSON(http.StatusCreated, models.Response{Data: models.DisputeCreationResponse{DisputeID: disputeId}})
	h.AuditLogger.LogDisputeProceedings(models.Disputes, map[string]interface{}{"user": claims, "message": "Dispute created and admin email sent"})
	logger.Info("Dispute created successfully: ", title)
}

func (h Dispute) UpdateStatus(c *gin.Context) {
	var disputeStatus models.DisputeStatusChange
	logger := utilities.NewLogger().LogWithCaller()
	if err := c.BindJSON(&disputeStatus); err != nil {
		logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request body"})
		return
	}

	disputeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("Invalid Dispute ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	err = h.Model.UpdateDisputeStatus(int64(disputeId), disputeStatus.Status)
	if err != nil {
		logger.WithError(err).Error("failed to update dispute status")
		utilities.InternalError(c)
		return
	}
	go h.Email.NotifyDisputeStateChanged(c, int64(disputeId), disputeStatus.Status)

	logger.Info("Dispute status updated successfully")

	jwtClaims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}
	h.AuditLogger.LogDisputeProceedings(models.Disputes, map[string]interface{}{"user": jwtClaims, "message": "Dispute status update successful"})
	c.JSON(http.StatusOK, models.Response{Data: "Dispute status update successful"})
}

// @Summary Update a dispute
// @Description Update a dispute
// @Tags dispute
// @Accept json
// @Produce json
// @Param id path string true "Dispute ID"
// @Success 200 {object} models.Response "Dispute Patch Endpoint"
// @Router /dispute/{id} [patch]
func (h Dispute) patchDispute(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")
	logger.Info("Dispute Patch Endpoint for ID: ", id)
	c.JSON(http.StatusOK, models.Response{Data: "Dispute Patch Endpoint for ID: " + id})
}

func (h Dispute) ExpertObjection(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	//get dispute id
	disputeId := c.Param("id")
	disputeIdInt, err := strconv.Atoi(disputeId)
	if err != nil {
		logger.WithError(err).Error("Cannot convert dispute ID to integer")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	//get info from body of post
	var req models.ExpertRejectRequest
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid json"})
		return
	}

	//check empty fields
	if req.ExpertID == nil || req.Reason == nil {
		logger.Error("Missing fields in request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing fields in request"})
		return
	}


	//get user properties from token
	claims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.Error("Unauthorized access attempt in function expertObjection")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	//get Admin
	admin, err := h.Model.GetUserById(*req.ExpertID)
	if err != nil {
		logger.WithError(err).Error("Failed to get Expert ID")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to get Expert ID"})
		return
	}

	//check that the user is assigned to the dispute
	assigned, err := h.Model.GetExperts(int64(disputeIdInt))
	if err != nil {
		logger.WithError(err).Error("Failed to get assigned experts")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to get assigned experts"})
		return
	}

	isAsssigned := false
	for _, expert := range assigned {
		if expert.ExpertID == *req.ExpertID {
			isAsssigned = true
			break
		}
	}

	if !isAsssigned {
		logger.Error("Expert is not assigned to dispute")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Expert is not assigned to dispute"})
		return
	}


	//create ticket
	titleTicket := "Objection against " + admin.FirstName + " " + admin.Surname + "On Dispute " + disputeId
	ticket, err := h.TicketModel.CreateTicket(claims.ID, int64(disputeIdInt), titleTicket, *req.Reason)
	if err != nil {
		logger.WithError(err).Error("Failed to create ticket")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to create ticket"})
		return
	}

	err = h.Model.ObjectExpert(int64(disputeIdInt), *req.ExpertID, ticket.ID)
	if err != nil {
		logger.WithError(err).Error("Failed to object to expert")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong"})
		return
	}

	h.AuditLogger.LogDisputeProceedings(models.Disputes, map[string]interface{}{"user": claims, "message": "Expert rejected suggestion"})
	logger.Info("Expert rejected suggestion")

	c.JSON(http.StatusOK, models.Response{Data: ticket.ID})
}

func (h Dispute) ExpertObjectionsReview(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	// Get dispute id
	objectionId := c.Param("id")
	objectionIdInt, err := strconv.Atoi(objectionId)
	if err != nil {
		logger.WithError(err).Error("Cannot convert dispute ID to integer")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	// Get info from token
	claims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.WithError(err).Error("Unauthorized access attempt", claims, err)
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Get body of post
	var req models.RejectExpertReview
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid body"})
		return
	}

	// Check empty fields
	if req.Status == nil {
		logger.Error("Missing fields in request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing fields in request"})
		return
	}

	err = h.Model.ReviewExpertObjection(int64(objectionIdInt), *req.Status)
	if err != nil {
		logger.WithError(err).Error("failed to review objection")
		c.JSON(http.StatusBadRequest, models.Response{Error: "failed to review objection"})
		return
	}

	logger.Info("Expert objections reviewed successfully")
	h.AuditLogger.LogDisputeProceedings(models.Disputes, map[string]interface{}{"user": claims, "message": "Expert objections reviewed successfully"})

	c.JSON(http.StatusNoContent, models.Response{Data: "Expert objections reviewed successfully"})
}

func (h Dispute) ViewExpertRejections(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	// get body of post
	var req models.ViewExpetRejectionsRequest
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Body"})
		return
	}

	//query the database
	rejections, err := h.Model.GetExpertRejections(req.Expert_id, req.Dispute_id, req.Limits, req.Offset)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve expert rejections")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	logger.Info("Expert rejections retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: rejections})
}
