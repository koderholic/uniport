import m from 'mithril';
import Siema from 'siema';

import menu from './#menu.js';
import footer from './#footer.js';
import Icons from './#icons.js';

import {appAlert} from './#utils.js';
import {checkRedirect} from './#utils.js';


var page = {
	classLogin:"",classRegister:"dn",
	classBtnLogin:"b--navy navy",
	classBtnRegister:"b--light-gray light-gray",

	switchFormLogin: function() {
		page.classLogin = "";
		page.classRegister = "dn";
		page.classBtnLogin = "b--navy navy";
		page.classBtnRegister = "b--light-gray light-gray";
	},
	switchFormRegister: function() {
		page.classLogin = "dn";
		page.classRegister = "";
		page.classBtnLogin = "b--light-gray light-gray";
		page.classBtnRegister = "b--navy navy";
	},

	Submit: function() {
		var actionFields = [
			{validationType : '', fieldID : 'username'},
			{validationType : '', fieldID : 'password'},
		]
		validateSubmit( "/api/login", actionFields);
	},

	oninit:function(vnode){
		m.mount(document.getElementById('appMenu'), menu)
	},
	view:function(vnode){
		return (
			<section style="" class="min-vh-100 mw8 center w-100 ">
				<div class="cf mv3"></div>
				<article class="min-vh-100 dt w-100">
				  <div class="dtc v-mid tc white ph2">
						<div class="center">
							<div class="fl w-100 w-50-l  ph3 pv3-l dn ">
								<article class="h5-l dt w-100  w-80-ns bg-home-banner br3">
									<div class="dtc v-mid tc white ph3 pv4-l ph4-l monaco ">
										<div class="cf mv2"></div>
										<div class="pb3-l f3 f2-l tracked fw5">
											Welcome back!
										</div>

										<div class="pv3 f5 f4-l tracked fw5">
											 Backpocket Secure Wallet
											 <br/>
											  at your service!!
										</div>

									</div>
								</article>
							</div>

						  <div class="dn fl w-100 w-50-l bg-home-banner shadow-0 br3 ">
								<article class="h5-l dt w-100">
									<div class="dtc v-mid tl white ph3 ph4-l monaco ">
										<div class="cf mv2"></div>
										<div class="pb3-l f3 f2-l tracked fw5">
											Welcome back!
										</div>

										<div class="pb3 f5 f4-l tracked fw5">
											 Backpocket Secure Wallet at your service!!
										</div>

										<p class="f5 tracked dn dib-l">
											Backpocket is an HD (Hierarchical Deterministic) wallet.
										</p>
										<p class="f5 tracked dn dib-l">
											This means that a unique address is generated every time you send or receive funds,
											making your transaction activity and total balance much harder to track,
											and keeping your financial business yours.
										</p>
						  		</div>
						  	</article>
						  </div>


						  <div class="fl w-100 w-50-ns pa3">
								<div class="ph3 w-100 bg-home-banner shadow-1 br3 br2 ba b--black-50">
									<article class="h5-l dt w-100">
										<div class="dtc v-mid tl white ph3 ph4-l monaco ">
											<div class="cf mv2"></div>
											<div class="pb3-l f4 f2-l tracked fw5">
												Welcome back!
											</div>

											<div class="pb3 f5 f4-l tracked fw5">
												 Backpocket Secure Wallet at your service!!
											</div>

											<p class="f5 tracked dn dib-l">
												Backpocket is an HD (Hierarchical Deterministic) wallet.
											</p>
											<p class="f5 tracked dn dib-l">
												This means that a unique address is generated every time you send or receive funds,
												making your transaction activity and total balance much harder to track,
												and keeping your financial business yours.
											</p>
										</div>
									</article>
								</div>
							</div>


							<div class="fl w-100 w-50-ns pa3">
								<div class="ph3 w-100 bg-white shadow-1">
									<div class="f6 avenir black tl cf pv2">
										<a href="/forgot" oncreate={m.route.link} class="gray no-underline ph1 br1">
											<span class="fw5 db">Restore account?</span>
											<span class="fw5 db navy">Import using account seed phrase</span>
										</a>
									</div>
								</div>
									<div class="w-100 pv2 cf "></div>
								<div class=" w-100 bg-white br2 shadow-1">
									<div class="monaco center black flex flex-row">
										<div class={"w-50 fw5 pa3 bb pointer "+page.classBtnLogin} onclick={page.switchFormLogin}>
											Login
										</div>
										<div class={"w-50 fw5 pa3 bb pointer "+page.classBtnRegister} onclick={page.switchFormRegister}>
											Register
										</div>
									</div>

									<div class={"ph3 f6 avenir black tl cf pv4 "+page.classLogin}>
										{m("div", {class:"br2 ba b--silver"} ,m("input",{ placeholder: "Password", type:"password", class: "w-100  bw0 br2 ph2 pv3 f6", id:"username",
											oninput: m.withAttr("value",function(value) {page.Username = value}),
											onkeyup: function(event) {if(event.key=="Enter"){action.Submit()}}
										 }))}
										 <div class="cf mv2"></div>
										<div class=" tc">
											<span class="bg-body near-white shadow-4 pointer fl w-100 dim pv3 br2" onclick={page.Submit}>Log me in » </span>
										</div>
									</div>

									<div class={"ph3 f6 avenir black tl cf pv4 "+page.classRegister}>
										{m("div", {class:"br2 ba b--silver"} ,m("input",{ placeholder: "New Password", type:"password", class: "w-100  bw0 br2 ph2 pv3 f6", id:"username",
											oninput: m.withAttr("value",function(value) {page.Username = value}),
											onkeyup: function(event) {if(event.key=="Enter"){action.Submit()}}
										 }))}
										 <div class="cf mv2"></div>
										 {m("div", {class:"br2 ba b--silver"} ,m("input",{ placeholder: "Confirm Password", type:"password", class: "w-100  bw0 br2 ph2 pv3 f6", id:"username",
												oninput: m.withAttr("value",function(value) {page.Username = value}),
												onkeyup: function(event) {if(event.key=="Enter"){action.Submit()}}
											 }))}
										 <div class="cf mv2"></div>
										<div class=" tc">
											<span class="bg-body near-white shadow-4 pointer fl w-100 dim pv3 br2" onclick={page.Submit}>Create Wallet » </span>
										</div>
									</div>
								</div>
						  </div>
						</div>
				  </div>
				</article>

			</section>
		)
	}
}

export default page;
