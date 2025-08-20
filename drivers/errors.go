package drivers

type Code string

const (
	CodeAlreadyExists      Code = "ALREADY_EXISTS"
	CodeNotFound           Code = "NOT_FOUND"
	CodePreconditionFailed Code = "PRECONDITION_FAILED"
	CodeInsufficientSpace  Code = "INSUFFICIENT_SPACE"
	CodeBusy               Code = "BUSY"
	CodeUnknown            Code = "UNKNOWN"
)

type DriverError struct {
	Code      Code
	Message   string
	Temporary bool
}

func (e *DriverError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Message
}

func NewError(code Code, msg string, temporary bool) *DriverError {
	return &DriverError{Code: code, Message: msg, Temporary: temporary}
}
