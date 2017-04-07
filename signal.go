package trader

//Signal - signal of strategy
type Signal struct {
	Strategy    *Strategy
	Init, Close Instruction
}

//Instruction for initiating/closing positions
type Instruction struct {
	BuySell bool
	Size    uint
}
