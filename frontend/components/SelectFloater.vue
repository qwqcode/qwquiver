<template>
  <transition name="fade">
    <div v-show="isShow" class="select-floater mini-scrollbar" style="animation-duration: 0.1s">
      <div
        v-for="(val, key) in options"
        :key="key"
        class="option"
        :class="{ 'selected': !!selected && selected === key }"
        @click="setVal(key)">{{ val }}</div>
    </div>
  </transition>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator"
import $ from 'jquery'

const OutClickEvtName = 'click.SelectFloater'

@Component({})
export default class SelectFloater extends Vue {
  @Prop({
    required: true
  })
  options!: { [key: string]: string }

  @Prop({})
  selected?: string

  @Prop({})
  change?: (val: string) => void

  isShow = false

  show (btnEl?: HTMLButtonElement) {
    this.isShow = true
    if (btnEl) {
      (this.$el as HTMLElement).style.left = `${btnEl.offsetLeft}px`;
      (this.$el as HTMLElement).style.top = `${btnEl.offsetTop}px`
    }
    this.bindOutClickEvt()
  }

  hide () {
    this.isShow = false
    this.unbindOutClickEvt()
  }

  private setVal (val: string) {
    if (this.change) {
      this.change(val)
    }
    this.hide()
  }

  bindOutClickEvt () {
    $(window).unbind(OutClickEvtName)
    window.setTimeout(() => {
      $(window).bind(OutClickEvtName, (evt) => {
        if (this.isShow && !$(evt.target).closest('.select-floater').length) {
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
.select-floater {
  position: absolute;
  z-index: 99999;
  pointer-events: all;
  background: #fff;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.2);
  max-height: 325px;
  width: 145px;
  border-radius: 2px;
  overflow-y: auto;

  .option {
    padding: 7px 10px;
    text-align: center;
    font-size: 13px;
    cursor: pointer;
    user-select: none;

    &:hover, &.selected {
      background: #f4f4f4;
    }
  }
}
</style>
