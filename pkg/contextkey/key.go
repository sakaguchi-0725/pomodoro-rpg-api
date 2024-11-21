package contextkey

type ContextKey string

const (
	Email     ContextKey = "email"
	UserID    ContextKey = "userID"
	RequestID ContextKey = "requestID"
)
