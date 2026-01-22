package commons

import "time"

func FormatTimeString(timeStr string) string {

	t, _ := time.Parse(time.RFC3339Nano, timeStr)
	newFTime := t.Format("2006-01-02 15:04")

	return newFTime
}
