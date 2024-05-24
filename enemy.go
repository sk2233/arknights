/*
@author: sk
@date: 2023/2/26
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
	"golang.org/x/image/colornames"
)

//========================Enemy=======================

type Enemy struct {
	*object.CollisionObject
	Data         *EnemyData
	path         *object.PathObject
	index        int
	hp           int64
	BuffHolder   *BuffHolder
	Die          bool
	attackTimer  int64
	moveDir      complex128
	BlockPlayer  *Player // 阻塞改敌人的 玩家 优先攻击  死亡时 注意解除阻塞
	attack       IAttack
	TargetPlayer *Player // 攻击目标
	infos        []*Info
	infoTimer    int
}

func (e *Enemy) ReleasePlayer(player *Player) {
	e.BlockPlayer = nil
}

func (e *Enemy) IsDie() bool {
	return e.Die
}

func (e *Enemy) Update() {
	e.CollisionObject.Update()
	if e.infoTimer > 0 {
		e.infoTimer--
	} else {
		e.infoTimer = 15
		if len(e.infos) > 0 {
			NewValueTip(e.infos[0].Text, e.infos[0].Color, e.Pos)
			e.infos = e.infos[1:]
		}
	}
	e.BuffHolder.Update()
	if e.TargetPlayer != nil && e.TargetPlayer.Die { // 清除攻击对象
		e.TargetPlayer = nil
	}
	if e.attackTimer > 0 { // 攻击判断
		e.attackTimer--
	} else {
		e.attackTimer = e.Data.AttackTime * 100 / e.BuffHolder.GetData(FieldAttackSpeed).Int()
		if e.BuffHolder.IsState(StateAttack) {
			e.attack.Attack()
		}
	} // 移动判断
	if e.BlockPlayer != nil || e.TargetPlayer != nil || !e.BuffHolder.IsState(StateMove) { // 阻塞存在  目标存在  或不允许移动
		return
	}
	if e.Data.CanBlock {
		grid := GridManager.GetGrid(ToIndex(e.Pos + e.moveDir*16))
		if grid != nil {
			player := grid.GetPlayer()
			if player != nil && player.Block(e) { // 注意双方死亡时  清楚 阻塞内容
				e.BlockPlayer = player
				return
			}
		}
	}
	moveSpeed := e.BuffHolder.GetData(FieldMoveSpeed).Float()
	if utils.VectorLen2(e.path.GetPoint(e.index)-e.Pos) < moveSpeed*moveSpeed {
		e.Pos = e.path.GetPoint(e.index)
		e.index++
		e.AdjustDir()
	} else {
		e.Pos += utils.ScaleVector(e.moveDir, moveSpeed)
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	e.CollisionObject.Draw(screen)
	utils.DrawAnchorText(screen, e.Data.Name, e.Pos, 0.5+0.5i, Font18, colornames.Red)
	if e.hp < e.Data.Hp {
		utils.FillRect(screen, e.Pos-GridSize/2, complex(64*float64(e.hp)/float64(e.Data.Hp), 4), ColorEnemyHp)
	}
}

func (e *Enemy) AdjustDir() {
	if e.index < e.path.GetLen() {
		e.moveDir = utils.UnitVector(e.path.GetPoint(e.index) - e.Pos)
		if dir := utils.Sign(real(e.moveDir)); dir != 0 {
			e.Scale = complex(dir, 1)
		}
	} else {
		e.Death()
		GameManager.HurtLife(e.Data.HurtLife)
	}
}

func (e *Enemy) Death() {
	e.Die = true
	if e.BlockPlayer != nil {
		e.BlockPlayer.ReleaseEnemy(e)
	}
	EnemyManager.ReduceCount()
}

func (e *Enemy) Hurt(player *Player, hurt int64) {
	if e.Die {
		return
	}
	e.AddInfo(fmt.Sprintf("-%d", hurt), colornames.Red)
	e.hp -= hurt
	if e.hp <= 0 {
		e.Death()
		PlayerManager.TriggerEvent(&Param{EventType: EventTypeEnemyDeath, Player: player, Enemy: e})
	}
}

func (e *Enemy) GetLast() float64 {
	res := utils.VectorLen(e.Pos + e.path.GetPoint(e.index))
	for i := e.index; i < e.path.GetLen()-1; i++ {
		res += e.path.GetDist(i)
	}
	return res
}

func (e *Enemy) AddInfo(info string, clr color.RGBA) {
	e.infos = append(e.infos, &Info{Text: info, Color: clr})
}

func NewEnemy(data *EnemyData, path *object.PathObject) *Enemy {
	res := &Enemy{Data: data, path: path, index: 1, BuffHolder: NewBuffHolder(), hp: data.Hp, Die: false}
	res.CollisionObject = object.NewCollisionObject()
	res.Offset = -GridSize / 2
	res.Size = GridSize
	res.Tag = EnemyTag
	res.BindSprite(GetMainSprite(data.Img), res)
	res.Pos = path.GetPoint(0)
	res.BuffHolder.InitData(FieldAttackSpeed, 100)
	res.BuffHolder.InitData(FieldMoveSpeed, data.MoveSpeed)
	res.BuffHolder.InitData(FieldAtk, float64(data.Atk))
	res.BuffHolder.InitData(FieldDef, float64(data.Def))
	res.attack = CreateEnemyAttack(data.Attack, res)
	res.AdjustDir()
	return res
}

//====================EnemyBullet=======================

type EnemyBullet struct {
	*object.PointObject
	hurt        int64
	player      *Player
	enemy       *Enemy
	speed       complex128
	die         bool
	adjustTimer int
}

func (e *EnemyBullet) UpdateRoom(rect model.IRect) {
	e.die = !utils.PointCollision(rect, e.Pos)
}

func (e *EnemyBullet) IsInRoom() bool {
	return !e.die
}

func (e *EnemyBullet) IsDie() bool {
	return e.die
}

func (e *EnemyBullet) Draw(screen *ebiten.Image) {
	utils.FillCircle(screen, e.Pos, 4, 8, config.ColorRed)
}

func (e *EnemyBullet) Update() {
	e.Pos += e.speed
	if e.player.Grid.CollisionPoint(e.Pos) {
		e.player.Hurt(e.enemy, e.hurt)
		e.die = true
	}
	if e.adjustTimer > 0 {
		e.adjustTimer--
	} else {
		if e.player.Die {
			e.die = true
		} else {
			e.adjustTimer = 15
			e.speed = utils.ScaleVector(utils.UnitVector(e.player.Grid.Pos+GridSize/2-e.Pos), ToMoveSpeed(3))
		}
	}
}

func NewEnemyBullet(pos complex128, hurt int64, player *Player, enemy *Enemy) *EnemyBullet {
	res := &EnemyBullet{hurt: hurt, player: player, die: false, enemy: enemy, adjustTimer: 15}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	res.speed = utils.ScaleVector(utils.UnitVector(player.Grid.Pos+GridSize/2-pos), ToMoveSpeed(3))
	utils.AddToLayer(R.LAYER.ENEMY, res)
	return res
}

//====================EnemyBoom=======================

type EnemyBoom struct {
	*object.PointObject
	hurt        int64
	player      *Player
	enemy       *Enemy
	speed       complex128
	die         bool
	adjustTimer int
	offsets     [][]int
}

func (e *EnemyBoom) UpdateRoom(rect model.IRect) {
	e.die = !utils.PointCollision(rect, e.Pos)
}

func (e *EnemyBoom) IsInRoom() bool {
	return !e.die
}

func (e *EnemyBoom) IsDie() bool {
	return e.die
}

func (e *EnemyBoom) Draw(screen *ebiten.Image) {
	utils.FillCircle(screen, e.Pos, 6, 8, config.ColorPurple)
}

func (e *EnemyBoom) Update() {
	e.Pos += e.speed
	if e.player.Grid.CollisionPoint(e.Pos) {
		pos := ToPos(ToIndex(e.Pos))
		NewRangeTip(pos, 0, e.offsets, ColorCardMask)
		players := CollisionPlayers(pos, 0, e.offsets)
		for i := 0; i < len(players); i++ {
			players[i].Hurt(e.enemy, e.hurt)
		}
		e.die = true
	}
	if e.adjustTimer > 0 {
		e.adjustTimer--
	} else {
		if e.player.Die {
			e.die = true
		} else {
			e.adjustTimer = 15
			e.speed = utils.ScaleVector(utils.UnitVector(e.player.Grid.Pos+GridSize/2-e.Pos), ToMoveSpeed(3))
		}
	}
}

func NewEnemyBoom(pos complex128, hurt int64, player *Player, enemy *Enemy) *EnemyBoom {
	res := &EnemyBoom{hurt: hurt, player: player, die: false, enemy: enemy, adjustTimer: 15, offsets: [][]int{{-1, 0},
		{1, 0}, {0, 1}, {0, -1}, {0, 0}}}
	res.PointObject = object.NewPointObject()
	res.Pos = pos
	res.speed = utils.ScaleVector(utils.UnitVector(player.Grid.Pos+GridSize/2-pos), ToMoveSpeed(3))
	utils.AddToLayer(R.LAYER.ENEMY, res)
	return res
}
