/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/config"
	"GameBase2/factory"
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	R "arknights/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.CURSOR, createCursor)
}

func createCursor(data *model.ObjectData) model.IObject {
	res := &cursor{}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	Cursor = res
	return res
}

type cursor struct {
	*object.PointObject
	data *CardData
	img  *ebiten.Image
}

func (c *cursor) UpdateClick(pos complex128) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		c.endCursor()
		return
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	if CardShow.CollisionPoint(pos) { // 返回卡牌
		c.endCursor()
		return
	}
	x, y := ToIndex(pos) // 尝试放下
	if x >= 0 && x < 20 && y >= 0 && y < 9 {
		grid := GridManager.GetGrid(x, y)
		if grid.GetPlace()&c.data.Place == 0 {
			Tip.AddTip("该位置无法放置")
			return
		}
		CardShow.RemoveCard(c.data)
		player := PlayerManager.CreatePlayer(c.data)
		GameManager.ChangePoint(-c.data.CostNum)
		GameManager.ChangeLastNum(-1)
		c.endCursor()
		grid.PushPlayer(player)
	}
}

func (c *cursor) Draw(screen *ebiten.Image) {
	if c.data == nil {
		return
	}
	pos := utils.GetCursorPos() - GridSize/2
	x, y := ToIndex(pos + GridSize/2)
	if x >= 0 && x < 20 && y >= 0 && y < 9 {
		grid := GridManager.GetGrid(x, y)
		if grid.GetPlace()&c.data.Place > 0 {
			pos = ToPos(x, y)
		}
	}
	utils.DrawImage(screen, c.img, pos)
}

func (c *cursor) SetCardData(cardData *CardData) {
	c.data = cardData
	c.img = GetMainImg(cardData.Img)
	ClickManager.Push(c)
	GridManager.ShowHint(cardData.Place)
}

func (c *cursor) endCursor() {
	c.data = nil
	c.img = nil
	ClickManager.Pop()
	GridManager.HideHint()
}
