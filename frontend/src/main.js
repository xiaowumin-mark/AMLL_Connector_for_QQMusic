import { createApp } from 'vue'
import App from './App.vue'
import "./style/base.css"
import 'mdui/mdui.css';
import 'mdui';
import router from './router'
import { Window } from "@wailsio/runtime";
const app = createApp(App)
app.use(router)
router.isReady().then(() => app.mount('#app'))

import { GetConfig } from "../bindings/AMLL_Connector_for_QQMusic/greetservice"
GetConfig().then(data => {
    data = JSON.parse(data)
    document.documentElement.classList.add("mdui-theme-dark");
    //const themeMedia = window.matchMedia("(prefers-color-scheme: light)");
    //// 监听系统主题变化
    //if (data.theme == "auto") {
    //    if (themeMedia.matches) {
    //        document.documentElement.classList.add("mdui-theme-light");
    //        Window.
    //    } else {
    //        document.documentElement.classList.add("mdui-theme-dark");
    //    }
    //} else {
    //    document.documentElement.classList.add("mdui-theme-" + data.theme)
    //}
    //themeMedia.addEventListener("change", (e) => {
    //    GetConfig().then(data => {
    //        data = JSON.parse(data)
    //        if (data.theme == "auto") {
    //            if (e.matches) {
    //                document.documentElement.classList.add("mdui-theme-layout-light");
    //            } else {
    //                document.documentElement.classList.add("mdui-theme-layout-dark");
    //            }
    //        } else {
    //            document.documentElement.classList.add("mdui-theme-" + data.theme)
    //        }
    //    })
    //
    //});
})