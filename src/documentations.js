import m from 'mithril';

import Icons from './#icons.js';
import {appAlert} from './#utils.js';
import {checkRedirect} from './#utils.js';

var searchXHR = null
var searchTimer;

var page = {
	Url: "/api/documentations", Form: {},
	searchText: "", searchView: "dn", formView: "dn", searchResult:[],
	searchDocumentations:function(){
		if (searchTimer){ clearTimeout(searchTimer); }
		if (searchXHR !== null) { searchXHR.abort() } searchXHR = null;

		page.searchResult = []
		searchTimer = setTimeout(function(){
			m.request({ method: 'GET', url: page.Url+"/list?search="+page.searchText,
				config: function(xhr) {searchXHR = xhr}, }).then(function(response) {

				var searchList = [];
				checkRedirect(response);

				if (response.Code == 200) {
					if (response.Body !== null && response.Body !== undefined ){
						response.Body.map(function(result) { if (result.ID > 0) {
							result.Updatedate = result.Updatedate.slice(0, 10);
							searchList.push( m(searchDocumentationsResult,
								{Code: result.Code, Category: result.Category, Title: result.Title,
									Author: result.Author, Updatedate: result.Updatedate}
							))
						}})
						page.formView = "dn"; page.searchView = "";
						page.searchResult = searchList
					}
				}
				if(searchList.length == 0) {
					page.searchResult = m("small",{class:"i f6"}, "No Result")
				}
			}).catch(function(error) {
				appAlert([{ type: 'bg-red', message: "Network Connectivity Error \n Please Check Your Network Access", }]);
			});
		}, 750);
	},

	readDocumentations:function(vnode){
		if (searchXHR !== null) { searchXHR.abort() } searchXHR = null;
		// page.formView = ""; page.searchView = "dn";
		m.request({ method: 'GET', url: page.Url+"/read?search="+vnode.attrs.path,
			config: function(xhr) {searchXHR = xhr}, }).then(function(response) {
			checkRedirect(response);

			if (response.Code == 200) {
				if (response.Body !== null && response.Body !== undefined ){
					page.Form = response.Body;
					page.Form.File = m.trust(page.Form.File)
					page.formView = ""; page.searchView = "dn";

					page.Form.Updatedate = page.Form.Updatedate.slice(0, 10);
				}
			}

		}).catch(function(error) {
			appAlert([{ type: 'bg-red', message: "Network Connectivity Error \n Please Check Your Network Access", }]);
		});
	},
	oninit:function(){ },
	oncreate:function(vnode){
		(vnode.attrs.path == undefined) ? page.searchDocumentations() : page.readDocumentations(vnode)
	},
	view:function(vnode){
	return (
		<section class="center bg-white min-vh-100 pb5">

			<div id="appAlert"></div>

			<article class="vh-25 dt w-100 documentationBG center">
				<div class="dtc v-mid center near-white bg-black-70">

				<section class="mw7 center ph3 tc">

					<h1 class="f4 f3-m f2-l fw4 tcnear-white">Welcome to Our Documentation</h1>
					<p class="f6 tc">
						Here you will find lots of helpful tips in the form of short tutorial videos, illustrated screenshots and FAQs.
					</p>


					<p class="tc f6 w-100">
						{m('input', {type:"text", class:"pa2 mb2 w-70 w-50-l dib bw0", placeholder:"Enter your search term here...",
							onchange: m.withAttr("value",function(value){page.searchText = value;})
						})}
						{m('button', {type:"text", class:"pa2 mb2 bg-dark-green near-white dib bw0", onclick: page.searchDocumentations},"Search")}
					</p>

				</section>


				</div>
			</article>

			<section class="mw7 center dark-gray pa3">

				<section class={page.searchView}>{page.searchResult}</section>

				<section class={page.formView}>
					<h1 class="f6 f5-m f4-l fw4 tc dark-gray b">{page.Form.Category}</h1>
					<h1 class="f4 f3-m f2-l fw4 tc dark-gray">{page.Form.Title}</h1>
					<p class="tc f6"> by {page.Form.Author} <br/> <small class="i">{page.Form.Updatedate}</small></p>

					<br/><br/>
					{page.Form.File}
					<a href="/documentations" class="no-underline fr ph2 pv1 f6 br1 bg-black near-white dim pointer">
						Back To Documentation
					</a>
				</section>

			</section>

		</section>
	)
  }
}


var searchDocumentationsResult = {view: function(vnode) {return(
	<p class="f6"> <Icons name="book" class="h1 dark-gray"/>
		<a href={"documentations/"+vnode.attrs.Code} class="no-underline dark-gray ph1 fw5">  {vnode.attrs.Category}:- <span class="fw3">{vnode.attrs.Title}</span> <br/> <small> by {vnode.attrs.Author} <i>{vnode.attrs.Updatedate}</i></small> </a>
	</p>
)}}

export default page;
