/**
 * Created by LittleGpNator on 01/05/2017.
 */
class Messages {
    constructor() {
        this.ws = new WebSocket("ws://localhost:3001");

        this.ws.onmessage = (event) =>{
            console.log(event.data)
            var jsonm = event.data;
            var bef = $('#messages-ter').html();
            var beep = "\n<li>\n<p>Res: " + jsonm
            $('#messages-ter').html(bef + beep);
        };
        this.ws.onerror = (error) =>{
            $("#connection_label").html("Not connected");
        };
        this.ws.onopen = () =>{
            $("#connection_label").html("Connected");
        };
        this.ws.onclose = function (message) {
            $("#connection_label").html("Not connected");
        };
    }

    on_send(event) {
        var name = $('#name').val();
        var msg = $('#msg').val();
        var color = $('#color').val();
        var date = new Date();
        var obj = {'name': name, 'msg': msg, 'date': date, 'color': color};
        console.log(color);
        this.ws.send(obj.name);
    }
}

var messages;
//When this file is fully loaded, initialize board with context
$(document).ready(function () {
    messages = new Messages($('#messages-ter'));
    $('#send').click((event) => {
        console.log("hei");
    messages.on_send(event);
    });
});