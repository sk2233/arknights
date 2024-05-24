/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/layer"
	"GameBase2/utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type PauseLayer struct {
	*layer.UILayer
}

func NewPauseLayer() *PauseLayer {
	res := &PauseLayer{}
	res.UILayer = layer.NewUILayer()
	return res
}

func (u *PauseLayer) Draw(screen *ebiten.Image) {
	utils.FillRect(screen, 0, GameSize, ColorTextMask)
	utils.DrawAnchorText(screen, "PAUSE", GameSize/2, 0.5+0.5i, Font64, colornames.White)
	utils.DrawAnchorText(screen, "----暂停中----", GameSize/2+64i, 0.5+0.5i, Font32, colornames.White)
}

type GameEndLayer struct {
	*layer.UILayer
	title string
}

func NewGameEndLayer(title string) *GameEndLayer {
	res := &GameEndLayer{title: title}
	res.UILayer = layer.NewUILayer()
	return res
}

func (u *GameEndLayer) Draw(screen *ebiten.Image) {
	utils.FillRect(screen, (720i-72i)/2, 1280+72i, ColorTextMask)
	utils.DrawAnchorText(screen, fmt.Sprintf("----%s----", u.title), GameSize/2, 0.5+0.5i, Font32, colornames.White)
}
