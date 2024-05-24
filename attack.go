/*
@author: sk
@date: 2023/3/4
*/
package main

import (
	"GameBase2/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	enemyAttackFactories  = make(map[string]func(*Enemy) IAttack)
	playerAttackFactories = make(map[string]func(*Player) IAttack)
)

func CreateEnemyAttack(name string, enemy *Enemy) IAttack {
	return enemyAttackFactories[name](enemy)
}

func CreatePlayerAttack(name string, player *Player) IAttack {
	return playerAttackFactories[name](player)
}

//********************EnemyAttack**********************

func init() {
	enemyAttackFactories["近战攻击"] = createEnemyMeleeAttack
	enemyAttackFactories["远程攻击"] = createEnemyRemoteAttack
	enemyAttackFactories["远程群攻"] = createEnemyGroupAttack
}

//===================EnemyGroupAttack==================

func createEnemyGroupAttack(enemy *Enemy) IAttack {
	return &EnemyGroupAttack{enemy: enemy}
}

type EnemyGroupAttack struct {
	enemy *Enemy
}

func (e *EnemyGroupAttack) Attack() {
	enemy := e.enemy
	player := enemy.TargetPlayer
	if player == nil {
		player = enemy.BlockPlayer
	}
	if player == nil { // 先尝试攻击 目标对象 再是被阻挡的  再 攻击  最近的
		player = GetMinPlayer(enemy.Pos, enemy.Data.AttackRange)
	}
	if player == nil {
		return
	}
	param := &Param{EventType: EventTypeEnemyAttack, Player: player, Enemy: enemy,
		Atk: enemy.BuffHolder.CloneData(FieldAtk), Def: player.BuffHolder.CloneData(FieldDef)}
	PlayerManager.TriggerEvent(param)
	if param.Invalid {
		return
	}
	enemy.TargetPlayer = player
	NewEnemyBoom(enemy.Pos, ParseHurt(param), player, e.enemy)
}

//====================EnemyRemoteAttack==================

func createEnemyRemoteAttack(enemy *Enemy) IAttack {
	return &EnemyRemoteAttack{enemy: enemy}
}

type EnemyRemoteAttack struct {
	enemy *Enemy
}

func (e *EnemyRemoteAttack) Attack() {
	enemy := e.enemy
	player := enemy.TargetPlayer
	if player == nil {
		player = enemy.BlockPlayer
	}
	if player == nil { // 先尝试攻击 目标对象 再是被阻挡的  再 攻击  最近的
		player = GetMinPlayer(enemy.Pos, enemy.Data.AttackRange)
	}
	if player == nil {
		return
	}
	param := &Param{EventType: EventTypeEnemyAttack, Player: player, Enemy: enemy,
		Atk: enemy.BuffHolder.CloneData(FieldAtk), Def: player.BuffHolder.CloneData(FieldDef)}
	PlayerManager.TriggerEvent(param)
	if param.Invalid {
		return
	}
	enemy.TargetPlayer = player
	NewEnemyBullet(enemy.Pos, ParseHurt(param), player, e.enemy)
}

//====================MeleeAttack======================

func createEnemyMeleeAttack(enemy *Enemy) IAttack {
	return &EnemyMeleeAttack{enemy: enemy}
}

type EnemyMeleeAttack struct {
	enemy *Enemy
}

func (m *EnemyMeleeAttack) Attack() { // 近战斗 仅攻击阻塞对象
	if m.enemy.BlockPlayer == nil {
		return
	}
	enemy := m.enemy
	player := enemy.BlockPlayer
	param := &Param{EventType: EventTypeEnemyAttack, Player: player, Enemy: enemy,
		Atk: enemy.BuffHolder.CloneData(FieldAtk), Def: player.BuffHolder.CloneData(FieldDef)}
	PlayerManager.TriggerEvent(param)
	if param.Invalid {
		return
	}
	player.Hurt(m.enemy, ParseHurt(param))
}

//*******************PlayerAttack******************

func init() {
	playerAttackFactories["近战攻击"] = createPlayerMeleeAttack
	playerAttackFactories["远程攻击"] = createPlayerRemoteAttack
	playerAttackFactories["群体治疗"] = createPlayerGroupRecover
	playerAttackFactories["远近攻击"] = createPlayerDistanceAttack
	playerAttackFactories["近战双攻"] = createPlayerDoubleAttack
	playerAttackFactories["单体治疗"] = createPlayerSingleRecover
	playerAttackFactories["区域治疗"] = createPlayerRangeRecover
	playerAttackFactories["区域攻击"] = createPlayerRangeAttack
}

//========================PlayerRangeAttack===========================

func createPlayerRangeAttack(player *Player) IAttack {
	return &PlayerRemoteAttack{player: player, range0: true}
}

//=======================PlayerRangeRecover============================

func createPlayerRangeRecover(player *Player) IAttack {
	return &PlayerRangeRecover{player: player}
}

type PlayerRangeRecover struct {
	player *Player
}

func (p *PlayerRangeRecover) PosDraw(pos complex128, screen *ebiten.Image) {
	DrawRange(screen, pos, 0, p.player.Data.Range, ColorRecover)
}

func (p *PlayerRangeRecover) Attack() {
	players := p.player.GetPlayers()
	param := &Param{EventType: EventTypePlayerRecover, Doctor: p.player}
	for i := 0; i < len(players); i++ {
		param.Atk = p.player.BuffHolder.CloneData(FieldAtk)
		param.Player = players[i]
		PlayerManager.TriggerEvent(param)
		if !param.Invalid {
			players[i].Recover(param.Atk.Int())
		}
	}
}

//=======================PlayerSingleRecover============================

func createPlayerSingleRecover(player *Player) IAttack {
	return &PlayerSingleRecover{player: player}
}

type PlayerSingleRecover struct {
	player *Player
}

func (p *PlayerSingleRecover) Attack() {
	players := p.player.GetPlayers()
	if len(players) <= 0 {
		return
	}
	player := players[0]
	max := player.Data.Hp - player.Hp
	for i := 0; i < len(players); i++ {
		if players[i].Data.Hp-players[i].Hp > max {
			player = players[i]
			max = player.Data.Hp - player.Hp
		}
	}
	param := &Param{EventType: EventTypePlayerRecover, Doctor: p.player, Atk: p.player.BuffHolder.CloneData(FieldAtk), Player: player}
	PlayerManager.TriggerEvent(param)
	if !param.Invalid {
		player.Recover(param.Atk.Int())
	}
}

//===================PlayerDoubleAttack=====================

func createPlayerDoubleAttack(player *Player) IAttack {
	return &PlayerDoubleAttack{player: player}
}

type PlayerDoubleAttack struct {
	player *Player
}

func (p *PlayerDoubleAttack) Attack() {
	enemy := GetMinLastEnemy(p.player.GetEnemies())
	if enemy == nil {
		return
	} // 里面会复制  不用担心  传值被改变
	skill := NewComboSkill(p.player, 10, 2, "近战双攻-Combo")
	skill.Enemy = enemy
	p.player.AddSkill("近战双攻-Combo", skill)
}

//===================PlayerDistanceAttack====================

func createPlayerDistanceAttack(player *Player) IAttack {
	return &PlayerDistanceAttack{player: player}
}

type PlayerDistanceAttack struct {
	player *Player
}

func (p *PlayerDistanceAttack) Attack() {
	enemy := GetMinLastEnemy(p.player.GetEnemies())
	if enemy == nil {
		return
	}
	param := &Param{EventType: EventTypePlayerAttack, Player: p.player, Enemy: enemy,
		Atk: p.player.BuffHolder.CloneData(FieldAtk), Def: enemy.BuffHolder.CloneData(FieldDef)}
	isMelee := p.player.BlockEnemies.Has(enemy)
	if !isMelee {
		param.Atk.Mul -= 0.2
	}
	PlayerManager.TriggerEvent(param)
	if param.Invalid {
		return
	}
	if isMelee {
		enemy.Hurt(p.player, ParseHurt(param))
	} else {
		NewPlayerBullet(p.player.Grid.Pos+GridSize/2, ParseHurt(param), enemy, p.player)
	}
}

//==================PlayerGroupRecover===================

func createPlayerGroupRecover(player *Player) IAttack {
	return &PlayerGroupRecover{player: player}
}

type PlayerGroupRecover struct {
	player *Player
}

func (p *PlayerGroupRecover) Attack() {
	players := p.player.GetPlayers()
	if len(players) <= 0 {
		return
	}
	SortPlayerByHurt(players)
	max := utils.Min(3, len(players))
	param := &Param{EventType: EventTypePlayerRecover, Doctor: p.player}
	for i := 0; i < max; i++ {
		param.Atk = p.player.BuffHolder.CloneData(FieldAtk)
		param.Player = players[i]
		PlayerManager.TriggerEvent(param)
		if !param.Invalid {
			players[i].Recover(param.Atk.Int())
		}
	}
}

//====================PlayerRemoteAttack=====================

func createPlayerRemoteAttack(player *Player) IAttack {
	return &PlayerRemoteAttack{player: player, range0: false}
}

type PlayerRemoteAttack struct {
	player *Player
	range0 bool
}

func (p *PlayerRemoteAttack) Attack() {
	enemy := GetMinLastEnemy(p.player.GetEnemies())
	if enemy == nil {
		return
	}
	param := &Param{EventType: EventTypePlayerAttack, Player: p.player, Enemy: enemy,
		Atk: p.player.BuffHolder.CloneData(FieldAtk), Def: enemy.BuffHolder.CloneData(FieldDef)}
	PlayerManager.TriggerEvent(param)
	if param.Invalid {
		return
	}
	NewPlayerBullet(p.player.Grid.Pos+GridSize/2, ParseHurt(param), enemy, p.player).Range = p.range0
}

//======================PlayerMeleeAttack====================

func createPlayerMeleeAttack(player *Player) IAttack {
	return &PlayerMeleeAttack{player: player}
}

type PlayerMeleeAttack struct {
	player *Player
}

func (p *PlayerMeleeAttack) Attack() {
	enemy := GetMinLastEnemy(p.player.GetEnemies())
	if enemy == nil {
		return
	}
	param := &Param{EventType: EventTypePlayerAttack, Player: p.player, Enemy: enemy,
		Atk: p.player.BuffHolder.CloneData(FieldAtk), Def: enemy.BuffHolder.CloneData(FieldDef)}
	PlayerManager.TriggerEvent(param)
	if param.Invalid {
		return
	}
	enemy.Hurt(p.player, ParseHurt(param))
}
