import { createRouter, createWebHistory } from 'vue-router'
import BodyView from '../layout/BodyView.vue'
import auth from "../auth/authView.vue"
import login from "../auth/loginPage.vue"

import IndexDefault from '../components/dashboard/default/IndexDefault.vue'
import indexEcommerce from '../components/dashboard/e-commerce/indexEcommerce'


import indexGeneral from '../components/widgets/general/indexGeneral'
import indexChart from '../components/widgets/charts/indexChart'


import indexList from '../components/project/projectlist/indexList'
import createProject from '../components/project/createproject/createProject'

import kanbanBoard from "../components/kanban/kanbanBoard"
import file_manager from "../components/filemaneger/file_manager"


import indexProduct from "../components/ecommerce/product/indexProduct"
import cartView from "../components/ecommerce/cartView"
import indexAdd from "../components/ecommerce/addproduct/indexAdd"
import productDetails from "../components/ecommerce/productDetails/productDetails"
import checkOut from "../components/ecommerce/chekout/checkOut"
import wishList from "../components/ecommerce/wishList"
import invoiceView from "../components/ecommerce/invoice/invoiceView"
import PaymentDetails from "../components/ecommerce/paymentDetails/PaymentDetails"
import orderHistory from "../components/ecommerce/orderhistory/orderHistory"


import indexEmail from "../components/email/indexEmail"


import indexChat from "../components/chat/chatList/indexChat"
import videoChat from "../components/chat/videoChat/videoChat"


import userProfile from "../components/users/profile/userProfile"
import editProfile from "../components/users/editProfile/editProfile"
import cardIndex from "../components/users/userCards/cardIndex"

import bookMark from "../components/bookmark/bookMark.vue"


import socialApp from "../components/socialApp/socialApp"


import todoIndex from "../components/todo/todoIndex"


import serchIndex from "../components/search/serchIndex"


import formValidation from "../components/forms/formValidetion/formValidation"
import base_Input from "../components/forms/baseInput/base_Input"
import checkbox_radio from "../components/forms/Checkbox&Radio/checkbox_radio"
import input_groups from "../components/forms/InputGroup/input_groups"
import megaOptions from "../components/forms/megaOptions/megaOptions"


import select2 from "../components/formWidgets/select2/select2"
import switch_From from "../components/formWidgets/switch/switch_From"
import touchspin_Form from "../components/formWidgets/touchspin/touchspin_Form"
import typeahead_Form from "../components/formWidgets/typeahead/typeahead_Form"
import clipboard_Form from "../components/formWidgets/clipboard/clipboard_Form"
import datepicker from "../components/formWidgets/datepicker/datepicker"


import form_wizard from "../components/formwizard/form_wizard"


import basic_tables from "../components/tables/bootstrapTable/basicTables/basic_tables"
import sizing_tables from "../components/tables/bootstrapTable/sizeTable/sizing_tables"
import border_Tables from "../components/tables/bootstrapTable/borderTables/border_Tables"
import styling_table from "../components/tables/bootstrapTable/stylingtable/styling_table"


import basic_Init from "../components/tables/dataTable/BasicInit/basic_Init"


import state_Color from "../components/uikits/statecolor/state_Color"
import typograPhy from "../components/uikits/Typography/typograPhy"
import avatars_Uikits from "../components/uikits/Avatars/avatars_Uikits"
import helper_classes from "../components/uikits/helper/helper_classes"
import grid_Uikits from "../components/uikits/grid/grid_Uikits"
import tag_pills from "../components/uikits/Tag&Pills/tag_pills"
import progress_bar from "../components/uikits/progressBar/progress_bar"
import modal_Uikit from "../components/uikits/modal/modal_Uikit"
import alert_Uikit from "../components/uikits/alert/alert_Uikit"
import popover_Uikit from "../components/uikits/Popover/popover_Uikit"
import tooltip_uikits from "../components/uikits/tooltip/tooltip_uikits"
import spinners_Uikit from "../components/uikits/Spinners/spinners_Uikit"
import accordion_Uikit from "../components/uikits/accordion/accordion_Uikit"
import box_shadow from "../components/uikits/Shadow/box_shadow"
import list_Uikit from "../components/uikits/list/list_Uikit"
import dropdown_Uikit from "../components/uikits/Dropdown/dropdown_Uikit"


import scrollable_advance from "../components/advance/Scrollable/scrollable_advance"
import pagination_advance from "../components/advance/Pagination/pagination_advance"
import sweetAlert from "../components/advance/SweetAlert/sweetAlert"
import ribbons_advance from "../components/advance/Ribbons/ribbons_advance"
import breadCrumb from "../components/advance/Breadcrumb/breadCrumb"
import cropper_advance from "../components/advance/cropper/cropper_advance"
import toaster_advance from "../components/advance/Toaster/toaster_advance"
import tour_advance from "../components/advance/Tour/tour_advance"
import rating_advance from "../components/advance/Rating/rating_advance"
import upload_advance from "../components/advance/upload/upload_advance"
import sticky_advance from "../components/advance/sticky/sticky_advance"
import range_advance from "../components/advance/rangeSlider/range_advance"
import dragdrop from "../components/advance/dragdrop/dragdrop"


import flag_Icon from "../components/icon/flag_Icon"
import feather_icon from "../components/icon/feather_icon"
import weather_icon from "../components/icon/weather_icon"
import themify_icon from "../components/icon/themify_icon"
import font_Awesome from "../components/icon/fontAwesome/font_Awesome"
import icon_Icon from "../components/icon/icon_Icon"


import default_button from "../components/button/Default/default_button"
import flat_button from "../components/button/Flat/flat_button"
import edge_button from "../components/button/Edge/edge_button"
import raised_button from "../components/button/Raised/raised_button"
import button_group from "../components/button/ButtonGroup/button_group"


import google_chart from "../components/charts/googleChart/google_chart"
import apex_chart from "../components/charts/ApexChart/apex_chart"
import chartist_chart from "../components/charts/Chartist/chartist_chart"


import grid_gallery from "../components/gallery/grid_gallery"
import grid_desc from "../components/gallery/grid_desc"
import hover_gallery from "../components/gallery/hover-gallery/hover_gallery"
import masonry_gallery from "../components/gallery/masonry-gallery/masonry_gallery"
import masonary_desc from "../components/gallery/masonary_desc"

import errorPage1 from '../errors/errorPage1.vue';
import errorPage2 from '../errors/errorPage2.vue';
import errorPage3 from '../errors/errorPage3.vue';
import errorPage4 from '../errors/errorPage4.vue';




import BasicView from '../components/cards/basicView/basicView.vue';
import DraggableView from '../components/cards/draggableView';
import creative_card from "../components/cards/creative/creative_card"
import tabbed_card from "../components/cards/Tabbed/tabbed_card"


import login_image from "../components/authentication/login_image"
import login_image2 from "../components/authentication/login_image2"
import ragister_simple from "../components/authentication/ragister_simple"
import ragister_image from "../components/authentication/ragister_image"
import ragister_image2 from "../components/authentication/ragister_image2"
import login_with_validation from "../components/authentication/login_with_validation"
import unlock from "../components/authentication/unlock"
import forget_password from "../components/authentication/forget_password"
import reset_password from "../components/authentication/reset_password"
import login_simple from "../components/authentication/login_simple.vue"


import timeline from "../components/timeline/timeline"


import google_map from "../components/maps/google_map"
import vue_leaflet from "../components/maps/vue_leaflet"

import ckediter from "../components/editer/ckediter"
import simple_editer from "../components/editer/simple_editer"


import blog_detail from "../components/blog/blog-details/blog_detail"
import blog_single from "../components/blog/blog-single/blog_single"


import job_list from "../components/job/job-list/job_list"
import job_details from "../components/job/job_details/job_details"
import job_apply from "../components/job/job_apply/job_apply"


import learning_list from "../components/learning/learninglist/learning_list"
import coursedetailed from "../components/learning/coursedetailed"


import faqindex from "../components/faq/faqindex"



import knowledgebase from "../components/Knowledgebase/knowledgebase"


import support from "../components/support/support"


import pricing from "../components/Pricing/pricing"


import maintenance from "../components/maintenance"

import sample_page from "../components/sample_page"


import calender from "../components/calendar/calender"

import ComingsoonImage from '../components/comingsoon/comingsoon_image';
import ComingsoonSimple from '../components/comingsoon/comingsoon_simple';
import ComingsoonVideo from '../components/comingsoon/comingsoon_video';
const routes = [
  {
    path: "",
    redirect: "/dashboard/default",
  },
  {
    path: '/',
    component: BodyView,
   
    children: [
    {
      path: '',
      name: 'defaultRoot',
      component: IndexDefault,
      meta: {
        title: 'Viho - Premium Vue Admin Template',
      }
    },
    
    ]
  },
  {
    path: '/dashboard',
    component: BodyView,
   
    children: [
    {
      path: 'default',
      name: 'defaultIndex',
      component: IndexDefault,
      meta: {
        title: 'Dashboard | Viho - Premium Vue Admin Template',
      }
    },
    {
      path: 'indexEcommerce',
      name:'E-Commerce',
      component: indexEcommerce,
      meta: {
        title: 'Dashboard Ecommerce | Viho - Premium Vue Admin Template',
      }
    },
    
    ]
  },
  {
    path:'/widgets',
    component:BodyView,
    children:[
      {
        path:'indexGeneral',
        name:'general',
        component:indexGeneral,
        meta: {
          title: 'Widgets General | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:'indexchart',
        name:'charts',
        component:indexChart ,
        meta: {
          title: 'Widgets Chart | Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:'/project',
    component:BodyView,
    children:[
      {
        path:'indexList',
        name:'projectlist',
        component:indexList,
        meta: {
          title: 'Project List| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:'createProject',
        name:'createProject',
        component:createProject,
        meta: {
          title: 'Create Project| Viho - Premium Vue Admin Template',
        }
      }
    ]
  },

  {
    path: '/error-page1',
    name: 'errorPage1',
    component: errorPage1,
    meta: {
      title: 'Error Page1| Viho - Premium Vue Admin Template',
    }
  },
  {
    path: '/error-page2',
    name: 'errorPage2',
    component: errorPage2,
    meta: {
      title: 'Error Page2| Viho - Premium Vue Admin Template',
    }
  },
  {
    path: '/error-page3',
    name: 'errorPage3',
    component: errorPage3,
    meta: {
      title: 'Error Page3| Viho - Premium Vue Admin Template',
    }
  },
  {
    path: '/error-page4',
    name: 'errorPage4',
    component: errorPage4,
    meta: {
      title: 'Error Page4| Viho - Premium Vue Admin Template',
    }
  },
  
  
  {
    path: '/cards',
    component: BodyView,
    children: [
      {
        path: 'basicView',
        name: 'BasicView',
        component: BasicView,
        meta: {
          title: 'BootstrapStyling Card | Viho - Premium Admin Template',
        }
      },
      {
        path: 'draggableView',
        name: 'DraggableView',
        component: DraggableView,
        meta: {
          title: 'Draggable Card | Viho - Premium Admin Template',
        }
      },
      {
        path:"creative",
        name:"creative",
        component:creative_card,
        meta: {
          title: 'Creative card | Viho - Premium Admin Template',
        }
      },
      {
        path:"tabbed",
        name:"tabbed",
        component:tabbed_card,
        meta: {
          title: 'Tabbed Card | Viho - Premium Admin Template',
        }
      }
    ]
  },
  {
    path:'/comingsoon/comingsoon-image',
    name:'ComingsoonImage',
    component:ComingsoonImage,
    meta: {
        title: 'ComingsoonImage | Viho - Premium Admin Template',
      }
  },
  {
    path:'/comingsoon/comingsoon-simple',
    name:'ComingsoonSimple',
    component:ComingsoonSimple,
    meta: {
        title: 'ComingsoonSimple | Viho - Premium Admin Template',
      }
  },
  {
    path:'/comingsoon/comingsoon-video',
    name:'ComingsoonVideo',
    component:ComingsoonVideo,
    meta: {
        title: 'ComingsoonVideo | Viho - Premium Admin Template',
      }
  },
  {
    path:"/maintenance",
    name:"maintenance",
    component:maintenance,
    meta: {
      title: 'Maintenance | Viho - Premium Vue Admin Template',
    }
  },
  {
    path:"/knowledgebase",
    component:BodyView,
    children:[
      {
        path:"",
       name:"knowledgebase",
       component:knowledgebase,
       meta: {
        title: 'Knowledgebase| Viho - Premium Vue Admin Template',
      }
      }
    ]
  },


    
  {
    path: '/app',
    component: BodyView,
    
    children: [
      {
        path: 'file_manager',
        name: 'defaultView3',
        component: file_manager,
        meta: {
          title: 'File Manager| Viho - Premium Vue Admin Template',
        }
      },
    {
      path: 'kanbanBoard',
      name: 'defaultView2',
      component: kanbanBoard,
      meta: {
        title: 'kanban Board| Viho - Premium Vue Admin Template',
      }
    },
    {
      path:"email",
      name:"email",
      component:indexEmail,
      meta: {
        title: 'Email | Viho - Premium Vue Admin Template',
      }
    },
    {
      path:"chat",
      name:"chat",
      component:indexChat,
      meta: {
        title: 'Chat App| Viho - Premium Vue Admin Template',
      }
    },
    {
      path:"videochat",
      name:"videoChat",
      component:videoChat,
      meta: {
        title: 'Video Chat| Viho - Premium Vue Admin Template',
      }
    },
    {
      path:"bookmark",
      name:"Bookmark",
      component:bookMark,
      meta: {
        title: 'Bookmark| Viho - Premium Vue Admin Template',
      }
    },
    {
        path:"todo",
        name:"todo",
        component:todoIndex,
        meta: {
          title: 'To Do| Viho - Premium Vue Admin Template',
        }
      },

      {
        path:"calendar",
        name:"calender",
        component:calender,
        meta: {
          title: 'Calendar| Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/pages",
    component:BodyView,
    children:[
      {
        path:"social-app",
        name:"socialapp",
        component:socialApp,
        meta: {
          title: 'Social App| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"search",
        name:"search",
        component:serchIndex,
        meta: {
          title: 'Search| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"sample-page",
        name:"samplepage",
        component:sample_page,
        meta: {
          title: 'Simple Page| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"faq",
        name:"faqindex",
        component:faqindex,
        meta: {
          title: 'Faq| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"support",
        name:"support",
        component:support,
        meta: {
          title: 'Support| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"pricing",
        name:"pricing",
        component:pricing,
        meta: {
          title: 'Pricing| Viho - Premium Vue Admin Template',
        }
      }
     
    ]
  },

  {
    path: '/ecommerce',
    component: BodyView,
    
    children: [
      {
        path: 'indexProduct',
        name: 'ecommerce',
        component: indexProduct,
        meta: {
          title: 'Product | Viho - Premium Vue Admin Template',
        }
      },
      {
        path: 'cartView',
        name: 'chart',
        component: cartView,
        meta: {
          title: 'Cart | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:'details/:id',
        name:'prouctDetails',
        component:productDetails,
        meta: {
          title: 'Product Page | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"checkOut",
        name:"checkout",
        component:checkOut,
        meta: {
          title: 'Checkout | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"wishList",
        name:"wisthlist",
        component:wishList,
        meta: {
          title: 'Wish List | Viho - Premium Vue Admin Template',
        }
      },
      {
       path:"invoiceView",
       name:"invoice",
       component:invoiceView,
       meta: {
        title: 'Invoice | Viho - Premium Vue Admin Template',
      }
      },
      {
        path:"payment/details",
        name:"payment",
        component:PaymentDetails,
        meta: {
          title: 'Payment Detail | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"order/history",
        name:"orderHistory",
        component:orderHistory,
        meta: {
          title: 'Order History | Viho - Premium Vue Admin Template',
        }

      },
      {
        path:"add-product",
        name:"addproduct",
        component:indexAdd,
        meta: {
          title: 'Add Product | Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/users",
    component:BodyView,
    children:[
      {
        path:"/users/profile",
        name:"userProfile",
        component:userProfile,
        meta: {
          title: 'User Profile| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"/users/edit",
        name:"edit",
        component:editProfile,
        meta: {
          title: 'User Edit| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"/users/cards",
        name:"userCard",
        component:cardIndex,
        meta: {
          title: 'User Cards| Viho - Premium Vue Admin Template',
        }
      }
    ]

  },
  {
    path:"/form",
    component:BodyView,
    children:[
      {
        path:"validation",
        name:"formValidation",
        component:formValidation,
        meta: {
          title: 'Form Controls Form Validation| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"inputs",
        name:"basicInput",
        component:base_Input,
        meta: {
          title: 'Form Controls Base Input| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"checkbox-radio",
        name:"checkbox & radio",
        component:checkbox_radio,
        meta: {
          title: 'Form Controls Checkbox & Radio| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"input-groups",
        name:"input Groups",
        component:input_groups,
        meta: {
          title: 'Form Controls Input Groups| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"mega-options",
        name:"megaOptions",
        component:megaOptions,
        meta: {
          title: 'Form Controls Mega Options| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"select2",
        name:"select2",
        component:select2,
        meta: {
          title: 'Form Widgets Select| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"switch",
        name:"switch",
        component:switch_From,
        meta: {
          title: 'Form Widgets Switch| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"touchspin",
        name:"touchspin",
        component:touchspin_Form,
        meta: {
          title: 'Form Widgets Touchspin| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"typeahead",
        name:"typeahead",
        component:typeahead_Form,
        meta: {
          title: 'Form Widgets Typeahead| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"clipboard",
        name:"clipboard",
        component:clipboard_Form,
        meta: {
          title: 'Form Widgets Clipboard| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"wizard",
        name:"formwizard",
        component:form_wizard,
        meta: {
          title: 'Form Layout Form Wizard 1| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"datepicker",
        name:"datepicker",
        component:datepicker,
        meta: {
          title: ' Datepicker| Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/table",
    component:BodyView,
    children:[
      {
        path:"basic",
        name:"basic1",
        component:basic_tables,
        meta: {
          title: 'Table Bootstrap Table| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"sizing",
        name:"sizing",
        component:sizing_tables,
        meta: {
          title: 'Sizing Table| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"border",
        name:"border",
        component:border_Tables,
        meta: {
          title: 'Border Table| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"styling",
        name:"styling",
        component:styling_table,
        meta: {
          title: 'Style Table| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"/datatable-basic",
        name:"datatable",
        component:basic_Init,
        meta: {
          title: 'Table Basic Init| Viho - Premium Vue Admin Template',
        }
      },
   
    ]
  },
  {
    path:"/uikits",
    component:BodyView,
    children:[
      {
        path:"typography",
        name:"typography",
        component:typograPhy,
        meta: {
          title: 'Uikits Typography | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"avatars",
        name:"avatars",
        component:avatars_Uikits,
        meta: {
          title: 'Uikits Avatars | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"helper-classes",
        name:"helper",
        component:helper_classes,
        meta: {
          title: 'Uikits Helper Classes | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"grid",
        name:"grid",
        component:grid_Uikits,
        meta: {
          title: 'Uikits Grid | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"tag-pills",
        name:"tag&pills",
        component:tag_pills,
        meta: {
          title: 'Uikits Tag & Pills | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"progress-bar",
        name:"progressbar",
        component:progress_bar,
        meta: {
          title: 'Uikits Progressbar | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"modal",
        name:"modal",
        component:modal_Uikit,
        meta: {
          title: 'Uikits Modal | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"alert",
        name:"alert",
        component:alert_Uikit,
        meta: {
          title: 'Uikits Alert | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"popover",
        name:"popover",
        component:popover_Uikit,
        meta: {
          title: 'Uikits popover | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"tooltip",
        name:"tooltip",
        component:tooltip_uikits,
        meta: {
          title: 'Uikits Tooltip | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"loader",
        name:"spinners",
        component:spinners_Uikit,
        meta: {
          title: 'Uikits Spinners | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"accordion",
        name:"accordion",
        component:accordion_Uikit,
        meta: {
          title: 'Uikits Accordion | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"box-shadow",
        name:"boxshadow",
        component:box_shadow,
        meta: {
          title: 'Uikits Shadow | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"lists",
        name:"list",
        component:list_Uikit,
        meta: {
          title: 'Uikits Lists | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"dropdown",
        name:"dropdown",
        component:dropdown_Uikit,
        meta: {
          title: 'Uikits Dropdown | Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/advance",
    component:BodyView,
    children:[
      {
        path:"scrollable",
        name:"scrollable",
        component:scrollable_advance,
        meta: {
          title: 'Bonus UI  Scrollable | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"pagination",
        name:"pagination",
        component:pagination_advance,
        meta: {
          title: 'Bonus UI  Pagenation | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"sweetalert",
        name:"sweetAlert",
        component:sweetAlert,
        meta: {
          title: 'Bonus UI  SweetAlert | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"ribbons",
        name:"ribbons",
        component:ribbons_advance,
        meta: {
          title: 'Bonus UI  Ribbons | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"breadcrumb",
        name:"breadCrumb",
        component:breadCrumb,
        meta: {
          title: 'Bonus UI  Breadcrumb | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"image-cropper",
        name:"cropper",
        component:cropper_advance,
        meta: {
          title: 'Bonus UI  ImageCropper | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"toastr",
        name:"toaster",
        component:toaster_advance,
        meta: {
          title: 'Bonus UI  Toaster | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"tour",
        name:"tour",
        component:tour_advance,
        meta: {
          title: 'Bonus UI  Tour | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"rating",
        name:"rating",
        component:rating_advance,
        meta: {
          title: 'Bonus UI  Rating | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"upload",
        name:"upload",
        component:upload_advance,
        meta: {
          title: 'Bonus UI  Upload | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"sticky",
        name:"sticky",
        component:sticky_advance,
        meta: {
          title: 'Bonus UI  Sticky | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"range-slider",
        name:"rangeslider",
        component:range_advance,
        meta: {
          title: 'Bonus UI  Range | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"dragdrop",
        name:"dragdrop",
        component:dragdrop,
        meta: {
          title: 'Bonus UI  Dragdrop | Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/icons",
    component:BodyView,
    children:[
      {
        path:"flag",
        name:"flag",
        component:flag_Icon,
        meta: {
          title: 'Icons Flag | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"feather_icon",
        name:"feather",
        component:feather_icon,
        meta: {
          title: 'Icons Feather | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"whether",
        name:"weather",
        component:weather_icon,
        meta: {
          title: 'Icons Whether | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"themify",
        name:"themify",
        component:themify_icon,
        meta: {
          title: 'Icons Themify | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"fontawesome",
        name:"FontAwesome",
        component:font_Awesome,
        meta: {
          title: 'Icons Fontawesome | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"ico",
        name:"icon",
        component:icon_Icon,
        meta: {
          title: 'Icons Icoicon | Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/buttons",
    component:BodyView,
    children:[
      {
        path:"default",
        name:"default",
        component:default_button,
        meta: {
          title: 'Default Button | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"flat",
        name:"flat",
        component:flat_button,
        meta: {
          title: 'Flat Button | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"edge",
        name:"edge",
        component:edge_button,
        meta: {
          title: 'Edge Button | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"raised",
        name:"raised",
        component:raised_button,
        meta: {
          title: 'Raised Button | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"group",
        name:"button_group",
        component:button_group,
        meta: {
          title: 'Group Button | Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/chart",
    component:BodyView,
    children:[
      {
        path:"google",
        name:"googlechart",
        component:google_chart,
        meta: {
          title: 'Chart Google Chart| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"apexChart",
        name:"apexchart",
        component:apex_chart,
        meta: {
          title: 'Chart Apex Chart| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"chartist",
        name:"chartist",
        component:chartist_chart,
        meta: {
          title: 'Chart Chartist Chart| Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/gallery",
    component:BodyView,
    children:[
      {
        path:"grid-gallery",
        name:"gridgallery",
        component:grid_gallery,
        meta: {
          title: 'Grid Gallery | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"gallery-desc",
        name:"griddesc",
        component:grid_desc,
        meta: {
          title: 'Grid Gallery With Desc | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"hover-effect",
        name:"hovergallery",
        component:hover_gallery,
        meta: {
          title: 'Hover Gallery| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"gallery-masonary",
        name:"masonry-gallery",
        component:masonry_gallery,
        meta: {
          title: 'Masonry | Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"gallery-masonary-desc",
        name:"masonary-desc",
        component:masonary_desc,
        meta: {
          title: 'Masonry Desc| Viho - Premium Vue Admin Template',
        }
      }

    ]
  },
 
  {
    path:"/auth",
    component: auth,
    children:[
      {
        path:"login",
        name:"login",
        component:login,
        meta: {
          title: 'Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/authentication/login/simple",
    name:"login-simple",
    component:login_simple,
    meta: {
      title: 'Login Simple| Viho - Premium Vue Admin Template',
    }
  },
   {
     path:"/authentication/login/one",
     name:"login-image",
     component:login_image,
     meta: {
      title: 'Login Image| Viho - Premium Vue Admin Template',
    }
   },
   {
     path:"/authentication/login/two",
     name:"login-image2",
     component:login_image2,
     meta: {
      title: 'Login Image Two| Viho - Premium Vue Admin Template',
    }
   },
   {
    path:"/authentication/register/image",
    name:"ragister-image",
    component:ragister_image,
    meta: {
      title: 'Register Image| Viho - Premium Vue Admin Template',
    }
   },
   {
     path:"/authentication/register/image2",
     name:"ragister-image2",
     component:ragister_image2,
     meta: {
      title: 'Register Image Two| Viho - Premium Vue Admin Template',
    }
   },
  {
    path:"/auth/register",
    name:"ragister_simple",
    component:ragister_simple,
    meta: {
      title: 'Register| Viho - Premium Vue Admin Template',
    }
  },
  {
    path:"/authentication/login/validate",
    name:"login-validation",
    component:login_with_validation,
    meta: {
      title: 'Login Validation| Viho - Premium Vue Admin Template',
    }
  },
  {
    path:"/authentication/unlockuser",
    name:"unlock",
    component:unlock,
    meta: {
      title: 'Unlock User| Viho - Premium Vue Admin Template',
    }
  },
  {
    path:"/authentication/forgetpassword",
    name:"forget-password",
    component:forget_password,
    meta: {
      title: 'Forget Password| Viho - Premium Vue Admin Template',
    }
  },
  {
    path:"/authentication/resetpassword",
    name:"reset-password",
    component:reset_password,
    meta: {
      title: 'Reset Password| Viho - Premium Vue Admin Template',
    }
  },
  {
    path:"/timeline",
    component:BodyView,
    children:[
      {
        path:"one",
        name:"timeline",
        component:timeline,
        meta: {
          title: 'Timeline | Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/maps",
    component:BodyView,
    children:[
      {
        path:"vue-google-maps",
        name:"google-map",
        component:google_map,
        meta: {
          title: 'Google Map| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"vue-leaflet-maps",
        name:"vueleaflet",
        component:vue_leaflet,
        meta: {
          title: 'Leaflet Map| Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/editor",
    component:BodyView,
    children:[
      {
        path:"ck-editor",
        name:"ckediter",
        component:ckediter,
        meta: {
          title: 'Ck Editor| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"simple-editor",
        name:"simple-editer",
        component:simple_editer,
        meta: {
          title: 'Simple Editor| Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/blog",
    component:BodyView,
    children:[
      {
        path:"details",
        name:"blog-detail",
        component:blog_detail,
        meta: {
          title: 'Blog Details| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"single",
        name:"blog_single",
        component:blog_single,
        meta: {
          title: 'Blog Single| Viho - Premium Vue Admin Template',
        }
      }
    ]
  },
  {
    path:"/job",
    component:BodyView,
    children:[
      {
        path:"list",
        name:"listview",
        component:job_list,
        meta: {
          title: 'Job List| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"/job/details/:id",
        name:"jobdetails",
        component:job_details,
        props:true,
        meta: {
          title: 'Job Details| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"apply/:id",
        name:"jobapply",
        component:job_apply,
        meta: {
          title: 'Job Apply| Viho - Premium Vue Admin Template',
        }
      }
    ]  
  },
  {
    path:"/learning",
    component:BodyView,
    children:[
      {
        path:"list",
        name:"learninglist",
        component:learning_list,
        meta: {
          title: 'Learning List| Viho - Premium Vue Admin Template',
        }
      },
      {
        path:"details/:id",
        name:"coursedetailed",
        component:coursedetailed,
        meta: {
          title: 'course Detailed| Viho - Premium Vue Admin Template',
        }
      }
    ]
  }

]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
  linkExactActiveClass: 'active',
  
})

router.beforeEach((to, from, next) => {
  if (typeof (to.meta.title) === 'string') {
    document.title = to.meta.title;
  }
    const  path = ['/auth/login','/auth/register'];
    if (path.includes(to.path) || localStorage.getItem('user')  ){
      return next();
    }
    next('/auth/login');
  });

export default router
