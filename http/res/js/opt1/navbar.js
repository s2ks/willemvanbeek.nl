$(document).ready(function(){
	$(".nav-item").each(function() {
		$(this).toggleClass("active", jQuery("a", this).attr("href") == window.location.pathname);
	});
});
