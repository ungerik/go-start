var gostart_media = gostart_media || {};

gostart_media.fillChooseDialog = function(thumbnailsSelector, thumbnailsURL, onClickFunc) {
	var thumbnails = jQuery(thumbnailsSelector);
	thumbnails.empty();
	jQuery.ajax({url: thumbnailsURL, dataType: "json"})
	.fail(function(jqXHR, textStatus) {
		alert("Request failed: " + textStatus);
	})
	.done(function(data) {
		jQuery.each(data, function(index, value) {
			var img = jQuery('<img src="'+value.url+'" alt="'+value.title+'" style="cursor:pointer;margin:5px;"/>');
			img.click(function(){ onClickFunc(value); });
			thumbnails.append(img);
		});
	});
};

