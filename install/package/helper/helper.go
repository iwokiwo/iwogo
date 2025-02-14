package helper

import (
	"encoding/json"
	"fmt"
	"iwogo/helper/paginator"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ResponseV2 struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	//Meta    entity.Meta `json:"meta"` // Make sure it's entity.Meta
	Meta interface{} `json:"meta,omitempty"`
}

type MetaBase struct {
	Pagination *paginator.Pagination `json:"pagination,omitempty"`
}

// Paginate is a generic struct for paginated responses
type PaginateBase[T any] struct {
	Meta   MetaBase `json:"meta"`
	Result []T      `json:"result"`
}

// FilterCondition represents a dynamic filter condition
type FilterCondition struct {
	Field    string
	Operator string // Example: "=", ">", "<", ">=", "<=", "LIKE"
	Value    interface{}
}

// ApplyFilters dynamically applies filtering to a GORM query
func ApplyFilters(db *gorm.DB, conditions []FilterCondition) *gorm.DB {
	// for _, condition := range conditions {
	// 	switch strings.ToUpper(condition.Operator) {
	// 	case "LIKE":
	// 		db = db.Where(fmt.Sprintf("%s LIKE ?", condition.Field), "%"+condition.Value.(string)+"%")
	// 	case "IN":
	// 		db = db.Where(fmt.Sprintf("%s IN (?)", condition.Field), condition.Value)
	// 	default:
	// 		db = db.Where(fmt.Sprintf("%s %s ?", condition.Field, condition.Operator), condition.Value)
	// 	}
	// }
	// return db

	for _, condition := range conditions {
		var value interface{}

		// Format time.Time to "MM/DD/YYYY" if needed
		if v, ok := condition.Value.(time.Time); ok {
			value = v.Format("01/02/2006") // Convert to MM/DD/YYYY
		} else {
			value = condition.Value
		}

		switch strings.ToUpper(condition.Operator) {
		case "LIKE":
			db = db.Where(fmt.Sprintf("%s LIKE ?", condition.Field), "%"+value.(string)+"%")
		case "IN":
			db = db.Where(fmt.Sprintf("%s IN (?)", condition.Field), value)
		default:
			db = db.Where(fmt.Sprintf("%s %s ?", condition.Field, condition.Operator), value)
		}
	}
	return db
}

func LoggerFile(message string, types string, idUser int, errors interface{}) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	//--------------------------setting file log from .env---------------------------
	path := os.Getenv("LOGFILE") + "application.log"
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger.SetOutput(file)

	//-------------------------switch type logger------------------------------------
	switch types {
	case "info":
		logger.WithField("userId", idUser).WithField("Errors", errors).Info(message)
	case "Warn":
		logger.WithField("userId", idUser).WithField("Errors", errors).Warnf(message)
	case "Error":
		logger.WithField("userId", idUser).WithField("Errors", errors).Error(message)
	}
}

func APIResponse(message string, code int, status string, data interface{}) ResponseV2 {
	// var metadata interface{}
	// if paginatedData, ok := data.(entity.Paginate); ok {
	// 	data = paginatedData.Result
	// 	metadata = paginatedData.Meta
	// }

	var metadata interface{}

	// Use reflection to detect generic Paginate[T]
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Struct {
		// Check if the struct is of type Paginate[T]
		if metaField := val.FieldByName("Meta"); metaField.IsValid() {
			if resultField := val.FieldByName("Result"); resultField.IsValid() {
				metadata = metaField.Interface() // Extract metadata
				data = resultField.Interface()   // Extract actual data
			}
		}
	}

	return ResponseV2{
		Message: message,
		Code:    code,
		Status:  status,
		Data:    data,
		Meta:    metadata,
	}
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	// Check if it's a validation error
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = e.Tag() // Adjust based on how you want to format the error
		}
		return errors
	}

	// Check if it's a JSON unmarshal type error
	if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
		errors[unmarshalErr.Field] = "Invalid type"
		return errors
	}

	// Generic error fallback
	errors["error"] = err.Error()
	return errors
}

func BindJSONAndValidate[T any](c *gin.Context, input *T) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		errors := FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, APIResponse("Invalid input", http.StatusUnprocessableEntity, "error", gin.H{"errors": errors}))
		return false
	}
	return true
}

func BindFormAndValidate[T any](c *gin.Context, input *T) bool {
	if err := c.ShouldBind(input); err != nil {
		errors := FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, APIResponse("Invalid input", http.StatusUnprocessableEntity, "error", gin.H{"errors": errors}))
		return false
	}
	return true
}
