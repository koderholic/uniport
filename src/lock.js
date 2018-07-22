import m from 'mithril';
import Siema from 'siema';

import menu from './#menu.js';
import footer from './#footer.js';
import Icons from './#icons.js';

import {appAlert} from './#utils.js';
import {checkRedirect} from './#utils.js';


var page = {
	classLogin:"",classRegister:"dn",
	classBtnLogin:"b--green green",
	classBtnRegister:"b--light-gray light-gray",

	switchFormLogin: function() {
		page.classLogin = "";
		page.classRegister = "dn";
		page.classBtnLogin = "b--green green";
		page.classBtnRegister = "b--light-gray light-gray";
	},
	switchFormRegister: function() {
		page.classLogin = "dn";
		page.classRegister = "";
		page.classBtnLogin = "b--light-gray light-gray";
		page.classBtnRegister = "b--green green";
	},

	Submit: function() {
		var actionFields = [
			{validationType : '', fieldID : 'username'},
			{validationType : '', fieldID : 'password'},
		]
		validateSubmit( "/api/login", actionFields);
	},

	oninit:function(vnode){
		// m.mount(document.getElementById('appMenu'), menu)
	},
	view:function(vnode){
		return (
			<section style="" class="vh-100 mw8 center w-100 ">
				<article class="min-vh-100 dt w-100">
				  <div class="dtc v-mid tc white ph2">
						<div class="mw7 center tc pa3">
							<img class="db center " src="../../assets/img/logo.png" />
							<div class="cf mv3"></div>
							<div class=" bg-white w-60-m w-50-l br2 shadow-1 pv3 center">
								<span class="center br2 relative">
									{m("input", {class:"f6 tracked bn black bg-light-gray pa3 br2 w-80", placeholder:"Enter password..."})}
									<Icons name="chevron-right" class="bg-body h1 w1 pv3 ph3 ph2-ns right-0 absolute white pointer br2 br--right"/>
								</span>
							</div>
						</div>
				  </div>
				</article>

			</section>
		)
	}
}

export default page;
