package lib

// so that I can return these errors when reconcile is running again and type check them

// rds resourcees creating in progress
type ErrorResourceCreatingInProgress struct {
	Message string
}

// rds resourcees creating in progress
func (e *ErrorResourceCreatingInProgress) Error() string {
	return e.Message
}

// rds resourcees deleting in progress
type ErrorResourceDeletingInProgress struct {
	Message string
}

// rds resourcees deleting in progress
func (e *ErrorResourceDeletingInProgress) Error() string {
	return e.Message
}

type ErrorKubernetesSecretDoesNotHaveKeyError struct {
	Message string
}

func (e *ErrorKubernetesSecretDoesNotHaveKeyError) Error() string {
	return e.Message
}
