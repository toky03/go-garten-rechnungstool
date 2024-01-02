package model

type Debtor struct {
	Name           string
	Address        string
	BuildingNumber string
	Zip            string
	City           string
	Country        string
}

type Creditor struct {
	Name           string
	Address        string
	BuildingNumber string
	Zip            string
	City           string
	Country        string
	Account        string
}

type DebtorData struct {
	Parzelle   string
	Are        float32
	IsVorstand bool
	Language   string
	LastName   string
	Debtor     Debtor
}

type VariableData struct {
	TextPachtzins         TranslatedText
	TextWasserbezug       TranslatedText
	TextGfAbonement       TranslatedText
	TextStrom             TranslatedText
	TextVersicherung      TranslatedText
	TextMitgliederbeitrag TranslatedText
	TextReparaturFonds    TranslatedText
	TextVerwaltungskosten TranslatedText

	Pachtzins         float32
	Wasserbezug       float32
	GfAbonement       float32
	Strom             float32
	Versicherung      float32
	Mitgliederbeitrag float32
	Reparaturfonds    float32
	Verwaltungskosten float32
}

type TranslatedText struct {
	De string
	Fr string
}

type InvoiceDetails struct {
	Creditor           Creditor
	Ueberschrift       TranslatedText
	TabelleAnzahl      TranslatedText
	TabelleEinheit     TranslatedText
	TabelleBezeichnung TranslatedText
	TabellePreis       TranslatedText
	TabelleBetrag      TranslatedText
	TabelleAaren       TranslatedText
	TabelleJahre       TranslatedText
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
