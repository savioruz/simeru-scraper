package repositories

import (
	"errors"
	"fmt"
	"github.com/savioruz/simeru-scraper/internal/adapters/cache"
	"github.com/savioruz/simeru-scraper/internal/cores/entities"
	"strings"
)

func (s *DB) GetSchedule(studyPrograms, day string) (*[]entities.RowData, error) {
	if studyPrograms == "" || day == "" {
		return nil, errors.New("study programs cannot be empty")
	}

	var schedule []entities.RowData

	key := fmt.Sprintf("schedule:studyPrograms:%s:day:%s", strings.ReplaceAll(studyPrograms, " ", "_"), day)
	err := s.cache.Get(key, &schedule)
	if errors.Is(err, cache.ErrCacheMiss) {
		return nil, cache.ErrCacheMiss
	}

	if errors.Is(err, cache.ErrCacheFailed) {
		return nil, cache.ErrCacheFailed
	}

	if err != nil {
		return nil, errors.New("could not retrieve schedule from cache")
	}

	return &schedule, nil
}

func (s *DB) GetStudyPrograms(faculty string) (*[]entities.StudyPrograms, error) {
	var key string
	if faculty == "" {
		key = "studyPrograms:all"
	} else {
		key = fmt.Sprintf("studyPrograms:faculty:%s", strings.ReplaceAll(faculty, " ", "_"))
	}

	var studyPrograms []entities.StudyPrograms

	err := s.cache.Get(key, &studyPrograms)
	if errors.Is(err, cache.ErrCacheMiss) {
		return nil, cache.ErrCacheMiss
	}

	if errors.Is(err, cache.ErrCacheFailed) {
		return nil, cache.ErrCacheFailed
	}

	if err != nil {
		return nil, errors.New("could not retrieve study programs from cache")
	}

	return &studyPrograms, nil
}
