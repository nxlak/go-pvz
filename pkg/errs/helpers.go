package errs

import (
	"encoding/json"

	"github.com/go-faster/jx"
	order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
)

func ToAPIErrorFields(fields []Field) order_v1.OptAppErrorFields {
	if len(fields) == 0 {
		return order_v1.OptAppErrorFields{}
	}

	apiFields := order_v1.AppErrorFields{} // map[string]jx.Raw

	for _, f := range fields {
		b, err := json.Marshal(f.Value)
		if err != nil {
			// можно залогировать и пропустить
			continue
		}
		apiFields[f.Key] = jx.Raw(b)
	}

	return order_v1.NewOptAppErrorFields(apiFields)
}

func parseFields(args ...any) []Field {
	var result []Field

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		result = append(result, Field{Key: key, Value: args[i+1]})
	}

	return result
}
