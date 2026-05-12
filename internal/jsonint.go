package internal

import "encoding/json"

// JSONInt is like json.Number for deserializing JSON numbers that are expected to be integers,
// but may be represented as strings in the JSON. Blank is deserialized as 0.
type JSONInt int

func (ji *JSONInt) UnmarshalJSON(data []byte) error {
	var tmp json.Number
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if len(tmp.String()) == 0 {
		*ji = 0
		return nil
	}
	num, err := tmp.Int64()
	if err != nil {
		return err
	}
	*ji = JSONInt(num)
	return nil
}

func (ji JSONInt) Int() int {
	return int(ji)
}

func JSONIntSliceInt(s []JSONInt) []int {
	ints := make([]int, len(s))
	for i, ji := range s {
		ints[i] = ji.Int()
	}
	return ints
}

func IntSliceJSONInt(s []int) []JSONInt {
	jsonInts := make([]JSONInt, len(s))
	for i, n := range s {
		jsonInts[i] = JSONInt(n)
	}
	return jsonInts
}

func UnmarshalJSONInt[T ~int](dst *T, data []byte) error {
	var tmp json.Number
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	num, err := tmp.Int64()
	if err != nil {
		return err
	}
	*dst = T(num)
	return nil
}
