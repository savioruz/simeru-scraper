package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/simeru-scraper/internal/cores/services"
	"github.com/savioruz/simeru-scraper/pkg/utils"
)

type ScheduleHandler struct {
	service   *services.ScheduleService
	validator *utils.Validator
}

func NewScheduleHandler(service *services.ScheduleService, validator *utils.Validator) *ScheduleHandler {
	return &ScheduleHandler{
		service:   service,
		validator: validator,
	}
}

// GetSchedule function is a handler to get schedule from the service
// It will return the schedule based on the request
// @Summary Get Schedule
// @Description Get schedule based on the request
// @Tags Schedule
// @Accept json
// @Produce json
// @Param data body ScheduleRequest true "Schedule Request"
// @Success 200 {object} ScheduleResponseSuccess
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/schedule [post]
func (h *ScheduleHandler) GetSchedule(c *fiber.Ctx) error {
	var req ScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return HandleError(c, fiber.StatusBadRequest, errors.New("invalid request"))
	}

	if err := h.validator.Validate(req); err != nil {
		return HandleError(c, fiber.StatusBadRequest, err)
	}

	schedule, err := h.service.GetSchedule(req.StudyPrograms, req.Day)
	if err != nil {
		return HandleError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(ScheduleResponseSuccess{
		Data: schedule,
	})
}

// GetStudyPrograms function is a handler to get study programs from the service
// It will return the study programs based on the request
// @Summary Get Study Programs
// @Description Get study programs based on the request
// @Tags StudyPrograms
// @Accept json
// @Produce json
// @Param data body StudyProgramsRequest true "Study Programs Request"
// @Success 200 {object} StudyProgramsResponseSuccess
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/study-programs [post]
func (h *ScheduleHandler) GetStudyPrograms(c *fiber.Ctx) error {
	var req StudyProgramsRequest
	if err := c.BodyParser(&req); err != nil {
		return HandleError(c, fiber.StatusBadRequest, errors.New("invalid request"))
	}

	studyPrograms, err := h.service.GetStudyPrograms(req.Faculty)
	if err != nil {
		return HandleError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(StudyProgramsResponseSuccess{
		Data: studyPrograms,
	})
}
