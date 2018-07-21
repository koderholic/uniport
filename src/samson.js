import m from 'mithril';
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
			<section class="menuCloudBG min-vh-100">
				<div id="appAlert"></div>
				<section class="mw9-ns center pa2 black-80 flex flex-row justify-center">

					<ul class="list pl0 mt0 measure center">
					  <li
					    class="flex items-center lh-copy pa3 ph0-l bb b--black-10">
					      <img class="w2 h2 w3-ns h3-ns br-100" src="http://tachyons.io/img/avatar-mrmrs.jpg" />
					      <div class="pl3 flex-auto">
					        <span class="f6 db black-70">Savar</span>
					        <span class="f6 db black-70">Victoria Island, lagos</span>
					      </div>
					      <div>
					        <a href="tel:" class="f6 link blue hover-dark-gray">+101100011</a>
					      </div>
					  </li>
					  <li
					    class="flex items-center lh-copy pa3 ph0-l bb b--black-10">
					      <img class="w2 h2 w3-ns h3-ns br-100" src="http://tachyons.io/img/avatar-jxnblk.jpg" />
					      <div class="pl3 flex-auto">
					        <span class="f6 db black-70">Jibs</span>
					        <span class="f6 db black-70">Lekki, lagos</span>
					      </div>
					      <div>
					        <a href="tel:" class="f6 link black hover-dark-gray">10111111001</a>
					      </div>
					  </li>
					  <li
					    class="flex items-center lh-copy pa3 ph0-l bb b--black-10">
					      <img class="w2 h2 w3-ns h3-ns br-100" src="http://tachyons.io/img/avatar-jasonli.jpg" />
					      <div class="pl3 flex-auto">
					        <span class="f6 db black-70">Authentic Sam</span>
					        <span class="f6 db black-70">Ajah, lagos</span>
					      </div>
					      <div>
					        <a href="tel:" class="f6 link blue hover-dark-gray">101101110</a>
					      </div>
					  </li>
					  <li
					    class="flex items-center lh-copy pa3 ph0-l bb b--black-10">
					      <img class="w2 h2 w3-ns h3-ns br-100" src="http://tachyons.io/img/avatar-yavor.jpg" />
					      <div class="pl3 flex-auto">
					        <span class="f6 db black-70">Scarflamez</span>
					        <span class="f6 db black-70">Banana Island</span>
					      </div>
					      <div>
					        <a href="tel:" class="f6 link blue hover-dark-gray">01110111101</a>
					      </div>
					  </li>
					</ul>				

				</section>
			</section>
		)
	}
}


export default page;
