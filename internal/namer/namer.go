package namer

type ImageNamer interface {
	GenerateFilename(string) (string, error)
}
