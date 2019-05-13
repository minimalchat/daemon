package person

/*
Person interface defines a common function GetFullName that is required for
anything to be considered a Person */
type Person interface {
	GetFullName() string
}
