<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { usePagedTable } from '@/composables/usePagedTable'
import * as userApi from '@/api/user'
import { listRoles } from '@/api/role'
import type { Role, User } from '@/types/api'

const { t } = useI18n()
const { rows, total, loading, query, load, search, onPage } = usePagedTable<User>(userApi.listUsers)

const dialog = reactive({ visible: false, isEdit: false, id: 0 })
const form = reactive({ username: '', password: '', nickname: '', email: '', phone: '', status: 1 })
const roles = ref<Role[]>([])

const roleDialog = reactive({ visible: false, userId: 0, roleIds: [] as number[] })

function openCreate() {
  Object.assign(form, { username: '', password: '', nickname: '', email: '', phone: '', status: 1 })
  dialog.isEdit = false; dialog.id = 0; dialog.visible = true
}
function openEdit(row: User) {
  Object.assign(form, { username: row.username, password: '', nickname: row.nickname, email: row.email, phone: row.phone, status: row.status })
  dialog.isEdit = true; dialog.id = row.id; dialog.visible = true
}
async function submit() {
  try {
    if (dialog.isEdit) {
      await userApi.updateUser(dialog.id, { nickname: form.nickname, email: form.email, phone: form.phone, status: form.status })
    } else {
      await userApi.createUser({ username: form.username, password: form.password, nickname: form.nickname, email: form.email, phone: form.phone })
    }
    ElMessage.success(t('common.saveSuccess')); dialog.visible = false; load()
  } catch (e) { ElMessage.error((e as Error).message) }
}
async function remove(row: User) {
  await ElMessageBox.confirm(t('user.confirmDelete', { name: row.username }), t('common.tip'), { type: 'warning' })
  try { await userApi.deleteUser(row.id); ElMessage.success(t('common.deleted')); load() }
  catch (e) { ElMessage.error((e as Error).message) }
}
async function resetPwd(row: User) {
  const { value } = await ElMessageBox.prompt(t('user.newPassword'), t('user.resetPassword'), { inputType: 'password' })
  try { await userApi.resetUserPassword(row.id, value); ElMessage.success(t('user.passwordReset')); load() }
  catch (e) { ElMessage.error((e as Error).message) }
}
async function unlock(row: User) {
  try { await userApi.unlockUser(row.id); ElMessage.success(t('user.unlocked')); load() }
  catch (e) { ElMessage.error((e as Error).message) }
}
async function openRoles(row: User) {
  if (!roles.value.length) roles.value = (await listRoles({ page: 1, size: 100 })).list
  roleDialog.userId = row.id
  roleDialog.roleIds = row.roles.map((r) => r.id)
  roleDialog.visible = true
}
async function submitRoles() {
  await userApi.assignUserRoles(roleDialog.userId, roleDialog.roleIds)
  ElMessage.success(t('user.rolesAssigned')); roleDialog.visible = false; load()
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="page__head">
      <h1>{{ t('user.title') }}</h1>
      <div class="page__actions">
        <el-input v-model="query.keyword" :placeholder="t('user.searchPlaceholder')" style="width:200px" clearable @keyup.enter="search" />
        <el-button @click="search">{{ t('common.search') }}</el-button>
        <el-button v-permission="'system:user:create'" type="primary" @click="openCreate">{{ t('common.create') }}</el-button>
      </div>
    </div>

    <el-table :data="rows" v-loading="loading" border>
      <el-table-column prop="username" :label="t('user.username')" />
      <el-table-column prop="nickname" :label="t('user.nickname')" />
      <el-table-column :label="t('user.roles')">
        <template #default="{ row }: { row: User }">
          <el-tag v-for="r in row.roles" :key="r.id" style="margin-right:4px">{{ r.name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.status')">
        <template #default="{ row }: { row: User }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? t('common.enabled') : t('common.disabled') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.operations')" width="320">
        <template #default="{ row }: { row: User }">
          <el-button v-permission="'system:user:update'" link type="primary" @click="openEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button v-permission="'system:user:assignrole'" link type="primary" @click="openRoles(row)">{{ t('user.roles') }}</el-button>
          <el-button v-permission="'system:user:resetpwd'" link type="primary" @click="resetPwd(row)">{{ t('user.resetPassword') }}</el-button>
          <el-button v-permission="'system:user:update'" link type="warning" @click="unlock(row)">{{ t('user.unlock') }}</el-button>
          <el-button v-permission="'system:user:delete'" link type="danger" @click="remove(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top:16px; justify-content:flex-end"
      layout="total, prev, pager, next"
      :total="total" :current-page="query.page" :page-size="query.size"
      @current-change="onPage"
    />

    <el-dialog v-model="dialog.visible" :title="dialog.isEdit ? t('user.editTitle') : t('user.createTitle')" width="460px">
      <el-form label-width="72px">
        <el-form-item :label="t('user.username')"><el-input v-model="form.username" :disabled="dialog.isEdit" /></el-form-item>
        <el-form-item v-if="!dialog.isEdit" :label="t('user.password')"><el-input v-model="form.password" type="password" show-password /></el-form-item>
        <el-form-item :label="t('user.nickname')"><el-input v-model="form.nickname" /></el-form-item>
        <el-form-item :label="t('user.email')"><el-input v-model="form.email" /></el-form-item>
        <el-form-item :label="t('user.phone')"><el-input v-model="form.phone" /></el-form-item>
        <el-form-item v-if="dialog.isEdit" :label="t('common.status')">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="submit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleDialog.visible" :title="t('user.assignRoles')" width="400px">
      <el-checkbox-group v-model="roleDialog.roleIds">
        <el-checkbox v-for="r in roles" :key="r.id" :value="r.id" :label="r.name" />
      </el-checkbox-group>
      <template #footer>
        <el-button @click="roleDialog.visible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="submitRoles">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page__head{ display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; }
.page__head h1{ font-size:20px; margin:0; }
.page__actions{ display:flex; gap:8px; }
</style>
