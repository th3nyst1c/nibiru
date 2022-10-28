package events

import (
	"github.com/NibiruChain/nibiru/feeder/types"
	"github.com/rs/zerolog"
)

// wsI exists for testing purposes.
type wsI interface {
	message() <-chan []byte
	close()
}

type Stream struct {
	stop, done chan struct{}

	ws wsI

	votingPeriod chan types.VotingPeriod
	params       chan types.Params
}

func (s *Stream) loop() {
	defer close(s.done)
	defer s.ws.close()

	for {
		select {
		case <-s.stop:
			return
		case msg := <-s.ws.message():
			panic(msg) // todo
		}
	}
}

func NewStream(tendermintRPC string, log zerolog.Logger) *Stream {
	stream := &Stream{
		stop:         make(chan struct{}),
		done:         make(chan struct{}),
		ws:           dial(tendermintRPC, nil, log.With().Str("component", "events.Stream.websocket").Logger()),
		votingPeriod: make(chan types.VotingPeriod),
		params:       make(chan types.Params, 1), // for initial params
	}

	go stream.loop()
	return stream
}
