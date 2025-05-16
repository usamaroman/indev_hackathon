package types

type ReservationType string

const (
	Confirmed  ReservationType = "confirmed"
	CheckedIn  ReservationType = "checked_in"
	CheckedOut ReservationType = "checked_out"
)

func (t ReservationType) String() string {
	return string(t)
}

var ReservationTypeToString = map[ReservationType]string{
	Confirmed:  Confirmed.String(),
	CheckedIn:  CheckedIn.String(),
	CheckedOut: CheckedOut.String(),
}

var StringToReservationType = map[string]ReservationType{
	"confirmed":   Confirmed,
	"checked_in":  CheckedIn,
	"checked_out": CheckedOut,
}
