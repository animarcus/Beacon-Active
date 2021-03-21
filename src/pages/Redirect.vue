<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
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
              <q-item :key="index" clickable :active="menuItem.label === 'Outbox'" v-ripple>
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


export default {
  name: 'MainLayout',
  // components: { EssentialLink },
  data () {
    return {
      firstTime: this.getStorage()
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
    this.getStorage()
  }
}
</script>
