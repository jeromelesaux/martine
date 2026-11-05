package main_test

import (
	"os"
	"os/exec"
	"testing"
)

const (
	flagIn           = "-in"
	flagFullscreen   = "-fullscreen"
	flagPlus         = "-plus"
	flagHeight       = "-height"
	flagIter         = "-iter"
	flagRla          = "-rla"
	flagKeephigh     = "-keephigh"
	flagKeeplow      = "-keeplow"
	flagMode         = "-mode"
	flagWidth        = "-width"
	flagRoll         = "-roll"
	flagNoUI         = "-noui"
	flagRra          = "-rra"
	flagOut          = "-out"
	exportHandler    = "export_handler.go"
	sampleBatman     = "../../samples/Batman-Neal-Adams.jpg"
	sampleRotate     = "../../samples/rotate.png"
	cmdGo            = "go"
	cmdRun           = "run"
	cmdMainGo        = "main.go"
	cmdProcessGo     = "process.go"
	errFmt           = "Expected no error and gets :%v"
)

// mask10000000 = 0xFF
// mask00000010 = 0x02
var mask4 = 0x04
var output = "../../test"

func TestMainBit(t *testing.T) {
	a := mask4

	t.Logf("%b", a)
	a >>= 1
	t.Logf("%b", a)

	t.Logf("%b", 6)
	t.Logf("4th :%b & %b = %b", 6, 0x0E, (6 & 8)) // 4th bit
	t.Logf("3rd :%b & %b = %b", 6, 0x0D, (6 & 4)) // 3rd bit
	t.Logf("2nd :%b & %b = %b", 6, 0x0B, (6 & 2)) // 2nd bit
	t.Logf("1st :%b & %b = %b", 6, 7, (6 & 1))    // 1st bit
	t.Logf("%b : decalage de 4 :%b : %d,%d", 0xef, (0xef & 0xf0 >> 4), 0xef, (0xef & 0xf0 >> 4))
}

func TestNormalScreenMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "0", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestNormalScreenMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestNormalScreenMode1Dsk(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "1", flagOut, output, "-dsk"}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestNormalScreenMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "2", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestFullScreenMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "0", flagFullscreen, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestFullScreenMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "1", flagFullscreen, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestFullScreenMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "2", flagFullscreen, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestFullScreenPlusMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "0", flagFullscreen, flagPlus, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestFullScreenPlusMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "1", flagFullscreen, flagPlus, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestFullScreenPlusMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleBatman, flagMode, "2", flagFullscreen, flagPlus, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestSpriteMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "0", flagWidth, "16", flagHeight, "16", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestSpriteMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "1", flagWidth, "16", flagHeight, "16", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestSpriteMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "2", flagWidth, "16", flagHeight, "16", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestSpritePlusMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "0", flagWidth, "16", flagHeight, "16", flagPlus, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestSpritePlusMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "1", flagWidth, "16", flagHeight, "16", flagPlus, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestSpritePlusMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "2", flagWidth, "16", flagHeight, "16", flagPlus, flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollRraMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "0", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagRra, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollRraMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "1", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagRra, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollRraMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "2", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagRra, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollRLaMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "0", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagRla, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollRLaMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "1", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagRla, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollRLaMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "2", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagRla, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollKeephighMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "0", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagKeephigh, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollKeephighMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "1", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagKeephigh, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollKeephighMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "2", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagKeephigh, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollKeeplowMode0(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "0", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagKeeplow, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollKeeplowMode1(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "1", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagKeeplow, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}

func TestRollKeeplowMode2(t *testing.T) {
	args := []string{cmdRun, cmdMainGo, cmdProcessGo, exportHandler, flagNoUI, flagIn, sampleRotate, flagMode, "2", flagWidth, "16", flagHeight, "16", flagRoll, flagIter, "16", flagKeeplow, "1", flagOut, output}
	cmd := exec.Command(cmdGo, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
}
