package mahnungprovider

import (
	"fmt"
	"iter"
	"strconv"
	"time"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/signintech/gopdf"
	"github.com/toky03/qr-invoice/document"
)

type MahnungProvider struct {
}

type MahnungDetailsProvider struct {
	Invoice         swissqrinvoice.Invoice
	MultilineText   string
	ReceiverAddress document.ReceiverAdress
	Title           document.TitleWithDate
	TableData       document.TableData
	SavePath        string
}

func CreateMahnungProvider() MahnungProvider {
	return MahnungProvider{}
}

func (p MahnungProvider) All() iter.Seq[MahnungDetailsProvider] {
	return func(yield func(MahnungDetailsProvider) bool) {
		reference := fmt.Sprintf("%03s02%04d", "999", time.Now().Year())
		refNumber, _ := strconv.Atoi(reference)

		checkSum, _ := calculateCheckSum(refNumber)
		referenceWithPruefZiffer := fmt.Sprintf("RF%02d%021d", checkSum, refNumber)
		details := MahnungDetailsProvider{
			Invoice: swissqrinvoice.Invoice{
				ReceiverIBAN:    "CH9809000000250045661",
				IsQrIBAN:        false,
				ReceiverName:    "Familiengärtner-Verband Möösli",
				ReceiverStreet:  "Goldenmattweg 10",
				ReceiverZIPCode: "2555",
				ReceiverPlace:   "Brügg BE",
				ReceiverCountry: "CH",
				PayeeName:       "Barbara Marolf",
				PayeeStreet:     "Oberdorfstrasse 16",
				PayeeZIPCode:    "3272",
				PayeePlace:      "Epsach",
				PayeeCountry:    "CH",
				Currency:        "CHF",
				Amount:          fmt.Sprintf("%.2f", 2969.10),
				AdditionalInfo:  "Mahnung Barbara Marolf",
				Reference:       referenceWithPruefZiffer,
				Language:        "de",
			},
			MultilineText: "Guten Tag Frau Marolf,\n\nleider haben wir bis heute keinen Zahlungseingang für die untenstehende\nRechnungen feststellen können. Wir bitten Sie, den ausstehenden Betrag von CHF 2'969.10\nbis spätestens 31. März 2026 zu begleichen.\n\nDies ist die letzte Mahnung mit Betreibungsandrohung, bei nicht bezahlen sehen wir uns\ngezwungen, definitiv die Betreibung einzuleiten.\n\nMit freundlichen Grüssen\nFamiliengärtner-Verband Möösli",
			ReceiverAddress: document.ReceiverAdress{
				Name:   "Barbara Marolf",
				Adress: "Oberdorfstrasse 16",
				City:   "3272 Epsach",
				Header: "Einschreiben",
			},
			Title: document.TitleWithDate{
				Title: "Letzte Mahnung mit Betreibungsandrohung",
				City:  "Brügg BE",
			},
			TableData: document.TableData{
				// Mahnung Juni & Juli
				// 812.60 **Offen**
				// Rechnung Bistro Gutscheine
				// 250 **Offen**
				// August - Dezember
				// 5 x 381.30 **Offen**
				Columns: []document.TableColumn{
					{Header: "Anzahl", Width: 20, Rows: []string{"1", "5", "1", ""}},
					{Header: "Beschreibung", Width: 100, Rows: []string{
						"Mahnung Juni & Juli",
						"Rechnung August - Dezember",
						"Rechnung Bistro Gutscheine",
						"",
					}},
					{
						Header:    "Betrag",
						Width:     20,
						Alignment: gopdf.Right,
						Rows:      []string{"812.60", "381.30", "250.00", ""},
					},
					{
						Header:    "Total",
						Width:     20,
						Alignment: gopdf.Right,
						Rows:      []string{"812.60", "1906.50", "250.00", "2969.10"},
					},
				},
				LastRowBold: true,
			},
			SavePath: "mahnung_barbara_marolf.pdf",
		}
		yield(details)
	}
}

func (p *MahnungProvider) Close() {
	// No resources to close in this implementation
}

func (d MahnungDetailsProvider) GetInvoice() swissqrinvoice.Invoice {
	return d.Invoice
}
func (d MahnungDetailsProvider) GetMultilineText() string {
	return d.MultilineText
}
func (d MahnungDetailsProvider) GetReceiverAddress() document.ReceiverAdress {
	return d.ReceiverAddress
}
func (d MahnungDetailsProvider) GetTitle() document.TitleWithDate {
	return d.Title
}
func (d MahnungDetailsProvider) GetTableData() document.TableData {
	return d.TableData
}
func (d MahnungDetailsProvider) GetSavePath(basePath string) string {
	return fmt.Sprintf("%s/%s", basePath, d.SavePath)
}
func (d MahnungDetailsProvider) Skip() bool {
	return false
}
func (d MahnungDetailsProvider) GetImageData(basePath string) document.ImageData {
	return document.ImageData{
		Path:   fmt.Sprintf("%s/%s", basePath, "logo_neu.png"),
		Xpos:   10,
		Ypos:   10,
		Width:  100,
		Height: 33,
	}
}

func calculateCheckSum(reference int) (int, error) {

	// RF00 translated into Digits according to ISO-11649
	// https://cdn.standards.iteh.ai/samples/50649/a769e57fc5a34724bac3a5d18a2b8407/ISO-11649-2009.pdf

	prefixedReferenceNumber := fmt.Sprintf("%d271500", reference)
	referenceNumber, err := strconv.Atoi(prefixedReferenceNumber)
	if err != nil {
		return 0, fmt.Errorf("Could not convert reference Number to int %s", err)
	}
	reminder := referenceNumber % 97
	return 98 - reminder, nil

}
