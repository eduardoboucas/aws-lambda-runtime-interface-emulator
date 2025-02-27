// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package metering

import (
	_ "runtime" //for nanotime() and walltime()
	"time"
	_ "unsafe" //for go:linkname
)

//go:linkname Monotime runtime.nanotime
func Monotime() int64

//go:linkname walltime runtime.walltime
//func walltime() (sec int64, nsec int32)

// MonoToEpoch converts monotonic time nanos to epoch time nanos.
func MonoToEpoch(t int64) int64 {
	// `runtime.walltime` is not available on Windows. We need an alternative way
	// of getting the wall clock, but for now we can live with less accuracy.
	monoNsec := Monotime()

	wallTime := time.Now()
	wallNsec := wallTime.UnixNano()

	clockOffset := wallNsec - monoNsec
	return t + clockOffset
}

type ExtensionsResetDurationProfiler struct {
	NumAgentsRegisteredForShutdown int
	AvailableNs                    int64
	extensionsResetStartTimeNs     int64
	extensionsResetEndTimeNs       int64
}

func (p *ExtensionsResetDurationProfiler) Start() {
	p.extensionsResetStartTimeNs = Monotime()
}

func (p *ExtensionsResetDurationProfiler) Stop() {
	p.extensionsResetEndTimeNs = Monotime()
}

func (p *ExtensionsResetDurationProfiler) CalculateExtensionsResetMs() (int64, bool) {
	var extensionsResetDurationNs = p.extensionsResetEndTimeNs - p.extensionsResetStartTimeNs
	var extensionsResetMs int64
	timedOut := false

	if p.NumAgentsRegisteredForShutdown == 0 || p.AvailableNs < 0 || extensionsResetDurationNs < 0 {
		extensionsResetMs = 0
	} else if extensionsResetDurationNs > p.AvailableNs {
		extensionsResetMs = p.AvailableNs / time.Millisecond.Nanoseconds()
		timedOut = true
	} else {
		extensionsResetMs = extensionsResetDurationNs / time.Millisecond.Nanoseconds()
	}

	return extensionsResetMs, timedOut
}
