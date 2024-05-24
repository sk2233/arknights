/*
@author: sk
@date: 2023/2/26
*/
package main

import (
	"GameBase2/model"
	"GameBase2/utils"
	R "arknights/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var (
	skillFactories = make(map[string]func(*Player) ISkill)
)

func CreateSkill(name string, player *Player) ISkill {
	return skillFactories[name](player)
}

func init() {
	// 共用方法
	skillFactories["撤退"] = createRetreatSkill
	skillFactories["迅捷打击A型"] = createXunJieDaJiTypeASkill
	skillFactories["攻击力强化A型"] = createGongJiLiQiangHuaTypeASkill
	skillFactories["治疗强化B型"] = createZhiLiaoQiangHuaTypeBSkill
	skillFactories["强力击A型"] = createQiangLiJiTypeASkill
	skillFactories["战术咏唱A型"] = createZhanShuYongChangTypeASkill
	skillFactories["治疗强化A型"] = createZhiLiaoQiangHuaTypeASkill
	skillFactories["战术咏唱B型"] = createZhanShuYongChangTypeBSkill
	// 特殊
	skillFactories["自毁"] = createZiHuiSkill
	// 星熊
	skillFactories["战术装甲"] = createZhanShuZhuangJiaSkill
	skillFactories["特种作战策略"] = createTeZhongZuoZhanCeLveSkill
	skillFactories["战意"] = createZhanYiSkill
	skillFactories["荆棘"] = createJingJiSkill
	skillFactories["力之锯"] = createLiZhiJuSkill
	// 斯卡蒂
	skillFactories["深海掠食者"] = createShenHaiLveShiZheSkill
	skillFactories["跃浪击"] = createYueLangJiSkill
	skillFactories["涌潮悲歌"] = createYongChaoBeiGeSkill
	// 白金
	skillFactories["蓄力攻击"] = createXuLiGongJiSkill
	skillFactories["天马视域"] = createTianMaShiYuSkill
	// 调香师
	skillFactories["熏衣草"] = createXunYiCaoSkill
	skillFactories["精调"] = createJingTiaoSkill
	// 雷蛇
	skillFactories["战术防御"] = createZhanShuFangYuSkill
	skillFactories["雷抗"] = createLeiKangSkill
	skillFactories["充能防御"] = createChongNengFangYuSkill
	skillFactories["反击电弧"] = createFanJiDianHuSkill
	// 银灰
	skillFactories["领袖"] = createLingXiuSkill
	skillFactories["雪境生存法则"] = createXueJingShengCunFaZeSkill
	skillFactories["真银斩"] = createZhenYinZhanSkill
	// 阿米娅
	skillFactories["情绪吸收"] = createQingXuXiShouSkill
	skillFactories["精神爆发"] = createJingShenBaoFaSkill
	skillFactories["奇美拉"] = createQiMeiLaSkill
	// 陈
	skillFactories["呵斥"] = createHeChiSkill
	skillFactories["持刀格斗术"] = createChiDaoGeDouShuSkill
	skillFactories["鞘击"] = createShaoJiSkill
	skillFactories["赤霄·拔刀"] = createChiXiaoBaDaoSkill
	skillFactories["赤霄·绝影"] = createChiXiaoJueYingSkill
	// 赫默
	skillFactories["强化注射"] = createQiangHuaZhuSheSkill
	skillFactories["医疗无人机"] = createYiLiaoWuRenJiSkill
	// 远山
	skillFactories["占卜"] = createZhanBuSkill
	skillFactories["命运"] = createMingYunSkill
	// 塞雷娅
	skillFactories["莱茵充能护服"] = createLaiYinChongNengHuFuSkill
	skillFactories["精神回复"] = createJingShenHuiFuSkill
	skillFactories["急救"] = createJiJiuSkill
	skillFactories["药物配置"] = createYaoWuPeiZhiSkill
	skillFactories["钙质化"] = createGaiZhiHuaSkill
	// 能天使
	skillFactories["快速弹匣"] = createKuaiSuTanXiaSkill
	skillFactories["天使的祝福"] = createTianShiDeZhuFuSkill
	skillFactories["冲锋模式"] = createChongFengMoShiSkill
	skillFactories["扫射模式"] = createSaoSheMoShiSkill
	skillFactories["过载模式"] = createGuoZaiMoShiSkill
}

//=======================GuoZaiMoShiSkill===========================

func createGuoZaiMoShiSkill(player *Player) ISkill {
	return &GuoZaiMoShiSkill{player: player}
}

type GuoZaiMoShiSkill struct {
	player      *Player
	timer       int64
	effectTimer int
}

func (j *GuoZaiMoShiSkill) Recover(recover int64) {
	j.timer += recover
}

func (j *GuoZaiMoShiSkill) Update() {
	j.timer++
	if j.timer > 50*60 {
		j.timer = 0
		j.effectTimer = 15 * 60
	}
	j.effectTimer--
}

func (j *GuoZaiMoShiSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == j.player && j.effectTimer > 0
}

func (j *GuoZaiMoShiSkill) Action(param *Param) {
	param.Invalid = true // 取消普攻 进行多人攻击
	skill := NewComboSkill(j.player, 10, 5, "远程五攻-Combo")
	skill.Melee = false
	j.player.AddSkill("远程五攻-Combo", skill)
}

//==================SaoSheMoShiSkill==================

func createSaoSheMoShiSkill(player *Player) ISkill {
	return NewInitiativeSkill("扫射模式", ToFrame(45), player, SaoSheMoShiHandler)
}

func SaoSheMoShiHandler(param *Param) {
	param.Player.AddSkill("扫射模式-连发", NewSaoSheMoShiSkill(param.Player))
}

type SaoSheMoShiSkill struct {
	*TimeSkill
}

func NewSaoSheMoShiSkill(player *Player) *SaoSheMoShiSkill {
	res := &SaoSheMoShiSkill{}
	res.TimeSkill = NewTimeSkill(player, ToFrame(15), "扫射模式-连发")
	return res
}

func (l *SaoSheMoShiSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == l.Player
}

func (l *SaoSheMoShiSkill) Action(param *Param) {
	param.Invalid = true
	skill := NewComboSkill(l.Player, 10, 4, "远程四攻-Combo")
	skill.Melee = false
	l.Player.AddSkill("远程四攻-Combo", skill)
}

//====================ChongFengMoShiSkill=====================

func createChongFengMoShiSkill(player *Player) ISkill {
	return &ChongFengMoShiSkill{player: player, timer: 0}
}

type ChongFengMoShiSkill struct {
	player *Player
	timer  int64
}

func (c *ChongFengMoShiSkill) Recover(recover int64) {
	c.timer += recover
}

func (c *ChongFengMoShiSkill) Update() {
	c.timer++
}

func (c *ChongFengMoShiSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == c.player && c.timer > 5*60
}

func (c *ChongFengMoShiSkill) Action(param *Param) {
	c.timer = 0
	param.Invalid = true
	skill := NewComboSkill(c.player, 10, 3, "远程三攻-Combo")
	skill.Atk.Mul = 1.05
	skill.Melee = false
	c.player.AddSkill("远程三攻-Combo", skill)
}

//======================TianShiDeZhuFuSkill=========================

func createTianShiDeZhuFuSkill(player *Player) ISkill {
	return NewLaunchSkill(player, TianShiDeZhuFuHandler)
}

func TianShiDeZhuFuHandler(param *Param) {
	buff1 := &Buff{
		Name:  "天使的祝福-攻击力",
		Timer: -1,
		Field: FieldAtk,
		Mul:   0.06,
	}
	buff2 := &Buff{
		Name:  "天使的祝福-防御力",
		Timer: -1,
		Field: FieldDef,
		Mul:   0.1,
	}
	param.Player.BuffHolder.AddBuff(buff1)
	param.Player.BuffHolder.AddBuff(buff2)
	player := utils.RandomItem(PlayerManager.GetPlayers(IsAny)...)
	player.BuffHolder.AddBuff(buff1)
	player.BuffHolder.AddBuff(buff2)
}

//======================LeiKangSkill============================

func createLeiKangSkill(player *Player) ISkill {
	return NewLaunchSkill(player, LeiKangHandler)
}

func LeiKangHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "雷抗-防御力",
		Timer: -1,
		Field: FieldDef,
		Add:   10,
	})
}

//=====================KuaiSuTanXiaSkill=====================

func createKuaiSuTanXiaSkill(player *Player) ISkill {
	return NewLaunchSkill(player, KuaiSuTanXiaHandler)
}

func KuaiSuTanXiaHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "快速弹匣-攻击速度",
		Timer: -1,
		Field: FieldAttackSpeed,
		Add:   6,
	})
}

//========================GaiZhiHuaSkill===========================

func createGaiZhiHuaSkill(player *Player) ISkill {
	return NewInitiativeSkill("钙质化", ToFrame(80), player, GaiZhiHuaHandler)
}

func GaiZhiHuaHandler(param *Param) {
	param.Player.AddSkill("钙质化-区域效果", NewGaiZhiHuaSkill(param.Player))
}

type GaiZhiHuaSkill struct {
	*TimeSkill
	offsets [][]int
	enemies *model.Set[*Enemy]
}

func (g *GaiZhiHuaSkill) PosDraw(pos complex128, screen *ebiten.Image) {
	DrawRange(screen, pos, 0, g.offsets, ColorGaiZhiHua)
}

func NewGaiZhiHuaSkill(player *Player) *GaiZhiHuaSkill {
	res := &GaiZhiHuaSkill{offsets: [][]int{{0, -3},
		{-1, -2}, {0, -2}, {1, -2},
		{-2, -1}, {-1, -1}, {0, -1}, {1, -1}, {2, -1},
		{-3, 0}, {-2, 0}, {-1, 0}, {0, 0}, {1, 0}, {2, 0}, {3, 0},
		{-2, 1}, {-1, 1}, {0, 1}, {1, 1}, {2, 1},
		{-1, 2}, {0, 2}, {1, 2},
		{0, 3}}, enemies: model.NewSet[*Enemy]()}
	res.TimeSkill = NewTimeSkill(player, ToFrame(10), "钙质化-区域效果")
	return res
}

func (g *GaiZhiHuaSkill) Update() {
	g.TimeSkill.Update()
	if g.Timer%60 == 0 { // 每秒的效果
		player := g.Player
		players := CollisionPlayers(player.Grid.Pos, 0, g.offsets) // 回血
		param := &Param{EventType: EventTypePlayerRecover, Doctor: player}
		for i := 0; i < len(players); i++ {
			param.Atk = player.BuffHolder.CloneData(FieldAtk)
			param.Atk.Mul = 0.1
			param.Player = players[i]
			PlayerManager.TriggerEvent(param)
			if !param.Invalid {
				players[i].Recover(param.Atk.Int())
			}
		}
		enemies := CollisionEnemies(player.Grid.Pos, 0, g.offsets)
		g.enemies.Clear()
		for i := 0; i < len(enemies); i++ {
			enemies[i].BuffHolder.AddBuff(&Buff{
				Name:  "钙质化-移动速度",
				Timer: ToFrame(1),
				Field: FieldMoveSpeed,
				Mul:   -0.6,
			})
			g.enemies.Add(enemies[i])
		}
	}
}

func (g *GaiZhiHuaSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && g.enemies.Has(param.Enemy)
}

func (g *GaiZhiHuaSkill) Action(param *Param) {
	param.Atk.Mul += 0.2
}

//========================YaoWuPeiZhiSkill=======================

func createYaoWuPeiZhiSkill(player *Player) ISkill {
	return &YaoWuPeiZhiSkill{player: player, timer: 0, offsets: [][]int{{-1, -2}, {0, -2}, {1, -2},
		{-2, -1}, {-1, -1}, {0, -1}, {1, -1}, {2, -1}, {-2, 0}, {-1, 0}, {0, 0}, {1, 0}, {2, 0}, {-2, 1}, {-1, 1}, {0, 1}, {1, 1}, {2, 1},
		{-1, 2}, {0, 2}, {1, 2}}}
}

type YaoWuPeiZhiSkill struct {
	player  *Player
	timer   int64
	offsets [][]int
}

func (y *YaoWuPeiZhiSkill) Recover(recover int64) {
	y.timer += recover
}

func (y *YaoWuPeiZhiSkill) Update() {
	if y.timer < 10*60 {
		y.timer++
	} else {
		y.timer = 0
		player := y.player
		NewRangeTip(player.Grid.Pos, 0, y.offsets, ColorHint)
		players := CollisionPlayers(player.Grid.Pos, 0, y.offsets)
		param := &Param{EventType: EventTypePlayerRecover, Doctor: player}
		for i := 0; i < len(players); i++ {
			param.Atk = player.BuffHolder.CloneData(FieldAtk)
			param.Atk.Mul = 0.8
			param.Player = players[i]
			PlayerManager.TriggerEvent(param)
			if !param.Invalid {
				players[i].Recover(param.Atk.Int())
			}
		}
	}
}

func (y *YaoWuPeiZhiSkill) CanHandle(param *Param) bool {
	return false
}

func (y *YaoWuPeiZhiSkill) Action(param *Param) {
}

//=====================JiJiuSkill======================

func createJiJiuSkill(player *Player) ISkill {
	return &JiJiuSkill{player: player, offsets: [][]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {0, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1}}}
}

type JiJiuSkill struct {
	player  *Player
	timer   int64
	offsets [][]int
}

func (j *JiJiuSkill) Recover(recover int64) {
	j.timer += recover
}

func (j *JiJiuSkill) Update() {
	j.timer++
}

func (j *JiJiuSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == j.player && j.timer > 6*60
}

func (j *JiJiuSkill) Action(param *Param) {
	players := CollisionPlayers(j.player.Grid.Pos, 0, j.offsets)
	for i := 0; i < len(players); i++ {
		if players[i].Hp < players[i].Data.Hp/2 {
			j.timer = 0
			param = &Param{EventType: EventTypePlayerRecover, Doctor: j.player,
				Atk: j.player.BuffHolder.CloneData(FieldAtk), Player: players[i]}
			param.Atk.Mul = 1.1
			PlayerManager.TriggerEvent(param)
			if !param.Invalid {
				players[i].Recover(param.Atk.Int())
			}
			return
		}
	}
}

//====================JingShenHuiFuSkill==================

func createJingShenHuiFuSkill(player *Player) ISkill {
	return &JingShenHuiFuSkill{player: player}
}

type JingShenHuiFuSkill struct {
	player *Player
}

func (j *JingShenHuiFuSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerRecover && param.Doctor == j.player
}

func (j *JingShenHuiFuSkill) Action(param *Param) {
	param.Player.RecoverSkill(60)
}

//=====================LaiYinChongNengHuFuSkill======================

func createLaiYinChongNengHuFuSkill(player *Player) ISkill {
	return &LaiYinChongNengHuFuSkill{player: player, timer: ToFrame(20), count: 0}
}

type LaiYinChongNengHuFuSkill struct {
	player *Player
	timer  int64
	count  int
}

func (l *LaiYinChongNengHuFuSkill) Update() {
	if l.timer > 0 {
		l.timer--
	} else {
		l.timer = ToFrame(20)
		l.count++
		l.player.BuffHolder.AddBuff(&Buff{
			Name:  "莱茵充能护服-攻击力",
			Timer: -1,
			Field: FieldAtk,
			Mul:   0.02 * float64(l.count),
		})
		l.player.BuffHolder.AddBuff(&Buff{
			Name:  "莱茵充能护服-防御力",
			Timer: -1,
			Field: FieldDef,
			Mul:   0.02 * float64(l.count),
		})
		if l.count >= 5 {
			l.player.RemoveSkill("莱茵充能护服")
		}
	}
}

func (l *LaiYinChongNengHuFuSkill) CanHandle(param *Param) bool {
	return false
}

func (l *LaiYinChongNengHuFuSkill) Action(param *Param) {
}

//======================MingYunSkill=========================

func createMingYunSkill(player *Player) ISkill {
	return NewInitiativeSkill("命运", ToFrame(100), player, MingYunHandler)
}

func MingYunHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "命运-攻击力",
		Timer: ToFrame(30),
		Field: FieldAtk,
		Mul:   0.3,
	})
	param.Player.Range.SetRange([][]int{{0, -1}, {1, -1}, {2, -1}, {3, -1}, {0, 0}, {1, 0}, {2, 0}, {3, 0},
		{0, 1}, {1, 1}, {2, 1}, {3, 1}})
	param.Player.AddSkill("命运-全打", NewMingYunSkill(param.Player))
}

type MingYunSkill struct {
	*TimeSkill
}

func NewMingYunSkill(player *Player) *MingYunSkill {
	res := &MingYunSkill{}
	res.TimeSkill = NewTimeSkill(player, ToFrame(30), "命运-全打")
	res.EndFunc = res.endHandler
	return res
}

func (l *MingYunSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == l.Player
}

func (l *MingYunSkill) Action(param *Param) {
	param.Invalid = true // 取消普攻 进行多人攻击
	player := l.Player
	enemies := player.GetEnemies()
	param = &Param{EventType: EventTypePlayerSkill, Player: player}
	for i := 0; i < len(enemies); i++ {
		param.Enemy = enemies[i]
		param.Atk = player.BuffHolder.CloneData(FieldAtk)
		param.Def = enemies[i].BuffHolder.CloneData(FieldDef)
		PlayerManager.TriggerEvent(param)
		if !param.Invalid {
			NewPlayerBullet(player.Grid.Pos+GridSize/2, ParseHurt(param), enemies[i], player).Range = true
		}
	}
}

func (l *MingYunSkill) endHandler() { // 恢复
	l.Player.Range.SetRange(l.Player.Data.Range)
}

//=======================ZhanShuYongChangTypeBSkill======================

func createZhanShuYongChangTypeBSkill(player *Player) ISkill {
	return NewInitiativeSkill("战术咏唱B型", ToFrame(45), player, ZhanShuYongChangTypeBHandler)
}

func ZhanShuYongChangTypeBHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "战术咏唱B型-攻击时间",
		Timer: ToFrame(25),
		Field: FieldAttackSpeed,
		Add:   15,
	})
}

//==================ZhanBuSkill==================

func createZhanBuSkill(player *Player) ISkill {
	return NewLaunchSkill(player, ZhanBuHandler)
}

func ZhanBuHandler(param *Param) {
	field := utils.RandomItem(FieldAtk, FieldDef, FieldMaxHp)
	if field == FieldMaxHp {
		param.Player.BuffHolder.AddBuff(&Buff{
			Name:  "占卜-随机增益",
			Timer: -1,
			Field: FieldMaxHp,
			Mul:   0.12,
		})
	} else {
		param.Player.BuffHolder.AddBuff(&Buff{
			Name:  "占卜-随机增益",
			Timer: -1,
			Field: field,
			Mul:   0.07,
		})
	}
}

//===================ZiHuiSkill====================

func createZiHuiSkill(player *Player) ISkill {
	return &ZiHuiSkill{player: player, timer: ToFrame(10)}
}

type ZiHuiSkill struct {
	player *Player
	timer  int64
}

func (z *ZiHuiSkill) Update() {
	if z.timer > 0 {
		z.timer--
	} else {
		z.player.Retreat()
	}
}

func (z *ZiHuiSkill) CanHandle(param *Param) bool {
	return false
}

func (z *ZiHuiSkill) Action(param *Param) {
}

//====================YiLiaoWuRenJiSkill=====================

func createYiLiaoWuRenJiSkill(player *Player) ISkill {
	return &YiLiaoWuRenJiSkill{timer: 0, has: false, data: &CardData{
		Name:        "医疗无人机",
		Img:         R.MAIN.PLAYER.PLAYER9,
		Place:       PlaceHighland | PlaceLand,
		ChangePlace: PlaceNone,
		Atk:         30,
		Def:         0,
		BlockNum:    0,
		CostNum:     5,
		Hp:          1000,
		AttackTime:  ToFrame(1),
		CoolTime:    ToFrame(70),
		Skills:      []string{"撤退", "自毁"},
		Range:       [][]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {0, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}},
		Attack:      "区域治疗",
		Career:      CareerSummon,
	}}
}

type YiLiaoWuRenJiSkill struct {
	timer int64
	has   bool
	data  *CardData
}

func (y *YiLiaoWuRenJiSkill) Recover(recover int64) {
	y.timer += recover
}

func (y *YiLiaoWuRenJiSkill) Update() {
	if y.timer < 30*60 {
		y.timer++
	} else {
		if !y.has {
			y.timer = 0
			CardShow.AddCard(NewCard(y.data, 0))
			y.has = true
		}
	}
}

func (y *YiLiaoWuRenJiSkill) CanHandle(param *Param) bool { // 监听召唤物退场
	return param.EventType == EventTypePlayerRetreat && param.Player.Data == y.data
}

func (y *YiLiaoWuRenJiSkill) Action(param *Param) {
	y.has = false
}

//========================ZhiLiaoQiangHuaTypeASkill=====================

func createZhiLiaoQiangHuaTypeASkill(player *Player) ISkill {
	return NewInitiativeSkill("治疗强化A型", ToFrame(40), player, ZhiLiaoQiangHuaTypeAHandler)
}

func ZhiLiaoQiangHuaTypeAHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "治疗强化A型-攻击力",
		Timer: ToFrame(30),
		Field: FieldAtk,
		Mul:   0.4,
	})
}

//===================QiangHuaZhuSheSkill====================

func createQiangHuaZhuSheSkill(player *Player) ISkill {
	return NewLifeBuffSkill(player, IsCareer(CareerMedical), &Buff{
		Name:  "强化注射-攻击速度",
		Timer: -1,
		Field: FieldAttackSpeed,
		Add:   6,
	})
}

//===================ChiXiaoJueYingSkill==================

func createChiXiaoJueYingSkill(player *Player) ISkill {
	return NewInitiativeSkill("赤霄*绝影", ToFrame(40), player, ChiXiaoJueYingHandler)
}

var (
	chiXiaoJueYingOffsets = [][]int{{0, -2}, {-1, -1}, {0, -1}, {1, -1}, {-2, 0}, {-1, 0}, {0, 0}, {1, 0},
		{2, 0}, {0, 2}, {-1, 1}, {0, 1}, {1, 1}}
)

func ChiXiaoJueYingHandler(param *Param) {
	player := param.Player
	NewRangeTip(player.Grid.Pos, player.RangeDir, chiXiaoJueYingOffsets, ColorCardMask)
	enemies := CollisionEnemies(player.Grid.Pos, player.RangeDir, chiXiaoJueYingOffsets)
	if len(enemies) <= 0 {
		return
	}
	enemy := enemies[0]
	minL2 := utils.VectorLen2(enemy.Pos - (player.Grid.Pos + GridSize/2))
	for i := 1; i < len(enemies); i++ {
		if l2 := utils.VectorLen2(enemies[i].Pos - (player.Grid.Pos + GridSize/2)); l2 < minL2 {
			enemy = enemies[i]
			minL2 = l2
		}
	}
	skill := NewComboSkill(player, 10, 10, "赤霄·绝影-Combo")
	skill.Enemy = enemy
	skill.Atk.Mul = 2
	skill.EndFunc = func() {
		enemy.BuffHolder.AddBuff(&Buff{
			Name:  "赤霄·绝影-眩晕",
			Timer: ToFrame(2),
			Field: FieldState,
			State: StateMove | StateAttack,
		})
	}
	player.AddSkill("赤霄·绝影-Combo", skill)
}

//================ChiXiaoBaDaoSkill================

func createChiXiaoBaDaoSkill(player *Player) ISkill {
	return NewInitiativeSkill("赤霄*拔刀", ToFrame(27), player, ChiXiaoBaDaoHandler)
}

var (
	chiXiaoBaDaoOffsets = [][]int{{0, -1}, {1, -1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {1, 1}}
)

func ChiXiaoBaDaoHandler(param *Param) {
	player := param.Player
	NewRangeTip(player.Grid.Pos, player.RangeDir, chiXiaoBaDaoOffsets, ColorCardMask)
	enemies := CollisionEnemies(player.Grid.Pos, player.RangeDir, chiXiaoBaDaoOffsets)
	if len(enemies) <= 0 {
		return
	}
	SortEnemyByLast(enemies)
	max := utils.Min(4, len(enemies))
	param = &Param{EventType: EventTypePlayerSkill, Player: player}
	for i := 0; i < max; i++ {
		param.Atk = player.BuffHolder.CloneData(FieldAtk)
		param.Atk.Mul = 6.6
		param.Enemy = enemies[i]
		param.Def = enemies[i].BuffHolder.CloneData(FieldDef)
		PlayerManager.TriggerEvent(param)
		if !param.Invalid {
			enemies[i].Hurt(player, ParseHurt(param))
		}
	}
}

//=================ShaoJiSkill===================

func createShaoJiSkill(player *Player) ISkill {
	return &ShaoJiSkill{player: player}
}

type ShaoJiSkill struct {
	player *Player
	timer  int64
}

func (s *ShaoJiSkill) Recover(recover int64) {
	s.timer += recover
}

func (s *ShaoJiSkill) Update() {
	s.timer++
}

func (s *ShaoJiSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == s.player && s.timer > 7*60
}

func (s *ShaoJiSkill) Action(param *Param) {
	s.timer = 0
	param.Atk.Mul = 2
	param.Enemy.BuffHolder.AddBuff(&Buff{
		Name:  "鞘击-眩晕",
		Timer: ToFrame(1),
		Field: FieldState,
		State: StateMove | StateAttack,
	})
}

//================ChiDaoGeDouShuSkill=================

func createChiDaoGeDouShuSkill(player *Player) ISkill {
	return &ChiDaoGeDouShuSkill{player: player}
}

type ChiDaoGeDouShuSkill struct {
	player *Player
}

func (c *ChiDaoGeDouShuSkill) CanHandle(param *Param) bool {
	if param.EventType == EventTypePlayerInit && param.Player == c.player {
		c.player.BuffHolder.AddBuff(&Buff{
			Name:  "持刀格斗术-攻击力",
			Timer: -1,
			Field: FieldAtk,
			Mul:   0.05,
		})
		c.player.BuffHolder.AddBuff(&Buff{
			Name:  "持刀格斗术-防御力",
			Timer: -1,
			Field: FieldDef,
			Mul:   0.05,
		})
	}
	return param.EventType == EventTypeEnemyAttack && param.Player == c.player
}

func (c *ChiDaoGeDouShuSkill) Action(param *Param) {
	if utils.Change(1, 10) {
		param.Invalid = true
		c.player.AddInfo("Miss", colornames.White)
	}
}

//===================HeChiSkill=================

func createHeChiSkill(player *Player) ISkill {
	return &HeChiSkill{timer: 5 * 60}
}

type HeChiSkill struct {
	timer int
}

func (h *HeChiSkill) Update() {
	if h.timer > 0 {
		h.timer--
	} else {
		h.timer = 5 * 60
		players := PlayerManager.GetPlayers(IsAny)
		for i := 0; i < len(players); i++ {
			players[i].RecoverSkill(60)
		}
	}
}

func (h *HeChiSkill) CanHandle(param *Param) bool {
	return false
}

func (h *HeChiSkill) Action(param *Param) {

}

//=====================QiMeiLaSkill===================

func createQiMeiLaSkill(player *Player) ISkill {
	return NewInitiativeSkill("奇美拉", ToFrame(120), player, QiMeiLaHandler)
}

func QiMeiLaHandler(param *Param) {
	player := param.Player
	player.BuffHolder.AddBuff(&Buff{
		Name:  "奇美拉-攻击力",
		Timer: ToFrame(30),
		Field: FieldAtk,
		Mul:   1,
	})
	player.BuffHolder.AddBuff(&Buff{
		Name:  "奇美拉-最大生命值",
		Timer: ToFrame(30),
		Field: FieldMaxHp,
		Mul:   0.25,
	})
	player.Recover(player.Data.Hp * 25 / 100)
	param.Player.AddSkill("奇美拉-真伤", NewQiMeiLaSkill(param.Player))
	param.Player.Range.AddRange([][]int{{0, -2}, {1, -2}, {2, -2}, {3, -1}, {3, 0}, {3, 1}, {0, 2}, {1, 2},
		{2, 2}})
}

type QiMeiLaSkill struct {
	*TimeSkill
}

func NewQiMeiLaSkill(player *Player) *QiMeiLaSkill {
	res := &QiMeiLaSkill{}
	res.TimeSkill = NewTimeSkill(player, ToFrame(30), "奇美拉-真伤")
	res.EndFunc = res.endHandler
	return res
}

func (q *QiMeiLaSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == q.Player
}

func (q *QiMeiLaSkill) Action(param *Param) {
	param.RealHurt = true
}

func (q *QiMeiLaSkill) endHandler() {
	q.Player.Range.SetRange(q.Player.Data.Range)
	q.Player.Retreat()
}

//=================JingShenBaoFaSkill===================

func createJingShenBaoFaSkill(player *Player) ISkill {
	return &JingShenBaoFaSkill{player: player}
}

type JingShenBaoFaSkill struct {
	player      *Player
	timer       int64
	effectTimer int
}

func (j *JingShenBaoFaSkill) Recover(recover int64) {
	j.timer += recover
}

func (j *JingShenBaoFaSkill) Update() {
	j.timer++
	if j.timer > 6000 {
		j.timer = 0
		j.effectTimer = 25 * 60
	}
	j.effectTimer--
}

func (j *JingShenBaoFaSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == j.player && j.effectTimer > 0
}

func (j *JingShenBaoFaSkill) Action(param *Param) {
	param.Invalid = true // 取消普攻 进行多人攻击
	skill := NewComboSkill(j.player, 10, 6, "精神爆发-连发")
	skill.Atk.Mul = 0.33
	j.player.AddSkill("精神爆发-连发", skill)
}

//==================ZhanShuYongChangTypeASkill============

func createZhanShuYongChangTypeASkill(player *Player) ISkill {
	return NewInitiativeSkill("战术咏唱A型", ToFrame(40), player, ZhanShuYongChangTypeAHandler)
}

func ZhanShuYongChangTypeAHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "战术咏唱A型-攻击时间",
		Timer: ToFrame(30),
		Field: FieldAttackSpeed,
		Add:   30,
	})
}

//===================QingXuXiShouSkill===================

func createQingXuXiShouSkill(player *Player) ISkill {
	return &QingXuXiShouSkill{player: player}
}

type QingXuXiShouSkill struct {
	player *Player
}

func (q *QingXuXiShouSkill) CanHandle(param *Param) bool {
	return (param.EventType == EventTypePlayerAttack || param.EventType == EventTypeEnemyDeath) && param.Player == q.player
}

func (q *QingXuXiShouSkill) Action(param *Param) {
	if param.EventType == EventTypePlayerAttack {
		q.player.RecoverSkill(60 * 2)
	} else {
		q.player.RecoverSkill(60 * 8)
	}
}

//===================ZhenYinZhanSkill=================

func createZhenYinZhanSkill(player *Player) ISkill {
	return NewInitiativeSkill("真银斩", ToFrame(90), player, ZhenYinZhanHandler)
}

func ZhenYinZhanHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "真银斩-防御力",
		Timer: ToFrame(20),
		Field: FieldDef,
		Mul:   -0.7,
	})
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "真银斩-攻击力",
		Timer: ToFrame(20),
		Field: FieldAtk,
		Mul:   1.1,
	})
	param.Player.AddSkill("真银斩-多打", NewZhenYinZhanSkill(param.Player))
	param.Player.Range.SetRange([][]int{{0, -3}, {0, -2}, {1, -2}, {0, -1}, {1, -1}, {2, -1}, {0, 0},
		{1, 0}, {2, 0}, {3, 0}, {0, 3}, {0, 2}, {1, 2}, {0, 1}, {1, 1}, {2, 1}})
}

type ZhenYinZhanSkill struct {
	*TimeSkill
	offsets [][]int
}

func NewZhenYinZhanSkill(player *Player) *ZhenYinZhanSkill {
	res := &ZhenYinZhanSkill{offsets: [][]int{{0, -3}, {0, -2}, {1, -2}, {0, -1}, {1, -1}, {2, -1}, {0, 0},
		{1, 0}, {2, 0}, {3, 0}, {0, 3}, {0, 2}, {1, 2}, {0, 1}, {1, 1}, {2, 1}}}
	res.TimeSkill = NewTimeSkill(player, ToFrame(20), "真银斩-多打")
	res.EndFunc = res.endHandler
	return res
}

func (l *ZhenYinZhanSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == l.Player
}

func (l *ZhenYinZhanSkill) Action(param *Param) {
	param.Invalid = true // 取消普攻 进行多人攻击
	player := l.Player
	enemies := player.GetEnemies()
	SortEnemyByLast(enemies)
	NewRangeTip(player.Grid.Pos, player.RangeDir, l.offsets, ColorCardMask)
	max := utils.Min(3, len(enemies))
	param = &Param{EventType: EventTypePlayerSkill, Player: player}
	for i := 0; i < max; i++ {
		param.Enemy = enemies[i]
		param.Atk = player.BuffHolder.CloneData(FieldAtk)
		param.Def = enemies[i].BuffHolder.CloneData(FieldDef)
		PlayerManager.TriggerEvent(param)
		if !param.Invalid {
			enemies[i].Hurt(l.Player, ParseHurt(param))
		}
	}
}

func (l *ZhenYinZhanSkill) endHandler() {
	l.Player.Range.SetRange(l.Player.Data.Range)
}

//=================XueJingShengCunFaZeSkill===============

func createXueJingShengCunFaZeSkill(player *Player) ISkill {
	res := &XueJingShengCunFaZeSkill{open: false, offsets: [][]int{{0, -1}, {0, 0}, {1, 0}, {0, 1}}}
	res.InitiativeSkill = NewInitiativeSkill("雪境生存法则", ToFrame(5), player, res.Handler)
	return res
}

type XueJingShengCunFaZeSkill struct {
	*InitiativeSkill
	offsets [][]int
	open    bool
	timer   int
}

func (x *XueJingShengCunFaZeSkill) Update() {
	x.InitiativeSkill.Update()
	if !x.open {
		return
	}
	if x.timer > 0 {
		x.timer--
	} else {
		x.timer = 60
		x.player.Recover(int64(float64(x.player.Data.Hp) * 0.03))
	}
}

func (x *XueJingShengCunFaZeSkill) Handler(param *Param) {
	x.open = !x.open
	if x.open {
		x.player.Range.SetRange(x.offsets)
		x.player.BuffHolder.AddBuff(&Buff{
			Name:  "雪境生存法则-防御力",
			Timer: -1,
			Field: FieldDef,
			Mul:   0.35,
		})
		x.timer = 60
	} else {
		x.player.Range.SetRange(x.player.Data.Range)
		x.player.BuffHolder.RemoveBuff("雪境生存法则-防御力")
	}
}

//==============QiangLiJiTypeASkill============

func createQiangLiJiTypeASkill(player *Player) ISkill {
	return &QiangLiJiTypeASkill{player: player}
}

type QiangLiJiTypeASkill struct {
	player *Player
	timer  int64
}

func (q *QiangLiJiTypeASkill) Recover(recover int64) {
	q.timer += recover
}

func (q *QiangLiJiTypeASkill) Update() {
	q.timer++
}

func (q *QiangLiJiTypeASkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == q.player && q.timer > 4*60
}

func (q *QiangLiJiTypeASkill) Action(param *Param) {
	q.timer = 0
	param.Atk.Mul = 1.9 // 提高到
}

//================LingXiuSkill===================

func createLingXiuSkill(player *Player) ISkill {
	return NewLaunchSkill(player, LingXiuHandler)
}

func LingXiuHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "领袖-攻击力",
		Timer: -1,
		Field: FieldAtk,
		Mul:   0.05,
	})
	CardShow.ReduceTime(0.05, 0)
}

//===================FanJiDianHuSkill====================

func createFanJiDianHuSkill(player *Player) ISkill {
	return NewInitiativeSkill("反击电弧", ToFrame(44), player, FanJiDianHuHandler)
}

func FanJiDianHuHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "反击电弧-攻击速度",
		Timer: ToFrame(20),
		Field: FieldAttackSpeed,
		Mul:   -40,
	})
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "反击电弧-攻击力",
		Timer: ToFrame(20),
		Field: FieldAtk,
		Mul:   0.85,
	})
	param.Player.AddSkill("反击电弧-眩晕", NewFanJiDianHuSkill(param.Player))
}

type FanJiDianHuSkill struct {
	*TimeSkill
}

func NewFanJiDianHuSkill(player *Player) *FanJiDianHuSkill {
	res := &FanJiDianHuSkill{}
	res.TimeSkill = NewTimeSkill(player, ToFrame(20), "反击电弧-眩晕")
	return res
}

func (l *FanJiDianHuSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == l.Player
}

func (l *FanJiDianHuSkill) Action(param *Param) {
	param.Invalid = true // 取消普攻 进行多人攻击
	player := l.Player
	enemies := player.GetEnemies()
	SortEnemyByLast(enemies)
	max := utils.Min(3, len(enemies))
	param = &Param{EventType: EventTypePlayerSkill, Player: player}
	for i := 0; i < max; i++ {
		param.Enemy = enemies[i]
		param.Atk = player.BuffHolder.CloneData(FieldAtk)
		param.Def = enemies[i].BuffHolder.CloneData(FieldDef)
		PlayerManager.TriggerEvent(param)
		if !param.Invalid {
			enemies[i].Hurt(l.Player, ParseHurt(param))
			if utils.Change(1, 10) { // 1/10 眩晕
				enemies[i].BuffHolder.AddBuff(&Buff{
					Name:  "反击电弧-眩晕",
					Timer: ToFrame(1),
					Field: FieldState,
					State: StateMove | StateAttack,
				})
			}
		}
	}
}

//================ChongNengFangYuSkill=================

func createChongNengFangYuSkill(player *Player) ISkill {
	return &ChongNengFangYuSkill{player: player}
}

type ChongNengFangYuSkill struct {
	player *Player
	timer  int64
}

func (c *ChongNengFangYuSkill) Recover(recover int64) {
	c.timer += recover
}

func (c *ChongNengFangYuSkill) Update() {
	c.timer++
}

func (c *ChongNengFangYuSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypeEnemyAttack && param.Player == c.player && c.timer > ToFrame(24)
}

func (c *ChongNengFangYuSkill) Action(param *Param) {
	c.timer = 0
	param.Invalid = true
	c.player.AddInfo("Miss", colornames.White)
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "充能防御-防御力",
		Timer: ToFrame(8),
		Field: FieldDef,
		Mul:   0.4,
	})
}

//=================ZhanShuFangYuSkill===============

func createZhanShuFangYuSkill(player *Player) ISkill {
	return &ZhanShuFangYuSkill{player: player, offsets: [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}}
}

type ZhanShuFangYuSkill struct {
	player  *Player
	offsets [][]int
}

func (z *ZhanShuFangYuSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerHurt && param.Player == z.player
}

func (z *ZhanShuFangYuSkill) Action(param *Param) {
	z.player.RecoverSkill(60)
	players := CollisionPlayers(z.player.Grid.Pos, 0, z.offsets)
	if len(players) > 0 {
		utils.RandomItem(players...).RecoverSkill(60)
	}
}

//==================JingTiaoSkill=================

func createJingTiaoSkill(player *Player) ISkill {
	return NewInitiativeSkill("精调", ToFrame(60), player, JingTiaoHandler)
}

func JingTiaoHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "精调-攻击速度",
		Timer: ToFrame(30),
		Field: FieldAttackSpeed,
		Mul:   -50,
	})
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "精调-攻击力",
		Timer: ToFrame(30),
		Field: FieldAtk,
		Mul:   1.5,
	})
}

//==================ZhiLiaoQiangHuaTypeBSkill================

func createZhiLiaoQiangHuaTypeBSkill(player *Player) ISkill {
	return NewInitiativeSkill("治疗强化B型", ToFrame(40), player, ZhiLiaoQiangHuaTypeBHandler)
}

func ZhiLiaoQiangHuaTypeBHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "治疗强化B型-攻击力",
		Timer: ToFrame(25),
		Field: FieldAtk,
		Mul:   0.15,
	})
}

//==================XunYiCaoSkill=============

func createXunYiCaoSkill(player *Player) ISkill {
	return &XunYiCaoSkill{player: player, timer: 60}
}

type XunYiCaoSkill struct {
	player *Player
	timer  int
}

func (x *XunYiCaoSkill) Update() {
	if x.timer > 0 {
		x.timer--
	} else {
		x.timer = 60
		players := PlayerManager.GetPlayers(IsAny)
		recoverValue := int64(x.player.BuffHolder.GetData(FieldAtk).Float() * 0.03)
		for i := 0; i < len(players); i++ {
			players[i].Recover(recoverValue)
		}
	}
}

func (x *XunYiCaoSkill) CanHandle(param *Param) bool {
	return false
}

func (x *XunYiCaoSkill) Action(param *Param) {

}

//======================TianMaShiYuSkill======================

func createTianMaShiYuSkill(player *Player) ISkill {
	return &TianMaShiYuSkill{player: player}
}

type TianMaShiYuSkill struct {
	player *Player
}

func (t *TianMaShiYuSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerInit && param.Player == t.player
}

func (t *TianMaShiYuSkill) Action(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "天马视域-攻击速度",
		Timer: -1,
		Field: FieldAttackSpeed,
		Add:   -20,
	})
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "天马视域-攻击力",
		Timer: -1,
		Field: FieldAtk,
		Mul:   0.5,
	})
	offsets := [][]int{{3, -1}, {4, -1}, {4, 0}, {3, 1}, {4, 1}}
	param.Player.Range.AddRange(offsets)
}

//====================GongJiLiQiangHuaTypeASkill====================

func createGongJiLiQiangHuaTypeASkill(player *Player) ISkill {
	return NewInitiativeSkill("攻击力强化A型", ToFrame(40), player, GongJiLiQiangHuaTypeAHandler)
}

func GongJiLiQiangHuaTypeAHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "攻击力强化A型-攻击力",
		Timer: ToFrame(30),
		Field: FieldAtk,
		Mul:   0.3,
	})
}

//======================XuLiGongJiSkill===========================

func createXuLiGongJiSkill(player *Player) ISkill {
	return &XuLiGongJiSkill{player: player}
}

type XuLiGongJiSkill struct {
	player *Player
	timer  int
}

func (x *XuLiGongJiSkill) Update() {
	x.timer++
}

func (x *XuLiGongJiSkill) CanHandle(param *Param) bool { // 自己攻击时
	return param.EventType == EventTypePlayerAttack && param.Player == x.player
}

func (x *XuLiGongJiSkill) Action(param *Param) {
	max := float64(ToFrame(2.5)) // 攻击 消耗蓄力
	param.Atk.Mul = 1.4 * math.Min(float64(x.timer), max) / max
	x.timer = 0
}

//====================YongChaoBeiGeSkill========================

func createYongChaoBeiGeSkill(player *Player) ISkill {
	return NewInitiativeSkill("涌潮悲歌", ToFrame(90), player, YongChaoBeiGeHandler)
}

func YongChaoBeiGeHandler(param *Param) {
	player := param.Player
	player.BuffHolder.AddBuff(&Buff{
		Name:  "涌潮悲歌-防御力",
		Timer: ToFrame(35),
		Field: FieldDef,
		Mul:   0.7,
	})
	player.BuffHolder.AddBuff(&Buff{
		Name:  "涌潮悲歌-攻击力",
		Timer: ToFrame(35),
		Field: FieldAtk,
		Mul:   0.7,
	})
	player.BuffHolder.AddBuff(&Buff{
		Name:  "涌潮悲歌-最大生命值",
		Timer: ToFrame(35),
		Field: FieldMaxHp,
		Mul:   0.7,
	})
	player.Recover(player.Data.Hp * 7 / 10)
}

type YongChaoBeiGeSkill struct { // 废弃
	*TimeSkill
	recoverValue int64
}

func NewYongChaoBeiGeSkill(player *Player) *YongChaoBeiGeSkill {
	res := &YongChaoBeiGeSkill{recoverValue: player.Data.Hp * 7 / 350}
	res.TimeSkill = NewTimeSkill(player, ToFrame(35), "涌潮悲歌-回血")
	return res
}

func (l *YongChaoBeiGeSkill) Update() {
	l.TimeSkill.Update()
	if l.Timer%60 == 0 {
		l.Player.Recover(l.recoverValue)
	}
}

func (l *YongChaoBeiGeSkill) CanHandle(param *Param) bool {
	return false
}

func (l *YongChaoBeiGeSkill) Action(param *Param) {

}

//====================YueLangJiSkill===================

func createYueLangJiSkill(player *Player) ISkill {
	return &YueLangJiSkill{player: player}
}

type YueLangJiSkill struct {
	player *Player
}

func (y *YueLangJiSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerInit && param.Player == y.player
}

func (y *YueLangJiSkill) Action(param *Param) {
	y.player.BuffHolder.AddBuff(&Buff{
		Name:  "跃浪击",
		Timer: ToFrame(15),
		Field: FieldAtk,
		Mul:   0.8,
	})
}

//=================XunJieDaJiTypeASkill================

func createXunJieDaJiTypeASkill(player *Player) ISkill {
	return NewInitiativeSkill("迅捷打击A型", ToFrame(45), player, XunJieDaJiTypeAHandler)
}

func XunJieDaJiTypeAHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "迅捷打击A型-攻击力",
		Timer: ToFrame(35),
		Field: FieldAtk,
		Mul:   0.2,
	})
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "迅捷打击A型-攻击速度",
		Timer: ToFrame(35),
		Field: FieldAttackSpeed,
		Add:   20,
	})
}

//==================ShenHaiLveShiZheSkill================

func createShenHaiLveShiZheSkill(player *Player) ISkill {
	return NewLifeBuffSkill(player, IsCareer(CareerGuards), &Buff{
		Name:  "深海掠食者",
		Timer: -1,
		Field: FieldAtk,
		Mul:   0.07,
	})
}

//==================LiZhiJuSkill===============

func createLiZhiJuSkill(player *Player) ISkill {
	return NewInitiativeSkill("力之锯", ToFrame(60), player, LiZhiJuHandler)
}

func LiZhiJuHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "力之锯-防御力",
		Timer: ToFrame(25),
		Field: FieldDef,
		Mul:   0.4,
	})
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "力之锯-攻击力",
		Timer: ToFrame(25),
		Field: FieldAtk,
		Mul:   0.65,
	})
	param.Player.AddSkill("力之锯-切割", NewLiZhiJuSkill(param.Player))
}

type LiZhiJuSkill struct {
	*TimeSkill
	offsets [][]int
}

func NewLiZhiJuSkill(player *Player) *LiZhiJuSkill {
	res := &LiZhiJuSkill{offsets: [][]int{{1, 0}}}
	res.TimeSkill = NewTimeSkill(player, ToFrame(25), "力之锯-切割")
	return res
}

func (l *LiZhiJuSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerAttack && param.Player == l.Player
}

func (l *LiZhiJuSkill) Action(param *Param) {
	param.Invalid = true // 取消普攻 进行范围攻击
	player := l.Player
	NewRangeTip(player.Grid.Pos, player.RangeDir, l.offsets, ColorCardMask)
	enemies := CollisionEnemies(player.Grid.Pos, player.RangeDir, l.offsets)
	param = &Param{EventType: EventTypePlayerSkill, Player: player}
	for i := 0; i < len(enemies); i++ {
		param.Enemy = enemies[i]
		param.Atk = player.BuffHolder.CloneData(FieldAtk)
		param.Def = enemies[i].BuffHolder.CloneData(FieldDef)
		PlayerManager.TriggerEvent(param)
		if !param.Invalid {
			enemies[i].Hurt(l.Player, ParseHurt(param))
		}
	}
}

//=====================JingJiSkill========================

func createJingJiSkill(player *Player) ISkill {
	return &JingJiSkill{player: player}
}

type JingJiSkill struct {
	player *Player
}

func (j *JingJiSkill) CanHandle(param *Param) bool { // 初始化时
	if param.EventType == EventTypePlayerInit && param.Player == j.player {
		j.player.BuffHolder.AddBuff(&Buff{
			Name:  "荆棘-防御力",
			Timer: -1,
			Field: FieldDef,
			Mul:   0.05,
		})
	} // 自己收到伤害
	return param.EventType == EventTypePlayerHurt && param.Player == j.player
}

func (j *JingJiSkill) Action(param *Param) { // 反甲
	param.Enemy.Hurt(j.player, j.player.Data.Atk/2)
}

//==================ZhanYiSkill======================

func createZhanYiSkill(player *Player) ISkill {
	return NewInitiativeSkill("战意", ToFrame(50), player, ZhanYiHandler)
}

func ZhanYiHandler(param *Param) {
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "战意-防御力",
		Timer: ToFrame(20),
		Field: FieldDef,
		Mul:   0.35,
	})
	param.Player.BuffHolder.AddBuff(&Buff{
		Name:  "战意-攻击力",
		Timer: ToFrame(20),
		Field: FieldAtk,
		Mul:   0.1,
	})
}

//=======================================

type TimeSkill struct { // 有时间限制的技能 到达时间进行移除
	Player  *Player
	Timer   int64
	name    string
	EndFunc func()
}

func NewTimeSkill(player *Player, timer int64, name string) *TimeSkill {
	return &TimeSkill{Player: player, Timer: timer, name: name}
}

func (t *TimeSkill) Update() {
	t.Timer--
	if t.Timer < 0 {
		t.Player.RemoveSkill(t.name)
		if t.EndFunc != nil {
			t.EndFunc()
		}
	}
}

//====================InitiativeSkill======================

type InitiativeSkill struct {
	name           string
	allTime, timer int64
	player         *Player
	handler        func(*Param)
}

func (i *InitiativeSkill) Recover(recover int64) {
	i.timer += recover
}

func NewInitiativeSkill(name string, allTime int64, player *Player, handler func(*Param)) *InitiativeSkill {
	return &InitiativeSkill{name: name, allTime: allTime, player: player, handler: handler, timer: 0}
}

func (i *InitiativeSkill) GetName() string {
	return i.name
}

func (i *InitiativeSkill) GetProgress() float64 {
	if i.timer >= i.allTime {
		return 1
	}
	return float64(i.timer) / float64(i.allTime)
}

func (i *InitiativeSkill) Update() {
	i.timer++
}

func (i *InitiativeSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerInitiative && param.Player == i.player
}

func (i *InitiativeSkill) Action(param *Param) {
	i.handler(param)
	i.timer = 0
}

//=======================LifeBuffSkill==========================

type LifeBuffSkill struct { // 存在的情况下给全体加Buff的技能
	player *Player
	filter func(*Player) bool
	buff   *Buff
}

func NewLifeBuffSkill(player *Player, filter func(*Player) bool, buff *Buff) *LifeBuffSkill {
	return &LifeBuffSkill{player: player, filter: filter, buff: buff}
}

func (l *LifeBuffSkill) CanHandle(param *Param) bool { // 自己登场 ，别人登场都要监听  且额外监听 自己退场事件
	return param.EventType == EventTypePlayerInit || (param.EventType == EventTypePlayerRetreat && param.Player == l.player)
}

func (l *LifeBuffSkill) Action(param *Param) {
	if param.EventType == EventTypePlayerInit {
		if param.Player == l.player { // 全场指定上Buff
			players := PlayerManager.GetPlayers(l.filter)
			for i := 0; i < len(players); i++ {
				players[i].BuffHolder.AddBuff(l.buff)
			}
		} else if l.filter(param.Player) { // 新来的满足 也加上
			param.Player.BuffHolder.AddBuff(l.buff)
		}
	} else { // 撤退 移除Buff
		players := PlayerManager.GetPlayers(l.filter)
		for i := 0; i < len(players); i++ {
			players[i].BuffHolder.RemoveBuff(l.buff.Name)
		}
	}
}

//========================TeZhongZuoZhanCeLveSkill========================

func createTeZhongZuoZhanCeLveSkill(player *Player) ISkill {
	return NewLifeBuffSkill(player, IsCareer(CareerReinstall), &Buff{
		Name:  "特种作战策略",
		Timer: -1,
		Field: FieldDef,
		Mul:   0.06,
	})
}

//=======================ZhanShuZhuangJiaSkill=====================

func createZhanShuZhuangJiaSkill(player *Player) ISkill {
	return &ZhanShuZhuangJiaSkill{player: player}
}

type ZhanShuZhuangJiaSkill struct {
	player *Player
}

func (z *ZhanShuZhuangJiaSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerHurt && param.Player == z.player
}

func (z *ZhanShuZhuangJiaSkill) Action(param *Param) {
	param.HurtValue *= (100 - 12) / 100
}

//=======================RetreatSkill=============================

func createRetreatSkill(player *Player) ISkill {
	return &RetreatSkill{player: player}
}

type RetreatSkill struct {
	player *Player
}

func (i *RetreatSkill) Update() {
}

func (i *RetreatSkill) GetName() string {
	return "撤退"
}

func (i *RetreatSkill) GetProgress() float64 { // 一直可用
	return 1
}

func (i *RetreatSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerInitiative && param.Player == i.player
}

func (i *RetreatSkill) Action(param *Param) {
	i.player.Retreat()
	GameManager.ChangePoint(i.player.Data.CostNum / 2)
}

//=====================LaunchSkill=====================

type LaunchSkill struct {
	player  *Player
	handler func(*Param)
}

func NewLaunchSkill(player *Player, handler func(*Param)) *LaunchSkill {
	return &LaunchSkill{player: player, handler: handler}
}

func (y *LaunchSkill) CanHandle(param *Param) bool {
	return param.EventType == EventTypePlayerInit && param.Player == y.player
}

func (y *LaunchSkill) Action(param *Param) {
	y.handler(param)
}

//=================ComboSkill=================

type ComboSkill struct {
	Melee          bool
	Atk            *DataWarp
	EndFunc        func()
	player         *Player
	Enemy          *Enemy
	timer, allTime int
	count          int
	name           string
}

func (t *ComboSkill) CanHandle(param *Param) bool {
	return false
}

func (t *ComboSkill) Action(param *Param) {
}

func NewComboSkill(player *Player, interval, count int, name string) *ComboSkill {
	return &ComboSkill{Melee: true, player: player, Atk: player.BuffHolder.CloneData(FieldAtk), allTime: interval, count: count, name: name}
}

func (t *ComboSkill) Update() {
	if t.timer > 0 {
		t.timer--
	} else {
		if t.Enemy != nil && t.Enemy.Die {
			t.player.RemoveSkill(t.name)
			return
		}
		t.timer = t.allTime
		t.trigger()
		t.count--
		if t.count <= 0 {
			t.player.RemoveSkill(t.name)
			if t.EndFunc != nil {
				t.EndFunc()
			}
		}
	}
}

func (t *ComboSkill) trigger() {
	enemy := t.Enemy
	if enemy == nil {
		enemy = GetMinLastEnemy(t.player.GetEnemies())
	}
	if enemy == nil {
		return
	}
	param := &Param{EventType: EventTypePlayerSkill, Player: t.player, Atk: t.Atk.Clone(),
		Enemy: enemy, Def: enemy.BuffHolder.CloneData(FieldDef)}
	PlayerManager.TriggerEvent(param)
	if param.Invalid {
		return
	}
	if t.Melee {
		enemy.Hurt(t.player, ParseHurt(param))
	} else {
		NewPlayerBullet(t.player.Grid.Pos+GridSize/2, ParseHurt(param), enemy, t.player)
	}
}
