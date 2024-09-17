package ports

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/savioruz/simeru-scraper/internal/cores/entities"
)

type ScrapeRepository interface {
	ScrapeStudyPrograms(ctx context.Context, opts ...chromedp.ExecAllocatorOption) error
	ScrapeSchedule(ctx context.Context, opts ...chromedp.ExecAllocatorOption) error
}

type ScheduleRepository interface {
	GetSchedule(studyPrograms, day string) (*[]entities.RowData, error)
	GetStudyPrograms(faculty string) (*[]entities.StudyPrograms, error)
}
