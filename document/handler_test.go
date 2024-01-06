package document

import (
	"reflect"
	"testing"
	"time"

	"github.com/signintech/gopdf"
)

type TextData struct {
	text  string
	style string
	size  int
}

type Cell struct {
	rect    gopdf.Rect
	text    string
	options gopdf.CellOption
}

type MockPdf struct {
	currentX, currentY float64
	currentFontSize    int
	currentStyle       string
	textData           map[float64]map[float64]TextData
	cells              map[float64]map[float64]Cell
}

func (p *MockPdf) AddText(x, y float64, content string) error {
	if p.textData[x] == nil {
		p.textData[x] = make(map[float64]TextData)
	}
	p.textData[x][y] = TextData{text: content, style: p.currentStyle, size: p.currentFontSize}
	return nil
}

func (p *MockPdf) AddMultilineText(x, y float64, content string) {
	if p.textData[x] == nil {
		p.textData[x] = make(map[float64]TextData)
	}
	p.textData[x][y] = TextData{text: content, style: p.currentStyle, size: p.currentFontSize}
}

func (p *MockPdf) AddFormattedText(x, y float64, content string, size int, style string) {
	if p.textData[x] == nil {
		p.textData[x] = make(map[float64]TextData)
	}
	p.textData[x][y] = TextData{text: content, style: style, size: size}
}

func (p *MockPdf) SetFontSize(size int) error {
	p.currentFontSize = size
	return nil
}
func (p *MockPdf) SetFontStyle(style string) error {
	p.currentStyle = style
	return nil
}
func (p *MockPdf) DefaultFontSize() {
	p.currentFontSize = 11

}
func (p *MockPdf) DefaultFontStyle() {
	p.currentStyle = "default"

}
func (p *MockPdf) SetPosition(x, y float64) {
	p.currentX = x
	p.currentY = y

}
func (p *MockPdf) CellWithOption(rectangle *gopdf.Rect, text string, opt gopdf.CellOption) error {
	if p.cells[p.currentX] == nil {
		p.cells[p.currentX] = make(map[float64]Cell)
	}
	p.cells[p.currentX][p.currentY] = Cell{rect: *rectangle, text: text, options: opt}
	return nil
}
func (p *MockPdf) Image(picPath string, x float64, y float64, rect *gopdf.Rect) error {

	return nil
}
func (p *MockPdf) WritePdf(pdfPath string) error {
	return nil
}

func TestAddAdressData(t *testing.T) {
	pdf := MockPdf{
		textData:        make(map[float64]map[float64]TextData),
		currentStyle:    "default",
		currentFontSize: 11,
	}
	type args struct {
		doc      PdfDoc
		receiver ReceiverAdress
	}
	tests := []struct {
		name string
		args args
		want MockPdf
	}{
		{"Create simple document",
			args{
				doc: &pdf,
				receiver: ReceiverAdress{
					Header: "Head",
					Name:   "Name",
					Adress: "Address",
					City:   "City",
				},
			},
			MockPdf{
				currentX:        0,
				currentY:        0,
				currentFontSize: 11,
				currentStyle:    "default",
				textData: map[float64]map[float64]TextData{
					130: {
						50: TextData{
							text:  "Head",
							style: "default",
							size:  11,
						},
						57: TextData{
							text:  "Name",
							style: "default",
							size:  11,
						},
						64: TextData{
							text:  "Address",
							style: "default",
							size:  11,
						},
						71: TextData{
							text:  "City",
							style: "default",
							size:  11,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddAdressData(tt.args.doc, tt.args.receiver)
			if !reflect.DeepEqual(pdf, tt.want) {
				t.Errorf("AddAdressData() = %v, want %v", pdf, tt.want)
			}
		})
	}
}

func TestAddTitle(t *testing.T) {
	now = func() time.Time {
		return time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC)
	}

	pdf := MockPdf{
		textData:        make(map[float64]map[float64]TextData),
		currentStyle:    "default",
		currentFontSize: 11,
	}
	type args struct {
		doc   PdfDoc
		input TitleWithDate
	}
	tests := []struct {
		name string
		args args
		want MockPdf
	}{
		{"Create title with date",
			args{
				doc: &pdf,
				input: TitleWithDate{
					Title: "Titel AAA",
					City:  "City",
				},
			},
			MockPdf{
				currentX:        0,
				currentY:        0,
				currentFontSize: 11,
				currentStyle:    "default",
				textData: map[float64]map[float64]TextData{
					130: {
						89: TextData{
							text:  "City, 01.01.2020",
							style: "default",
							size:  11,
						},
					},
					20: {
						100: TextData{
							text:  "Titel AAA",
							style: "bold",
							size:  14,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddTitle(tt.args.doc, tt.args.input)
			if !reflect.DeepEqual(pdf, tt.want) {
				t.Errorf("AddTitle() = \n%v, \nwant \n%v", pdf, tt.want)
			}
		})
	}
}

func TestAddTable(t *testing.T) {

	pdf := MockPdf{
		cells:           make(map[float64]map[float64]Cell),
		currentStyle:    "default",
		currentFontSize: 11,
	}
	type args struct {
		doc   PdfDoc
		input TableData
	}
	tests := []struct {
		name string
		args args
		want MockPdf
	}{
		{"Create title with date",
			args{
				doc: &pdf,
				input: TableData{
					Columns: []TableColumn{
						{
							Header:    "Zeile 1",
							Alignment: gopdf.Left,
							Width:     10,
							Rows: []string{
								"A", "b", "C",
							},
						},
						{
							Header:    "Zeile 2",
							Alignment: gopdf.Right,
							Width:     30,
							Rows: []string{
								"X", "y", "Z",
							},
						},
					},
					LastRowBold: true,
				},
			},
			MockPdf{
				currentX:        30,
				currentY:        148,
				currentFontSize: 11,
				currentStyle:    "default",
				cells: map[float64]map[float64]Cell{
					20: {
						130: {
							rect: gopdf.Rect{
								W: 10,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Left,
								Border: gopdf.Bottom,
							},
							text: "Zeile 1",
						},
						136: {
							rect: gopdf.Rect{
								W: 10,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Left,
								Border: 0,
							},
							text: "A",
						},
						142: {
							rect: gopdf.Rect{
								W: 10,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Left,
								Border: 0,
							},
							text: "b",
						},
						148: {
							rect: gopdf.Rect{
								W: 10,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Left,
								Border: 0,
							},
							text: "C",
						},
					},
					30: {
						130: {
							rect: gopdf.Rect{
								W: 30,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Right,
								Border: gopdf.Bottom,
							},
							text: "Zeile 2",
						},
						136: {
							rect: gopdf.Rect{
								W: 30,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Right,
								Border: 0,
							},
							text: "X",
						},
						142: {
							rect: gopdf.Rect{
								W: 30,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Right,
								Border: 0,
							},
							text: "y",
						},
						148: {
							rect: gopdf.Rect{
								W: 30,
								H: 4.5,
							},
							options: gopdf.CellOption{
								Align:  gopdf.Right,
								Border: 0,
							},
							text: "Z",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddTable(tt.args.doc, tt.args.input)
			if !reflect.DeepEqual(pdf, tt.want) {
				t.Errorf("AddTable() = \n%v, \nwant \n%v", pdf, tt.want)
			}
		})
	}
}

func TestAddText(t *testing.T) {

	pdf := MockPdf{
		textData:        make(map[float64]map[float64]TextData),
		currentStyle:    "default",
		currentFontSize: 11,
	}
	type args struct {
		doc           PdfDoc
		multilineText string
	}
	tests := []struct {
		name string
		args args
		want MockPdf
	}{
		{
			"Add Multiline Text",
			args{
				doc:           &pdf,
				multilineText: "Input Text multiline \n next line",
			},
			MockPdf{
				currentX:        0,
				currentY:        0,
				currentFontSize: 11,
				currentStyle:    "default",
				textData: map[float64]map[float64]TextData{
					20: {
						110: TextData{
							text:  "Input Text multiline \n next line",
							style: "default",
							size:  11,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddText(tt.args.doc, tt.args.multilineText)
			if !reflect.DeepEqual(pdf, tt.want) {
				t.Errorf("AddMultilinetext() = \n%v, \nwant \n%v", pdf, tt.want)
			}
		})
	}
}
