package main

import (
	"os"
	"os/exec"
	"testing"
)

var (
	mask10000000 = 0xFF
	mask00000010 = 0x02
	mask4        = 0x04
)

func TestInit(t *testing.T) {
	os.Mkdir("test", os.ModePerm)
}

func TestMainBit(t *testing.T) {
	a := mask4

	t.Logf("%b", a)
	a = a >> 1
	t.Logf("%b", a)

	t.Logf("%b", 6)
	t.Logf("4th :%b & %b = %b", 6, 0x0E, (6 & 8)) // 4th bit
	t.Logf("3rd :%b & %b = %b", 6, 0x0D, (6 & 4)) // 3rd bit
	t.Logf("2nd :%b & %b = %b", 6, 0x0B, (6 & 2)) // 2nd bit
	t.Logf("1st :%b & %b = %b", 6, 7, (6 & 1))    // 1st bit
	t.Logf("%b : decalage de 4 :%b : %d,%d", 0xef, (0xef & 0xf0 >> 4), 0xef, (0xef & 0xf0 >> 4))
}

func TestNormalScreenMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "0", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestNormalScreenMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestNormalScreenMode1Dsk(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "1", "-o", "test", "-dsk"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}
func TestNormalScreenMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "2", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestFullScreenMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "0", "-f", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestFullScreenMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "1", "-f", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestFullScreenMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "2", "-f", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}
func TestFullScreenPlusMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "0", "-f", "-p", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestFullScreenPlusMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "1", "-f", "-p", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestFullScreenPlusMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/Batman-Neal-Adams.jpg", "-m", "2", "-f", "-p", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestSpriteMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "0", "-w", "16", "-h", "16", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestSpriteMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "1", "-w", "16", "-h", "16", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestSpriteMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "2", "-w", "16", "-h", "16", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestSpritePlusMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "0", "-w", "16", "-h", "16", "-p", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestSpritePlusMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "1", "-w", "16", "-h", "16", "-p", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestSpritePlusMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "2", "-w", "16", "-h", "16", "-p", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollRraMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "0", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-rra", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollRraMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "1", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-rra", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollRraMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "2", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-rra", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollRLaMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "0", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-rla", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollRLaMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "1", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-rla", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollRLaMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "2", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-rla", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollKeephighMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "0", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-keephigh", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollKeephighMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "1", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-keephigh", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollKeephighMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "2", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-keephigh", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollKeeplowMode0(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "0", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-keeplow", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollKeeplowMode1(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "1", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-keeplow", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestRollKeeplowMode2(t *testing.T) {
	args := []string{"run", "main.go", "-i", "samples/rotate.png", "-m", "2", "-w", "16", "-h", "16", "-roll", "-iter", "16", "-keeplow", "1", "-o", "test"}
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error and gets :%v", err)
	}
}

func TestEnded(t *testing.T) {
	os.RemoveAll("test")
}
