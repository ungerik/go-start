var gostart_media = {

	upload: function(url, thumbnailFrameID) {
		var thumbnailFrame = jQuery("#"+thumbnailFrameID);
	},

	fillChooser: function(thumbnailsSelector, thumbnailsURL, onClickFunc) {
		var thumbnails = jQuery(thumbnailsSelector);
		thumbnails.empty();
		jQuery.ajax(thumbnailsURL)
		.fail(function(jqXHR, textStatus) {
			alert("Request failed: " + textStatus);
		})
		.done(function(data) {
			jQuery.each(data, function(index, value) {
				var img = jQuery('<img src="'+value.url+'" alt="'+value.title+'"/>');
				img.click(function(){ onClickFunc(value); });
				thumbnails.append(img);
			});
		});
	}

};

