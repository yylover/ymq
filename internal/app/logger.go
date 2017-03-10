package app

//Logger 接口
type Logger interface {
	Output(maxdepth int, s string) error
}
