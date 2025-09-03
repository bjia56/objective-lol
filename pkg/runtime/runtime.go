package runtime

import "github.com/bjia56/objective-lol/pkg/types"

// ReturnValue is used to control return flow in function execution
type ReturnValue struct {
	Value types.Value
}

func (r ReturnValue) Error() string {
	return "return"
}

// Helper function to check if an error is a return value
func IsReturnValue(err error) bool {
	_, ok := err.(ReturnValue)
	return ok
}

// Helper function to extract return value from error
func GetReturnValue(err error) types.Value {
	if retVal, ok := err.(ReturnValue); ok {
		return retVal.Value
	}
	return types.NOTHIN
}

// Exception represents a thrown exception with string message
type Exception struct {
	Message string
}

func (e Exception) Error() string {
	return e.Message
}

// Helper function to check if an error is an exception
func IsException(err error) bool {
	_, ok := err.(Exception)
	return ok
}

// Helper function to extract exception message from error
func GetExceptionMessage(err error) string {
	if exc, ok := err.(Exception); ok {
		return exc.Message
	}
	return err.Error()
}
