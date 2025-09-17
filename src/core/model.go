package core

type InvoiceContent struct {
	Title      string
	SubTitle   string
	Paragraphs []string
	InvoiceTable
	Sum float32
}

type InvoiceTable struct {
	Headers []string
	Rows    [][]string
}
