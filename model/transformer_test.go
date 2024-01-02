package model

import (
	"reflect"
	"testing"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
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
		debtor DebtorData
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
			args: args{debtor: DebtorData{
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
			args: args{debtor: DebtorData{
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
		Creditor     Creditor
		Ueberschrift TranslatedText
	}
	type args struct {
		debtor         DebtorData
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
				Creditor: Creditor{
					Name:    "Toky Tok",
					Address: "Evergreen 10",
					City:    "Biel",
					Account: "CH12 1234 5678 9101 12",
					Country: "CH",
					Zip:     "1233",
				},
				Ueberschrift: TranslatedText{
					De: "Ueberschrift DE",
				},
			},
			args: args{
				debtor: DebtorData{
					Debtor: Debtor{
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
