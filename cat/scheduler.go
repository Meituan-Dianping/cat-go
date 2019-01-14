package cat

type nameGetter interface {
	GetName() string
}

type signalsMixer interface {
	nameGetter
	shutdown()
	getExitSignal() int
}

type signalsMixin struct {
	isAlive    bool
	signals    chan int
	exitSignal int
}

func makeSignalsMixedIn(exitSignal int) signalsMixin {
	return signalsMixin{
		isAlive:    true,
		signals:    make(chan int),
		exitSignal: exitSignal,
	}
}

func (p *signalsMixin) stop() {
	p.isAlive = false
}

func (p *signalsMixin) exit() {
	close(p.signals)
	scheduler.signals <- p.exitSignal
}

func (p *signalsMixin) getExitSignal() int {
	return p.exitSignal
}

func (p *signalsMixin) shutdown() {
	p.signals <- signalShutdown
}

type catScheduler struct {
	signalsMixin
}

var scheduler = catScheduler{
	makeSignalsMixedIn(signal0),
}

func (p *catScheduler) shutdownAndWaitGroup(items []signalsMixer) {
	var expectedSignals = make(map[int]string)
	var count = 0

	for _, v := range items {
		v.shutdown()
		expectedSignals[v.getExitSignal()] = v.GetName()
		count++
	}

	for signal := range p.signals {
		if name, ok := expectedSignals[signal]; ok {
			count--
			logger.Info("%s exited.", name)
		} else {
			logger.Warning("Unpredicted signal received: %d", signal)
		}
		if count == 0 {
			break
		}
	}
}

func (p *catScheduler) shutdown() {
	group1 := []signalsMixer{&router, &monitor}
	group2 := []signalsMixer{aggregator.transaction, aggregator.event, aggregator.metric}
	group3 := []signalsMixer{&sender}

	logger.Info("Received shutdown request, scheduling...")
	// TODO disable cat api.

	p.shutdownAndWaitGroup(group1)
	p.shutdownAndWaitGroup(group2)
	p.shutdownAndWaitGroup(group3)

	logger.Info("All systems down.")
}
