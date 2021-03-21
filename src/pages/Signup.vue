<template>
  <q-page class="flex flex-center column text-center">
    <h4>Welcome!</h4>
    <h4 class="displayText">Before proceeding to the Beacon Active web app, please input a username to identify yourself (you can change it at anytime)</h4>
    <q-input
      filled
      v-model="username"
      @keyup.enter="authenticate"
      label="Username" />
  </q-page>
</template>





<script>
import { LocalStorage } from 'quasar'
import { api } from 'boot/axios'
export default {
  data() {
    return {
      username: ''
    }
  },
  methods: {
    async authenticate() {
      // console.log("aaaaaaaaaaaaaaaa")
      // console.log(MainLayoutVue.data)
      if (!this.username.length) return
      const user = await checkUsername(this.username)
      if (!user) return
      LocalStorage.set('userInfo', user)
      // console.log("request response:", t, this.username)
      this.username = ''

      this.$router.push({path: '/'})
    }
  }

}






async function checkUsername(username) {
  return 42
  // const data = { username }
  // api.post('http://philips-macbook.local:8080/', username)
  //     .then(response => {
  //       console.log(response.data.user_id)
  //       return response.data.user_id
  //     });  //userId    
}
</script>


<style lang="scss">
  .displayText {
    font-size: 30px;
    padding: 0px
  }
</style>


