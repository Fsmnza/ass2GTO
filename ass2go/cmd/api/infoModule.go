package main

import (
	"ass2/internal/data"
	"ass2/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ModuleName     string `json:"module_name"`
		ModuleDuration int    `json:"module_duration"`
		ExamType       string `json:"exam_type"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	module := &data.ModuleInfo{
		ModuleName:     input.ModuleName,
		ModuleDuration: input.ModuleDuration,
		ExamType:       input.ExamType,
	}
	v := validator.New()
	if data.ValidateModuleInfo(v, module); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.InfoModel.Insert(module)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/module_infos/%d", module.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"module_info": module}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	module, err := app.models.InfoModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": module}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	module, err := app.models.InfoModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		ModuleName     string `json:"module_name"`
		ModuleDuration int    `json:"module_duration"`
		ExamType       string `json:"exam_type"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	module.ModuleName = input.ModuleName
	module.ModuleDuration = input.ModuleDuration
	module.ExamType = input.ExamType
	v := validator.New()
	if data.ValidateModuleInfo(v, module); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.InfoModel.Update(module)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": module}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.InfoModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "module_info successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	modules, err := app.models.InfoModel.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": modules}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
