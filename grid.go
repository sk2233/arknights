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
	config.ObjectFactory.RegisterPointFactory(R.CLASS.GRID, createGrid)
}

func createGrid(data *model.ObjectData) model.IObject {
	res := &Grid{players: model.NewStack[*Player](), showHint: false}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	return res
}

type Grid struct {
	*object.PointObject
	place    int
	players  *model.Stack[*Player]
	showHint bool
	Rank     int
}

func (g *Grid) GetMin() complex128 {
	return g.Pos
}

func (g *Grid) GetMax() complex128 {
	return g.Pos + GridSize
}

func (g *Grid) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(g, pos)
}

func (g *Grid) Update() {
	items := g.players.Items()
	for i := 0; i < len(items); i++ {
		items[i].Update()
	}
}

func (g *Grid) Order() int {
	return g.Rank
}

func (g *Grid) UpdateClick(pos complex128) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	player := g.players.Peek()
	player.SetRangeDir(GetDir(pos - g.Pos))
	player.ShowRange = false
	player.Init()
	g.Rank = 0
	ClickManager.Pop()
}

func (g *Grid) Init() {
	g.place = g.GetInt(R.PROP.PLACE, PlaceHighland) // 默认高地
}

func (g *Grid) BeforeDraw(screen *ebiten.Image) {
	g.drawSelf(screen)
	items := g.players.Items()
	for i := 0; i < len(items); i++ {
		items[i].BeforeDraw(screen)
	}
}

func (g *Grid) Draw(screen *ebiten.Image) {
	items := g.players.Items()
	for i := 0; i < len(items); i++ {
		items[i].Draw(screen)
	}
}

func (g *Grid) GetPlace() int {
	if g.players.IsEmpty() {
		return g.place
	}
	return g.players.Peek().GetPlace()
}

func (g *Grid) drawSelf(screen *ebiten.Image) {
	offset := 1 + 1i
	switch g.place {
	case PlaceLand:
		utils.FillRect(screen, g.Pos+offset, GridSize-offset*2, ColorPlaceLand)
	case PlaceHighland:
		utils.FillRect(screen, g.Pos+offset, GridSize-offset*2, ColorPlaceHighLand)
	case PlaceStart:
		utils.StrokeRect(screen, g.Pos+offset, GridSize-offset*2, 2, config.ColorRed)
		utils.DrawLine(screen, g.Pos+offset, g.Pos+GridSize-offset, 2, config.ColorRed)
		utils.DrawLine(screen, g.Pos+62+offset, g.Pos+62i+offset, 2, config.ColorRed)
	case PlaceEnd:
		utils.StrokeRect(screen, g.Pos+offset, GridSize-offset*2, 2, config.ColorBlue)
		utils.StrokeCircle(screen, g.Pos+GridSize/2, 16, 12, 2, config.ColorBlue)
	}
	if g.showHint {
		utils.FillRect(screen, g.Pos+offset, GridSize-offset*2, ColorHint)
	}
}

func (g *Grid) ShowHint(place int) {
	g.showHint = g.GetPlace()&place > 0
}

func (g *Grid) HideHint() {
	g.showHint = false
}

func (g *Grid) PushPlayer(player *Player) {
	player.Grid = g
	player.ShowRange = true
	g.players.Push(player)
	g.Rank = 2233
	ClickManager.Push(g)
}

func (g *Grid) GetPlayer() *Player {
	if g.players.IsEmpty() {
		return nil
	}
	return g.players.Peek()
}

func (g *Grid) PopPlayer() {
	PlayerManager.RemovePlayer(g.players.Pop())
}
