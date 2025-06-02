<template>
    <main  ref="mainc">
        <div class="m" ref="lyricsView" @mousemove="handleMouseMove" @mouseup="handleMouseUp" @mousedown="handleMouseDown">

        </div>
    </main>
</template>

<script setup>
import { Events, Window } from '@wailsio/runtime'
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { GetNowLyrics } from '../../bindings/AMLL_Connector_for_QQMusic/greetservice'
const mainc = ref(null)
document.body.style.visibility = 'hidden'
// 窗口居中
Window.Center()
const dbc = async () =>  {
    mainc.value.animate([
        {
            // 透明度
            opacity: 0,
            filter: 'blur(10px)',
        },
        {
            opacity: 0,
            filter: 'blur(10px)',
            transform: 'scale(0.7)',
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
            borderRadius: '999px',
            transformOrigin: 'center',
            offset: 0.5,
        },
        {
            opacity: 1,
            filter: 'blur(0px)',
            backgroundColor: 'rgba(0, 0, 0, 0)',
            transform: 'scale(1)',
            borderRadius: '10px',

        }
    ], {
        duration: 500,
    })
    setTimeout(() =>
        Window.ToggleMaximise(), 100)
}


window.addEventListener('dblclick', dbc)
let isMouseDown = ref(false)
// 新增状态用于记录拖动起始位置
const dragStartPos = ref({ x: 0, y: 0 })
const windowStartPos = ref({ x: 0, y: 0 })
// 在ref声明部分新增拖拽状态
const isDragging = ref(false)

// 修改后的鼠标事件处理函数
const handleMouseDown = (e) => {
    isMouseDown.value = true
    isDragging.value = false // 重置拖拽状态
    dragStartPos.value = { x: e.screenX, y: e.screenY }
    Window.Position().then(pos => {
        windowStartPos.value = { x: pos.x, y: pos.y }
    })
}

const handleMouseMove = (e) => {
    if (!isMouseDown.value) return

    // 计算偏移量
    const deltaX = e.screenX - dragStartPos.value.x
    const deltaY = e.screenY - dragStartPos.value.y

    // 计算移动距离（使用平方比较避免开方运算）
    const distanceSq = deltaX * deltaX + deltaY * deltaY

    // 当移动超过10px或已开始拖动时执行
    if (distanceSq > 100 || isDragging.value) {
        isDragging.value = true
        Window.SetPosition(
            windowStartPos.value.x + deltaX,
            windowStartPos.value.y + deltaY
        )
    }
}

// 新增mouseup处理函数
const handleMouseUp = () => {
    isMouseDown.value = false
    isDragging.value = false
}

const lyrics = ref([])
const previousIndex = ref([])
const currentPrc = ref(0)
const lyricsView = ref(null)
let lyricsEle = []

// 处理歌词变化
const handleLyricsChange = (currentIndex, currentTime) => {
    const added = currentIndex.filter(i => !previousIndex.value.includes(i))
    const removed = previousIndex.value.filter(i => !currentIndex.includes(i))

    removed.forEach(index => {
        console.log('removed', index);
        //lyricsEle[index].classList.remove('on-play')
        //lyricsEle[index].classList.add('off-play')


    })

    added.forEach(index => {
        console.log('added', index);
        for (let i = 0; i < lyricsEle.length; i++) {
            if (i === index) {
                lyricsEle[i].classList.add('on-play')
                lyricsEle[i].classList.remove('off-play')
            } else {
                lyricsEle[i].classList.remove('on-play')
                lyricsEle[i].classList.add('off-play')
            }
        }
        let chils = lyricsEle[index].children

        for (let i = 0; i < chils.length; i++) {

            let animates = chils[i].getAnimations()
            animates.forEach(animate => {
                animate.cancel()
            })
            let chil = chils[i]
            let sty = getLyricsStyle(chil.offsetWidth, chil.offsetHeight, 0.05)
            //console.log(sty);
            chil.style.backgroundImage = sty.backgroundImage
            chil.style.backgroundSize = sty.backgroundSize
            chil.style.backgroundPositionX = sty.backgroundPositionX
            chil.style.backgroundRepeat = sty.backgroundRepeat

            //console.log(chil);
            //
            chil.animate([
                {}, {
                    backgroundPositionX: "0px"
                }
            ], {
                duration: lyrics.value[index].words[i].end_time - lyrics.value[index].words[i].start_time,
                delay: lyrics.value[index].words[i].start_time - currentTime,
                fill: 'forwards'
            }
            )
        }

    })

    previousIndex.value = [...currentIndex]
}

// 核心高亮逻辑
const highlightLyrics = (currentTime) => {
    const currentIndex = []
    lyrics.value.forEach((line, index) => {
        if (currentTime >= line.start_time && currentTime <= line.end_time) {
            currentIndex.push(index)
        }
    })

    if (currentIndex.length !== previousIndex.value.length ||
        !currentIndex.every((val, i) => val === previousIndex.value[i])) {
        handleLyricsChange(currentIndex, currentTime)
    }
}
function getScale(duration) {
    if (duration < 1333) return 1.1;
    if (duration > 2000) return 1.3;
    return 1 + (duration - 1333) / 1000 * 0.3;
}
function getTrY(duration) {
    if (duration < 1000) return 0;
    if (duration > 2000) return -8;
    return -(duration - 1000) / 1000 * 8;
}
onMounted(() => {
    document.body.addEventListener("mousemove",handleMouseMove)
    document.body.addEventListener("mouseup",handleMouseUp)
    document.body.addEventListener("mousedown",handleMouseDown)

    GetNowLyrics().then(res => {

        lyrics.value = res
    })
    Events.On("set_lyric_from_qrc", (data) => {
        lyrics.value = data.data[0] // 假设传递的是解析后的歌词数组


    })

    Events.On("amll_play_progress", (data) => {
        currentPrc.value = data.data[0]
        //highlightLyrics(data.data[0] + 500) // 假设进度是毫秒数
    })
    Events.On("amll_lyrics_add", (data) => {
        let index = data.data[0]
        let currentTime = currentPrc.value
        console.log('added', index);
        for (let i = 0; i < lyricsEle.length; i++) {
            if (i === index) {
                lyricsEle[i].classList.add('on-play')
                lyricsEle[i].classList.remove('off-play')
            } else {
                lyricsEle[i].classList.remove('on-play')
                lyricsEle[i].classList.add('off-play')
            }
        }
        let chils = lyricsEle[index].children


        for (let i = 0; i < chils.length; i++) {

            if (chils[i].children.length > 0) {
                const du = lyrics.value[index].words[i].end_time - lyrics.value[index].words[i].start_time
                const scale = getScale(du)
                const ty = getTrY(du)
                for (let j = 0; j < chils[i].children.length; j++) {
                    const chil = chils[i].children[j]
                    let doms = chils[i].children
                    let animates = chil.getAnimations()


                    animates.forEach(animate => {
                        animate.cancel()
                    })


                    let sty = getLyricsStyle(chil.offsetWidth, chil.offsetHeight, 0)
                    const pjsj = (lyrics.value[index].words[i].end_time - lyrics.value[index].words[i].start_time) / chils[i].children.length
                    const baseDelay = lyrics.value[index].words[i].start_time - currentTime;
                    const charDelay = j * pjsj; // 字符级基础延迟
                    // console.log(sty);
                    chil.style.backgroundImage = sty.backgroundImage
                    chil.style.backgroundSize = sty.backgroundSize
                    chil.style.backgroundPositionX = sty.backgroundPositionX
                    chil.style.backgroundRepeat = sty.backgroundRepeat
                    chil.animate([
                        {
                            backgroundPositionX: sty.backgroundPositionX
                        }, {
                            backgroundPositionX: "0px",
                        }
                    ], {
                        duration: pjsj,
                        delay: lyrics.value[index].words[i].start_time - currentTime + j * pjsj,
                        fill: 'forwards',
                        easing: 'linear'
                    }
                    )
                    chil.animate([
                        {
                            transform: "scale(1)",
                            easing: "ease-out",
                        }, {
                            transform: `scale(${scale}) translateX(${getScaleOffset(j, scale, doms)}px) translateY(${ty}%)`,
                            easing: "ease-in",
                        }, {
                            transform: "scale(1)",
                            easing: "ease",
                        }
                    ], {
                        duration: (lyrics.value[index].words[i].end_time - lyrics.value[index].words[i].start_time) * 1.5,
                        delay: baseDelay + charDelay - charDelay * 0.6, // 通过系数控制重叠比例
                        fill: 'forwards',
                        easing: "ease"
                    }
                    )
                }
            } else {

                let animates = chils[i].getAnimations()
                animates.forEach(animate => {
                    animate.cancel()
                })
                let chil = chils[i]
                let sty = getLyricsStyle(chil.offsetWidth, chil.offsetHeight, 0.05)
                //console.log(sty);
                chil.style.backgroundImage = sty.backgroundImage
                chil.style.backgroundSize = sty.backgroundSize
                chil.style.backgroundPositionX = sty.backgroundPositionX
                chil.style.backgroundRepeat = sty.backgroundRepeat

                //console.log(chil);
                //
                chil.animate([
                    {
                        backgroundPositionX: sty.backgroundPositionX
                    }, {
                        backgroundPositionX: "0px"
                    }
                ], {
                    duration: lyrics.value[index].words[i].end_time - lyrics.value[index].words[i].start_time,
                    delay: lyrics.value[index].words[i].start_time - currentTime,
                    fill: 'forwards'
                }
                )
            }

        }

    })
    Events.On("amll_lyrics_remove", (data) => {
        console.log(data.data[0]);
    })
    Events.On("lyric_window_will_update_visibility", (data) => {
        if (data.data[0]) {
            document.body.style.visibility = "visible";
        } else {
            document.body.style.visibility = "hidden";
        }
    })
})

onUnmounted(() => {
    document.body.removeEventListener("mousemove",handleMouseMove)
    document.body.removeEventListener("mouseup",handleMouseUp)
    document.body.removeEventListener("mousedown",handleMouseDown)
    Events.Off("set_lyric_from_qrc")
    Events.Off("amll_play_progress")
    Events.Off("amll_lyrics_add")
    Events.Off("amll_lyrics_remove")
    window.removeEventListener('dblclick', dbc)
    Events.Off("lyric_window_will_update_visibility")
})
watch(lyrics, (newValue) => {
    lyricsView.value.innerHTML = ''
    lyricsEle = []
    console.log(newValue);

    for (let i = 0; i < newValue.length; i++) {
        const line = newValue[i]

        const lyricLineEle = document.createElement('p')
        lyricLineEle.className = 'lyric-line off-play'

        for (let j = 0; j < line.words.length; j++) {
            const word = line.words[j]

            const lyricWordEle = document.createElement('span')
            lyricWordEle.className = 'lyric-word'

            if (word.end_time - word.start_time > 1000) {
                for (let k = 0; k < word.word.length; k++) {
                    const ele = document.createElement('span')
                    ele.innerHTML = word.word[k].replace(/\s/g, '&nbsp;');
                    ele.className = 'lyric-letter'
                    lyricWordEle.appendChild(ele)
                }
            } else {
                // 将word.word中的所有空格 替换为&nbsp;
                lyricWordEle.innerHTML = word.word.replace(/\s/g, '&nbsp;');
            }



            lyricLineEle.appendChild(lyricWordEle)
        }
        lyricsEle.push(lyricLineEle)

        lyricsView.value.appendChild(lyricLineEle)


    }

    //for (let k = 0; k < lyricsEle.length; k++) {
    //
    //    for (let l = 0; l < lyricsEle[k].children; l++) {
    //        
    //        
    //        const ele = lyricsEle[k].children[l]
    //        console.log(ele);
    //        let sty = getLyricsStyle(ele.offsetWidth, ele.offsetHeight, 0)
    //        console.log(sty);
    //        ele.style = { ...ele.style, ...sty }
    //
    //
    //    }
    //}
}
)
function getScaleOffset(index, scale, doms) {
    const centerIndex = (doms.length - 1) / 2;

    // Calculate the cumulative width up to the current character
    let cumulativeWidth = 0;
    for (let i = 0; i < index; i++) {
        cumulativeWidth += doms[i].offsetWidth;
    }

    // Calculate the cumulative width up to the center character
    let centerCumulativeWidth = 0;
    for (let i = 0; i < centerIndex; i++) {
        centerCumulativeWidth += doms[i].offsetWidth;
    }

    // The offset is the difference between current position and center position,
    // multiplied by the scale factor
    return (cumulativeWidth - centerCumulativeWidth) * (scale - 1) * 0.5;
}

/*func generateBackgroundFadeStyle(elementWidth, elementHeight, fadeRatio float64) (string, string, float64, float64) {

    //const fadeWidth = elementHeight * fadeRatio;
    //  const widthRatio = fadeWidth / elementWidth;

    //  // 使用源码同款算法
    //  const totalAspect = 2 + widthRatio;
    //  const widthInTotal = widthRatio / totalAspect;
    //  const leftPos = (1 - widthInTotal) / 2;

    //  const from = leftPos * 100;
    //  const to = (leftPos + widthInTotal) * 100;

    //  const backgroundImage = `linear-gradient(to right,
    //    rgba(255, 255, 255, 1.0) ${from.toFixed(6)}%,
    //    rgba(255, 255, 255, 0.0) ${to.toFixed(6)}%)`;

    //  const backgroundSize = `${(totalAspect * 100).toFixed(3)}% 100%`;
    //  const totalPxWidth = elementWidth + fadeWidth;

    //  return {
    //    backgroundImage,
    //    backgroundSize,
    //    backgroundRepeat: 'no-repeat',
    //    backgroundOrigin: 'left',
    //    backgroundPositionX: `${-totalPxWidth}px`,
    //    finalPositionX: `0px`,
    //    transitionDistance: totalPxWidth,
    //  };

    fadeWidth := elementHeight * fadeRatio
    widthRatio := fadeWidth / elementWidth

    totalAspect := 2 + widthRatio
    widthInTotal := widthRatio / totalAspect
    leftPos := (1 - widthInTotal) / 2

    from := leftPos * 100
    to := (leftPos + widthInTotal) * 100

    backgroundImage := fmt.Sprintf("linear-gradient(to right,rgba(255, 255, 255, var(--rcolor)) %f%%,rgba(255, 255, 255, var(--color)) %f%%)", from, to)
    backgroundSize := fmt.Sprintf("%f%% 100%%", totalAspect*100)
    totalPxWidth := elementWidth + fadeWidth
    return backgroundImage, backgroundSize, -totalPxWidth, totalPxWidth
}
 */

function getLyricsStyle(elementWidth, elementHeight, fadeRatio) {
    const fadeWidth = elementHeight * fadeRatio
    const widthRatio = fadeWidth / elementWidth

    const totalAspect = 2 + widthRatio
    const widthInTotal = widthRatio / totalAspect
    const leftPos = (1 - widthInTotal) / 2

    const from = leftPos * 100
    const to = (leftPos + widthInTotal) * 100
    const backgroundImage = `linear-gradient(to right,
	rgb(111, 198, 255) ${from.toFixed(6)}%,
	rgba(255, 255, 255, 1) ${to.toFixed(6)}%)`;

    const backgroundSize = `${(totalAspect * 100).toFixed(3)}% 100%`;
    const totalPxWidth = elementWidth + fadeWidth;

    return {
        backgroundImage,
        backgroundSize,
        backgroundRepeat: 'no-repeat',
        backgroundOrigin: 'left',
        backgroundPositionX: `${-totalPxWidth}px`,
        finalPositionX: `0px`,
        transitionDistance: totalPxWidth,
    };
}
</script>

<style lang="css" scoped>

main {
    width: 100%;
    height:100%;
    overflow: hidden;
    /*--wails-draggable: drag;*/
    
    transition: box-shadow 0.3s, background-color 0.3s;
    border-radius: 10px;
    box-sizing: border-box;

    display: flex;
    align-items: center;
    justify-content: center;
}

.m{
    width: calc(100% - 20px);
    height: calc(100% - 20px);
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
}

main:hover {
    border: 1px solid rgba(131, 131, 131, 0.5);
    box-shadow: 0 0 100px rgba(58, 58, 58, 0.5) inset;
}

main:active {
    background-color: rgba(0, 0, 0, 0.3);
}
</style>

<style>
body{
    overflow: hidden;
}
.lyric-line {
    font-size: 4vw;
    text-align: center;
    transition: transform 0.5s, opacity 0.5s, filter 0.5s;
    position: absolute;
    font-weight: 600;

}

.lyric-word {
    font-family: "微软雅黑";
    display: inline-block;
    position: relative;

    background-repeat: no-repeat;
    background-clip: text;
    color: transparent;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    z-index: 1;
    -webkit-text-stroke: 0.15vw #000000d2;
}

.lyric-letter {
    display: inline-block;
    background-repeat: no-repeat;
    background-clip: text;
    color: transparent;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    will-change: transform;
}

.on-play {
    transform: scale(1);
}

.off-play {
    transform: scale(0.5);
    opacity: 0;
    filter: blur(10px);
}
</style>