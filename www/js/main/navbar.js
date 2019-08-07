$(document).ready(function() {
	$(".nav-item").each(function() {
		$(this).toggleClass("active", jQuery("a", this).attr("href") == window.location.pathname);
	});
});

$(document).ready(function() {
	var spacer = $(".wvb-spacer");
	var navbar = $("#wvb-navbar");
	spacer.height(navbar.outerHeight(true));
	$(window).resize(function() {
		spacer.height(navbar.outerHeight(true));
	});
});
