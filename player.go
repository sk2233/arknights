/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/config"
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	R "arknights/res"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

//==================Player====================

type Player struct {
	img          *ebiten.Image
	Data         *CardData
	Grid         *Grid
	ShowRange    bool
	Range        *Range
	RangeDir     int
	Hp           int64
	skills       map[string]ISkill
	attackTimer  int64
	BuffHolder   *BuffHolder
	BlockEnemies *model.Set[*Enemy]
	dir          float64
	Die          bool
	attack       IAttack
	skillBtns    *SkillButtons
	infos        []*Info
	infoTimer    int
}

func (p *Player) Init() {
	p.Die = false
	PlayerManager.TriggerEvent(&Param{EventType: EventTypePlayerInit, Player: p})
}

func (p *Player) SetRangeDir(dir int) {
	p.RangeDir = dir
	p.dir = utils.If(dir > 1, -1.0, 1.0)
}

func NewPlayer(data *CardData) *Player {
	res := &Player{Data: data, Hp: data.Hp, BuffHolder: NewBuffHolder(),
		BlockEnemies: model.NewSet[*Enemy](), dir: 1, Die: true, infos: make([]*Info, 0)}
	res.BuffHolder.InitData(FieldAtk, float64(data.Atk))
	res.BuffHolder.InitData(FieldDef, float64(data.Def))
	res.BuffHolder.InitData(FieldAttackSpeed, 100)
	res.BuffHolder.InitData(FieldMaxHp, float64(data.Hp))
	res.img = GetMainImg(data.Img)
	res.Range = NewRange(data.Range)
	res.attack = CreatePlayerAttack(data.Attack, res)
	res.skills = make(map[string]ISkill)
	for i := 0; i < len(data.Skills); i++ {
		res.skills[data.Skills[i]] = CreateSkill(data.Skills[i], res)
	}
	return res
}

func (p *Player) GetPlace() int {
	return p.Data.ChangePlace
}

func (p *Player) BeforeDraw(screen *ebiten.Image) {
	for _, skill := range p.skills {
		InvokePosDraw(skill, p.Grid.Pos, screen)
	}
	InvokePosDraw(p.attack, p.Grid.Pos, screen)
	if p.ShowRange {
		if p.Die {
			dir := GetDir(utils.GetCursorPos() - p.Grid.Pos)
			p.Range.DirDraw(screen, p.Grid.Pos, dir)
			p.dir = utils.If(dir > 1, -1.0, 1.0)
		} else {
			p.Range.DirDraw(screen, p.Grid.Pos, p.RangeDir)
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	utils.DrawScaleImage(screen, p.img, p.Grid.Pos, GridSize/2, complex(p.dir, 1))
	utils.DrawAnchorText(screen, p.Data.Name, p.Grid.Pos+GridSize/2, 0.5+0.5i, Font18, colornames.Aqua)
	utils.FillRect(screen, p.Grid.Pos+60i, complex(64*float64(p.Hp)/p.BuffHolder.GetData(FieldMaxHp).Float(), 4), ColorPlayerHp)
}

func (p *Player) Update() {
	if p.Die {
		return
	}
	if p.infoTimer > 0 {
		p.infoTimer--
	} else {
		p.infoTimer = 15
		if len(p.infos) > 0 {
			NewValueTip(p.infos[0].Text, p.infos[0].Color, p.Grid.Pos+GridSize/2)
			p.infos = p.infos[1:]
		}
	}
	p.BuffHolder.Update()
	for _, skill := range p.skills {
		utils.InvokeUpdate(skill)
	}
	if p.attackTimer > 0 {
		p.attackTimer--
	} else {
		p.attackTimer = p.Data.AttackTime * 100 / p.BuffHolder.GetData(FieldAttackSpeed).Int()
		p.attack.Attack()
	}
}

func (p *Player) GetSkills(param *Param) []ISkill {
	res := make([]ISkill, 0)
	for _, skill := range p.skills {
		if skill.CanHandle(param) {
			res = append(res, skill)
		}
	}
	return res
}

func (p *Player) Block(enemy *Enemy) bool {
	if p.BlockEnemies.Size() >= p.Data.BlockNum {
		return false
	}
	p.BlockEnemies.Add(enemy)
	return true
}

func (p *Player) ReleaseEnemy(enemy *Enemy) {
	p.BlockEnemies.Remove(enemy)
}

func (p *Player) Hurt(enemy *Enemy, hurt int64) {
	if p.Die {
		return
	} // 这里攻击已经不能再无效了
	param := &Param{EventType: EventTypePlayerHurt, Player: p, Enemy: enemy, HurtValue: hurt}
	PlayerManager.TriggerEvent(param)
	p.AddInfo(fmt.Sprintf("-%d", param.HurtValue), colornames.Red)
	p.Hp -= param.HurtValue
	if p.Hp <= 0 {
		p.Retreat()
	}
}

func (p *Player) Recover(recover int64) {
	recover = utils.Min(recover, p.BuffHolder.GetData(FieldMaxHp).Int()-p.Hp)
	p.AddInfo(fmt.Sprintf("+%d", recover), colornames.Green)
	p.Hp += recover
}

func (p *Player) Retreat() {
	PlayerManager.TriggerEvent(&Param{EventType: EventTypePlayerRetreat, Player: p})
	for enemy := range p.BlockEnemies.Items() {
		enemy.ReleasePlayer(p)
	}
	p.Die = true
	p.Grid.PopPlayer()
	GameManager.ChangeLastNum(1)
	if p.Data.Career != CareerSummon { // 不是召唤物才返回 卡组
		CardShow.AddCardData(p.Data)
	}
}

func (p *Player) GetEnemies() []*Enemy {
	return p.Range.CollisionEnemies(p.Grid.Pos, p.RangeDir)
}

func (p *Player) ShowSkill() {
	temps := p.GetSkills(&Param{EventType: EventTypePlayerInitiative, Player: p})
	if len(temps) <= 0 {
		ClickManager.Pop()
		Tip.AddTip("该干员没有主动技能")
		return
	}
	skills := make([]IInitiativeSkill, 0)
	for i := 0; i < len(temps); i++ {
		skills = append(skills, temps[i].(IInitiativeSkill))
	}
	p.skillBtns = NewSkillButtons(p.Grid.Pos+32, skills)
	p.Grid.Rank = 2233
	p.ShowRange = true
}

func (p *Player) UpdateClick(pos complex128) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	skill := p.skillBtns.ClickSkill(pos)
	if skill != nil {
		if skill.GetProgress() < 1 {
			Tip.AddTip("该技能还未冷却")
			return
		}
		skill.Action(&Param{EventType: EventTypePlayerInitiative, Player: p})
	}
	ClickManager.Pop()
	p.Grid.Rank = 0
	p.skillBtns.Die = true
	p.skillBtns = nil
	p.ShowRange = false
}

func (p *Player) RemoveSkill(name string) {
	delete(p.skills, name)
}

func (p *Player) AddSkill(name string, skill ISkill) {
	p.skills[name] = skill
}

func (p *Player) GetPlayers() []*Player {
	return p.Range.CollisionPlayers(p.Grid.Pos, p.RangeDir)
}

func (p *Player) RecoverSkill(recover int64) {
	p.AddInfo(fmt.Sprintf("+%d", recover/60), colornames.Blue)
	for _, skill := range p.skills {
		if temp, ok := skill.(IRecover); ok {
			temp.Recover(recover)
		}
	}
}

func (p *Player) AddInfo(info string, clr color.RGBA) {
	p.infos = append(p.infos, &Info{Text: info, Color: clr})
}

//==============PlayerBullet=============

type PlayerBullet struct {
	*object.PointObject
	hurt        int64
	enemy       *Enemy
	player      *Player
	speed       complex128
	die         bool
	adjustTimer int
	Range       bool
}

func (e *PlayerBullet) UpdateRoom(rect model.IRect) {
	e.die = !utils.PointCollision(rect, e.Pos)
}

func (e *PlayerBullet) IsInRoom() bool {
	return !e.die
}

func (e *PlayerBullet) IsDie() bool {
	return e.die
}

func (e *PlayerBullet) Draw(screen *ebiten.Image) {
	utils.FillCircle(screen, e.Pos, 4, 8, config.ColorBlue)
}

func (e *PlayerBullet) Update() {
	e.Pos += e.speed
	if e.enemy.CollisionPoint(e.Pos) {
		e.die = true
		if e.Range {
			pos := ToPos(ToIndex(e.Pos))
			offset := [][]int{{0, 0}}
			enemies := CollisionEnemies(pos, 0, offset)
			for i := 0; i < len(enemies); i++ {
				enemies[i].Hurt(e.player, e.hurt)
			}
			NewRangeTip(pos, 0, offset, ColorCardMask)
		} else {
			e.enemy.Hurt(e.player, e.hurt)
		}
	}
	if e.adjustTimer > 0 {
		e.adjustTimer--
	} else {
		if e.player.Die {
			e.die = true
		} else {
			e.adjustTimer = 15
			e.speed = utils.ScaleVector(utils.UnitVector(e.enemy.Pos-e.Pos), ToMoveSpeed(3))
		}
	}
}

func NewPlayerBullet(pos complex128, hurt int64, enemy *Enemy, player *Player) *PlayerBullet {
	res := &PlayerBullet{hurt: hurt, enemy: enemy, die: false, player: player, adjustTimer: 15, Range: false}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	res.speed = utils.ScaleVector(utils.UnitVector(enemy.Pos-pos), ToMoveSpeed(3))
	utils.AddToLayer(R.LAYER.PLAYER, res)
	return res
}
