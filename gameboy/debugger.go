package gameboy

type Debugger struct {
	gameboy *Gameboy
}

func NewDebugger(gb *Gameboy) *Debugger {
	return &Debugger{
		gameboy: &Gameboy{},
	}
}
