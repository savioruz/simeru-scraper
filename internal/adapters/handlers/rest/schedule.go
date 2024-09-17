package rest

import "github.com/savioruz/simeru-scraper/internal/cores/entities"

type ScheduleRequest struct {
	StudyPrograms string `json:"study_programs" validate:"required,min=3,max=255"`
	Day           string `json:"day" validate:"required,day"`
}

type ScheduleResponseSuccess struct {
	Data *[]entities.RowData `json:"data"`
}

type StudyProgramsRequest struct {
	Faculty string `json:"faculty" validate:"required,alphanum,min=3,max=255"`
}

type StudyProgramsResponseSuccess struct {
	Data *[]entities.StudyPrograms `json:"data"`
}
