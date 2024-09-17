package repositories

import (
	"errors"
	"fmt"
	"github.com/savioruz/simeru-scraper/internal/cores/entities"
	"log"
	"strings"
)

func (s *DB) GetSchedule(studyPrograms, day string) (*[]entities.RowData, error) {
	if studyPrograms == "" || day == "" {
		return nil, errors.New("study programs cannot be empty")
	}

	var schedule []entities.RowData

	key := fmt.Sprintf("schedule:studyPrograms:%s:day:%s", strings.ReplaceAll(studyPrograms, " ", "_"), day)
	log.Printf("key: %s", key)
	err := s.cache.Get(key, &schedule)
	if err != nil {
		return nil, errors.New("could not retrieve schedule from cache")
	}

	return &schedule, nil
}

func (s *DB) GetStudyPrograms(faculty string) (*[]entities.StudyPrograms, error) {
	if faculty == "" {
		return nil, errors.New("faculty cannot be empty")
	}

	var studyPrograms []entities.StudyPrograms

	key := fmt.Sprintf("studyPrograms:faculty:%s", strings.ReplaceAll(faculty, " ", "_"))
	err := s.cache.Get(key, &studyPrograms)
	if err != nil {
		return nil, errors.New("could not retrieve study programs from cache")
	}

	return &studyPrograms, nil
}
