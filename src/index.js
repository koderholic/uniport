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
		// m.mount(document.getElementById('appMenu'), menu)
	},
	view:function(vnode){
		return (
			<section style="" class="min-vh-100 mw8 center w-100 ">
				<article class="min-vh-100 dt w-100">
				  <div class="dtc v-mid tc white ph2">
						<div class="center">
							<img class="db center h3" src="../../assets/img/logo.png" />

						  <div class="fl w-100 w-50-ns pa3">
								<div class="ph3 w-100 bg-home-banner shadow-1 br3 br2 ba b--black-50">
									<article class="h5-l dt w-100">
										<div class="dtc v-mid tl white ph3 ph4-l monaco ">
											<div class="cf mv2"></div>
											<div class="pb3 f5 f4-l tracked fw5">
												 uniport Secure Wallet!!
											</div>

											<p class="f5 tracked dn dib-l">
												uniport secure wallet is a blockchain powered smart wallet for dApps, tokens, cryptocurrencies, digital files, private and public records.
											</p>
											<p class="f5 tracked dn dib-ns">
												#TRULY PRIVATE
												Private keys are under client control, they are never sent or stored outside of your device.
											</p>
											<p class="f5 tracked dn dib-ns">
												#TRULY PRIVATE
												Private keys are under client control, they are never sent or stored outside of your device.
											</p>
											<p class="f5 tracked dn dib-ns">
												#DAPPS PLAYER
												Interact with smartcontracts that adhere to either the ERC20 or ERC721 standard or any standard in future
											</p>
										</div>
									</article>
								</div>
							</div>

							<div class="fl w-100 w-50-ns pa3">
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
										<span class="fl w-100 center br2 relative">
											{m("input", {class:"f6 tracked bn black bg-light-gray pa3 br2 w-90-l w-80", placeholder:"Enter password..."})}
											<Icons name="chevron-right" class="bg-body h1 w1 pv3 ph3 ph2-ns right-0 absolute white pointer br2 br--right"/>
										</span>
									</div>

									<div class={"ph3 f6 avenir black tl cf pv4 "+page.classRegister}>
										<span class="fl w-100 center br2 relative">
											{m("input", {class:"f6 tracked bn black bg-light-gray pa3 br2 w-90-l w-80", placeholder:"Enter Password"})}
											<Icons name="lock-locked" class="bg-body h1 w1 pv3 ph3 ph2-ns right-0 absolute white pointer br2 br--right"/>
										</span>
										 <div class="cf mv2"></div>
										 <span class="fl w-100 center br2 relative">
 											{m("input", {class:"f6 tracked bn black bg-light-gray pa3 br2 w-90-l w-80", placeholder:"Confirm Password"})}
 											<Icons name="lock-locked" class="bg-body h1 w1 pv3 ph3 ph2-ns right-0 absolute white pointer br2 br--right"/>
 										</span>
										 <div class="cf mv3"></div>
										<div class=" tc">
											<span class="bg-body near-white shadow-4 pointer fl w-100 dim pv3 br2" onclick={page.Submit}>Create Your Wallet </span>
										</div>
									</div>
								</div>

								<div class={page.classLogin}>
									<div class="w-100 pv2 cf "></div>
									<div class="ph3 w-100 bg-white shadow-1 br2">
										<div class="f6 avenir black tl cf pv2">
											<a href="/import" oncreate={m.route.link} class="gray no-underline ph1 br1">
												<span class="fw5 db">Restore account?</span>
												<span class="fw5 db navy">Import using account seed phrase</span>
											</a>
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
