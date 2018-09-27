package observer

// Stream represents the list of values a property is updated to.  For every
// property update, that value is appended to the list in the order they
// happen. The value is discarded once you advance the stream.  Please note
// that Stream is not goroutine safe: you cannot use the same stream on
// multiple goroutines concurrently. If you want to use multiple streams for
// the same property, either use Property.Observe (goroutine-safe) or use
// Stream.Clone (before passing it to another goroutine).
type Stream interface {
	// Wait 等待 Stream 更新，
	// 会在 Stream 没有发生更新时，发生阻塞
	Wait()

	// Value 可以获取 Stream 当前的值
	Value() interface{}

	// Next = Wait + Value
	// 注意： Next 无法获取到 Stream 生成时的第一个值
	Next() interface{}

	// Clone creates a new independent stream from this one but sharing the same
	// Property. Updates to the property will be reflected in both streams but
	// they may have different values depending on when they advance the stream
	// with Next.
	Clone() Stream
}

type stream struct {
	state *state
}

func (s *stream) Value() interface{} {
	return s.state.value
}

func (s *stream) Wait() {
	<-s.state.done
	s.state = s.state.next
}

func (s *stream) Next() interface{} {
	s.Wait()
	return s.Value()
}

func (s *stream) Clone() Stream {
	return &stream{state: s.state}
}
