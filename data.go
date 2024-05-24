/*
@author: sk
@date: 2023/2/25
*/
package main

import R "arknights/res"

func GetCardDatas() []*CardData {
	// 带实现  记录特点
	//星熊,
	//斯卡蒂
	//白金,
	//调香师,
	//雷蛇,     有眩晕
	//银灰,
	//阿米娅,    有强制退出
	//陈,
	//赫默,      可以召唤医疗无人机
	//远山,
	//塞雷娅,    状态，技能，减移动速度(每秒记录区域内的敌人，新的改变移动速度，不存在的恢复移动速度)
	//能天使
	res := make([]*CardData, 0)
	res = append(res, &CardData{
		Name:        "星熊",
		Img:         R.MAIN.PLAYER.PLAYER1,
		Place:       PlaceLand,
		ChangePlace: PlaceNone,
		Atk:         221,
		Def:         257,
		BlockNum:    3,
		CostNum:     19,
		Hp:          1602,
		AttackTime:  ToFrame(1.2),
		CoolTime:    ToFrame(70),
		Skills:      []string{"战术装甲", "特种作战策略", "战意", "荆棘", "力之锯", "撤退"},
		Range:       [][]int{{0, 0}, {1, 0}},
		Attack:      "近战攻击",
		Career:      CareerReinstall,
	})
	res = append(res, &CardData{
		Name:        "斯卡蒂",
		Img:         R.MAIN.PLAYER.PLAYER2,
		Place:       PlaceLand,
		ChangePlace: PlaceNone,
		Atk:         452,
		Def:         116,
		BlockNum:    1,
		CostNum:     17,
		Hp:          1521,
		AttackTime:  ToFrame(1.5),
		CoolTime:    ToFrame(60),
		Skills:      []string{"深海掠食者", "迅捷打击A型", "跃浪击", "涌潮悲歌", "撤退"},
		Range:       [][]int{{0, 0}, {1, 0}},
		Attack:      "近战攻击",
		Career:      CareerGuards,
	})
	res = append(res, &CardData{
		Name:        "白金",
		Img:         R.MAIN.PLAYER.PLAYER3,
		Place:       PlaceHighland,
		ChangePlace: PlaceNone,
		Atk:         171,
		Def:         58,
		BlockNum:    1,
		CostNum:     10,
		Hp:          693,
		AttackTime:  ToFrame(1),
		CoolTime:    ToFrame(70),
		Skills:      []string{"蓄力攻击", "攻击力强化A型", "天马视域", "撤退"},
		Range:       [][]int{{0, -1}, {1, -1}, {2, -1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {1, 1}, {2, 1}},
		Attack:      "远程攻击",
		Career:      CareerSniper,
	})
	res = append(res, &CardData{
		Name:        "调香师",
		Img:         R.MAIN.PLAYER.PLAYER4,
		Place:       PlaceHighland,
		ChangePlace: PlaceNone,
		Atk:         117,
		Def:         69,
		BlockNum:    1,
		CostNum:     13,
		Hp:          710,
		AttackTime:  ToFrame(2.85),
		CoolTime:    ToFrame(70),
		Skills:      []string{"熏衣草", "治疗强化B型", "精调", "撤退"},
		Range:       [][]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {0, 0}, {1, 0}, {2, 0}, {-1, 1}, {0, 1}, {1, 1}},
		Attack:      "群体治疗",
		Career:      CareerMedical,
	})
	res = append(res, &CardData{
		Name:        "雷蛇",
		Img:         R.MAIN.PLAYER.PLAYER5,
		Place:       PlaceLand,
		ChangePlace: PlaceNone,
		Atk:         219,
		Def:         256,
		BlockNum:    3,
		CostNum:     17,
		Hp:          1307,
		AttackTime:  ToFrame(1.2),
		CoolTime:    ToFrame(70),
		Skills:      []string{"战术防御", "充能防御", "反击电弧", "撤退"},
		Range:       [][]int{{0, 0}, {1, 0}, {2, 0}},
		Attack:      "近战攻击",
		Career:      CareerReinstall,
	})
	res = append(res, &CardData{
		Name:        "银灰",
		Img:         R.MAIN.PLAYER.PLAYER6,
		Place:       PlaceLand,
		ChangePlace: PlaceNone,
		Atk:         297,
		Def:         189,
		BlockNum:    2,
		CostNum:     17,
		Hp:          1075,
		AttackTime:  ToFrame(1.3),
		CoolTime:    ToFrame(70),
		Skills:      []string{"领袖", "强力击A型", "雪境生存法则", "真银斩", "撤退"},
		Range:       [][]int{{0, -1}, {1, -1}, {0, 0}, {1, 0}, {2, 0}, {0, 1}, {1, 1}},
		Attack:      "远近攻击",
		Career:      CareerGuards,
	})
	res = append(res, &CardData{
		Name:        "阿米娅",
		Img:         R.MAIN.PLAYER.PLAYER7,
		Place:       PlaceHighland,
		ChangePlace: PlaceNone,
		Atk:         276,
		Def:         48,
		BlockNum:    1,
		CostNum:     18,
		Hp:          888,
		AttackTime:  ToFrame(1.6),
		CoolTime:    ToFrame(70),
		Skills:      []string{"情绪吸收", "战术咏唱A型", "精神爆发", "奇美拉", "撤退"},
		Range:       [][]int{{0, -1}, {1, -1}, {2, -1}, {0, 0}, {1, 0}, {2, 0}, {0, 1}, {1, 1}, {2, 1}},
		Attack:      "远程攻击",
		Career:      CareerMagian,
	})
	res = append(res, &CardData{
		Name:        "陈",
		Img:         R.MAIN.PLAYER.PLAYER8,
		Place:       PlaceLand,
		ChangePlace: PlaceNone,
		Atk:         249,
		Def:         154,
		BlockNum:    2,
		CostNum:     18,
		Hp:          1229,
		AttackTime:  ToFrame(1.3),
		CoolTime:    ToFrame(70),
		Skills:      []string{"呵斥", "持刀格斗术", "鞘击", "赤霄·拔刀", "赤霄·绝影", "撤退"},
		Range:       [][]int{{0, 0}, {1, 0}},
		Attack:      "近战双攻",
		Career:      CareerGuards,
	})
	res = append(res, &CardData{
		Name:        "赫默",
		Img:         R.MAIN.PLAYER.PLAYER9,
		Place:       PlaceHighland,
		ChangePlace: PlaceNone,
		Atk:         166,
		Def:         62,
		BlockNum:    1,
		CostNum:     16,
		Hp:          845,
		AttackTime:  ToFrame(2.85),
		CoolTime:    ToFrame(70),
		Skills:      []string{"强化注射", "治疗强化A型", "医疗无人机", "撤退"},
		Range:       [][]int{{0, -1}, {1, -1}, {2, -1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {1, 1}, {2, 1}},
		Attack:      "单体治疗",
		Career:      CareerMedical,
	})
	res = append(res, &CardData{
		Name:        "远山",
		Img:         R.MAIN.PLAYER.PLAYER10,
		Place:       PlaceHighland,
		ChangePlace: PlaceNone,
		Atk:         332,
		Def:         47,
		BlockNum:    1,
		CostNum:     28,
		Hp:          653,
		AttackTime:  ToFrame(2.9),
		CoolTime:    ToFrame(70),
		Skills:      []string{"占卜", "战术咏唱B型", "命运", "撤退"},
		Range:       [][]int{{0, -1}, {1, -1}, {0, 0}, {1, 0}, {2, 0}, {0, 1}, {1, 1}},
		Attack:      "区域攻击",
		Career:      CareerMagian,
	})
	res = append(res, &CardData{
		Name:        "塞雷娅",
		Img:         R.MAIN.PLAYER.PLAYER11,
		Place:       PlaceLand,
		ChangePlace: PlaceNone,
		Atk:         200,
		Def:         248,
		BlockNum:    3,
		CostNum:     17,
		Hp:          1309,
		AttackTime:  ToFrame(1.2),
		CoolTime:    ToFrame(70),
		Skills:      []string{"莱茵充能护服", "精神回复", "急救", "药物配置", "钙质化", "撤退"},
		Range:       [][]int{{0, 0}},
		Attack:      "近战攻击",
		Career:      CareerReinstall,
	})
	res = append(res, &CardData{
		Name:        "能天使",
		Img:         R.MAIN.PLAYER.PLAYER12,
		Place:       PlaceHighland,
		ChangePlace: PlaceNone,
		Atk:         183,
		Def:         57,
		BlockNum:    1,
		CostNum:     11,
		Hp:          711,
		AttackTime:  ToFrame(1),
		CoolTime:    ToFrame(70),
		Skills:      []string{"快速弹匣", "天使的祝福", "冲锋模式", "扫射模式", "过载模式", "撤退"},
		Range:       [][]int{{0, -1}, {1, -1}, {2, -1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {1, 1}, {2, 1}},
		Attack:      "远程攻击",
		Career:      CareerSniper,
	})
	return res
}

func GetEnemiesDatas() map[string]*EnemyData {
	res := make(map[string]*EnemyData)
	// 源石虫,猎狗,士兵, 重装防御者，伐木机
	// 弩手 远程攻击
	// 鸡尾酒投掷者 物理群伤
	// 弑君者 不可阻挡
	res["源石虫"] = &EnemyData{
		Name:       "源石虫",
		Img:        R.MAIN.ENEMY.ENEMY1,
		Hp:         550,
		Atk:        130,
		Def:        0,
		MoveSpeed:  ToMoveSpeed(1),
		AttackTime: ToFrame(1.7),
		CanBlock:   true,
		HurtLife:   1,
		Attack:     "近战攻击",
	}
	res["猎狗"] = &EnemyData{
		Name:       "猎狗",
		Img:        R.MAIN.ENEMY.ENEMY2,
		Hp:         820,
		Atk:        190,
		Def:        0,
		MoveSpeed:  ToMoveSpeed(1.9),
		AttackTime: ToFrame(1.4),
		CanBlock:   true,
		HurtLife:   1,
		Attack:     "近战攻击",
	}
	res["士兵"] = &EnemyData{
		Name:       "士兵",
		Img:        R.MAIN.ENEMY.ENEMY3,
		Hp:         1650,
		Atk:        200,
		Def:        100 / 2,
		MoveSpeed:  ToMoveSpeed(1.1),
		AttackTime: ToFrame(2),
		CanBlock:   true,
		HurtLife:   1,
		Attack:     "近战攻击",
	}
	res["弩手"] = &EnemyData{
		Name:        "弩手",
		Img:         R.MAIN.ENEMY.ENEMY4,
		Hp:          1400,
		Atk:         140, //  240
		Def:         100 / 2,
		MoveSpeed:   ToMoveSpeed(0.9),
		AttackTime:  ToFrame(2.4),
		AttackRange: ToAttackRange(1.9),
		CanBlock:    true,
		HurtLife:    1,
		Attack:      "远程攻击",
	}
	res["重装防御者"] = &EnemyData{
		Name:       "重装防御者",
		Img:        R.MAIN.ENEMY.ENEMY5,
		Hp:         6000,
		Atk:        600 / 2,
		Def:        800 / 4,
		MoveSpeed:  ToMoveSpeed(0.75),
		AttackTime: ToFrame(2.6),
		CanBlock:   true,
		HurtLife:   1,
		Attack:     "近战攻击",
	}
	res["鸡尾酒投掷者"] = &EnemyData{
		Name:        "鸡尾酒投掷者",
		Img:         R.MAIN.ENEMY.ENEMY6,
		Hp:          1550,
		Atk:         180,
		Def:         50 / 2,
		MoveSpeed:   ToMoveSpeed(1),
		AttackTime:  ToFrame(2.7),
		AttackRange: ToAttackRange(1.75),
		CanBlock:    true,
		HurtLife:    1,
		Attack:      "远程群攻",
	}
	res["伐木机"] = &EnemyData{
		Name:       "伐木机",
		Img:        R.MAIN.ENEMY.ENEMY7,
		Hp:         8000,
		Atk:        750 / 2,
		Def:        80 / 2,
		MoveSpeed:  ToMoveSpeed(0.75),
		AttackTime: ToFrame(3.3),
		CanBlock:   true,
		HurtLife:   1,
		Attack:     "近战攻击",
	}
	res["弑君者"] = &EnemyData{
		Name:       "弑君者",
		Img:        R.MAIN.ENEMY.ENEMY8,
		Hp:         6000,
		Atk:        400,
		Def:        120 / 2,
		MoveSpeed:  ToMoveSpeed(1.4),
		AttackTime: ToFrame(2.8),
		CanBlock:   false,
		HurtLife:   1,
		Attack:     "近战攻击",
	}
	return res
}

func GetEnemyWaves() [][]string { // 获取敌人 出地图  最开始1个  之后4次再变换组合  每次变换组合加1个  出兵时间固定9s一波
	arr := []string{"源石虫", "猎狗", "士兵", "弩手", "重装防御者", "鸡尾酒投掷者", "伐木机", "弑君者"}
	res := make([][]string, 0)
	temp := make([]string, 0)
	for i := 0; i < len(arr); i++ {
		temp = append(temp, arr[i])
		for j := 0; j < 8-i; j++ {
			res = append(res, temp)
		}
	}
	return res
}
