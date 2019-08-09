$(document).ready(function() {
	$(".nav-item").each(function() {
		$(this).toggleClass("active", jQuery("a", this).attr("href") == window.location.pathname);
	});

	var spacer = $(".wvb-spacer");
	var navbar = $("#wvb-navbar");
	spacer.attr("style", "display: block !important;");
	spacer.height(navbar.outerHeight(true));
	$(window).resize(function() {
		spacer.height(navbar.outerHeight(true));
	});

	$(".wvb-toggle").click(function() {
		//$(".navbar-collapse").collapse("hide");

	});

	var scroll = $(".wvb-scroll");
	var target = $(scroll.attr("href").replace('/', ''));
	scroll.click(function() {
		$("#wvb-root").animate({
			scrollTop: target.position().top
		}, 1000);
	});
});

