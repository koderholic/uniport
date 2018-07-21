var m = require("mithril")
import {appAlert} from './#utils.js';

export validateImage(fileId, imageSrc) {

	var imageInputFile = document.getElementById('#' + imageID + 'File').get(0).files;

	if (imageInputFile.length == 0) {
		appAlert([{ type: 'bg-red', message: "Please Select File(s) to Upload", }])
		return;
	}

	if (imageInputFile.length > 1) {
		appAlert([{ type: 'bg-red', message: "Exceeded Maximum Allowed Files", }])
		return;
	}

	for (var counter = 0; counter < imageInputFile.length; counter++) {
		var reader = new FileReader();
		reader.onload = function(file_name) {
			return function(event) {
				var the_url = event.target.result;

				if (the_url.indexOf("base64") == 5) {
					mime_type = "image/jpeg";
					the_url = the_url.replace("data:base64","data:image/jpeg;base64");
				}

				var mime_type = "";

				if (the_url.indexOf("image/jpeg;base64") == 5) { mime_type = "image/jpeg"; }
				if (the_url.indexOf("image/png;base64") == 5) { mime_type = "image/png"; }
				if (the_url.indexOf("image/gif;base64") == 5) { mime_type = "image/gif"; }

				if (mime_type.length == 0) {
					appAlert([{ type: 'bg-red', message: "File " + file_name + " Not Allowed!", }])
				} else {

					var source_img_obj = document.createElement('img');
					source_img_obj.src = the_url;

					var callback = function(source) {
						if(!source) source = this;

						var quality = 80;
						var maxWidth = 600;
						maxWidth = maxWidth || 1000;
						var natW = source_img_obj.naturalWidth;
						var natH = source_img_obj.naturalHeight;
						var ratio = natH / natW;
						
						natW = maxWidth;
						natH = ratio * maxWidth;

						cvs = document.createElement('canvas');
						cvs.width = natW;
						cvs.height = natH;

						var ctx = cvs.getContext("2d").drawImage(source_img_obj, 0, 0, natW, natH);
						var compressed_image_url = cvs.toDataURL(mime_type, quality / 100);
						
						$('#' + imageID).val(compressed_image_url);
						$('#' + imageID + 'Src').attr('src', compressed_image_url);
						$('#' + imageID + 'Src').css('height',"100% !important");
						$('#' + imageID + 'Name').val(file_name);
					}

					source_img_obj.onload = callback;
				}				
			};
		}(imageInputFile[counter].name);
		reader.readAsDataURL(imageInputFile[counter]);
	}
}