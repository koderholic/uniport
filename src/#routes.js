var m = require("mithril")


//Generic Website Pages
import indexPage from './index.js';
import lockPage from './lock.js';
import importPage from './import.js';
import walletPage from './wallet.js';

// import contactPage from './contactus.js';
// import newsPage from './news.js';
// import aboutusPage from './aboutus.js';
// import galleryPage from './gallery.js';
// import documentationPage from './documentations.js';


m.route.setOrig = m.route.set;
m.route.set = function(path, data, options){
	m.route.setOrig(path, data, options);
	window.scrollTo(0,0);
}

m.route.linkOrig = m.route.link;
m.route.link = function(vnode){
	m.route.linkOrig(vnode);
	window.scrollTo(0,0);
}

m.route.prefix("")
m.route.mode = "pathname"
m.route(document.getElementById('appContent'), "/", {
	"/":{ view: function(vnode) { return m(indexPage);},},
	"/lock":{ view: function(vnode) { return m(lockPage);},},
	"/import":{ view: function(vnode) { return m(importPage);},},
	"/wallet":{ view: function(vnode) { return m(walletPage);},},
});
