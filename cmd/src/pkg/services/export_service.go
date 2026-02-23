package services

import (
	"bytes"
	"github.com/fogleman/gg"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"image/png"
	"strings"
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

var PngImage = struct {
	RowHeight     float64
	ColWidthShift float64
	ColWidthUsers float64
	Padding       float64
	FontSize      float64
}{
	RowHeight:     45.0,
	ColWidthShift: 160.0,
	ColWidthUsers: 400.0,
	Padding:       15.0,
	FontSize:      16.0,
}

func (e *ExportService) AssignmentsToPng(rosterID uint) ([]byte, error) {
	savedShifts, _, err := e.rosterService.GetSavedRoster(rosterID)
	if err != nil {
		return nil, err
	}

	tempDc := gg.NewContext(0, 0)
	if err := tempDc.LoadFontFace("cmd/src/static/fonts/arial.ttf", PngImage.FontSize); err != nil {
		log.Err(err).Msg("failed to load font for measurement")
	}

	// 2. Calculate the required width for the Users column
	maxUserWidth := PngImage.ColWidthUsers
	for _, shift := range savedShifts {
		names := []string{}
		for _, u := range shift.Users {
			names = append(names, u.Name)
		}
		userText := strings.Join(names, ", ")

		// Measure how wide this specific string is in pixels
		textW, _ := tempDc.MeasureString(userText)

		// Add padding to the measurement
		totalRowTextWidth := textW + (PngImage.Padding * 2)

		if totalRowTextWidth > maxUserWidth {
			maxUserWidth = totalRowTextWidth
		}
	}

	// 3. Final Image Dimensions
	width := int(PngImage.ColWidthShift + maxUserWidth)
	height := (len(savedShifts) + 1) * int(PngImage.RowHeight)

	dc := gg.NewContext(width, height)
	if err := dc.LoadFontFace("cmd/src/static/fonts/arial.ttf", PngImage.FontSize); err != nil {
		log.Err(err).Msg("failed to load font")
	}

	// --- Drawing Logic ---
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Header
	dc.SetHexColor("#f3f4f6")
	dc.DrawRectangle(0, 0, float64(width), PngImage.RowHeight)
	dc.Fill()

	dc.SetHexColor("#374151")
	dc.DrawStringAnchored("SHIFT", PngImage.Padding, PngImage.RowHeight/2, 0, 0.5)
	dc.DrawStringAnchored("ASSIGNED USERS", PngImage.ColWidthShift+PngImage.Padding, PngImage.RowHeight/2, 0, 0.5)

	for i, shift := range savedShifts {
		y := float64(i+1) * PngImage.RowHeight

		if i%2 == 0 {
			dc.SetHexColor("#f9fafb")
			dc.DrawRectangle(0, y, float64(width), PngImage.RowHeight)
			dc.Fill()
		}

		// Bottom Border
		dc.SetHexColor("#e5e7eb")
		dc.DrawLine(0, y+PngImage.RowHeight, float64(width), y+PngImage.RowHeight)
		dc.Stroke()

		// Shift Name
		dc.SetHexColor("#111827")
		dc.DrawStringAnchored(shift.RosterShift.Name, PngImage.Padding, y+(PngImage.RowHeight/2), 0, 0.5)

		// Users List
		names := []string{}
		for _, u := range shift.Users {
			names = append(names, u.Name)
		}
		userText := strings.Join(names, ", ")

		dc.SetHexColor("#4b5563")
		dc.DrawStringAnchored(userText, PngImage.ColWidthShift+PngImage.Padding, y+(PngImage.RowHeight/2), 0, 0.5)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
