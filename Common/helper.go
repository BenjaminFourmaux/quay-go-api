package Common

import (
	"fmt"
	"reflect"
	"strings"
)

type UpdateFieldMapping struct {
	ModelFieldName string
	Value          interface{}
}

func BuildUpdatedFields[TModel any, TDto any](dto TDto, mappings map[string]UpdateFieldMapping) map[string]interface{} {
	updatedFields := make(map[string]interface{})

	dtoValue := reflect.ValueOf(dto)
	if dtoValue.Kind() == reflect.Ptr {
		if dtoValue.IsNil() {
			return updatedFields
		}
		dtoValue = dtoValue.Elem()
	}

	if dtoValue.Kind() != reflect.Struct {
		return updatedFields
	}

	dtoType := dtoValue.Type()
	for i := 0; i < dtoValue.NumField(); i++ {
		dtoFieldMeta := dtoType.Field(i)
		dtoFieldValue := dtoValue.Field(i)

		if dtoFieldValue.Kind() != reflect.Ptr || dtoFieldValue.IsNil() {
			continue
		}

		mapping, hasMapping := mappings[dtoFieldMeta.Name]
		modelFieldName := dtoFieldMeta.Name
		if hasMapping && mapping.ModelFieldName != "" {
			modelFieldName = mapping.ModelFieldName
		}

		columnName := GetColumnName[TModel](modelFieldName)
		if columnName == "" {
			continue
		}

		value := dtoFieldValue.Elem().Interface()
		if hasMapping && mapping.Value != nil {
			value = mapping.Value
		}

		updatedFields[columnName] = value
	}

	return updatedFields
}

func GetColumnName[T any](fieldName string) string {
	modelType := reflect.TypeOf((*T)(nil)).Elem()
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return ""
	}

	field, ok := modelType.FieldByName(fieldName)
	if !ok {
		return ""
	}

	if columnName := getGormTagValue(field.Tag.Get("gorm"), "column"); columnName != "" {
		return columnName
	}

	return field.Name
}

func GetColumnNameByPointer[T any](model *T, fieldPtr any) string {
	if model == nil || fieldPtr == nil {
		return ""
	}

	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.IsNil() {
		return ""
	}

	structValue := modelValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return ""
	}

	fieldPointerValue := reflect.ValueOf(fieldPtr)
	if fieldPointerValue.Kind() != reflect.Ptr || fieldPointerValue.IsNil() {
		return ""
	}

	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		candidateField := structValue.Field(i)
		if !candidateField.CanAddr() {
			continue
		}

		if candidateField.Addr().Pointer() != fieldPointerValue.Pointer() {
			continue
		}

		fieldMeta := structType.Field(i)
		if columnName := getGormTagValue(fieldMeta.Tag.Get("gorm"), "column"); columnName != "" {
			return columnName
		}

		return fieldMeta.Name
	}

	return ""
}

func getGormTagValue(tag string, key string) string {
	for _, part := range strings.Split(tag, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		name, value, found := strings.Cut(part, ":")
		if found && strings.TrimSpace(name) == key {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

/*
SplitRepositoryNamespaced splits a namespaced repository string into its namespace and name components.
ex: org/my-image -> 'org' as namespace; 'my-image' as name
ex: my-image -> ” as namespace; 'my-image' as name
Returns an error if the input string is not in the expected format.
*/
func SplitRepositoryNamespaced(repositoryNamespaced string) (*string, string, error) {
	parts := strings.SplitN(repositoryNamespaced, "/", 2)
	if len(parts) == 1 { // repo name only, a non org/user scoped repository
		if IsValidRepositoryName(parts[0]) {
			return nil, parts[0], nil
		}
	} else {
		if IsValidOrganizationOrUserName(parts[0]) && IsValidRepositoryName(parts[1]) {
			return &parts[0], parts[1], nil
		}
	}

	return nil, "", fmt.Errorf("invalid repository namespaced %s", repositoryNamespaced)
}

/*
InlineIf Return a if the condition is true, otherwise return b.
This function is generic and can be used with any type.
*/
func InlineIf[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}
