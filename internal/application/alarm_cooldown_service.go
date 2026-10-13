package application

import (
	"fmt"
	"sync"
	"time"
)

// AlarmCooldownService prevents alert storms by enforcing sliding cooldown windows per camera and event type.
type AlarmCooldownService struct {
	mu            sync.RWMutex
	lastTriggered map[string]time.Time
}

// NewAlarmCooldownService instantiates a thread-safe alarm cooldown coordinator.
func NewAlarmCooldownService() *AlarmCooldownService {
	return &AlarmCooldownService{
		lastTriggered: make(map[string]time.Time),
	}
}

// ShouldTrigger checks if an alarm event should be dispatched or suppressed.
// If the time since the last event of this type on this camera exceeds the cooldown window,
// it updates the timestamp and returns true. Otherwise, it returns false (suppressed).
func (s *AlarmCooldownService) ShouldTrigger(cameraID, eventType string, cooldownWindow time.Duration) bool {
	if cooldownWindow <= 0 {
		cooldownWindow = 30 * time.Second
	}

	key := fmt.Sprintf("%s:%s", cameraID, eventType)
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	lastTime, exists := s.lastTriggered[key]
	if !exists || now.Sub(lastTime) >= cooldownWindow {
		s.lastTriggered[key] = now
		return true
	}

	return false
}

// Reset clears the cooldown record for a given camera and event type.
func (s *AlarmCooldownService) Reset(cameraID, eventType string) {
	key := fmt.Sprintf("%s:%s", cameraID, eventType)
	s.mu.Lock()
	delete(s.lastTriggered, key)
	s.mu.Unlock()
}

// PurgeExpired removes entries that have been idle longer than maxIdleDuration to prevent memory growth.
func (s *AlarmCooldownService) PurgeExpired(maxIdleDuration time.Duration) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, lastTime := range s.lastTriggered {
		if now.Sub(lastTime) > maxIdleDuration {
			delete(s.lastTriggered, k)
		}
	}
}
