package util

func HandleException(e error) {
	if e != nil {
		panic(e)
	}
}
