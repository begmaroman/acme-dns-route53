package notifier

// Notifier represents interface for sending notification
type Notifier interface {
	// Notify sends a notification with a given topic and message
	Notify(topic, message string) error
}
