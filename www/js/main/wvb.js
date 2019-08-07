
function get_remaining_height(e) {
	return $(window).height() - e.offset().top;
}

function set_max_height(e, max) {
	e.css("max-height", max + "px");
}

var elem;

$(document).ready(function() {
	elem = $(".wvb-fill-remaining");
	set_max_height(elem, get_remaining_height(elem));
	$(window).resize(function() {
		elem = $(".wvb-fill-remaining");
		set_max_height(elem, get_remaining_height(elem));
	});
});
