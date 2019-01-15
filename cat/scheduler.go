package cat

type scheduleMixin struct {
	isAlive    bool
	signals    chan int
	exitSignal int
}

type scheduleMixer interface {
	GetName() string

	handle(signal int)

	process()
	afterStart()
	beforeStop()

	getScheduleMixin() *scheduleMixin
}

func (p *scheduleMixin) handle(signal int) {
	switch signal {
	case signalShutdown:
		p.isAlive = false
	}
}

func (p *scheduleMixin) process() {
	return
}

func (p *scheduleMixin) afterStart() {
	return
}

func (p *scheduleMixin) beforeStop() {
	return
}

func (p *scheduleMixin) getScheduleMixin() *scheduleMixin {
	return p
}

func background(p scheduleMixer) {
	mixin := p.getScheduleMixin()

	mixin.isAlive = true
	p.afterStart()

	for mixin.isAlive {
		p.process()
	}

	p.beforeStop()

	close(mixin.signals)
	scheduler.signals <- mixin.exitSignal
}

func makeScheduleMixedIn(exitSignal int) scheduleMixin {
	return scheduleMixin{
		isAlive:    false,
		signals:    make(chan int),
		exitSignal: exitSignal,
	}
}

type catScheduler struct {
	signals chan int
}

var scheduler = catScheduler{
	signals: make(chan int),
}

func (p *catScheduler) shutdownAndWaitGroup(items []scheduleMixer) {
	var expectedSignals = make(map[int]string)
	var count = 0

	for _, v := range items {
		mixin := v.getScheduleMixin()
		if mixin.isAlive {
			mixin.signals <- signalShutdown
			expectedSignals[mixin.exitSignal] = v.GetName()
			count++
		}
	}

	if count == 0 {
		return
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
	group1 := []scheduleMixer{&router, &monitor}
	group2 := []scheduleMixer{aggregator.transaction, aggregator.event, aggregator.metric}
	group3 := []scheduleMixer{&sender}

	disable()

	logger.Info("Received shutdown request, scheduling...")

	p.shutdownAndWaitGroup(group1)
	p.shutdownAndWaitGroup(group2)
	p.shutdownAndWaitGroup(group3)

	logger.Info("All systems down.")
}
