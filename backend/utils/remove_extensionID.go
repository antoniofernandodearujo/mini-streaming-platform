package utils

import "strings"

func RemoveExtensionID(videoID string) string {
	
	if strings.HasSuffix(videoID, ".mp4") {
		return strings.TrimSuffix(videoID, ".mp4")
	}

	lastDot := strings.LastIndex(videoID, ".")
	if lastDot == -1 {
		return videoID
	}

	return videoID[:lastDot]
}