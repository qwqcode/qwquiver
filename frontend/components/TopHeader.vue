<template>
  <div class="top-header">
    <div class="main-navbar" :class="{ 'hide': !isShow }">
      <div class="left">
        <div class="sidebar-toggle-btn" @click="$sidebar.toggle()">
          <i class="zmdi zmdi-menu"></i>
        </div>
        <span class="brand">
          <a href="/" class="app-name">市 2018/2021届 初高中 成绩统计</a>
        </span>
      </div>
      <div class="right">
        <ul class="actions-btn-bar main-bar">
          <li>
            <a @click="$searchLayer.toggle()">
              <i class="zmdi zmdi-search"></i>
              <span>搜索</span>
            </a>
          </li>
        </ul>
        <ul class="actions-btn-bar scene-bar" style="display: none;"></ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'nuxt-property-decorator'

@Component({})
export default class TopHeader extends Vue {
  isShow = true

  created () {
    Vue.prototype.$topHeader = this
  }

  show () {
    this.isShow = true
  }

  hide () {
    this.isShow = false
  }

  toggle () {
    if (this.isShow) {
      this.hide()
    } else {
      this.show()
    }
  }
}
</script>

<style scoped lang="scss">
.top-header {
  z-index: 9998;
  position: fixed;
  left: 0;
  top: 0;
  width: 100%;
}

.main-navbar {
  height: 55px;
  overflow: hidden;
  position: fixed;
  width: 100%;
  background: var(--mainColor);
  color: #FFFFFF;
  padding: 0 15px;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  transition: background 0.15s ease-in-out, transform 450ms cubic-bezier(0.23, 1, 0.32, 1) 0ms;

  &.hide {
    transform: translate(0, -70px);
  }

  @media screen and (max-width: 559px) {
    & { padding: 0 5px 0 10px; }
  }

  & > .left,
  & > .right {
    position: relative;
    overflow: hidden;
    /* height: 60px; */
  }

  & > .left {
    float: left;
    display: flex;
    align-items: center;

    .sidebar-toggle-btn {
      cursor: pointer;
      color: #FFF;
      position: relative;
      transition: background-color .3s;
      font-size: 15px;
      border-radius: 2px;
      margin-right: 15px;

      &:active {
        background-color: rgba(255,255,255,.2);
      }

      & > i {
        padding: 8px 8px;
      }
    }
  }

  & > .right {
    float: right;
    -webkit-box-flex: 1;
    -ms-flex: 1;
    flex: 1;

    .actions-btn-bar {
      list-style: none;
      padding: 0;
      margin: 0 0 0 auto;
      float: right;

      & > li {
        display: inline-block;
        vertical-align: middle;
        position: relative;
      }

      & > li > a {
        cursor: pointer;
        padding: 8px 14px;
        min-width: 50px;
        display: block;
        color: #FFF;
        border-radius: 2px;
        text-align: center;
        line-height: 100%;
        position: relative;
        transition: background-color .3s;
        font-size: 14px;
        vertical-align: middle;

        &:hover {
          background-color: rgba(255,255,255,.2)
        }


        @media screen and (max-width: 559px) {
          & {
              padding: 8px 10px;
              min-width: inherit;
          }

          & > span {
              display: none;
          }
        }

        & > i {
          font-size: 18px;
          vertical-align: middle;
          margin-right: 2px;
        }

      }
    }
  }

  .brand {
    font-size: 17px;
    font-weight: normal;
    margin: 0;
    color: rgba(255, 255, 255, 0.94);
  }

  .app-name {
    color: #fff;
    text-decoration: none;
    outline: none;
    transition: transform .3s,-webkit-transform .3s;

    &.changing {
      -webkit-transform: translateY(-100%);
      transform: translateY(-100%);
    }

    i {
      vertical-align: middle;
      width: 30px;
      line-height: 25px;
      font-size: 1.6rem;
    }
  }
}
</style>
