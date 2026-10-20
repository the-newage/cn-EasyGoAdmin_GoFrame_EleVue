package captcha

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"easygoadmin/internal/consts"
	"easygoadmin/utility/cache"
)

// ============================== 图形验证码 ==============================
//
// 输出 PNG data URL，使用 TTF 字体渲染，可直接由 img 标签展示。
// 前端契约：GET /api/captcha 返回 data: {key, captcha}。
// 字体缺失或解析失败时降级为内置 SVG 矢量渲染，接口不失败。

// defaultChars 默认字符集，已排除易混淆字符（0/O、1/l/I）
// 可用配置 captcha.chars 与 captcha.length 覆盖
const defaultChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// SVG 降级画布尺寸
// 仅在字体路径缺失时使用
const svgWidth, svgHeight = 120, 40

var (
	fontOnce   sync.Once      // 保证字体仅加载一次
	parsedFont *opentype.Font // 已解析的字体实例
)

// loadFont 加载并缓存 TTF 字体
//
// 进程内只读取一次，失败后永久降级为 SVG 渲染
func loadFont(ctx context.Context, path string) *opentype.Font {
	fontOnce.Do(func() {
		data, err := os.ReadFile(path)
		if err != nil {
			g.Log().Warningf(ctx, "captcha: 字体文件不可读 %s: %v（降级 SVG 渲染）", path, err)
			return
		}
		f, err := opentype.Parse(data)
		if err != nil {
			g.Log().Warningf(ctx, "captcha: 字体解析失败 %s: %v（降级 SVG 渲染）", path, err)
			return
		}
		parsedFont = f
	})
	return parsedFont
}

// Generate 生成验证码
//
// 返回验证码标识与 data URL，标识写入缓存，TTL 取配置 captcha.expire
func Generate(ctx context.Context) (key, dataURL string, err error) {
	length := g.Cfg().MustGet(ctx, "captcha.length").Int()
	if length < 4 || length > 8 {
		length = 4
	}
	code := randomCode(length)
	key = randomKey()
	expire := g.Cfg().MustGet(ctx, "captcha.expire").Int()
	if expire <= 0 {
		expire = 300
	}
	if err = cache.Set(ctx, consts.CacheKeyCaptcha(key), code, expire); err != nil {
		return "", "", err
	}
	fontPath := g.Cfg().MustGet(ctx, "captcha.fontPath").String()
	if fontPath == "" {
		fontPath = "resource/static/fonts/Vera.ttf"
	}
	if f := loadFont(ctx, fontPath); f != nil {
		if s, e := renderPNG(ctx, f, code); e == nil {
			return key, s, nil
		} else {
			g.Log().Warningf(ctx, "captcha: PNG 渲染失败: %v（降级 SVG 渲染）", e)
		}
	}
	return key, renderSVG(code), nil
}

// Verify 校验并消费验证码
//
// 验证码为一次性使用；标识为空时返回 true，走宽松模式
func Verify(ctx context.Context, key, code string) bool {
	if key == "" {
		return true
	}
	// 万能验证码（配置 captcha.bypassCode，仅开发调试用，生产留空）
	if bypass := g.Cfg().MustGet(ctx, "captcha.bypassCode").String(); bypass != "" &&
		strings.EqualFold(bypass, strings.TrimSpace(code)) {
		ck := consts.CacheKeyCaptcha(key)
		_ = cache.Del(ctx, ck)
		return true
	}
	ck := consts.CacheKeyCaptcha(key)
	expected, ok := cache.Get(ctx, ck)
	_ = cache.Del(ctx, ck) // 无论成败均消费，防爆破重放
	if !ok || expected == "" {
		return false
	}
	return strings.EqualFold(expected, strings.TrimSpace(code))
}

// randomCode 生成指定长度的随机验证码
func randomCode(n int) string {
	charset := g.Cfg().MustGet(context.Background(), "captcha.chars").String()
	if charset == "" {
		charset = defaultChars
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		sb.WriteByte(charset[idx.Int64()])
	}
	return sb.String()
}

// randomKey 生成随机验证码标识
func randomKey() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// ----------------------------------------------------------------------------
// PNG 渲染（主路径）
// 使用 TTF 字体，配合随机旋转、颜色与位置，叠加干扰线与噪点。
// 相关配置：captcha.imageWidth、captcha.imageHeight、captcha.fontSize、
// captcha.noiseLines、captcha.noiseDots
// ----------------------------------------------------------------------------

// pngColors 字符颜色池（绿色系）
var pngColors = []color.RGBA{
	{0, 128, 0, 255},   // 绿色
	{0, 150, 0, 255},   // 稍亮绿
	{34, 139, 34, 255}, // 森林绿
	{50, 205, 50, 255}, // 酸橙绿
	{0, 100, 0, 255},   // 深绿
}

// renderPNG 渲染验证码为 PNG data URL
func renderPNG(ctx context.Context, f *opentype.Font, code string) (string, error) {
	w := cfgInt(ctx, "captcha.imageWidth", 160)
	h := cfgInt(ctx, "captcha.imageHeight", 50)
	fs := cfgInt(ctx, "captcha.fontSize", 28)
	lines := cfgInt(ctx, "captcha.noiseLines", 8)
	dots := cfgInt(ctx, "captcha.noiseDots", 500)

	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: float64(fs), DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return "", err
	}
	defer func() { // opentype face 实现了 io.Closer，防御性断言
		if c, ok := face.(io.Closer); ok {
			_ = c.Close()
		}
	}()

	// 字符总宽（含字距），画布不够时自动加宽——防长验证码右侧被裁剪
	const gap = 4
	total := 0
	for _, ch := range code {
		adv, _ := face.GlyphAdvance(ch)
		total += adv.Ceil() + gap
	}
	const margin = 16
	if need := total + margin*2; w < need {
		w = need
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	fillRGBA(img, color.RGBA{247, 249, 252, 255})

	// 干扰线：深绿色
	lineColor := color.RGBA{0, 100, 0, 255}
	for i := 0; i < lines; i++ {
		drawLine(img, randInt(w), randInt(h), randInt(w), randInt(h), lineColor)
	}

	// 字符（逐字独立小画布 → 旋转 → 混合到主画布）
	cellH := fs * 2
	x := (w - total) / 2
	for _, ch := range code {
		adv, _ := face.GlyphAdvance(ch)
		cw := adv.Ceil() + 8
		cell := image.NewRGBA(image.Rect(0, 0, cw, cellH))
		col := pngColors[randInt(len(pngColors))]
		d := font.Drawer{
			Dst:  cell,
			Src:  image.NewUniform(col),
			Face: face,
			Dot:  fixed.P(4, cellH/2+fs/3),
		}
		d.DrawString(string(ch))
		angle := float64(randInt(50)-25) * math.Pi / 180
		dx := x - 4 + randInt(5) - 2
		dy := h/2 - cellH/2 + randInt(5) - 2
		rotateBlend(img, cell, dx, dy, cw/2, cellH/2, angle)
		x += adv.Ceil() + gap
	}

	// 噪点：半透明绿色
	dotColor := color.RGBA{0, 128, 0, 90}
	for i := 0; i < dots; i++ {
		blendPixel(img, randInt(w), randInt(h), dotColor)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// cfgInt 读取正整数配置
//
// 配置缺失或非正时返回默认值
func cfgInt(ctx context.Context, key string, def int) int {
	if v := g.Cfg().MustGet(ctx, key).Int(); v > 0 {
		return v
	}
	return def
}

// fillRGBA 用指定颜色填充整张图像
func fillRGBA(img *image.RGBA, c color.RGBA) {
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i] = c.R
		img.Pix[i+1] = c.G
		img.Pix[i+2] = c.B
		img.Pix[i+3] = 255
	}
}

// rotateBlend 旋转并混合
//
// 将 src（透明底）绕自身中心 (scx,scy) 旋转 angle 后，
// 放到 dst 的 (dx,dy) 处并做 alpha 混合
func rotateBlend(dst, src *image.RGBA, dx, dy, scx, scy int, angle float64) {
	sin, cos := math.Sin(angle), math.Cos(angle)
	sb := src.Bounds()
	cx, cy := float64(dx+scx), float64(dy+scy)
	for y := dy; y < dy+sb.Dy(); y++ {
		for x := dx; x < dx+sb.Dx(); x++ {
			// 逆旋转求源坐标
			relX, relY := float64(x)-cx, float64(y)-cy
			sx := int(relX*cos+relY*sin) + scx
			sy := int(-relX*sin+relY*cos) + scy
			if sx < 0 || sy < 0 || sx >= sb.Dx() || sy >= sb.Dy() {
				continue
			}
			si := src.PixOffset(sx, sy)
			a := src.Pix[si+3]
			if a == 0 {
				continue
			}
			di := imgOffset(dst, x, y)
			if di < 0 || di+3 >= len(dst.Pix) {
				continue
			}
			// src 为预乘 alpha（image.RGBA 约定）：dst = src + dst*(1-a)
			fa := 255 - uint32(a)
			dst.Pix[di] = uint8(uint32(src.Pix[si]) + uint32(dst.Pix[di])*fa/255)
			dst.Pix[di+1] = uint8(uint32(src.Pix[si+1]) + uint32(dst.Pix[di+1])*fa/255)
			dst.Pix[di+2] = uint8(uint32(src.Pix[si+2]) + uint32(dst.Pix[di+2])*fa/255)
			dst.Pix[di+3] = 255
		}
	}
}

// imgOffset 计算像素在图像数据中的偏移量
//
// 坐标越界时返回 -1
func imgOffset(img *image.RGBA, x, y int) int {
	if !image.Pt(x, y).In(img.Bounds()) {
		return -1
	}
	return (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*4
}

// drawLine 使用 Bresenham 算法绘制直线
//
// 不透明绘制，用于生成干扰线
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		if i := imgOffset(img, x0, y0); i >= 0 {
			img.Pix[i] = c.R
			img.Pix[i+1] = c.G
			img.Pix[i+2] = c.B
			img.Pix[i+3] = 255
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

// blendPixel 单像素 alpha 混合
//
// 用于绘制噪点，传入颜色为直 alpha，先转预乘再叠加
func blendPixel(img *image.RGBA, x, y int, c color.RGBA) {
	i := imgOffset(img, x, y)
	if i < 0 {
		return
	}
	av := uint32(c.A)
	pr, pg, pb := uint32(c.R)*av/255, uint32(c.G)*av/255, uint32(c.B)*av/255
	fa := 255 - av
	img.Pix[i] = uint8(pr + uint32(img.Pix[i])*fa/255)
	img.Pix[i+1] = uint8(pg + uint32(img.Pix[i+1])*fa/255)
	img.Pix[i+2] = uint8(pb + uint32(img.Pix[i+2])*fa/255)
}

// abs 计算整数绝对值
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ----------------------------------------------------------------------------
// SVG 渲染（降级路径）
// 字体缺失时兜底，保持接口可用。
// ----------------------------------------------------------------------------

// renderSVG 渲染带干扰的 SVG 验证码
func renderSVG(code string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		svgWidth, svgHeight, svgWidth, svgHeight))
	sb.WriteString(`<rect width="100%" height="100%" fill="#f7f9fc"/>`)
	// 干扰线
	for i := 0; i < 4; i++ {
		x1 := randInt(svgWidth)
		y1 := randInt(svgHeight)
		x2 := randInt(svgWidth)
		y2 := randInt(svgHeight)
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#c8d3e0" stroke-width="1"/>`, x1, y1, x2, y2))
	}
	// 字符
	for i, ch := range code {
		x := 14 + i*26 + randInt(6)
		y := 27 + randInt(6)
		rot := randInt(50) - 25
		colorIdx := randInt(len(pngColors))
		c := pngColors[colorIdx]
		sb.WriteString(fmt.Sprintf(
			`<text x="%d" y="%d" font-size="24" font-family="Arial, sans-serif" font-weight="bold" fill="#%02x%02x%02x" transform="rotate(%d %d %d)">%c</text>`,
			x, y, c.R, c.G, c.B, rot, x, y, ch))
	}
	// 噪点
	for i := 0; i < 18; i++ {
		sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="1" fill="#aab7c4"/>`, randInt(svgWidth), randInt(svgHeight)))
	}
	sb.WriteString(`</svg>`)
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(sb.String()))
}

// randInt 生成指定上限内的随机整数
func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	v, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
	return int(v.Int64())
}
