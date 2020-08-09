<template>
  <div class="notify-layer" />
</template>

<script lang="ts">
import { Vue, Component } from 'nuxt-property-decorator'

@Component
export default class NotifyLayer extends Vue {
  created () {
    Vue.prototype.$notify = this
  }

  add (message: string, level?: string, timeout: number = 2000): void {
    message = String(message)
    window.console.log(`[notify][${level}][${new Date().toLocaleString()}] ${message}`)

    const notifyElem = document.createElement('div')
    notifyElem.className = `notify-item anim-fade-in ${(level ? 'type-' + level : '')}`
    notifyElem.innerHTML = `<p class="notify-content">${message.replace('\n', '<br/>')}</p>`
    this.$el.prepend(notifyElem)

    const notifyRemove = function () {
      notifyElem.className += ' anim-fade-out'
      window.setTimeout(function () {
        notifyElem.remove()
      }, 200)
    }

    const timeOutFn = window.setTimeout(() => {
      notifyRemove()
    }, timeout)

    notifyElem.onclick = () => {
      window.clearTimeout(timeOutFn)
      notifyRemove()
    }
  }

  success (message: string) {
    this.add(message, 's')
  }

  error (message: string) {
    this.add(message, 'e')
  }

  warning (message: string) {
    this.add(message, 'w')
  }

  info (message: string) {
    this.add(message, 'i')
  }

  clearAll () {
    this.$el.innerHTML = ''
  }
}
</script>

<style lang="scss" scoped>
.notify-layer {
  $width: 430px;
  position: fixed;
  z-index: 9999;
  top: 70px;
  left: calc(50vw - (#{$width} / 2));
  width: $width;
  pointer-events: none;

  @media screen and (max-width: 600px) {
    & {
      top: 55px;
      left: 0;
      width: 100%;
    }
  }

  /deep/ .notify-item {
    display: block;
    overflow: hidden;
    background-color: #2c2c2c;
    border-color: #2c2c2c;
    color: #FFF;
    cursor: pointer;
    pointer-events: all;
    border-radius: 2px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.075);

    @media screen and (max-width: 600px) {
      border-radius: 0;
    }

    &:not(:last-child) {
      margin-bottom: 13px;
    }

    &.type-s {
      color: #fff;
      background-color: #32c787;
    }

    &.type-e {
      color: #fff;
      background-color: #ff6b68;
    }

    &.type-w {
      color: #fff;
      background-color: #ffc721;
    }

    &.type-i {
      color: #fff;
      background-color: #03a9f4;
    }

    p {
      line-height: 1.8;
      padding: 10px 17px;
      margin: 0;
      text-align: center;
    }
  }
}
</style>
