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
