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
// It will return the study programs based on the faculty provided in the query parameters
// @Summary Get Study Programs
// @Description Get study programs based on the faculty in the query parameters
// @Tags StudyPrograms
// @Accept json
// @Produce json
// @Param faculty query string false "Faculty"
// @Success 200 {object} StudyProgramsResponseSuccess
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/study-programs [get]
func (h *ScheduleHandler) GetStudyPrograms(c *fiber.Ctx) error {
	faculty := c.Query("faculty")
	if err := h.validator.Validate(StudyProgramsRequest{Faculty: faculty}); err != nil {
		return HandleError(c, fiber.StatusBadRequest, err)
	}

	// Fetch study programs based on the faculty
	studyPrograms, err := h.service.GetStudyPrograms(faculty)
	if err != nil {
		return HandleError(c, fiber.StatusInternalServerError, err)
	}

	// Return the study programs in the response
	return c.Status(fiber.StatusOK).JSON(StudyProgramsResponseSuccess{
		Data: studyPrograms,
	})
}
