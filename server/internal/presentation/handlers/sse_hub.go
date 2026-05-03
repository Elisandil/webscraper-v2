package handlers

import "sync"

// SSEHub manages per-user Server-Sent Events connections.
type SSEHub struct {
	mu      sync.RWMutex
	clients map[int64][]chan struct{}
}

func NewSSEHub() *SSEHub {
	return &SSEHub{clients: make(map[int64][]chan struct{})}
}

func (h *SSEHub) subscribe(userID int64) (chan struct{}, func()) {
	ch := make(chan struct{}, 1)
	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], ch)
	h.mu.Unlock()

	return ch, func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		chs := h.clients[userID]
		for i, c := range chs {
			if c == ch {
				h.clients[userID] = append(chs[:i], chs[i+1:]...)
				break
			}
		}
		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
		}
	}
}

// Notify signals all active SSE connections for the given user.
func (h *SSEHub) Notify(userID int64) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, ch := range h.clients[userID] {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
