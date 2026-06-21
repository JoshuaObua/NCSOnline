import bookmark from "../../data/bookmark"
import projects from '../../data/project';

const state = {
	allprojects : projects.all,
	doingprojects : projects.Doing,
	doneprojects : projects.Done,
    bookmark : bookmark.data,
};



export default {
	namespaced: true,
	state,
	
};