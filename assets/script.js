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

Element.prototype.remove = function() {
    this.parentElement.removeChild(this);
}
NodeList.prototype.remove = HTMLCollection.prototype.remove = function() {
    for(var i = this.length - 1; i >= 0; i--) {
        if(this[i] && this[i].parentElement) {
            this[i].parentElement.removeChild(this[i]);
        }
    }
}

// ------------------------------------------------------

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
                let cb = this.eventListener[data.event];
                if (cb)
                    cb(data.data);
            }
        } catch (e) {
            console.log(e);
        }
    }
}

var ws = new WsClient(
    window.location.href.replace(/((http)|(https)):\/\//gm, "ws://") + "ws"
);

// ---------------------------------

var $ = (query) => document.querySelector(query);

var tb_message = $('#tb_message');
var tb_name = $('#tb_name');
var tb_password = $('#tb_password');
var btn_send = $('#btn_send');
var div_responses = $('#div_responses');
var f_input = $('#f_input');
var div_login = $('#login');
var div_acc_create = $('#acc_create');
var lb_connected = $('#lb_connected_counter');
var ul_users_list = $('#ul_users_list');

var lastMsgUsername = "";

var myUsername = "";

// ---------------------------------

ws.ws.onclose = () => {
    window.alert('Websocket connection closed.');
    window.location = "/";
};

ws.on('message', (data) => {
    appendMessage(data);
});

ws.on('connected', (data) => {
    let elem = document.createElement('p');
    elem.innerText = `[${data.name} CONNECTED]`;
    elem.className = 'message_tile status_msg';
    div_responses.appendChild(elem);
    updateUsersList(data.clients);
    if (data.name == myUsername) {
        if (data.history) {
            console.log(data.history)
            data.history.forEach(msg => appendMessage(msg.data));
        }
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

ws.on('messageDeleted', (data) => {
    let msg = $('#message' + data.data.id);
    let container = $('#container' + data.data.id);
    msg.innerHTML = '<p class="status_msg">deleted</p>';
    container.getElementsByClassName('deleteLink')[0].remove();
});

ws.on('disconnected', (data) => {
    let elem = document.createElement('p');
    elem.innerText = `[${data.name} DISCONNECTED]`;
    elem.className = 'message_tile status_msg';
    div_responses.appendChild(elem);
    updateUsersList(data.clients);
    lb_connected.innerText = data.nclients;
});

ws.on('usernameState', (data) => {
    div_acc_create.style.top = data ? '50px' : '0px';
    div_acc_create.style.opacity = data ? 0 : 1;

});

f_name.onsubmit = (e) => {
    e.preventDefault();
    if (tb_name.value.length < 1 || tb_password.value.length < 1) {
        animateRejected();
        return;
    }
    myUsername = tb_name.value;
    ws.emit({
        event: 'login',
        data: {
            username: myUsername,
            password: tb_password.value,
        },
    });
};

tb_name.oninput = (e) => {
    let cvalue = e.target.value;
    setTimeout(() => {
        if (tb_name.value.length < 1) {
            div_acc_create.style.top = '50px';
            div_acc_create.style.opacity = 0;
            return;
        }
        if (tb_name.value == cvalue) {
            ws.emit({
                event: 'checkUsername',
                data: cvalue
            });
        }
    }, 300);
};

f_input.onsubmit = (e) => {
    e.preventDefault();
    if (tb_message.value.length < 1)
        return;
    ws.emit({
        event: 'message',
        data: tb_message.value,
    });
    tb_message.value = '';
};

// ---------------------------------

function animateRejected() {
    tb_name.style.animation = 'shake .15s ease-in-out';
    setTimeout(() => tb_name.style.animation = '', 150);
}

function appendMessage(msgEvent) {
    var converter = new showdown.Converter({headerLevelStart: 3, strikethrough: true, emoji: true, underline: true,});
    let div = document.createElement('div');
    div.className = 'message_tile';
    div.id = 'container' + msgEvent.id;

    let divTitle = document.createElement('div');
    divTitle.className = 'head';

    if (lastMsgUsername != msgEvent.author.username) {
        let uname = document.createElement('div');
        uname.innerText = msgEvent.author.username;
        uname.className = 'username';
        uname.style.color = msgEvent.author.color;
        divTitle.appendChild(uname);

        let time = document.createElement('div');
        time.innerText = getTime(msgEvent.timestamp);
        time.className = 'time';
        divTitle.appendChild(time);

        div.appendChild(divTitle);
    }
    let message = document.createElement('div');
    message.id = 'message' + msgEvent.id;
    let converted = converter.makeHtml(msgEvent.content);
    message.innerHTML = converted;
    message.className = 'message';
    if (msgEvent.content.includes('@' + myUsername)) {
        message.className += ' highlighted';
        let ops = {
            body: message,
            icon: "https://camo.githubusercontent.com/6490d7a3b892fd102a9ef1718aec3c41838639ee/68747470733a2f2f7a656b726f2e64652f7372632f676f5f636861745f6c6f676f2e706e67"
        }
        if (!("Notification" in window)) {
            alert("This browser does not support desktop notification");
        } else if (Notification.permission === "granted") {
            var notification = new Notification("You got pinged", ops);
        } else if (Notification.permission !== "denied") {
            Notification.requestPermission(function (permission) {
                if (permission === "granted") {
                    var notification = new Notification("You got pinged", ops);
                }
            });
        }
    }

    div.appendChild(message);

    if (msgEvent.author.username == myUsername) {
        let messageActionDiv = document.createElement('div');
        let deleteLink = document.createElement('a');
        //deleteLink.href = "javasctipt:{}";
        deleteLink.onclick = () => {
            ws.emit({
                event: 'deleteMessage',
                data: msgEvent.id
            });
        };
        deleteLink.innerText = "remove";
        messageActionDiv.className = "deleteLink";
        messageActionDiv.appendChild(deleteLink);
        div.appendChild(messageActionDiv);
    }

    div_responses.appendChild(div);
    div.scrollIntoView();
    lastMsgUsername = msgEvent.author.username;
}

function updateUsersList(usersMap) {
    ul_users_list.innerHTML = '';
    Object.keys(usersMap).forEach(uname => {
        let elem = document.createElement('li');
        elem.style.color = usersMap[uname];
        elem.innerText = uname;
        ul_users_list.appendChild(elem);
    });
}

function replaceLinks(message) {
    return message.replace(/((https?:\/\/)?(www.)?\w+\.\w+)/gm, function (a) {
        if (!a.startsWith('http://') && !a.startsWith('https://'))
            a = "https://" + a;
        return `<a href="${a}" target="_blank">${a}</a>`;
    });
}

function getTime(timestamp) {
    function btf(inp) {
        if (inp < 10)
            return '0' + inp;
        return inp;
    }
    var date = new Date(timestamp * 1000),
        y = date.getFullYear(),
        m = btf(date.getMonth() + 1),
        d = btf(date.getDate()),
        h = btf(date.getHours()),
        min = btf(date.getMinutes()),
        s = btf(date.getSeconds());
    return `${d}.${m}.${y} - ${h}:${min}:${s}`;
}
