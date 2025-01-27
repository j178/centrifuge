package natsbroker

import (
	"strconv"
	"testing"

	"github.com/centrifugal/centrifuge"
)

func newTestNatsBroker() *NatsBroker {
	return NewTestNatsBrokerWithPrefix("centrifuge-test")
}

func NewTestNatsBrokerWithPrefix(prefix string) *NatsBroker {
	n, _ := centrifuge.New(centrifuge.Config{})
	b, _ := New(n, Config{Prefix: prefix})
	n.SetBroker(b)
	err := n.Run()
	if err != nil {
		panic(err)
	}
	return b
}

func BenchmarkNatsBrokerPublish(b *testing.B) {
	broker := newTestNatsBroker()
	rawData := []byte(`{"bench": true}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := broker.Publish("channel", rawData, centrifuge.PublishOptions{})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkNatsBrokerPublishParallel(b *testing.B) {
	broker := newTestNatsBroker()
	rawData := []byte(`{"bench": true}`)
	b.SetParallelism(128)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := broker.Publish("channel", rawData, centrifuge.PublishOptions{})
			if err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkNatsBrokerSubscribe(b *testing.B) {
	broker := newTestNatsBroker()
	j := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j++
		err := broker.Subscribe("subscribe" + strconv.Itoa(j))
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkNatsBrokerSubscribeParallel(b *testing.B) {
	broker := newTestNatsBroker()
	i := 0
	b.SetParallelism(128)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			err := broker.Subscribe("subscribe" + strconv.Itoa(i))
			if err != nil {
				panic(err)
			}
		}
	})
}
