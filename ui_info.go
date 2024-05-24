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
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.UI_INFO, createUIInfo)
}

func createUIInfo(data *model.ObjectData) model.IObject {
	res := &uiInfo{}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	UIInfo = res
	return res
}

type uiInfo struct {
	*object.PointObject
}

func (u *uiInfo) Draw(screen *ebiten.Image) {
	waveStr := fmt.Sprintf("Wave:%d/%d", EnemyManager.WaveIndex, EnemyManager.GetAllWave())
	utils.DrawAnchorText(screen, waveStr, 1280/2-128+16i, 0.5+0.5i, Font32, colornames.White)
	lifeStr := fmt.Sprintf("Life:%d", GameManager.Life)
	utils.DrawAnchorText(screen, lifeStr, 1280/2+128+16i, 0.5+0.5i, Font32, colornames.Red)
	pointStr := fmt.Sprintf("P:%02d", GameManager.Point)
	utils.DrawAnchorText(screen, pointStr, complex(1280-128/2, 64*9+32+48/2), 0.5+0.5i, Font48,
		colornames.White)
	utils.FillRect(screen, complex(1280-128, 64*9+32+48-4), complex(128*GameManager.Progress, 4),
		config.ColorWhite)
	lastStr := fmt.Sprintf("剩余可放置角色:%d", GameManager.LastNum)
	utils.DrawAnchorText(screen, lastStr, complex(1280-128*3, 64*9+32+48-32/2), 0.5i, Font32, colornames.White)
}
