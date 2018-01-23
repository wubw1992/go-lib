/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package appinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/BurntSushi/xgbutil"
)

var (
	pid      int
	prog     string
	hostname string
)

func init() {
	pid = os.Getpid()
	prog = filepath.Base(os.Args[0])
	hostname, _ = os.Hostname()
}

type AppLaunchContext struct {
	sync.Mutex
	xu          *xgbutil.XUtil
	count       uint
	timestamp   uint32
	cmdPrefixes []string
}

func NewAppLaunchContext(xu *xgbutil.XUtil) *AppLaunchContext {
	return &AppLaunchContext{
		xu: xu,
	}
}

func (ctx *AppLaunchContext) SetTimestamp(timestamp uint32) {
	ctx.timestamp = timestamp
}

func (ctx *AppLaunchContext) GetTimestamp() uint32 {
	return ctx.timestamp
}

func (ctx *AppLaunchContext) SetCmdPrefixes(v []string) {
	ctx.cmdPrefixes = v
}

func (ctx *AppLaunchContext) GetCmdPrefixes() []string {
	return ctx.cmdPrefixes
}

func (ctx *AppLaunchContext) GetStartupNotifyId(appInfo AppInfo, files []string) (string, error) {
	execBase := filepath.Base(appInfo.GetExecutable())
	snId := fmt.Sprintf("%s-%d-%s-%s-%d_TIME%d", prog, pid, hostname, execBase, ctx.count, ctx.timestamp)
	ctx.count++

	screenStr := strconv.Itoa(ctx.xu.Conn().DefaultScreen)
	// send new msg
	msg := &StartupNotifyMessage{
		Type: "new",
		KeyValues: map[string]string{
			"ID":     snId,
			"SCREEN": screenStr,
			"NAME":   appInfo.GetName(),
		},
	}
	startupWMClass := appInfo.GetStartupWMClass()
	if startupWMClass != "" {
		msg.KeyValues["WMCLASS"] = startupWMClass
	}

	err := msg.Broadcast(ctx.xu)
	if err != nil {
		return "", err
	}
	return snId, nil
}

func (ctx *AppLaunchContext) LaunchFailed(startupNotifyId string) error {
	// send remove msg
	msg := &StartupNotifyMessage{
		Type: "remove",
		KeyValues: map[string]string{
			"ID": startupNotifyId,
		},
	}
	return msg.Broadcast(ctx.xu)
}
