<template>
  <div></div>
</template>

<script lang="ts">
import { Component, Vue } from 'nuxt-property-decorator'

@Component({})
export default class QuickStats extends Vue {}
</script>

<style scoped lang="scss">
/* quick-stats */
.quick-stats {
}

@media screen and (max-width: 559px) {
  .quick-stats {
    display: none;
  }
}

.quick-stats-item {
  padding: 1.5rem 1.5rem 1.45rem;
  border-radius: 2px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-align: baseline;
  -ms-flex-align: baseline;
  align-items: baseline;
}

.quick-stats-item .quick-stats-info > h2 {
  font-weight: normal;
  margin: 0;
  font-size: 1.6rem;
  color: #ffffff;
  text-shadow: 2px 2px 2px rgba(255, 255, 255, 0.34);
}

.quick-stats-item .quick-stats-info > small {
  font-size: 1rem;
  display: block;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 0.6rem;
  line-height: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-stats-item-blue {
  background-color: #03a9f4;
}

.quick-stats-item-amber {
  background-color: #f5c942;
}

.quick-stats-item-purple {
  background-color: #d066e2;
}

.quick-stats-item-red {
  background-color: #ff6b68;
}
</style>
