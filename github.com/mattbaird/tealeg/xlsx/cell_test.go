package xlsx

import (
	. "gopkg.in/check.v1"
)

type CellSuite struct{}

var _ = Suite(&CellSuite{})

// Test that we can set and get a Value from a Cell
func (s *CellSuite) TestValueSet(c *C) {
	// Note, this test is fairly pointless, it serves mostly to
	// reinforce that this functionality is important, and should
	// the mechanics of this all change at some point, to remind
	// us not to lose this.
	cell := Cell{}
	cell.Value = "A string"
	c.Assert(cell.Value, Equals, "A string")
}

// Test that GetStyle correctly converts the xlsxStyle.Fonts.
func (s *CellSuite) TestGetStyleWithFonts(c *C) {
	font := NewFont(10, "Calibra")
	style := *NewStyle()
	style.Font = *font

	cell := &Cell{Value: "123", style: style}
	style = cell.GetStyle()
	c.Assert(style, NotNil)
	c.Assert(style.Font.Size, Equals, 10)
	c.Assert(style.Font.Name, Equals, "Calibra")
}

// Test that SetStyle correctly translates into a xlsxFont element
func (s *CellSuite) TestSetStyleWithFonts(c *C) {
	file := NewFile()
	sheet := file.AddSheet("Test")
	row := sheet.AddRow()
	cell := row.AddCell()
	font := NewFont(12, "Calibra")
	style := *NewStyle()
	style.Font = *font
	cell.SetStyle(style)
	style = cell.GetStyle()
	_, xFont, _, _, _, _ := style.makeXLSXStyleElements()
	c.Assert(xFont.Sz.Val, Equals, "12")
	c.Assert(xFont.Name.Val, Equals, "Calibra")
}

// Test that GetStyle correctly converts the xlsxStyle.Fills.
func (s *CellSuite) TestGetStyleWithFills(c *C) {
	fill := *NewFill("solid", "FF000000", "00FF0000")
	style := *NewStyle()
	style.Fill = fill
	cell := &Cell{Value: "123", style: style}
	style = cell.GetStyle()
	_, _, xFill, _, _, _ := style.makeXLSXStyleElements()
	c.Assert(xFill.PatternFill.PatternType, Equals, "solid")
	c.Assert(xFill.PatternFill.BgColor.RGB, Equals, "00FF0000")
	c.Assert(xFill.PatternFill.FgColor.RGB, Equals, "FF000000")
}

// Test that SetStyle correctly updates xlsxStyle.Fills.
func (s *CellSuite) TestSetStyleWithFills(c *C) {
	file := NewFile()
	sheet := file.AddSheet("Test")
	row := sheet.AddRow()
	cell := row.AddCell()
	fill := NewFill("solid", "00FF0000", "FF000000")
	style := *NewStyle()
	style.Fill = *fill
	cell.SetStyle(style)
	style = cell.GetStyle()
	_, _, xFill, _, _, _ := style.makeXLSXStyleElements()
	xPatternFill := xFill.PatternFill
	c.Assert(xPatternFill.PatternType, Equals, "solid")
	c.Assert(xPatternFill.FgColor.RGB, Equals, "00FF0000")
	c.Assert(xPatternFill.BgColor.RGB, Equals, "FF000000")
}

// Test that GetStyle correctly converts the xlsxStyle.Borders.
func (s *CellSuite) TestGetStyleWithBorders(c *C) {
	border := *NewBorder("thin", "thin", "thin", "thin")
	style := *NewStyle()
	style.Border = border
	cell := Cell{Value: "123", style: style}
	style = cell.GetStyle()
	_, _, _, xBorder, _, _ := style.makeXLSXStyleElements()
	c.Assert(xBorder.Left.Style, Equals, "thin")
	c.Assert(xBorder.Right.Style, Equals, "thin")
	c.Assert(xBorder.Top.Style, Equals, "thin")
	c.Assert(xBorder.Bottom.Style, Equals, "thin")
}

// We can return a string representation of the formatted data
func (l *CellSuite) TestFormattedValue(c *C) {
	cell := Cell{Value: "37947.7500001"}
	negativeCell := Cell{Value: "-37947.7500001"}
	smallCell := Cell{Value: "0.007"}
	earlyCell := Cell{Value: "2.1"}

	cell.numFmt = "general"
	c.Assert(cell.FormattedValue(), Equals, "37947.7500001")
	negativeCell.numFmt = "general"
	c.Assert(negativeCell.FormattedValue(), Equals, "-37947.7500001")

	cell.numFmt = "0"
	c.Assert(cell.FormattedValue(), Equals, "37947")

	cell.numFmt = "#,##0" // For the time being we're not doing
	// this comma formatting, so it'll fall back to the related
	// non-comma form.
	c.Assert(cell.FormattedValue(), Equals, "37947")

	cell.numFmt = "0.00"
	c.Assert(cell.FormattedValue(), Equals, "37947.75")

	cell.numFmt = "#,##0.00" // For the time being we're not doing
	// this comma formatting, so it'll fall back to the related
	// non-comma form.
	c.Assert(cell.FormattedValue(), Equals, "37947.75")

	cell.numFmt = "#,##0 ;(#,##0)"
	c.Assert(cell.FormattedValue(), Equals, "37947")
	negativeCell.numFmt = "#,##0 ;(#,##0)"
	c.Assert(negativeCell.FormattedValue(), Equals, "(37947)")

	cell.numFmt = "#,##0 ;[red](#,##0)"
	c.Assert(cell.FormattedValue(), Equals, "37947")
	negativeCell.numFmt = "#,##0 ;[red](#,##0)"
	c.Assert(negativeCell.FormattedValue(), Equals, "(37947)")

	cell.numFmt = "0%"
	c.Assert(cell.FormattedValue(), Equals, "3794775%")

	cell.numFmt = "0.00%"
	c.Assert(cell.FormattedValue(), Equals, "3794775.00%")

	cell.numFmt = "0.00e+00"
	c.Assert(cell.FormattedValue(), Equals, "3.794775e+04")

	cell.numFmt = "##0.0e+0" // This is wrong, but we'll use it for now.
	c.Assert(cell.FormattedValue(), Equals, "3.794775e+04")

	cell.numFmt = "mm-dd-yy"
	c.Assert(cell.FormattedValue(), Equals, "11-22-03")

	cell.numFmt = "d-mmm-yy"
	c.Assert(cell.FormattedValue(), Equals, "22-Nov-03")
	earlyCell.numFmt = "d-mmm-yy"
	c.Assert(earlyCell.FormattedValue(), Equals, "1-Jan-00")

	cell.numFmt = "d-mmm"
	c.Assert(cell.FormattedValue(), Equals, "22-Nov")
	earlyCell.numFmt = "d-mmm"
	c.Assert(earlyCell.FormattedValue(), Equals, "1-Jan")

	cell.numFmt = "mmm-yy"
	c.Assert(cell.FormattedValue(), Equals, "Nov-03")

	cell.numFmt = "h:mm am/pm"
	c.Assert(cell.FormattedValue(), Equals, "6:00 pm")
	smallCell.numFmt = "h:mm am/pm"
	c.Assert(smallCell.FormattedValue(), Equals, "12:14 am")

	cell.numFmt = "h:mm:ss am/pm"
	c.Assert(cell.FormattedValue(), Equals, "6:00:00 pm")
	smallCell.numFmt = "h:mm:ss am/pm"
	c.Assert(smallCell.FormattedValue(), Equals, "12:14:47 am")

	cell.numFmt = "h:mm"
	c.Assert(cell.FormattedValue(), Equals, "18:00")
	smallCell.numFmt = "h:mm"
	c.Assert(smallCell.FormattedValue(), Equals, "00:14")

	cell.numFmt = "h:mm:ss"
	c.Assert(cell.FormattedValue(), Equals, "18:00:00")
	// This is wrong, but there's no eary way aroud it in Go right now, AFAICT.
	smallCell.numFmt = "h:mm:ss"
	c.Assert(smallCell.FormattedValue(), Equals, "00:14:47")

	cell.numFmt = "m/d/yy h:mm"
	c.Assert(cell.FormattedValue(), Equals, "11/22/03 18:00")
	smallCell.numFmt = "m/d/yy h:mm"
	c.Assert(smallCell.FormattedValue(), Equals, "12/30/99 00:14") // Note, that's 1899
	earlyCell.numFmt = "m/d/yy h:mm"
	c.Assert(earlyCell.FormattedValue(), Equals, "1/1/00 02:24") // and 1900

	cell.numFmt = "mm:ss"
	c.Assert(cell.FormattedValue(), Equals, "00:00")
	smallCell.numFmt = "mm:ss"
	c.Assert(smallCell.FormattedValue(), Equals, "14:47")

	cell.numFmt = "[h]:mm:ss"
	c.Assert(cell.FormattedValue(), Equals, "18:00:00")
	smallCell.numFmt = "[h]:mm:ss"
	c.Assert(smallCell.FormattedValue(), Equals, "14:47")

	cell.numFmt = "mmss.0" // I'm not sure about these.
	c.Assert(cell.FormattedValue(), Equals, "00.8640")
	smallCell.numFmt = "mmss.0"
	c.Assert(smallCell.FormattedValue(), Equals, "1447.999997")

	cell.numFmt = "yyyy\\-mm\\-dd"
	c.Assert(cell.FormattedValue(), Equals, "2003\\-11\\-22")

	cell.numFmt = "dd/mm/yy"
	c.Assert(cell.FormattedValue(), Equals, "22/11/03")
	earlyCell.numFmt = "dd/mm/yy"
	c.Assert(earlyCell.FormattedValue(), Equals, "01/01/00")

	cell.numFmt = "hh:mm:ss"
	c.Assert(cell.FormattedValue(), Equals, "18:00:00")
	smallCell.numFmt = "hh:mm:ss"
	c.Assert(smallCell.FormattedValue(), Equals, "00:14:47")

	cell.numFmt = "dd/mm/yy\\ hh:mm"
	c.Assert(cell.FormattedValue(), Equals, "22/11/03\\ 18:00")

	cell.numFmt = "yy-mm-dd"
	c.Assert(cell.FormattedValue(), Equals, "03-11-22")

	cell.numFmt = "d-mmm-yyyy"
	c.Assert(cell.FormattedValue(), Equals, "22-Nov-2003")
	earlyCell.numFmt = "d-mmm-yyyy"
	c.Assert(earlyCell.FormattedValue(), Equals, "1-Jan-1900")

	cell.numFmt = "m/d/yy"
	c.Assert(cell.FormattedValue(), Equals, "11/22/03")
	earlyCell.numFmt = "m/d/yy"
	c.Assert(earlyCell.FormattedValue(), Equals, "1/1/00")

	cell.numFmt = "m/d/yyyy"
	c.Assert(cell.FormattedValue(), Equals, "11/22/2003")
	earlyCell.numFmt = "m/d/yyyy"
	c.Assert(earlyCell.FormattedValue(), Equals, "1/1/1900")

	cell.numFmt = "dd-mmm-yyyy"
	c.Assert(cell.FormattedValue(), Equals, "22-Nov-2003")

	cell.numFmt = "dd/mm/yyyy"
	c.Assert(cell.FormattedValue(), Equals, "22/11/2003")

	cell.numFmt = "mm/dd/yy hh:mm am/pm"
	c.Assert(cell.FormattedValue(), Equals, "11/22/03 06:00 pm")

	cell.numFmt = "mm/dd/yyyy hh:mm:ss"
	c.Assert(cell.FormattedValue(), Equals, "11/22/2003 18:00:00")
	smallCell.numFmt = "mm/dd/yyyy hh:mm:ss"
	c.Assert(smallCell.FormattedValue(), Equals, "12/30/1899 00:14:47")

	cell.numFmt = "yyyy-mm-dd hh:mm:ss"
	c.Assert(cell.FormattedValue(), Equals, "2003-11-22 18:00:00")
	smallCell.numFmt = "yyyy-mm-dd hh:mm:ss"
	c.Assert(smallCell.FormattedValue(), Equals, "1899-12-30 00:14:47")
}
