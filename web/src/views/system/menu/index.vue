<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as menuApi from '@/api/menu'
import type { MenuNode } from '@/types/api'

const { t } = useI18n()

const tree = ref<MenuNode[]>([])
const loading = ref(false)

const dialog = reactive({ visible: false, isEdit: false, id: 0 })
const form = reactive({ parentId: 0, name: '', type: 'C', path: '', component: '', perm: '', icon: '', sort: 0, visible: 1, status: 1 })

async function load() {
  loading.value = true
  try { tree.value = await menuApi.getMenuTree() } finally { loading.value = false }
}
function openCreate() {
  Object.assign(form, { parentId: 0, name: '', type: 'C', path: '', component: '', perm: '', icon: '', sort: 0, visible: 1, status: 1 })
  dialog.isEdit = false; dialog.id = 0; dialog.visible = true
}
function openChild(row: MenuNode) {
  openCreate(); form.parentId = row.id
}
function openEdit(row: MenuNode) {
  Object.assign(form, { parentId: row.parentId, name: row.name, type: row.type, path: row.path, component: row.component, perm: row.perm, icon: row.icon, sort: row.sort, visible: row.visible ?? 1, status: row.status ?? 1 })
  dialog.isEdit = true; dialog.id = row.id; dialog.visible = true
}
async function submit() {
  try {
    if (dialog.isEdit) await menuApi.updateMenu(dialog.id, { ...form })
    else await menuApi.createMenu({ ...form })
    ElMessage.success(t('common.saveSuccess')); dialog.visible = false; load()
  } catch (e) { ElMessage.error((e as Error).message) }
}
async function remove(row: MenuNode) {
  await ElMessageBox.confirm(t('menuMgmt.confirmDelete', { name: row.name }), t('common.tip'), { type: 'warning' })
  try { await menuApi.deleteMenu(row.id); ElMessage.success(t('common.deleted')); load() }
  catch (e) { ElMessage.error((e as Error).message) }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="page__head">
      <h1>{{ t('menuMgmt.title') }}</h1>
      <el-button v-permission="'system:menu:create'" type="primary" @click="openCreate">{{ t('menuMgmt.createRoot') }}</el-button>
    </div>
    <el-table :data="tree" v-loading="loading" row-key="id" border default-expand-all
      :tree-props="{ children: 'children' }">
      <el-table-column prop="name" :label="t('menuMgmt.name')" />
      <el-table-column :label="t('menuMgmt.type')" width="90">
        <template #default="{ row }: { row: MenuNode }">
          <el-tag size="small">{{ row.type === 'M' ? t('menuMgmt.typeM') : row.type === 'C' ? t('menuMgmt.typeC') : t('menuMgmt.typeF') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" :label="t('menuMgmt.path')" />
      <el-table-column prop="perm" :label="t('menuMgmt.perm')" />
      <el-table-column prop="sort" :label="t('common.sort')" width="80" />
      <el-table-column :label="t('common.operations')" width="220">
        <template #default="{ row }: { row: MenuNode }">
          <el-button v-permission="'system:menu:create'" link type="primary" @click="openChild(row)">{{ t('menuMgmt.addChild') }}</el-button>
          <el-button v-permission="'system:menu:update'" link type="primary" @click="openEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button v-permission="'system:menu:delete'" link type="danger" @click="remove(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialog.visible" :title="dialog.isEdit ? t('menuMgmt.editTitle') : t('menuMgmt.createTitle')" width="480px">
      <el-form label-width="84px">
        <el-form-item :label="t('menuMgmt.parentId')"><el-input-number v-model="form.parentId" :min="0" /></el-form-item>
        <el-form-item :label="t('menuMgmt.name')"><el-input v-model="form.name" /></el-form-item>
        <el-form-item :label="t('menuMgmt.type')">
          <el-select v-model="form.type">
            <el-option :label="t('menuMgmt.typeM')" value="M" /><el-option :label="t('menuMgmt.typeC')" value="C" /><el-option :label="t('menuMgmt.typeF')" value="F" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('menuMgmt.path')"><el-input v-model="form.path" /></el-form-item>
        <el-form-item :label="t('menuMgmt.component')"><el-input v-model="form.component" :placeholder="t('menuMgmt.componentPlaceholder')" /></el-form-item>
        <el-form-item :label="t('menuMgmt.perm')"><el-input v-model="form.perm" :placeholder="t('menuMgmt.permPlaceholder')" /></el-form-item>
        <el-form-item :label="t('common.sort')"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item :label="t('menuMgmt.visible')"><el-switch v-model="form.visible" :active-value="1" :inactive-value="0" /></el-form-item>
        <el-form-item :label="t('common.status')"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="submit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page__head{ display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; }
.page__head h1{ font-size:20px; margin:0; }
</style>
