// Copyright 2023 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt

package streamingest

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/ccl/streamingccl/streamclient"
	"github.com/cockroachdb/cockroach/pkg/util/ctxgroup"
)

// MergedSubscription combines multiple subscriptions into a single
// merged stream of events.
type MergedSubscription struct {
	cg       ctxgroup.Group
	cgCancel context.CancelFunc
	eventCh  chan PartitionEvent
}

func MergeSubscriptions(
	ctx context.Context, subscriptions map[string]streamclient.Subscription,
) *MergedSubscription {
	ctx, cancel := context.WithCancel(ctx)
	m := &MergedSubscription{
		cg:       ctxgroup.WithContext(ctx),
		cgCancel: cancel,
		eventCh:  make(chan PartitionEvent),
	}
	for partition, sub := range subscriptions {
		partition := partition
		sub := sub
		m.cg.GoCtx(func(ctx context.Context) error {
			ctxDone := ctx.Done()
			for {
				select {
				case event, ok := <-sub.Events():
					if !ok {
						return sub.Err()
					}

					pe := PartitionEvent{
						Event:     event,
						partition: partition,
					}

					select {
					case m.eventCh <- pe:
					case <-ctxDone:
						return ctx.Err()
					}
				case <-ctxDone:
					return ctx.Err()
				}
			}
		})
	}
	return m
}

// Run blocks until the merged stream is closed.
func (m *MergedSubscription) Run() error {
	err := m.cg.Wait()
	close(m.eventCh)
	return err
}

// Close stops the merged stream. Note that the underlying
// subscriptions are not closed.
func (m *MergedSubscription) Close() {
	m.cgCancel()
}

// Events returns the merged event channel.
func (m *MergedSubscription) Events() chan PartitionEvent {
	return m.eventCh
}
