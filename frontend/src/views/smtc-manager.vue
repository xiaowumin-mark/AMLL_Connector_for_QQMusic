<template>
    <main class="scroll_bar">
        <mdui-list>
            <mdui-list-subheader>SMTC管理器</mdui-list-subheader>
            <mdui-list-item v-for="item in smtcList" :headline="item.appId">
                <div slot="description">
                    <p>标题：{{ item.title }}</p>
                    <p>艺术家：{{ item.artist }}</p>
                    <p>专辑：{{ item.album }}</p>
                    <p>播放状态：{{ item.playbackStatus }}</p>
                    <p>歌曲总时长：{{ formatMilliseconds(item.timelineEndTimeMS) }}  ({{ item.timelineEndTimeMS }})</p>
                    <p>播放进度：{{ formatMilliseconds(item.timelinePositionMS) }}  ({{ item.timelinePositionMS }})</p>
                </div>

            </mdui-list-item>
        </mdui-list>
    </main>
</template>

<script setup>
import { reactive, onMounted, onUnmounted } from 'vue';
import { GetAllSmtc } from '../../bindings/AMLL_Connector_for_QQMusic/greetservice';
import { Events } from "@wailsio/runtime"
import{formatMilliseconds} from '../scripts/time'

const smtcList = reactive([])

onMounted(() => {
    GetAllSmtc().then(res => {
        for (let i = 0; i < res.length; i++) {
            smtcList.push(res[i])
        }
    })
    Events.On("smtc_added", (data) => {
        console.log("ad",data.data[0]);
        for (let i = 0; i < smtcList.length; i++) {
            if (smtcList[i].appId === data.data[0].appId) {
                smtcList.splice(i, 1)
                break
            }
        }
        smtcList.push(data.data[0])
    })
    Events.On("smtc_removed", (data) => {
        console.log("re",data.data[0]);
        
        for (let i = 0; i < smtcList.length; i++) {
            if (smtcList[i].appId === data.data[0].appId) {
                smtcList.splice(i, 1)
                break
            }
        }
    })
    Events.On("smtc_changed", (data) => {
        console.log("ch",data.data[0]);
        for (let i = 0; i < smtcList.length; i++) {
            if (smtcList[i].appId === data.data[0].appId) {
                smtcList[i] = data.data[0]
                console.log(smtcList[i]);
                
                break
            }
        }
    })

})
onUnmounted(() => {
    Events.Off("smtc_added")
    Events.Off("smtc_removed")
    Events.Off("smtc_changed")
})
</script>

<style scoped>
main {
    width: 100%;
    height: 100%;
    overflow-y: auto;
}
</style>