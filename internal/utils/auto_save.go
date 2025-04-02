package utils

import (
	"log"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

func StartAutoSave(storage *storage.MemStorage, filePath string, interval int, stopChan chan struct{}) {
	if interval == 0 {
		log.Println("Auto-save disabled (STORE_INTERVAL=0)")
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := storage.SaveLoadMetrics(filePath, "save"); err != nil {
				log.Printf("Failed to save metrics to file: %v", err)
			} else {
				log.Println("Metrics successfully saved to a file.")
			}
		case <-stopChan:
			return
		}
	}
}
