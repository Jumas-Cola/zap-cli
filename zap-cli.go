package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/jumas-cola/zap-cli/player"
	"github.com/jumas-cola/zap-cli/ui"
)

func drawButtons(dir string, entries []os.DirEntry, s tcell.Screen, boxStyle tcell.Style) {
	s.Clear()
	width, _ := s.Size()
	btnsCount := width / 20
	if btnsCount < 1 {
		btnsCount = 1
	}
	btnWidth := width / btnsCount
	btnPos := 0

	var currBtn btn

	btns = []btn{}
	startX := 0
	startY := 0
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".mp3" {
			continue
		}
		startX = btnPos * btnWidth
		currBtn = btn{
			x1:   startX + 1,
			y1:   startY,
			x2:   startX + btnWidth - 3,
			y2:   startY + 3,
			path: filepath.Join(dir, e.Name())}
		ui.DrawBox(s, currBtn.x1, currBtn.y1, currBtn.x2, currBtn.y2, boxStyle, e.Name())
		btns = append(btns, currBtn)
		btnPos = (btnPos + 1) % btnsCount
		if btnPos == 0 {
			startY += 5
		}
	}

	ui.DrawBox(s, 1, startY+5, 18, startY+7, boxStyle, "Press q to quit")
}

func checkFilesCount(entries []os.DirEntry) int {
	mp3Count := 0
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".mp3" {
			mp3Count++
		}
	}
	return mp3Count
}

var defStyle tcell.Style

type btn struct {
	x1, y1, x2, y2 int
	path           string
}

var btns []btn

func main() {
	dir := flag.String("dir", "./", "Directory to scan.")
	flag.Parse()

	entries, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatal(err)
	}

	mp3Count := checkFilesCount(entries)
	if mp3Count == 0 {
		log.Fatal("No .mp3 files in directory")
		return
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	drawButtons(*dir, entries, s, boxStyle)

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			drawButtons(*dir, entries, s, boxStyle)
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'q' {
				return
			}
		case *tcell.EventMouse:
			x, y := ev.Position()

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:
				for _, b := range btns {
					if x > b.x1 && x < b.x2 && y > b.y1 && y < b.y2 {
						player.Play(b.path, s)
					}
				}
			}
		}
	}
}
