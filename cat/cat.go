package cat

import (
	"os"
)

func Init(domain string) {
	config.Init(domain)

	go router.Background()
	go monitor.Background()
	go aggregator.Background()
	go sender.Background()
}

func Shutdown() {
	scheduler.shutdown()
}

func DebugOn() {
	logger.logger.SetOutput(os.Stdout)
}
