<template>
    <main>
        <div class="topbar">
            <div class="item">
                ACFQ
            </div>
        </div>
        <div class="container">
            <div class="logo">
                <img src="../../public/icon/icon_1x.png" alt="" srcset="">
            </div>
            <div class="text">
                适用于QQ音乐的AMLL连接器
            </div>
        </div>
        <div class="bottombar">
            正在启动中...
        </div>
    </main>
</template>
<style lang="css" scoped>
main {
    width: 100%;
    height: 100%;
}

.topbar {
    width: 100%;
    height: 50px;
    display: flex;
    align-items: center;
    .item {
        height: 25px;
        border-radius: 5px;
        box-sizing: border-box;
        padding: 5px 10px;
        margin-left: 10px;
        color: #fff;
        background-color: rgba(0, 26, 255, 0.212);
        font-size: 12px;
    }
}

.container {
    width: 100%;
    height: calc(100% - 80px);
    display: flex;
    flex-direction: column;
    align-items: center;

    .logo{
        width: 100px;
        height: 100px;
        img{
            width: 100px;
            height: 100px;
        }
    }
    .text{
        font-size: 16px;
        color: #fff;
        margin-top: 20px;
    }
}

.bottombar {
    width: 100%;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: start;
    color: #ffffff8a;
    padding-left: 10px;
    box-sizing: border-box;
    font-size: 15px;
    animation: loading 1s steps(1);
}
@keyframes loading {
    0% {
        opacity: 0;
    }
    100% {
        opacity: 1;
    }

}
</style>
