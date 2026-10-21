package main

import "testing"

func TestReturnTo(t *testing.T) {
	tests := []struct {
		name, provider, workflow, campaign, want string
		invalid                                  bool
	}{
		{"donor receipt", "google", "receipt", "spring 2026", "https://nonprofit.example/donor/receipts?campaign=spring+2026", false},
		{"volunteer reminder", "github", "reminder", "food-drive", "https://nonprofit.example/volunteer/reminders?campaign=food-drive", false},
		{"campaign reporting", "google", "report", "fall", "https://nonprofit.example/campaign/reports?campaign=fall", false},
		{"unknown workflow", "google", "admin", "fall", "", true},
		{"unknown provider", "other", "receipt", "fall", "", true},
		{"missing campaign", "google", "receipt", "", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := returnTo("https://nonprofit.example", loginRequest{tc.provider, tc.workflow, tc.campaign})
			if (err != nil) != tc.invalid || got != tc.want {
				t.Fatalf("returnTo() = %q, %v; want %q, invalid=%v", got, err, tc.want, tc.invalid)
			}
		})
	}
}
