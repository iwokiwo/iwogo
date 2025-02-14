package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the structure of the JSON configuration file.
type Config struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

// Field represents a struct field (name, type, and JSON tag).
type Field struct {
	Name string
	Type string
	Tag  string
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

// parseModelStruct extracts struct fields from the Go model file.
func parseModelStruct(modelPath string) ([]Field, string, error) {
	src, err := os.ReadFile(modelPath)
	if err != nil {
		return nil, "", err
	}

	// Parse the Go source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, modelPath, src, parser.AllErrors)
	if err != nil {
		return nil, "", err
	}

	var fields []Field
	var structName string

	// Traverse the AST to find the struct definition
	ast.Inspect(node, func(n ast.Node) bool {
		// Check for type declarations
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if it's a struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Save struct name (e.g., "Rfid")
		structName = typeSpec.Name.Name

		// Extract struct fields
		for _, field := range structType.Fields.List {
			// Skip embedded structs
			if len(field.Names) == 0 {
				continue
			}

			fieldName := field.Names[0].Name
			var fieldType string

			switch t := field.Type.(type) {
			case *ast.Ident:
				fieldType = t.Name
			case *ast.SelectorExpr:
				fieldType = fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
			case *ast.ArrayType:
				if ident, ok := t.Elt.(*ast.Ident); ok {
					fieldType = "[]" + ident.Name
				}
			}

			// Extract JSON tag
			jsonTag := ""
			gormTag := ""

			if field.Tag != nil {
				tagValue := strings.Trim(field.Tag.Value, "`") // Remove backticks

				// Extract JSON tag
				jsonParts := strings.Split(tagValue, "json:\"")
				if len(jsonParts) > 1 {
					jsonTag = strings.Split(jsonParts[1], "\"")[0]
				}

				// Extract GORM tag
				gormParts := strings.Split(tagValue, "gorm:\"")
				if len(gormParts) > 1 {
					gormTag = strings.Split(gormParts[1], "\"")[0]
				}
			}

			// **Skip fields that have "foreignKey" in their GORM tag**
			if strings.Contains(gormTag, "foreignKey") {
				continue
			}

			fields = append(fields, Field{Name: fieldName, Type: fieldType, Tag: jsonTag})
		}

		// for _, field := range structType.Fields.List {
		// 	// Extract field name
		// 	if len(field.Names) == 0 {
		// 		// Handle embedded structs
		// 		if ident, ok := field.Type.(*ast.Ident); ok {
		// 			fields = append(fields, Field{Name: ident.Name, Type: ident.Name})
		// 		}
		// 		continue // Ignore embedded fields
		// 	}
		// 	fieldName := field.Names[0].Name

		// 	// Extract field type
		// 	var fieldType string
		// 	switch t := field.Type.(type) {
		// 	case *ast.Ident:
		// 		fieldType = t.Name
		// 	case *ast.SelectorExpr:
		// 		fieldType = fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
		// 	case *ast.ArrayType:
		// 		if ident, ok := t.Elt.(*ast.Ident); ok {
		// 			fieldType = "[]" + ident.Name
		// 		}
		// 	}

		// 	// Extract JSON tag
		// 	jsonTag := ""
		// 	if field.Tag != nil {
		// 		tagValue := field.Tag.Value
		// 		tagValue = strings.Trim(tagValue, "`") // Remove backticks
		// 		jsonParts := strings.Split(tagValue, "json:\"")
		// 		if len(jsonParts) > 1 {
		// 			jsonTag = strings.Split(jsonParts[1], "\"")[0]
		// 		}
		// 	}

		// 	fields = append(fields, Field{Name: fieldName, Type: fieldType, Tag: jsonTag})
		// }
		return false
	})

	if structName == "" {
		return nil, "", fmt.Errorf("no struct found in %s", modelPath)
	}

	return fields, structName, nil
}

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

// generateFiles creates DTO and entity files based on the extracted fields.
func generateFiles(config Config, structName string, fields []Field) error {
	// Define output paths
	dtoPath := filepath.Join(config.Name, "dto")
	//entityPath := filepath.Join(config.Name, "entity")

	// Create directories
	os.MkdirAll(dtoPath, os.ModePerm)
	//os.MkdirAll(entityPath, os.ModePerm)

	// Generate DTO and Entity filenames
	dtoFile := filepath.Join(dtoPath, "dto_"+config.Name+".go")
	//entityFile := filepath.Join(entityPath, config.Name+".go")

	// Check if we need to import "time"
	needsTimeImport := false
	for _, field := range fields {
		if field.Type == "time.Time" {
			needsTimeImport = true
			break
		}
	}

	// Generate DTO content
	dtoContent := fmt.Sprintf(`package dto

%s
type %s struct {
	Models.Base
`, addImports(needsTimeImport), structName)

	entityContent := fmt.Sprintf(`package entity

%s
type %s struct {
`, addImports(needsTimeImport), structName)

	paramsContent := fmt.Sprintf(`
type Params%s struct {
	%s
	paginator.PaginateReq
}
`, structName, structName)

	for _, field := range fields {
		jsonTag := fmt.Sprintf("`json:\"%s\"`", field.Tag)
		dtoContent += fmt.Sprintf("\t%s %s %s\n", field.Name, field.Type, jsonTag)
		entityContent += fmt.Sprintf("\t%s %s\n", field.Name, field.Type)
	}

	dtoContent += "}\n\n" + paramsContent
	entityContent += "}\n\n" + generateTableNameFunction(structName, config.Name)

	// Write files
	os.WriteFile(dtoFile, []byte(dtoContent), 0644)
	//os.WriteFile(entityFile, []byte(entityContent), 0644)

	fmt.Println("Generated DTO:", dtoFile)
	//fmt.Println("Generated Entity:", entityFile)

	return nil
}

// addImports returns the appropriate import statement.
//
//	func addImports(needsTime bool) string {
//		if needsTime {
//			return `import "time"`
//		}
//		return ""
//	}
func addImports(needsTime bool) string {
	imports := []string{}
	if needsTime {
		imports = append(imports, `"time"`)
	}
	imports = append(imports, `"iwogo/helper/paginator"`) // Ensure proper import formatting
	imports = append(imports, `"iwogo/Models"`)
	if len(imports) > 0 {
		return "import (\n\t" + strings.Join(imports, "\n\t") + "\n)"
	}
	return ""
}

// generateTableNameFunction returns the TableName() method.
func generateTableNameFunction(structName, tableName string) string {
	return fmt.Sprintf(`func (%s) TableName() string {
	return "%s"
}
`, structName, tableName)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <json-config>")
	}

	configFile := os.Args[1]
	config, err := readConfig(configFile)
	if err != nil {
		log.Fatal("Failed to read config:", err)
	}

	fields, structName, err := parseModelStruct(config.Model)
	if err != nil {
		log.Fatal("Failed to parse model:", err)
	}

	err = generateFiles(config, structName, fields)
	if err != nil {
		log.Fatal("Failed to generate files:", err)
	}
}
