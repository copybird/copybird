package notifier

type notifier interface {
	sendNotification(string) error
}

