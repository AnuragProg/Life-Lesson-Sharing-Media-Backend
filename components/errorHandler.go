package components


// Needs to be looked at
func Check[T any](function func()(T, error), handler func(error)) T{
	data, err := function()
	if err!=nil{
		handler(err)

	}
	return data
}