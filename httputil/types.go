package httputil

type Method string

const (
	Get     Method = `GET`
	Post           = `POST`
	Put            = `PUT`
	Delete         = `DELETE`
	Head           = `HEAD`
	Options        = `OPTIONS`
	Patch          = `PATCH`
)
