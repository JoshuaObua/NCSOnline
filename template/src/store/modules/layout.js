import Layout from '../../data/layout.json';

const state = {
	layout: Layout,
	sidebarType: localStorage.getItem('SidebarType') || 'default',
	boxlayout: true
};


const getters = {

};


const mutations = {
	set: (state) => {
		document.body.className = state.layout.color.mix_layout;
		document.body.setAttribute('main-theme-layout', state.layout.settings.layout_type);
		document.getElementsByTagName('html')[0].setAttribute('dir', state.layout.settings.layout_type);
		var primaryColor = localStorage.getItem('--theme-deafult') || state.layout.color.primary_color;
		var secondaryColor = localStorage.getItem('--theme-secondary') || state.layout.color.secondary_color;
		var layoutVersion = localStorage.getItem('layoutVersion') || state.layout.color.layout_version;
		if (primaryColor || secondaryColor) {
			addStyle(primaryColor, secondaryColor);
			if (layoutVersion)
				document.body.className = layoutVersion;
		}
	},
	setLayout: (state, payload) => {
	
		if((document.body.className=='rtl' || document.body.className=='dark-only' || document.body.className=='light-only rtl') && payload=='dark-only' ){
			document.body.className = payload+ ' rtl';
		}

		else if(payload=='light-only' && document.body.getAttribute('main-theme-layout')=='rtl')
	{
		document.body.className = 'rtl'
	}
	else if(payload=='light-only' && document.body.getAttribute('main-theme-layout')=='ltr')
	{
		document.body.className = 'ltr'
	}
	else if(payload=='light-only' && document.body.getAttribute('main-theme-layout')=='box-layout')
	{
		document.body.className = 'box-layout'
	}
		else if (document.body.className=='box-layout' || document.body.className=='dark-only' || document.body.className=='light-only box-layout' ||document.body.className=='light-only'){
			document.body.className = payload + ' box-layout';
		}
		else{
			document.body.className = payload ;
		}
		
	},
	setCustomizeSidebarType: (state, payload) => {
		localStorage.setItem('SidebarType', payload);
	},
	setColorScheme: (state, payload) => {
		setColor(state, payload);
		state.primaryColor = payload.primary;
		state.secondaryColor = payload.secondary;
		state.layout.color.layout_version = 'light';
		localStorage.setItem('layoutVersion', state.layout.color.layout_version);
	},
	setColorDarkScheme: (state, payload) => {
		setColor(state, payload);
		state.layout.color.layout_version = 'dark-only';
		localStorage.setItem('layoutVersion', state.layout.color.layout_version);
	},
	setLayoutType: (state, payload) => {
		
		document.body.removeAttribute('class','box-layout dark-only')
		document.body.removeAttribute('class','dark-only')
		document.body.removeAttribute('class','dark-only rtl')
		document.body.removeAttribute('class','ltr')	
		document.body.removeAttribute('class','rtl')

		payload.val && document.body.setAttribute('class', payload.val);
		payload.val && document.body.setAttribute('main-theme-layout', payload.val);
	
		state.layout.settings.layout_type = payload.val;
		document.getElementsByTagName('html')[0].setAttribute('dir', payload.val);
	},

};


const actions = {
	set: (context) => {
		context.commit('set');
	},
	setLayoutType: (context, payload) => {
		context.commit('setLayoutType', payload);
	},
	setLayout: (context, payload) => {
		context.commit('setLayout', payload);
	},
	setCustomizeSidebarType: (context, payload) => {
		context.commit('setCustomizeSidebarType', payload);
	},
	setColorScheme: (context, payload) => {
		context.commit('setColorScheme', payload);
	},
	setColorDarkScheme: (context, payload) => {
		context.commit('setColorDarkScheme', payload);
	},

};

function addStyle(primary, secondary) {
	document.documentElement.style.setProperty('--theme-deafult', primary);
	document.documentElement.style.setProperty('--theme-secondary', secondary);
}

function setColor(state, color) {
	addStyle(color.primary, color.secondary);
	localStorage.setItem('--theme-deafult', color.primary);
	localStorage.setItem('--theme-secondary', color.secondary);
	window.location.reload();
}

export default {
	namespaced: true,
	state,
	getters,
	actions,
	mutations
};
