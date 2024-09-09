package code

const (
	SUCCESS        = 0
	REQUEST_FAILED = 100
)

func CodeToMessage(code uint) string {
	switch code {
	case SUCCESS:
		return "success"
	case REQUEST_FAILED:
		return "ERROR"
	default:
		return "UNKNOWN"
	}

}
