package gartenexcelprovider

import (
	"reflect"
	"testing"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/signintech/gopdf"
	"github.com/toky03/qr-invoice/document"
)

func TestVariableData_ToCalculatedTableData(t *testing.T) {
	type fields struct {
		Pachtzins         float32
		Wasserbezug       float32
		GfAbonement       float32
		Strom             float32
		Versicherung      float32
		Mitgliederbeitrag float32
		Reparaturfonds    float32
		Verwaltungskosten float32
	}
	type args struct {
		debtor PaechterData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CalculatedData
	}{
		{
			name: "calculate sum for total for non vorstandmitglied",
			fields: fields{
				Pachtzins:         10,
				Wasserbezug:       2,
				GfAbonement:       20,
				Strom:             1,
				Versicherung:      2,
				Mitgliederbeitrag: 50,
				Reparaturfonds:    8,
				Verwaltungskosten: 9,
			},
			args: args{debtor: PaechterData{
				Are:        1.5,
				IsVorstand: false,
			}},
			want: CalculatedData{
				Pachtzins:         15,
				Wasserbezug:       3,
				GfAbonement:       20,
				Strom:             1,
				Versicherung:      2,
				Mitgliederbeitrag: 50,
				Reparaturfonds:    8,
				Verwaltungskosten: 9,
				Total:             108,
				Are:               1.5,
			},
		},
		{
			name: "calculate sum for total for vorstandmitglied",
			fields: fields{
				Pachtzins:         10,
				Wasserbezug:       2,
				GfAbonement:       20,
				Strom:             1,
				Versicherung:      2,
				Mitgliederbeitrag: 50,
				Reparaturfonds:    8,
				Verwaltungskosten: 9,
			},
			args: args{debtor: PaechterData{
				Are:        1.5,
				IsVorstand: true,
			}},
			want: CalculatedData{
				Pachtzins:         15,
				Wasserbezug:       3,
				GfAbonement:       20,
				Strom:             1,
				Versicherung:      2,
				Mitgliederbeitrag: 0,
				Reparaturfonds:    8,
				Verwaltungskosten: 9,
				Total:             58,
				Are:               1.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VariableData{
				Pachtzins:         tt.fields.Pachtzins,
				Wasserbezug:       tt.fields.Wasserbezug,
				GfAbonement:       tt.fields.GfAbonement,
				Strom:             tt.fields.Strom,
				Versicherung:      tt.fields.Versicherung,
				Mitgliederbeitrag: tt.fields.Mitgliederbeitrag,
				Reparaturfonds:    tt.fields.Reparaturfonds,
				Verwaltungskosten: tt.fields.Verwaltungskosten,
			}
			if got := v.ToCalculatedTableData(tt.args.debtor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VariableData.ToCalculatedTableData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvoiceDetails_ToInvoiceDetails(t *testing.T) {
	type fields struct {
		Creditor     document.Creditor
		Ueberschrift document.TranslatedText
	}
	type args struct {
		debtor         PaechterData
		calculatedData CalculatedData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   swissqrinvoice.Invoice
	}{
		{
			name: "transform to invoice Details",
			fields: fields{
				Creditor: document.Creditor{
					Name:    "Toky Tok",
					Address: "Evergreen 10",
					City:    "Biel",
					Account: "CH12 1234 5678 9101 12",
					Country: "CH",
					Zip:     "1233",
				},
				Ueberschrift: document.TranslatedText{
					De: "Ueberschrift DE",
				},
			},
			args: args{
				debtor: PaechterData{
					Debtor: document.Debtor{
						Name:    "Rechnungsempfaenger Name",
						Address: "Rechnungsempfaenger Adresse",
						Zip:     "3333",
						City:    "Bern",
						Country: "CH",
					},
					Parzelle: "99",
					Language: "de",
					LastName: "Nachname",
				},
				calculatedData: CalculatedData{
					Pachtzins: 12,
					Total:     200.1,
				},
			},
			want: swissqrinvoice.Invoice{
				ReceiverIBAN:    "CH12 1234 5678 9101 12",
				IsQrIBAN:        false,
				ReceiverName:    "Toky Tok",
				ReceiverStreet:  "Evergreen 10",
				ReceiverZIPCode: "1233",
				ReceiverPlace:   "Biel",
				ReceiverCountry: "CH",
				PayeeName:       "Rechnungsempfaenger Name",
				PayeeStreet:     "Rechnungsempfaenger Adresse",
				PayeeZIPCode:    "3333",
				Amount:          "200.10",
				AdditionalInfo:  "Parzelle 99",
				Language:        "de",
				Currency:        "CHF",
				PayeePlace:      "Bern",
				PayeeCountry:    "CH",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := InvoiceDetails{
				Creditor:     tt.fields.Creditor,
				Ueberschrift: tt.fields.Ueberschrift,
			}
			if got := i.ToInvoiceDetails(tt.args.debtor, tt.args.calculatedData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvoiceDetails.ToInvoiceDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebtorData_ToReceiverAdress(t *testing.T) {
	type fields struct {
		Parzelle string
		Debtor   document.Debtor
		Language string
	}
	tests := []struct {
		name   string
		fields fields
		want   document.ReceiverAdress
	}{
		{"debtor with all data german", fields{
			Parzelle: "33",
			Debtor: document.Debtor{
				Name:    "Vor und Nachname",
				Address: "Strasse xy",
				Zip:     "3011",
				City:    "Bern",
			},
			Language: "de",
		},
			document.ReceiverAdress{
				Header: "Parzelle 33",
				Name:   "Vor und Nachname",
				Adress: "Strasse xy",
				City:   "3011 Bern",
			},
		},
		{"debtor with all data french", fields{
			Parzelle: "33",
			Debtor: document.Debtor{
				Name:    "Vor und Nachname",
				Address: "Strasse xy",
				Zip:     "3011",
				City:    "Bern",
			},
			Language: "fr",
		},
			document.ReceiverAdress{
				Header: "Parcelle 33",
				Name:   "Vor und Nachname",
				Adress: "Strasse xy",
				City:   "3011 Bern",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debtor := PaechterData{
				Parzelle: tt.fields.Parzelle,
				Debtor:   tt.fields.Debtor,
				Language: tt.fields.Language,
			}
			if got := debtor.ToReceiverAdress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebtorData.ToReceiverAdress() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestInvoiceDetails_ToTitle(t *testing.T) {
	type fields struct {
		Creditor     document.Creditor
		Ueberschrift document.TranslatedText
	}
	type args struct {
		language string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   document.TitleWithDate
	}{
		{
			"Create Title German",
			fields{
				Creditor: document.Creditor{
					City: "Biel",
				},
				Ueberschrift: document.TranslatedText{
					De: "Ueberschrift abc",
					Fr: "Title xyz",
				},
			},
			args{language: "de"},
			document.TitleWithDate{
				Title: "Ueberschrift abc",
				City:  "Biel"},
		},
		{
			"Create Title French",
			fields{
				Creditor: document.Creditor{
					City: "Bienne",
				},
				Ueberschrift: document.TranslatedText{
					De: "Ueberschrift abc",
					Fr: "Title xyz",
				},
			},
			args{language: "fr"},
			document.TitleWithDate{
				Title: "Title xyz",
				City:  "Bienne"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoiceDetails := InvoiceDetails{
				Creditor:     tt.fields.Creditor,
				Ueberschrift: tt.fields.Ueberschrift,
			}
			if got := invoiceDetails.ToTitle(tt.args.language); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvoiceDetails.ToTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvoiceDetails_ToTableData(t *testing.T) {
	type fields struct {
		TabelleAnzahl      document.TranslatedText
		TabelleEinheit     document.TranslatedText
		TabelleBezeichnung document.TranslatedText
		TabellePreis       document.TranslatedText
		TabelleBetrag      document.TranslatedText
		TabelleAaren       document.TranslatedText
		TabelleJahre       document.TranslatedText
		Zusatz             document.TranslatedText
	}
	type args struct {
		language       string
		debtorData     PaechterData
		variableData   VariableData
		calculatedData CalculatedData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   document.TableData
	}{
		{
			"Create table german",
			fields{
				TabelleAnzahl: document.TranslatedText{
					De: "Anzahl",
					Fr: "Nombre",
				},
				TabelleEinheit: document.TranslatedText{
					De: "Einheit",
					Fr: "Unité",
				},
				TabelleBezeichnung: document.TranslatedText{
					De: "Bezeichnung",
					Fr: "Description",
				},
				TabellePreis: document.TranslatedText{
					De: "Preis",
					Fr: "Prix",
				},
				TabelleBetrag: document.TranslatedText{
					De: "Betrag",
					Fr: "Montant",
				},
				TabelleAaren: document.TranslatedText{
					De: "Aren",
					Fr: "Ares",
				},
				TabelleJahre: document.TranslatedText{
					De: "Jahre",
					Fr: "Années",
				},
			},
			args{
				language: "de",
				debtorData: PaechterData{
					Are: 2.2123,
				},
				variableData: VariableData{
					TextPachtzins: document.TranslatedText{
						De: "Pachtzins",
						Fr: "Fr: Pachtzins",
					},
					TextWasserbezug: document.TranslatedText{
						De: "Wasserbezug",
						Fr: "Fr: Wasserbezug",
					},
					TextGfAbonement: document.TranslatedText{
						De: "GF Abonement",
						Fr: "Fr: GF Abonement",
					},
					TextStrom: document.TranslatedText{
						De: "Strom",
						Fr: "Fr: Strom",
					},
					TextVersicherung: document.TranslatedText{
						De: "Versicherung",
						Fr: "Fr: Versicherung",
					},
					TextMitgliederbeitrag: document.TranslatedText{
						De: "Mitgliederbeitrag",
						Fr: "Fr: Mitgliederbeitrag",
					},
					TextReparaturFonds: document.TranslatedText{
						De: "Reparaturfonds",
						Fr: "Fr: Reparaturfonds",
					},
					TextVerwaltungskosten: document.TranslatedText{
						De: "Verwaltungskosten",
						Fr: "Fr: Verwaltungskosten",
					},
					Pachtzins:         5,
					Wasserbezug:       1,
					GfAbonement:       2,
					Strom:             3,
					Versicherung:      4,
					Mitgliederbeitrag: 6,
					Reparaturfonds:    7,
					Verwaltungskosten: 8,
				},
				calculatedData: CalculatedData{
					Pachtzins:         10,
					Wasserbezug:       2,
					GfAbonement:       4,
					Strom:             6,
					Versicherung:      8,
					Mitgliederbeitrag: 12,
					Reparaturfonds:    14,
					Verwaltungskosten: 16,
					Total:             200,
				},
			},
			document.TableData{
				Columns: []document.TableColumn{
					{
						Header:    "Anzahl",
						Alignment: gopdf.Left,
						Width:     35,
						Rows:      []string{"2.21", "2.21", "1", "1", "1", "1", "1", "1", "Total"},
					},
					{
						Header:    "Einheit",
						Alignment: gopdf.Left,
						Width:     30,
						Rows:      []string{"Aren", "Aren", "Jahre", "Jahre", "Jahre", "Jahre", "Jahre", "Jahre", ""},
					},
					{
						Header:    "Bezeichnung",
						Alignment: gopdf.Left,
						Width:     45,
						Rows:      []string{"Pachtzins", "Wasserbezug", "GF Abonement", "Strom", "Versicherung", "Mitgliederbeitrag", "Reparaturfonds", "Verwaltungskosten", ""},
					},
					{
						Header:    "Preis",
						Alignment: gopdf.Right,
						Width:     30,
						Rows:      []string{"CHF 5.00", "CHF 1.00", "CHF 2.00", "CHF 3.00", "CHF 4.00", "CHF 6.00", "CHF 7.00", "CHF 8.00", ""},
					},
					{
						Header:    "Betrag",
						Alignment: gopdf.Right,
						Width:     30,
						Rows:      []string{"CHF 10.00", "CHF 2.00", "CHF 4.00", "CHF 6.00", "CHF 8.00", "CHF 12.00", "CHF 14.00", "CHF 16.00", "CHF 200.00"},
					},
				},
				LastRowBold: true,
			},
		},
		{
			"Create table french",
			fields{
				TabelleAnzahl: document.TranslatedText{
					De: "Anzahl",
					Fr: "Nombre",
				},
				TabelleEinheit: document.TranslatedText{
					De: "Einheit",
					Fr: "Unité",
				},
				TabelleBezeichnung: document.TranslatedText{
					De: "Bezeichnung",
					Fr: "Description",
				},
				TabellePreis: document.TranslatedText{
					De: "Preis",
					Fr: "Prix",
				},
				TabelleBetrag: document.TranslatedText{
					De: "Betrag",
					Fr: "Montant",
				},
				TabelleAaren: document.TranslatedText{
					De: "Aren",
					Fr: "Ares",
				},
				TabelleJahre: document.TranslatedText{
					De: "Jahre",
					Fr: "Années",
				},
				Zusatz: document.TranslatedText{
					De: "Zusatz",
					Fr: "FR. Zusatz",
				},
			},
			args{
				language: "fr",
				debtorData: PaechterData{
					Are: 2.2123,
				},
				variableData: VariableData{
					TextPachtzins: document.TranslatedText{
						De: "Pachtzins",
						Fr: "Fr: Pachtzins",
					},
					TextWasserbezug: document.TranslatedText{
						De: "Wasserbezug",
						Fr: "Fr: Wasserbezug",
					},
					TextGfAbonement: document.TranslatedText{
						De: "GF Abonement",
						Fr: "Fr: GF Abonement",
					},
					TextStrom: document.TranslatedText{
						De: "Strom",
						Fr: "Fr: Strom",
					},
					TextVersicherung: document.TranslatedText{
						De: "Versicherung",
						Fr: "Fr: Versicherung",
					},
					TextMitgliederbeitrag: document.TranslatedText{
						De: "Mitgliederbeitrag",
						Fr: "Fr: Mitgliederbeitrag",
					},
					TextReparaturFonds: document.TranslatedText{
						De: "Reparaturfonds",
						Fr: "Fr: Reparaturfonds",
					},
					TextVerwaltungskosten: document.TranslatedText{
						De: "Verwaltungskosten",
						Fr: "Fr: Verwaltungskosten",
					},
					Pachtzins:         5,
					Wasserbezug:       1,
					GfAbonement:       2,
					Strom:             3,
					Versicherung:      4,
					Mitgliederbeitrag: 6,
					Reparaturfonds:    7,
					Verwaltungskosten: 8,
				},
				calculatedData: CalculatedData{
					Pachtzins:         10,
					Wasserbezug:       2,
					GfAbonement:       4,
					Strom:             6,
					Versicherung:      8,
					Mitgliederbeitrag: 0,
					Reparaturfonds:    14,
					Verwaltungskosten: 16,
					Total:             200,
				},
			},
			document.TableData{
				Columns: []document.TableColumn{
					{
						Header:    "Nombre",
						Alignment: gopdf.Left,
						Width:     35,
						Rows:      []string{"2.21", "2.21", "1", "1", "1", "1", "1", "1", "Total"},
					},
					{
						Header:    "Unité",
						Alignment: gopdf.Left,
						Width:     30,
						Rows:      []string{"Ares", "Ares", "Années", "Années", "Années", "Années", "Années", "Années", ""},
					},
					{
						Header:    "Description",
						Alignment: gopdf.Left,
						Width:     45,
						Rows:      []string{"Fr: Pachtzins", "Fr: Wasserbezug", "Fr: GF Abonement", "Fr: Strom", "Fr: Versicherung", "Fr: Mitgliederbeitrag", "Fr: Reparaturfonds", "Fr: Verwaltungskosten", ""},
					},
					{
						Header:    "Prix",
						Alignment: gopdf.Right,
						Width:     30,
						Rows:      []string{"CHF 5.00", "CHF 1.00", "CHF 2.00", "CHF 3.00", "CHF 4.00", "CHF 6.00", "CHF 7.00", "CHF 8.00", ""},
					},
					{
						Header:    "Montant",
						Alignment: gopdf.Right,
						Width:     30,
						Rows:      []string{"CHF 10.00", "CHF 2.00", "CHF 4.00", "CHF 6.00", "CHF 8.00", "CHF     -", "CHF 14.00", "CHF 16.00", "CHF 200.00"},
					},
				},
				LastRowBold: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoiceDetails := InvoiceDetails{
				TabelleAnzahl:      tt.fields.TabelleAnzahl,
				TabelleEinheit:     tt.fields.TabelleEinheit,
				TabelleBezeichnung: tt.fields.TabelleBezeichnung,
				TabellePreis:       tt.fields.TabellePreis,
				TabelleBetrag:      tt.fields.TabelleBetrag,
				TabelleAaren:       tt.fields.TabelleAaren,
				TabelleJahre:       tt.fields.TabelleJahre,
				Zusatz:             tt.fields.Zusatz,
			}
			if got := invoiceDetails.ToTableData(tt.args.language, tt.args.debtorData, tt.args.variableData, tt.args.calculatedData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvoiceDetails.ToTableData() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
