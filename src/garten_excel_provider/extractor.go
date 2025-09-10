package gartenexcelprovider

import (
	"fmt"
	"log"
	"strconv"

	"github.com/toky03/qr-invoice/document"
	"github.com/xuri/excelize/v2"
)

func ReadPaechterData(workbook *excelize.File) []PaechterData {

	var debtors []PaechterData = make([]PaechterData, 0)

	rows, err := workbook.GetRows("Mitgliederliste")
	if err != nil {
		log.Printf("could not read Sheet %s", err)
	}

	for i, row := range rows {

		// skip first row as this is the header Row
		if i == 0 {
			continue

		}

		if len(row) < 10 {
			continue
		}

		var isVorstand bool

		if len(row) < 13 {
			isVorstand = false
		} else {
			isVorstand = row[12] == "J"
		}

		zip := row[4]
		lang := row[11]
		if zip == "" || lang == "" {
			continue
		}

		lastName := row[1]
		firstName := row[2]
		address := row[3]
		city := row[5]

		parzelle := row[0]

		are, err := strconv.ParseFloat(row[7], 32)
		if err != nil {
			cellName, _ := excelize.CoordinatesToCellName(7+1, i+1)
			log.Printf(
				"could not convert %s from sheet Mitgliederliste on Cell %s to number\n",
				row[7],
				cellName,
			)
			continue
		}

		var language string

		if lang == "F" {
			language = "fr"
		} else if lang == "D" {
			language = "de"
		} else if lang == "I" {
			language = "it"
		} else {
			language = "en"
		}

		debtor := PaechterData{
			Parzelle:   parzelle,
			Are:        float32(are),
			IsVorstand: isVorstand,
			Language:   language,
			LastName:   lastName,
			Debtor: document.Debtor{
				Name:    fmt.Sprintf("%s %s", firstName, lastName),
				Address: address,
				Zip:     zip,
				City:    city,
				Country: "CH",
			},
		}

		debtors = append(debtors, debtor)
	}

	return debtors

}

func ReadVariableData(workbook *excelize.File) (VariableData, error) {

	pachtzinsText, err := extractTranslation(workbook, 1)
	wasserBezugText, err := extractTranslation(workbook, 2)
	abonementGfText, err := extractTranslation(workbook, 3)
	stromText, err := extractTranslation(workbook, 4)
	versicherungText, err := extractTranslation(workbook, 5)
	mitgliederbeitragText, err := extractTranslation(workbook, 6)
	reparaturfondsText, err := extractTranslation(workbook, 7)
	verwaltungskostenText, err := extractTranslation(workbook, 8)

	pachtzins, err := extractCellValueAsFloat32(workbook, "A3")
	wasserbezug, err := extractCellValueAsFloat32(workbook, "B3")
	gfAbonement, err := extractCellValueAsFloat32(workbook, "C3")
	strom, err := extractCellValueAsFloat32(workbook, "D3")
	versicherung, err := extractCellValueAsFloat32(workbook, "E3")
	mitgliederbeitrag, err := extractCellValueAsFloat32(workbook, "F3")
	reparaturfonds, err := extractCellValueAsFloat32(workbook, "G3")
	verwaltungskosten, err := extractCellValueAsFloat32(workbook, "H3")

	if err != nil {
		return VariableData{}, fmt.Errorf("could not read variable Data from Betraege %s", err)
	}

	return VariableData{
		TextPachtzins:         pachtzinsText,
		TextWasserbezug:       wasserBezugText,
		TextGfAbonement:       abonementGfText,
		TextStrom:             stromText,
		TextVersicherung:      versicherungText,
		TextMitgliederbeitrag: mitgliederbeitragText,
		TextReparaturFonds:    reparaturfondsText,
		TextVerwaltungskosten: verwaltungskostenText,
		Pachtzins:             pachtzins,
		Wasserbezug:           wasserbezug,
		GfAbonement:           gfAbonement,
		Strom:                 strom,
		Versicherung:          versicherung,
		Mitgliederbeitrag:     mitgliederbeitrag,
		Reparaturfonds:        reparaturfonds,
		Verwaltungskosten:     verwaltungskosten,
	}, nil

}

func ReadInvoiceDetails(workbook *excelize.File) (InvoiceDetails, error) {

	name, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"B1",
		excelize.Options{RawCellValue: true},
	)
	address, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"B2",
		excelize.Options{RawCellValue: true},
	)
	buildingNumber, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"B3",
		excelize.Options{RawCellValue: true},
	)
	zip, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"B4",
		excelize.Options{RawCellValue: true},
	)
	city, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"B5",
		excelize.Options{RawCellValue: true},
	)
	iban, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"B6",
		excelize.Options{RawCellValue: true},
	)
	ueberschriftDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C7",
		excelize.Options{RawCellValue: true},
	)
	ueberschriftFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D7",
		excelize.Options{RawCellValue: true},
	)

	anzahlDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C8",
		excelize.Options{RawCellValue: true},
	)
	anzahlFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D8",
		excelize.Options{RawCellValue: true},
	)

	einheitDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C9",
		excelize.Options{RawCellValue: true},
	)
	einheitFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D9",
		excelize.Options{RawCellValue: true},
	)

	bezeichnungDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C10",
		excelize.Options{RawCellValue: true},
	)
	bezeichnungFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D10",
		excelize.Options{RawCellValue: true},
	)

	preisDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C11",
		excelize.Options{RawCellValue: true},
	)
	preisFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D11",
		excelize.Options{RawCellValue: true},
	)

	betragDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C12",
		excelize.Options{RawCellValue: true},
	)
	betragFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D12",
		excelize.Options{RawCellValue: true},
	)

	arenDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C13",
		excelize.Options{RawCellValue: true},
	)
	arenFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D13",
		excelize.Options{RawCellValue: true},
	)

	jahreDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C14",
		excelize.Options{RawCellValue: true},
	)
	jahreFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D14",
		excelize.Options{RawCellValue: true},
	)

	zusatzDe, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"C15",
		excelize.Options{RawCellValue: true},
	)
	zusatzFr, err := workbook.GetCellValue(
		"Rechnungsdetails",
		"D15",
		excelize.Options{RawCellValue: true},
	)

	if err != nil {
		return InvoiceDetails{}, fmt.Errorf(
			"could not read invoice Data from Rechnungsdetails %s",
			err,
		)
	}

	if name == "" || address == "" || zip == "" || city == "" || iban == "" {
		return InvoiceDetails{}, fmt.Errorf("some required values are missing")
	}

	return InvoiceDetails{
		Creditor: document.Creditor{
			Name:           name,
			Address:        address,
			BuildingNumber: buildingNumber,
			Zip:            zip,
			City:           city,
			Country:        "CH",
			Account:        iban,
		},
		Ueberschrift: document.TranslatedText{
			De: ueberschriftDe,
			Fr: ueberschriftFr,
		},
		TabelleAnzahl: document.TranslatedText{
			De: anzahlDe,
			Fr: anzahlFr,
		},
		TabelleEinheit: document.TranslatedText{
			De: einheitDe,
			Fr: einheitFr,
		},
		TabelleBezeichnung: document.TranslatedText{
			De: bezeichnungDe,
			Fr: bezeichnungFr,
		},
		TabellePreis: document.TranslatedText{
			De: preisDe,
			Fr: preisFr,
		},
		TabelleBetrag: document.TranslatedText{
			De: betragDe,
			Fr: betragFr,
		},
		TabelleAaren: document.TranslatedText{
			De: arenDe,
			Fr: arenFr,
		},
		TabelleJahre: document.TranslatedText{
			De: jahreDe,
			Fr: jahreFr,
		},
		Zusatz: document.TranslatedText{
			De: zusatzDe,
			Fr: zusatzFr,
		},
	}, nil
}

func extractCellValueAsFloat32(workbook *excelize.File, cellName string) (float32, error) {
	strValue, err := workbook.GetCellValue(
		"Betraege",
		cellName,
		excelize.Options{RawCellValue: true},
	)
	if err != nil {
		return 0, err
	}

	floatValue, err := strconv.ParseFloat(strValue, 32)
	if err != nil {
		log.Printf("could not convert string %s into number from Betraege %s", strValue, cellName)
	}
	return float32(floatValue), err
}

func extractTranslation(workbook *excelize.File, columnIndex int) (document.TranslatedText, error) {
	cellNameDe, _ := excelize.CoordinatesToCellName(columnIndex, 1)
	cellNameFr, _ := excelize.CoordinatesToCellName(columnIndex, 2)
	de, err := workbook.GetCellValue("Betraege", cellNameDe, excelize.Options{RawCellValue: true})
	if err != nil || de == "" {
		return document.TranslatedText{}, fmt.Errorf(
			"could not extract value from Betraege %s %s",
			cellNameDe,
			err,
		)
	}
	fr, err := workbook.GetCellValue("Betraege", cellNameFr, excelize.Options{RawCellValue: true})
	if err != nil || fr == "" {
		return document.TranslatedText{}, fmt.Errorf(
			"could not extract value from Betraege %s %s",
			cellNameDe,
			err,
		)
	}

	return document.TranslatedText{
		De: de,
		Fr: fr,
	}, nil
}
