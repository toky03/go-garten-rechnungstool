package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestReadDebtorData(t *testing.T) {
	type args struct {
		workbook *excelize.File
	}
	tests := []struct {
		name string
		args args
		want []DebtorData
	}{
		{
			name: "empty data", args: args{
				workbook: createExampleFileWithSheet("Mitgliederliste", [][]string{
					{"Parzelle", "Nachname", "Vorname", "Adresse", "plz", "Ort", "", "Aare", "", "", "Sprache"},
				})},
			want: []DebtorData{}},
		{
			name: "single row with no value in last column",
			args: args{workbook: createExampleFileWithSheet("Mitgliederliste", [][]string{
				{"Parzelle", "Nachname", "Vorname", "Adresse", "plz", "Ort", "", "Aare", "", "", "Sprache"},
				{"1", "Simpson", "Homer", "Evergreen Terrace 742", "3011", "Bern", "", "2.50", "", "", "D"}})},
			want: []DebtorData{
				{
					Parzelle:   "1",
					Are:        2.5,
					IsVorstand: false,
					Language:   "de",
					LastName:   "Simpson",
					Debtor: Debtor{
						Name:    "Homer Simpson",
						Address: "Evergreen Terrace 742",
						Zip:     "3011",
						City:    "Bern",
						Country: "CH",
					},
				},
			}},
		{
			name: "multiple rows",
			args: args{workbook: createExampleFileWithSheet("Mitgliederliste", [][]string{
				{"Parzelle", "Nachname", "Vorname", "Adresse", "plz", "Ort", "", "Aare", "", "", "Sprache"},
				{"1", "Simpson", "Homer", "Evergreen Terrace 742", "3011", "Bern", "", "2.50", "", "", "D"},
				{"2", "Simpson", "Marge", "Evergreen Terrace 742", "3011", "Bern", "", "2.50", "", "", "D", "J"},
			})},
			want: []DebtorData{
				{
					Parzelle:   "1",
					Are:        2.5,
					IsVorstand: false,
					Language:   "de",
					LastName:   "Simpson",
					Debtor: Debtor{
						Name:    "Homer Simpson",
						Address: "Evergreen Terrace 742",
						Zip:     "3011",
						City:    "Bern",
						Country: "CH",
					},
				},
				{
					Parzelle:   "2",
					Are:        2.5,
					IsVorstand: true,
					Language:   "de",
					LastName:   "Simpson",
					Debtor: Debtor{
						Name:    "Marge Simpson",
						Address: "Evergreen Terrace 742",
						Zip:     "3011",
						City:    "Bern",
						Country: "CH",
					},
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadDebtorData(tt.args.workbook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDebtorData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadVariableData(t *testing.T) {
	type args struct {
		workbook *excelize.File
	}
	tests := []struct {
		name string
		args args
		want VariableData
	}{
		{
			name: "read data from second line", args: args{
				workbook: createExampleFileWithSheet("Betraege", [][]string{
					{"pachtzins", "wasserbezug", "GF Abonement", "Strom", "Versicherung", "Mitgliederbeitrag", "Reparaturfonds", "Verwaltungskosten"},
					{"Fr pachtzins", "Fr wasserbezug", "Fr GF Abonement", "Fr Strom", "Fr Versicherung", "Fr Mitgliederbeitrag", "Fr Reparaturfonds", "Fr Verwaltungskosten"},
					{"10", "20", "5.5", "8", "9", "50", "11.9", "22"},
				})},
			want: VariableData{
				TextPachtzins: TranslatedText{
					De: "pachtzins",
					Fr: "Fr pachtzins",
				},
				TextWasserbezug: TranslatedText{
					De: "wasserbezug",
					Fr: "Fr wasserbezug",
				},
				TextGfAbonement: TranslatedText{
					De: "GF Abonement",
					Fr: "Fr GF Abonement",
				},
				TextStrom: TranslatedText{
					De: "Strom",
					Fr: "Fr Strom",
				},
				TextVersicherung: TranslatedText{
					De: "Versicherung",
					Fr: "Fr Versicherung",
				},
				TextMitgliederbeitrag: TranslatedText{
					De: "Mitgliederbeitrag",
					Fr: "Fr Mitgliederbeitrag",
				},
				TextReparaturFonds: TranslatedText{
					De: "Reparaturfonds",
					Fr: "Fr Reparaturfonds",
				},
				TextVerwaltungskosten: TranslatedText{
					De: "Verwaltungskosten",
					Fr: "Fr Verwaltungskosten",
				},
				Pachtzins:         10,
				Wasserbezug:       20,
				GfAbonement:       5.5,
				Strom:             8,
				Versicherung:      9,
				Mitgliederbeitrag: 50,
				Reparaturfonds:    11.9,
				Verwaltungskosten: 22,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ReadVariableData(tt.args.workbook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadVariableData() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}

func TestReadVariableDataWithError(t *testing.T) {
	type args struct {
		workbook *excelize.File
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "empty data should panic", args: args{
				workbook: createExampleFileWithSheet("Betraege", [][]string{})},
			want: fmt.Errorf("could not read variable Data from Betraege %s", "strconv.ParseFloat: parsing \"\": invalid syntax")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := ReadVariableData(tt.args.workbook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadVariableData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadInvoiceDetailsWithError(t *testing.T) {
	type args struct {
		workbook *excelize.File
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "sheet does not exist", args: args{
				workbook: createExampleFileWithSheet("Betraege", [][]string{})},
			want: fmt.Errorf("could not read invoice Data from Rechnungsdetails sheet Rechnungsdetails does not exist"),
		},
		{
			name: "empty data return error", args: args{
				workbook: createExampleFileWithSheet("Rechnungsdetails", [][]string{
					{"", ""},
					{"", ""},
					{"", ""},
					{"", ""},
					{"", ""},
					{"", ""},
					{"", ""},
				})},
			want: fmt.Errorf("some required values are missing"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := ReadInvoiceDetails(tt.args.workbook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadInvoiceDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadInvoiceDetails(t *testing.T) {
	type args struct {
		workbook *excelize.File
	}
	tests := []struct {
		name string
		args args
		want InvoiceDetails
	}{
		{
			name: "extract value successfully",
			args: args{
				workbook: createExampleFileWithSheet("Rechnungsdetails", [][]string{
					{"Name", "Tester Test"},
					{"Adresse", "Testerweg"},
					{"Adresse Nummer", "42"},
					{"Postleitzahl", "3011"},
					{"Stadt", "Bern"},
					{"iban Nummer", "CH04 0077 7001 7282 6000 2"},
					{"Ueberschrift", "", "Rechnung 2024", "Facture 2024"},
					{"Anzahl", "", "Anzahl", "Nombre"},
					{"Einheit", "", "Einheit", "Unité"},
					{"Bezeichnung", "", "Bezeichnung", "Dénomination"},
					{"Preis", "", "Preis", "Prix"},
					{"Betrag", "", "Betrag", "Montant"},
					{"Aren", "", "Aren", "Are"},
					{"Jahre", "", "Jahr", "Année"},
				})},
			want: InvoiceDetails{
				Creditor: Creditor{
					Name:           "Tester Test",
					Address:        "Testerweg",
					BuildingNumber: "42",
					Zip:            "3011",
					City:           "Bern",
					Country:        "CH",
					Account:        "CH04 0077 7001 7282 6000 2",
				},
				Ueberschrift: TranslatedText{
					De: "Rechnung 2024",
					Fr: "Facture 2024",
				},
				TabelleAnzahl: TranslatedText{
					De: "Anzahl",
					Fr: "Nombre",
				},
				TabelleEinheit: TranslatedText{
					De: "Einheit",
					Fr: "Unité",
				},
				TabelleBezeichnung: TranslatedText{
					De: "Bezeichnung",
					Fr: "Dénomination",
				},
				TabellePreis: TranslatedText{
					De: "Preis",
					Fr: "Prix",
				},
				TabelleBetrag: TranslatedText{
					De: "Betrag",
					Fr: "Montant",
				},
				TabelleAaren: TranslatedText{
					De: "Aren",
					Fr: "Are",
				},
				TabelleJahre: TranslatedText{
					De: "Jahr",
					Fr: "Année",
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ReadInvoiceDetails(tt.args.workbook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadInvoiceDetails() = \n%v\nwant \n%v", got, tt.want)
			}
		})
	}
}

func createExampleFileWithSheet(sheetName string, content [][]string) *excelize.File {
	f := excelize.NewFile()
	f.NewSheet(sheetName)
	for rowIndex, row := range content {
		for colIndex, colValue := range row {
			cellName, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			f.SetCellValue(sheetName, cellName, colValue)
		}
	}
	return f

}
