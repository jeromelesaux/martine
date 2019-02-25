package gfx

type Size struct {
	Width int
	High int
	LinesNumber int 
	ColumnsNumber int
}

var (
	Mode0 = Size{Width:160,High:200, LinesNumber:200, ColumnsNumber:80}
	Mode1 = Size{Width:320,High:200, LinesNumber:200, ColumnsNumber:80}
	Mode2 = Size{Width:640,High:200, LinesNumber:200, ColumnsNumber:80}
	Overscan = Size{Width:640, High:400, LinesNumber:272, ColumnsNumber: 96}
)
