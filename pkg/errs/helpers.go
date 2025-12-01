package errs

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
