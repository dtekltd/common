package mail

import (
	"time"

	"github.com/dtekltd/common/system"
)

var limitMail int

/**
 * ticker.Stop()
 * done <- true
 */
func ScheduleSend(interval time.Duration, limit int) (*time.Ticker, chan<- bool) {
	limitMail = limit
	ticker := time.NewTicker(interval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				system.Logger.Tracef("Task send mail at: %s", t.String())
				taskSend()
			}
		}
	}()

	system.Logger.Info("Mailing Send mail task was scheduled with interval: ", interval)

	return ticker, done
}

func taskSend() {
	instances := []Instance{}
	db.Model(Instance{}).
		Preload("Message").
		Order("priority ASC").
		Limit(limitMail).
		Find(&instances, "sent_at=0 AND failed_count < 3")
	if len(instances) > 0 {
		failed := 0
		success := 0
		for _, ins := range instances {
			if err := SendInstance(&ins); err != nil {
				system.Logger.Errorf("Mailing sent failed: ", err.Error())
				failed += 1
			} else {
				success += 1
			}
		}
		system.Logger.Infof("Mailing sent: %d, success: %d, failed: %d", len(instances), success, failed)
	}
}
