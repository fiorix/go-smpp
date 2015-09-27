// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"sync"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

// deliveryPool let peers register themselves to receive broadcast
// notifications with delivery receipts.
type deliveryPool struct {
	mu sync.Mutex
	m  map[string]chan *DeliveryReceipt
}

// Handler handles DeliverSM coming from a Transceiver SMPP connection.
// It broadcasts received delivery receipt to all registered peers.
func (pool *deliveryPool) Handler(p pdu.Body) {
	switch p.Header().ID {
	case pdu.DeliverSMID:
		f := p.Fields()
		dr := &DeliveryReceipt{
			Src:  f[pdufield.SourceAddr].String(),
			Dst:  f[pdufield.DestinationAddr].String(),
			Text: f[pdufield.ShortMessage].String(),
		}
		pool.Broadcast(dr)
	}
}

// Register adds peer k to receiving delivery receipts over the
// returned channel.
func (pool *deliveryPool) Register(k string) chan *DeliveryReceipt {
	c := make(chan *DeliveryReceipt, 10)
	pool.mu.Lock()
	pool.m[k] = c
	pool.mu.Unlock()
	return c
}

// Unregister removes peer k from the delivery receipt broadcast,
// and closes the channel previously returned by Register.
func (pool *deliveryPool) Unregister(k string) {
	pool.mu.Lock()
	c := pool.m[k]
	if c != nil {
		delete(pool.m, k)
		close(c)
	}
	pool.mu.Unlock()
}

// Broadcast broadcasts the given delivery receipt to all registered peers.
func (pool *deliveryPool) Broadcast(r *DeliveryReceipt) {
	pool.mu.Lock()
	for _, c := range pool.m {
		select {
		case c <- r:
		default:
		}
	}
	pool.mu.Unlock()
}
