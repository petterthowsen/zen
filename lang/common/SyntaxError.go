package common

type SyntaxError struct {
	Message  string
	Location *SourceLocation
}

func NewSyntaxError(message string, location *SourceLocation) *SyntaxError {
	return &SyntaxError{
		Message:  message,
		Location: location,
	}
}

func (e *SyntaxError) Error() string {
	var str = e.Message + " at " + e.Location.String() + "\n"
	return str + e.Location.GetLineWithMarker()
}
