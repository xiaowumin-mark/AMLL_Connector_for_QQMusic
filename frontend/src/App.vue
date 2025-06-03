<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { Window, Events, } from "@wailsio/runtime";
import { RouterLink, RouterView, useRouter, useRoute } from 'vue-router'
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiCogOutline, mdiCog, mdiDesktopTower, mdiLinkCircleOutline, mdiLinkCircle, mdiFormatListText, mdiTextBoxOutline, mdiTextBox } from '@mdi/js'
import { TriggerSnapLayout, ConnectAmll, DisconnectAmll, ShowLyricWindow, HideLyricWindow, IsLyricWindowShow } from "../bindings/AMLL_Connector_for_QQMusic/greetservice"
document.body.style.visibility = "visible";
const route = useRoute()
const router = useRouter()

const routeName = ref(route.name);
router.afterEach((to, from) => {
  routeName.value = to.name;
});

let timer = null;

const isCon = ref(false);
const isShowLyricWindow = ref(false);
IsLyricWindowShow().then((res) => {
  isShowLyricWindow.value = res;
});
watch(isShowLyricWindow, (val) => {
  if (val) {
    ShowLyricWindow();
  } else {
    HideLyricWindow();
  }
})
// 调用后端方法
function onMaximizeHover(event) {
  // 调用 Go 后端的方法
  timer = setTimeout(async () => {
    await TriggerSnapLayout();
  }, 700);
}
function onMaximizeleave(event) {
  if (timer) {
    clearTimeout(timer);
  }
}


function TSNpush(d) {
  if (document.startViewTransition) {
    document.startViewTransition(() => {
      router.push(d);
    });
  } else {
    router.push(d);
  }
};

function Con(e) {
  console.log("amll_ws_state");

  if (isCon.value) {
    DisconnectAmll().then((e) => {
      console.log(e);
    }).catch((e) => {
      console.error(e);
    });
  } else {
    ConnectAmll().then((e) => {
      console.log(e);
    }).catch((e) => {
      console.error(e);
    });
  }
}

onMounted(() => {
  Events.On("amll_ws_state", (e) => {
    console.log(e.data)
    isCon.value = e.data[0]
  })
  Events.On("window_will_update_visibility", (e) => {
    console.log(e.data);

    //isShowLyricWindow.value = e.data[0]
    if (e.data[0].window=="main"){
      if (['home', 'setting', 'smtcs'].includes(routeName.value)) {
        document.body.style.visibility = e.data[0].visible ? "visible" : "hidden";
      }
    }else if (e.data[0].window=="lyric") {
      isShowLyricWindow.value = e.data[0].visible;
      if (routeName.value == "lyrics") {
        document.body.style.visibility = e.data[0].visible ? "visible" : "hidden";
      }
    }
  })
  
})
onUnmounted(() => {
  Events.Off("amll_ws_state");
  Events.Off("window_will_update_visibility");
})
</script>

<template>

  <div class="main" :style="{ 'padding-left': route.meta.show_nav ? '80px' : '0px' }">
    <mdui-navigation-rail contained alignment="center" class="nobg" :value="routeName" v-show="route.meta.show_nav"
      style="border-right: 1px solid rgba(255,255,255,0.1);">
      <mdui-tooltip content="桌面歌词" slot="bottom">
        <mdui-button-icon selectable @change="isShowLyricWindow = $event.target.selected" :selected="isShowLyricWindow">
          <svg-icon type="mdi" :path="mdiTextBoxOutline"></svg-icon>
          <svg-icon type="mdi" :path="mdiTextBox" slot="selected-icon"></svg-icon>
        </mdui-button-icon>
      </mdui-tooltip>


      <mdui-navigation-rail-item value="home" @click="TSNpush({ name: 'home' })">仪表盘
        <svg-icon type="mdi" :path="mdiDesktopTower" slot="icon"></svg-icon>
        <svg-icon type="mdi" :path="mdiDesktopTower" slot="active-icon"></svg-icon>
      </mdui-navigation-rail-item>
      <mdui-navigation-rail-item value="smtcs" @click="TSNpush({ name: 'smtcs' })">SMTCS
        <svg-icon type="mdi" :path="mdiFormatListText" slot="icon"></svg-icon>
        <svg-icon type="mdi" :path="mdiFormatListText" slot="active-icon"></svg-icon>
      </mdui-navigation-rail-item>
      <mdui-navigation-rail-item value="setting" @click="TSNpush({ name: 'setting' })">设置
        <svg-icon type="mdi" :path="mdiCogOutline" slot="icon"></svg-icon>
        <svg-icon type="mdi" :path="mdiCog" slot="active-icon"></svg-icon>
      </mdui-navigation-rail-item>
    </mdui-navigation-rail>

    <div :style="{ 'width': `calc(100%)`, 'height': '100vh' }">
      <div class="titlebar" v-show="route.meta.show_top_bar">


        <div class="other">
          <mdui-button-icon @click="Con" :selected="isCon" title="连接AMLL">
            <svg-icon type="mdi" :path="mdiLinkCircleOutline"></svg-icon>
            <svg-icon type="mdi" :path="mdiLinkCircle" slot="selected-icon"></svg-icon>
          </mdui-button-icon>
          <div class="drag" @dblclick="Window.ToggleMaximise()">

          </div>
        </div>


        <div class="button" @click="Window.Minimise()">
          <svg width="10" height="2" viewBox="0 0 10 2" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M0.498047 1.39099C0.429688 1.39099 0.364583 1.37797 0.302734 1.35193C0.244141 1.32589 0.192057 1.29008 0.146484 1.24451C0.100911 1.19893 0.0651042 1.14685 0.0390625 1.08826C0.0130208 1.02641 0 0.961304 0 0.892944C0 0.824585 0.0130208 0.761108 0.0390625 0.702515C0.0651042 0.640666 0.100911 0.586955 0.146484 0.541382C0.192057 0.492554 0.244141 0.455119 0.302734 0.429077C0.364583 0.403035 0.429688 0.390015 0.498047 0.390015H9.50195C9.57031 0.390015 9.63379 0.403035 9.69238 0.429077C9.75423 0.455119 9.80794 0.492554 9.85352 0.541382C9.89909 0.586955 9.9349 0.640666 9.96094 0.702515C9.98698 0.761108 10 0.824585 10 0.892944C10 0.961304 9.98698 1.02641 9.96094 1.08826C9.9349 1.14685 9.89909 1.19893 9.85352 1.24451C9.80794 1.29008 9.75423 1.32589 9.69238 1.35193C9.63379 1.37797 9.57031 1.39099 9.50195 1.39099H0.498047Z"
              fill="currentColor" fill-opacity="0.8956"></path>
          </svg>
        </div>
        <div class="button" @mouseenter="onMaximizeHover" @mouseleave="onMaximizeleave"
          @click="Window.ToggleMaximise()">
          <svg width="10" height="11" viewBox="0 0 10 11" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M1.47461 10.391C1.2793 10.391 1.09212 10.3519 0.913086 10.2738C0.734049 10.1924 0.576172 10.085 0.439453 9.95154C0.30599 9.81482 0.198568 9.65694 0.117188 9.47791C0.0390625 9.29887 0 9.11169 0 8.91638V1.8656C0 1.67029 0.0390625 1.48311 0.117188 1.30408C0.198568 1.12504 0.30599 0.968791 0.439453 0.835327C0.576172 0.698608 0.734049 0.591187 0.913086 0.513062C1.09212 0.431681 1.2793 0.390991 1.47461 0.390991H8.52539C8.7207 0.390991 8.90788 0.431681 9.08691 0.513062C9.26595 0.591187 9.4222 0.698608 9.55566 0.835327C9.69238 0.968791 9.7998 1.12504 9.87793 1.30408C9.95931 1.48311 10 1.67029 10 1.8656V8.91638C10 9.11169 9.95931 9.29887 9.87793 9.47791C9.7998 9.65694 9.69238 9.81482 9.55566 9.95154C9.4222 10.085 9.26595 10.1924 9.08691 10.2738C8.90788 10.3519 8.7207 10.391 8.52539 10.391H1.47461ZM8.50098 9.39001C8.56934 9.39001 8.63281 9.37699 8.69141 9.35095C8.75326 9.32491 8.80697 9.2891 8.85254 9.24353C8.89811 9.19796 8.93392 9.14587 8.95996 9.08728C8.986 9.02543 8.99902 8.96033 8.99902 8.89197V1.89001C8.99902 1.82166 8.986 1.75818 8.95996 1.69958C8.93392 1.63774 8.89811 1.58403 8.85254 1.53845C8.80697 1.49288 8.75326 1.45707 8.69141 1.43103C8.63281 1.40499 8.56934 1.39197 8.50098 1.39197H1.49902C1.43066 1.39197 1.36556 1.40499 1.30371 1.43103C1.24512 1.45707 1.19303 1.49288 1.14746 1.53845C1.10189 1.58403 1.06608 1.63774 1.04004 1.69958C1.014 1.75818 1.00098 1.82166 1.00098 1.89001V8.89197C1.00098 8.96033 1.014 9.02543 1.04004 9.08728C1.06608 9.14587 1.10189 9.19796 1.14746 9.24353C1.19303 9.2891 1.24512 9.32491 1.30371 9.35095C1.36556 9.37699 1.43066 9.39001 1.49902 9.39001H8.50098Z"
              fill="currentColor" fill-opacity="0.8956"></path>
          </svg>
        </div>
        <div class="button close" @click="Window.Hide()">
          <svg width="10" height="11" viewBox="0 0 10 11" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M5 6.099L0.854492 10.2445C0.756836 10.3422 0.639648 10.391 0.50293 10.391C0.359701 10.391 0.239258 10.3438 0.141602 10.2494C0.0472005 10.1517 0 10.0313 0 9.88806C0 9.75134 0.0488281 9.63416 0.146484 9.5365L4.29199 5.39099L0.146484 1.24548C0.0488281 1.14783 0 1.02901 0 0.889038C0 0.820679 0.0130208 0.755575 0.0390625 0.693726C0.0651042 0.631877 0.100911 0.579793 0.146484 0.537476C0.192057 0.491903 0.245768 0.456095 0.307617 0.430054C0.369466 0.404012 0.43457 0.390991 0.50293 0.390991C0.639648 0.390991 0.756836 0.439819 0.854492 0.537476L5 4.68298L9.14551 0.537476C9.24316 0.439819 9.36198 0.390991 9.50195 0.390991C9.57031 0.390991 9.63379 0.404012 9.69238 0.430054C9.75423 0.456095 9.80794 0.491903 9.85352 0.537476C9.89909 0.583049 9.9349 0.636759 9.96094 0.698608C9.98698 0.757202 10 0.820679 10 0.889038C10 1.02901 9.95117 1.14783 9.85352 1.24548L5.70801 5.39099L9.85352 9.5365C9.95117 9.63416 10 9.75134 10 9.88806C10 9.95642 9.98698 10.0215 9.96094 10.0834C9.9349 10.1452 9.89909 10.1989 9.85352 10.2445C9.8112 10.2901 9.75911 10.3259 9.69727 10.3519C9.63542 10.378 9.57031 10.391 9.50195 10.391C9.36198 10.391 9.24316 10.3422 9.14551 10.2445L5 6.099Z"
              fill="currentColor" fill-opacity="0.8956"></path>
          </svg>
        </div>
      </div>
      <div class="main-app" :style="{ 'height': `calc(100% - ${route.meta.show_top_bar ? '42px' : '0px'})` }">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<style scoped>
.main {
  height: 100%;
}

.main-app {
  width: 100%;
  /*height: calc(100% - 32px);*/
}

.titlebar {
  width: 100%;
  height: 42px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  /* 让整个titlebar可以拖动 */
  z-index: 999;
}

.other {
  width: calc(100% - 42px * 3);
  height: 100%;
  padding-left: 10px;
  box-sizing: border-box;
  display: flex;
  justify-content: space-between;
}

.drag {
  width: calc(100% - 40px);
  height: 100%;
  --wails-draggable: drag;
}

.button {
  width: 50px;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  /* 按钮不可拖动 */
  cursor: pointer;
  color: white;
  font-size: 18px;
  transition: background 0.2s;
}

.button:hover {
  background: #5555557c;
}

.button.close:hover {
  background: #C42B1C;

}

.button.close::after {
  content: "";
  position: absolute;
  top: 0;
  right: 0;
  width: 5px;
  height: 5px;
  z-index: 999;
  cursor: nesw-resize;
}

html[class=mdui-theme-light] {
  .titlebar svg {
    filter: invert(1);
  }

  .close:hover svg {
    filter: invert(0);

  }
}
</style>
