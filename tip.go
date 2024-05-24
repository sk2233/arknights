/*
@author: sk
@date: 2023/2/26
*/
package main

import (
	"GameBase2/config"
	"GameBase2/factory"
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	R "arknights/res"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

//===================Tip====================

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.TIP, createTip)
}

func createTip(data *model.ObjectData) model.IObject {
	res := &tip{}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	Tip = res
	return res
}

type Msg struct {
	Text   string
	Offset float64
}

type tip struct {
	*object.PointObject
	msgs []*Msg
}

func (t *tip) Update() {
	for i := 0; i < len(t.msgs); i++ {
		t.msgs[i].Offset += 2
	}
	if len(t.msgs) > 0 && t.msgs[0].Offset > 360 { // 只有第一个有出界可能性
		t.msgs = t.msgs[1:]
	}
}

func (t *tip) Draw(screen *ebiten.Image) {
	for _, msg := range t.msgs {
		utils.DrawAnchorText(screen, msg.Text, GameSize/2-complex(0, msg.Offset), 0.5+0.5i, Font48, colornames.White)
	}
}

func (t *tip) AddTip(msg string) {
	t.msgs = append(t.msgs, &Msg{Text: msg, Offset: 0})
	for i := len(t.msgs) - 2; i >= 0; i-- { // 调整位置
		if t.msgs[i].Offset < t.msgs[i+1].Offset+48 {
			t.msgs[i].Offset = t.msgs[i+1].Offset + 48
		}
	}
	if t.msgs[0].Offset > 360 { // 只有第一个有出界可能性
		t.msgs = t.msgs[1:]
	}
}

//=====================ValueTip===================

type ValueTip struct {
	*object.PointObject
	show  string
	clr   color.Color
	timer int
}

func (v *ValueTip) IsDie() bool {
	return v.timer <= 0
}

func (v *ValueTip) Update() {
	v.timer--
	v.Pos -= 64i / 60
}

func (v *ValueTip) Draw(screen *ebiten.Image) {
	utils.DrawAnchorText(screen, v.show, v.Pos, 0.5+0.5i, Font18, v.clr)
}

func NewValueTip(show string, clr color.Color, pos complex128) *ValueTip {
	res := &ValueTip{show: show, clr: clr, timer: 60}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	utils.AddToLayer(R.LAYER.UI, res)
	return res
}

//=====================RangeTip===================

type RangeTip struct {
	*object.PointObject
	offsets [][]int
	dir     int
	clr     *model.Color
	timer   int
}

func (v *RangeTip) IsDie() bool {
	return v.timer <= 0
}

func (v *RangeTip) Update() {
	v.timer--
}

func (v *RangeTip) Draw(screen *ebiten.Image) {
	DrawRange(screen, v.Pos, v.dir, v.offsets, v.clr)
}

func NewRangeTip(pos complex128, dir int, offsets [][]int, clr *model.Color) *RangeTip {
	res := &RangeTip{dir: dir, offsets: offsets, clr: clr, timer: 10}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	utils.AddToLayer(R.LAYER.UI, res)
	return res
}
