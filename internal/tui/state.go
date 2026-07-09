package tui

type State int

const (
    MenuState State = iota
    FormState
    RunningState
    FinishedState
)