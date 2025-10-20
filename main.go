package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type DownloadTask struct {
	URL      string
	FilePath string
}

func main() {
	tasks := []DownloadTask{
		{URL: "https://example.com/file1.zip", FilePath: "file1.zip"},
		{URL: "https://example.com/file2.zip", FilePath: "file2.zip"},
	}
	var wg sync.WaitGroup
	taskChannel := make(chan DownloadTask)
	var download_Cnt int = 5
	fmt.Scanf("%d", &download_Cnt)
	for i := 0; i < download_Cnt; i++ {
		go downloadWorker(taskChannel, &wg)
	}
	for _, task := range tasks {
		wg.Add(1)
		taskChannel <- task
	}
	close(taskChannel)
	wg.Wait()
}

// 下载工作函数
func downloadWorker(taskChannel chan DownloadTask, wg *sync.WaitGroup) {
	for task := range taskChannel {
		if err := downloadFile(task.URL, task.FilePath); err != nil {
			fmt.Printf("下载失败：%s, 错误：%v\n", task.URL, err)
		} else {
			fmt.Printf("下载完成：%s\n", task.FilePath)
		}
		wg.Done()
	}
}

// 下载文件
func downloadFile(url, filePath string) error {
	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 获取文件内容
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 将文件写入本地
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
