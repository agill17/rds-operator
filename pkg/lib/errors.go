package lib

// so that I can return these errors when reconcile is running again and type check them

type ErrorResourceCreatingInProgress struct {
	Message string
}

func (e *ErrorResourceCreatingInProgress) Error() string {
	return e.Message
}

type ErrorResourceDeletingInProgress struct {
	Message string
}

func (e *ErrorResourceDeletingInProgress) Error() string {
	return e.Message
}
