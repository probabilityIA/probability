package mocks

import "context"

// AIPauseCheckerMock es un mock de ports.IAIPauseChecker
type AIPauseCheckerMock struct {
	Paused bool
}

func (m *AIPauseCheckerMock) IsAIPaused(_ context.Context, _ string) bool {
	return m.Paused
}
