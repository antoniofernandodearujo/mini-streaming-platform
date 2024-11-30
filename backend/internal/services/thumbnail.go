package services

import (
	"fmt"
	"os/exec"
)

func GenerateThumbnail(videoPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01.000", "-vframes", "1", outputPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate thumbnail: %v", err)
	}
	return nil
}
