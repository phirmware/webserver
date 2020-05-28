package views

import "log"

const (
	// AlertLvlError danger class
	AlertLvlError = "danger"
	// AlertLvlWarning warning class
	AlertLvlWarning = "warning"
	// AlertLvlInfo info class
	AlertLvlInfo = "info"
	// AlertLvlSuccess success class
	AlertLvlSuccess = "success"
	// AlertMsgGeneric is a generic message for unknown errors
	AlertMsgGeneric = "Something went wrong. Please try again and contact us if the problem persists."
)

// PublicError defines the error interface for public errors
type PublicError interface {
	error
	Public() string
}

// Alert defines the shape of our alert data
type Alert struct {
	Level   string
	Message string
}

// Data defines the shape of our views data
type Data struct {
	Alert *Alert
	Yield interface{}
}

// SetAlert sets the data field on the data struct
func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// AlertError constructs a custom error message
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
