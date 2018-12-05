package api

func jsonMessage(msg string) []byte {
	return []byte(`{"message":"` + msg + `"}`)
}
