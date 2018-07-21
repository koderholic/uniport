var m = require("mithril");
import Icons from './#icons.js';

// export function menu() {
// 	m.render(document.getElementById('appMenu'), m(menu))
// }

export var menu = {
	menuFixed: "bg-transparent",
	oninit: function() {
		window.addEventListener('scroll', function() {
			var menuFixed;
			var shrinkOn = 160;
			var distanceY = window.pageYOffset || document.documentElement.scrollTop;
			if (distanceY > shrinkOn) { menuFixed = "top-nav fadeIn fixed"; }
			else { menuFixed = "bg-transparent"; }

			if (menuFixed !== menu.menuFixed) {
				menu.menuFixed = menuFixed;
				m.redraw();
			}
		});
	},
	linkItem : {
		view: function(vnode) {
			return(
				<a class="link" href={vnode.attrs.href}>
					<li class="tr" onclick={menu.toggle}>
						<p class="ph2 pv3 mv0 white hover-bg-white hover-blue fw5 tracked">
							{vnode.children}
						</p>
					</li>
				</a>
			)
		}
	},
	menuItem : {
		view: function(vnode) {
			return(
				<a class="link f5" oncreate={m.route.link} href={vnode.attrs.href}>
					<li class="tr" onclick={menu.toggle}>
						<p class="ph2 pv3 mv0 dark-green hover-bg-gradient hover-white fw5 tracked">
							{vnode.children}
						</p>
					</li>
				</a>
			)
		}
	},
	toggle: function() {
		var appmenuToggle = document.getElementById("menuToggle");
		var appmenuCover = document.getElementById("menuCover");
		appmenuCover.classList.toggle('dn');
		appmenuToggle.classList.toggle('animated');
		appmenuToggle.classList.toggle('bounceInRight');

		// document.getElementById("nav").classList.toggle('dn');
		// document.getElementById("menuBlur").classList.toggle('vh-100');
		document.getElementById("html").classList.toggle('overflow-hidden');
	},
	view: function(vnode) {
		return (
			<section id="menuBlur" class={"z-max w-100  "+menu.menuFixed}>
				<div id="menuCover"  class=" absolute right-0 w-100 vh-100 fr dn pa0" style="">
					<ul id="menuToggle" class="fr list pl0 w-70 w-40-m vh-100 ma0 bg-black-70" style="">
						<li class="tr">
							<p class="ph2 mv0 gray hover-red">
								<Icons name="cancel" class=" mh2 mv3 h1 dim dib white" onclick={menu.toggle}/>
							</p>
						</li>
						<li class="tr" onclick={menu.toggle}>
							<a oncreate={m.route.link} class="link mh2 link white f6 fr" href="/webshop">
									<span class="pa2 ph3 bw1 ba br4 b--white-20 fw5">Rinkeby Test Network </span>
								</a>
						</li>

					</ul>
				</div>

				<nav id="nav" class="w-100 mw9 mw7-m center black z-5">
					<div class="w-100 mw9 mw7-m center cf tc inline-flex items-center ph3 ph5-l pv2">
						<img class="h2" src="../../assets/img/logo.png" />
						<span class="ph2 fw6 white tracked"> uniport </span>
						<a oncreate={m.route.link} class="dn dib-ns link mh2 link white f6 fr" href="#">
							<span class="pa2 ph3 bw1 ba br4 b--white-20 fw5">Rinkeby Test Network </span>
						</a>
						<div class=" fr">
								<Icons name="menu" class=" white h1 fr" onclick={menu.toggle}/>
						</div>
					</div>
				</nav>




				<div id="appAlert"></div>
			</section>
		)
	}
}

export default menu;
