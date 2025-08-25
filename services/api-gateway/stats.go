package main

import (
	"sync"
	"time"
)

// StatsCollector collects and manages gateway statistics
type StatsCollector struct {
	mu              sync.RWMutex
	totalRequests   int64
	successRequests int64
	errorRequests   int64
	totalLatency    time.Duration
	serviceStats    map[string]*ServiceStats
}

// NewStatsCollector creates a new stats collector
func NewStatsCollector() *StatsCollector {
	return &StatsCollector{
		serviceStats: make(map[string]*ServiceStats),
	}
}

// RecordRequest records a request with its outcome and latency
func (s *StatsCollector) RecordRequest(service string, success bool, latency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update total stats
	s.totalRequests++
	s.totalLatency += latency

	if success {
		s.successRequests++
	} else {
		s.errorRequests++
	}

	// Update service-specific stats
	if s.serviceStats[service] == nil {
		s.serviceStats[service] = &ServiceStats{}
	}

	serviceStats := s.serviceStats[service]
	serviceStats.Requests++
	serviceStats.LastRequest = time.Now()

	if success {
		serviceStats.Successes++
	} else {
		serviceStats.Errors++
	}

	// Calculate running average for latency
	if serviceStats.Requests == 1 {
		serviceStats.AvgLatency = float64(latency.Nanoseconds()) / 1e6 // Convert to milliseconds
	} else {
		// Use exponential moving average for better performance
		alpha := 0.1 // Smoothing factor
		newLatencyMs := float64(latency.Nanoseconds()) / 1e6
		serviceStats.AvgLatency = alpha*newLatencyMs + (1-alpha)*serviceStats.AvgLatency
	}
}

// GetStats returns current gateway statistics
func (s *StatsCollector) GetStats() GatewayStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Calculate average response time
	var avgResponseTime float64
	if s.totalRequests > 0 {
		avgResponseTime = float64(s.totalLatency.Nanoseconds()) / float64(s.totalRequests) / 1e6 // Convert to milliseconds
	}

	// Copy service stats
	serviceStatsCopy := make(map[string]ServiceStats)
	for service, stats := range s.serviceStats {
		serviceStatsCopy[service] = ServiceStats{
			Requests:    stats.Requests,
			Successes:   stats.Successes,
			Errors:      stats.Errors,
			AvgLatency:  stats.AvgLatency,
			LastRequest: stats.LastRequest,
		}
	}

	return GatewayStats{
		TotalRequests:   s.totalRequests,
		SuccessRequests: s.successRequests,
		ErrorRequests:   s.errorRequests,
		AvgResponseTime: avgResponseTime,
		ServiceStats:    serviceStatsCopy,
	}
}

// Reset resets all statistics
func (s *StatsCollector) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.totalRequests = 0
	s.successRequests = 0
	s.errorRequests = 0
	s.totalLatency = 0
	s.serviceStats = make(map[string]*ServiceStats)
}

// GetServiceStats returns statistics for a specific service
func (s *StatsCollector) GetServiceStats(service string) (ServiceStats, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats, exists := s.serviceStats[service]
	if !exists {
		return ServiceStats{}, false
	}

	return ServiceStats{
		Requests:    stats.Requests,
		Successes:   stats.Successes,
		Errors:      stats.Errors,
		AvgLatency:  stats.AvgLatency,
		LastRequest: stats.LastRequest,
	}, true
}