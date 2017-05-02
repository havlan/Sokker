/**
 * Created by LittleGpNator on 01/05/2017.
 */
var id=0;

class Messages {
    constructor() {
        this.ws = new WebSocket("ws://localhost:3001");

        this.ws.onmessage = (event) =>{
            console.log(event.data);
            var jsonm = JSON.parse(event.data);
            var side = "left";
            var colorA = "#00c5c9";
            if (jsonm.id == id){
                side = "right";
                colorA = "#bbc8c9";
            }

            $("<div class='messageBoxOver'><div class='messageBox' style='float: "+side+"; background-color: "+colorA+"'><div class='navn'>"+jsonm.name+"</div><div class='messageBoxText'>"+jsonm.msg+"</div></div></div>").appendTo($(".scroll-area"));
            $("#messages-div").scrollTop($("#messages-div")[0].scrollHeight);


        };
        this.ws.onerror = (error) =>{
            $("footer").style.backgroundColor ='red';
            $("#connection_label").html("Not connected");
        };
        this.ws.onopen = () =>{
            $("#connection_label").html("Connected");
            id=Math.floor(Math.random() * (99999 - 4 + 1)) + 4;
        };
        this.ws.onclose = function (message) {
            $("#connection_label").html("Not connected");
        };
    }

    on_send(event) {
        var name = $('#name').val();
        var msg = $('#msg').val();
        var ide  = id;
        var data = {'name':name,'msg':msg,'id':ide};
        var strinfy = JSON.stringify(data);
        if (strinfy.length>130){
            strinfy="";
        }
        this.ws.send(strinfy);
    }
}

var messages;
//When this file is fully loaded, initialize board with context
$(document).ready(function () {
    messages = new Messages($('#messages-ter'));
    $('#send').click((event) => {
        messages.on_send(event);
    });
});