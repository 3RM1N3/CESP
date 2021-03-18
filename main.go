/*
Copyright 2021 王瑀
版权所有 2021 王瑀

This file is part of CESP.

    CESP is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    CESP is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with CESP.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/flopp/go-findfont"
)

var qst100 [][]string        // 抽取的100道题
var dict = map[int]bool{}    // 用于判断题目做过与否
var wrongList = [][]string{} // 错题列表

func init() {
	fontList := findfont.List() // 查找系统中的黑体中文字体，设置为程序字体。
	for _, font := range fontList {
		if strings.Contains(font, "simhei.ttf") {
			os.Setenv("FYNE_FONT", font)
			break
		}
	}
}

func main() {
	log.Println("START")
	a := app.New()

	w := a.NewWindow("医学类专业基础理论综合水平模拟考试平台-UNOFFICIAL-v0.2")

	loadDataWindow := a.NewWindow("综合考试模拟平台")
	loadDataWindow.SetIcon(resourceWindowIconPng)
	loadDataWindow.Resize(fyne.NewSize(400, 200)) // 设置窗口大小
	loadDataWindow.SetFixedSize(true)             // 禁止改变窗口大小

	infoLabel := widget.NewLabel("加载题库中......\n请勿关闭此窗口。")
	infoLabel.Alignment = fyne.TextAlignCenter
	infoLabel.Wrapping = fyne.TextWrapWord
	infoLabel.Resize(fyne.NewSize(400, 150))
	confirmButton := widget.NewButton("确定", func() { os.Exit(1) })
	confirmButton.Hide()

	loadDataWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		infoLabel,
		confirmButton,
		layout.NewSpacer(),
	))

	// 随机抽取100道题
	go func() {
		time.Sleep(time.Second)
		f, err := excelize.OpenFile("联考题库.xlsx")
		if err != nil {
			infoLabel.SetText(fmt.Sprint("题库文件打开失败：", err))
			confirmButton.Show()
			return
		}
		rows, err := f.GetRows("题库")
		if err != nil {
			infoLabel.SetText(fmt.Sprint("读取表格失败：", err))
			confirmButton.Show()
			return
		}
		length := len(rows)
		if length < 101 {
			infoLabel.SetText(fmt.Sprint("题库内题目数量少于100道！", err))
			confirmButton.Show()
			return
		}
		if length == 101 {
			qst100 = rows[1:]
			//goto MAIN
		}
		for len(dict) < 100 {
			rand.Seed(time.Now().UnixNano())
			randInt := rand.Intn(length-1) + 1
			dict[randInt] = false
			time.Sleep(time.Microsecond)
		}
		for k := range dict {
			qst100 = append(qst100, rows[k])
		}
		// 顺利加载
		infoLabel.SetText("加载完毕。")
		//time.Sleep(time.Second)
		createMainWindow(&w)
		loadDataWindow.Close()
		w.Show()
	}()
	loadDataWindow.ShowAndRun()
}

func createMainWindow(w *fyne.Window) {
	(*w).SetMaster()
	(*w).SetIcon(resourceWindowIconPng)
	(*w).Resize(fyne.NewSize(650, 450)) // 设置窗口大小
	(*w).SetFixedSize(true)             // 禁止改变窗口大小

	(*w).SetContent(container.NewAppTabs(
		container.NewTabItem("练习", examPage(*w)),
		container.NewTabItem("使用帮助&关于", showHelpPage()),
	))

	warning := `请在执行任何操作前均保证您已仔细阅读此说明。
该软件非官方版本，全部题目来源于网络。
如因对此软件使用不当造成考试失利，使用者自行承担一切后果。
*点击“确定”按钮即代表您同意上述所有条约。*`
	cnf := dialog.NewConfirm("警告", warning, func(r bool) {
		if !r {
			os.Exit(0)
		}
	}, *w)
	cnf.SetDismissText("取消")
	cnf.SetConfirmText("确定")

	cnf.Show()
}
