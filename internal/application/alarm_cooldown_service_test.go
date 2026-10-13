package application

import (
	"sync"
	"testing"
	"time"
)

func TestAlarmCooldown_ImmediateTrigger(t *testing.T) {
	cd := NewAlarmCooldownService()
	if !cd.ShouldTrigger("cam_01", "motion", 50*time.Millisecond) {
		t.Fatalf("First event should always trigger")
	}
}

func TestAlarmCooldown_SuppressionWithinWindow(t *testing.T) {
	cd := NewAlarmCooldownService()
	cam := "cam_entrance"
	evt := "intrusion"
	window := 100 * time.Millisecond

	// 1st: true
	if !cd.ShouldTrigger(cam, evt, window) {
		t.Fatalf("Expected 1st trigger to be true")
	}

	// 2nd immediately: false
	if cd.ShouldTrigger(cam, evt, window) {
		t.Fatalf("Expected immediate duplicate trigger to be suppressed")
	}

	// Wait for window to expire
	time.Sleep(120 * time.Millisecond)

	// 3rd after window: true
	if !cd.ShouldTrigger(cam, evt, window) {
		t.Fatalf("Expected trigger after cooldown window to be true")
	}
}

func TestAlarmCooldown_DifferentCamerasIndependent(t *testing.T) {
	cd := NewAlarmCooldownService()
	window := 100 * time.Millisecond

	if !cd.ShouldTrigger("cam_01", "vmd", window) {
		t.Fatalf("cam_01 should trigger")
	}

	// cam_02 should trigger even though cam_01 just triggered
	if !cd.ShouldTrigger("cam_02", "vmd", window) {
		t.Fatalf("cam_02 should trigger independently")
	}
}

func TestAlarmCooldown_Concurrency(t *testing.T) {
	cd := NewAlarmCooldownService()
	var wg sync.WaitGroup
	allowedCount := 0
	var mu sync.Mutex

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if cd.ShouldTrigger("cam_shared", "tamper", 500*time.Millisecond) {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if allowedCount != 1 {
		t.Fatalf("Expected exactly 1 trigger across concurrent calls, got %d", allowedCount)
	}
}
