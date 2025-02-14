package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Struct to parse JSON file
type Config struct {
	Name string `json:"name"`
}

type RouteData struct {
	ModuleName string
}

// Repository template
const repoTemplate = `package {{ .PackageName }}

import (
	"context"
	"iwogo/helper"
	"iwogo/helper/paginator"
	"iwogo/{{ .PackageName }}/dto"
	"iwogo/Models"
	"gorm.io/gorm"
)

type Repository interface {
	{{- range .Methods }}
	{{ .Name }}(ctx context.Context, param {{ .ParamType }}) ({{ .ReturnType }}, error)
	{{- end }}
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

{{ range .Methods }}
func (r *repository) {{ .Name }}(ctx context.Context, param {{ .ParamType }}) ({{ .ReturnType }}, error) {
	{{- if eq .Name "Create" }}
	result := r.db.Create(&param)
	if result.Error != nil {
		return param, result.Error
	}
	return param, nil
	{{- else if eq .Name "Update" }}
	err := r.db.Save(&param).Error
	if err != nil {
		return param, err
	}
	return param, nil
	{{- else if eq .Name "FindOne" }}
	var record Models.{{ .EntityName }}
	err := r.db.Where("id = ?", param.ID).Find(&record).Error
	if err != nil {
		return record, err
	}
	return record, nil
	{{- else if eq .Name "FindAll" }}
	res := helper.PaginateBase[Models.{{ .EntityName }}]{
		Meta: helper.MetaBase{
			Pagination: &paginator.Pagination{
				PaginateReq: param.PaginateReq,
			},
		},
	}
	var conditions []helper.FilterCondition

	if param.ID != 0 {
		conditions = append(conditions, helper.FilterCondition{Field: "id", Operator: "=", Value: param.ID})
	}

	db := helper.ApplyFilters(r.db, conditions)
	
	//--------------add this one for show data user if you have relation in model data--------------------
	//db = db.Preload("User")

	err := db.Scopes(paginator.Paginate(&res.Result, res.Meta.Pagination, db)).Find(&res.Result).Error
	if err != nil {
		return helper.PaginateBase[Models.{{ .EntityName }}]{}, err
	}
	return res, nil
	{{- else if eq .Name "Remove" }}
	err := r.db.Where("id = ?", param.ID).Delete(&param).Error
	if err != nil {
		return param, err
	}
	return param, nil
	{{- else }}
	// TODO: Implement logic
	return {{ .ReturnValue }}, nil
	{{- end }}
}
{{ end }}
`

// Service template
const serviceTemplate = `package {{ .PackageName }}

import (
	"context"
	"iwogo/helper"
	"iwogo/{{ .PackageName }}/dto"
	"iwogo/Models"
)

type Service interface {
	{{- range .Methods }}
	{{ .Name }}(ctx context.Context, param {{ .ParamType }}) ({{ .ReturnType }}, error)
	{{- end }}
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

{{ range .Methods }}
func (s *service) {{ .Name }}(ctx context.Context, param {{ .ParamType }}) ({{ .ReturnType }}, error) {
	{{- if eq .Name "Create" }}
	data, err := s.repository.Create(ctx, param)
	if err != nil {
		return param, err
	}
	return data, nil
	{{- else if eq .Name "Update" }}
	data, err := s.repository.Update(ctx, param)
	if err != nil {
		return param, err
	}
	return data, nil
	{{- else if eq .Name "FindOne" }}
	data, err := s.repository.FindOne(ctx, param)
	if err != nil {
		return data, err
	}
	return data, nil
	{{- else if eq .Name "FindAll" }}
	data, err := s.repository.FindAll(ctx, param)
	if err != nil {
		return data, err
	}
	return data, nil
	{{- else if eq .Name "Remove" }}
	data, err := s.repository.Remove(ctx, param)
	if err != nil {
		return param, err
	}
	return data, nil
	{{- else }}
	// TODO: Implement logic
	return {{ .ReturnValue }}, nil
	{{- end }}
}
{{ end }}
`

const controllerTemplate = `package {{ .PackageName }}

import (
	"iwogo/helper"
	"iwogo/{{ .PackageName }}/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

{{ range .Methods }}
func (cn *Controller) {{ .Name }}(c *gin.Context) {
	{{- if eq .Name "Create" }}
	var input {{ .ParamType }}
	if !helper.BindJSONAndValidate(c, &input) {
		return
	}

	data, err := cn.service.Create(c, input)
	if err != nil {
		response := helper.APIResponse("Create failed", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Data has been registered", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
	{{- else if eq .Name "Update" }}
	var input {{ .ParamType }}
	if !helper.BindJSONAndValidate(c, &input) {
		return
	}

	data, err := cn.service.Update(c, input)
	if err != nil {
		response := helper.APIResponse("Update failed", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Data updated", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
	{{- else if eq .Name "FindOne" }}
	var param {{ .ParamType }}
	if !helper.BindFormAndValidate(c, &param) {
		return
	}

	data, err := cn.service.FindOne(c, param)
	if err != nil {
		response := helper.APIResponse("Failed to get detail", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Data", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
	{{- else if eq .Name "FindAll" }}
	var param {{ .ParamType }}
	if !helper.BindFormAndValidate(c, &param) {
		return
	}

	data, err := cn.service.FindAll(c, param)
	if err != nil {
		response := helper.APIResponse("Failed to get detail", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Data", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
	{{- else if eq .Name "Remove" }}
	var input {{ .ParamType }}
	if !helper.BindFormAndValidate(c, &input) {
		return
	}

	data, err := cn.service.Remove(c, input)
	if err != nil {
		response := helper.APIResponse("Delete failed", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Data has been deleted", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
	{{- else }}
	// TODO: Implement logic
	{{- end }}
}
{{ end }}
`

const routeTemplate = `package routes

import (
	"iwogo/auth"
	"iwogo/middleware"
	"iwogo/{{ .ModuleName }}"
	"iwogo/modules/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func {{ .ModuleName | title }}Router(db *gorm.DB, api *gin.RouterGroup, auth auth.Service) *gin.RouterGroup {

	userService := user.NewService(user.NewRepository(db))
	controller := {{ .ModuleName }}.NewController({{ .ModuleName }}.NewService({{ .ModuleName }}.NewRepository(db)))

	api.GET("/{{ .ModuleName }}/list", middleware.AuthMiddleware(auth, userService), controller.FindAll)
	api.GET("/{{ .ModuleName }}/detail", middleware.AuthMiddleware(auth, userService), controller.FindOne)
	api.POST("/{{ .ModuleName }}/create", middleware.AuthMiddleware(auth, userService), controller.Create)
	api.PUT("/{{ .ModuleName }}/update", middleware.AuthMiddleware(auth, userService), controller.Update)
	api.DELETE("/{{ .ModuleName }}/delete", middleware.AuthMiddleware(auth, userService), controller.Remove)

	return api
}
`

type Method struct {
	Name        string
	ParamType   string
	ReturnType  string
	ReturnValue string
	EntityName  string
}

// Generate creates a route file from a JSON config.
func Generate(moduleName string) error {
	routeFileName := fmt.Sprintf("routes/%s_routes.go", strings.ToLower(moduleName))

	tmpl, err := template.New("route").Funcs(template.FuncMap{
		"title": strings.Title,
	}).Parse(routeTemplate)
	if err != nil {
		return err
	}

	// Create file
	file, err := os.Create(routeFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, RouteData{ModuleName: strings.ToLower(moduleName)})
	if err != nil {
		return err
	}

	fmt.Println("Generated route file:", routeFileName)
	return nil
}

// Function to parse JSON file and get folder name
func getFolderNameFromJSON(filePath string) (string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return "", err
	}

	return config.Name, nil
}

func createFile(filePath, content string) error {
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// readConfig reads and parses the JSON configuration file.
func readConfig(jsonPath string) (Config, error) {
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	return config, err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run generateAutoFile/generate_repo.go <json-file>")
	}

	jsonFile := os.Args[1]
	folderName, err := getFolderNameFromJSON(jsonFile)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
	}

	fmt.Println("Generating repository and service for:", folderName)

	err = os.MkdirAll(folderName, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	entityName := strings.Title(folderName)
	fileRepositoru := "repository"
	fileService := "service"
	fileController := "controller"

	repoFile := filepath.Join(folderName, fileRepositoru+".go")
	serviceFile := filepath.Join(folderName, fileService+".go")
	controllerFile := filepath.Join(folderName, fileController+".go")
	//dtoFile := filepath.Join(folderName, "dto", "dto_"+folderName+".go")
	// entityFile := filepath.Join(folderName, "entity", folderName+".go")

	methods := []Method{
		{"Create", "dto." + entityName, "dto." + entityName, "param", entityName},
		{"Update", "dto." + entityName, "dto." + entityName, "param", entityName},
		{"FindAll", "dto.Params" + entityName, "helper.PaginateBase[Models." + entityName + "]", "helper.PaginateBase[Models." + entityName + "]{}", entityName},
		{"FindOne", "dto.Params" + entityName, "Models." + entityName, "Models." + entityName + "{}", entityName},
		{"Remove", "dto." + entityName, "dto." + entityName, "param", entityName},
	}

	tmplRepo, _ := template.New("repository").Parse(repoTemplate)
	file, _ := os.Create(repoFile)
	defer file.Close()
	_ = tmplRepo.Execute(file, struct {
		PackageName string
		Methods     []Method
	}{folderName, methods})

	tmplService, _ := template.New("service").Parse(serviceTemplate)
	file, _ = os.Create(serviceFile)
	defer file.Close()
	_ = tmplService.Execute(file, struct {
		PackageName string
		Methods     []Method
	}{folderName, methods})

	tmplController, _ := template.New("controller").Parse(controllerTemplate)
	file, _ = os.Create(controllerFile)
	defer file.Close()
	_ = tmplController.Execute(file, struct {
		PackageName string
		Methods     []Method
	}{folderName, methods})

	// Generate Routes
	if err := Generate(entityName); err != nil {
		log.Fatal("Error generating Routes:", err)
	}
	// _ = createFile(dtoFile, strings.ReplaceAll(dtoTemplate, "{{ .EntityName }}", entityName))
	// _ = createFile(entityFile, strings.ReplaceAll(entityTemplate, "{{ .EntityName }}", entityName))

	log.Println("Generated files:", repoFile, serviceFile, controllerFile)
}
