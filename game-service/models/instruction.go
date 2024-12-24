package models

// Instruction represents the message structure for communication between TV and Phone
type Instruction struct {
	Role    string `json:"role"`    // "tv" or "phone"
	Command string `json:"command"` // Command to execute
}
