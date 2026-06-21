import masonry from "../../data/masonry.json"

const state = {
    imgs : masonry.imgs,
    imagearray:masonry.imagearray
};

export default {
    namespaced: true,
    state,
}