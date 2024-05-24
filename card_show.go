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
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.CARD_SHOW, createCardShow)
}

func createCardShow(data *model.ObjectData) model.IObject {
	res := &cardShow{cards: make([]*Card, 0)}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	CardShow = res
	return res
}

type cardShow struct {
	*object.PointObject
	cards []*Card
}

func (c *cardShow) Update() {
	for i := 0; i < len(c.cards); i++ {
		c.cards[i].Update()
	}
}

func (c *cardShow) CollisionPoint(pos complex128) bool {
	for i := 0; i < len(c.cards); i++ {
		if c.cards[i].CollisionPoint(pos) {
			return true
		}
	}
	return false
}

func (c *cardShow) Init() {
	cardDatas := GetCardDatas()
	for i := 0; i < len(cardDatas); i++ {
		c.cards = append(c.cards, NewCard(cardDatas[i], 0))
	}
	c.tidyCard()
}

func (c *cardShow) Draw(screen *ebiten.Image) {
	for i := 0; i < len(c.cards); i++ {
		c.cards[i].Draw(screen)
	}
}

func (c *cardShow) tidyCard() {
	for i := 0; i < len(c.cards); i++ {
		c.cards[i].Pos = complex(float64(1280-(i+1)*64), 720-64)
	}
}

func (c *cardShow) UpdateClick(pos complex128) {
	if GameManager.LastNum <= 0 {
		Tip.AddTip("没有放置容量")
		return
	}
	for i := 0; i < len(c.cards); i++ {
		if c.cards[i].CollisionPoint(pos) {
			if c.cards[i].CanUse() {
				Cursor.SetCardData(c.cards[i].Data)
			} else {
				Tip.AddTip("费用不足或未完全冷却")
			}
			return
		}
	}
}

func (c *cardShow) RemoveCard(data *CardData) {
	c.cards = utils.FilterSlice(c.cards, func(card *Card) bool {
		return card.Data != data
	})
	c.tidyCard()
}

func (c *cardShow) AddCardData(data *CardData) {
	c.AddCard(NewCard(data, data.CoolTime))
}

func (c *cardShow) AddCard(card *Card) {
	c.cards = append(c.cards, card)
	c.tidyCard()
}

func (c *cardShow) UpdateCard() {
	for i := 0; i < len(c.cards); i++ {
		c.cards[i].UpdateCanUse()
	}
}

func (c *cardShow) ReduceTime(mul float64, add int64) {
	for i := 0; i < len(c.cards); i++ {
		c.cards[i].ReduceTime(mul, add)
	}
}
