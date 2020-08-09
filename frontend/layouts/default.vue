<template>
  <div>
    <!-- eslint-disable-next-line vue/attributes-order -->
    <component v-for="layer in AllLayersNameList" :is="layer" :key="layer" />
    <TopHeader />

    <div class="wrap">
      <Sidebar />

      <div :class="{ 'full': contFullScreen }" class="main-cont-area">
        <Explorer v-show="$route.name === 'index'" />
        <nuxt class="content-inner" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'nuxt-property-decorator'
import $ from 'jquery'
import _ from 'lodash'
import F from '../types/Field'
import * as ApiT from '../types/ApiTypes'
import TopHeader from '@/components/TopHeader.vue'
import Sidebar from '@/components/Sidebar.vue'
import Explorer from '@/components/Explorer.vue'
import Layers from '@/components/layers'

@Component({
  components: { TopHeader, Sidebar, Explorer, ...Layers }
})
export default class Default extends Vue {
  Conf: ApiT.ConfData|null = null
  contFullScreen = false

  created () {
    Vue.prototype.$app = this
  }

  mounted () {
    // 载入最新的考试数据
    let params: ApiT.QueryParams = {
      page: 1,
      pageSize: 60,
      init: true
    }
    if (this.$route.query) params = { ...params, ...this.$route.query }
    this.$explorer.onRouteQueryChanged(params as ApiT.QueryParams)
  }

  get ExamMap () {
    if (this.Conf === null) return null
    return this.Conf.examMap
  }

  get ExamMapSorted () {
    if (this.ExamMap === null) return null
    const arr = _.sortBy(this.$app.ExamMap, o => o.Date ? -(new Date(o.Date).getTime()) : -1)
    return arr
  }

  get ExamNameToLabelObj () {
    if (this.ExamMap === null) return null
    const obj = {}
    _.forEach(this.ExamMapSorted, (e) => { obj[e.Name] = e.Label })
    return obj
  }

  get FieldTransDict () {
    if (this.Conf === null) return null
    return this.Conf.fieldTransDict
  }

  transField (f: F) {
    if (!this.FieldTransDict) return f
    return this.FieldTransDict[f] || f
  }

  get AllLayersNameList () {
    return Object.keys(Layers)
  }

  getContentHeight () {
    return (
      ($(window).height() || 0) -
      ($('.main-navbar').outerHeight(true) || 0) -
      ($('.card .card-header').outerHeight(true) || 0) -
      80
    )
  }

  setContFullScreen (val: boolean) {
    this.contFullScreen = val
  }
}
</script>

<style scoped lang="scss"></style>
