const ActionError = "00",
    ActionPing = "01",
    ActionGameSetup = "02",
    ActionChat = "03",
    ActionPlayerStatus = "04",
    ActionMove = "05",
    ActionFire = "06",
    ActionSplit = "07",
    ActionLeaderBoard = "08";

const global = {
    debug: false,
    gameWidth: 0,
    gameHeight: 0,
    screenWidth: window.innerWidth,
    screenHeight: window.innerHeight,
    lineColor: '#000000',
    backgroundColor: '#f2fbff',
    virusColor: '#7bff66',
    ripWaitSeconds: 3,
};

const client = {
    targetX: 0,
    targetY: 0,
    player: undefined,
    ws: undefined,
    pingTime: undefined,
    ping: undefined,
    animLoopHandle: undefined,
    gameLoopCount: 0,
};

const canvas = document.getElementById("game"),
    graph = canvas.getContext("2d");

const controller = {
    init() {
        document.getElementById('startBtn').onclick = controller.start;
        window.onkeypress = controller.onKeypress;
        canvas.addEventListener("mousemove", controller.onMouseMove);
        canvas.addEventListener("mouseout", controller.onMouseOut);

        window.requestAnimFrame = (function () {
            return window.requestAnimationFrame ||
                window.webkitRequestAnimationFrame ||
                window.mozRequestAnimationFrame ||
                window.msRequestAnimationFrame ||
                window.oRequestAnimationFrame ||
                window.msRequestAnimationFrame ||
                function (callback) {
                    window.setTimeout(callback, 1000 / 60);
                };
        })();

        window.cancelAnimFrame = (function (handle) {
            return window.cancelAnimationFrame ||
                window.mozCancelAnimationFrame;
        })();
    },
    start() {
        let name = document.getElementById('playerNameInput').value;
        if (name === '') {
            return
        }
        let protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
        let ws = new WebSocket(`${protocol}://${window.location.host}/game?name=${name}`);
        ws.onopen = evt => {
            document.getElementById('gameAreaWrapper').style.opacity = 1;
            document.getElementById('startMenuWrapper').style.maxHeight = '0px';

            client.ws = ws;
            controller.gameLoop();
        };
        ws.onmessage = evt => handler.handle(evt.data);
        ws.onclose = evt => {
            client.ws = undefined;
            client.player = undefined;
            drawer.drawBackground();
            drawer.drawRIP();
            if (client.animLoopHandle) {
                window.cancelAnimationFrame(client.animLoopHandle);
                client.animLoopHandle = undefined;
            }
            window.setTimeout(function () {
                document.getElementById('gameAreaWrapper').style.opacity = 0;
                document.getElementById('startMenuWrapper').style.maxHeight = '1000px';
                messager.clear();
            }, global.ripWaitSeconds * 1000);
        };
    },
    gameLoop() {
        if (!client.ws) {
            return
        }
        if (client.gameLoopCount > 60) {
            client.pingTime = (new Date()).valueOf();
            sender.ping();
            client.gameLoopCount = 0;
        } else {
            client.gameLoopCount++;
        }
        client.animLoopHandle = window.requestAnimFrame(controller.gameLoop);
        drawer.drawBackground();
        let player = client.player;
        if (player) {
            drawer.drawPlayer();
            sender.send(ActionMove, client.targetX + ',' + client.targetY)
        }
        if (global.debug) {
            drawer.drawDebugInfo();
        }
    },
    onKeypress(event) {
        let key = event.which || event.keyCode;
        switch (key) {
            case 88:
            case 122:
                sender.send(ActionFire);
                break;
            case 90:
            case 120:
                sender.send(ActionSplit);
                break;
        }
    },
    onMouseMove(event) {
        let cRect = canvas.getBoundingClientRect();
        let x = Math.round(event.clientX - cRect.left);
        let y = Math.round(event.clientY - cRect.top);
        client.targetX = x - global.screenWidth / 2;
        client.targetY = y - global.screenHeight / 2;
    },
    onMouseOut(event) {
        client.targetX = 0;
        client.targetY = 0;
    },
};

const messager = {
    chatInput: document.getElementById('chatInput'),
    chatList: document.getElementById('chatList'),
    init() {
        let input = messager.chatInput;
        input.addEventListener('keypress', messager.onKeypress);
        input.addEventListener('keyup', messager.onKeyup);
    },
    onKeypress(event) {
        let input = messager.chatInput;
        let key = event.which || event.keyCode;
        if (key === 13 || key === 108) {
            let text = input.value.replace(/(<([^>]+)>)/ig, '');
            if (text !== '') {
                sender.send(ActionChat, text);
                input.value = '';
                canvas.focus();
            }
        }
    },
    onKeyup(event) {
        let input = messager.chatInput;
        let key = event.which || event.keyCode;
        if (key === 27) {
            input.value = '';
            canvas.focus();
        }
    },
    append(message, isSystem) {
        let line = document.createElement('li');
        if (isSystem) {
            line.className = 'player';
        } else {
            line.className = 'system';
        }
        line.innerHTML = message;
        let chatList = messager.chatList;
        if (chatList.childNodes.length > 10) {
            chatList.removeChild(chatList.childNodes[0]);
        }
        chatList.appendChild(line);
    },
    clear() {
        messager.chatList.innerHTML = ''
    },
};

const drawer = {
    drawBackground() {
        graph.fillStyle = global.backgroundColor;
        graph.fillRect(0, 0, global.screenWidth, global.screenHeight);
    },
    drawDebugInfo() {
        let infos = [];
        infos.push(['targetX:', client.targetX,
            'targetY:', client.targetY,
            'connected:', !!client.ws,
            'ping', client.ping ? client.ping : 0].join(' '));
        graph.fillStyle = '#000000';
        let player = client.player;
        if (player) {
            graph.fillStyle = '#000000';
            infos.push(['x:', player.x,
                'y:', player.y,
                'mass:', player.massTotal].join(' '));
            infos.push(['cell:', player.visibleCells.length,
                'food:', player.visibleFoods.length,
                'mass food:', player.visibleMassFoods.length,
                'virus:', player.visibleViruses.length].join(' '));
        }
        for (let i = 0; i < infos.length; i++) {
            graph.fillText(infos[i], 0, 10 + i * 20);
        }
    },
    drawCircle(x, y, radius) {
        graph.beginPath();
        graph.arc(x, y, radius, 0, 2 * Math.PI, false);
        graph.closePath();
        graph.stroke();
        graph.fill();
    },
    drawPlayer() {
        let player = client.player;
        let font = graph.font;
        player.visibleCells.forEach(item => {
            let x = item.x - player.x + global.screenWidth / 2;
            let y = item.y - player.y + global.screenHeight / 2;
            let r = item.radius;
            graph.fillStyle = item.background;
            drawer.drawCircle(x, y, r);
            graph.fillStyle = item.textColor;
            let fontSize = Math.max(r / 3, 12);
            graph.font = 'bold ' + fontSize + 'px sans-serif';
            graph.fillText(item.name, x - fontSize * item.name.length / 2, y + 3);
        });
        graph.font = font;

        player.visibleFoods.forEach(item => {
            graph.fillStyle = item.color;
            drawer.drawCircle(item.x - player.x + global.screenWidth / 2,
                item.y - player.y + global.screenHeight / 2,
                item.radius);
        });

        player.visibleMassFoods.forEach(item => {
            graph.fillStyle = item.color;
            drawer.drawCircle(item.x - player.x + global.screenWidth / 2,
                item.y - player.y + global.screenHeight / 2,
                item.radius);
        });

        graph.fillStyle = global.virusColor;
        player.visibleViruses.forEach(item => {
            drawer.drawCircle(item.x - player.x + global.screenWidth / 2,
                item.y - player.y + global.screenHeight / 2,
                item.radius);
        });
    },
    drawRIP() {
        let font = graph.font;
        graph.fillStyle = '#ff0000';
        graph.font = 'bold 22px sans-serif';
        graph.fillText('you are eaten or lose connect!', global.screenWidth / 2 - 60, global.screenHeight / 2 - 15);
        graph.fillText('game will exit after ' + global.ripWaitSeconds + ' seconds ...', global.screenWidth / 2 - 120, global.screenHeight / 2 + 15);
        graph.font = font;
    },
};

const handler = {
    handle(message) {
        let msgType = message.substring(0, 2);
        let payload = message.substring(3);
        switch (msgType) {
            case ActionError:
                handler.handleError(payload);
                break;
            case ActionPing:
                handler.handlePing(payload);
                break;
            case ActionGameSetup:
                handler.handleGameSetup(payload);
                break;
            case ActionChat:
                handler.handleChat(payload);
                break;
            case ActionPlayerStatus:
                handler.handlePlayerStatus(payload);
                break;
            case ActionLeaderBoard:
                handler.handleLeaderBoard(payload);
                break;
        }
    },
    handlePing(data) {
        if (client.pingTime) {
            client.ping = (new Date()).valueOf() - client.pingTime;
            client.pingTime = undefined;
        }
    },
    handleGameSetup(data) {
        let split = data.split("|");
        global.gameWidth = parseFloat(split[0]);
        global.gameHeight = parseFloat(split[1]);
        global.screenWidth = parseFloat(split[2]);
        global.screenHeight = parseFloat(split[3]);
        global.virusColor = split[4];
        canvas.setAttribute('width', global.screenWidth);
        canvas.setAttribute('height', global.screenHeight);
    },
    handleChat(data) {
        let type = data.substring(0, 1);
        if ('0' === type) {
            messager.append(data.substring(1), true)
        } else {
            messager.append(data.substring(1), false)
        }
    },
    handleError(data) {

    },
    handlePlayerStatus(data) {
        let parsePlayer = data => {
            let datas = data.split('|');
            let index = 0;
            let pDatas = datas[index].split(',');
            let player = {
                name: pDatas[0],
                x: parseFloat(pDatas[1]),
                y: parseFloat(pDatas[2]),
                massTotal: parseFloat(pDatas[3]),
                visibleCells: [],
                visibleFoods: [],
                visibleMassFoods: [],
                visibleViruses: [],
            };
            index++;

            let cLen = parseInt(datas[index]);
            index++;
            for (let i = index; i < index + cLen; i++) {
                let cDatas = datas[i].split(',');
                player.visibleCells.push({
                    name: cDatas[0],
                    background: cDatas[1],
                    textColor: cDatas[2],
                    x: parseFloat(cDatas[3]),
                    y: parseFloat(cDatas[4]),
                    radius: parseFloat(cDatas[5]),
                });
            }
            index += cLen;

            let fLen = parseInt(datas[index]);
            index++;
            for (let i = index; i < index + fLen; i++) {
                let fDatas = datas[i].split(',');
                player.visibleFoods.push({
                    x: parseFloat(fDatas[0]),
                    y: parseFloat(fDatas[1]),
                    radius: parseFloat(fDatas[2]),
                    color: fDatas[3],
                });
            }
            index += fLen;

            let mfLen = parseInt(datas[index]);
            index++;
            for (let i = index; i < index + mfLen; i++) {
                let mfDatas = datas[i].split(',');
                player.visibleMassFoods.push({
                    x: parseFloat(mfDatas[0]),
                    y: parseFloat(mfDatas[1]),
                    radius: parseFloat(mfDatas[2]),
                    color: mfDatas[3],
                });
            }
            index += mfLen;

            let vLen = parseInt(datas[index]);
            index++;
            for (let i = index; i < index + vLen; i++) {
                let vDatas = datas[i].split(',');
                player.visibleViruses.push({
                    x: parseFloat(vDatas[0]),
                    y: parseFloat(vDatas[1]),
                    radius: parseFloat(vDatas[2]),
                });
            }
            return player;
        };

        client.player = parsePlayer(data);
    },
    handleLeaderBoard(data) {
        let split = data.split(',');
        let status = '<span class="title">LeaderBoard</span>';
        for (let i = 0; i < split.length; i++) {
            status += '<br />';
            if (split[i] === client.player.name) {
                status += '<span class="me">' + (i + 1) + '. ' + split[i] + "</span>"
            } else {
                status += (i + 1) + '. ' + split[i];
            }
        }
        document.getElementById('status').innerHTML = status
    },
};

const sender = {
    ping() {
        sender.send(ActionPing)
    },
    send(msgType, payload) {
        let ws = client.ws;
        if (ws) {
            ws.send(msgType + "|" + payload)
        }
    }
};

function bytesToSize(bytes) {
    if (bytes === 0) return '0 B';
    let k = 1000, // or 1024
        sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
        i = Math.floor(Math.log(bytes) / Math.log(k));

    return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i];
}

controller.init();
messager.init();