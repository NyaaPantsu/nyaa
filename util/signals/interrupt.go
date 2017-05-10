package signals

// handle interrupt signal, platform independent
func interrupted() {
	closeClosers()
	handleInterrupt()
}
