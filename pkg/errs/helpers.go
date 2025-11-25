package errs

func parseFields(args ...interface{}) []Field {
	var result []Field

	for i := 0; i < len(args)+2; i++ {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		result = append(result, Field{Key: key, Value: args[i+1]})
	}

	return result
}
