var m = require("mithril")
import {appAlert} from './#utils.js';
import {checkRedirect} from './#utils.js';

export function validateForm(postUrl, objValidate, objForm) { 

	var  Alert = []
	for (var objectField of objValidate) {
		var error = false;	
		if (objForm[objectField.fieldID] == undefined) {
			var errorMessage = objectField.fieldID + " is missing";
			if (objectField.fieldName !== "") {
				 errorMessage = objectField.fieldName + " is missing";
			}
			Alert.push({ type: 'bg-red', message: errorMessage, }); 

		} else {

			switch (objectField.fieldID) {
				case "Username":
					objForm[objectField.fieldID] = objForm[objectField.fieldID].replace(/\s/g, "")
					break;
			}

			if (objectField.fieldType !== "number") {
				objForm[objectField.fieldID] = objForm[objectField.fieldID].trim();
			}



			switch(objectField.fieldType) {
				default: if(objForm[objectField.fieldID].trim() == "") {error = true;}  break;
				case "sentence": if(!objForm[objectField.fieldID].match(/^\s*\S+(?:\s+\S+){2,}\s*$/)) {error = true;} break;
				case "email": if(!objForm[objectField.fieldID].match(/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/)) {error = true;} break;
				case "mobile": if(objForm[objectField.fieldID].length < 10) {error = true;}  break;
				case "image": if(objForm[objectField.fieldID].trim() == "" || objForm[objectField.fieldID].trim() == "default.jpg") {error = true;}  break;
				case "number": if(objForm[objectField.fieldID]==0) {  error = true;} break;
			}

			
		
			switch (objectField.fieldID) {
				case "ConfirmPassword":
					if(!error && objForm.Password !== objForm.ConfirmPassword ) { 
						Alert.push({ type: 'bg-red', message: "Password and Confirm Password do not Match", }); 
					}
					break;
			}

			if(error) { 
				var errorMessage = objectField.fieldID + " is Required";
				if (objectField.fieldName !== "") {
					 errorMessage = objectField.fieldName + " is Required";
				}
				Alert.push({ type: 'bg-red', message: errorMessage, }); 
			}

		}
	}

	// if (Alert.length === 0) {
	// 	startLoader();
	// 	m.request({ method: 'POST', url: postUrl, data: objForm, }).then(function(response) {

	// 		var Alert = [];
	// 		var alertType = 'bg-dark-green';


	// 		if (response.Code !== 200) { alertType = 'bg-red';  }

	// 		if (response.Message !== null && response.Message !== undefined && response.Message !== "" ){
	// 			Alert.push({ type: alertType, message: response.Message });
	// 		}

	// 		if (response.Error !== null && response.Error !== undefined && response.Error !== "" ){
	// 			Alert.push({ type: 'bg-red', message: response.Error });
	// 		}


	// 		if(Alert.length > 0) { appAlert(Alert); }
	// 		checkRedirect(response);

	// 		stopLoader();
	// 	}).catch(function(error) {
	// 		appAlert([{ type: 'bg-red', message: "Network Connectivity Error \n Please Check Your Network Access \n "+error, }]);
	// 		stopLoader();
	// 	});
	// }

	if (Alert.length != 0) {
		appAlert(Alert)
		return false;
	} else {
		return true;
	}
}