package fwatch

import "os"

// EventType - file modify event
type EventType int

const (
	// Create - file create event
	Create EventType = iota + 1
	// Remove - file remove event
	Remove
	// Rename - file rename event
	Rename
	// Write - file write event
	Write
	// Chmod - file chmod event
	Chmod
)

// String - to string
func (et EventType) String() string {
	switch et {
	case Create:
		return "create"
	case Remove:
		return "remove"
	case Rename:
		return "rename"
	case Write:
		return "write"
	case Chmod:
		return "chmod"
	}
	return ""
}

// Eq -
func (et EventType) Eq(t string) bool {
	return t == et.String()
}

// Verbose - show log flag
var Verbose bool = false

// IsService -
var IsService bool = false

// PidFile -
var PidFile *os.File
