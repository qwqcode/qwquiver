<template>
  <transition name="fade">
    <div
      v-if="isShow"
      class="dialog"
      style="animation-duration: 0.2s"
    >
      <div ref="inner" class="inner">
        <div v-if="showCloseBtn" class="close-btn" @click="onCloseBtnClick">
          <i class="zmdi zmdi-close" />
        </div>
        <div :class="{ body: !noSlotStyle }">
          <slot ref="bodyContent" />
        </div>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
import Vue from 'vue'
import { Prop, Component } from 'vue-property-decorator'
import $ from 'jquery'

@Component({})
export default class Dialog extends Vue {
  isShow: boolean = false

  @Prop({
    default: false
  }) readonly outClickClose!: boolean

  @Prop({
    default: true
  }) readonly showCloseBtn!: boolean

  @Prop({
    default: false
  }) readonly noSlotStyle!: boolean

  show () {
    this.isShow = true
    if (this.outClickClose) {
      $(document).bind('click.dialogOutClickHide', (e) => {
        if ($(e.target).is(this.$refs.inner as Element) || !$(e.target).closest(this.$refs.bodyContent as Element)) {
          this.hide()
          $(document).unbind('click.dialogOutClickHide')
        }
      })
    }
  }

  hide () {
    this.isShow = false
    if (this.outClickClose) {
      $(document).unbind('click.dialogOutClickHide')
    }
  }

  onCloseBtnClick () {
    this.isShow = false
    this.$emit('closed')
  }

  beforeDestroy () {
    this.hide()
  }
}
</script>

<style lang="scss" scoped>
.dialog {
  z-index: 99;
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;

  & > .inner {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    overflow-y: auto;
    overflow-x: hidden;
    background: rgba(0, 0, 0, 0.35);

    & > .body {
      padding: 0 10vw;
      display: flex;
      height: 100%;
      align-items: center;
      justify-content: center;
    }
  }

  .close-btn {
    z-index: 9999;
    position: fixed;
    right: 35px;
    top: 70px;

    & > i {
      font-size: 2.3em;
      color: #fff;
      cursor: pointer;
      user-select: none;

      &:hover {
        opacity: 0.7;
      }

      padding: 10px;
    }
  }
}
</style>
