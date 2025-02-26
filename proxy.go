package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func ProtocolOptions() []string {
	return []string{
		"HTTP", "HTTPS", "SOCKS5",
	}
}

func ProtocolIndex(protocol string) int {
	switch protocol {
	case "HTTP":
		return 0
	case "HTTPS":
		return 1
	case "SOCKS5":
		return 2
	default:
		return 0
	}
}

func ProxyConfigTest(url string, config ProxyConfig) (time.Duration, error) {
	now := time.Now()

	// if !engin.IsConnect(item.Address, 5) {
	// 	return 0, fmt.Errorf("remote address connnect %s fail", item.Address)
	// }

	// urls, err := url.Parse(testhttps)
	// if err != nil {
	// 	logs.Error("%s raw url parse fail, %s", testhttps, err.Error())
	// 	return 0, err
	// }

	// var auth *engin.AuthInfo
	// if item.Auth {
	// 	auth = &engin.AuthInfo{
	// 		User:  item.User,
	// 		Token: item.Password,
	// 	}
	// }

	// var tls bool
	// if strings.ToLower(item.Protocol) == "https" {
	// 	tls = true
	// }

	// forward, err := engin.NewHttpProxyForward(item.Address, 5, auth, tls, "", "")
	// if err != nil {
	// 	logs.Error("new remote http proxy fail, %s", err.Error())
	// 	return 0, err
	// }

	// defer forward.Close()

	// request, err := http.NewRequest("GET", testhttps, nil)
	// if err != nil {
	// 	logs.Error("%s raw url parse fail, %s", testhttps, err.Error())
	// 	return 0, err
	// }

	// if strings.ToLower(urls.Scheme) == "https" {
	// 	_, err = forward.Https(engin.Address(urls), request)
	// } else {
	// 	_, err = forward.Http(request)
	// }

	// if err != nil && err.Error() != "EOF" {
	// 	logs.Error("remote server %s forward %s fail, %s",
	// 		item.Address, urls.RawPath, err.Error())
	// 	return 0, err
	// }

	return time.Since(now), nil
}

func DialogProxyConfig(parent walk.Form) {
	var dialog *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var protocol *walk.ComboBox
	var enable, auth *walk.CheckBox
	var port *walk.NumberEdit
	var user, passwd, address, testurl *walk.LineEdit
	var testButton *walk.PushButton

	proxy := ConfigGet().Proxy

	_, err := Dialog{
		AssignTo:      &dialog,
		Title:         "Proxy Config",
		Icon:          walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Size:          Size{Width: 250, Height: 300},
		MinSize:       Size{Width: 250, Height: 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Using Proxy: ",
					},
					CheckBox{
						AssignTo: &enable,
						Checked:  proxy.Enable,
					},

					Label{
						Text: "Address: ",
					},
					LineEdit{
						AssignTo: &address,
						Text:     proxy.Address,
					},

					Label{
						Text: "Port: ",
					},
					NumberEdit{
						AssignTo: &port,
						Value:    float64(proxy.Port),
						MaxValue: 65535,
						MinValue: 0,
					},

					Label{
						Text: "Protocol: ",
					},
					ComboBox{
						AssignTo:     &protocol,
						Model:        ProtocolOptions(),
						CurrentIndex: ProtocolIndex(proxy.Protocol),
						OnBoundsChanged: func() {
							protocol.SetCurrentIndex(ProtocolIndex(proxy.Protocol))
						},
					},

					Label{
						Text: "Auth: ",
					},
					CheckBox{
						AssignTo: &auth,
						Checked:  proxy.Auth,
						OnClicked: func() {
							user.SetEnabled(auth.Checked())
							passwd.SetEnabled(auth.Checked())
						},
					},

					Label{
						Text: "Username: ",
					},
					LineEdit{
						AssignTo: &user,
						Text:     proxy.User,
						Enabled:  proxy.Auth,
					},

					Label{
						Text: "Password: ",
					},
					LineEdit{
						AssignTo: &passwd,
						Text:     proxy.Password,
						Enabled:  proxy.Auth,
					},

					PushButton{
						AssignTo: &testButton,
						Text:     "Testing",
						OnClicked: func() {
							go func() {
								testButton.SetEnabled(false)
								delay, err := ProxyConfigTest(testurl.Text(), proxy)
								if err != nil {
									ErrorBoxAction(dialog, err.Error())
								} else {
									info := fmt.Sprintf("%s, %s %dms",
										"Test Pass",
										"Delay", delay/time.Millisecond)
									InfoBoxAction(dialog, info)
								}
								testButton.SetEnabled(true)
							}()
						},
					},
					LineEdit{
						AssignTo: &testurl,
						Text:     "https://www.google.com",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "Accept",
						OnClicked: func() {
							if auth.Checked() {
								if user.Text() == "" || passwd.Text() == "" {
									ErrorBoxAction(dialog, "Please input username and passwd")
									return
								}
							}

							if address.Text() == "" {
								ErrorBoxAction(dialog, "Please input address")
								return
							}

							if int(port.Value()) == 0 {
								ErrorBoxAction(dialog, "Please input port")
								return
							}

							proxy.Enable = enable.Checked()
							proxy.Address = address.Text()
							proxy.Port = int(port.Value())
							proxy.Protocol = protocol.Text()
							proxy.Auth = auth.Checked()
							proxy.User = user.Text()
							proxy.Password = passwd.Text()

							if err := ProxyConfigSave(proxy); err != nil {
								ErrorBoxAction(dialog, "Save proxy config fail, "+err.Error())
								return
							}

							dialog.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dialog.Cancel()
						},
					},
				},
			},
		},
	}.Run(parent)

	if err != nil {
		logs.Error("run proxy config dialog fail, %s", err.Error())
	}
}
