/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	R "arknights/res"
	"fmt"
	"reflect"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//===================gameManager======================

type gameManager struct {
	Life     int
	Point    int     // 最大  99
	Progress float64 // 0 ~ 1
	LastNum  int     // 剩余可放置角色
}

func (g *gameManager) Update() {
	if !StackRoom.IsEmpty() { // 暂停状态 判断
		return
	}
	// 正常逻辑
	g.Progress += 0.01 * GameSpeed
	if g.Progress > 1 {
		g.Progress = 0
		g.Point = utils.If(g.Point+1 > 99, 99, g.Point+1)
		CardShow.UpdateCard()
	}
}

func (g *gameManager) HurtLife(life int) {
	g.Life -= life
	if g.Life <= 0 {
		StackRoom.PushLayer(NewGameEndLayer("任务失败"))
	}
}

func (g *gameManager) ChangePoint(value int) {
	g.Point = utils.If(g.Point+value > 99, 99, g.Point+value)
}

func (g *gameManager) ChangeLastNum(value int) {
	g.LastNum += value
	CardShow.UpdateCard()
}

func NewGameManager() *gameManager {
	return &gameManager{Point: 10, Progress: 0, LastNum: 10, Life: 15}
}

//====================clickManager===================

type clickManager struct { // 使用栈管理点击事件  防止 点击事件随意触发
	clicks             *model.Stack[IUpdateClick]
	speedBtn, pauseBtn *Button
}

func (g *clickManager) Update() {
	pos := utils.GetCursorPos()
	if !g.clicks.IsEmpty() { // 有 优先级 更高的先执行
		g.clicks.Peek().UpdateClick(pos)
		return
	}
	// 通常事件  现在都是 需要鼠标 左键
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	if g.speedBtn.CollisionPoint(pos) {
		GameSpeedIndex++
		GameSpeed = GameSpeeds[GameSpeedIndex%len(GameSpeeds)]
		g.speedBtn.Text = fmt.Sprintf("X%.0f", GameSpeed)
		return
	}
	if g.pauseBtn.CollisionPoint(pos) {
		if StackRoom.IsEmpty() {
			StackRoom.PushLayer(NewPauseLayer())
			g.pauseBtn.Text = "RESUME"
		} else {
			StackRoom.PopLayer()
			g.pauseBtn.Text = "PAUSE"
		}
		return
	}
	CardShow.UpdateClick(pos)      // 选择 卡牌
	PlayerManager.UpdateClick(pos) // 选择干员，创建技能显示技能事件
}

func (g *clickManager) Init() {
	g.speedBtn = NewButton(1280-256, 128+32i, "X1")
	g.pauseBtn = NewButton(1280-128, 128+32i, "PAUSE")
}

func (g *clickManager) Push(updateClick IUpdateClick) {
	g.clicks.Push(updateClick)
}

func (g *clickManager) Pop() {
	g.clicks.Pop()
}

func NewClickManager() *clickManager {
	return &clickManager{clicks: model.NewStack[IUpdateClick]()}
}

//=================gridManager===================

type gridManager struct {
	grids [][]*Grid
}

func (g *gridManager) Init() {
	g.grids = make([][]*Grid, 20)
	for i := 0; i < 20; i++ {
		g.grids[i] = make([]*Grid, 9)
	}
	temps := utils.GetObjectLayer(R.LAYER.PLAYER).GetObjectsByType(reflect.TypeOf(&Grid{}))
	for i := 0; i < len(temps); i++ {
		grid := temps[i].(*Grid)
		x, y := ToIndex(grid.Pos)
		g.grids[x][y] = grid
	}
}

func (g *gridManager) ShowHint(place int) {
	for i := 0; i < len(g.grids); i++ {
		for j := 0; j < len(g.grids[i]); j++ {
			g.grids[i][j].ShowHint(place)
		}
	}
}

func (g *gridManager) GetGrid(x, y int) *Grid {
	return g.grids[x][y]
}

func (g *gridManager) HideHint() {
	for i := 0; i < len(g.grids); i++ {
		for j := 0; j < len(g.grids[i]); j++ {
			g.grids[i][j].HideHint()
		}
	}
}

func (g *gridManager) GetGrids(filter func(*Grid) bool) []*Grid {
	res := make([]*Grid, 0)
	for i := 0; i < len(g.grids); i++ {
		for j := 0; j < len(g.grids[i]); j++ {
			if filter(g.grids[i][j]) {
				res = append(res, g.grids[i][j])
			}
		}
	}
	return res
}

func NewGridManager() *gridManager {
	return &gridManager{}
}

//======================playerManager=======================

type playerManager struct {
	players []*Player
}

func (m *playerManager) TriggerEvent(param *Param) {
	skills := make([]ISkill, 0)
	for i := 0; i < len(m.players); i++ {
		skills = append(skills, m.players[i].GetSkills(param)...)
	} // 后续 有 排序 需求  在这里 排序
	sort.Slice(skills, func(i, j int) bool {
		return utils.InvokeOrder(skills[i]) < utils.InvokeOrder(skills[j])
	})
	for i := 0; i < len(skills); i++ {
		skills[i].Action(param)
	}
}

func (m *playerManager) CreatePlayer(data *CardData) *Player {
	res := NewPlayer(data)
	m.players = append(m.players, res)
	return res
}

func (m *playerManager) RemovePlayer(player *Player) {
	m.players = utils.RemoveSliceItem(m.players, player)
}

func (m *playerManager) GetPlayers(filter func(player *Player) bool) []*Player {
	return utils.FilterSlice(m.players, filter)
}

func (m *playerManager) UpdateClick(pos complex128) {
	for i := 0; i < len(m.players); i++ {
		if m.players[i].Grid.CollisionPoint(pos) {
			ClickManager.Push(m.players[i])
			m.players[i].ShowSkill()
			return
		}
	}
}

func (m *playerManager) GetPlayer(filter func(player *Player) bool) *Player {
	for i := 0; i < len(m.players); i++ {
		if filter(m.players[i]) {
			return m.players[i]
		}
	}
	return nil
}

func NewPlayerManager() *playerManager {
	return &playerManager{players: make([]*Player, 0)}
}

//=====================enemyManager=====================

type enemyManager struct {
	enemyWaves            [][]string
	enemyDatas            map[string]*EnemyData
	WaveIndex, enemyIndex int
	timer                 int
	path                  *object.PathObject
	enemyCount            int
}

func (e *enemyManager) Init() {
	e.path = utils.GetObjectLayer(R.LAYER.ENEMY).GetObject(R.OBJECT.PATH).(*object.PathObject)
}

func (e *enemyManager) Update() {
	if !StackRoom.IsEmpty() || e.WaveIndex >= len(e.enemyWaves) { // 暂停状态 判断
		return
	}
	if e.timer > 0 {
		e.timer--
		if e.timer%30 == 0 && e.enemyIndex < len(e.enemyWaves[e.WaveIndex]) {
			data := e.enemyDatas[e.enemyWaves[e.WaveIndex][e.enemyIndex]]
			utils.AddToLayer(R.LAYER.ENEMY, NewEnemy(data, e.path))
			e.enemyIndex++
			e.enemyCount++
		}
	} else {
		e.WaveIndex++
		e.enemyIndex = 0
		e.timer = 60 * 9
	}
}

func (e *enemyManager) GetAllWave() int {
	return len(e.enemyWaves)
}

func (e *enemyManager) ReduceCount() {
	e.enemyCount--
	if e.WaveIndex >= len(e.enemyWaves) && e.enemyCount == 0 {
		StackRoom.PushLayer(NewGameEndLayer("游戏胜利"))
	}
}

func NewEnemyManager() *enemyManager {
	return &enemyManager{enemyWaves: GetEnemyWaves(), WaveIndex: -1, enemyDatas: GetEnemiesDatas(), enemyCount: 0}
}
