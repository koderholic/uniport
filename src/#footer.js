var m = require("mithril");
import {appAlert} from './#utils.js';

export var footer = {
	clickedNewsletter: false,
	IpAddress:"", UserAgent:"",
	FormNewsletter : {Firstname:"",Lastname:"",Email:""},
	submitNewsletter: function() {
		var alert = []
		if (footer.clickedNewsletter) {
			console.log("footer.clickedNewsletter: "+footer.clickedNewsletter)
			appAlert([{ message: "Signed up already!!" }]); return
		}


		footer.FormNewsletter.IpAddress = footer.IpAddress
		footer.FormNewsletter.UserAgent = footer.UserAgent

		if (footer.FormNewsletter.Firstname.length == 0) { alert.push({ message: "First Name is required" }); }
		else if (footer.FormNewsletter.Lastname.length < 3) { alert.push({ message: "Last Name is required" });}

		if (footer.FormNewsletter.Lastname.length == 0) { alert.push({ message: "Last Name is required" }); }
		else if (footer.FormNewsletter.Lastname.length < 3) { alert.push({ message: "Last Name is too short" }); }

		if (footer.FormNewsletter.Email.length == 0) { alert.push({ message: "Email is required" }); }
		else if(!footer.FormNewsletter.Email.match(/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/)) {
			alert.push({ message: "Email is invalid" });
		}

		if (alert.length > 0) {
			footer.clickedNewsletter = false
			appAlert(alert)
			return
		}

		// startLoader();
		footer.FormNewsletter.List = "signup-newsletter"
		footer.FormNewsletter.Notify = "signup-newsletter-notify"
		m.request({ method: 'POST', url: "/api/signup-newsletter", data: footer.FormNewsletter, }).then(function(response) {
			var lStoploader = true;
			if (response.Message !== null &&  response.Message !== "") {
				appAlert([{ message: response.Message }]);

			}
			// if(lStoploader) { stopLoader();}
		}).catch(function(error) {
			appAlert([{ message: error }]);
			// stopLoader();
		});

		footer.clickedNewsletter = true
		footer.FormNewsletter = {Firstname:"",Lastname:"",Email:""};
	},
	oninit: function(){
		m.request({method:'GET', url: "https://icanhazip.com/",
			deserialize: function(value) {return value}}).then(function(response){
			footer.IpAddress = response;
			footer.UserAgent = navigator.userAgent;
		});
	},
	view: function(vnode) {
		return (
			<footer class="mw8 center ">
				<div class="fl w-100 pv4 ph3 ph5-m ph6-l washed-green">
				  <small class="f6 db tc">© 2018 <b class="ttu tracked">Green House</b> - All Rights Reserved</small>
				  <div class="tc mt3">
				    <a href="/terms/" title="Terms" class="f6 dib ph2 link washed-green dim">Terms of Use</a>
				    <a href="/privacy/" title="Privacy" class="f6 dib ph2 link washed-green dim">Privacy</a>
				  </div>
				</div>
			</footer>
		)
	}
}

export default footer;
