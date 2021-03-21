<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated style="height: 100px">
      <q-toolbar class="row">
        <q-btn
          v-if="!$route.meta.hideDrawer"
          flat
          dense
          round
          
          side=left
          size=30px
          icon="menu"
          aria-label="Menu"
          class=""
          @click="isDrawer = !isDrawer"
        />

      </q-toolbar>
    </q-header>
    <q-drawer
        v-model="isDrawer"
        
        :width="200"
        :breakpoint="500"
        bordered
        content-class="bg-grey-3"
      >
        <q-scroll-area class="fit" >
          <q-list>
            <template v-for="(menuItem, index) in menuList">
              <q-item
                :to="menuItem.redirect"
                :key="index"
                exact
                clickable
                :active="menuItem.label === 'Outbox'"
                v-ripple>
                <q-item-section avatar>
                  <q-icon :name="menuItem.icon" />
                </q-item-section>
                <q-item-section>
                  {{ menuItem.label }}
                </q-item-section>
              </q-item>
              <q-separator :key="'sep' + index"  v-if="menuItem.separator" />
            </template>

          </q-list>
        </q-scroll-area>
      </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
import { LocalStorage } from 'quasar'
// import EssentialLink from 'components/EssentialLink.vue'


const menuList = [
  {
    icon: 'settings_remote',
    label: 'Activate Beacon',
    redirect: '/beacon',
    separator: true
  },
  {
    icon: 'info',
    label: 'What is this web app about?',
    redirect: '/',
    separator: false
  },
  {
    icon: 'groups',
    label: 'Who are the Developpers',
    redirect: '/devinfo',
    separator: false
  }
]

export default {
  name: 'MainLayout',
  // components: { EssentialLink },
  data () {
    return {
      isDrawer: false,
      menuList
    }
  },
  methods: {
    getStorage() {
      const empty = LocalStorage.isEmpty()
      return empty
      // if (!empty) {
      //   LocalStorage.set('firstTime', '')
      // }
    }
  },
  created () {
    this.firstTime = this.getStorage()
    this.hideDrawer = this.$route.meta.hideDrawer
  }
}
</script>
