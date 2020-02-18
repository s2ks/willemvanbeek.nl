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

$('#img-modal').modal({
	backdrop: false
});

function displayImgModal(src) {
	var m = $('#img-modal');
	var i = $('#img-modal-img');

	i.attr("src", src);
	m.modal('show');
}
