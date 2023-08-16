package captcha

import (
	"context"
	"fmt"
	"image/color"
	"strings"

	"github.com/leor-w/kid/utils"

	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/plugin"

	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	store         Store `inject:""`
	driver        base64Captcha.Driver
	embeddedFonts *base64Captcha.EmbeddedFontsStorage
	options       Options
}

type Option func(*Options)

func (c *Captcha) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(plugin.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("captcha%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithEngine(VerifyType(config.GetInt(utils.GetConfigurationItem(confPrefix, "engine")))),
		WithOptionShowHollowLine(config.GetBool(utils.GetConfigurationItem(confPrefix, "option.showHollowLine"))),
		WithOptionShowSlimeLine(config.GetBool(utils.GetConfigurationItem(confPrefix, "option.showSlimeLine"))),
		WithOptionShowSineLine(config.GetBool(utils.GetConfigurationItem(confPrefix, "option.showSineLine"))),
		WithHeight(config.GetInt(utils.GetConfigurationItem(confPrefix, "height"))),
		WithWidth(config.GetInt(utils.GetConfigurationItem(confPrefix, "width"))),
		WithCaptchaLen(config.GetInt(utils.GetConfigurationItem(confPrefix, "captchaLen"))),
		WithNoiseCount(config.GetInt(utils.GetConfigurationItem(confPrefix, "noiseCount"))),
		WithBgColor(config.GetString(utils.GetConfigurationItem(confPrefix, "bgColor"))),
		WithLanguage(config.GetString(utils.GetConfigurationItem(confPrefix, "audio.language"))),
	)
}

func (c *Captcha) Generate() (key string, code string, err error) {
	key, codeContent, answer := c.driver.GenerateIdQuestionAnswer()
	item, err := c.driver.DrawCaptcha(codeContent)
	if err != nil {
		return
	}
	err = c.store.Set(key, answer)
	code = item.EncodeB64string()
	return
}

func (c *Captcha) Verify(key, answer string, clear bool) bool {
	vv := c.store.Get(key, clear)
	vv = strings.TrimSpace(vv)
	return vv == strings.TrimSpace(answer)
}

func New(opts ...Option) *Captcha {
	options := Options{
		Engine:               1,
		OptionShowSineLine:   false,
		OptionShowSlimeLine:  false,
		OptionShowHollowLine: false,
		Height:               80,
		Width:                240,
		CaptchaLen:           5,
		NoiseCount:           20,
		BgColor:              "#000000",
		Audio: Audio{
			Language: "en",
		},
	}
	for _, opt := range opts {
		opt(&options)
	}

	captcha := &Captcha{
		options: options,
	}

	if err := captcha.Init(); err != nil {
		panic(err)
	}

	return captcha
}

func (c *Captcha) Init() error {
	var (
		bgColor *color.RGBA
		err     error
	)
	if len(c.options.BgColor) > 0 {
		bgColor, err = utils.HexToRGBA(c.options.BgColor)
		if err != nil {
			return err
		}
	}
	switch c.options.Engine {
	case VerifyTypeAudio:
		c.driver = base64Captcha.NewDriverAudio(c.options.CaptchaLen, c.options.Audio.Language)
	case VerifyTypeCharacter:
		driver := base64Captcha.NewDriverString(
			c.options.Height,
			c.options.Width,
			c.options.NoiseCount,
			c.options.GetShowLine(),
			c.options.CaptchaLen,
			"",
			bgColor,
			base64Captcha.DefaultEmbeddedFonts,
			[]string{
				"3Dumb.ttf",
				"ApothecaryFont.ttf",
				"Comismsh.ttf",
				"DENNEthree-dee.ttf",
				"DeborahFancyDress.ttf",
				"Flim-Flam.ttf",
				"RitaSmith.ttf",
				"actionj.ttf",
				"chromohv.ttf",
				"wqy-microhei.ttc",
			})
		c.driver = driver.ConvertFonts()
	case VerifyTypeMath:
		c.driver = base64Captcha.NewDriverMath(
			c.options.Height,
			c.options.Width,
			c.options.NoiseCount,
			c.options.GetShowLine(),
			bgColor,
			nil,
			[]string{"3Dumb.ttf"},
		)
	case VerifyTypeChinese:
		c.driver = base64Captcha.NewDriverChinese(
			c.options.Height,
			c.options.Width,
			c.options.NoiseCount,
			c.options.GetShowLine(),
			c.options.CaptchaLen,
			"",
			bgColor,
			nil,
			nil,
		)
	case VerifyTypeDigit:
		c.driver = base64Captcha.DefaultDriverDigit
	case VerifyTypeLanguage:
		c.driver = base64Captcha.NewDriverLanguage(
			c.options.Height,
			c.options.Width,
			c.options.NoiseCount,
			c.options.GetShowLine(),
			c.options.CaptchaLen,
			bgColor,
			nil,
			nil,
			"en",
		)
	default:
		return fmt.Errorf("invalid captcha engine: %d", c.options.Engine)
	}
	return nil
}
