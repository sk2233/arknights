/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/config"
	"GameBase2/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

//===================Card====================

type Card struct {
	Pos       complex128
	img       *ebiten.Image
	Data      *CardData
	coolTimer int64
	canUse    bool
}

func (c *Card) GetMin() complex128 {
	return c.Pos
}

func (c *Card) GetMax() complex128 {
	return c.Pos + GridSize
}

func (c *Card) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(c, pos)
}

func (c *Card) Draw(screen *ebiten.Image) {
	utils.DrawImage(screen, c.img, c.Pos)
	utils.FillRect(screen, c.Pos+64/4, 32+18i, ColorTextMask)
	utils.DrawAnchorText(screen, strconv.Itoa(c.Data.CostNum), c.Pos+32+18i/2, 0.5+0.5i, Font18, colornames.White)
	utils.StrokeRect(screen, c.Pos, GridSize, 1, config.ColorWhite)
	utils.DrawAnchorText(screen, c.Data.Name, c.Pos+GridSize/2, 0.5+0.5i, Font18, colornames.Aqua)
	if c.coolTimer > 0 {
		utils.FillRect(screen, c.Pos, GridSize, ColorCardMask)
		utils.DrawAnchorText(screen, strconv.FormatFloat(float64(c.coolTimer)/60, 'f', 1, 64), c.Pos+GridSize/2, 0.5+0.5i, Font18, colornames.White)
	} else if !c.canUse {
		utils.FillRect(screen, c.Pos, GridSize, ColorTextMask)
	}
}

func (c *Card) Update() {
	c.coolTimer--
}

func (c *Card) CanUse() bool {
	return c.coolTimer <= 0 && c.canUse
}

func (c *Card) UpdateCanUse() {
	c.canUse = GameManager.Point >= c.Data.CostNum && GameManager.LastNum > 0
}

func (c *Card) ReduceTime(mul float64, add int64) {
	c.coolTimer -= int64(float64(c.Data.CoolTime)*mul) + add
}

func NewCard(data *CardData, coolTimer int64) *Card {
	res := &Card{Data: data, Pos: 0, coolTimer: coolTimer}
	res.img = GetMainImg(data.Img)
	return res
}

//================SkillButton==================

type SkillButton struct {
	Skill IInitiativeSkill
	Pos   complex128
	Size  complex128
}

func (s *SkillButton) GetMin() complex128 {
	return s.Pos
}

func (s *SkillButton) GetMax() complex128 {
	return s.Pos + s.Size
}

func (s *SkillButton) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(s, pos)
}

func (s *SkillButton) Draw(screen *ebiten.Image) {
	utils.FillRect(screen, s.Pos, s.Size, config.ColorRed)
	utils.FillRect(screen, s.Pos, complex(real(s.Size)*s.Skill.GetProgress(), imag(s.Size)), config.ColorBlue)
	utils.DrawAnchorText(screen, s.Skill.GetName(), s.Pos+s.Size/2, 0.5+0.5i, Font18, colornames.White)
	utils.StrokeRect(screen, s.Pos, s.Size, 1, config.ColorWhite)
}

func NewSkillButton(skill IInitiativeSkill) *SkillButton {
	res := &SkillButton{Skill: skill, Pos: 0}
	bound := text.BoundString(Font18, skill.GetName())
	res.Size = utils.Int2Vector(bound.Dx()+8, bound.Dy()+8) // 留 4 边距
	return res
}
