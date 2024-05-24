/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/config"
	"GameBase2/object"
	"GameBase2/utils"
	R "arknights/res"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

//=====================Button==================

type Button struct {
	*object.PointObject
	size complex128
	Text string
}

func (b *Button) GetMin() complex128 {
	return b.Pos
}

func (b *Button) GetMax() complex128 {
	return b.Pos + b.size
}

func (b *Button) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(b, pos)
}

func NewButton(pos complex128, size complex128, text string) *Button {
	res := &Button{size: size, Text: text}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	utils.AddToLayer(R.LAYER.UI, res)
	return res
}

func (b *Button) Draw(screen *ebiten.Image) {
	utils.DrawRect(screen, b.Pos, b.size, 2, ColorButton, config.ColorWhite)
	utils.DrawAnchorText(screen, b.Text, b.Pos+b.size/2, 0.5+0.5i, Font32, colornames.White)
}

//========================SkillButtons========================

type SkillButtons struct {
	*object.PointObject
	btns []*SkillButton
	Die  bool
}

func (s *SkillButtons) Draw(screen *ebiten.Image) {
	for i := 0; i < len(s.btns); i++ {
		s.btns[i].Draw(screen)
	}
}

func (s *SkillButtons) IsDie() bool {
	return s.Die
}

func (s *SkillButtons) adjustBtns(pos complex128) {
	w, h := 0.0, 0.0
	for i := 0; i < len(s.btns); i++ {
		if i > 0 {
			w += 4 // 间隔 4
		}
		w += real(s.btns[i].Size)
		h = imag(s.btns[i].Size)
	}
	x, y := real(pos)-w/2, imag(pos)-h
	for i := 0; i < len(s.btns); i++ {
		s.btns[i].Pos = complex(x, y)
		x += real(s.btns[i].Size) + 4
	}
}

func (s *SkillButtons) ClickSkill(pos complex128) IInitiativeSkill {
	for i := 0; i < len(s.btns); i++ {
		if s.btns[i].CollisionPoint(pos) {
			return s.btns[i].Skill
		}
	}
	return nil
}

func NewSkillButtons(pos complex128, skills []IInitiativeSkill) *SkillButtons {
	res := &SkillButtons{btns: make([]*SkillButton, 0), Die: false}
	res.PointObject = object.NewPointObject()
	for i := 0; i < len(skills); i++ {
		res.btns = append(res.btns, NewSkillButton(skills[i]))
	}
	res.adjustBtns(pos)
	utils.AddToLayer(R.LAYER.UI, res)
	return res
}
