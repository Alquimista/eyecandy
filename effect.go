// Package eyecandy read/write SSA/ASS Subtitle Script
package eyecandy

import (
	"github.com/Alquimista/eyecandy/reader"
	"github.com/Alquimista/eyecandy/writer"
)

// ALIGN SSA/ASS alignment map
var ALIGN = writer.ALIGN

// ENC SSA/ASS encoding map
var ENC = writer.ENC

type Effect struct {
	scriptIn           *reader.Script
	scriptOut          *writer.Script
	Resolution         [2]int // WIDTH, HEIGHT
	VideoPath          string
	VideoZoom          float64
	VideoPosition      int
	VideoAR            float64
	MetaFilename       string
	MetaTitle          string
	MetaOriginalScript string
	MetaTranslation    string
	MetaTiming         string
	Audio              string
}

func (fx *Effect) Dialogs() (dialogs []*writer.Dialog) {
	for _, dlg := range fx.scriptIn.Dialog.NotCommented() {
		d := writer.NewDialog(dlg.Text)
		d.Layer = dlg.Layer
		d.Start = dlg.Start
		d.End = dlg.End
		d.StyleName = dlg.StyleName
		d.Actor = dlg.Actor
		d.Effect = dlg.Effect
		d.Margin = dlg.Margin
		d.Comment = dlg.Comment
		dialogs = append(dialogs, d)
	}
	return dialogs
}

func (fx *Effect) Styles() (dialogs map[string]*reader.Style) {
	return fx.scriptIn.StyleUsed
}

func (fx *Effect) AddStyle(sty *writer.Style) {
	fx.scriptOut.AddStyle(sty)
}

func (fx *Effect) AddDialog(d *writer.Dialog) {
	fx.scriptOut.AddDialog(d)
}

func (fx *Effect) Save(fn string) {
	fx.scriptOut.Resolution = fx.Resolution
	fx.scriptOut.VideoPath = fx.VideoPath
	fx.scriptOut.VideoZoom = fx.VideoZoom
	fx.scriptOut.VideoPosition = fx.VideoPosition
	fx.scriptOut.VideoAR = fx.VideoAR
	fx.scriptOut.MetaFilename = fx.MetaFilename
	fx.scriptOut.MetaTitle = fx.MetaTitle
	fx.scriptOut.MetaOriginalScript = fx.MetaOriginalScript
	fx.scriptOut.MetaTranslation = fx.MetaTranslation
	fx.scriptOut.MetaTiming = fx.MetaTiming
	fx.scriptOut.Audio = fx.Audio
	fx.scriptOut.Save(fn)
}

func NewEffect(inFN string) *Effect {
	input := reader.Read(inFN)
	output := writer.NewScript()

	for _, style := range input.StyleUsed {
		s := NewStyle(style.Name)
		s.Name = style.Name
		s.FontName = style.FontName
		s.FontSize = style.FontSize
		s.Color = style.Color
		s.Bold = style.Bold
		s.Italic = style.Italic
		s.Underline = style.Underline
		s.StrikeOut = style.StrikeOut
		s.Scale = style.Scale
		s.Spacing = style.Spacing
		s.Angle = style.Angle
		s.OpaqueBox = style.OpaqueBox
		s.Bord = style.Bord
		s.Shadow = style.Shadow
		s.Alignment = style.Alignment
		s.Margin = style.Margin
		s.Encoding = style.Encoding
		output.AddStyle(s)
	}

	// Add the original karaoke commented by default in the script"""
	// This help to jump to the wanted line in the preview in Aegisub,
	// and/or keep a backup of the timed subs
	dok := NewDialog("### Original Karaoke ###")
	dok.Comment = true
	output.AddDialog(dok)
	for _, dlg := range input.Dialog {
		d := writer.NewDialog(dlg.Text)
		d.Layer = dlg.Layer
		d.Start = dlg.Start
		d.End = dlg.End
		d.StyleName = dlg.StyleName
		d.Actor = dlg.Actor
		d.Effect = dlg.Effect
		d.Margin = dlg.Margin
		d.Comment = true
		output.AddDialog(d)
	}
	dke := NewDialog("### Karaoke Effect ###")
	dke.Comment = true
	output.AddDialog(dke)

	return &Effect{
		scriptIn:           &input,
		scriptOut:          output,
		Resolution:         input.Resolution,
		VideoPath:          input.VideoPath,
		VideoZoom:          input.VideoZoom,
		VideoPosition:      input.VideoPosition,
		VideoAR:            input.VideoAR,
		MetaFilename:       input.MetaFilename,
		MetaTitle:          input.MetaTitle,
		MetaOriginalScript: input.MetaOriginalScript,
		MetaTranslation:    input.MetaTranslation,
		MetaTiming:         input.MetaTiming,
		Audio:              input.Audio,
	}
}

func NewStyle(name string) *writer.Style {
	return writer.NewStyle(name)
}

func NewDialog(text string) *writer.Dialog {
	return writer.NewDialog(text)
}
