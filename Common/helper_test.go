package Common

import "testing"

type helperTestModel struct {
	ID   int    `gorm:"primaryKey;column:team_id;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`
}

func TestGetColumnNameReturnsGormColumnTag(t *testing.T) {
	columnName := GetColumnName[helperTestModel]("ID")

	if columnName != "team_id" {
		t.Fatalf("expected team_id, got %q", columnName)
	}
}

func TestGetColumnNameFallsBackToFieldName(t *testing.T) {
	columnName := GetColumnName[helperTestModel]("Name")

	if columnName != "Name" {
		t.Fatalf("expected Name, got %q", columnName)
	}
}

func TestGetColumnNameSupportsPointerModel(t *testing.T) {
	columnName := GetColumnName[*helperTestModel]("ID")

	if columnName != "team_id" {
		t.Fatalf("expected team_id, got %q", columnName)
	}
}

func TestGetColumnNameReturnsEmptyStringWhenFieldDoesNotExist(t *testing.T) {
	columnName := GetColumnName[helperTestModel]("Unknown")

	if columnName != "" {
		t.Fatalf("expected empty string, got %q", columnName)
	}
}

func TestGetColumnNameByPointerReturnsGormColumnTag(t *testing.T) {
	model := helperTestModel{}
	columnName := GetColumnNameByPointer(&model, &model.ID)

	if columnName != "team_id" {
		t.Fatalf("expected team_id, got %q", columnName)
	}
}

func TestGetColumnNameByPointerFallsBackToFieldName(t *testing.T) {
	model := helperTestModel{}
	columnName := GetColumnNameByPointer(&model, &model.Name)

	if columnName != "Name" {
		t.Fatalf("expected Name, got %q", columnName)
	}
}

