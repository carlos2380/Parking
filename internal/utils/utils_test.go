package utils

import (
	"testing"
	"time"
)

func TestCentsIntToEurStr(t *testing.T) {

	tests := []struct {
		name   string
		cents  int
		expect string
	}{
		{"100 cents", 100, "1.00"},
		{"50 cents", 50, "0.50"},
		{"0 cents", 0, "0.00"},
		{"299 cents", 299, "2.99"},
		{"-150 cents", -150, "-1.50"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CentsIntToEurStr(test.cents)
			if test.expect != result {
				t.Errorf("CentsIntToEurStr(%d) = %s; want %s", test.cents, result, test.expect)
			}
		})
	}
}

func TestDateToStr(t *testing.T) {
	tests := []struct {
		name   string
		date   time.Time
		expect string
	}{
		{
			name:   "Date in 2024",
			date:   time.Date(2024, time.August, 15, 14, 45, 30, 0, time.UTC),
			expect: "15 Aug 24 14:45 UTC",
		},
		{
			name:   "Date in 2000",
			date:   time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			expect: "01 Jan 00 00:00 UTC",
		},
		{
			name:   "Date in 1999",
			date:   time.Date(1999, time.December, 31, 23, 59, 59, 0, time.UTC),
			expect: "31 Dec 99 23:59 UTC",
		},
	}
	for _, test := range tests {
		result := DateToStr(test.date)
		if result != test.expect {
			t.Errorf("Test:%s; Expected:%s; Result:%s", test.name, test.expect, result)
		}
	}
}

func TestStrDateToDate(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		expected time.Time
		retErr   bool
	}{
		{
			name:     "Valid RFC822 date",
			dateStr:  "15 Aug 24 14:45 UTC",
			expected: time.Date(2024, time.August, 15, 14, 45, 0, 0, time.UTC),
			retErr:   false,
		},
		{
			name:     "Another valid RFC822 date",
			dateStr:  "01 Jan 00 00:00 UTC",
			expected: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			retErr:   false,
		},
		{
			name:     "Invalid RFC822 date",
			dateStr:  "31/12/1999 23:59", // Incorrect format
			expected: time.Time{},        // Expect zero time since it's an invalid input
			retErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := StrDateToDate(test.dateStr)
			if (err != nil) != test.retErr {
				t.Errorf("Test: %s; Expected Error: %v; Got Error: %v", test.name, test.retErr, err != nil)
				return
			}
			if !result.Equal(test.expected) {
				t.Errorf("Test: %s; Expected: %s; Result: %s", test.name, test.expected.Format(time.RFC822), result.Format(time.RFC822))
			}
		})
	}
}

func TestGetMinutes(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		expected  int
	}{
		{
			name:      "Exactly 60 minutes difference",
			startTime: time.Date(2024, time.August, 15, 14, 0, 0, 0, time.UTC),
			endTime:   time.Date(2024, time.August, 15, 15, 0, 0, 0, time.UTC),
			expected:  60,
		},
		{
			name:      "30 minutes difference",
			startTime: time.Date(2024, time.August, 15, 14, 0, 0, 0, time.UTC),
			endTime:   time.Date(2024, time.August, 15, 14, 30, 0, 0, time.UTC),
			expected:  30,
		},
		{
			name:      "Negative difference (end time before start time)",
			startTime: time.Date(2024, time.August, 15, 15, 0, 0, 0, time.UTC),
			endTime:   time.Date(2024, time.August, 15, 14, 0, 0, 0, time.UTC),
			expected:  -60,
		},
		{
			name:      "Zero minutes difference",
			startTime: time.Date(2024, time.August, 15, 14, 0, 0, 0, time.UTC),
			endTime:   time.Date(2024, time.August, 15, 14, 0, 0, 0, time.UTC),
			expected:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetMinutes(&test.startTime, &test.endTime)
			if result != test.expected {
				t.Errorf("Test: %s; Expected: %d; Result: %d", test.name, test.expected, result)
			}
		})
	}
}
