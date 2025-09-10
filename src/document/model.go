package document

type ReceiverAdress struct {
	Header string
	Name   string
	Adress string
	City   string
}

type TitleWithDate struct {
	Title string
	City  string
}

type TableColumn struct {
	Header    string
	Alignment int
	Width     float64
	Rows      []string
}

type TableData struct {
	Columns     []TableColumn
	LastRowBold bool
}

type ImageData struct {
	Path   string
	Xpos   float64
	Ypos   float64
	Width  float64
	Height float64
}

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

type TranslatedText struct {
	De string
	Fr string
}
