import m from 'mithril';
import Icons from './#icons.js';
import {validateSubmit} from './#validateSubmit.js';

var action = {
	Submit: function() {
		var actionFields = [
			{validationType : 'email', fieldID : 'email'},
		]
		validateSubmit( "/api/forgot", actionFields);
	},
};


export var page = {
	oninit:function(vnode){},
	view:function(vnode){
		return (
			<section style="" class="min-vh-100 mw8 center w-100 ">
				<article class="min-vh-100 dt w-100">
				  <div class="dtc v-mid tc white ph2">
						<div class="center w-50-l w-70-m">



							<div class="fl w-100 pa3">
								<div class=" w-100 bg-white br2 shadow-1">

									<div class="monaco center black flex flex-row">
										<div class={"w-100 fw5 pa3 bb pointer "+page.classBtnLogin} onclick={page.switchFormLogin}>
											Restore Your Account
										</div>
									</div>


									<div class={"ph3 f6 avenir black tl cf pv3 "}>
										<span class="fl w-100 center br2 relative">
											{m("textarea", {class:"f6 h4 tracked bn black bg-light-gray pa3 br2 w-100", placeholder:"Enter Account Seed"})}
										</span>
										<div class="cf mv2"></div>
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
											<span class="bg-body near-white shadow-4 pointer fl w-100 dim pv3 br2" onclick={page.Submit}>Restore Account</span>
										</div>
									</div>
								</div>
								<div class="w-100 pv2 cf "></div>
								<div class="ph3 w-100 bg-white shadow-1 br2">
									<div class="f6 avenir black tl cf pv2">
										<a href="/" oncreate={m.route.link} class="gray no-underline ph1 br1">
											<span class="fw5 db">Already Registered?</span>
											<span class="fw5 db navy">Login using password</span>
										</a>
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
