package state

// Read safely reads state
func (s *State) Read(fn func(*State)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s)
}

// Write safley writes to state
func (s *State) Write(fn func(*State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fn(s)
}
