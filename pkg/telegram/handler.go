package telegram

type HandlerFunc func(c *Context) bool

type UpdateHandler struct {
	command  string
	handlers []HandlerFunc
}

type HandlerError struct {
	Message       string
	MessageForLog string
}

func (e *HandlerError) Error() string {
	return e.Message
}
