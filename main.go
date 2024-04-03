package main

import (
	"errors"
	"fmt"
	"fyne-SIA/SIA"
	"fyne-SIA/icon"
	"fyne-SIA/network"
	"net"
	"strconv"

	_ "log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var siaConn *network.Connection
var iSeq uint16

var button4Connect *widget.Button

var logsBox *widget.Entry

var inputRrcvr *widget.Entry
var inputLPref *widget.Entry
var inputAcct *widget.Entry
var inputData *widget.Entry

func addSystemTray(a fyne.App, w fyne.Window) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("SIA test tool", func() { w.Show() /* 显示窗口 */ })
		h.Icon = theme.HomeIcon() //图标
		m := fyne.NewMenu("MyApp", h)
		desk.SetSystemTrayMenu(m)
	}
	w.SetCloseIntercept(func() {
		w.Hide()
	})
}

func main() {

	a := app.New()
	w := a.NewWindow("SIA test tool")

	addSystemTray(a, w)

	//demo act as a SIA device or SIA platrom
	typeSelect := widget.NewSelect([]string{"Device", "Platform"}, func(value string) {
		//fmt.Println("type:", value)
		if value == "Platform" {

			err := errors.New("Not support Platform")
			dialog.ShowError(err, w)
			return
		}
	})

	// defalut vaule of typeSelect
	typeSelect.SetSelected("Device")

	typeBox := container.NewHBox(widget.NewLabel("Device/Platform"), typeSelect)

	//SIA platform IP and port
	inputIP := widget.NewEntry()
	inputIP.SetPlaceHolder("Input HCP address like 127.0.0.1")
	//inputIP.Resize(fyne.NewSize(40, 100))

	//inputIP.SetText("127.0.0.1")

	inputPort := widget.NewEntry()
	inputPort.SetPlaceHolder("Input HCP port for SIA like 15664")
	//inputPort.SetText("15664")

	button4Connect = widget.NewButton("Connect HCP", func() {

		//set button to disconnect HCP
		if (siaConn != nil) && (!siaConn.GetConnStatus()) {
			siaConn.Stop()
			siaConn = nil

			if button4Connect != nil {
				button4Connect.SetText("Connect HCP")
			}

			//println("Connect HCP")

		} else {

			//else if (siaConn != nil) && (siaConn.GetConnStatus()) {
			//the SIA server may reject the data and the SIA connection will be closed automatically.

			//start to connect to SIA platform, only one connection with id always equeals to 1
			_, err := strconv.Atoi(inputPort.Text)
			if err != nil {
				err := errors.New(fmt.Sprintf("Invalid port:%s", err.Error()))
				dialog.ShowError(err, w)
				return
			}

			//fmt.Println("inputIP.Text:", inputIP.Text)
			//var rAddr *net.TCPAddr = &net.TCPAddr{IP: []byte(inputIP.Text), Port: nPort}
			rAddr, err := net.ResolveTCPAddr("tcp4", inputIP.Text+":"+inputPort.Text)
			if err != nil {
				err := errors.New(fmt.Sprintf("connect error:%s", err.Error()))
				dialog.ShowError(err, w)
				return
			}

			conn, err := net.DialTCP("tcp4", nil, rAddr)
			if err != nil {
				err := errors.New(fmt.Sprintf("connect error:%s", err.Error()))
				dialog.ShowError(err, w)
				return
			}
			siaConn = network.NewConnection(conn, 1)

			siaConn.Start()

			if button4Connect != nil {
				button4Connect.SetText("Disconnect HCP")
			}
			//println("Disconnect HCP")

			//for
			if logsBox != nil {
				go func(input *widget.Entry) {
					//fmt.Println("11111")

					for {
						//fmt.Println("11111")

						if siaConn != nil {
							//fmt.Println("22222")
							select {
							case data := <-siaConn.MsgReadChan:
								//有数据要写给客户端
								//fmt.Println("33333")
								logsBox.Append("Received:" + string(data) + "\n")
								//fmt.Println("55555")

							case <-siaConn.ExitChan:
								//代表Reader已经退出，此时Writer也要退出
								//fmt.Println("66666")
								return

							}
						}

					}

					//fmt.Println("77777")
				}(logsBox)

			}

		}

	})

	button4Connect.Importance = widget.HighImportance
	button4Connect.SetIcon(icon.LoadButtonIcon())

	iputIpAndPort := container.NewGridWithRows(1, widget.NewLabel("HCP IP"), inputIP, layout.NewSpacer(), widget.NewLabel("Port"), inputPort, layout.NewSpacer(), button4Connect)

	//id
	idSelect := widget.NewSelect([]string{"SIA-DCS", "ADM-CID"}, func(value string) {
		//fmt.Println("type:", value)
		if value == "ADM-CID" {
			err := errors.New("Not support ADM-CID")
			dialog.ShowError(err, w)
			return
		}
	})

	//default value
	DefaultCheckBox := widget.NewCheck("Use test value", func(b bool) {
		if b {
			if inputRrcvr != nil {
				inputRrcvr.SetText("579BD")
			}
			if inputLPref != nil {
				inputLPref.SetText("E1DF0")
			}
			if inputAcct != nil {
				inputAcct.SetText("Test123")
			}
			if inputData != nil {
				inputData.SetText("NF1234/NPAZone123")
			}
		} else {

			if inputRrcvr != nil {
				inputRrcvr.SetText("")
			}
			if inputLPref != nil {
				inputLPref.SetText("")
			}
			if inputAcct != nil {
				inputAcct.SetText("")
			}
			if inputData != nil {
				inputData.SetText("")
			}
		}

	})

	//default value of id
	idSelect.SetSelected("SIA-DCS")

	idBox := container.NewHBox(widget.NewLabel("id"), idSelect, DefaultCheckBox)

	//Rrcvr
	inputRrcvr = widget.NewEntry()
	inputRrcvr.SetPlaceHolder("1-6 HEX ASCII digits like 579BD")

	//LPref
	inputLPref = widget.NewEntry()
	inputLPref.SetPlaceHolder("1-6 HEX ASCII digits like E1DF0")

	//acct
	inputAcct = widget.NewEntry()
	inputAcct.SetPlaceHolder("3-16 ASCII characters like Test123")

	//data
	inputData = widget.NewEntry()
	inputData.SetPlaceHolder("SIA data like NF1234/NPAZone123")

	siaParamBox := container.NewVBox(idBox, container.NewGridWithColumns(4, widget.NewLabel("Rrcvr"), inputRrcvr, widget.NewLabel("LPref"), inputLPref, widget.NewLabel("acct"), inputAcct, widget.NewLabel("data"), inputData))
	//siaParamBox := container.NewGridWithRows(2, idBox, container.NewGridWithColumns(2, container.NewGridWithRows(1, widget.NewLabel("Rrcvr"), inputRrcvr), container.NewHBox(widget.NewLabel("LPref"), inputLPref), container.NewHBox(widget.NewLabel("acct"), inputAcct), container.NewHBox(widget.NewLabel("data"), inputData)))

	//card to hold SIA params.
	siaDataCard := widget.NewCard("", "SIA params", siaParamBox)

	/*
		//response
		rspDataStyle := widget.TextSegment{
			Style: widget.RichTextStyleBlockquote,
			Text:  siaData,
		}
		rspDataTxt := widget.NewRichText(&rspDataStyle)

		//rspDataTxt.

		//scroll
		scroll4RspDataTxt := container.NewScroll(rspDataTxt)
	*/

	/*
		var componentsList = []string{"test1: test1"}

		scroll4RspDataTxt := widget.NewList(
			func() int { return len(componentsList) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(componentsList[i]) },
		)

		//create a go currency to update
		go func(list *[]string, listWidget *widget.List) {

			var count int
			for {
				*list = append(*list, fmt.Sprintf("current count:%d", count))
				listWidget.Refresh()
				listWidget.Resize()
				time.Sleep(time.Duration(1) * time.Second)
				fmt.Println("current count:", count)
			}
		}(&componentsList, scroll4RspDataTxt)
	*/

	//
	logsBox = widget.NewMultiLineEntry()
	logsBox.Wrapping = fyne.TextTruncate
	logsBox.SetMinRowsVisible(20)
	logsBox.SetPlaceHolder("Waiting for logs...")
	//logsBox.Disable()
	logsBox.OnChanged = func(newMsg string) {

	}

	//focusedItem := logsW

	//card for response data
	rspDatacard := widget.NewCard("", "Response data", container.NewGridWithColumns(1, logsBox))

	//button for sending data
	button4Data := widget.NewButton("Send", func() {
		//send data by the sia connection
		if siaConn == nil || siaConn.GetConnStatus() {
			err := errors.New("Please connect HCP first!")
			dialog.ShowError(err, w)
			return
		}

		//Generate SIA data
		var sia *SIA.Sia = &SIA.Sia{ID: "SIA-DCS", Rrcvr: inputRrcvr.Text, Lpref: inputLPref.Text, Acct: inputAcct.Text, Data: inputData.Text, ISeq: iSeq}
		siaData, err := sia.GenSiaData()
		if err != nil {
			errTmp := errors.New(fmt.Sprintf("connect error:%s", err.Error()))
			dialog.ShowError(errTmp, w)
			return
		}

		err = siaConn.SendMsg(siaData)
		if err != nil {
			errTmp := errors.New(fmt.Sprintf("send data error:%s", err.Error()))
			dialog.ShowError(errTmp, w)
		} else {
			if logsBox != nil {
				logsBox.Append("Send:" + string(siaData) + "\n")
			}
		}

	})

	button4Data.Importance = widget.HighImportance
	//button4Data.SetIcon(icon.LoadButtonIcon())

	w.SetContent(container.NewVBox(typeBox, iputIpAndPort, siaDataCard, rspDatacard, layout.NewSpacer(), container.NewGridWithRows(1, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), button4Data)))
	//w.SetContent(container.NewVBox(typeBox, iputIpAndPort, siaDataCard, container.NewGridWithColumns(1, rspDatacard), layout.NewSpacer(), container.NewGridWithRows(1, layout.NewSpacer(), layout.NewSpacer(), button4Data)))
	//w.SetContent(container.NewVBox(typeBox, iputIpAndPort, siaDataCard, rspDatacard, container.NewGridWithRows(1, layout.NewSpacer(), layout.NewSpacer(), button4Data)))

	w.SetIcon(icon.LoadWindowIcon())

	w.ShowAndRun()
}
