package context

type Key string

func (c Key) String() string {
	return string(c)
}

const (
	ContextKeyDatastore Key = Key("datastore")
	ContextKeyCache     Key = Key("cache")
	ContextKeyClientID  Key = Key("clientID")
)
