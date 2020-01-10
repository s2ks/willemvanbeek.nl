/*
$('.grid').colcade({
	columns: '.grid-col',
	items: '.grid-item'
});
*/

var $grid = $('.grid').masonry({
	itemSelector: '.grid-item',
	columnWidth: '.grid-sizer',
	percentPosition: true
});

$grid.imagesLoaded().progress(function() {
	$grid.masonry('layout');
});
