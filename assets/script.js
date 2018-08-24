/*
Copyright (c) 2018 Ringo Hoffmann (zekro Development)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

function WsClient(url) {
    this.ws = new WebSocket(url);

    this.eventListener = {};

    this.on = (event, cb) => this.eventListener[event] = cb;

    this.emit = (event, data) => {
        let rawData = JSON.stringify(event, data);
        this.ws.send(rawData);
    }

    this.ws.onmessage = (response) => {
        try {
            let data = JSON.parse(response.data);
            if (data) {
                let cb = this.eventListener[data.event]
                if (cb)
                    cb(data.data)
            }
        } catch (e) {
            console.log(e)
        }
    }
}

var ws = new WsClient('ws://zekro.de:7778/ws');

// ---------------------------------

var $ = (query) => document.querySelector(query);

var tb_message    = $('#tb_message');
var tb_name       = $('#tb_name');
var btn_send      = $('#btn_send');
var div_responses = $('#div_responses');
var f_input       = $('#f_input');
var div_login     = $('#login');
var lb_connected  = $('#lb_connected_counter');

var lastMsgUsername = "";
var userNameColor = getRandomColor();

// ---------------------------------

ws.on('message', (data) => {
    let div = document.createElement('div');
    div.className = "message_tile";
    if (lastMsgUsername != data.username) {
        let uname = document.createElement('p');
        uname.innerText = data.username;
        uname.className = "username";
        uname.style.color = userNameColor;
        div.appendChild(uname);
    }
    let message = document.createElement('p');
    message.innerText = data.message;
    message.className = "message";
    if (data.message.includes('@' + data.username)) {
        message.className += " highlighted";
    }
    div.appendChild(message);
    div_responses.appendChild(div);
    lastMsgUsername = data.username;
});

ws.on('connected', (data) => {
    let elem = document.createElement('p');
    elem.innerText = `[${data.name} CONNECTED]`;
    elem.className = "message_tile status_msg";
    div_responses.appendChild(elem);
    if (data.name == tb_name.value) {
        div_login.style.top = `-${window.innerHeight}px`;
        setTimeout(() => {
            div_login.style.display = "none";
        }, 750);
        tb_name.autofocus = false;
        tb_message.autofocus = true;
        tb_message.focus();
        lb_connected.innerText = data.nclients;
    }
});

ws.on('connect_rejected', () => {
    animateRejected();
});

ws.on('disconnected', (data) => {
    let elem = document.createElement('p');
    elem.innerText = `[${data.name} DISCONNECTED]`;
    elem.className = "message_tile status_msg";
    div_responses.appendChild(elem);
    lb_connected.innerText = data.nclients;
})

f_name.onsubmit = (e) => {
    e.preventDefault();
    if (tb_name.value.length < 1) {
        animateRejected();
        return;
    }
    ws.emit({
        event: 'username',
        data:  tb_name.value
    });
}

f_input.onsubmit = (e) => {
    e.preventDefault();
    if (tb_message.value.length < 1)
        return;
    ws.emit({
        event: 'message',
        data:  tb_message.value
    });
    tb_message.value = "";
}

// ---------------------------------

function animateRejected() {
    tb_name.style.animation = "shake .15s ease-in-out";
    setTimeout(() => tb_name.style.animation = "", 150);
}

function getRandomColor() {
    let colors = ["#ffebee","#ffcdd2","#ef9a9a","#e57373","#ef5350","#f44336","#e53935","#d32f2f","#c62828","#b71c1c","#ff8a80","#ff5252","#ff1744","#d50000","#fce4ec","#f8bbd0","#f48fb1","#f06292","#ec407a","#e91e63","#d81b60","#c2185b","#ad1457","#880e4f","#ff80ab","#ff4081","#f50057","#c51162","#f3e5f5","#e1bee7","#ce93d8","#ba68c8","#ab47bc","#9c27b0","#8e24aa","#7b1fa2","#6a1b9a","#4a148c","#ea80fc","#e040fb","#d500f9","#aa00ff","#ede7f6","#d1c4e9","#b39ddb","#9575cd","#7e57c2","#673ab7","#5e35b1","#512da8","#4527a0","#311b92","#b388ff","#7c4dff","#651fff","#6200ea","#e8eaf6","#c5cae9","#9fa8da","#7986cb","#5c6bc0","#3f51b5","#3949ab","#303f9f","#283593","#1a237e","#8c9eff","#536dfe","#3d5afe","#304ffe","#e3f2fd","#bbdefb","#90caf9","#64b5f6","#42a5f5","#2196f3","#1e88e5","#1976d2","#1565c0","#0d47a1","#82b1ff","#448aff","#2979ff","#2962ff","#e1f5fe","#b3e5fc","#81d4fa","#4fc3f7","#29b6f6","#03a9f4","#039be5","#0288d1","#0277bd","#01579b","#80d8ff","#40c4ff","#00b0ff","#0091ea","#e0f7fa","#b2ebf2","#80deea","#4dd0e1","#26c6da","#00bcd4","#00acc1","#0097a7","#00838f","#006064","#84ffff","#18ffff","#00e5ff","#00b8d4","#e0f2f1","#b2dfdb","#80cbc4","#4db6ac","#26a69a","#009688","#00897b","#00796b","#00695c","#004d40","#a7ffeb","#64ffda","#1de9b6","#00bfa5","#e8f5e9","#c8e6c9","#a5d6a7","#81c784","#66bb6a","#4caf50","#43a047","#388e3c","#2e7d32","#1b5e20","#b9f6ca","#69f0ae","#00e676","#00c853","#f1f8e9","#dcedc8","#c5e1a5","#aed581","#9ccc65","#8bc34a","#7cb342","#689f38","#558b2f","#33691e","#ccff90","#b2ff59","#76ff03","#64dd17","#f9fbe7","#f0f4c3","#e6ee9c","#dce775","#d4e157","#cddc39","#c0ca33","#afb42b","#9e9d24","#827717","#f4ff81","#eeff41","#c6ff00","#aeea00","#fffde7","#fff9c4","#fff59d","#fff176","#ffee58","#ffeb3b","#fdd835","#fbc02d","#f9a825","#f57f17","#ffff8d","#ffff00","#ffea00","#ffd600","#fff8e1","#ffecb3","#ffe082","#ffd54f","#ffca28","#ffc107","#ffb300","#ffa000","#ff8f00","#ff6f00","#ffe57f","#ffd740","#ffc400","#ffab00","#fff3e0","#ffe0b2","#ffcc80","#ffb74d","#ffa726","#ff9800","#fb8c00","#f57c00","#ef6c00","#e65100","#ffd180","#ffab40","#ff9100","#ff6d00","#fbe9e7","#ffccbc","#ffab91","#ff8a65","#ff7043","#ff5722","#f4511e","#e64a19","#d84315","#bf360c","#ff9e80","#ff6e40","#ff3d00","#dd2c00","#efebe9","#d7ccc8","#bcaaa4","#a1887f","#8d6e63","#795548","#6d4c41","#5d4037","#4e342e","#3e2723","#fafafa","#f5f5f5","#eeeeee","#e0e0e0","#bdbdbd","#9e9e9e","#757575","#616161","#424242","#212121","#eceff1","#cfd8dc","#b0bec5","#90a4ae","#78909c","#607d8b","#546e7a","#455a64","#37474f","#263238","rgba(0, 0, 0, 0.87)","rgba(0, 0, 0, 0.54)","rgba(0, 0, 0, 0.38)","rgba(0, 0, 0, 0.12)","rgba(0, 0, 0, 0.54)","rgba(0, 0, 0, 0.38)","rgba(255, 255, 255, 1)","rgba(255, 255, 255, 0.7)","rgba(255, 255, 255, 0.5)","rgba(255, 255, 255, 0.12)","rgba(255, 255, 255, 1)","rgba(255, 255, 255, 0.5)","#ffffff","#000000"];
    return colors[Math.floor(Math.random() * (colors.length - 1))];
}
  