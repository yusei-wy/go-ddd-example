package env

type Stage int

const (
	Debug Stage = iota
	Local
	Staging
	Production
)

func NewStage(stage string) Stage {
	switch stage {
	case "dbg":
		return Debug
	case "local":
		return Local
	case "stg":
		return Staging
	case "prod":
		return Production
	default:
		return Debug
	}
}
