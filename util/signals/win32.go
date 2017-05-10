// +build win32

package signals

func Handle() {
	// windows has no sighup LOOOOL, this does nothing
	// TODO: Something about SIGHUP for Windows
}
