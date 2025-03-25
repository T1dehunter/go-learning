package types

type EventType string

const (
	ServerEvent EventType = "server"
	ClientEvent EventType = "client"
)

type LogEvent struct {
	Type      EventType
	Title     string
	CreatedAt string
}
