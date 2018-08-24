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

var ws = new WsClient('ws://zekro.de:7777/ws');

// ---------------------------------

var $ = (query) => document.querySelector(query);

var tb_message    = $('#tb_message');
var tb_name       = $('#tb_name');
var btn_send      = $('#btn_send');
var div_responses = $('#div_responses');
var f_input       = $('#f_input');
var div_login     = $('#login');
var lb_connected  = $('#lb_connected_counter');

// ---------------------------------

ws.on('message', (data) => {
    let elem = document.createElement('p');
    elem.innerText = data;
    elem.className = "message_tile";
    div_responses.appendChild(elem);
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
