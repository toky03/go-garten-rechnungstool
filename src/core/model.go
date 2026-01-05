package core

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
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

	// Calculate total for each row
	totalValues := make([]float64, t.Dataframe.Nrow())
	for i := 0; i < t.Dataframe.Nrow(); i++ {
		multVal := multiplicator.Elem(i).Float()
		priceVal := price.Elem(i).Float()
		totalValues[i] = multVal * priceVal
	}

	// Create new series with calculated totals
	total := series.New(totalValues, series.Float, t.TotalColumnName)
	t.Dataframe = t.Dataframe.Mutate(total)

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
