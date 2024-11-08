package captcha

import "github.com/mojocn/base64Captcha"

type VerifyType uint8

const (
	VerifyTypeDigit     VerifyType = iota + 1 // 1 数字验证码
	VerifyTypeCharacter                       // 2 字符验证码
	VerifyTypeAudio                           // 3 音频验证码
	VerifyTypeLanguage                        // 4 语言验证码
	VerifyTypeMath                            // 5 数学验证码
	VerifyTypeChinese                         // 6 中文验证码
)

type Options struct {
	Engine               VerifyType
	OptionShowSineLine   bool
	OptionShowSlimeLine  bool
	OptionShowHollowLine bool
	Height               int
	Width                int
	CaptchaLen           int
	NoiseCount           int
	BgColor              string
	Audio                Audio
	Digit                // 数字验证码
}

type Digit struct {
	Length   int
	MaxSkew  float64
	DotCount int
}

type Audio struct {
	Language string
}

func WithEngine(engine VerifyType) Option {
	return func(o *Options) {
		o.Engine = engine
	}
}

func WithOptionShowSineLine(showSineLine bool) Option {
	return func(options *Options) {
		options.OptionShowSineLine = showSineLine
	}
}

func WithOptionShowSlimeLine(showSlimeLine bool) Option {
	return func(o *Options) {
		o.OptionShowSlimeLine = showSlimeLine
	}
}

func WithOptionShowHollowLine(showHollowLine bool) Option {
	return func(o *Options) {
		o.OptionShowHollowLine = showHollowLine
	}
}

func WithCaptchaLen(captchaLen int) Option {
	return func(o *Options) {
		o.CaptchaLen = captchaLen
	}
}

func WithLanguage(language string) Option {
	return func(o *Options) {
		o.Audio.Language = language
	}
}

func WithHeight(height int) Option {
	return func(o *Options) {
		o.Height = height
	}
}

func WithWidth(width int) Option {
	return func(o *Options) {
		o.Width = width
	}
}

func WithNoiseCount(noiseCount int) Option {
	return func(o *Options) {
		o.NoiseCount = noiseCount
	}
}

func WithCharacterCaptchaLen(captchaLen int) Option {
	return func(o *Options) {
		o.CaptchaLen = captchaLen
	}
}

func WithBgColor(bgColor string) Option {
	return func(o *Options) {
		o.BgColor = bgColor
	}
}

func (o *Options) GetShowLine() int {
	var optionShowLine int
	if o.OptionShowHollowLine {
		optionShowLine |= base64Captcha.OptionShowHollowLine
	}
	if o.OptionShowSlimeLine {
		optionShowLine |= base64Captcha.OptionShowSlimeLine
	}
	if o.OptionShowSineLine {
		optionShowLine |= base64Captcha.OptionShowSineLine
	}
	return optionShowLine
}

func WithDigitLength(length int) Option {
	return func(o *Options) {
		o.Digit.Length = length
	}
}

func WithDigitMaxSkew(maxSkew float64) Option {
	return func(o *Options) {
		o.Digit.MaxSkew = maxSkew
	}
}

func WithDigitDotCount(dotCount int) Option {
	return func(o *Options) {
		o.Digit.DotCount = dotCount
	}
}
