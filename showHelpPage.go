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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func showHelpPage() fyne.CanvasObject {
	usage := `
请在执行任何操作前均保证您已仔细阅读此说明。
该软件非官方版本，全部题目来源于网络。如因对此软件使用不当造成考试失利，使用者自行承担一切后果。

使用帮助：

	在“练习”页面下进行模拟考试练习，答题完毕后点击右下角“结束答题”进行提交；
	错题会保存在 “错题.xlsx” 中方便查看与反复练习；
	修改 “联考题库.xlsx” 文件可以自定义题库；
	欢迎通过下方电子邮箱或 QQ:2200474484 汇报程序 BUG 或提出功能建议。

待添加功能：

	- 答题一轮后继续答题

已知BUG：

	- 无E选项的题目会显示可选择的空白E选项。

更新日志：
	2021-3-17 v0.2: 添加窗口图标，添加导出错题功能，添加题库相关错误提示，优化算法，修复bug，提升用户体验
	2021-3-16 v0.1: 编写基础功能

Author: 3RM1N3 @ 18级麻醉5班·牡丹江医学院
E-mail: wangyu7439@hotmail.com

捐赠通道：
开发软件不易，如果此软件帮到了你，请考虑通过扫描下方二维码捐赠通道给我买一杯咖啡。
`

	usageLabel := widget.NewLabel(usage)
	usageLabel.Wrapping = fyne.TextWrapWord

	alipay := canvas.NewImageFromResource(resourceALiPayJpg)
	alipay.FillMode = canvas.ImageFillOriginal
	wechatpay := canvas.NewImageFromResource(resourceWeChatPayJpg)
	wechatpay.FillMode = canvas.ImageFillOriginal

	return container.NewBorder(
		nil,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("版权所有 © 2021 @github.com/3RM1N3。保留所有权利。\nCopyright © 2021 @github.com/3RM1N3. All Rights Reserved."),
			layout.NewSpacer(),
		),
		nil,
		nil,
		container.NewVScroll(container.NewVBox(
			usageLabel,
			container.NewHBox(layout.NewSpacer(), alipay, layout.NewSpacer(), wechatpay, layout.NewSpacer()),
		)),
	)
}
