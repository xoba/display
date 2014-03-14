var commas = function(e){e=e.toString();var t=/(-?\d+)(\d{3})/;while(t.test(e))e=e.replace(t,"$1,$2");return e};

var display = function(sel,h,cb) {
    if (h) {
	$(sel).html(h.join("\n"));
    }
    if (cb) {
	$.each(cb,function(_,f) {
	    display(null,null,f())
	});
    }
};

var esc = function(x) {
    if (x && typeof x === 'string') {
	x = x.replace(/&/g,"&amp;");
	x = x.replace(/</g,"&lt;");
	x = x.replace(/ /g,"&nbsp;");
    }
    return x;
}

var genid = (function() {
    var uuid = (function() {
        var x = function() {
	    return (((1+Math.random())*0x10000)|0).toString(16).substring(1);
        };
        return function() {
	    return (x()+x()+"-"+x()+"-"+x()+"-"+x()+"-"+x()+x()+x());
        };
    })(); 
    return function() {
        return "I" + uuid().substring(0,8);
    };
})();

var color = function(x,c) {
    return "<span style='color:#"+c+";'>"+x+"</span>";
}

var cpxi = "CPX" + color('i','5EA0CC') + color(".","F2CB31") + color('ch','D52B1E');

