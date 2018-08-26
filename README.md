<div align="center">
     <img src="https://zekro.de/src/go_chat_logo.png" width="400"/>
     <h1>~ goChat ~</h1>
     <strong>Simple websocket based chat server + client written in Go</strong><br><br>
     <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/uses-html.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/uses-css.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/uses-js.svg" height="30" />&nbsp;
     <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
</div>

---

## First Things First

**Attention:** In the master branch is only the code for the default standard setup without further features. For the current features code, switch to branch [new_features](https://github.com/zekroTJA/goChat/tree/new_features).

---

## Motivation

Currently, I am trying to get into Go and gorilla/websocket, so, I created this little project to test around with a self-created event system and socket communication. Also for some fontend training, of course ;)

---

## Installation

This chat is based on a self-hostet server component and browser-driven web client.

> If you want to compile the binaries by yourself, take a look [below](#self-compiling).

Download the pre-compiled binaries and place them somewhere on your server. **Attention:** The web client assets provided by the integrated HTTP server are not included in the binaries and deed to be placed right next to the binary like in following example setup:

```
goChat/
  |-- precompiled_131_linux_x86_64
  |-- assets/
        |-- index.html
        |-- script.js
        |-- slyte.css
```

After that, you can run the server component by entering following command:
```
./gochat_server 7777
```
Of course, you need to use the name of your binary downloaded. The port the server is accessable on is set as first argument. Defaultly, if no agrument is passed, the port 7777 will be used.

---

## Self Compiling

To compile it by yourself, use following commands:

```bash
$ git clone https://github.com/zekroTJA/goChat
$ cd goChat
# if you want to use the latest features:
$ git checkout new_features
# of course, GOPATH needs to be set before
$ go get github.com/gorilla/websocket
$ go build -o gochat_server
```

---

## Screenshots

<img src="http://zekro.de/ss/firefox_2018-08-25_19-44-50.png" width="880"/>
<img src="http://zekro.de/ss/firefox_2018-08-25_19-56-15.png" width="880"/>

---

## Whished Features & Ideas

- [ ] Emojis + Emoji keywords
- [ ] Multiple Chat Rooms
- [ ] Logo as Favicon
- [ ] Deletable Messages
- [ ] Editable Messages
- [ ] Chat Administration
- [ ] SSL Certification
- [ ] Light Design Mode *haha, nice joke :^)* 
- [x] Connected Users List `(1.4)`
- [x] Link Highlighting `(1.4)`
- [x] Mentions `(1.3)`
- [x] Accounts and Login `(1.3)`

---

## 3rd Party Dependencies & Credits

- [gorilla/websocket](https://github.com/gorilla/websocket)

Chat Client Font:
- [Source Code Pro](https://fonts.google.com/specimen/Source+Code+Pro) *by Paul D. Hunt*

Logo Font:
- [Montserrat](https://fonts.google.com/specimen/Montserrat) *by Julieta Ulanovsky, Sol Matas, Juan Pablo del Peral & Jacques Le Bailly*

---

© 2018 zekro Development  

[zekro.de](https://zekro.de) | contact[at]zekro.de


