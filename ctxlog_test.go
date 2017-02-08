package ctxlog

import (
	"testing"
	"context"
	"github.com/apex/log"
	"github.com/stretchr/testify/suite"
)

type MockHandler struct {
	LatestEntry *log.Entry
}

type TestSuite struct {
	suite.Suite
	ctx context.Context
	handler *MockHandler
	logger log.Interface
}

func (m *MockHandler) HandleLog(entry *log.Entry) error {
	m.LatestEntry = entry
	return nil
}

func (s *TestSuite) SetupTest() {
	s.ctx = context.Background()
	s.handler = &MockHandler{}
	s.logger = &log.Logger{
		Handler: s.handler,
		Level: log.InfoLevel,
	}
}

func TestCtxLogSuite(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) TestGetLoggerEmptyContext() {
	logger := GetContextualLogger(s.ctx, s.logger)
	logger.Info("test")
	s.Len(s.handler.LatestEntry.Fields, 0)
}

func (s *TestSuite) TestReturnedLoggerHasFields() {
	_, logger := GetUpdatedLoggingContext(s.ctx, s.logger, log.Fields{
		"f1": "v1",
		"f2": "v2",
	})
	logger.Info("test")
	s.Equal(log.Fields{"f1": "v1", "f2": "v2"}, s.handler.LatestEntry.Fields)
}

func (s *TestSuite) TestContextCarriesFields() {
	ctx, _ := GetUpdatedLoggingContext(s.ctx, s.logger, log.Fields{
		"f1": "v1",
		"f2": "v2",
	})

	logger := GetContextualLogger(ctx, s.logger)

	logger.Info("test")
	s.Equal(log.Fields{"f1": "v1", "f2": "v2"}, s.handler.LatestEntry.Fields)
}

func (s *TestSuite) TestMergingFields() {
	ctx, _ := GetUpdatedLoggingContext(s.ctx, s.logger, log.Fields{
		"f1": "v1",
		"f2": "v2",
	})

	ctx2, _ := GetUpdatedLoggingContext(ctx, s.logger, log.Fields{
		"f2": "v2-new",
		"f3": "v3",
	})

	logger := GetContextualLogger(ctx2, s.logger)

	logger.Info("test")
	s.Equal(log.Fields{"f1": "v1", "f2": "v2-new", "f3": "v3"}, s.handler.LatestEntry.Fields)
}

func (s *TestSuite) TestEmptyUpdate() {
	ctx, _ := GetUpdatedLoggingContext(s.ctx, s.logger, log.Fields{})
	logger := GetContextualLogger(ctx, s.logger)

	logger.Info("test")
	s.Len(s.handler.LatestEntry.Fields, 0)
}

func (s *TestSuite) TestContextGarbageGet() {
	ctx := context.WithValue(s.ctx, logFieldsContextKey, "x")
	logger := GetContextualLogger(ctx, s.logger)
	logger.Info("test")
	s.Len(s.handler.LatestEntry.Fields, 0)
}

func (s *TestSuite) TestContextGarbageUpdate() {
	ctx := context.WithValue(s.ctx, logFieldsContextKey, "x")
	ctx2, _ := GetUpdatedLoggingContext(ctx, s.logger, log.Fields{"f1": "v1"})

	logger := GetContextualLogger(ctx2, s.logger)
	logger.Info("test")
	s.Equal(log.Fields{"f1": "v1"}, s.handler.LatestEntry.Fields)
}
