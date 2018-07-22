import m from 'mithril';
import menu from './#menu.js';
import footer from './#footer.js';

import Icons from './#icons.js';
import {appAlert} from './#utils.js';
import {checkRedirect} from './#utils.js';
import {defaultImage} from './#utils.js';
import {displayImage} from './#utils.js';

var searchXHR = null
var searchTimer;

var page = {
	oninit:function(){
		m.mount(document.getElementById('appMenu'), menu)
		m.mount(document.getElementById('appFooter'), footer);
	},
	oncreate:function(vnode){},
	view:function(vnode){
	return (
		<section class="center min-vh-100">

			<section class="mw8 center dt w-100 bg-white pt5" style="min-height: 180px;">
				<div class="dtc v-mid tc black ph3 ph4-l">
					<h1 class="f2 f1-l fw6 tc athelas i">About Us</h1>
				</div>
			</section>



			<section class="mw8 cf center">
				<div class="fl w-100 w-30-l tc pv3 ph2">
					<img class="br3 shadow-5" src="assets/img/africa.jpg"/>
				</div>

				<div class=" bg-white fl bg-white  w-100 w-70-l pt2 pb2 ph2 tj">
					<p class="f3 lh-copy">
						The Greenhouse provides a safe space for women and youths, where they can come to gain access to resources that enable their creative expressions.
					</p>
					<p class="f5 lh-copy">
						Greenhouse provides an enabling environment for women and youths to meet, work, learn and collaborate,
						while ensuring that they become innovators.
						Comprehensive programs which include practical ICT sessions, Entrepreneurship workshops,
						Financial Literacy inclusion, Social Advocacy - online and offline peer meets, Career Talks
						as well as interactive extra-curricular activities e.g. book reading, gaming and dancing sessions
						are all included in the hub activities to help the young women in their growth to running
						productive and sustainable businesses.
					</p>
				</div>
			</section>

			<section class="bg-gold">
			<section class="mw8 center ph3-ns tracked pv1 fw4 black tc">
				<div class="fl w-100 pv3">
					<div class="pa2">
						<img src="../../assets/img/goals.svg" class="h4 dn" />
						<h1 class="mt3 ">OUR GOALS</h1>
						<p class="f6 f5-l tc">is to add value to potential clients and grow into being a </p>
						<p class="f6 f5-l tc">recognized black owned company that services the needs of all stake holders.</p>
					</div>
				</div>
				<div class="cf pv3 w-100">
					<div class="fl w-100 w-50-ns">
						<div class="pa2">
							<img src="../../assets/img/vision.svg" class="h4" />
							<h1 class="mt3">VISION</h1>
							<p class="f6 tl">
								<ul>
									<li class="pv2">To grow and be a recognized company, also acquire endorsement from reputable bodies in the country</li>
									<li class="pv2">We also want to                          grow by expanding our area of operation to other regions in the country and abroad</li>
									<li class="pv2">Engage ourselves in large operations and opening branches</li>
								</ul>
							</p>
						</div>
					</div>
					<div class="fl w-100 w-50-ns">
						<div class="pa2">
							<img src="../../assets/img/mission.svg" class="h4" />
							<h1 class="mt3">MISSION</h1>
							<p class="f6 tl">
								<ul>
									<li class="pv2">To contribute in the social development of our communities</li>
									<li class="pv2">To be one of the leading companies in the Supply chain and construction business</li>
									<li class="pv2">To provide an excellent service to our clients across various sectors of our economy</li>
								</ul>
							</p>
						</div>
					</div>
				</div>
			</section>
			</section>

			<section class="mw8 bg-white cf center pv3">
				<div class="br2-l br--bottom-l fl bg-white  w-100 w-70-l pt2 pb4 ph2 tj">
					<p class=" lh-copy">
						The Greenhouse provides a safe space for women and youths, where they can come to gain access to resources that enable their creative expressions.
					</p>
					<p class=" lh-copy">
						Greenhouse provides an enabling environment for women and youths to meet, work, learn and collaborate,
						while ensuring that they become innovators.
						Comprehensive programs which include practical ICT sessions, Entrepreneurship workshops,
						Financial Literacy inclusion, Social Advocacy - online and offline peer meets, Career Talks
						as well as interactive extra-curricular activities e.g. book reading, gaming and dancing sessions
						are all included in the hub activities to help the young women in their growth to running
						productive and sustainable businesses.
					</p>
				</div>
				<div class="fl w-100 w-30-l pb3 ph2">
					<img class="br3 shadow-5" src="assets/img/africa.jpg"/>
				</div>
			</section>


		</section>
	)
  }
}

export default page;
