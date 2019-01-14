package cat

func Init(domain string) {
	config.Init(domain)

	go router.Background()
	go monitor.Background()
	go aggregator.Background()
	go sender.Background()
}

func Shutdown() {
	monitor.Shutdown()
}

func Wait() {
}
