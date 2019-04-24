package internal

import "context"

// String gets a string with the given key from the given context's values
func String(ctx context.Context, key interface{}) string {
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	return ""
}
