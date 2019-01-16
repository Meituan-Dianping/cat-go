package cat

import (
	"github.com/Meituan-Dianping/cat-go/message"
	"sync"
)

type catMessageTree struct {
	message.Transaction

	mu sync.Mutex
	wg sync.WaitGroup

	messageId       string
	parentMessageId string
	rootMessageId   string
}

func NewMessageTree() *catMessageTree {

	tree := &catMessageTree{
		Transaction: *message.NewTransaction("System", "Tree", nil),

		mu: sync.Mutex{},
		wg: sync.WaitGroup{},

		messageId:       "",
		parentMessageId: "",
		rootMessageId:   "",
	}
	return tree
}

func (p *catMessageTree) Wait() {
	p.wg.Wait()
}

func (p *catMessageTree) hasProblem() bool {
	for _, m := range p.GetChildren() {
		if m.GetStatus() != SUCCESS {
			return true
		}
	}
	return false
}

func (p *catMessageTree) Complete() {
	if p.IsCompleted() {
		// do nothing.
	} else {
		p.Transaction.Complete()
		p.GetMessageId()
		manager.flush(p)
	}
}

func (p *catMessageTree) flush(message.Messager) {
	p.wg.Done()
}

func (p *catMessageTree) NewTransaction(mtype, name string) message.Transactor {
	if p.IsCompleted() {
		return NewTransaction(mtype, name)
	}
	p.wg.Add(1)

	trans := message.NewTransaction(mtype, name, p.flush)
	p.AddChild(trans)

	return trans
}

func (p *catMessageTree) GetMessageId() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.messageId == "" {
		p.messageId = manager.nextId()
	}
	return p.messageId
}

func (p *catMessageTree) SetParentMessageId(parentMessageId string) {
	p.parentMessageId = parentMessageId
}

func (p *catMessageTree) SetRootMessageId(rootMessageId string) {
	p.rootMessageId = rootMessageId
}
