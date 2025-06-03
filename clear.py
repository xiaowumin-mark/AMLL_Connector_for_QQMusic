import re

# 定义元数据关键词和正则表达式
metadata_patterns = [
    r'^\s*(作词|作曲|编曲|演唱|歌手|专辑|出品|监制|录音|混音|母带|吉他|贝斯|鼓|键盘|弦乐|和声|版权|制作人|原唱|翻唱)\s*[:：]',
    r'^\s*(OP|SP|Lyrics|Composed|Produced|Vocals|Engineer|Publisher)\s*[:：]',
    r'^\s*(发行|特别鸣谢|宣传|录音棚|Studio|感谢)\s*[:：]',
    r'^\s*未经.*?许可.*?不得.*?使用.*?$',
    r'^\s*未经.*?许可.*?不得.*?方式.*?$',
    r'^\s*未经著作权人书面许可，\s*不得以任何方式\s*（?包括.*?等）?\s*使用\s*$',
    r'^\s*发行方\s*[:：].*?$',
    r'^\s*.*?(工作室|特别企划).*?$'
]

# 编译正则表达式
compiled_patterns = [re.compile(pattern, re.IGNORECASE) for pattern in metadata_patterns]

def is_metadata_line(line):
    """判断一行是否为元数据行"""
    for pattern in compiled_patterns:
        if pattern.match(line):
            print(f"Matched metadata line: {line.strip()}")
            return True
        print(f"Did not match metadata line: {line.strip()}")
    return False

def clean_lyrics(text):
    """清除歌词文本中的元数据，保留纯净的歌词内容"""
    lines = text.strip().split('\n')
    # 去除头部元数据行
    while lines and is_metadata_line(lines[0]):
        lines.pop(0)
    # 去除尾部元数据行
    while lines and is_metadata_line(lines[-1]):
        lines.pop()
    return '\n'.join(lines)

# 示例使用
sample_text = """暮色回响 (《默杀》电影推广曲) - 张韶涵
词：T１
曲：T１
制作人：Tim姜皓天
编曲：Tim姜皓天
吉他：伍凌枫FlinWu
贝斯：彭宏立Simon
钢琴：Tim姜皓天/马乐
和声：王子@Soulkidz
弦乐编写：何山/黄宇弘
弦乐监制：胡静成/黄宇弘
弦乐演奏：国际首席爱乐乐团
弦乐录音：王小四@金田录音棚
配唱制作：杨钧尧 Bryan Yang
人声录音：潘尧泓 Hendrik Pan
录音室：白金录音室 Platinum Studio
音频编辑：李宗远@Studio２１A
混音母带：周天澈@Studio２１A
项目支持：李金蕾/彭珂/雷宇/郑薇/万馨蔚/沈雯
词曲OP：乐宣时代（西安）
营销策划：王淼
总企划：张馨艺 Summer
监制：LSD/唐潇
出品人：师晶
鸣谢：猫眼宣传中心/猫眼阿尔法短视频工作室/猫眼物料营销团队
特别鸣谢：Angmuzik天涵有限公司/乐宣时代（西安）
传说浩瀚银河有颗星是他
走出时间后仍然选择留下
漆昼中温柔的不像话
静守着他的遗憾啊
旧的摇椅吱吱呀呀停不下"""

cleaned_lyrics = clean_lyrics(sample_text)
print(cleaned_lyrics)
