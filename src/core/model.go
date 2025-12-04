package core

import (
	"github.com/go-gota/gota/dataframe"
)

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

type TableData struct {
	PriceColumnName         string
	MultiplicatorColumnName string
	TotalColumnName         string
	Dataframe               dataframe.DataFrame
}

func (t TableData) ToInvoiceTable() InvoiceTable {

	multiplicator := t.Dataframe.Col(t.MultiplicatorColumnName)
	price := t.Dataframe.Col(t.PriceColumnName)

	t.Dataframe.(func(s dataframe.Series) dataframe.Series) {
		if s.Name() == t.TotalColumnName {
			return dataframe.NewSeriesFloat64(t.TotalColumnName, nil)
		}
		return s
	}, t.TotalColumnName)
	t.Dataframe = t.Dataframe.Mutate(total.Name(t.TotalColumnName))

	headers := t.Dataframe.Names()
	rows := make([][]string, t.Dataframe.Nrow())
	for i := 0; i < t.Dataframe.Nrow(); i++ {
		row := make([]string, t.Dataframe.Ncol())
		for j := 0; j < t.Dataframe.Ncol(); j++ {
			row[j] = t.Dataframe.Elem(i, j).String()
		}
		rows[i] = row
	}

	return InvoiceTable{
		Headers: headers,
		Rows:    rows,
	}
}
