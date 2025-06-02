<script setup>
import { ref, onMounted, onUnmounted, watch, handleError } from 'vue'
import { Events } from '@wailsio/runtime'
import { formatMilliseconds } from '../scripts/time'
import { GetNowMusicInfo } from '../../bindings/AMLL_Connector_for_QQMusic/greetservice'
const musicInfo = ref({
    "title": "---",
    "artist": "---",
    "position": 0,
    "duration": 0,
    "album": "---",
    "status": "Paused",
    "mid": "---",
    "id": "---",
    "aid": "---"
})
const log = ref([])
const logContent = ref(null)
const logLength = ref(0)
const autoScroll = ref(true) // 新增：控制是否自动滚动
const LyricStr = ref("")
const pic = ref("")

// 新增：处理滚动事件
const handleScroll = (e) => {
    if (e.target.scrollTop + e.target.clientHeight >= e.target.scrollHeight) {
        autoScroll.value = true;
    } else {
        autoScroll.value = false;
    }
}

watch(log, () => {
    if (autoScroll.value) {
        setTimeout(() => {
            logContent.value.scrollTop = logContent.value.scrollHeight;
        }, 0);
    }
}, {
    deep: true
})

onMounted(() => {

    GetNowMusicInfo().then(res => {
        console.log(res);

        musicInfo.value.title = res.title
        musicInfo.value.artist = res.artist
        musicInfo.value.album = res.album
        musicInfo.value.duration = res.duration
        musicInfo.value.status = "Paused"
        musicInfo.value.mid = res.mid
        musicInfo.value.id = res.id
        musicInfo.value.aid = res.albumId

        LyricStr.value = res.lyricRaw

        pic.value = "data:image/png;base64,"+res.pic
    })

    Events.On("set_music_info_more", (e) => {
        pic.value = "data:image/png;base64,"+e.data[0].pic;
    })

    Events.On("amll_play_progress", (e) => {
        musicInfo.value.position = e.data[0]
    })

    Events.On("amll_music_info", (e) => {
        musicInfo.value.title = e.data[0].title
        musicInfo.value.artist = e.data[0].artist
        musicInfo.value.album = e.data[0].album
        musicInfo.value.duration = e.data[0].duration
        musicInfo.value.status = e.data[0].status
    })

    Events.On("logger", (e) => {
        logLength.value += 1
        if (log.value.length > 100) {
            log.value.shift()
        }
        log.value.push({
            id: logLength.value,
            data: e.data
        })
    })

    Events.On("set_music_info_is_id", (e) => {
        console.log(e.data);
        musicInfo.value.mid = e.data[0].mid
        musicInfo.value.id = e.data[0].id
        musicInfo.value.aid = e.data[0].aid

    })
    Events.On("set_lyric", (e) => {
        console.log(e.data);
        LyricStr.value = e.data[0]
    })

})

onUnmounted(() => {
    Events.Off("set_music_info_more")
    Events.Off("amll_play_progress")
    Events.Off("amll_music_info")
    Events.Off("logger")
    Events.Off("set_music_info_is_id")
    Events.Off("set_lyric")
})


</script>

<template>
    <main>
        <div class="row">
            <mdui-card class="item">
                <div class="title">详细信息</div>
                <div class="content" style="line-height: 22px;">
                    <p class="can_select">歌曲名称：{{ musicInfo.title }}</p>
                    <p class="can_select">歌手：{{ musicInfo.artist }}</p>
                    <p class="can_select">专辑：{{ musicInfo.album }}</p>
                    <p>播放进度：{{ formatMilliseconds(musicInfo.position) }}</p>
                    <p class="can_select">歌曲时长：{{ formatMilliseconds(musicInfo.duration) }}</p>
                    <p>播放状态：{{ musicInfo.status }}</p>
                    <p>歌曲MID：<span class="can_select">{{ musicInfo.mid }}</span></p>
                    <p>歌曲ID：<span class="can_select">{{ musicInfo.id }}</span></p>
                    <p>专辑ID：<span class="can_select">{{ musicInfo.aid }}</span></p>
                </div>
            </mdui-card>
            <mdui-card class="item">
                <div class="title">歌词</div>
                <div class="content can_select" style="font-family: 'Consolas';" v-html="LyricStr">
                </div>
            </mdui-card>

            <mdui-card class="item">
                <div class="title">封面</div>
                <div class="content pic" :style="{backgroundImage: 'url(' + pic + ')'}">
                    <!--<img style="width: 100%;height: 100%;object-fit: contain;transition: all 0.3s;"
                        :src="'data:image/png;base64,' + pic" alt="">-->
                </div>
            </mdui-card>
        </div>
        <div class="logger">
            <mdui-card>
                <div class="title">日志({{ logLength }})</div>
                <div class="content" ref="logContent" @scroll="handleScroll"
                    style="padding-left: 10px;padding-right: 10px;box-sizing: border-box;font-family: 'Consolas';line-height: 22px;">
                    <p v-for="(item, index) in log" class="log-item" :key="item.id">
                        <span v-for="item in item.data">{{ item }}</span>
                    </p>
                </div>
            </mdui-card>
        </div>
    </main>
</template>

<style lang="css" scoped>
main {
    width: 100%;
    height: 100%;
}

.row {
    width: 100%;
    height: 250px;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding-left: 10px;
    padding-right: 10px;
    box-sizing: border-box;
    align-items: center;
    gap: 10px;

    .item {
        width: calc(50% - 5px);
        height: 100%;
    }
}

.logger {
    width: 100%;
    height: calc(100% - 250px);
    padding: 10px;
    box-sizing: border-box;

    mdui-card {
        width: 100%;
        height: 100%;
    }
}

.title {
    font-size: 15px;
    width: 100%;
    height: 30px;
    text-align: center;
    line-height: 30px;
}

.content {
    width: 100%;
    height: calc(100% - 30px);
    overflow-y: auto;
    overflow-x: hidden;
    padding-left: 10px;
    padding-right: 10px;
    box-sizing: border-box;

    p {
        margin: 0;
    }

    &::-webkit-scrollbar {
        width: 5px;
        background-color: rgba(255, 255, 255, 0.05);
        border-radius: 2.5px;
    }

    &::-webkit-scrollbar-thumb {
        background-color: rgba(255, 255, 255, 0.2);
        border-radius: 2.5px;
    }

    &::-webkit-scrollbar-thumb:hover {
        background-color: rgba(255, 255, 255, 0.3);
    }
}


.log-item {
    white-space: pre-wrap;
    opacity: 0;
    transform-origin: left center;
    animation: fade-in 1s forwards;
}

@keyframes fade-in {
    0% {
        opacity: 0;
        filter: blur(2px);
        transform: scaleX(0.95);
    }

    100% {
        opacity: 1;
        filter: blur(0);
        transform: scaleX(1);
    }
}

.log-item span:first-child {
    color: greenyellow;
    margin-right: 5px;
    transition: color 0.2s ease-in-out;

    &:hover {
        color: green;
    }
}

mdui-card {
    background-color: rgba(29, 27, 32, 0.7);
}

.pic {
  background-size: contain; /* 关键属性 */
  background-repeat: no-repeat;
  background-position: center;
  transition: all 0.3s;
}
</style>