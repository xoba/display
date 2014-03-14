$(function() {

    var connection = new ReconnectingWebSocket('ws://'+location.hostname+':'+location.port+'/ws');

    var log = function(x) {
        console.log(x);
    };

    connection.onopen = function() {
        log("connected");
    };
    
    connection.onclose = function() {
        log("lost connection");
    };
    
    connection.onerror = function(error) {
	log("error",error);
    };

    connection.onmessage = function(e) {
	var img = location.origin + "/proxy?url="+encodeURIComponent(JSON.parse(e.data))
	$('html').css("background","url("+img+") no-repeat center center fixed").css("background-size","cover");
    };
    
    $('html').click(function() {
        var e = $('html')[0];
        if (e.requestFullscreen) {
            e.requestFullscreen();
        } else if (e.mozRequestFullScreen) {
            e.mozRequestFullScreen();
        } else if (e.webkitRequestFullScreen) {
            e.webkitRequestFullScreen();
        }
    });
    
})