package data

import (
	"ass2/internal/validator"
	"database/sql"
	"errors"
	"time"
)

type ModuleInfo struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ModuleName     string    `json:"module_name"`
	ModuleDuration int       `json:"module_duration"`
	ExamType       string    `json:"exam_type"`
	Version        int32     `json:"version"`
}

func ValidateModuleInfo(v *validator.Validator, module *ModuleInfo) {
	v.Check(module.ModuleName != "", "module_name", "must be provided")
	v.Check(len(module.ModuleName) <= 500, "module_name", "must not be more than 500 bytes long")
	v.Check(module.ModuleDuration != 0, "module_duration", "must be provided")
	v.Check(module.ModuleDuration > 0, "module_duration", "must be a positive integer")
	v.Check(module.ExamType != "", "exam_type", "must be provided")
}

type ModuleInfoModel struct {
	DB *sql.DB
}

func (m ModuleInfoModel) Insert(module *ModuleInfo) error {
	query := `
		INSERT INTO module_info (module_name, module_duration, exam_type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`
	args := []interface{}{module.ModuleName, module.ModuleDuration, module.ExamType}
	return m.DB.QueryRow(query, args...).Scan(&module.ID, &module.CreatedAt, &module.Version)
}

func (m ModuleInfoModel) Get(id int64) (*ModuleInfo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, module_name, module_duration, exam_type, version
FROM module_info
WHERE id = $1`
	var module ModuleInfo
	err := m.DB.QueryRow(query, id).Scan(
		&module.ID,
		&module.CreatedAt,
		&module.ModuleName,
		&module.ModuleDuration,
		&module.ExamType,
		&module.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &module, nil
}

func (m ModuleInfoModel) Update(module *ModuleInfo) error {
	query := `
UPDATE module_info
SET module_name = $1, module_duration = $2, exam_type = $3, updated_at = NOW(), version = version + 1
WHERE id = $4
RETURNING updated_at, version`
	args := []interface{}{
		module.ModuleName,
		module.ModuleDuration,
		module.ExamType,
		module.ID,
	}
	return m.DB.QueryRow(query, args...).Scan(&module.UpdatedAt, &module.Version)
}

func (m ModuleInfoModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
DELETE FROM module_info
WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m ModuleInfoModel) GetAll() ([]*ModuleInfo, error) {
	query := `
SELECT id, created_at, module_name, module_duration, exam_type, version
FROM module_info`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modules := []*ModuleInfo{}
	for rows.Next() {
		var module ModuleInfo
		err := rows.Scan(
			&module.ID,
			&module.CreatedAt,
			&module.ModuleName,
			&module.ModuleDuration,
			&module.ExamType,
			&module.Version,
		)
		if err != nil {
			return nil, err
		}
		modules = append(modules, &module)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}
