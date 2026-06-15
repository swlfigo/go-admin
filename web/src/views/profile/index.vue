<script setup lang="ts">
import { reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import * as profileApi from '@/api/profile'

const { t } = useI18n()

const user = useUserStore()
const info = reactive({
  nickname: user.user?.nickname ?? '',
  email: user.user?.email ?? '',
  phone: user.user?.phone ?? '',
})
const pwd = reactive({ oldPassword: '', newPassword: '' })

async function saveInfo() {
  try {
    await profileApi.updateProfile({ nickname: info.nickname, email: info.email, phone: info.phone })
    await user.loadMe()
    ElMessage.success(t('profile.infoSaved'))
  } catch (e) { ElMessage.error((e as Error).message) }
}
async function savePwd() {
  try {
    await profileApi.changePassword({ oldPassword: pwd.oldPassword, newPassword: pwd.newPassword })
    pwd.oldPassword = ''; pwd.newPassword = ''
    ElMessage.success(t('profile.passwordChanged'))
  } catch (e) { ElMessage.error((e as Error).message) }
}
</script>

<template>
  <div class="profile">
    <el-card :header="t('profile.basicInfo')">
      <el-form label-width="72px" style="max-width:420px">
        <el-form-item :label="t('profile.username')"><el-input :model-value="user.user?.username" disabled /></el-form-item>
        <el-form-item :label="t('profile.nickname')"><el-input v-model="info.nickname" /></el-form-item>
        <el-form-item :label="t('profile.email')"><el-input v-model="info.email" /></el-form-item>
        <el-form-item :label="t('profile.phone')"><el-input v-model="info.phone" /></el-form-item>
        <el-button type="primary" @click="saveInfo">{{ t('profile.saveInfo') }}</el-button>
      </el-form>
    </el-card>
    <el-card :header="t('profile.changePassword')" style="margin-top:16px">
      <el-form label-width="72px" style="max-width:420px">
        <el-form-item :label="t('profile.oldPassword')"><el-input v-model="pwd.oldPassword" type="password" show-password /></el-form-item>
        <el-form-item :label="t('profile.newPassword')"><el-input v-model="pwd.newPassword" type="password" show-password /></el-form-item>
        <el-button type="primary" @click="savePwd">{{ t('profile.changePassword') }}</el-button>
      </el-form>
    </el-card>
  </div>
</template>
