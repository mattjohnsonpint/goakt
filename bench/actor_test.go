/*
 * MIT License
 *
 * Copyright (c) 2022-2025  Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package bench

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	actors "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/bench/benchpb"
	"github.com/tochemey/goakt/v3/internal/pause"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/passivation"
)

const receivingTimeout = 100 * time.Millisecond

func BenchmarkActor(b *testing.B) {
	b.Run("Tell(api:default mailbox)", func(b *testing.B) {
		ctx := context.TODO()

		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// define the benchmark actor
		actor := &Actor{}

		// create the actor ref
		pid, _ := actorSystem.Spawn(ctx, "test", actor)

		// wait for actors to start properly
		pause.For(1 * time.Second)

		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := actors.Tell(ctx, pid, new(benchpb.BenchTell)); err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})

		b.StopTimer()

		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")

		_ = pid.Shutdown(ctx)
		_ = actorSystem.Stop(ctx)
	})
	b.Run("Tell(actor-2-actor)", func(b *testing.B) {
		ctx := context.TODO()

		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// create the actors
		sender, _ := actorSystem.Spawn(ctx, "sender", new(Actor))
		receiver, _ := actorSystem.Spawn(ctx, "receiver", new(Actor))

		// wait for actors to start properly
		pause.For(1 * time.Second)
		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				err := sender.Tell(ctx, receiver, new(benchpb.BenchTell))
				if err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})
		b.StopTimer()
		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")
		_ = actorSystem.Stop(ctx)
	})
	b.Run("Tell(bounded mailbox)", func(b *testing.B) {
		ctx := context.TODO()

		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// create the actors
		sender, _ := actorSystem.Spawn(ctx, "sender", new(Actor), actors.WithMailbox(actors.NewBoundedMailbox(b.N)))
		receiver, _ := actorSystem.Spawn(ctx, "receiver", new(Actor), actors.WithMailbox(actors.NewBoundedMailbox(b.N)))

		// wait for actors to start properly
		pause.For(1 * time.Second)
		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				err := sender.Tell(ctx, receiver, new(benchpb.BenchTell))
				if err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})
		b.StopTimer()
		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")
		_ = actorSystem.Stop(ctx)
	})
	b.Run("Tell(priority mailbox)", func(b *testing.B) {
		ctx := context.TODO()

		priorityFunc := func(msg1, msg2 proto.Message) bool {
			p1 := msg1.(*benchpb.BenchPriorityMailbox)
			p2 := msg2.(*benchpb.BenchPriorityMailbox)
			return p1.Priority > p2.Priority
		}

		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// create the actors
		sender, _ := actorSystem.Spawn(ctx, "sender", new(Actor), actors.WithMailbox(actors.NewUnboundedPriorityMailBox(priorityFunc)))
		receiver, _ := actorSystem.Spawn(ctx, "receiver", new(Actor), actors.WithMailbox(actors.NewUnboundedPriorityMailBox(priorityFunc)))

		// wait for actors to start properly
		pause.For(1 * time.Second)
		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		var priorityCounter atomic.Int64
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				priorityCounter.Add(1)
				priorityMessage := &benchpb.BenchPriorityMailbox{
					Priority: priorityCounter.Load(),
				}
				err := sender.Tell(ctx, receiver, priorityMessage)
				if err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})
		b.StopTimer()
		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")
		_ = actorSystem.Stop(ctx)
	})
	b.Run("SendAsync(default mailbox)", func(b *testing.B) {
		ctx := context.TODO()

		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// create the actors
		sender, _ := actorSystem.Spawn(ctx, "sender", new(Actor))
		receiver, _ := actorSystem.Spawn(ctx, "receiver", new(Actor))

		// wait for actors to start properly
		pause.For(1 * time.Second)
		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := sender.SendAsync(ctx, receiver.Name(), new(benchpb.BenchTell)); err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})
		b.StopTimer()

		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")

		_ = actorSystem.Stop(ctx)
	})
	b.Run("Ask(api:default mailbox)", func(b *testing.B) {
		ctx := context.TODO()
		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// define the benchmark actor
		actor := &Actor{}

		// create the actor ref
		pid, _ := actorSystem.Spawn(ctx, "test", actor,
			actors.WithPassivationStrategy(
				passivation.NewTimeBasedStrategy(5*time.Second),
			))

		// wait for actors to start properly
		pause.For(1 * time.Second)
		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if _, err := actors.Ask(ctx, pid, new(benchpb.BenchRequest), receivingTimeout); err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})
		b.StopTimer()

		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")

		_ = pid.Shutdown(ctx)
		_ = actorSystem.Stop(ctx)
	})
	b.Run("Ask(actor-2-actor)", func(b *testing.B) {
		ctx := context.TODO()

		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// create the actors
		sender, _ := actorSystem.Spawn(ctx, "sender", new(Actor))
		receiver, _ := actorSystem.Spawn(ctx, "receiver", new(Actor))

		// wait for actors to start properly
		pause.For(1 * time.Second)

		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if _, err := sender.Ask(ctx, receiver, new(benchpb.BenchRequest), receivingTimeout); err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})
		b.StopTimer()

		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")

		_ = actorSystem.Stop(ctx)
	})
	b.Run("Ask(bounded mailbox)", func(b *testing.B) {
		ctx := context.TODO()

		// create the actor system
		actorSystem, _ := actors.NewActorSystem("bench",
			actors.WithLogger(log.DiscardLogger),
			actors.WithActorInitMaxRetries(1))

		// start the actor system
		_ = actorSystem.Start(ctx)

		// wait for system to start properly
		pause.For(1 * time.Second)

		// create the actors
		sender, _ := actorSystem.Spawn(ctx, "sender", new(Actor), actors.WithMailbox(actors.NewBoundedMailbox(b.N)))
		receiver, _ := actorSystem.Spawn(ctx, "receiver", new(Actor), actors.WithMailbox(actors.NewBoundedMailbox(b.N)))

		// wait for actors to start properly
		pause.For(1 * time.Second)

		var counter int64
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if _, err := sender.Ask(ctx, receiver, new(benchpb.BenchRequest), receivingTimeout); err != nil {
					b.Fatal(err)
				}
				atomic.AddInt64(&counter, 1)
			}
		})
		b.StopTimer()

		messagesPerSec := float64(atomic.LoadInt64(&counter)) / b.Elapsed().Seconds()
		b.ReportMetric(messagesPerSec, "messages/sec")

		_ = actorSystem.Stop(ctx)
	})
}
