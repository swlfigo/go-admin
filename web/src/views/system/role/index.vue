<script setup lang="ts">
import { nextTick, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox, ElTree } from 'element-plus'
import { usePagedTable } from '@/composables/usePagedTable'
import * as roleApi from '@/api/role'
import { getMenuTree } from '@/api/menu'
import { toElTree, type ElTreeNode } from '@/utils/tree'
import type { Role } from '@/types/api'

const { t } = useI18n()
const { rows, total, loading, query, load, search, onPage } = usePagedTable<Role>(roleApi.listRoles)

const dialog = reactive({ visible: false, isEdit: false, id: 0 })
const form = reactive({ name: '', code: '', sort: 0, remark: '', status: 1 })

const permDialog = reactive({ visible: false, roleId: 0 })
const treeData = ref<ElTreeNode[]>([])
const checkedKeys = ref<number[]>([])
const treeRef = ref<InstanceType<typeof ElTree>>()

function openCreate() {
  Object.assign(form, { name: '', code: '', sort: 0, remark: '', status: 1 })
  dialog.isEdit = false; dialog.id = 0; dialog.visible = true
}
function openEdit(row: Role) {
  Object.assign(form, { name: row.name, code: row.code, sort: row.sort, remark: row.remark, status: row.status })
  dialog.isEdit = true; dialog.id = row.id; dialog.visible = true
}
async function submit() {
  try {
    const body = { name: form.name, code: form.code, sort: form.sort, remark: form.remark, status: form.status }
    if (dialog.isEdit) await roleApi.updateRole(dialog.id, body)
    else await roleApi.createRole(body)
    ElMessage.success(t('common.saveSuccess')); dialog.visible = false; load()
  } catch (e) { ElMessage.error((e as Error).message) }
}
async function remove(row: Role) {
  await ElMessageBox.confirm(t('role.confirmDelete', { name: row.name }), t('common.tip'), { type: 'warning' })
  try { await roleApi.deleteRole(row.id); ElMessage.success(t('common.deleted')); load() }
  catch (e) { ElMessage.error((e as Error).message) }
}
async function openPerm(row: Role) {
  if (!treeData.value.length) treeData.value = toElTree(await getMenuTree())
  const detail = await roleApi.getRole(row.id)
  permDialog.roleId = row.id
  permDialog.visible = true
  // 等弹窗+tree 渲染后回填勾选
  await nextTick()
  treeRef.value?.setCheckedKeys(detail.menus.map((m) => m.id))
}
async function submitPerm() {
  const tree = treeRef.value
  if (!tree) return
  const ids = [...(tree.getCheckedKeys() as number[]), ...(tree.getHalfCheckedKeys() as number[])]
  await roleApi.assignRoleMenus(permDialog.roleId, ids)
  ElMessage.success(t('role.permSaved')); permDialog.visible = false
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="page__head">
      <h1>{{ t('role.title') }}</h1>
      <div class="page__actions">
        <el-input v-model="query.keyword" :placeholder="t('role.searchPlaceholder')" style="width:200px" clearable @keyup.enter="search" />
        <el-button @click="search">{{ t('common.search') }}</el-button>
        <el-button v-permission="'system:role:create'" type="primary" @click="openCreate">{{ t('common.create') }}</el-button>
      </div>
    </div>

    <el-table :data="rows" v-loading="loading" border>
      <el-table-column prop="name" :label="t('role.name')" />
      <el-table-column prop="code" :label="t('role.code')" />
      <el-table-column prop="sort" :label="t('common.sort')" width="80" />
      <el-table-column prop="remark" :label="t('common.remark')" />
      <el-table-column :label="t('common.operations')" width="240">
        <template #default="{ row }: { row: Role }">
          <el-button v-permission="'system:role:update'" link type="primary" @click="openEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button v-permission="'system:role:assignmenu'" link type="primary" @click="openPerm(row)">{{ t('role.assignPerm') }}</el-button>
          <el-button v-permission="'system:role:delete'" link type="danger" @click="remove(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top:16px; justify-content:flex-end"
      layout="total, prev, pager, next"
      :total="total" :current-page="query.page" :page-size="query.size"
      @current-change="onPage"
    />

    <el-dialog v-model="dialog.visible" :title="dialog.isEdit ? t('role.editTitle') : t('role.createTitle')" width="440px">
      <el-form label-width="72px">
        <el-form-item :label="t('role.name')"><el-input v-model="form.name" /></el-form-item>
        <el-form-item :label="t('role.code')"><el-input v-model="form.code" :disabled="dialog.isEdit" /></el-form-item>
        <el-form-item :label="t('common.sort')"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item :label="t('common.remark')"><el-input v-model="form.remark" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="submit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="permDialog.visible" :title="t('role.assignPerm')" width="420px">
      <el-tree
        ref="treeRef" :data="treeData" show-checkbox node-key="id"
        :default-checked-keys="checkedKeys" :props="{ label: 'label', children: 'children' }"
      />
      <template #footer>
        <el-button @click="permDialog.visible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="submitPerm">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page__head{ display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; }
.page__head h1{ font-size:20px; margin:0; }
.page__actions{ display:flex; gap:8px; }
</style>
