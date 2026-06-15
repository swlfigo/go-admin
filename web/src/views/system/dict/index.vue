<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as dictApi from '@/api/dict'
import type { DictData, DictType } from '@/types/api'

const { t } = useI18n()

const types = ref<DictType[]>([])
const activeType = ref<string>('')
const data = ref<DictData[]>([])

const typeDialog = reactive({ visible: false, isEdit: false, id: 0 })
const typeForm = reactive({ name: '', type: '', remark: '', status: 1 })
const dataDialog = reactive({ visible: false, isEdit: false, id: 0 })
const dataForm = reactive({ label: '', value: '', tagType: '', sort: 0, status: 1 })

async function loadTypes() {
  types.value = (await dictApi.listDictTypes({ page: 1, size: 100 })).list
  if (!activeType.value && types.value.length) selectType(types.value[0].type)
}
async function selectType(t: string) {
  activeType.value = t
  data.value = await dictApi.listDictData(t)
}
function openType() { Object.assign(typeForm, { name: '', type: '', remark: '', status: 1 }); typeDialog.isEdit = false; typeDialog.id = 0; typeDialog.visible = true }
function editType(row: DictType) { Object.assign(typeForm, { name: row.name, type: row.type, remark: row.remark, status: row.status }); typeDialog.isEdit = true; typeDialog.id = row.id; typeDialog.visible = true }
async function submitType() {
  try {
    if (typeDialog.isEdit) await dictApi.updateDictType(typeDialog.id, { name: typeForm.name, status: typeForm.status, remark: typeForm.remark })
    else await dictApi.createDictType({ name: typeForm.name, type: typeForm.type, remark: typeForm.remark })
    ElMessage.success(t('common.saveSuccess')); typeDialog.visible = false; loadTypes()
  } catch (e) { ElMessage.error((e as Error).message) }
}
async function removeType(row: DictType) {
  await ElMessageBox.confirm(t('dict.confirmDeleteType', { name: row.name }), t('common.tip'), { type: 'warning' })
  try { await dictApi.deleteDictType(row.id); ElMessage.success(t('common.deleted')); activeType.value = ''; loadTypes() }
  catch (e) { ElMessage.error((e as Error).message) }
}
function openData() { Object.assign(dataForm, { label: '', value: '', tagType: '', sort: 0, status: 1 }); dataDialog.isEdit = false; dataDialog.id = 0; dataDialog.visible = true }
function editData(row: DictData) { Object.assign(dataForm, { label: row.label, value: row.value, tagType: row.tagType, sort: row.sort, status: row.status }); dataDialog.isEdit = true; dataDialog.id = row.id; dataDialog.visible = true }
async function submitData() {
  try {
    if (dataDialog.isEdit) await dictApi.updateDictData(dataDialog.id, { label: dataForm.label, value: dataForm.value, tagType: dataForm.tagType, sort: dataForm.sort, status: dataForm.status })
    else await dictApi.createDictData({ dictType: activeType.value, label: dataForm.label, value: dataForm.value, tagType: dataForm.tagType, sort: dataForm.sort })
    ElMessage.success(t('common.saveSuccess')); dataDialog.visible = false; selectType(activeType.value)
  } catch (e) { ElMessage.error((e as Error).message) }
}
async function removeData(row: DictData) {
  await ElMessageBox.confirm(t('dict.confirmDeleteData', { name: row.label }), t('common.tip'), { type: 'warning' })
  try { await dictApi.deleteDictData(row.id); ElMessage.success(t('common.deleted')); selectType(activeType.value) }
  catch (e) { ElMessage.error((e as Error).message) }
}

onMounted(loadTypes)
</script>

<template>
  <div class="dict">
    <div class="dict__col">
      <div class="page__head"><h2>{{ t('dict.typeTitle') }}</h2><el-button v-permission="'system:dict:create'" size="small" type="primary" @click="openType">{{ t('common.create') }}</el-button></div>
      <el-table :data="types" border highlight-current-row @row-click="(r: DictType) => selectType(r.type)">
        <el-table-column prop="name" :label="t('dict.name')" />
        <el-table-column prop="type" :label="t('dict.typeKey')" />
        <el-table-column :label="t('common.operations')" width="120">
          <template #default="{ row }: { row: DictType }">
            <el-button v-permission="'system:dict:update'" link type="primary" @click.stop="editType(row)">{{ t('common.edit') }}</el-button>
            <el-button v-permission="'system:dict:delete'" link type="danger" @click.stop="removeType(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="dict__col">
      <div class="page__head"><h2>{{ t('dict.dataTitle') }} <span v-if="activeType" class="muted">/ {{ activeType }}</span></h2><el-button v-permission="'system:dict:create'" size="small" type="primary" :disabled="!activeType" @click="openData">{{ t('common.create') }}</el-button></div>
      <el-table :data="data" border>
        <el-table-column prop="label" :label="t('dict.label')" />
        <el-table-column prop="value" :label="t('dict.value')" />
        <el-table-column prop="sort" :label="t('common.sort')" width="80" />
        <el-table-column :label="t('common.operations')" width="120">
          <template #default="{ row }: { row: DictData }">
            <el-button v-permission="'system:dict:update'" link type="primary" @click="editData(row)">{{ t('common.edit') }}</el-button>
            <el-button v-permission="'system:dict:delete'" link type="danger" @click="removeData(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="typeDialog.visible" :title="typeDialog.isEdit ? t('dict.typeEditTitle') : t('dict.typeCreateTitle')" width="400px">
      <el-form label-width="72px">
        <el-form-item :label="t('dict.name')"><el-input v-model="typeForm.name" /></el-form-item>
        <el-form-item :label="t('dict.typeKey')"><el-input v-model="typeForm.type" :disabled="typeDialog.isEdit" :placeholder="t('dict.typeKeyPlaceholder')" /></el-form-item>
        <el-form-item :label="t('common.remark')"><el-input v-model="typeForm.remark" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="typeDialog.visible=false">{{ t('common.cancel') }}</el-button><el-button type="primary" @click="submitType">{{ t('common.save') }}</el-button></template>
    </el-dialog>

    <el-dialog v-model="dataDialog.visible" :title="dataDialog.isEdit ? t('dict.dataEditTitle') : t('dict.dataCreateTitle')" width="400px">
      <el-form label-width="72px">
        <el-form-item :label="t('dict.label')"><el-input v-model="dataForm.label" /></el-form-item>
        <el-form-item :label="t('dict.value')"><el-input v-model="dataForm.value" /></el-form-item>
        <el-form-item :label="t('common.sort')"><el-input-number v-model="dataForm.sort" :min="0" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dataDialog.visible=false">{{ t('common.cancel') }}</el-button><el-button type="primary" @click="submitData">{{ t('common.save') }}</el-button></template>
    </el-dialog>
  </div>
</template>

<style scoped>
.dict{ display:grid; grid-template-columns:1fr 1fr; gap:16px; }
.page__head{ display:flex; justify-content:space-between; align-items:center; margin-bottom:12px; }
.page__head h2{ font-size:16px; margin:0; }
.muted{ color: var(--fg-muted); font-weight:400; font-size:13px; }
</style>
