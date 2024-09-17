package services

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/savioruz/simeru-scraper/internal/cores/ports"
)

type ScrapeService struct {
	ScrapeRepository ports.ScrapeRepository
}

func NewScrapeService(repo ports.ScrapeRepository) *ScrapeService {
	return &ScrapeService{
		ScrapeRepository: repo,
	}
}

func (s *ScrapeService) ScrapeStudyPrograms(ctx context.Context, opts ...chromedp.ExecAllocatorOption) error {
	return s.ScrapeRepository.ScrapeStudyPrograms(ctx, opts...)
}

func (s *ScrapeService) ScrapeSchedule(ctx context.Context, opts ...chromedp.ExecAllocatorOption) error {
	return s.ScrapeRepository.ScrapeSchedule(ctx, opts...)
}
