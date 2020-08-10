<template>
  <div class="explorer card" :class="{ 'card-fullscreen': isFullScreen }">
    <div v-if="data !== null" class="card-header">
      <h2 ref="tTitle" class="card-title">
        <span class="exam-label" @click="$refs.examSelect.show($event.target)">{{ data.examConf.Label || data.examName }}</span> - {{ data.dataDesc }}
        <span
          style="font-size: 13px;vertical-align: bottom;"
        >[页码 {{ data.page }}/{{ data.lastPage }}]</span>
      </h2>
      <SelectFloater
        ref="examSelect"
        :options="$app.ExamNameToLabelObj || {}"
        :selected="curtExamName"
        :change="switchExam" />
      <small class="card-subtitle">{{ data.total }} 人</small>
      <!-- Actions -->
      <div class="actions">
        <span class="item show-top-badge" :class="{ 'active': !!$searchLayer.isShow }" @click="$searchLayer.toggle()">
          <i class="zmdi zmdi-search" />
          <span>搜索</span>
        </span>
        <span class="item" @click="$refs.tableDataDownloadDialog.show()">
          <i class="zmdi zmdi-download" />
          <span>下载</span>
        </span>
        <span class="item" @click="printTable()">
          <i class="zmdi zmdi-print" />
          <span>打印</span>
        </span>
        <span class="item" @click="$refs.tableDataCounterDialog.show()">
          <i class="zmdi zmdi-flash" />
          <span>平均分</span>
        </span>
        <span v-if="fieldRankOn" class="item active anim-fade-in" @click="$refs.settingDialog.show()">
          <i class="zmdi zmdi-swap-vertical" />
          <span>单科排名 [{{ FieldRankTypeNameDict[fieldRankType] }}视角]</span>
        </span>
        <span class="item" @click="$refs.settingDialog.show()">
          <i class="zmdi zmdi-format-paint" />
          <span>设置</span>
        </span>
        <span class="item" @click="toggleFullScreen()">
          <i :class="`zmdi zmdi-fullscreen${isFullScreen ? '-exit' : ''}`" />
          <span>{{ !isFullScreen ? '全屏显示' : '退出全屏' }}</span>
        </span>
      </div>
    </div>

    <div v-if="data !== null" ref="tWrap" class="score-table-wrap" :class="{'field-rank-on': fieldRankOn}" style="padding: 0;">
      <div
        ref="tContainer"
        class="wly-table-container"
        data-toggle="wlyTable"
      >
        <!-- Table -->
        <div ref="tHeader" class="wly-table-header">
          <table class="table table-striped table-hover">
            <thead>
              <tr>
                <th
                  v-for="(f, i) in ViewFieldList"
                  :key="i"
                  @click="switchSort(f)"
                >
                  <span
                    :class="getFieldItemClass(f)"
                    :title="getFieldItemHoverTitle(f)"
                  >{{ transField(f) }}</span>
                </th>
              </tr>
            </thead>
          </table>
        </div>
        <div ref="tBody" class="wly-table-body">
          <table class="table table-striped table-hover">
            <thead>
              <tr>
                <th
                  v-for="(f, i) in ViewFieldList"
                  :key="i"
                >{{ transField(f) }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, i) in data.list" :key="i" class="table-item">
                <th v-for="f in ViewFieldList" :key="f">
                  <span v-if="f === 'NAME'" class="clickable-text" @click="goChart(item)">{{ item[f] }}</span>
                  <span v-else>{{ item[f] }}</span>
                  <span v-if="fieldRankOn && TargetRankField.includes(f)" class="field-rank-print print-only">[{{ getItemFieldRank(item, f) }}]</span>
                  <span
                    v-if="fieldRankOn && TargetRankField.includes(f)"
                    class="field-rank anim-fade-in"
                    :title="getItemFieldRankHoverTitle(item, f)"
                    @click="fieldRankClickSwitch()"
                  >{{ getItemFieldRank(item, f) }}</span>
                </th>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div ref="tPagination" class="wly-table-pagination">
        <div class="paginate-simple">
          <a
            v-if="!visiblePageBtn.includes(1)"
            class="paginate-button"
            title="第一页"
            @click="switchPage(1)"
          >1</a>
          <a
            :class="{ disabled: data.page-1 <= 0 }"
            class="paginate-button previous"
            title="上一页"
            @click="switchPage(data.page-1)"
          ></a>
          <span>
            <a
              v-for="(pageNum, i) in visiblePageBtn"
              :key="i"
              :class="{ current: pageNum === data.page }"
              class="paginate-button"
              @click="switchPage(pageNum)"
            >{{ pageNum }}</a>
          </span>
          <a
            :class="{ disabled: data.page+1 > data.lastPage }"
            class="paginate-button next"
            title="下一页"
            @click="switchPage(data.page+1)"
          ></a>
          <a
            v-if="!visiblePageBtn.includes(data.lastPage)"
            class="paginate-button"
            title="最后一页"
            @click="switchPage(data.lastPage)"
          >{{ data.lastPage }}</a>
        </div>
      </div>
    </div>

    <LoadingLayer ref="tLoading" />

    <ExplorerDialog ref="settingDialog" title="设置">
      <div v-if="data !== null" class="table-ctrl-dialog">
          <span class="dialog-label">点按下列方块来 显示 / 隐藏 字段</span>
          <div class="field-list">
            <span
              v-for="(f, i) in FieldList"
              :key="i"
              :class="{ 'active': HideFieldList.includes(f) }"
              class="checkbox"
              @click="toggleFieldView(f)"
            >{{ transField(f) }}</span>
            <div
              class="checkbox"
              :class="{ 'active': !fieldRankOn }"
              @click="() => { fieldRankOn = !fieldRankOn }"
            >单科排名</div>
          </div>
          <div v-if="fieldRankOn" class="anim-fade-in">
            <span class="dialog-label">单科排名 视角</span>
            <div
              v-for="(name, type) in FieldRankTypeNameDict"
              :key="type"
              class="checkbox"
              :class="{ 'active': fieldRankType === type }"
              @click="setFieldRankType(type)"
            >{{ name }}</div>
          </div>
          <span class="dialog-label">每页显示项目数量 （数字不宜过大）</span>
          <div class="page-per-show">
            <input type="number" class="page-per-show-input" placeholder="每页显示数" min="1" :value="data.pageSize" />
          </div>
          <span class="dialog-label">表格字体大小调整</span>
          <div class="table-font-size-control">
            <span class="font-size-minus" @click="adjustTableFontSize(-2)">-</span>
            <span ref="tFontSize" class="font-size-value">15</span>
            <span class="font-size-plus" @click="adjustTableFontSize(+2)">+</span>
          </div>
      </div>
    </ExplorerDialog>

    <ExplorerDialog ref="tableDataDownloadDialog" title="保存数据为电子表格">
      <span v-if="data !== null">
        <span class="dialog-label">保存 "{{ data.examConf.Label }}" 的 "{{ data.dataDesc }}" 为电子表格</span>
        <span class="dialog-btn" data-dialog-func="save-now">保存 仅第 {{ data.page }} 页 数据</span>
        <span class="dialog-btn" data-dialog-func="save-now-noPaging">保存 第 1~{{ data.lastPage }} 页 数据</span>
        <span class="dialog-label">保存 "{{ data.examConf.Label }}" 的全部数据为电子表格</span>
        <span class="dialog-btn" data-dialog-func="save-noPaging">保存 全部成绩</span>
      </span>
    </ExplorerDialog>

    <ExplorerDialog ref="tableDataCounterDialog" title="数据统计">
      <div v-if="data !== null" class="table-data-counter">
      <span class="dialog-label">数据 "{{ data.dataDesc }}" 平均值</span>
      <span v-for="(avg, f) in data.avgList" :key="f" class="data-item">
        <span class="data-name">{{ transField(f) }}</span>
        <span class="data-value">{{ avg }}</span>
      </span>
      </div>
    </ExplorerDialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from 'nuxt-property-decorator'
import { Persist } from 'vue-local-storage-decorator'
import $ from 'jquery'
import _ from 'lodash'
import F, { ScoreData } from '../types/Field'
import * as FG from '../types/Field/Grp'
import * as ApiT from '../types/ApiTypes'
import LoadingLayer from './LoadingLayer.vue'
import SelectFloater from './SelectFloater.vue'
import ExplorerDialog from './ExplorerDialog.vue'

type FIELD_RANK_TYPE = 'all'|'class'|'school'

@Component({
  name: 'explorer',
  components: { LoadingLayer, ExplorerDialog, SelectFloater }
})
export default class Explorer extends Vue {
  data: ApiT.QueryData | null = null
  params: ApiT.QueryParams | null = null
  isFullScreen = false
  loading!: LoadingLayer

  created () {
    Vue.prototype.$explorer = this
  }

  mounted () {
    this.loading = this.$refs.tLoading as LoadingLayer
    $(window).resize(() => {
      this.adjustDisplay()
    })
  }

  /** 当前 Exam 名 */
  get curtExamName () {
    if (!this.data)
      return this.params ? this.params.exam || null : null
    return this.data.examConf.Name || null
  }

  /** 当前 Exam 标签 */
  get curtExamLabel () {
    return this.data ? this.data.examConf.Label || null : null
  }

  /** 隐藏的字段 */
  HideFieldList: F[] = [F.SCHOOL_RANK, F.CLASS_RANK,...FG.F_EXT_SUM]

  /** 全部字段 */
  get FieldList () {
    if (!this.data) return []
    const rawFieldList = this.data.fieldList // 源字段

    // 允许出现的字段
    let F_MAIN = FG.F_MAIN
    if (this.isClassSearchMode) {
      F_MAIN = FG.F_MAIN_CLASS
    } else if (this.isSchoolSearchMode) {
      F_MAIN = FG.F_MAIN_SCHOOL
    }

    const allowFields = [...F_MAIN, ...FG.F_SUBJ, ...FG.F_EXT_SUM]

    // 构建有序的字段 & 排除字段
    const fieldList: F[] = []
    _.forEach(allowFields, (fieldName) => {
      if (rawFieldList.includes(fieldName))
        fieldList.push(fieldName)
    })

    return fieldList
  }

  get TargetRankField () {
    return FG.F_TARGET_RANK
  }

  get paramsWhereObj (): {[f in F]?: any} {
    if (!this.params) return {}
    if (!this.params.where) return {}
    try {
      return JSON.parse(this.params.where)
    } catch {
      return {}
    }
  }

  /** 显示的字段 */
  get ViewFieldList () {
    return _.filter(this.FieldList, (f) => !this.HideFieldList.includes(f))
  }

  toggleFieldView (field: F) {
    const isShow = !this.HideFieldList.includes(field)
    this.setFieldView(field, !isShow)
  }

  setFieldView (field: F|F[]|{[f in F]?: boolean}, show?: boolean) {
    const setOneFieldView = (field: F, show: boolean) => {
      if (show && this.HideFieldList.includes(field)) {
        this.HideFieldList.splice(this.HideFieldList.indexOf(field), 1)
      }
      if (!show && !this.HideFieldList.includes(field)) {
        this.HideFieldList.push(field)
      }
    }
    if (_.isArray(field) && show !== undefined) {
      _.forEach(field, (f) => { setOneFieldView(f, show) })
    } else if (_.isObject(field)) {
      _.forEach(field, (show, f) => { setOneFieldView(f as F, show) })
    } else if(show !== undefined) setOneFieldView(field, show)
    this.$nextTick(() => {
      this.adjustDisplay()
    })
  }

  @Watch('data')
  onDataChanged () {
    this.$nextTick(() => {
      this.adjustDisplay()
      $(this.$refs.tBody).bind('scroll.table-scroll-sync', (e) => {
        $(this.$refs.tHeader).scrollLeft(($(e.target) as any).scrollLeft())
        $(this.$refs.tHeader).scrollTop(($(e.target) as any).scrollTop())
        this.setFullScreen(true)
      })
    })

    // if (this.data === null) return
  }

  @Watch('$route.query')
  public async onRouteQueryChanged (query: object) {
    if (query === this.params) return

    this.loading.show()
    this.params = query
    if (!this.params.exam) {
      this.params.init = true
    }

    let resp:any
    try {
      resp = await this.$axios.$get('./api/query', {
        params: this.params
      })
    } catch (err) {
      this.$notify.error(String(err))
    } finally {
      this.loading.hide()
    }

    if (resp.success && !!resp.data) {
      this.data = resp.data
      // 初始化配置装载
      if (!!this.data && !!this.data.initConf) {
        if (this.params) {
          delete this.params.init // 删除初始化请求参数
          this.params.exam = this.data.examConf.Name
        }
        this.$app.Conf = this.data.initConf
      }

      // 自动设置字段 显示/隐藏
      const whereObj = this.paramsWhereObj
      if (this.isClassSearchMode) {
        this.setFieldView({
          [F.SCHOOL]: false, [F.SCHOOL_RANK]: false,
          [F.CLASS]: false, [F.CLASS_RANK]: true
        })
        this.fieldRankType = 'class'
      } else if (this.isSchoolSearchMode) {
        this.setFieldView({
          [F.SCHOOL]: false, [F.SCHOOL_RANK]: true,
          [F.CLASS]: true, [F.CLASS_RANK]: false
        })
        this.fieldRankType = 'school'
      } else {
        this.setFieldView({
          [F.SCHOOL]: true, [F.SCHOOL_RANK]: false,
          [F.CLASS]: true, [F.CLASS_RANK]: false
        })
        this.fieldRankType = 'all'
      }
    } else {
      this.$notify.error(resp.msg)
    }
  }

  fetchData (params: ApiT.QueryParams, initialize = false) {
    const reqParams: ApiT.QueryParams = !initialize
      ? { ...this.params, ...params }
      : params
    this.$router.push({ name: 'index', query: reqParams as any })
  }

  switchExam (examName: string, initialize = false) {
    this.fetchData({ exam: examName, page: 1 }, initialize)
  }

  switchPage (pageNum: number) {
    if (!this.data || pageNum <= 0 || pageNum > this.data.lastPage) return
    this.fetchData({ page: pageNum })
  }

  switchSort (fieldName: F) {
    if (!this.data) return
    this.fetchData({
      page: 1,
      sort: JSON.stringify({
        [fieldName]: this.data.sortList[fieldName] === -1 ? 1 : -1
      })
    })
  }

  get isNormalSearchMode () {
    const whereObj = this.paramsWhereObj
    return (!whereObj.CLASS && !whereObj.SCHOOL)
  }

  get isClassSearchMode () {
    const whereObj = this.paramsWhereObj
    return (!!whereObj.CLASS && !!whereObj.SCHOOL)
  }

  get isSchoolSearchMode () {
    return !!this.paramsWhereObj.SCHOOL
  }

  get visiblePageBtn () {
    if (this.data === null) return []
    const arr: number[] = []
    const lItemNum = 3
    const rItemNum = 3
    for (let i = lItemNum; i > 0; i--) {
      const pg = this.data.page - i
      if (pg > 0) arr.push(pg)
    }
    for (let i = 0; i < rItemNum + 1; i++) {
      const pg = this.data.page + i
      if (pg <= this.data.lastPage) arr.push(pg)
    }
    return arr
  }

  setFullScreen (val: boolean) {
    if (this.isFullScreen === val) return
    if (val === true) {
      this.$topHeader.hide()
      // this.browserRequestFullScreen(document.body)
    } else {
      this.$topHeader.show()
      // this.browserCancelFullScreen(document)
    }

    this.isFullScreen = val
    this.$nextTick(() => {
      this.adjustDisplay()
    })
  }

  browserRequestFullScreen (el: any) {
    const requestMethod = el.requestFullScreen || el.webkitRequestFullScreen || el.mozRequestFullScreen || el.msRequestFullScreen

    if (requestMethod) { // Native full screen.
      requestMethod.call(el)
    } else if (typeof (window as any).ActiveXObject !== "undefined") { // Older IE
      const wscript = new (window as any).ActiveXObject("WScript.Shell")
      if (wscript !== null) {
        wscript.SendKeys("{F11}")
      }
    }
  }

  browserCancelFullScreen (el: any) {
    const requestMethod = el.cancelFullScreen || el.webkitCancelFullScreen || el.mozCancelFullScreen || el.exitFullscreen
    if (requestMethod) { // cancel full screen.
      requestMethod.call(el)
    } else if (typeof (window as any).ActiveXObject !== "undefined") { // Older IE.
      const wscript = new (window as any).ActiveXObject("WScript.Shell")
      if (wscript !== null) {
        wscript.SendKeys("{F11}")
      }
    }
  }

  toggleFullScreen () {
    this.setFullScreen(!this.isFullScreen)
  }

  adjustDisplay () {
    const wrapEl = $(this.$refs.tWrap)
    const containerEl = $(this.$refs.tContainer)
    const headerEl = $(this.$refs.tHeader)
    const headerTableEl = headerEl.find('table')
    const bodyEl = $(this.$refs.tBody)
    const bodyTableEl = bodyEl.find('table')
    const paginationEl = $(this.$refs.tPagination)
    // 设置悬浮样式
    const curtHeadHeight = (bodyEl.find('table thead').outerHeight() || 1) + 2
    bodyTableEl.css(
      'margin-top',
      `-${curtHeadHeight}px`
    )
    bodyEl.css('height', `calc(100% - ${curtHeadHeight}px)`)

    // 获取 body table thead tr 中每个 th 对象
    const bodyThItems = bodyTableEl.find(
      '> thead > tr:first-child:not(.no-records-found) > *'
    )
    $.each(bodyThItems, (i: number, item: any) => {
      // 逐个设置 head table 中每个 th 的宽度 === body th 的宽度
      headerTableEl
        .find(`> thead th:nth-child(${Number(i) + 1})`)
        .width($(item).width() || '')
    })
    headerTableEl.width((bodyTableEl.outerWidth(true) || 0) - 2) // minus the 2px border-width
  }

  tableFontSize: number = 0

  adjustTableFontSize (number: number) {
    const curtFontSize = this.tableFontSize || Number($('.table').first().css('font-size').replace(/px$/, ''))
    const calcSize = curtFontSize + number
    if (calcSize <= 1) return
    this.tableFontSize = calcSize
    let style = $('#TableFontSize')
    if (!style.length) {
      style = $('<style id="TableFontSize"></style>').appendTo('head')
    }
    style.html(`.table {font-size: ${this.tableFontSize}px !important;}`);
    (this.$refs.tFontSize as HTMLElement).innerHTML = String(this.tableFontSize)
    this.adjustDisplay()
  }

  transField (f: F) {
    return this.$app.transField(f)
  }

  getFieldItemClass (fieldName: F) {
    if (!this.data) return ''
    const sortType = this.data.sortList[fieldName]
    if (typeof sortType !== 'number') return ''
    return sortType === 1 ? 'sort-asc' : 'sort-desc'
  }

  getFieldItemHoverTitle (f: F) {
    if (!this.data) return ''
    const sortType = this.data.sortList[f]
    let title = `依 ${this.transField(f)} ${sortType === -1 ? '升序' : '降序'} `
    if (sortType !== undefined)
      title += `[当前为 ${sortType === -1 ? '降序' : '升序'}]`
    return title
  }

  goChart (item: ScoreData) {
    if (this.data === null) return

    const query: ApiT.ChartParams = {
      examGrp: this.data.examConf.Grp,
      where: JSON.stringify({
        NAME: item.NAME,
        SCHOOL: item.SCHOOL,
        CLASS: item.CLASS
      } as {[k in F]?: string})
    }

    this.$router.push({
      name: 'chart',
      query: query as any
    })
  }

  printTable () {
    ($(this.$refs.tBody).find('table').css('margin-top', '') as any).print({
      globalStyles: true,
      mediaPrint: false,
      iframe: false,
      prepend: `
      <style>html, body {background-color: transparent !important;}*{font-size: 10px !important;}</style>
      <h2 style="text-align: center;margin-bottom: 20px">${$(this.$refs.tTitle).text()}</h2>
      `
    })
    this.adjustDisplay()
  }

  @Persist()
  fieldRankOn = true
  fieldRankType: FIELD_RANK_TYPE = 'all'
  FieldRankTypeNameDict: {[key in FIELD_RANK_TYPE]: string} = {'all': '总', 'class': '班级', 'school': '学校'}
  get FieldRankTypes () { return Object.keys(this.FieldRankTypeNameDict) as FIELD_RANK_TYPE[] }

  @Watch('fieldRankOn')
  onFieldRankOnChanged (isOn: boolean) {
    if (isOn === true) { // 打开初始化
      if (this.isClassSearchMode) {
        this.fieldRankType = 'class'
      } else if (this.isSchoolSearchMode) {
        this.fieldRankType = 'school'
      } else {
        this.fieldRankType = 'all'
      }
    }
    this.$nextTick(() => {
      this.adjustDisplay()
    })
  }

  setFieldRankType (type: FIELD_RANK_TYPE) {
    this.fieldRankType = type
    this.$notify.clearAll()
    this.$notify.info(`单科排名 已调整为 “${this.FieldRankTypeNameDict[type]}视角”`)
  }

  fieldRankClickSwitch () {
    const allTypes = this.FieldRankTypes
    let nxtTypeIndex = allTypes.indexOf(this.fieldRankType) + 1
    if (nxtTypeIndex >= allTypes.length) nxtTypeIndex = 0
    this.setFieldRankType(this.FieldRankTypes[nxtTypeIndex])
  }

  getItemFieldRank (item: ScoreData, f: F) {
    if (!FG.F_TARGET_RANK.includes(f)) return null
    let rankF
    if (this.fieldRankType === 'all') {
      rankF = `${f}_RANK`
    } else if (this.fieldRankType === 'class') {
      rankF = `${f}_CLASS_RANK`
    } else if (this.fieldRankType === 'school') {
      rankF = `${f}_SCHOOL_RANK`
    }
    // School & Class Field
    if (f === F.SCHOOL) {
      rankF = `SCHOOL_RANK`
    } else if (f === F.CLASS) {
      rankF = `CLASS_RANK`
    }

    return item[rankF] ? Number(item[rankF]) : null
  }

  getItemFieldRankHoverTitle (item: ScoreData, f: F) {
    return `“${item.NAME}” 的 “${this.transField(f)}” 在 “${this.FieldRankTypeNameDict[this.fieldRankType]}成绩” 中排 ${this.getItemFieldRank(item, f)}名`
  }
}
</script>

<style scoped lang="scss">
.explorer {
  .card-title {
    .exam-label {
      cursor: pointer;
      &:hover {
        color: var(--mainColor)
      }
    }
  }
}

.card-fullscreen {
  .score-table-wrap {
    height: calc(100vh - #{87px});
  }
}

/* table */
.score-table-wrap {
  height: calc(100vh - #{55px+87px+15px*2});
  display: flex;
  flex-direction: column;
}

.wly-table-container {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.wly-table-header {
  overflow: hidden;
  /* margin-right: 17px; */

  table {
    margin-bottom: 0;
  }
  thead {
    overflow: hidden;

    th span {
      cursor: pointer;

      &.select {
        color: var(--mainColor);
      }

      &.sort-desc,
      &.sort-asc {
        color: var(--mainColor);

        &:after {
          font-family: Material-Design-Iconic-Font;
          position: absolute;
          width: 20px;
        }
      }

      &.sort-desc:after {
        content: '\f2fe';
      }
      &.sort-asc:after {
        content: '\f303';
      }
    }
  }
}

@media screen and (max-width: 559px) {
  .wly-table-header thead th span.select:after {
    position: initial;
  }
}

@media print {
  .wly-table-header thead th span.select {
    color: #fff;
  }

  .wly-table-header thead th span:after {
    display: none;
  }
}

.wly-table-pagination {
  overflow: hidden;
  background: #fff;
  height: 55px;
  display: flex;
  justify-content: center;
  align-items: center;
}

table {
  border-spacing: 0;
  border-collapse: collapse;
}

.field-rank-on {
  table tr th {
    padding-right: 30px !important;
  }
}

.table {
  width: 100%;
  max-width: 100%;
  margin-bottom: 20px;
  font-size: 15px;

  tr.table-item {
    position: relative;
    th {
      position: relative;
    }

    .clickable-text {
      cursor: pointer;
      &:hover {
        color: var(--mainColor)
      }
    }

    .field-rank {
      position: absolute;
      font-size: 12px;
      margin-left: 6px;
      background: #F4F4F4;
      padding: 0 6px;
      border-radius: 2px;
      margin-top: 2px;
      animation-duration: 0.5s;
      cursor: pointer;
      user-select: none;

      &:hover {
        background: #eeeeee;
      }
    }

    .field-rank-print {
      font-size: 13px;
    }
  }
}

.table > thead > tr > th,
.table > tbody > tr > th,
.table > tfoot > tr > th,
.table > thead > tr > td,
.table > tbody > tr > td,
.table > tfoot > tr > td {
  padding: 0.9rem 1.2rem;
  vertical-align: top;
  color: #707070;
  border-top: 1px solid #f2f2f2;
  text-align: center;
}

@media screen and (max-width: 559px) {
  .table > thead > tr > th,
  .table > tbody > tr > th,
  .table > tfoot > tr > th,
  .table > thead > tr > td,
  .table > tbody > tr > td,
  .table > tfoot > tr > td {
    padding: 8px 5px;
    min-width: 100px;
  }

  .table > thead:first-child > tr:first-child > th:first-child,
  .table > tbody > tr > td:first-child {
    min-width: 50px;
  }
}

.table > caption + thead > tr:first-child > th,
.table > colgroup + thead > tr:first-child > th,
.table > thead:first-child > tr:first-child > th,
.table > caption + thead > tr:first-child > td,
.table > colgroup + thead > tr:first-child > td,
.table > thead:first-child > tr:first-child > td {
  border-top: 0;
}

.table > caption + thead > tr:first-child > th,
.table > colgroup + thead > tr:first-child > th,
.table > thead:first-child > tr:first-child > th,
.table > caption + thead > tr:first-child > td,
.table > colgroup + thead > tr:first-child > td,
.table > thead:first-child > tr:first-child > td {
  color: #707070;
  background-color: #fff;
  border-bottom: 2px solid #f2f2f2;
}

.table-striped > tbody > tr:nth-of-type(odd) {
  background: #fff;
}

.table-striped > tbody > tr:nth-of-type(even) {
  background: #fcfcfc;
}

.table th,
label {
  font-weight: 500;
}

.wly-table-body {
  overflow-x: auto;
  overflow-y: auto;
  height: 100%;
  transition: filter 0.15s ease-in-out;

  table {
    border: 0;

    thead {
    }

    tbody td {
      position: relative;

      .ranking {
        font-size: 12px;
        vertical-align: text-top;
        color: var(--mainColor);
        position: absolute;
        margin-left: 3px;
      }
    }
  }

  .table-link {
    cursor: pointer;
  }
}

/* pagination */
.paginate-simple {
  display: flex;
  flex-direction: row;

  .paginate-button {
    background-color: #efefef;
    display: inline-block;
    color: #8a8a8a;
    vertical-align: top;
    border-radius: 50%;
    margin: 0 1px 0 2px;
    font-size: 1rem;
    cursor: pointer;
    width: 2.5rem;
    height: 2.5rem;
    line-height: 2.5rem;
    text-align: center;
    user-select: none;
  }

  .paginate-button {
    &.current {
      background-color: var(--mainColor);
      color: #fff;
      cursor: default;
    }

    &.current,
    &.disabled {
      cursor: default;
    }

    &:not(.current):not(.disabled):focus,
    &:not(.current):not(.disabled):hover {
      background-color: #e2e2e2;
      color: #575757;
      text-decoration: none;
    }

    &.next,
    &.previous,
    &.first-page,
    &.last-page {
      font-size: 0;
      position: relative;
    }

    @media screen and (-ms-high-contrast: active), (-ms-high-contrast: none) {
      &.next,
      &.previous,
      &.first-page,
      &.last-page {
        font-size: 1rem;
      }
    }

    &.previous:before,
    &.next:before,
    &.first-page:before,
    &.last-page:before {
      font-family: Material-Design-Iconic-Font;
      font-size: 1rem;
      line-height: 2.55rem;
    }

    &.previous:before {
      content: '\f2fa';
    }

    &.next:before {
      content: '\f2fb';
    }

    &.disabled {
      opacity: 0.6;
    }

    &.disabled:focus,
    &.disabled:hover {
      color: #8a8a8a;
    }
  }
}
</style>
