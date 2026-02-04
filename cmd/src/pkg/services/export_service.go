package services

import (
	"bytes"
	"github.com/fogleman/gg"
	"gorm.io/gorm"
	"image/png"
)

type ExportServiceInterface interface {
	AssignmentsToPng(rosterID uint) ([]byte, error)
}

type ExportService struct {
	rosterService *RosterService
	db            *gorm.DB
}

func NewExportService(rs *RosterService, db *gorm.DB) *ExportService {
	return &ExportService{rs, db}
}

func (e *ExportService) AssignmentsToPng(rosterID uint) ([]byte, error) {
	savedShifts, _, err := e.rosterService.GetSavedRoster(rosterID)

	if err != nil {
		return nil, err
	}

	const rowHeight = 40
	const colWidthShift = 150
	const colWidthUsers = 400
	width := colWidthShift + colWidthUsers
	height := len(savedShifts) * rowHeight

	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	for i, shift := range savedShifts {
		y := float64(i * rowHeight)

		dc.SetRGB(0.8, 0.8, 0.8)
		dc.DrawRectangle(0, y, float64(width), rowHeight)
		dc.Stroke()

		dc.SetRGB(0.2, 0.2, 0.2)
		dc.DrawStringAnchored(shift.RosterShift.Name, 10, y+(rowHeight/2), 0, 0.5)

		userText := ""
		for j, user := range shift.Users {
			userText += user.Name
			if j < len(shift.Users)-1 {
				userText += ", "
			}
		}
		dc.SetRGB(0.3, 0.3, 0.9) // Blue-ish for names
		dc.DrawStringAnchored(userText, float64(colWidthShift)+10, y+(rowHeight/2), 0, 0.5)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
