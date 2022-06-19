package events

type TestsEventBus struct {
	wantErr error
	called  bool
}

func NewTestsEventBus(wantErr error) *TestsEventBus {
	return &TestsEventBus{wantErr: wantErr}
}

func (t *TestsEventBus) Enqueue(_ Event) error {
	t.called = true

	if t.wantErr != nil {
		return t.wantErr
	}

	return nil
}

func (t *TestsEventBus) Called() bool {
	return t.called
}
