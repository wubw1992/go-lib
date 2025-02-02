/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package display

import (
	"github.com/linuxdeepin/go-lib/backlight/common"
	"strings"
)

const controllersDir = "/sys/class/backlight"

type Controller struct {
	*common.Controller
	Type       ControllerType
	DeviceEDID []byte
}

func NewController(path string) (*Controller, error) {
	basec, err := common.NewController(path)
	if err != nil {
		return nil, err
	}

	c := &Controller{
		Controller: basec,
		Type:       ControllerTypeUnknown,
	}

	typeStr, err := c.GetPropString("type")
	if err != nil {
		return nil, err
	}
	c.Type = ControllerTypeFromString(strings.TrimSpace(typeStr))

	c.DeviceEDID, _ = c.GetPropBytes("device/edid")
	return c, nil
}

func (c *Controller) GetActualBrightness() (int, error) {
	brightness, err := c.GetPropInt("actual_brightness")
	if err != nil {
		return 0, err
	}
	return brightness, nil
}

type Controllers []*Controller

func List() (Controllers, error) {
	return list(controllersDir)
}

func list(dir string) (Controllers, error) {
	paths, err := common.ListControllerPaths(dir)
	if err != nil {
		return nil, err
	}
	controllers := make(Controllers, 0, len(paths))
	for _, path := range paths {
		c, err := NewController(path)
		if err != nil {
			continue
		}
		controllers = append(controllers, c)
	}
	return controllers, nil
}

func (cs Controllers) GetByEDID(edid []byte) *Controller {
	if len(edid) == 0 {
		return nil
	}
	for _, c := range cs {
		if byteSliceEqual(c.DeviceEDID, edid) {
			return c
		}
	}
	return nil
}

func (cs Controllers) GetByName(name string) *Controller {
	for _, c := range cs {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func byteSliceEqual(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
