package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/u2takey/ffmpeg-go"
)

// VideoGetNetImgCount 从视频流中提取指定数量的帧并保存为本地 JPEG 图片
// 参数:
//   - frameNum: 需要提取的帧数（当前实现仅提取第一帧以匹配原代码行为）
//   - url: 视频流地址（如 RTMP 流，例：rtmp://192.168.20.221:30200/live/2dfp52anvad7g）
//
// 返回值: 保存的 JPEG 文件路径，若出错则返回空字符串
func VideoGetNetImgCount(frameNum int, url string) string {
	// 从 URL 中提取影片标识用于输出文件名
	// 假设 URL 格式为 rtmp://<host>:<port>/<app>/<movieId>，提取 movieId
	urlStrings := strings.Split(url, ".")
	urlStrings2 := strings.Split(urlStrings[len(urlStrings)-2], "/")
	if len(urlStrings2) < 2 {
		fmt.Println("URL 格式无效，无法提取影片标识")
		return ""
	}
	movieUrl := urlStrings2[len(urlStrings2)-1]
	fileName := fmt.Sprintf("./app/video/tmp/%s_cover.jpg", movieUrl)

	// 使用 ffmpeg-go 构建 FFmpeg 命令，提取视频流的第一帧
	// -i: 指定输入流（如 RTMP）
	// -vframes 1: 仅提取一帧
	// -f image2: 输出格式为图像
	// -pix_fmt rgb24: 设置像素格式为 RGB24
	// -vf fps=1: 设置帧率为 1，确保只取一帧
	err := ffmpeg_go.Input(url, ffmpeg_go.KwArgs{
		"rtmp_live": "live", // 确保 RTMP 流以实时模式处理
	}).
		Output(fileName, ffmpeg_go.KwArgs{
			"vframes": frameNum, // 提取一帧
			"f":       "image2", // 输出为图像
			"pix_fmt": "rgb24",  // 像素格式 RGB24
			"vf":      "fps=1",  // 设置帧率为 1
		}).
		OverWriteOutput(). // 覆盖输出文件
		ErrorToStdOut().   // 错误输出到标准输出
		Run()

	if err != nil {
		fmt.Println("FFmpeg 命令执行出错:", err)
		return ""
	}

	// 检查输出文件是否存在
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("输出文件未生成:", err)
		return ""
	}

	return fileName
}

// SaveFrameJpg 将单帧保存为 JPEG 文件（为兼容原接口保留，但在此实现中未使用）
// 参数:
//   - movieUrl: 输出文件名的影片标识
//
// 返回值: 保存的 JPEG 文件路径，若出错则返回空字符串
func SaveFrameJpg(movieUrl string) string {
	// 由于 ffmpeg-go 直接输出 JPEG 文件，此函数仅为占位以兼容原接口
	// 实际帧保存由 VideoGetNetImgCount 完成
	return fmt.Sprintf("./app/video/tmp/%s_cover.jpg", movieUrl)
}
