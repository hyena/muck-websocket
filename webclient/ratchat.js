// Client side logic for RatChat
// @author: Strain-113
// Requires ansi_up.js and jquery.

// TODO list:
// * Make sure newlines are sane.
// * Reconnection logic.
// * Seamless reconnections.

var ratChat = function () {
    var host = "ws://0.0.0.0:8000/";

    function showMessage (msg) {
        // TODO: Handle newlines, etc.
        var historyBox = $("#chathistory");
        var atBottom = (historyBox.scrollTop() >= historyBox[0].scrollHeight - historyBox.height());
        historyBox.append(msg);
        // If we were scrolled to the bottom before this call, remain there.
        if (atBottom) {
            historyBox.scrollTop(historyBox[0].scrollHeight - historyBox.height())
        }
    }

    // Register listeners.
    $(document).ready(function () {
        // Check for websocket availability.
        if (!("WebSocket" in window)) {
            showMessage("Sorry. RatChat requires a browser with WebSockets support.");
		    return;
        }

        try{
             var socket = new WebSocket(host);

             socket.onopen = function(){
                 showMessage('<p class="text-success">Connected.</p>');
             }

             socket.onmessage = function(msg){
                 var msg = ansi_up.linkify(ansi_up.ansi_to_html(ansi_up.escape_for_html(msg.data)));

                 showMessage(msg);
             }

             socket.onclose = function(){
                 showMessage('<p class="text-danger">Disconnected.</p>');
             }

         } catch(exception){
             showMessage('<p class="text-danger">Error connecting.</p>');
         }

        $("#textbox").keyup(function (e) {
            if (e.keyCode == 13) {
                if (!e.shiftKey) {
                    var content = this.value;
                    try {
                        socket.send(content);
                    } catch(exception) {
                        showMessage('<p class="text-warning">Couldn\'t send message.</p>');
                    }
                    this.value = "";
                    e.stopPropagation();
                }
            }
        });
    });
}();
