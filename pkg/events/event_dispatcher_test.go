package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

type TestEventHandler struct {
	ID string
}

func (h *TestEventHandler) Handle(event EventInterface) {
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event1          TestEvent
	event2          TestEvent
	handler1        TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (s *EventDispatcherTestSuite) SetupTest() {
	s.eventDispatcher = NewEventDispatcher()
	s.handler1 = TestEventHandler{ID: "1"}
	s.handler2 = TestEventHandler{ID: "2"}
	s.handler3 = TestEventHandler{ID: "3"}
	s.event1 = TestEvent{
		Name:    "event1",
		Payload: "payload1",
	}
	s.event2 = TestEvent{
		Name:    "event2",
		Payload: "payload2",
	}
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := s.eventDispatcher.Register(s.event1.GetName(), &s.handler1)

	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	err = s.eventDispatcher.Register(s.event1.GetName(), &s.handler2)

	s.Nil(err)
	s.Equal(2, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	assert.Equal(s.T(), &s.handler1, s.eventDispatcher.handlers[s.event1.GetName()][0])
	assert.Equal(s.T(), &s.handler2, s.eventDispatcher.handlers[s.event1.GetName()][1])
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := s.eventDispatcher.Register(s.event1.GetName(), &s.handler1)

	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	err = s.eventDispatcher.Register(s.event1.GetName(), &s.handler1)

	s.NotNil(err)
	s.Equal(ErrHandlerAlreadyRegistered, err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event1.GetName()]))
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := s.eventDispatcher.Register(s.event1.GetName(), &s.handler1)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	err = s.eventDispatcher.Register(s.event1.GetName(), &s.handler2)
	s.Nil(err)
	s.Equal(2, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	err = s.eventDispatcher.Register(s.event2.GetName(), &s.handler3)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event2.GetName()]))

	err = s.eventDispatcher.Clear()

	s.Nil(err)
	s.Equal(0, len(s.eventDispatcher.handlers[s.event1.GetName()]))
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	err := s.eventDispatcher.Register(s.event1.GetName(), &s.handler1)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	err = s.eventDispatcher.Register(s.event1.GetName(), &s.handler2)
	s.Nil(err)
	s.Equal(2, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	assert.True(s.T(), s.eventDispatcher.Has(s.event1.GetName(), &s.handler1))
	assert.True(s.T(), s.eventDispatcher.Has(s.event1.GetName(), &s.handler2))
	assert.False(s.T(), s.eventDispatcher.Has(s.event1.GetName(), &s.handler3))
}

type mockEventHandler struct {
	mock.Mock
}

func (m *mockEventHandler) Handle(event EventInterface) {
	m.Called(event)
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eh := &mockEventHandler{}
	eh.On("Handle", &s.event1)
	s.eventDispatcher.Register(s.event1.GetName(), eh)

	s.eventDispatcher.Dispatch(&s.event1)

	eh.AssertExpectations(s.T())
	eh.AssertNumberOfCalls(s.T(), "Handle", 1)
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	err := s.eventDispatcher.Register(s.event1.GetName(), &s.handler1)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	err = s.eventDispatcher.Register(s.event1.GetName(), &s.handler2)
	s.Nil(err)
	s.Equal(2, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	err = s.eventDispatcher.Register(s.event2.GetName(), &s.handler3)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event2.GetName()]))

	s.eventDispatcher.Remove(s.event1.GetName(), &s.handler1)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event1.GetName()]))
	assert.Equal(s.T(), &s.handler2, s.eventDispatcher.handlers[s.event1.GetName()][0])

	s.eventDispatcher.Remove(s.event1.GetName(), &s.handler2)
	s.Equal(0, len(s.eventDispatcher.handlers[s.event1.GetName()]))

	s.eventDispatcher.Remove(s.event2.GetName(), &s.handler3)
	s.Equal(0, len(s.eventDispatcher.handlers[s.event2.GetName()]))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
