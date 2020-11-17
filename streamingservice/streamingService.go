package streamingservice

// StreamingService struct
// - file → absolute path of the streamed music
type StreamingService struct {
	file string
}

// New return a pointer to a new streamingservice
func New() *StreamingService {
	return new(StreamingService)
}

// Start a streaming given the absolute path of the music
func (ss *StreamingService) Start(song string) {
	ss.file = song
}

// Stop a streaming
func (ss *StreamingService) Stop() {
	ss.file = ""
}

// Pause a streaming
func (ss *StreamingService) Pause() {
}

// Continue a streaminf after a pause
func (ss *StreamingService) Continue() {
}
