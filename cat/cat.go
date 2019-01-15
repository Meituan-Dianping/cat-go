package cat

import (
	"os"
)

func Init(domain string) {
	if err := config.Init(domain); err != nil {
		// TODO disable cat.
		return
	}

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
