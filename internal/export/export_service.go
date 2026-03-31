package export

import (
	"GEWIS-Rooster/internal/models"
	"bytes"
	_ "embed"
	"fmt"
	"image/png"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/font"
	"gorm.io/gorm"
)

//go:embed static/fonts/arial.ttf
var ArialRegular []byte

//go:embed static/fonts/arial-bold.ttf
var ArialBold []byte

type Service interface {
	AssignmentsToPng(rosterID uint) ([]byte, error)
}
type RosterProvider interface {
	GetSavedRoster(uint) ([]*models.SavedShift, []*models.SavedShiftOrdering, error)
}

type OrganProvider interface {
	GetMembersSettings(organID uint) ([]*models.UserOrgan, error)
}

type service struct {
	rosterService RosterProvider
	organService  OrganProvider
	db            *gorm.DB
}

func NewExportService(rs RosterProvider, op OrganProvider, db *gorm.DB) Service {
	return &service{rs, op, db}
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

func (e *service) AssignmentsToPng(rosterID uint) ([]byte, error) {
	savedShifts, _, err := e.rosterService.GetSavedRoster(rosterID)
	if err != nil {
		return nil, err
	}

	normalFont, boldFont, err := loadFonts()
	if err != nil {
		log.Err(err).Msg("failed to load embedded fonts")
		return nil, err
	}

	if len(savedShifts) == 0 {
		return nil, fmt.Errorf("roster %d contains no shifts", rosterID)
	}

	var roster models.Roster
	if err := e.db.First(&roster, "id = ?", rosterID).Error; err != nil {
		return nil, err
	}

	memberSettings, err := e.organService.GetMembersSettings(roster.OrganID)
	if err != nil {
		return nil, err
	}

	userNicknames := make(map[uint]string)
	for _, uo := range memberSettings {
		if uo.Username != "" {
			userNicknames[uo.UserID] = uo.Username
		}
	}

	tempDc := gg.NewContext(0, 0)
	tempDc.SetFontFace(normalFont)

	maxUserWidth := PngImage.ColWidthUsers
	for _, shift := range savedShifts {
		var rowTextW float64
		for i, u := range shift.Users {
			var textToMeasure string
			if nickname, exists := userNicknames[u.ID]; exists {
				tempDc.SetFontFace(normalFont)
				textToMeasure = nickname
			} else {
				tempDc.SetFontFace(boldFont)
				parts := strings.Split(strings.TrimSpace(u.Name), " ")
				if len(parts) > 0 {
					textToMeasure = parts[0]
				}
			}

			w, _ := tempDc.MeasureString(textToMeasure)
			rowTextW += w

			if i < len(shift.Users)-1 {
				tempDc.SetFontFace(normalFont)
				cw, _ := tempDc.MeasureString(", ")
				rowTextW += cw
			}
		}

		totalRowTextWidth := rowTextW + (PngImage.Padding * 2)
		if totalRowTextWidth > maxUserWidth {
			maxUserWidth = totalRowTextWidth
		}
	}

	width := int(PngImage.ColWidthShift + maxUserWidth)
	height := (len(savedShifts) + 1) * int(PngImage.RowHeight)

	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetHexColor("#f3f4f6")
	dc.DrawRectangle(0, 0, float64(width), PngImage.RowHeight)
	dc.Fill()

	dc.SetFontFace(normalFont)
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

		dc.SetHexColor("#e5e7eb")
		dc.DrawLine(0, y+PngImage.RowHeight, float64(width), y+PngImage.RowHeight)
		dc.Stroke()

		dc.SetFontFace(boldFont)
		dc.SetHexColor("#111827")
		dc.DrawStringAnchored(shift.RosterShift.Name, PngImage.Padding, y+(PngImage.RowHeight/2), 0, 0.5)

		currentX := PngImage.ColWidthShift + PngImage.Padding
		centerY := y + (PngImage.RowHeight / 2)

		for j, u := range shift.Users {
			parts := strings.Split(strings.TrimSpace(u.Name), " ")
			if len(parts) == 0 {
				continue
			}

			if nickname, exists := userNicknames[u.ID]; exists {
				dc.SetFontFace(boldFont)
				dc.SetHexColor("#4b5563")
				dc.DrawStringAnchored(nickname, currentX, centerY, 0, 0.5)
				nw, _ := dc.MeasureString(nickname)
				currentX += nw
			} else {
				firstName := parts[0]
				dc.SetFontFace(boldFont)
				dc.SetHexColor("#4b5563")
				dc.DrawStringAnchored(firstName, currentX, centerY, 0, 0.5)
				fw, _ := dc.MeasureString(firstName)
				currentX += fw
			}

			if j < len(shift.Users)-1 {
				comma := ", "
				dc.SetFontFace(boldFont)
				dc.SetHexColor("#4b5563")
				dc.DrawStringAnchored(comma, currentX, centerY, 0, 0.5)
				cw, _ := dc.MeasureString(comma)
				currentX += cw
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func loadFonts() (font.Face, font.Face, error) {
	fReg, err := truetype.Parse(ArialRegular)
	if err != nil {
		return nil, nil, err
	}

	fBold, err := truetype.Parse(ArialBold)
	if err != nil {
		return nil, nil, err
	}

	normalFont := truetype.NewFace(fReg, &truetype.Options{
		Size:    PngImage.FontSize,
		Hinting: font.HintingFull,
	})
	boldFont := truetype.NewFace(fBold, &truetype.Options{
		Size:    PngImage.FontSize,
		Hinting: font.HintingFull,
	})

	return normalFont, boldFont, nil
}
