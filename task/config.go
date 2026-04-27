package task

type ConfigTemplate interface {
	ParseFromInterface(config map[string]interface{}) any
	GetDefault() any
}
