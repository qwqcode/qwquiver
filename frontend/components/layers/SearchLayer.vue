<template>
  <transition name="fade">
    <div v-if="isShow" class="search-panel" style="animation-duration: 0.14s">
      <div class="inner">
        <div class="tab-bar">
          <span
            v-for="(iLabel, iType) in searchTypeList"
            :key="iType"
            class="type-switch"
            :class="{ 'active': iType === searchType }"
            @click="searchType = iType"
          >{{ iLabel }}</span>
          <span v-if="!!searchExamName" class="curt-exam" @click="$refs.examSelect.show($event.target)">{{ searchExamLabel }}</span>
          <SelectFloater
            ref="examSelect"
            :options="$app.ExamNameToLabelObj || {}"
            :selected="searchExamName"
            :change="switchExam" />
        </div>
        <form :class="`search-type-${searchType}`" class="search-form" @submit.prevent="submit">
          <div v-if="searchType === 'Name'">
            <button type="submit">
              <i class="zmdi zmdi-search"></i>
            </button>
            <input
              ref="SearchInput"
              v-model="searchData.NAME"
              type="text"
              placeholder="搜索..."
              autocomplete="off"
              required
              autofocus
            />
          </div>

          <div v-if="searchType === 'SchoolClass'">
            <LoadingLayer ref="scLoading" />
            <div v-if="!!sc && !!sc.data" class="school-class-list mini-scrollbar">
              <div class="list school-list">
                <span
                  v-for="school in Object.keys(sc.data.school)"
                  :key="school"
                  :class="{ active: school === sc.openedSchool }"
                  class="item"
                  @click="sc.openedSchool = school"
                >{{ school }}</span>
              </div>
              <div class="list class-list">
                <span class="school-name">{{ sc.openedSchool }}</span>
                <span class="item" @click="scSubmit(sc.openedSchool)">全校成绩</span>
                <span
                  v-for="className in scOpenedSchoolClassList"
                  :key="className"
                  class="item"
                  @click="scSubmit(sc.openedSchool, className)"
                >{{ className }}</span>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'nuxt-property-decorator'
import _ from 'lodash'
import $ from 'jquery'
import LoadingLayer from '../LoadingLayer.vue'
import SelectFloater from '../SelectFloater.vue'
import F, { ScoreData } from '../../types/Field'
import * as ApiT from '../../types/ApiTypes'

type SearchType = 'Name' | 'SchoolClass'
const OutClickEvtName = 'click.SearchLayer'

@Component({
  components: { LoadingLayer, SelectFloater }
})
export default class SearchLayer extends Vue {
  created () {
    Vue.prototype.$searchLayer = this
  }

  mounted () {
    this.bindOutClickEvt()
  }

  isShow = false
  searchType: SearchType = 'Name'
  searchTypeList: { [key in SearchType]: string } = {
    Name: '姓名',
    SchoolClass: '学校班级'
  }
  searchData: { [key in F]?: string } = {}
  searchExamName: string|null = null

  scLoading!: LoadingLayer
  sc: {
    data: ApiT.AllSchoolData | null
    openedSchool: string | null
  } | null = null

  get searchExamLabel () {
    if (this.searchExamName === null) return null
    return this.$app.ExamNameToLabelObj ? this.$app.ExamNameToLabelObj[this.searchExamName] || null : null
  }

  switchExam (examName: string) {
    this.searchExamName = examName
  }

  @Watch('searchType')
  onSearchTypeChanged (searchType: SearchType) {
    this.searchData = {}

    if (searchType === 'Name') {
      this.focusSearchInput()
    }

    if (searchType === 'SchoolClass') {
      this.$nextTick(() => {
        this.scLoadSchoolClass()
      })
    }
  }

  @Watch('searchExamName')
  onSearchExamNameChanged () {
    if (this.searchType === 'SchoolClass') {
      this.$nextTick(() => {
        this.scLoadSchoolClass()
      })
    }
  }

  focusSearchInput () {
    this.$nextTick(() => {
      const input = (this.$refs.SearchInput as HTMLInputElement)
      if (!input || !input.focus) return
      input.focus()
    })
  }

  submit () {
    const reqParams: ApiT.QueryParams = { where: JSON.stringify(this.searchData), page: 1 }
    if (this.searchExamName) reqParams.exam = this.searchExamName
    this.$explorer.fetchData(reqParams)
    this.$nextTick(() => {
      this.searchData = {}
      this.hide()
    })
  }

  scSubmit (schoolName: string, className?: string) {
    this.searchData.SCHOOL = schoolName
    if (className)
      this.searchData.CLASS = className
    this.submit()
  }

  get scOpenedSchoolClassList () {
    const sc = this.sc
    if (sc === null || !sc.data || !sc.openedSchool || !sc.data.school[sc.openedSchool])
      return []
    return _.orderBy(sc.data.school[sc.openedSchool], (o) => {
      const num = o.match(/[0-9]+/)
      return num ? Number(num[0] || 0) : 0
    }, ['asc'])
  }


  async scLoadSchoolClass () {
    if (this.searchType !== 'SchoolClass') return
    this.sc = null
    this.scLoading = this.$refs.scLoading as LoadingLayer
    this.scLoading.show()
    let respData
    try {
      respData = await this.$axios.$get('./api/school/all', {
        params: { exam: this.searchExamName }
      })
    } catch (err) {
      this.$notify.error(String(err))
    } finally {
      this.scLoading.hide()
    }
    if (!respData.success) {
      this.$notify.error(respData.msg)
      return
    }
    const data = (respData.data || null) as ApiT.AllSchoolData
    this.sc = { data, openedSchool: Object.keys(data.school)[0] || null }
  }

  show () {
    this.isShow = true
    this.bindOutClickEvt()
    this.focusSearchInput()
    this.searchExamName = this.$explorer.curtExamName
  }

  hide () {
    this.isShow = false
    this.unbindOutClickEvt()
  }

  toggle () {
    this.isShow ? this.hide() : this.show()
  }

  bindOutClickEvt () {
    $(window).unbind(OutClickEvtName)
    window.setTimeout(() => {
      $(window).bind(OutClickEvtName, (evt) => {
        if (this.isShow && !$(evt.target).closest('.search-panel').length) {
          this.hide()
        }
      })
    }, 80)
  }

  unbindOutClickEvt () {
    $(window).unbind(OutClickEvtName)
  }
}
</script>

<style scoped lang="scss">
%card {
  display: flex;
  background: #fff;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.2);
  border-radius: 3px;
  overflow: hidden;
  position: relative;
}

/* wly-search-panel */
.search-panel {
  position: fixed;
  justify-content: center;
  display: flex;
  top: 0;
  left: 0;
  width: 100%;
  margin-top: calc(30vh - 20px);
  z-index: 999;
  pointer-events: none;
}

.inner {
  width: 580px;
  position: relative;
}

.tab-bar {
  pointer-events: none;
  padding-left: 12px;
  display: flex;
  flex-direction: row;

  & > .type-switch {
    @extend %card;
    background: rgba(255, 255, 255, 0.8);
    pointer-events: all;
    border-radius: 2px 3px 0 0;
    padding: 5px 20px;
    margin-left: -3px;
    cursor: pointer;
    transition: background-color .2s,border .2s,box-shadow .2s, color .2s;

    &.active {
      z-index: 2;
      background: #fff;
      color: var(--mainColor);
    }

    &:hover {
      background: #fff;
    }
  }

  & > .curt-exam {
    position: relative;
    pointer-events: all;
    padding-left: 10px;
    padding-right: 20px;
    margin-left: auto;
    margin-right: 10px;
    margin-top: 6px;
    background: rgba(255, 255, 255, 0.8);
    height: 24px;
    line-height: 24px;
    border-radius: 2px;
    font-size: 12px;
    cursor: pointer;
    box-shadow: 0 1px 10px rgba(0, 0, 0, 0.2);
    user-select: none;
    transition: background-color .2s,border .2s,box-shadow .2s;

    &:hover {
      background: #fff
    }

    &:after {
      position: absolute;
      right: 0;
      top: 0;
      height: 100%;
      padding-right: 8px;
      content: '\f2f2';
      font-family: 'Material-Design-Iconic-Font';
    }
  }
}

.search-form {
  @extend %card;
  pointer-events: all;
  z-index: 4;
  height: fit-content;
  width: 100%;
  transition: height 0.2s ease-in-out;

  &.search-type-Name {
    height: 60px;

    & > div {
      display: flex;
      flex: 1;
      flex-direction: row;
      place-items: center;

      input {
        flex: 1;
        height: 100%;
        padding-left: 55px;
        padding-right: 30px;
        outline: none;
        border: none;

        &:focus {
        }
      }

      button {
        cursor: pointer;
        background: transparent;
        border: 0;
        box-shadow: none;
        outline: none;
        left: 10px;
        position: absolute;
        width: 40px;
        height: 40px;
        font-size: 20px;
      }
    }
  }

  &.search-type-SchoolClass {
    height: 300px;

    & > div {
      display: flex;
      width: 100%;
    }

    & > div > .school-class-list {
      display: flex;
      width: 100%;

      .list {
        font-size: 13px;
        height: 300px;
        overflow-y: auto;
        overflow-x: hidden;

        &.school-list {
          flex: 30%;
          box-shadow: 0 -8px 6px rgba(0, 0, 0, 0.2);

          .item {
            position: relative;
            display: block;
            padding: 10px 21px;
            border-left: 1px solid transparent;
            cursor: pointer;
            transition: background-color .2s,border .2s,box-shadow .2s;

            &.active {
              color: var(--mainColor);
            }

            &.active:before, &:hover:before {
              background: var(--mainColor);
              content: ' ';
              position: absolute;
              left: -2px;
              top: 10px;
              height: calc(100% - 20px);
              width: 3px;
              box-shadow: 0 2px 15px rgba(0, 131, 255, 0.22);
              border-left: 1px solid var(--mainColor);
            }
          }
        }

        &.class-list {
          flex: 70%;
          padding: 8px 13px 8px 23px;

          .school-name {
            display: block;
            font-size: 14px;
            padding: 8px 8px 11px 8px;
            border-bottom: 1px solid #F4F4F4;
          }

          .item {
            background: #f4f4f4;
            border-radius: 40px;
            display: inline-block;
            padding: 5px 13px;
            margin: 14px 10px 0 0;
            cursor: pointer;
            transition: background-color .2s,border .2s,box-shadow .2s;

            &:hover {
              background: #EEE;
            }
          }
        }
      }
    }
  }
}
</style>
