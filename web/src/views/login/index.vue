<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCaptcha } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const user = useUserStore()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const form = reactive({ username: 'admin', password: 'admin123', captchaCode: '' })
const captchaId = ref('')
const captchaImg = ref('')
const captchaOn = ref(false)
const loading = ref(false)

async function loadCaptcha() {
  const res = await getCaptcha()
  captchaOn.value = res.enabled
  captchaId.value = res.captchaId ?? ''
  captchaImg.value = res.img ?? ''
}

async function onSubmit() {
  loading.value = true
  try {
    await user.login({
      username: form.username,
      password: form.password,
      captchaId: captchaId.value,
      captchaCode: form.captchaCode,
    })
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (e) {
    ElMessage.error((e as Error).message)
    if (captchaOn.value) loadCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(loadCaptcha)
</script>

<template>
  <div class="login">
    <div class="login__card">
      <div class="login__brand">
        <img src="/logo.svg" alt="Go Admin" width="46" height="46" />
        <div class="login__brandtext">
          <b>Go Admin</b>
          <small>{{ t('shell.brandTagline') }}</small>
        </div>
      </div>
      <el-form @submit.prevent="onSubmit">
        <el-form-item>
          <el-input v-model="form.username" :placeholder="t('login.username')" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" :placeholder="t('login.password')" show-password />
        </el-form-item>
        <el-form-item v-if="captchaOn">
          <div class="login__captcha">
            <el-input v-model="form.captchaCode" :placeholder="t('login.captcha')" />
            <img :src="captchaImg" alt="captcha" @click="loadCaptcha" />
          </div>
        </el-form-item>
        <el-button type="primary" :loading="loading" style="width:100%" @click="onSubmit">{{ t('login.submit') }}</el-button>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login{ height:100%; display:flex; align-items:center; justify-content:center; background: var(--bg); }
.login__card{ width:360px; padding:32px; background: var(--surface); border:1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-lg); }
.login__brand{ display:flex; align-items:center; gap:13px; justify-content:center; margin:0 0 26px; }
.login__brandtext b{ display:block; font-size:21px; font-weight:800; letter-spacing:-.02em; line-height:1; }
.login__brandtext small{ display:block; font-size:10.5px; font-weight:600; letter-spacing:.16em; text-transform:uppercase; color:var(--fg-subtle); margin-top:6px; }
.login__captcha{ display:flex; gap:10px; align-items:center; }
/* 验证码固定白底：图片透明背景在暗色下会透出深色卡片导致数字看不清，白底保证明暗都清晰可读 */
.login__captcha img{ height:38px; border-radius: var(--radius-sm); cursor:pointer; border:1px solid var(--border); background:#fff; }
</style>
