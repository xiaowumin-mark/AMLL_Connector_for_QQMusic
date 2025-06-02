<script setup>
import { ref, onMounted, watch } from 'vue'
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiFolderOutline } from '@mdi/js'
import { Dialogs, Events, Window } from "@wailsio/runtime";
import { GetConfig, SetAllConfig, ChoseDir } from "../../bindings/AMLL_Connector_for_QQMusic/greetservice"
const config = ref({})

onMounted(() => {
    GetConfig().then(res => {
        res = JSON.parse(res)
        config.value = res

    })
})
watch(() => config.value, () => {
    let data = Object.assign({}, config.value)
    SetAllConfig(JSON.stringify(data)).then(res => {
        res = JSON.parse(res)
        console.log(res);
    })
}, {
    deep: true
})

function ChangeDir(type) {
    //Dialogs.OpenFile({
    //    CanChooseDirectories: true,
    //    CanChooseFiles: false,
    //}).then(res => {
    //    console.log(res);
    //    
    //    //config.value[type] = res
    //}).catch(err => {
    //    console.log(err);
    //    
    //})
    ChoseDir(config.value[type]).then(res => {
        if (res) {
            console.log(res);
            config.value[type] = res
        }
    }).catch(err => {
        console.log(err);

    })
}
</script>
<template>
    <span
        style="position: fixed;bottom: 5px;width: 100vw;left: 0;text-align: center;pointer-events: none;opacity: 0.5;">APP版本:{{
            config.app_version }}</span>
    <mdui-list>
        <mdui-list-subheader>设置</mdui-list-subheader>
        <!--<mdui-list-item rounded>主题
            <mdui-radio-group :value="config.theme" slot="end-icon" @change="config.theme = $event.target.value">
                <mdui-radio value="auto">自动</mdui-radio>
                <mdui-radio value="light">浅色</mdui-radio>
                <mdui-radio value="dark">深色</mdui-radio>
            </mdui-radio-group>
        </mdui-list-item>-->
        <mdui-list-item rounded>自动连接AMLL
            <mdui-switch :checked="config.auto_connec" slot="end-icon"
                @change="config.auto_connec = $event.target.checked"></mdui-switch>
        </mdui-list-item>
        <mdui-list-item rounded>AMLL WS 主机
            <mdui-text-field variant="outlined" slot="end-icon" :value="config.auto_connect_address"
                @change="config.auto_connect_address = $event.target.value" style="height: 40px;">

            </mdui-text-field>
        </mdui-list-item>
        <mdui-list-item rounded>AMLL WS 端口
            <mdui-text-field variant="outlined" type="number" step="1" slot="end-icon" :value="config.auto_connect_port"
                @change="config.auto_connect_port = +$event.target.value" style="height: 40px;">

            </mdui-text-field>
        </mdui-list-item>
        <mdui-list-item rounded>歌词保存目录
            <mdui-text-field variant="outlined" type="text" slot="end-icon" :value="config.lyrics_path"
                @change="config.lyrics_path = $event.target.value" style="height: 40px;">

                <mdui-button-icon slot="end-icon" @click="ChangeDir('lyrics_path')">
                    <svg-icon type="mdi" :path="mdiFolderOutline">

                    </svg-icon>
                </mdui-button-icon>
            </mdui-text-field>
        </mdui-list-item>
        <mdui-list-item rounded>专辑图片保存目录
            <mdui-text-field variant="outlined" type="text" slot="end-icon" :value="config.album_art_path"
                @change="config.album_art_path = $event.target.value" style="height: 40px;">
                <mdui-button-icon slot="end-icon" @click="ChangeDir('album_art_path')">
                    <svg-icon type="mdi" :path="mdiFolderOutline"></svg-icon>
                </mdui-button-icon>
            </mdui-text-field>
        </mdui-list-item>
        <mdui-list-item rounded>在有其他软件播放时降低音量
            <mdui-switch :checked="config.auto_change_volume" slot="end-icon"
                @change="config.auto_change_volume = $event.target.checked"></mdui-switch>
        </mdui-list-item>
    </mdui-list>
</template>