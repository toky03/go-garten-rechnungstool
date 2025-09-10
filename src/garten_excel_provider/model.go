package gartenexcelprovider

import "github.com/toky03/qr-invoice/document"

type PaechterData struct {
	Parzelle   string
	Are        float32
	IsVorstand bool
	Language   string
	LastName   string
	Debtor     document.Debtor
}

type VariableData struct {
	TextPachtzins         document.TranslatedText
	TextWasserbezug       document.TranslatedText
	TextGfAbonement       document.TranslatedText
	TextStrom             document.TranslatedText
	TextVersicherung      document.TranslatedText
	TextMitgliederbeitrag document.TranslatedText
	TextReparaturFonds    document.TranslatedText
	TextVerwaltungskosten document.TranslatedText

	Pachtzins         float32
	Wasserbezug       float32
	GfAbonement       float32
	Strom             float32
	Versicherung      float32
	Mitgliederbeitrag float32
	Reparaturfonds    float32
	Verwaltungskosten float32
}

type InvoiceDetails struct {
	Creditor           document.Creditor
	Ueberschrift       document.TranslatedText
	TabelleAnzahl      document.TranslatedText
	TabelleEinheit     document.TranslatedText
	TabelleBezeichnung document.TranslatedText
	TabellePreis       document.TranslatedText
	TabelleBetrag      document.TranslatedText
	TabelleAaren       document.TranslatedText
	TabelleJahre       document.TranslatedText
	Zusatz             document.TranslatedText
}

type CalculatedData struct {
	Pachtzins         float32
	Wasserbezug       float32
	GfAbonement       float32
	Strom             float32
	Versicherung      float32
	Mitgliederbeitrag float32
	Reparaturfonds    float32
	Verwaltungskosten float32
	Total             float32
	Are               float32
}
