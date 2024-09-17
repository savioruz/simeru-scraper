package services

import (
	"github.com/savioruz/simeru-scraper/internal/cores/entities"
	"github.com/savioruz/simeru-scraper/internal/cores/ports"
)

type ScheduleService struct {
	ScheduleRepository ports.ScheduleRepository
}

func NewScheduleService(repo ports.ScheduleRepository) *ScheduleService {
	return &ScheduleService{
		ScheduleRepository: repo,
	}
}

func (s *ScheduleService) GetSchedule(studyPrograms, day string) (*[]entities.RowData, error) {
	return s.ScheduleRepository.GetSchedule(studyPrograms, day)
}

func (s *ScheduleService) GetStudyPrograms(faculty string) (*[]entities.StudyPrograms, error) {
	return s.ScheduleRepository.GetStudyPrograms(faculty)
}
