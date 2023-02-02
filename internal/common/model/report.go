package model

type Event struct {
	Time  int64  `json:"time"`
	Msg   string `json:"msg"`
	Level string `json:"level"`
}
