package apiv2

import "github.com/rackcorpcloud/rackcorp-api-go/internal"

type RegionID int

func (id *RegionID) UnmarshalJSON(data []byte) error {
	return internal.UnmarshalJSONInt(id, data)
}
