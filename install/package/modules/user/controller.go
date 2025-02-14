package user

import (
	"iwogo/auth"
	"iwogo/helper"
	"iwogo/modules/user/dto"
	"iwogo/modules/user/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userController struct {
	userService Service
	authService auth.Service
}

func NewUserController(userService Service, authService auth.Service) *userController {
	return &userController{userService, authService}
}

func (h *userController) RegisterUser(c *gin.Context) {
	// tangkap input
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service

	validate := validator.New()

	var input dto.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Validation account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	errValidate := validate.VarWithValue(input.Password, input.Phone, "eqfield")

	if errValidate != nil {
		response := helper.APIResponse("Register failed", http.StatusBadRequest, "error", errValidate.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// FIND EMAIL
	email := input.Email
	isAvailable, err := h.userService.IsEmailAvailable(email)
	if err != nil {
		response := helper.APIResponse("Register email failed", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if isAvailable == false {
		response := helper.APIResponse("Register email failed", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Generate token failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := entity.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userController) Login(c *gin.Context) {
	var input dto.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// token
	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := entity.FormatUser(loggedinUser, token)
	response := helper.APIResponse("Login success", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userController) GetUserProfile(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(entity.User)
	userId := currentUser.ID

	profile, err := h.userService.GetUserbyId(userId)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Please try again", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := entity.FormatProfile(profile)
	response := helper.APIResponse("Get Profile", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userController) GetUserByID(c *gin.Context) {
	var input dto.UserIdInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("User Validation failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	profile, err := h.userService.GetUserbyId(input.ID)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Please try again", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := entity.FormatProfile(profile)
	response := helper.APIResponse("Get User", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func (h *userController) ChekEmailAvailability(c *gin.Context) {
	// input email dari user
	// input email di mapping ke struct
	// struct input di passing ke service
	// service akan manggil repository
	// repository akan melakukan query ke database
	var input dto.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	email := input.Email
	isEmailAvailable, err := h.userService.IsEmailAvailable(email)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"
	metaStatus := "error"

	if isEmailAvailable {
		metaMessage = "Email is available"
		metaStatus = "success"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, metaStatus, data)
	c.JSON(http.StatusOK, response)

}

func (h *userController) ChangeEmailHandler(c *gin.Context) {
	var input dto.ChangeEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Change email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	currentUser := c.MustGet("currentUser").(entity.User)
	userId := currentUser.ID

	_, err = h.userService.ChangeEmailService(userId, input)
	if err != nil {
		response := helper.APIResponse("Change email failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}
	response := helper.APIResponse("Change email success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)

}

func (h *userController) FetchUser(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(entity.User)

	formatter := entity.FormatUser(currentUser, "")

	response := helper.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}

func (h *userController) ChangeNameHandler(c *gin.Context) {
	var input dto.ChangeNameInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Change name failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	currentUser := c.MustGet("currentUser").(entity.User)
	userId := currentUser.ID

	_, err = h.userService.ServiceChangeName(userId, input)
	if err != nil {
		response := helper.APIResponse("Change name failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	response := helper.APIResponse("Change name success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *userController) ChangeDetailHandler(c *gin.Context) {
	var input dto.ChangeDetailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Change Detail failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	_, err = h.userService.ChangeDetailService(input)
	if err != nil {
		response := helper.APIResponse("Change Detail failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	response := helper.APIResponse("Change Detail success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)

}

func (h *userController) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		response := helper.APIResponse("Get users failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}
	formatter := entity.FormatUsers(users)
	response := helper.APIResponse("Get Profile", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userController) ChangePassword(c *gin.Context) {
	var input dto.ChangePassword
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Change password failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	if input.Password != input.PasswordCheck {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Change password failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	currentUser := c.MustGet("currentUser").(entity.User)
	userId := currentUser.ID
	_, err = h.userService.ChangePassword(userId, input)

	if err != nil {
		response := helper.APIResponse("Change password failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	response := helper.APIResponse("Change password success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)

}

func (h *userController) DeleteUser(c *gin.Context) {
	var input dto.DeleteInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Delete user failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	del, err := h.userService.Delete(input)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Please try again", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if del != true {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Please try again", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Delete user", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
