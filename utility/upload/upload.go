// Package upload 文件上传
//
// 采用两段式流程：先落入临时目录，提交保存时转正至正式目录。
//  1. SaveTemp：上传到 {uploadDir}/{tempDir}/日期/ 目录，含白名单与大小校验
//  2. SaveFromTemp：表单提交保存时，将临时文件移动到正式子目录（如 user/article）
//
// 路径契约：数据库中统一存储相对 uploadDir 的路径，例如 "tenant/20260923/x.png"；
// 对外访问地址由 GetFileUrl 拼接 file.domainUrl 生成。
package upload

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// cfgType 文件上传配置
//
// 对应 manifest/config/config.yaml 中的 file 配置项
type cfgType struct {
	UploadDir string // 上传根目录
	TempDir   string // 临时目录名
	DomainUrl string // 文件访问域名
	ImageExt  string // 允许的图片扩展名
	FileExt   string // 允许的文件扩展名
	MaxSizeMB int    // 单文件大小上限（MB）
}

// loadCfg 加载文件上传配置
//
// 临时目录名缺省为 temp，大小上限缺省为 10MB
func loadCfg(ctx context.Context) cfgType {
	c := cfgType{
		UploadDir: g.Cfg().MustGet(ctx, "file.uploadDir").String(),
		TempDir:   g.Cfg().MustGet(ctx, "file.tempDir").String(),
		DomainUrl: g.Cfg().MustGet(ctx, "file.domainUrl").String(),
		ImageExt:  g.Cfg().MustGet(ctx, "file.imageExt").String(),
		FileExt:   g.Cfg().MustGet(ctx, "file.fileExt").String(),
		MaxSizeMB: g.Cfg().MustGet(ctx, "file.maxSize").Int(),
	}
	if c.TempDir == "" {
		c.TempDir = "temp"
	}
	if c.MaxSizeMB <= 0 {
		c.MaxSizeMB = 10
	}
	return c
}

// Result 上传结果
type Result struct {
	Url  string // 可访问地址：domainUrl + /temp/20261013/xxx.png
	Path string // 相对 uploadDir 的路径，入库形态
	Name string // 原始文件名
	Size int64  // 文件大小（字节）
	Ext  string // 扩展名
	Mime string // MIME 类型（按文件头内容嗅探）
}

// BaseDir 上传根目录
//
// 取配置 file.uploadDir，文件统一落盘于此
func BaseDir(ctx context.Context) string {
	return loadCfg(ctx).UploadDir
}

// TempBaseDir 临时目录绝对路径
//
// 供清理任务使用
func TempBaseDir(ctx context.Context) string {
	c := loadCfg(ctx)
	return filepath.Join(c.UploadDir, c.TempDir)
}

// GetAccessBase 访问基址
//
// 取配置 file.domainUrl，即 uploadDir 的对外访问地址，不含末尾斜杠
func GetAccessBase(ctx context.Context) string {
	return strings.TrimRight(loadCfg(ctx).DomainUrl, "/")
}

// GetFileUrl 相对路径转可访问地址
//
// 库中存储相对 uploadDir 的路径，输出为 file.domainUrl 加相对路径；
// 空值与已是完整 URL 的值原样返回
func GetFileUrl(ctx context.Context, path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return GetAccessBase(ctx) + "/" + strings.TrimLeft(path, "/")
}

// SaveTemp 上传到临时目录
//
// typ 为 file 或 image，分别按 fileExt 与 imageExt 白名单及 maxSize 校验；
// 图片类型额外做内容嗅探
func SaveTemp(ctx context.Context, file *ghttp.UploadFile, typ string) (*Result, error) {
	if file == nil {
		return nil, fmt.Errorf("请选择文件")
	}
	c := loadCfg(ctx)

	// 1) 大小校验
	if file.Size > int64(c.MaxSizeMB)*1024*1024 {
		return nil, fmt.Errorf("文件大小超过限制(%dMB)", c.MaxSizeMB)
	}

	// 2) 扩展名白名单
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
	if ext == "" {
		return nil, fmt.Errorf("无法识别文件扩展名")
	}
	allow := c.FileExt
	if typ == "image" {
		allow = c.ImageExt
	}
	if !extIn(ext, allow) {
		return nil, fmt.Errorf("不允许的文件类型: .%s", ext)
	}

	// 3) 读取内容做图片嗅探（type=image 时）
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	head = head[:n]
	if typ == "image" && !isImage(head) {
		return nil, fmt.Errorf("图片文件内容无效或已损坏")
	}
	// MIME 嗅探：按文件头内容识别，识别不出时回退通用二进制类型
	mime := "application/octet-stream"
	if n > 0 {
		if m := http.DetectContentType(head); m != "" {
			mime = m
		}
	}
	if _, err = src.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// 4) 落盘 temp/{日期}/{32位随机}.{ext}
	dateDir := time.Now().Format("20060102")
	relDir := filepath.Join(c.TempDir, dateDir)
	fullDir := filepath.Join(c.UploadDir, relDir)
	if err = os.MkdirAll(fullDir, 0o755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %v", err)
	}
	name := randomHex(16) + "." + ext
	full := filepath.Join(fullDir, name)
	dst, err := os.Create(full)
	if err != nil {
		return nil, err
	}
	size, err := io.Copy(dst, src)
	if cerr := dst.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(full)
		return nil, err
	}

	rel := filepath.ToSlash(filepath.Join(relDir, name))
	return &Result{
		Url:  GetFileUrl(ctx, rel),
		Path: rel,
		Name: file.Filename,
		Size: size,
		Ext:  ext,
		Mime: mime,
	}, nil
}

// SaveFromTemp 将临时文件转正
//
// 返回值即入库形态，为相对 uploadDir 的路径：
//   - 空值返回空串
//   - 本站 URL 或相对路径：剥离域名前缀，临时文件移动到 {uploadDir}/{saveDir}/{日期}/{时分秒}_{随机}.{ext}
//   - 外部 URL：原样返回
func SaveFromTemp(ctx context.Context, fileUrl, saveDir string) string {
	if fileUrl == "" {
		return ""
	}
	c := loadCfg(ctx)

	rel := fileUrl
	// 剥离本站域名（domainUrl 配置），得到相对 uploadDir 的路径
	stripped := false
	if c.DomainUrl != "" && strings.HasPrefix(rel, c.DomainUrl) {
		rel = strings.TrimPrefix(rel, c.DomainUrl)
		stripped = true
	}
	rel = strings.TrimPrefix(rel, "/")

	// 仅 temp 需要转移；已是正式相对路径则归一化后返回（域名已剥离）
	if !strings.HasPrefix(rel, c.TempDir+"/") {
		if !stripped && (strings.HasPrefix(rel, "http://") || strings.HasPrefix(rel, "https://")) {
			return fileUrl // 外部URL（未匹配本站域名）原样保留
		}
		return rel
	}
	srcFull := filepath.Join(c.UploadDir, filepath.FromSlash(rel))
	if _, err := os.Stat(srcFull); err != nil {
		return rel // 文件不存在则维持相对路径（幂等）
	}

	ext := strings.TrimPrefix(filepath.Ext(rel), ".")
	if ext == "" {
		ext = "bin"
	}
	dateDir := time.Now().Format("20060102")
	name := time.Now().Format("150405") + "_" + randomHex(4) + "." + ext
	newRelDir := filepath.ToSlash(filepath.Join(strings.Trim(saveDir, "/"), dateDir))
	if err := os.MkdirAll(filepath.Join(c.UploadDir, filepath.FromSlash(newRelDir)), 0o755); err != nil {
		return rel
	}
	newRel := newRelDir + "/" + name
	dstFull := filepath.Join(c.UploadDir, filepath.FromSlash(newRel))
	if err := os.Rename(srcFull, dstFull); err != nil {
		// 跨卷回退：复制后删源
		if data, rerr := os.ReadFile(srcFull); rerr == nil {
			if werr := os.WriteFile(dstFull, data, 0o644); werr == nil {
				_ = os.Remove(srcFull)
			} else {
				return rel
			}
		} else {
			return rel
		}
	}
	return newRel
}

// SaveFromTempMulti 逗号分隔多文件批量转正
//
// 适用于 article.images 这类多图字段
func SaveFromTempMulti(ctx context.Context, fileUrls, saveDir string) string {
	if fileUrls == "" {
		return ""
	}
	parts := strings.Split(fileUrls, ",")
	for i, p := range parts {
		parts[i] = SaveFromTemp(ctx, strings.TrimSpace(p), saveDir)
	}
	return strings.Join(parts, ",")
}

// extIn 判断扩展名是否在逗号分隔的白名单内
func extIn(ext, allow string) bool {
	for _, a := range strings.Split(allow, ",") {
		if strings.TrimSpace(a) == ext {
			return true
		}
	}
	return false
}

// isImage 按魔数嗅探常见图片格式
func isImage(b []byte) bool {
	if len(b) < 4 {
		return false
	}
	mime := http.DetectContentType(b)
	return strings.HasPrefix(mime, "image/")
}

// randomHex 生成随机十六进制串
//
// n 字节随机数转为 2n 位十六进制字符，使用加密随机源防猜测
func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
