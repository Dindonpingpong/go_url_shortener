package secretary

type Secretary interface {
	Encode(data string) (string)
	Decode(msg string) (string, error)
}