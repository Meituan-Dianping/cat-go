package cat

import (
	"os"
	"sync/atomic"
)

var enabled uint32 = 0

var started uint32 = 0

func Init(domain string) {
	if err := config.Init(domain); err != nil {
		logger.Warning("Cat initialize failed.")
		return
	}

	if !atomic.CompareAndSwapUint32(&started, 0, 1) {
		// Cat goroutines has already been started.
		return
	}
	enable()

	go background(&router)
	go background(&monitor)
	go background(&sender)
	aggregator.Background()
}

func enable() {
	if atomic.SwapUint32(&enabled, 1) == 0 {
		logger.Info("Cat has been enabled.")
	}
}

func disable() {
	if atomic.SwapUint32(&enabled, 0) == 1 {
		logger.Info("Cat has been disabled.")
	}
}

func IsEnabled() bool {
	return atomic.LoadUint32(&enabled) > 0
}

func Shutdown() {
	scheduler.shutdown()
}

func Wait() {
	// outdated
}

func DebugOn() {
	logger.logger.SetOutput(os.Stdout)
}
