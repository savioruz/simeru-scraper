package rest

import "github.com/savioruz/simeru-scraper/internal/cores/entities"

type ScheduleRequest struct {
	StudyPrograms string `json:"study_programs" validate:"required,min=3,max=255" example:"matematika"`
	Day           string `json:"day" validate:"required,day" example:"senin"`
}

type ScheduleResponseSuccess struct {
	Data *[]entities.RowData `json:"data"`
}

type StudyProgramsRequest struct {
	Faculty string `json:"faculty,omitempty" validate:"omitempty,alphanum" example:"pascasarjana"`
}

type StudyProgramsResponseSuccess struct {
	Data *[]entities.StudyPrograms `json:"data"`
}
