package Subscription

import (
	"context"
	"ozon/internal/transport/graph/model"
	"sync"
)

type Subscription struct {
	commentSubscriptions map[string][]chan *model.Comment
	lock                 sync.Mutex
}

func New() *Subscription {
	return &Subscription{
		commentSubscriptions: make(map[string][]chan *model.Comment),
		lock:                 sync.Mutex{},
	}
}

func (p *Subscription) Subscribe(ctx context.Context, postId string) chan *model.Comment {
	p.lock.Lock()
	defer p.lock.Unlock()

	ch := make(chan *model.Comment, 1)
	p.commentSubscriptions[postId] = append(p.commentSubscriptions[postId], ch)

	return ch
}

func (p *Subscription) Publish(ctx context.Context, comment *model.Comment) {
	go func() {
		p.lock.Lock()
		defer p.lock.Unlock()

		if subscribers, ok := p.commentSubscriptions[comment.PostID]; ok {
			for _, ch := range subscribers {
				ch <- comment
			}
		}
	}()
}

func (p *Subscription) Check(postId string) bool {
	_, exists := p.commentSubscriptions[postId]
	return exists
}

func (p *Subscription) Unsubscribe(ctx context.Context, postId string, ch chan *model.Comment) {
	p.lock.Lock()
	defer p.lock.Unlock()

	subs := p.commentSubscriptions[postId]
	newSubs := make([]chan *model.Comment, 0, len(subs))
	for _, sub := range subs {
		if sub != ch {
			newSubs = append(newSubs, sub)
		}
	}
	p.commentSubscriptions[postId] = newSubs

	close(ch)
}
