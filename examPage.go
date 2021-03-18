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
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func examPage(win fyne.Window) fyne.CanvasObject {
	qstNum := 1
	selected := ""
	currentPage := qst100[0]
	rightAsr := currentPage[7]
	rightNum, wrongNum := 0, 0 // 定义正确与错误题目数量

	statusBar := widget.NewLabel(fmt.Sprintf("当前为第 %d 道题，还剩 %d 道。", qstNum, 100-qstNum))

	question := widget.NewLabel(fmt.Sprintf("%d.%s", qstNum, currentPage[1])) // 题目
	question.Wrapping = fyne.TextWrapWord                                     // 自动拆行

	radio := widget.NewRadioGroup(currentPage[2:7], func(s string) { // 选项
		selected = s
	})

	preQst := widget.NewButton("上一题", func() {
		currentQstIndex, _ := strconv.Atoi(currentPage[0])
		if qstNum == 1 {
			return
		} else if !dict[currentQstIndex] {
			if selected == "" {
				dialog.ShowInformation("错误", "尚未选择任何选项！", win)
				return
			}
			if selected[:1] == rightAsr { // 判断选择是否正确
				log.Println("回答正确")
				rightNum++
			} else {
				log.Println("回答错误")
				currentPage = append(currentPage, selected[:1])
				wrongList = append(wrongList, currentPage)
				wrongNum++
			}
			selected = ""
			dict[currentQstIndex] = true
		}
		qstNum--
		currentPage = qst100[qstNum-1]
		question.SetText(fmt.Sprintf("%d.%s", qstNum, currentPage[1])) // 设置上一题目
		radio.Options = currentPage[2:7]                               // 设置上一题选项
		radio.Refresh()                                                // 刷新选项
		radio.Disable()                                                // 禁用选项
		rightAsr = currentPage[7]                                      // 设置上一题正确答案
		statusBar.SetText(fmt.Sprintf("当前为第 %d 道题，还剩 %d 道。", qstNum, 100-qstNum))
	})
	nextQst := widget.NewButton("下一题", func() {
		currentQstIndex, _ := strconv.Atoi(currentPage[0])
		if qstNum == 100 { // 判断是否为最后一道题
			if selected == "" {
				return
			}
			if judge(rightAsr, selected) { // 判断对错
				rightNum++
			} else {
				wrongNum++
				currentPage = append(currentPage, selected[:1])
				wrongList = append(wrongList, currentPage)
			}
			dict[currentQstIndex] = true // 设置为已做过
			selected = ""
			return
		}

		nextQstIndex, _ := strconv.Atoi(qst100[qstNum][0]) // 判断下一题做过与否
		if !dict[nextQstIndex] {
			defer radio.Enable() // 函数返回后取消禁用单选框
		}

		if radio.Disabled() { // 判断该题做过与否
			goto LABEL1
		}
		if selected == "" {
			dialog.ShowInformation("错误", "尚未选择任何选项！", win)
			return
		}
		if judge(rightAsr, selected) { // 判断对错
			rightNum++
		} else {
			wrongNum++
			currentPage = append(currentPage, selected[:1])
			wrongList = append(wrongList, currentPage)
		}
		dict[currentQstIndex] = true
	LABEL1:
		qstNum++
		currentPage = qst100[qstNum-1]
		selected = ""
		question.SetText(fmt.Sprintf("%d.%s", qstNum, currentPage[1])) // 设置下一题目
		radio.Options = currentPage[2:7]                               // 设置下一题选项
		radio.Refresh()                                                // 刷新选项
		rightAsr = currentPage[7]                                      // 设置下一题正确答案
		statusBar.SetText(fmt.Sprintf("当前为第 %d 道题，还剩 %d 道。", qstNum, 100-qstNum))
	})

	quit := widget.NewButton("结束答题", func() { // 结束答题按钮
		if rightNum == 0 && wrongNum == 0 {
			dialog.ShowInformation("提示", "还没有作答哦！", win)
			return
		}
		info := ""
		if rightNum+wrongNum != 100 {
			info = "尚有题目未作答！\n"
		}
		info += fmt.Sprintf("此次答题正确%d个，错误%d个。\n按确定以导出错题，按取消键继续答题。", rightNum, wrongNum)
		cnf := dialog.NewConfirm("提示", info, func(r bool) {
			if !r {
				return
			}
			outputWrong(win) // 导出错题
		}, win)
		cnf.SetDismissText("取消")
		cnf.SetConfirmText("确定")
		cnf.Show()
	})

	return container.NewBorder(
		question,
		container.NewVBox(
			container.NewHBox(preQst, nextQst, layout.NewSpacer(), quit),
			statusBar,
		),
		nil,
		nil,
		container.NewHScroll(radio),
	)
}

func judge(rightAsr, usersChoice string) bool {
	if rightAsr == usersChoice[:1] {
		return true
	}
	return false
}

func outputWrong(win fyne.Window) {
	err := errors.New("")
	title := []string{}
	if len(wrongList) == 0 {
		os.Exit(0)
	}
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", "错题")
	title = []string{"题库题号", "题目", "A选项", "B选项", "C选项", "D选项", "E选项", "正确答案", "我的答案"}
	for x, cell := range title {
		xName, _ := excelize.ColumnNumberToName(x + 1)
		axis := fmt.Sprintf("%s%d", xName, 1)
		f.SetCellValue("错题", axis, cell)
	}

	for y, lines := range wrongList {
		for x, cell := range lines {
			xName, _ := excelize.ColumnNumberToName(x + 1)
			axis := fmt.Sprintf("%s%d", xName, y+2)
			f.SetCellValue("错题", axis, cell)
		}
	}
	err = f.SaveAs("错题.xlsx")
	if err != nil {
		dialog.ShowInformation("注意", "错题文件保存失败！请检查“错题.xlsx”是否被占用！", win)
		return
	}
	os.Exit(0)
}
