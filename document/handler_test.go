package document

import (
	"testing"

	gopdf_wrapper "github.com/72nd/gopdf-wrapper"
)

func TestAddAdressData(t *testing.T) {
	type args struct {
		doc      *gopdf_wrapper.Doc
		receiver ReceiverAdress
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddAdressData(tt.args.doc, tt.args.receiver)
		})
	}
}
