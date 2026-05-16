package Common

import (
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
