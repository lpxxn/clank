package internal

type httpServer struct {
}
type Handler struct {
	Method   string
	FullPath string
	respBody map[string]struct{}
}
