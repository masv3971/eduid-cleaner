package worker_manager

import "testing"

func TestEmptyChannels(t *testing.T) {
	tts := []struct {
		name string
		have *workers
	}{
		{
			name: "OK",
			have: &workers{
				skv:   nil,
				ladok: nil,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {

		})
	}

}
