package apiv2

import "github.com/rackcorpcloud/rackcorp-api-go/internal"

type CustomerID int

func (id *CustomerID) UnmarshalJSON(data []byte) error {
	return internal.UnmarshalJSONInt(id, data)
}
