package captcha

import "github.com/mojocn/base64Captcha"

type Store struct {
	driver base64Captcha.Driver
	store  base64Captcha.Store
}

func NewStore() *Store {
	// 数字算式验证码，内存 store（一期；二期可换 redis store）
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	return &Store{driver: driver, store: base64Captcha.DefaultMemStore}
}

// Generate 返回 captchaId、base64 图片、answer（answer 仅供测试）
func (s *Store) Generate() (id, b64, answer string, err error) {
	c := base64Captcha.NewCaptcha(s.driver, s.store)
	id, b64, answer, err = c.Generate()
	return
}

// Verify 校验后清除（一次性）
func (s *Store) Verify(id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	return s.store.Verify(id, answer, true)
}
