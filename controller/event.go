package controller

type event interface {
	Type() string
}
